package layout

import (
	"mljr-web/ui/token"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type PageProps struct {
	Title       string
	Description string
	Theme       token.Theme // default swissbrut
	Mode        token.Mode  // default light
	CSSPath     string      // default /static/app.css
	JSPath      string      // default /static/datastar.js
	HeadExtra   []g.Node
	BodyAttrs   []g.Node
}

// PageShell renders the full HTML document. The pre-paint inline script reads
// localStorage to set data-theme/data-mode on <html> before first paint so the
// chosen theme survives reloads without FOUC.
func PageShell(p PageProps, body ...g.Node) g.Node {
	if p.Theme == "" {
		p.Theme = token.ThemeSwissBrut
	}
	if p.Mode == "" {
		p.Mode = token.ModeLight
	}
	if p.CSSPath == "" {
		p.CSSPath = "/static/app.css?v=20260610-snake"
	}
	if p.JSPath == "" {
		p.JSPath = "/static/datastar.js"
	}

	// Pre-paint: read localStorage, set data-theme/data-mode on <html>,
	// expose values as window.__mljr* so Datastar signal seeds can read them.
	prepaint := `(function(){var t='` + string(p.Theme) + `',m='` + string(p.Mode) + `';try{t=localStorage.getItem('mljr-theme')||t;m=localStorage.getItem('mljr-mode')||m;}catch(e){}var r=document.documentElement;r.setAttribute('data-theme',t);r.setAttribute('data-mode',m);window.__mljrTheme=t;window.__mljrMode=m;})();`

	return g.Group{
		g.Raw("<!doctype html>"),
		h.HTML(
			h.Lang("en"),
			g.Attr("data-theme", string(p.Theme)),
			g.Attr("data-mode", string(p.Mode)),
			h.Head(
				h.Meta(h.Charset("utf-8")),
				h.Meta(h.Name("viewport"), h.Content("width=device-width,initial-scale=1")),
				h.TitleEl(g.Text(p.Title)),
				g.If(p.Description != "",
					h.Meta(h.Name("description"), h.Content(p.Description))),
				h.Link(h.Rel("icon"), h.Href("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 64 64'%3E%3Crect width='64' height='64' fill='%23e9473f'/%3E%3Cpath d='M14 46V18h8l10 15 10-15h8v28h-8V31L32 46 22 31v15z' fill='%23fffaf0'/%3E%3C/svg%3E")),
				h.Link(h.Rel("stylesheet"), h.Href(p.CSSPath)),
				// Pre-paint runs synchronously (no defer/async) to set
				// <html data-*> and window.__mljr* before first paint and
				// before Datastar initializes. Datastar must be type=module
				// (the v1.0.2 bundle uses ES module export{}).
				h.Script(g.Raw(prepaint)),
				h.Script(h.Type("module"), h.Src(p.JSPath)),
				g.Group(p.HeadExtra),
			),
			h.Body(
				g.Group(p.BodyAttrs),
				g.Group(body),
			),
		),
	}
}
