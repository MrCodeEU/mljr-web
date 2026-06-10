package config

import "testing"

func TestMailConfigured(t *testing.T) {
	cfg := Config{
		SMTP: SMTPConfig{
			Host: "smtp.example.com",
			Port: "587",
			From: "site@example.com",
		},
		ContactTo: "owner@example.com",
	}

	if !cfg.MailConfigured() {
		t.Fatal("MailConfigured() = false, want true")
	}

	cfg.ContactTo = ""
	if cfg.MailConfigured() {
		t.Fatal("MailConfigured() = true with empty ContactTo, want false")
	}
}
