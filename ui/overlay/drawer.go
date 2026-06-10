package overlay

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type DrawerPlacement string

const (
	DrawerRight DrawerPlacement = "right"
	DrawerLeft  DrawerPlacement = "left"
)

type DrawerProps struct {
	ID        string // default "drawer"
	Title     string
	OpenExpr  string          // Datastar open expression (default "$drawerOpen")
	Placement DrawerPlacement // default: right
	Size      string          // sm | md (default) | lg
	Attrs     []g.Node
}

// Drawer renders a Datastar-gated side panel. Pair with a trigger that sets the open signal.
func Drawer(p DrawerProps, children ...g.Node) g.Node {
	if p.ID == "" {
		p.ID = "drawer"
	}
	if p.OpenExpr == "" {
		p.OpenExpr = "$drawerOpen"
	}
	placement := string(p.Placement)
	if placement == "" {
		placement = "right"
	}

	closeExpr := p.OpenExpr + " = false"

	header := g.If(p.Title != "",
		h.Div(
			g.Attr("data-slot", "header"),
			h.Span(g.Text(p.Title)),
			h.Button(
				g.Attr("data-component", "button"),
				g.Attr("data-variant", "ghost"),
				g.Attr("data-size", "icon"),
				g.Attr("aria-label", "Close"),
				g.Attr("data-on:click", closeExpr),
				g.Text("×"),
			),
		),
	)

	drawerAttrs := []g.Node{
		h.ID(p.ID),
		g.Attr("data-component", "drawer"),
		g.Attr("data-placement", placement),
		g.If(p.Size != "", g.Attr("data-size", p.Size)),
		g.Attr("data-on:click__stop", "void 0"),
	}
	drawerAttrs = append(drawerAttrs, p.Attrs...)
	drawerAttrs = append(drawerAttrs, header)
	drawerAttrs = append(drawerAttrs, h.Div(g.Attr("data-slot", "body"), g.Group(children)))

	return h.Div(
		h.ID(p.ID+"-scrim"),
		g.Attr("data-component", "drawer-scrim"),
		g.Attr("data-show", p.OpenExpr),
		h.Style("display:none"),
		g.Attr("data-on:click", closeExpr),
		h.Div(drawerAttrs...),
	)
}
