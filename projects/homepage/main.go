package main

import (
	"context"
	"embed"
	"log"
	"time"

	"mljr-web/internal/config"
	"mljr-web/internal/mail"
	"mljr-web/internal/web"
	hpdata "mljr-web/projects/homepage/data"
	"mljr-web/projects/homepage/homelab"
	"mljr-web/projects/homepage/pages"

	"github.com/labstack/echo/v4"
)

//go:embed all:assets/static
var assets embed.FS

func main() {
	cfg := config.Load()
	e := web.NewEcho()

	if err := registerAnalyticsProxy(e, cfg.Analytics); err != nil {
		log.Fatal(err)
	}
	web.MountStatic(e, assets, "projects/homepage/assets/static", web.IsDev())
	web.MountLogos(e, assets, "projects/homepage/assets/static", web.IsDev())
	web.MountDataAssets(e, cfg.Data.File)

	dataStore := hpdata.NewStore(cfg.Data.File, cfg.Data.ReloadSeconds)
	analytics := pages.AnalyticsConfig{
		UmamiScriptSrc: cfg.Analytics.UmamiScriptSrc,
		UmamiWebsiteID: cfg.Analytics.UmamiWebsiteID,
		UmamiHostURL:   cfg.Analytics.UmamiHostURL,
		UmamiDomains:   cfg.Analytics.UmamiDomains,
	}
	// Live homelab panel: background poller, 60s cadence.
	hlPoller := homelab.New(cfg.Homelab.KumaURL, cfg.Homelab.KumaSlug, cfg.Homelab.PromURL)
	hlPoller.Start(context.Background(), 60*time.Second)
	hlSnapshot := func() homelab.Snapshot {
		snap := hlPoller.Snapshot()
		if !snap.KumaOK && web.IsDev() {
			return homelab.Sample()
		}
		return snap
	}

	e.GET("/", func(c echo.Context) error {
		return web.Render(c, 200, pages.Home(dataStore.Current(), web.Lang(c), analytics, hlSnapshot()))
	})
	e.GET("/impressum", func(c echo.Context) error {
		return web.Render(c, 200, pages.Impressum(web.Lang(c), analytics))
	})
	e.GET("/datenschutz", func(c echo.Context) error {
		return web.Render(c, 200, pages.Datenschutz(web.Lang(c), analytics))
	})
	e.GET("/robots.txt", func(c echo.Context) error {
		return c.String(200, "User-agent: *\nAllow: /\nSitemap: https://mljr.eu/sitemap.xml\n")
	})
	e.HEAD("/robots.txt", func(c echo.Context) error { return c.NoContent(200) })
	e.GET("/sitemap.xml", func(c echo.Context) error {
		c.Response().Header().Set("Content-Type", "application/xml; charset=utf-8")
		return c.String(200, `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  <url><loc>https://mljr.eu/</loc><changefreq>weekly</changefreq><priority>1.0</priority></url>
  <url><loc>https://mljr.eu/impressum</loc><changefreq>monthly</changefreq><priority>0.3</priority></url>
  <url><loc>https://mljr.eu/datenschutz</loc><changefreq>monthly</changefreq><priority>0.3</priority></url>
</urlset>`)
	})
	e.HEAD("/sitemap.xml", func(c echo.Context) error { return c.NoContent(200) })
	e.GET("/healthz", func(c echo.Context) error { return c.String(200, "ok") })
	e.GET("/favicon.ico", func(c echo.Context) error { return c.NoContent(204) })

	contactMailer := mail.ContactMailer(mail.LogMailer{})
	if cfg.MailConfigured() {
		smtpMailer, err := mail.NewSMTPMailer(mail.SMTPConfig{
			Host: cfg.SMTP.Host,
			Port: cfg.SMTP.Port,
			User: cfg.SMTP.User,
			Pass: cfg.SMTP.Pass,
			From: cfg.SMTP.From,
			To:   cfg.ContactTo,
		})
		if err != nil {
			log.Fatal(err)
		}
		contactMailer = smtpMailer
	} else if cfg.Env == "prod" {
		log.Fatal("mail is not configured: set SMTP_HOST, SMTP_PORT, SMTP_FROM, and CONTACT_TO")
	} else {
		log.Printf("contact mail disabled: set SMTP_HOST, SMTP_FROM, and CONTACT_TO to send email")
	}

	registerHandlers(e, cfg.AltchaKey, contactMailer)
	registerHomelabHandler(e, hlSnapshot)

	log.Printf("homepage listening on :%s (env=%s)", cfg.Port, cfg.Env)
	if err := e.Start(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
