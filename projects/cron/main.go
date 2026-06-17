package main

import (
	"embed"
	"log"

	"mljr-web/internal/config"
	"mljr-web/internal/web"
	"mljr-web/projects/cron/pages"

	"github.com/labstack/echo/v4"
	"github.com/starfederation/datastar-go/datastar"
)

//go:embed all:assets/static
var assets embed.FS

func main() {
	cfg := config.Load()
	if cfg.Port == "8090" {
		cfg.Port = "8093"
	}

	e := web.NewEcho()
	web.MountStatic(e, assets, "projects/cron/assets/static", web.IsDev())

	e.GET("/", func(c echo.Context) error {
		return web.Render(c, 200, pages.Home())
	})

	e.POST("/api/eval", func(c echo.Context) error {
		var s struct {
			Expr    string `json:"expr"`
			WithSec bool   `json:"withSec"`
		}
		if err := datastar.ReadSignals(c.Request(), &s); err != nil {
			return err
		}
		result := pages.EvalCron(pages.CronInput{
			Expression: s.Expr,
			Count:      10,
			WithSec:    s.WithSec,
		})
		sse := datastar.NewSSE(c.Response().Writer, c.Request())
		return sse.PatchElements(web.RenderToString(pages.OutputFragment(result)))
	})

	e.POST("/api/example", func(c echo.Context) error {
		var s struct {
			Example string `json:"example"`
		}
		if err := datastar.ReadSignals(c.Request(), &s); err != nil {
			return err
		}
		ex, ok := pages.FindExample(s.Example)
		if !ok {
			return c.NoContent(204)
		}

		sse := datastar.NewSSE(c.Response().Writer, c.Request())
		if err := sse.MarshalAndPatchSignals(struct {
			Expr string `json:"expr"`
		}{ex.Expr}); err != nil {
			return err
		}

		result := pages.EvalCron(pages.CronInput{Expression: ex.Expr, Count: 10})
		return sse.PatchElements(web.RenderToString(pages.OutputFragment(result)))
	})

	e.GET("/healthz", func(c echo.Context) error { return c.String(200, "ok") })
	e.GET("/favicon.ico", func(c echo.Context) error { return c.NoContent(204) })

	log.Printf("cron explorer listening on :%s (env=%s)", cfg.Port, cfg.Env)
	if err := e.Start(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
