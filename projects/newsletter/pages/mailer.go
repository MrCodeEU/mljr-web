package pages

import (
	"net/mail"

	"mljr-web/internal/config"

	pbmailer "github.com/pocketbase/pocketbase/tools/mailer"
)

// mailSender is the minimal interface scheduler.Mailer satisfies — declared
// locally so pages doesn't import scheduler (which would import pages back,
// via no path today, but keeps the dependency direction one-way regardless).
type mailSender interface {
	Send(msg *pbmailer.Message) error
}

var (
	activeMailer mailSender
	mailFrom     mail.Address
	publicAppURL string
)

// SetMailer wires the scheduler's bootstrapped Mailer into pages, so handlers
// like HandleCreateInvite can send mail without pages depending on scheduler.
func SetMailer(m mailSender, cfg config.Config) {
	activeMailer = m
	if addr, err := mail.ParseAddress(cfg.SMTP.From); err == nil {
		mailFrom = *addr
	} else {
		mailFrom = mail.Address{Address: cfg.SMTP.From}
	}
	publicAppURL = cfg.Newsletter.PublicAppURL
}

func sendMail(msg *pbmailer.Message) error {
	if activeMailer == nil {
		return nil
	}
	return activeMailer.Send(msg)
}
