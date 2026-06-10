//go:build showcase

package pages

import (
	"fmt"

	"mljr-web/ui/registry"
	"mljr-web/ui/special"
	"mljr-web/ui/token"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

// ComponentPreview renders a minimal iframe shell for a component.
// Theme and mode come from URL query params so the parent detail page can sync them.
func ComponentPreview(c *registry.Component, props map[string]string, theme token.Theme, mode token.Mode) g.Node {
	if theme == "" {
		theme = token.ThemeSwissBrut
	}
	if mode == "" {
		mode = token.ModeLight
	}
	seedScript := fmt.Sprintf(`window.__mljrTheme=%q;window.__mljrMode=%q;`, string(theme), string(mode))
	return g.Group{
		g.Raw("<!doctype html>"),
		h.HTML(
			h.Lang("en"),
			g.Attr("data-theme", string(theme)),
			g.Attr("data-mode", string(mode)),
			h.Head(
				h.Meta(h.Charset("utf-8")),
				h.Meta(h.Name("viewport"), h.Content("width=device-width,initial-scale=1")),
				h.Link(h.Rel("stylesheet"), h.Href("/static/app.css?v=20260610-snake")),
				h.Script(g.Raw(seedScript)),
				h.Script(h.Type("module"), h.Src("/static/datastar.js")),
				h.Script(h.Src("/static/motion.min.js")),
				g.El("style", g.Raw(`html{background:var(--bg);}body{margin:0;padding:clamp(1rem,4vw,2.5rem);box-sizing:border-box;background:var(--bg);min-height:100vh;overflow:auto;position:relative;}`)),
			),
			h.Body(
				special.ThemeToggleRoot(theme, mode),
				c.Render(props),
			),
		),
	}
}
