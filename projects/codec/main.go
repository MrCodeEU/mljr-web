package main

import (
	"embed"
	"log"

	"mljr-web/internal/config"
	"mljr-web/internal/web"
	"mljr-web/projects/codec/pages"

	"github.com/labstack/echo/v4"
	"github.com/starfederation/datastar-go/datastar"
)

//go:embed all:assets/static
var assets embed.FS

func main() {
	cfg := config.Load()
	if cfg.Port == "8090" {
		cfg.Port = "8094"
	}

	e := web.NewEcho()
	web.MountStatic(e, assets, "projects/codec/assets/static", web.IsDev())

	e.GET("/", func(c echo.Context) error {
		return web.Render(c, 200, pages.Home())
	})

	e.POST("/api/eval", func(c echo.Context) error {
		var s struct {
			Mode  string `json:"mode"`
			Op    string `json:"op"`
			Input string `json:"input"`
		}
		if err := datastar.ReadSignals(c.Request(), &s); err != nil {
			return err
		}
		result := pages.EvalCodec(pages.CodecInput{
			Mode:  s.Mode,
			Op:    s.Op,
			Input: s.Input,
		})
		sse := datastar.NewSSE(c.Response().Writer, c.Request())
		return sse.PatchElements(web.RenderToString(pages.OutputFragment(result)))
	})

	e.GET("/healthz", func(c echo.Context) error { return c.String(200, "ok") })
	e.GET("/favicon.ico", func(c echo.Context) error { return c.NoContent(204) })

	log.Printf("codec listening on :%s (env=%s)", cfg.Port, cfg.Env)
	if err := e.Start(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
