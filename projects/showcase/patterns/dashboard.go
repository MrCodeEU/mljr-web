//go:build showcase

package patterns

import (
	"mljr-web/ui/icon"
	"mljr-web/ui/primitive"
	"mljr-web/ui/registry"
	"mljr-web/ui/special"
	"mljr-web/ui/token"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.RegisterPattern(&registry.Pattern{
		Slug:        "app-dashboard",
		Name:        "App Dashboard",
		Category:    "app",
		Description: "Analytics dashboard with sidebar nav, stat cards, and data table.",
		Render: func(theme, mode string) g.Node {
			th := token.Theme(theme)
			mo := token.Mode(mode)
			if th == "" {
				th = token.ThemeSwissBrut
			}
			if mo == "" {
				mo = token.ModeLight
			}

			sidebar := h.Aside(
				h.Style("width:220px;flex-shrink:0;border-right:var(--bw-2) solid var(--ink);padding:var(--sp-5);display:flex;flex-direction:column;gap:var(--sp-2);background:var(--surface)"),
				h.Div(
					h.Style("font-weight:800;font-size:var(--t-lg);margin-bottom:var(--sp-4)"),
					g.Text("mljr app"),
				),
				navLink("lucide:layout-dashboard", "Dashboard", true),
				navLink("lucide:bar-chart-2", "Analytics", false),
				navLink("lucide:users", "Users", false),
				navLink("lucide:settings", "Settings", false),
				h.Div(h.Style("flex:1")),
				special.UserMenu(special.UserMenuProps{
					Name:     "Alex Dev",
					Email:    "alex@mljr.eu",
					Initials: "AD",
					Signal:   "_dashUm",
					Align:    "left",
					Items: []special.UserMenuItem{
						{Label: "Profile", Icon: "lucide:user", Href: "#"},
						{Divider: true},
						{Label: "Sign out", Icon: "lucide:log-out", Danger: true},
					},
				}),
			)

			stats := h.Div(
				h.Style("display:grid;grid-template-columns:repeat(4,1fr);gap:var(--sp-4)"),
				statCard("Revenue", "$128.5K", "+12.4%", true),
				statCard("Active Users", "42,317", "+8.1%", true),
				statCard("Conversion", "3.24%", "-0.6%", false),
				statCard("Avg. Session", "4m 12s", "+22s", true),
			)

			tableRows := []struct{ email, name, plan, joined, status string }{
				{"jane@example.com", "Jane Smith", "Pro", "2026-06-01", "Active"},
				{"alex@mljr.eu", "Alex Dev", "Enterprise", "2026-05-28", "Active"},
				{"bob@test.io", "Bob Builder", "Free", "2026-05-15", "Churned"},
				{"carol@firm.co", "Carol White", "Pro", "2026-06-03", "Active"},
				{"dave@corp.net", "Dave Corp", "Enterprise", "2026-06-07", "Trial"},
			}
			trows := make([]g.Node, len(tableRows))
			for i, r := range tableRows {
				trows[i] = h.Tr(
					h.Td(g.Text(r.email)),
					h.Td(g.Text(r.name)),
					h.Td(g.Text(r.plan)),
					h.Td(g.Text(r.joined)),
					h.Td(h.Span(g.Attr("data-component", "badge"), g.Text(r.status))),
				)
			}

			tableNode := h.Div(
				g.Attr("data-component", "table"),
				h.Table(
					h.THead(h.Tr(
						h.Th(g.Text("Email")), h.Th(g.Text("Name")), h.Th(g.Text("Plan")),
						h.Th(g.Text("Joined")), h.Th(g.Text("Status")),
					)),
					h.TBody(g.Group(trows)),
				),
			)

			main := h.Div(
				h.Style("flex:1;overflow-y:auto;padding:var(--sp-6)"),
				h.Div(
					h.Style("display:flex;align-items:center;justify-content:space-between;margin-bottom:var(--sp-6)"),
					h.H1(h.Style("font-size:var(--t-xl);font-weight:800"), g.Text("Dashboard")),
					h.Div(h.Style("display:flex;gap:var(--sp-2)"),
						primitive.Button(primitive.ButtonProps{Variant: token.Outline, Size: token.SizeSM}, g.Text("Export")),
						primitive.Button(primitive.ButtonProps{Variant: token.Primary, Size: token.SizeSM},
							icon.Icon("lucide:plus", icon.Props{Size: "0.9rem"}), g.Text("New Report"),
						),
					),
				),
				stats,
				h.Div(h.Style("margin-top:var(--sp-6)"),
					h.H2(h.Style("font-size:var(--t-base);font-weight:700;margin-bottom:var(--sp-3)"), g.Text("Recent Users")),
					tableNode,
				),
			)

			return fullPage(th, mo,
				h.Div(
					h.Style("display:flex;height:100vh;overflow:hidden"),
					sidebar,
					main,
				),
			)
		},
	})
}

func navLink(ic, label string, active bool) g.Node {
	style := "display:flex;align-items:center;gap:var(--sp-2);padding:var(--sp-2) var(--sp-3);border-radius:var(--radius);font-size:var(--t-sm);font-weight:600;text-decoration:none"
	if active {
		style += ";background:var(--accent);color:var(--accent-ink)"
	} else {
		style += ";color:var(--muted)"
	}
	return h.A(h.Href("#"), h.Style(style),
		icon.Icon(ic, icon.Props{Size: "1rem"}),
		g.Text(label),
	)
}

func statCard(label, value, delta string, positive bool) g.Node {
	deltaColor := "var(--success)"
	if !positive {
		deltaColor = "var(--danger)"
	}
	return h.Div(
		g.Attr("data-component", "card"),
		h.Div(h.Style("font-size:var(--t-xs);color:var(--muted);font-weight:600;margin-bottom:var(--sp-1)"), g.Text(label)),
		h.Div(h.Style("font-size:var(--t-xl);font-weight:800;line-height:1"), g.Text(value)),
		h.Div(h.Style("font-size:var(--t-xs);font-weight:700;color:"+deltaColor+";margin-top:var(--sp-1)"), g.Text(delta+" vs last month")),
	)
}
