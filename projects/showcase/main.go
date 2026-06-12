//go:build showcase

package main

import (
	"embed"
	"fmt"
	"html"
	"log"
	"net/http"
	"strings"
	"time"

	"mljr-web/internal/config"
	"mljr-web/internal/web"
	"mljr-web/projects/showcase/pages"
	"mljr-web/ui/registry"
	"mljr-web/ui/token"

	// blank imports trigger init() in *_showcase.go files
	_ "mljr-web/projects/showcase/patterns"
	_ "mljr-web/ui/data"
	_ "mljr-web/ui/datastar"
	_ "mljr-web/ui/feedback"
	_ "mljr-web/ui/form"

	"github.com/starfederation/datastar-go/datastar"

	altchalib "github.com/altcha-org/altcha-lib-go"
	"github.com/labstack/echo/v4"
)

//go:embed all:assets/static
var assets embed.FS

func main() {
	cfg := config.Load()
	if cfg.Port == "8090" {
		cfg.Port = "8091"
	}
	e := web.NewEcho()

	web.MountStatic(e, assets, "projects/showcase/assets/static", web.IsDev())

	e.GET("/", func(c echo.Context) error {
		return web.Render(c, 200, pages.Catalogue())
	})
	e.GET("/healthz", func(c echo.Context) error { return c.String(200, "ok") })
	e.GET("/favicon.ico", func(c echo.Context) error { return c.NoContent(204) })

	e.GET("/components/:slug", func(c echo.Context) error {
		comp, ok := registry.Get(c.Param("slug"))
		if !ok {
			return echo.NewHTTPError(http.StatusNotFound, "component not found")
		}
		return web.Render(c, 200, pages.ComponentDetail(comp))
	})
	e.GET("/components/:slug/preview", func(c echo.Context) error {
		comp, ok := registry.Get(c.Param("slug"))
		if !ok {
			return echo.NewHTTPError(http.StatusNotFound, "component not found")
		}
		props := registry.DefaultProps(comp)
		for k, v := range c.QueryParams() {
			if len(v) > 0 {
				props[k] = v[0]
			}
		}
		theme := token.Theme(c.QueryParam("theme"))
		mode := token.Mode(c.QueryParam("mode"))
		return web.Render(c, 200, pages.ComponentPreview(comp, props, theme, mode))
	})

	// ── Patterns ─────────────────────────────────────────────────────────────

	e.GET("/patterns", func(c echo.Context) error {
		return web.Render(c, 200, pages.PatternsListing())
	})
	e.GET("/patterns/:slug", func(c echo.Context) error {
		p, ok := registry.GetPattern(c.Param("slug"))
		if !ok {
			return echo.NewHTTPError(http.StatusNotFound, "pattern not found")
		}
		theme := token.Theme(c.QueryParam("theme"))
		mode := token.Mode(c.QueryParam("mode"))
		return web.Render(c, 200, pages.PatternDetail(p, theme, mode))
	})
	e.GET("/patterns/:slug/preview", func(c echo.Context) error {
		p, ok := registry.GetPattern(c.Param("slug"))
		if !ok {
			return echo.NewHTTPError(http.StatusNotFound, "pattern not found")
		}
		theme := token.Theme(c.QueryParam("theme"))
		mode := token.Mode(c.QueryParam("mode"))
		return web.Render(c, 200, pages.PatternPreview(p, theme, mode))
	})

	// ── Datastar showcase demo routes ────────────────────────────────────────

	// POST /demo/echo — reads $q signal, patches #echo-result + updates $echoResult
	e.POST("/demo/echo", func(c echo.Context) error {
		var s struct {
			Q string `json:"_q"`
		}
		_ = datastar.ReadSignals(c.Request(), &s)
		msg := strings.TrimSpace(s.Q)
		if msg == "" {
			msg = "(empty message)"
		}
		sse := datastar.NewSSE(c.Response().Writer, c.Request())
		fragment := fmt.Sprintf(`<div id="echo-result">Server received at %s: <strong>%s</strong></div>`,
			time.Now().Format("15:04:05"), html.EscapeString(msg))
		sse.PatchElements(fragment)
		sse.MarshalAndPatchSignals(map[string]any{
			"echoResult": "Server received: " + msg,
			"_fetching":  false,
		})
		return nil
	})

	// GET /demo/search — filters a static fruit list by $q signal, patches #search-results
	e.GET("/demo/search", func(c echo.Context) error {
		var s struct {
			Q string `json:"_q"`
		}
		_ = datastar.ReadSignals(c.Request(), &s)
		q := strings.ToLower(strings.TrimSpace(s.Q))

		fruits := []string{
			"Apple", "Apricot", "Avocado", "Banana", "Blueberry",
			"Cherry", "Coconut", "Grape", "Guava", "Kiwi",
			"Lemon", "Lime", "Lychee", "Mango", "Melon",
			"Orange", "Papaya", "Peach", "Pear", "Pineapple",
			"Plum", "Pomegranate", "Raspberry", "Strawberry", "Watermelon",
		}

		var sb strings.Builder
		sb.WriteString(`<div id="search-results" style="display:grid;grid-template-columns:repeat(auto-fill,minmax(120px,1fr));gap:4px;margin-top:8px">`)
		count := 0
		for _, f := range fruits {
			if q == "" || strings.Contains(strings.ToLower(f), q) {
				sb.WriteString(`<div style="padding:6px 10px;border-radius:4px;font-size:.85rem;background:var(--surface-2)">`)
				sb.WriteString(html.EscapeString(f))
				sb.WriteString(`</div>`)
				count++
			}
		}
		if count == 0 {
			sb.WriteString(`<p style="opacity:.4;font-size:.85rem">No results</p>`)
		}
		sb.WriteString(`</div>`)

		sse := datastar.NewSSE(c.Response().Writer, c.Request())
		sse.PatchElements(sb.String())
		return nil
	})

	// GET /demo/time — persistent SSE stream, patches $serverTime every second for 15s
	e.GET("/demo/time", func(c echo.Context) error {
		sse := datastar.NewSSE(c.Response().Writer, c.Request())
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		deadline := time.After(15 * time.Second)
		for {
			select {
			case <-c.Request().Context().Done():
				return nil
			case <-deadline:
				sse.MarshalAndPatchSignals(map[string]any{
					"_streamActive": false,
				})
				return nil
			case t := <-ticker.C:
				sse.MarshalAndPatchSignals(map[string]any{
					"serverTime": t.Format("15:04:05"),
				})
			}
		}
	})

	// ── Altcha captcha ───────────────────────────────────────────────────────

	e.GET("/api/altcha", func(c echo.Context) error {
		ch, err := altchalib.CreateChallenge(altchalib.ChallengeOptions{
			HMACKey:   cfg.AltchaKey,
			MaxNumber: 200000,
		})
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "challenge creation failed")
		}
		return c.JSON(http.StatusOK, ch)
	})

	log.Printf("showcase listening on :%s (env=%s)", cfg.Port, cfg.Env)
	if err := e.Start(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
