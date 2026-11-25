package main

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
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
	e.Pre(web.LangPrefix())

	if err := registerAnalyticsProxy(e, cfg.Analytics); err != nil {
		log.Fatal(err)
	}
	web.MountStatic(e, assets, "projects/homepage/assets/static", web.IsDev())
	web.MountFonts(e)
	web.MountLogos(e, assets, "projects/homepage/assets/static", web.IsDev())
	web.MountDataAssets(e, cfg.Data.File)

	dataStore := hpdata.NewStore(cfg.Data.File, cfg.Data.ReloadSeconds)
	analytics := pages.AnalyticsConfig{
		UmamiScriptSrc: cfg.Analytics.UmamiScriptSrc,
		UmamiWebsiteID: cfg.Analytics.UmamiWebsiteID,
		UmamiHostURL:   cfg.Analytics.UmamiHostURL,
		UmamiDomains:   cfg.Analytics.UmamiDomains,
	}
	// unknown pages get the site's 404 page; everything else the default
	defaultErr := e.HTTPErrorHandler
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		var he *echo.HTTPError
		if errors.As(err, &he) && he.Code == http.StatusNotFound && !c.Response().Committed &&
			strings.Contains(c.Request().Header.Get("Accept"), "text/html") {
			_ = web.Render(c, http.StatusNotFound, pages.NotFound(web.Lang(c), analytics))
			return
		}
		defaultErr(err, c)
	}
	// Live homelab panel: background poller, 60s cadence.
	hlPoller := homelab.New(homelab.Options{
		KumaURL:  cfg.Homelab.KumaURL,
		KumaSlug: cfg.Homelab.KumaSlug,
		PromURL:  cfg.Homelab.PromURL,
	})
	hlPoller.Start(context.Background(), 60*time.Second)
	hlSnapshot := func() homelab.Snapshot {
		snap := hlPoller.Snapshot()
		if !web.IsDev() {
			return snap
		}
		if !snap.KumaOK {
			return homelab.Sample()
		}
		// Dev without HOMELAB_PROM_URL: live Kuma, sample PromQL parts, so
		// bans, CPU and the attack heatmap can still be reviewed locally.
		if snap.ActiveBans < 0 {
			sm := homelab.Sample()
			snap.ActiveBans, snap.Attacks24h, snap.SecurityEvents, snap.HostsOnline = sm.ActiveBans, sm.Attacks24h, sm.SecurityEvents, sm.HostsOnline
			snap.BansCommunity, snap.BansLocal, snap.TopThreats = sm.BansCommunity, sm.BansLocal, sm.TopThreats
			snap.AttackDays, snap.CPUUtil, snap.CPULabels = sm.AttackDays, sm.CPUUtil, sm.CPULabels
		}
		if !snap.MeshOK {
			sm := homelab.Sample()
			snap.MeshOK, snap.Mesh, snap.MeshLinks, snap.MeshCircuits = true, sm.Mesh, sm.MeshLinks, sm.MeshCircuits
		}
		return snap
	}

	e.GET("/", func(c echo.Context) error {
		return web.Render(c, 200, pages.Home(dataStore.Current(), web.Lang(c), analytics, hlSnapshot()))
	})
	e.GET("/cv", func(c echo.Context) error {
		return web.Render(c, 200, pages.CV(dataStore.Current(), web.Lang(c), analytics))
	})
	e.GET("/notes/zero-trust", func(c echo.Context) error {
		return web.Render(c, 200, pages.ZeroTrustNote(dataStore.Current(), web.Lang(c), analytics))
	})
	e.GET("/impressum", func(c echo.Context) error {
		return web.Render(c, 200, pages.Impressum(web.Lang(c), analytics))
	})
	e.GET("/datenschutz", func(c echo.Context) error {
		return web.Render(c, 200, pages.Datenschutz(web.Lang(c), analytics))
	})
	e.GET("/games/minenfeld", func(c echo.Context) error {
		return web.Render(c, 200, pages.Minenfeld(web.Lang(c), analytics))
	})
	e.GET("/games/snake", func(c echo.Context) error {
		return web.Render(c, 200, pages.Snake(web.Lang(c), analytics))
	})
	e.GET("/games/ai-snake", func(c echo.Context) error {
		return web.Render(c, 200, pages.AISnake(web.Lang(c), analytics))
	})
	e.GET("/games/ai-snake-live", func(c echo.Context) error {
		return web.Render(c, 200, pages.AISnakeLive(web.Lang(c), analytics))
	})
	e.GET("/games/double-pendulum", func(c echo.Context) error {
		return web.Render(c, 200, pages.DoublePendulum(web.Lang(c), analytics))
	})
	e.GET("/games/fourier", func(c echo.Context) error {
		return web.Render(c, 200, pages.Fourier(web.Lang(c), analytics))
	})
	e.GET("/games/boids", func(c echo.Context) error {
		return web.Render(c, 200, pages.Boids(web.Lang(c), analytics))
	})
	e.GET("/games/particle-life", func(c echo.Context) error {
		return web.Render(c, 200, pages.ParticleLife(web.Lang(c), analytics))
	})
	e.GET("/games/maze", func(c echo.Context) error {
		return web.Render(c, 200, pages.Maze(web.Lang(c), analytics))
	})
	e.GET("/games/sorting", func(c echo.Context) error {
		return web.Render(c, 200, pages.Sorting(web.Lang(c), analytics))
	})
	e.GET("/projects/:id", func(c echo.Context) error {
		d := dataStore.Current()
		p, ok := d.ProjectByID(c.Param("id"))
		if !ok || !pages.HasProjectDetail(p) {
			return echo.ErrNotFound
		}
		return web.Render(c, 200, pages.ProjectDetail(d, p, web.Lang(c), analytics))
	})
	e.GET("/robots.txt", func(c echo.Context) error {
		return c.String(200, "User-agent: *\nAllow: /\nSitemap: https://mljr.eu/sitemap.xml\n")
	})
	e.HEAD("/robots.txt", func(c echo.Context) error { return c.NoContent(200) })
	// RFC 9116 security contact info. Expires is generated dynamically
	// (now + 1 year) so the file never silently goes stale.
	e.GET("/.well-known/security.txt", func(c echo.Context) error {
		expires := time.Now().UTC().AddDate(1, 0, 0).Format(time.RFC3339)
		body := "Contact: mailto:michael-reinegger@tuta.io\n" +
			"Expires: " + expires + "\n" +
			"Preferred-Languages: en\n"
		return c.String(200, body)
	})
	e.HEAD("/.well-known/security.txt", func(c echo.Context) error { return c.NoContent(200) })
	e.GET("/sitemap.xml", func(c echo.Context) error {
		c.Response().Header().Set("Content-Type", "application/xml; charset=utf-8")
		var b strings.Builder
		b.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
		b.WriteString(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9" xmlns:xhtml="http://www.w3.org/1999/xhtml">` + "\n")
		// every page twice, English and German, each listing both versions
		add := func(path, freq, prio string) {
			de := "/de" + path
			if path == "/" {
				de = "/de/"
			}
			alt := fmt.Sprintf(`<xhtml:link rel="alternate" hreflang="en" href="%[1]s%[2]s"/><xhtml:link rel="alternate" hreflang="de" href="%[1]s%[3]s"/><xhtml:link rel="alternate" hreflang="x-default" href="%[1]s%[2]s"/>`, web.SiteURL, path, de)
			for _, loc := range []string{path, de} {
				fmt.Fprintf(&b, "  <url><loc>%s%s</loc>%s<changefreq>%s</changefreq><priority>%s</priority></url>\n", web.SiteURL, loc, alt, freq, prio)
			}
		}
		add("/", "weekly", "1.0")
		add("/cv", "monthly", "0.8")
		add("/notes/zero-trust", "monthly", "0.7")
		add("/impressum", "monthly", "0.3")
		add("/datenschutz", "monthly", "0.3")
		for _, it := range pages.PlaygroundURLPaths() {
			add(it, "yearly", "0.4")
		}
		for _, p := range dataStore.Current().GitHub {
			if !p.HasDetailPage() {
				continue
			}
			add("/projects/"+p.ID, "monthly", "0.6")
		}
		b.WriteString(`</urlset>`)
		return c.String(200, b.String())
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
	registerSnakeLive(e)

	log.Printf("homepage listening on :%s (env=%s)", cfg.Port, cfg.Env)
	if err := e.Start(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
