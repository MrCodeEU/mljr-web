package scheduler

import (
	"log"
	"net/mail"
	"strconv"

	"mljr-web/internal/config"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/mailer"
)

// Mailer is the minimal interface the scheduler and invite handlers send
// through, so dev environments (no SMTP_HOST configured) can fall back to
// logging instead of trying — and failing — to actually deliver mail.
type Mailer interface {
	Send(msg *mailer.Message) error
}

type logMailer struct{}

func (logMailer) Send(msg *mailer.Message) error {
	to := msg.To
	if len(to) == 0 {
		to = msg.Bcc
	}
	log.Printf("newsletter mail (SMTP not configured, logging only): to=%v subject=%q", to, msg.Subject)
	return nil
}

type pbMailer struct {
	app core.App
}

func (m pbMailer) Send(msg *mailer.Message) error {
	return m.app.NewMailClient().Send(msg)
}

// debugLoggingMailer wraps another Mailer and logs to=/subject= before
// forwarding the send, so the e2e suite can assert on subject text (e.g.
// language) regardless of whether real SMTP or the log-only fallback is
// configured. Only wired in when NEWSLETTER_E2E_DEBUG=1 (main.go) — never in
// a normal deploy.
type debugLoggingMailer struct {
	inner Mailer
}

func (m debugLoggingMailer) Send(msg *mailer.Message) error {
	to := msg.To
	if len(to) == 0 {
		to = msg.Bcc
	}
	log.Printf("newsletter mail (e2e debug): to=%v subject=%q", to, msg.Subject)
	return m.inner.Send(msg)
}

// WrapDebugMailer wraps m so every send is logged before being forwarded.
func WrapDebugMailer(m Mailer) Mailer {
	return debugLoggingMailer{inner: m}
}

// BootstrapMailer writes cfg.SMTP into PocketBase's app Settings (so
// app.NewMailClient() picks it up) and returns the Mailer to send through.
// If SMTP_HOST is unset, it leaves SMTP disabled in Settings and returns a
// log-only Mailer instead, matching every other projects/* server's dev
// fallback behavior.
func BootstrapMailer(app core.App, cfg config.SMTPConfig) (Mailer, error) {
	if cfg.Host == "" {
		return logMailer{}, nil
	}

	port, err := strconv.Atoi(cfg.Port)
	if err != nil {
		port = 587
	}

	settings := app.Settings()
	settings.SMTP.Enabled = true
	settings.SMTP.Host = cfg.Host
	settings.SMTP.Port = port
	settings.SMTP.Username = cfg.User
	settings.SMTP.Password = cfg.Pass
	// Port 465 is implicit TLS (mailyak.NewWithTLS); every other port (587,
	// 25, ...) uses mailyak.New, which negotiates STARTTLS with the server
	// itself. Forcing TLS=true unconditionally broke STARTTLS-only port 587
	// servers with "tls: first record does not look like a TLS handshake".
	settings.SMTP.TLS = port == 465
	if err := app.Save(settings); err != nil {
		return nil, err
	}

	return pbMailer{app: app}, nil
}

// fromAddress parses cfg.From (e.g. "Newsletter <noreply@mljr.eu>") into a
// mail.Address, falling back to a bare address if it doesn't parse as a
// display-name form.
func fromAddress(from string) mail.Address {
	if addr, err := mail.ParseAddress(from); err == nil {
		return *addr
	}
	return mail.Address{Address: from}
}
