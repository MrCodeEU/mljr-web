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

func TestLoadReadsAnalyticsConfig(t *testing.T) {
	t.Setenv("UMAMI_SCRIPT_SRC", "/umami/script.js")
	t.Setenv("UMAMI_WEBSITE_ID", "site-id")
	t.Setenv("UMAMI_HOST_URL", "/umami")
	t.Setenv("UMAMI_DOMAINS", "mljr.eu,www.mljr.eu")
	t.Setenv("UMAMI_PROXY_TARGET", "https://stats.example.com")

	cfg := Load()
	if cfg.Analytics.UmamiScriptSrc != "/umami/script.js" {
		t.Fatalf("UmamiScriptSrc = %q", cfg.Analytics.UmamiScriptSrc)
	}
	if cfg.Analytics.UmamiWebsiteID != "site-id" {
		t.Fatalf("UmamiWebsiteID = %q", cfg.Analytics.UmamiWebsiteID)
	}
	if cfg.Analytics.UmamiHostURL != "/umami" {
		t.Fatalf("UmamiHostURL = %q", cfg.Analytics.UmamiHostURL)
	}
	if cfg.Analytics.UmamiDomains != "mljr.eu,www.mljr.eu" {
		t.Fatalf("UmamiDomains = %q", cfg.Analytics.UmamiDomains)
	}
	if cfg.Analytics.UmamiProxyTarget != "https://stats.example.com" {
		t.Fatalf("UmamiProxyTarget = %q", cfg.Analytics.UmamiProxyTarget)
	}
}
