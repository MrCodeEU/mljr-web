package config

import "os"

type Config struct {
	Env       string // "dev" | "prod"
	Port      string
	AltchaKey string // HMAC key for altcha challenge signing; generate with: openssl rand -hex 32
	SMTP      SMTPConfig
	ContactTo string
}

type SMTPConfig struct {
	Host string
	Port string
	User string
	Pass string
	From string
}

func Load() Config {
	return Config{
		Env:       envOr("MLJR_ENV", "dev"),
		Port:      envOr("PORT", "8090"),
		AltchaKey: envOr("ALTCHA_HMAC_KEY", "dev-insecure-change-in-prod-please"),
		SMTP: SMTPConfig{
			Host: os.Getenv("SMTP_HOST"),
			Port: envOr("SMTP_PORT", "587"),
			User: os.Getenv("SMTP_USER"),
			Pass: os.Getenv("SMTP_PASS"),
			From: os.Getenv("SMTP_FROM"),
		},
		ContactTo: os.Getenv("CONTACT_TO"),
	}
}

func (c Config) MailConfigured() bool {
	return c.SMTP.Host != "" && c.SMTP.Port != "" && c.SMTP.From != "" && c.ContactTo != ""
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
