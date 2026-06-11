package main

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

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
		ch, err := altchalib.CreateChallenge(altchalib.ChallengeOptions{
			HMACKey:   key,
			MaxNumber: 200000,
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

func contactSubmit(key string, contactMailer mail.ContactMailer) echo.HandlerFunc {
	if contactMailer == nil {
		contactMailer = mail.LogMailer{}
	}
	return func(c echo.Context) error {
		var s contactSignals
		if err := datastar.ReadSignals(c.Request(), &s); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		sse := datastar.NewSSE(c.Response().Writer, c.Request())

		// honeypot
		if s.Honeypot != "" {
			log.Printf("contact: honeypot triggered from %s", c.RealIP())
			return patchSuccess(sse)
		}

		// validation
		valid := true
		errs := map[string]any{"nameError": "", "emailError": "", "msgError": ""}
		if strings.TrimSpace(s.Name) == "" {
			errs["nameError"] = "Name is required"
			valid = false
		}
		if !strings.Contains(s.Email, "@") || !strings.Contains(s.Email, ".") {
			errs["emailError"] = "Valid email required"
			valid = false
		}
		if len(strings.TrimSpace(s.Message)) < 10 {
			errs["msgError"] = "Message must be at least 10 characters"
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
			return patchFormError(sse, "Please complete the captcha.")
		}
		ok, err := altchalib.VerifySolution(s.Altcha, key, true)
		if err != nil || !ok {
			return patchFormError(sse, "Captcha verification failed. Please try again.")
		}

		ctx, cancel := context.WithTimeout(c.Request().Context(), 20*time.Second)
		defer cancel()
		if err := contactMailer.SendContact(ctx, mail.ContactMessage{
			Name:      strings.TrimSpace(s.Name),
			Email:     strings.TrimSpace(s.Email),
			Message:   strings.TrimSpace(s.Message),
			RemoteIP:  c.RealIP(),
			UserAgent: c.Request().UserAgent(),
		}); err != nil {
			log.Printf("contact: send mail failed from %s <%s>: %v", s.Name, s.Email, err)
			return patchFormError(sse, "Message delivery failed. Please email me directly.")
		}

		log.Printf("contact: mail accepted from %s <%s>: %.80s", s.Name, s.Email, s.Message)
		return patchSuccess(sse)
	}
}

func patchSuccess(sse *datastar.ServerSentEventGenerator) error {
	node := h.Div(
		h.ID("contact-form"),
		g.Attr("data-component", "contact-result"),
		g.Attr("data-variant", "success"),
		h.Div(g.Attr("data-slot", "icon"), g.Text("✓")),
		primitive.Heading(primitive.HeadingProps{Level: 3}, g.Text("Message sent!")),
		h.P(g.Text("I'll get back to you as soon as possible.")),
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
