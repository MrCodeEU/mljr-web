//go:build showcase

package layout

import (
	"mljr-web/ui/icon"
	"mljr-web/ui/primitive"
	"mljr-web/ui/registry"
	"mljr-web/ui/token"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "app-shell", Name: "App Shell", Category: "layout",
		PreviewHeight: "420px",
		Summary:       "Sidebar + main content layout. Wraps layout.Sidebar with a flex container.",
		Code: `layout.AppShell(layout.AppShellProps{},
    layout.Sidebar(layout.SidebarProps{DefaultOpen: true},
        layout.SidebarSection("Nav", /* items */),
    ),
    h.Div(/* main content */),
)`,
		Render: func(p map[string]string) g.Node {
			return h.Div(
				h.Style("height:380px;border:var(--bw-2) solid var(--line);border-radius:var(--radius);overflow:hidden"),
				AppShell(AppShellProps{MinHeight: "100%"},
					Sidebar(SidebarProps{DefaultOpen: true},
						SidebarSection("Main",
							SidebarItem(SidebarItemProps{Href: "#", Active: true}, g.Text("Dashboard"), icon.Icon("lucide:home")),
							SidebarItem(SidebarItemProps{Href: "#"}, g.Text("Projects"), icon.Icon("lucide:folder")),
							SidebarItem(SidebarItemProps{Href: "#"}, g.Text("Analytics"), icon.Icon("lucide:bar-chart-2")),
						),
						SidebarSection("Config",
							SidebarItem(SidebarItemProps{Href: "#"}, g.Text("Settings"), icon.Icon("lucide:settings")),
						),
					),
					h.Div(
						h.Style("padding:var(--sp-5)"),
						primitive.Heading(primitive.HeadingProps{Level: 2}, g.Text("Dashboard")),
						h.P(h.Style("color:var(--muted);font-size:var(--t-sm)"), g.Text("Main content area. Flex:1, min-width:0, overflow-y:auto.")),
					),
				),
			)
		},
	})

	registry.Register(&registry.Component{
		Slug: "auth-layout", Name: "Auth Layout", Category: "layout",
		Summary: "Centered single-column card layout for login and register pages.",
		Code: `layout.AuthLayout(layout.AuthLayoutProps{MaxWidth: "420px"},
    primitive.Card(primitive.CardProps{},
        h.Form(/* login form */),
    ),
)`,
		Render: func(p map[string]string) g.Node {
			return AuthLayout(AuthLayoutProps{},
				primitive.Card(primitive.CardProps{},
					h.Div(h.Style("display:flex;flex-direction:column;gap:var(--sp-5)"),
						primitive.Heading(primitive.HeadingProps{Level: 2}, g.Text("Sign in")),
						h.P(h.Style("color:var(--muted);margin:0"), g.Text("Welcome back — enter your credentials to continue.")),
						h.Div(
							h.Style("display:flex;flex-direction:column;gap:var(--sp-3)"),
							h.Input(g.Attr("data-component", "input"), h.Type("email"), h.Placeholder("Email")),
							h.Input(g.Attr("data-component", "input"), h.Type("password"), h.Placeholder("Password")),
						),
						primitive.Button(primitive.ButtonProps{Variant: token.Primary},
							h.Style("width:100%"),
							g.Text("Sign in"),
						),
						h.P(h.Style("text-align:center;font-size:var(--t-sm);color:var(--muted)"),
							g.Text("No account? "),
							h.A(h.Href("#"), g.Text("Create one")),
						),
					),
				),
			)
		},
	})
}
