package main

import (
	"embed"
	"io/fs"
	"log"
	"os"

	"mljr-web/internal/config"
	"mljr-web/internal/web"
	_ "mljr-web/projects/newsletter/migrations"
	"mljr-web/projects/newsletter/pages"
	"mljr-web/projects/newsletter/scheduler"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

//go:embed all:assets/static
var assets embed.FS

func main() {
	cfg := config.Load()
	port := cfg.Port
	if port == "8090" {
		port = "8096"
	}

	app := pocketbase.NewWithConfig(pocketbase.Config{
		DefaultDataDir: cfg.Newsletter.DataDir,
		DefaultDev:     cfg.Env == "dev",
	})

	app.OnServe().BindFunc(func(e *core.ServeEvent) error {
		mailerClient, err := scheduler.BootstrapMailer(app, cfg.SMTP)
		if err != nil {
			return err
		}
		pages.SetMailer(mailerClient, cfg)

		app.Cron().MustAdd("newsletter_scan", "*/5 * * * *", func() {
			if err := scheduler.RunScan(app, mailerClient, cfg); err != nil {
				log.Printf("newsletter scan failed: %v", err)
			}
		})

		var staticFS fs.FS
		if web.IsDev() {
			staticFS = os.DirFS("projects/newsletter/assets/static")
		} else {
			sub, err := fs.Sub(assets, "assets/static")
			if err != nil {
				return err
			}
			staticFS = sub
		}
		e.Router.GET("/static/{path...}", apis.Static(staticFS, false))

		if err := pages.RegisterRoutes(e); err != nil {
			return err
		}
		return e.Next()
	})

	pages.RegisterHooks(app)

	// Inject a default listen address when run without explicit CLI args
	// (the normal "go run"/binary case for this project), so PORT/.env
	// config controls it like every other projects/* server.
	if len(os.Args) < 2 {
		os.Args = append(os.Args, "serve", "--http=0.0.0.0:"+port)
	}

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
