package main

import (
	"embed"
	"log"

	"mljr-web/internal/config"
	"mljr-web/internal/mail"
	"mljr-web/internal/web"
	hpdata "mljr-web/projects/homepage/data"
	"mljr-web/projects/homepage/pages"

	"github.com/labstack/echo/v4"
)

//go:embed all:assets/static
var assets embed.FS

func main() {
	cfg := config.Load()
	e := web.NewEcho()

	web.MountStatic(e, assets, "projects/homepage/assets/static", web.IsDev())

	d := hpdata.Load()
	e.GET("/", func(c echo.Context) error {
		return web.Render(c, 200, pages.Home(d))
	})
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

	log.Printf("homepage listening on :%s (env=%s)", cfg.Port, cfg.Env)
	if err := e.Start(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
