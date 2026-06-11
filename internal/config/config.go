package config

import "os"

type Config struct {
	Env       string // "dev" | "prod"
	Port      string
	AltchaKey string // HMAC key for altcha challenge signing; generate with: openssl rand -hex 32
	SMTP      SMTPConfig
	ContactTo string
	Analytics AnalyticsConfig
	Homelab   HomelabConfig
}

// HomelabConfig points the live homelab panel at its data sources.
// Both sources are optional; the panel degrades per-source.
type HomelabConfig struct {
	KumaURL  string // Uptime Kuma base URL (public status page API, no auth)
	KumaSlug string // status page slug (default "all")
	PromURL  string // Prometheus/VictoriaMetrics base URL via Tailscale; empty disables PromQL stats
}

type SMTPConfig struct {
	Host string
	Port string
	User string
	Pass string
	From string
}

type AnalyticsConfig struct {
	UmamiScriptSrc   string
	UmamiWebsiteID   string
	UmamiHostURL     string
	UmamiDomains     string
	UmamiProxyTarget string
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
		Analytics: AnalyticsConfig{
			UmamiScriptSrc:   os.Getenv("UMAMI_SCRIPT_SRC"),
			UmamiWebsiteID:   os.Getenv("UMAMI_WEBSITE_ID"),
			UmamiHostURL:     os.Getenv("UMAMI_HOST_URL"),
			UmamiDomains:     os.Getenv("UMAMI_DOMAINS"),
			UmamiProxyTarget: os.Getenv("UMAMI_PROXY_TARGET"),
		},
		Homelab: HomelabConfig{
			KumaURL:  envOr("HOMELAB_KUMA_URL", "https://uptime.mljr.eu"),
			KumaSlug: envOr("HOMELAB_KUMA_SLUG", "all"),
			PromURL:  os.Getenv("HOMELAB_PROM_URL"),
		},
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
