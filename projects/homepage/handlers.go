package main

import (
	"context"
	"crypto/sha256"
	"errors"
	"log"
	"net/http"
	netmail "net/mail"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"mljr-web/internal/i18n"
	"mljr-web/internal/mail"
	"mljr-web/internal/web"
	"mljr-web/projects/homepage/homelab"
	"mljr-web/projects/homepage/pages"
	"mljr-web/ui/layout"
	"mljr-web/ui/primitive"
	"mljr-web/ui/token"

	altchalib "github.com/altcha-org/altcha-lib-go"
	"github.com/labstack/echo/v4"
	"github.com/starfederation/datastar-go/datastar"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

// registerHandlers wires all API endpoints onto the Echo instance.
func registerHandlers(e *echo.Echo, altchaKey string, contactMailer mail.ContactMailer) {
	e.GET("/api/altcha", altchaChallenge(altchaKey))
	e.POST("/api/contact", contactSubmit(altchaKey, contactMailer))
}

// ── /api/altcha ───────────────────────────────────────────────────────────────

func altchaChallenge(key string) echo.HandlerFunc {
	return func(c echo.Context) error {
		expires := time.Now().Add(10 * time.Minute)
		ch, err := altchalib.CreateChallenge(altchalib.ChallengeOptions{
			HMACKey:   key,
			MaxNumber: 200000,
			Expires:   &expires,
		})
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "challenge creation failed")
		}
		return c.JSON(http.StatusOK, ch)
	}
}

// ── /api/contact ──────────────────────────────────────────────────────────────

type contactSignals struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Message  string `json:"message"`
	Altcha   string `json:"altcha"`
	Honeypot string `json:"_hp"`
}

const (
	contactMaxBodyBytes     = 32 << 10
	contactMaxNameLength    = 100
	contactMaxEmailBytes    = 254
	contactMaxMessageLength = 5000
	contactRateWindow       = 10 * time.Minute
	contactRateMax          = 5
	contactTrackingMax      = 4096
	contactReplayWindow     = 15 * time.Minute
)

type contactRateEntry struct {
	start time.Time
	count int
}

type contactRateLimiter struct {
	mu            sync.Mutex
	byIP          map[string]contactRateEntry
	usedSolutions map[[sha256.Size]byte]time.Time
}

func (l *contactRateLimiter) allow(ip string, now time.Time) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.byIP == nil {
		l.byIP = make(map[string]contactRateEntry)
	}
	if _, exists := l.byIP[ip]; !exists && len(l.byIP) >= contactTrackingMax {
		l.dropOldestIP()
	}
	entry := l.byIP[ip]
	if entry.start.IsZero() || now.Sub(entry.start) >= contactRateWindow {
		l.byIP[ip] = contactRateEntry{start: now, count: 1}
		return true
	}
	if entry.count >= contactRateMax {
		return false
	}
	entry.count++
	l.byIP[ip] = entry
	return true
}

func (l *contactRateLimiter) dropOldestIP() {
	var oldestIP string
	var oldest time.Time
	for ip, entry := range l.byIP {
		if oldest.IsZero() || entry.start.Before(oldest) {
			oldestIP, oldest = ip, entry.start
		}
	}
	delete(l.byIP, oldestIP)
}

func (l *contactRateLimiter) consumeSolution(solution string, now time.Time) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.usedSolutions == nil {
		l.usedSolutions = make(map[[sha256.Size]byte]time.Time)
	}
	for hash, usedAt := range l.usedSolutions {
		if now.Sub(usedAt) >= contactReplayWindow {
			delete(l.usedSolutions, hash)
		}
	}
	hash := sha256.Sum256([]byte(solution))
	if _, used := l.usedSolutions[hash]; used {
		return false
	}
	if len(l.usedSolutions) >= contactTrackingMax {
		var oldestHash [sha256.Size]byte
		var oldest time.Time
		for candidate, usedAt := range l.usedSolutions {
			if oldest.IsZero() || usedAt.Before(oldest) {
				oldestHash, oldest = candidate, usedAt
			}
		}
		delete(l.usedSolutions, oldestHash)
	}
	l.usedSolutions[hash] = now
	return true
}

func contactSubmit(key string, contactMailer mail.ContactMailer) echo.HandlerFunc {
	if contactMailer == nil {
		contactMailer = mail.LogMailer{}
	}
	limiter := &contactRateLimiter{}
	return func(c echo.Context) error {
		lang := web.Lang(c)
		if !limiter.allow(c.RealIP(), time.Now()) {
			sse := datastar.NewSSE(c.Response().Writer, c.Request())
			return patchFormError(sse, i18n.T(lang, "contact.rate_limited"))
		}

		c.Request().Body = http.MaxBytesReader(c.Response().Writer, c.Request().Body, contactMaxBodyBytes)

		var s contactSignals
		if err := datastar.ReadSignals(c.Request(), &s); err != nil {
			var tooLarge *http.MaxBytesError
			if errors.As(err, &tooLarge) {
				return echo.NewHTTPError(http.StatusRequestEntityTooLarge, "contact payload too large")
			}
			return echo.NewHTTPError(http.StatusBadRequest, "invalid contact payload")
		}

		sse := datastar.NewSSE(c.Response().Writer, c.Request())

		// honeypot
		if s.Honeypot != "" {
			log.Printf("contact: honeypot triggered from %s", c.RealIP())
			return patchSuccess(sse, lang)
		}

		// validation
		valid := true
		errs := map[string]any{"nameError": "", "emailError": "", "msgError": ""}
		name := strings.TrimSpace(s.Name)
		email := strings.TrimSpace(s.Email)
		message := strings.TrimSpace(s.Message)
		if name == "" {
			errs["nameError"] = i18n.T(lang, "contact.error_name_required")
			valid = false
		} else if utf8.RuneCountInString(name) > contactMaxNameLength || !utf8.ValidString(name) || hasHeaderControl(name) {
			errs["nameError"] = i18n.T(lang, "contact.error_name_invalid")
			valid = false
		}
		if !validContactEmail(email) {
			errs["emailError"] = i18n.T(lang, "contact.error_email_invalid")
			valid = false
		}
		if len(message) < 10 {
			errs["msgError"] = i18n.T(lang, "contact.error_message_short")
			valid = false
		} else if utf8.RuneCountInString(message) > contactMaxMessageLength || !utf8.ValidString(message) || strings.ContainsRune(message, '\x00') {
			errs["msgError"] = i18n.T(lang, "contact.error_message_long")
			valid = false
		}
		if !valid {
			errs["sending"] = false
			return sse.MarshalAndPatchSignals(errs)
		}
		// clear errors
		if err := sse.MarshalAndPatchSignals(errs); err != nil {
			return err
		}

		// altcha
		if s.Altcha == "" {
			return patchFormError(sse, i18n.T(lang, "contact.captcha_required"))
		}
		ok, err := altchalib.VerifySolution(s.Altcha, key, true)
		if err != nil || !ok {
			return patchFormError(sse, i18n.T(lang, "contact.captcha_failed"))
		}
		if !limiter.consumeSolution(s.Altcha, time.Now()) {
			return patchFormError(sse, i18n.T(lang, "contact.captcha_failed"))
		}

		ctx, cancel := context.WithTimeout(c.Request().Context(), 20*time.Second)
		defer cancel()
		if err := contactMailer.SendContact(ctx, mail.ContactMessage{
			Name:      name,
			Email:     email,
			Message:   message,
			RemoteIP:  c.RealIP(),
			UserAgent: c.Request().UserAgent(),
		}); err != nil {
			log.Printf("contact: send mail failed from %q <%q>: %v", name, email, err)
			return patchFormError(sse, i18n.T(lang, "contact.delivery_failed"))
		}

		log.Printf("contact: mail accepted from %q <%q>: %.80q", name, email, message)
		return patchSuccess(sse, lang)
	}
}

func hasHeaderControl(value string) bool {
	return strings.IndexFunc(value, func(r rune) bool { return r < 0x20 || r == 0x7f }) >= 0
}

func validContactEmail(value string) bool {
	if value == "" || len(value) > contactMaxEmailBytes || !utf8.ValidString(value) || hasHeaderControl(value) {
		return false
	}
	addr, err := netmail.ParseAddress(value)
	return err == nil && addr.Address == value
}

func patchSuccess(sse *datastar.ServerSentEventGenerator, lang string) error {
	node := h.Div(
		h.ID("contact-form"),
		g.Attr("data-component", "contact-result"),
		g.Attr("data-variant", "success"),
		h.Div(g.Attr("data-slot", "icon"), g.Text("✓")),
		primitive.Heading(primitive.HeadingProps{Level: 3}, g.Text(i18n.T(lang, "contact.success_title"))),
		h.P(g.Text(i18n.T(lang, "contact.success_body"))),
	)
	if err := sse.PatchElements(web.RenderToString(node)); err != nil {
		return err
	}
	return sse.MarshalAndPatchSignals(map[string]any{"sending": false})
}

func patchFormError(sse *datastar.ServerSentEventGenerator, msg string) error {
	node := layout.Stack(layout.StackProps{Attrs: []g.Node{
		h.ID("contact-error"),
	}},
		primitive.Tag(primitive.TagProps{Tone: token.ToneNone}, g.Text("⚠ "+msg)),
	)
	if err := sse.PatchElements(web.RenderToString(node)); err != nil {
		return err
	}
	return sse.MarshalAndPatchSignals(map[string]any{"sending": false})
}

// ── /api/homelab ──────────────────────────────────────────────────────────────

// registerHomelabHandler serves the live homelab panel fragment. The homepage
// re-fetches it every 60s via data-on-interval and patches #homelab-panel.
func registerHomelabHandler(e *echo.Echo, snapshot func() homelab.Snapshot) {
	e.GET("/api/homelab", func(c echo.Context) error {
		sse := datastar.NewSSE(c.Response().Writer, c.Request())
		return sse.PatchElements(web.RenderToString(pages.HomelabPanel(snapshot())))
	})
}
