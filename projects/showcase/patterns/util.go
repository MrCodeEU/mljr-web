//go:build showcase

package patterns

import (
	"mljr-web/ui/layout"
	"mljr-web/ui/special"
	"mljr-web/ui/token"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

// fullPage wraps pattern content in a complete HTML document with theme + Datastar.
func fullPage(theme token.Theme, mode token.Mode, content g.Node) g.Node {
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
				h.Script(h.Type("module"), h.Src("/static/datastar.js")),
			),
			h.Body(
				h.Style("margin:0;background:var(--bg);color:var(--ink);min-height:100vh"),
				special.ThemeToggleRoot(theme, mode),
				content,
			),
		),
	}
}

// ensure layout is imported (used in some patterns)
var _ = layout.Container
