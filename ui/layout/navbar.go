package layout

import (
	"mljr-web/ui/icon"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type NavbarProps struct {
	Attrs []g.Node
}

// Navbar takes 3 slot groups: brand, nav, actions. Below the menu breakpoint
// (≤640px), nav + actions collapse into a single dropdown panel behind a
// hamburger button, toggled via a Datastar signal scoped to this header.
func Navbar(p NavbarProps, brand, nav, actions g.Node) g.Node {
	return h.Header(
		g.Attr("data-component", "navbar"),
		g.Attr("data-signals", "{navOpen:0}"),
		g.Group(p.Attrs),
		h.Div(g.Attr("data-slot", "brand"), brand),
		h.Div(
			g.Attr("data-slot", "menu"),
			g.Attr("data-class:open", "$navOpen===1"),
			h.Nav(g.Attr("data-slot", "nav"), nav),
			h.Div(g.Attr("data-slot", "actions"), actions),
		),
		h.Button(
			h.Type("button"),
			g.Attr("data-slot", "menu-toggle"),
			g.Attr("aria-label", "Toggle navigation menu"),
			g.Attr("data-on:click", "$navOpen=$navOpen===1?0:1"),
			h.Span(h.Style("display:none"), g.Attr("data-show", "$navOpen===0"), icon.Icon("lucide:menu")),
			h.Span(h.Style("display:none"), g.Attr("data-show", "$navOpen===1"), icon.Icon("lucide:x")),
		),
	)
}
