package mail

import (
	"context"
	"log"
)

type ContactMessage struct {
	Name      string
	Email     string
	Message   string
	RemoteIP  string
	UserAgent string
}

type ContactMailer interface {
	SendContact(ctx context.Context, msg ContactMessage) error
}

type LogMailer struct{}

func (LogMailer) SendContact(_ context.Context, msg ContactMessage) error {
	log.Printf("contact mail disabled: from %s <%s>: %.80s", msg.Name, msg.Email, msg.Message)
	return nil
}
