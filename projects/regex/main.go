package main

import (
	"embed"
	"log"

	"mljr-web/internal/config"
	"mljr-web/internal/web"
	"mljr-web/projects/regex/pages"

	"github.com/labstack/echo/v4"
	"github.com/starfederation/datastar-go/datastar"
)

//go:embed all:assets/static
var assets embed.FS

func main() {
	cfg := config.Load()
	if cfg.Port == "8090" {
		cfg.Port = "8092"
	}

	e := web.NewEcho()
	web.MountStatic(e, assets, "projects/regex/assets/static", web.IsDev())

	e.GET("/", func(c echo.Context) error {
		return web.Render(c, 200, pages.Home())
	})

	e.POST("/api/eval", func(c echo.Context) error {
		var s struct {
			Pattern string `json:"pattern"`
			FlagI   bool   `json:"flagI"`
			FlagM   bool   `json:"flagM"`
			FlagS   bool   `json:"flagS"`
			Input   string `json:"input"`
			Replace string `json:"replace"`
		}
		if err := datastar.ReadSignals(c.Request(), &s); err != nil {
			return err
		}
		result := pages.EvalRegex(pages.EvalInput{
			Pattern: s.Pattern,
			FlagI:   s.FlagI,
			FlagM:   s.FlagM,
			FlagS:   s.FlagS,
			Input:   s.Input,
			Replace: s.Replace,
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

		// Update all input signals so the inputs reflect the example.
		if err := sse.MarshalAndPatchSignals(struct {
			Pattern string `json:"pattern"`
			FlagI   bool   `json:"flagI"`
			FlagM   bool   `json:"flagM"`
			FlagS   bool   `json:"flagS"`
			Input   string `json:"input"`
			Replace string `json:"replace"`
		}{ex.Pattern, ex.FlagI, ex.FlagM, ex.FlagS, ex.Input, ex.Replace}); err != nil {
			return err
		}

		// Evaluate and patch output.
		result := pages.EvalRegex(pages.EvalInput{
			Pattern: ex.Pattern,
			FlagI:   ex.FlagI,
			FlagM:   ex.FlagM,
			FlagS:   ex.FlagS,
			Input:   ex.Input,
			Replace: ex.Replace,
		})
		return sse.PatchElements(web.RenderToString(pages.OutputFragment(result)))
	})

	e.GET("/healthz", func(c echo.Context) error { return c.String(200, "ok") })
	e.GET("/favicon.ico", func(c echo.Context) error { return c.NoContent(204) })

	log.Printf("regex lab listening on :%s (env=%s)", cfg.Port, cfg.Env)
	if err := e.Start(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
