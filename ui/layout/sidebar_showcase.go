//go:build showcase

package layout

import (
	"mljr-web/ui/icon"
	"mljr-web/ui/registry"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "sidebar", Name: "Sidebar", Category: "layout",
		PreviewHeight: "400px",
		Summary: "Collapsible left navigation sidebar. Datastar signal controls open/collapsed state. Width transitions smoothly.",
		Code: `layout.Sidebar(layout.SidebarProps{DefaultOpen: true},
    layout.SidebarSection("Main",
        layout.SidebarItem(layout.SidebarItemProps{Href: "/", Active: true},
            g.Text("Dashboard"), icon.Icon("lucide:home")),
        layout.SidebarItem(layout.SidebarItemProps{Href: "/projects"},
            g.Text("Projects"), icon.Icon("lucide:folder")),
    ),
    layout.SidebarSection("Settings",
        layout.SidebarItem(layout.SidebarItemProps{Href: "/settings"},
            g.Text("Settings"), icon.Icon("lucide:settings")),
    ),
)`,
		Render: func(p map[string]string) g.Node {
			return h.Div(
				h.Style("display:flex;height:360px;border:var(--bw-2) solid var(--line);border-radius:var(--radius);overflow:hidden"),
				Sidebar(SidebarProps{DefaultOpen: true},
					SidebarSection("Main",
						SidebarItem(SidebarItemProps{Href: "#", Active: true}, g.Text("Dashboard"), icon.Icon("lucide:home")),
						SidebarItem(SidebarItemProps{Href: "#"}, g.Text("Projects"), icon.Icon("lucide:folder")),
						SidebarItem(SidebarItemProps{Href: "#"}, g.Text("Analytics"), icon.Icon("lucide:bar-chart-2")),
						SidebarItem(SidebarItemProps{Href: "#"}, g.Text("Users"), icon.Icon("lucide:users")),
					),
					SidebarSection("System",
						SidebarItem(SidebarItemProps{Href: "#"}, g.Text("Settings"), icon.Icon("lucide:settings")),
						SidebarItem(SidebarItemProps{Href: "#"}, g.Text("Log out"), icon.Icon("lucide:log-out")),
					),
				),
				h.Div(
					h.Style("flex:1;padding:var(--sp-5);background:var(--bg)"),
					h.P(h.Style("font-weight:700;margin:0 0 var(--sp-2)"), g.Text("Dashboard")),
					h.P(h.Style("color:var(--muted);font-size:var(--t-sm)"), g.Text("Click the ☰ button in the sidebar to collapse it.")),
				),
			)
		},
	})
}
