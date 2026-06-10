package layout

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type NavbarProps struct {
	Attrs []g.Node
}

// Navbar takes 3 slot groups: brand, nav, actions.
func Navbar(p NavbarProps, brand, nav, actions g.Node) g.Node {
	return h.Header(
		g.Attr("data-component", "navbar"),
		g.Group(p.Attrs),
		h.Div(g.Attr("data-slot", "brand"), brand),
		h.Nav(g.Attr("data-slot", "nav"), nav),
		h.Div(g.Attr("data-slot", "actions"), actions),
	)
}
