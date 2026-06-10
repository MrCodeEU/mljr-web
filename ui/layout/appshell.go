package layout

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type AppShellProps struct {
	// SidebarWidth is the CSS width of the sidebar when open (default "240px").
	SidebarWidth string
	// MinHeight is the min-height of the shell (default "100vh").
	MinHeight string
}

// AppShell renders a sidebar + main content layout.
// Pass Sidebar() as the sidebar slot and any content as main.
// Both slots are plain g.Node — use layout.Sidebar for a collapsible sidebar.
func AppShell(p AppShellProps, sidebarSlot g.Node, mainSlot ...g.Node) g.Node {
	if p.MinHeight == "" {
		p.MinHeight = "100vh"
	}

	return h.Div(
		g.Attr("data-component", "app-shell"),
		h.Style("display:flex;min-height:"+p.MinHeight),
		h.Div(g.Attr("data-slot", "sidebar"), sidebarSlot),
		h.Main(
			g.Attr("data-slot", "main"),
			h.Style("flex:1;min-width:0;overflow-y:auto"),
			g.Group(mainSlot),
		),
	)
}

// AuthLayout renders a centered single-column card layout for login/register pages.
type AuthLayoutProps struct {
	MaxWidth string // CSS max-width (default "420px")
	Logo     g.Node // optional logo above the card
}

func AuthLayout(p AuthLayoutProps, card g.Node) g.Node {
	if p.MaxWidth == "" {
		p.MaxWidth = "420px"
	}
	return h.Div(
		g.Attr("data-component", "auth-layout"),
		h.Div(
			g.Attr("data-slot", "inner"),
			h.Style("width:100%;max-width:"+p.MaxWidth),
			g.If(p.Logo != nil, h.Div(g.Attr("data-slot", "logo"), p.Logo)),
			card,
		),
	)
}
