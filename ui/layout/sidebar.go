package layout

import (
	"mljr-web/ui/icon"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type SidebarProps struct {
	Width          string // CSS width when open (default "240px")
	DefaultOpen    bool
	SignalName     string // Datastar signal name (default "_sidebarOpen")
	CollapsedWidth string // CSS width when collapsed (default "56px")
}

// Sidebar renders a collapsible left navigation sidebar driven by a Datastar signal.
// Place inside a flex-row container alongside your main content.
// Use SidebarItem() / SidebarSection() to build nav items.
func Sidebar(p SidebarProps, children ...g.Node) g.Node {
	if p.Width == "" {
		p.Width = "240px"
	}
	if p.CollapsedWidth == "" {
		p.CollapsedWidth = "56px"
	}
	if p.SignalName == "" {
		p.SignalName = "_sidebarOpen"
	}

	openDefault := "false"
	if p.DefaultOpen {
		openDefault = "true"
	}

	sig := p.SignalName
	toggleExpr := "$" + sig + "=!$" + sig
	isOpen := "$" + sig
	widthExpr := `{"style":'width:'+($` + sig + `?"` + p.Width + `":"` + p.CollapsedWidth + `")}`

	return h.Aside(
		g.Attr("data-component", "sidebar"),
		g.Attr("data-signals", `{"`+sig+`":`+openDefault+`}`),
		g.Attr("data-attr", widthExpr),
		h.Style("width:"+func() string {
			if p.DefaultOpen {
				return p.Width
			}
			return p.CollapsedWidth
		}()),

		// Toggle button at top
		h.Div(
			g.Attr("data-slot", "header"),
			h.Button(
				g.Attr("data-component", "button"),
				g.Attr("data-variant", "ghost"),
				g.Attr("data-size", "icon"),
				h.Type("button"),
				g.Attr("data-on:click", toggleExpr),
				g.Attr("aria-label", "Toggle sidebar"),
				icon.Icon("lucide:menu"),
			),
		),

		// Nav content — hide labels when collapsed
		h.Nav(
			g.Attr("data-slot", "nav"),
			g.Attr("data-show", isOpen),
			g.Group(children),
		),
	)
}

type SidebarItemProps struct {
	Href   string
	Active bool
}

func SidebarItem(p SidebarItemProps, label g.Node, ico ...g.Node) g.Node {
	nodes := []g.Node{}
	if len(ico) > 0 {
		nodes = append(nodes, ico[0])
	}
	nodes = append(nodes, h.Span(g.Attr("data-slot", "label"), label))
	attrs := []g.Node{
		g.Attr("data-slot", "item"),
		h.Href(p.Href),
	}
	if p.Active {
		attrs = append(attrs, g.Attr("data-active", ""))
	}
	return h.A(append(attrs, g.Group(nodes))...)
}

func SidebarSection(label string, children ...g.Node) g.Node {
	return h.Div(
		g.Attr("data-slot", "section"),
		h.Span(g.Attr("data-slot", "section-label"), g.Text(label)),
		g.Group(children),
	)
}
