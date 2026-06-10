package mail

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"mime"
	"net"
	"net/mail"
	"net/smtp"
	"net/textproto"
	"strings"
	"time"
)

type SMTPConfig struct {
	Host string
	Port string
	User string
	Pass string
	From string
	To   string
}

type SMTPMailer struct {
	cfg SMTPConfig
}

func NewSMTPMailer(cfg SMTPConfig) (*SMTPMailer, error) {
	if strings.TrimSpace(cfg.Host) == "" {
		return nil, errors.New("SMTP_HOST is required")
	}
	if strings.TrimSpace(cfg.Port) == "" {
		return nil, errors.New("SMTP_PORT is required")
	}
	if _, err := mail.ParseAddress(cfg.From); err != nil {
		return nil, fmt.Errorf("SMTP_FROM is invalid: %w", err)
	}
	if _, err := mail.ParseAddress(cfg.To); err != nil {
		return nil, fmt.Errorf("CONTACT_TO is invalid: %w", err)
	}
	return &SMTPMailer{cfg: cfg}, nil
}

func (m *SMTPMailer) SendContact(ctx context.Context, msg ContactMessage) error {
	if ctx == nil {
		ctx = context.Background()
	}
	from, err := mail.ParseAddress(m.cfg.From)
	if err != nil {
		return err
	}
	to, err := mail.ParseAddress(m.cfg.To)
	if err != nil {
		return err
	}
	replyTo, err := mail.ParseAddress(msg.Email)
	if err != nil {
		return fmt.Errorf("reply-to email is invalid: %w", err)
	}

	body := contactBody(msg)
	raw := buildMessage(message{
		From:    from.String(),
		To:      to.String(),
		ReplyTo: replyTo.String(),
		Subject: "New contact form message from " + msg.Name,
		Body:    body,
	})

	addr := net.JoinHostPort(m.cfg.Host, m.cfg.Port)
	auth := smtp.Auth(nil)
	if m.cfg.User != "" || m.cfg.Pass != "" {
		auth = smtp.PlainAuth("", m.cfg.User, m.cfg.Pass, m.cfg.Host)
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- smtp.SendMail(addr, auth, from.Address, []string{to.Address}, raw)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errCh:
		return err
	case <-time.After(15 * time.Second):
		return errors.New("SMTP send timed out")
	}
}

type message struct {
	From    string
	To      string
	ReplyTo string
	Subject string
	Body    string
}

func buildMessage(msg message) []byte {
	headers := textproto.MIMEHeader{}
	headers.Set("From", msg.From)
	headers.Set("To", msg.To)
	headers.Set("Reply-To", msg.ReplyTo)
	headers.Set("Subject", mime.QEncoding.Encode("utf-8", msg.Subject))
	headers.Set("MIME-Version", "1.0")
	headers.Set("Content-Type", `text/plain; charset="utf-8"`)
	headers.Set("Content-Transfer-Encoding", "8bit")

	var buf bytes.Buffer
	for _, key := range []string{"From", "To", "Reply-To", "Subject", "MIME-Version", "Content-Type", "Content-Transfer-Encoding"} {
		buf.WriteString(key)
		buf.WriteString(": ")
		buf.WriteString(headers.Get(key))
		buf.WriteString("\r\n")
	}
	buf.WriteString("\r\n")
	buf.WriteString(msg.Body)
	if !strings.HasSuffix(msg.Body, "\n") {
		buf.WriteString("\r\n")
	}
	return buf.Bytes()
}

func contactBody(msg ContactMessage) string {
	var b strings.Builder
	b.WriteString("New contact form message\n\n")
	b.WriteString("Name: ")
	b.WriteString(msg.Name)
	b.WriteString("\nEmail: ")
	b.WriteString(msg.Email)
	b.WriteString("\nRemote IP: ")
	b.WriteString(msg.RemoteIP)
	b.WriteString("\nUser-Agent: ")
	b.WriteString(msg.UserAgent)
	b.WriteString("\n\nMessage:\n")
	b.WriteString(msg.Message)
	b.WriteString("\n")
	return b.String()
}
