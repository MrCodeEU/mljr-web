//go:build showcase

package patterns

import (
	"mljr-web/ui/form"
	"mljr-web/ui/icon"
	"mljr-web/ui/layout"
	"mljr-web/ui/primitive"
	"mljr-web/ui/registry"
	"mljr-web/ui/special"
	"mljr-web/ui/token"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.RegisterPattern(&registry.Pattern{
		Slug:        "app-settings",
		Name:        "Settings Page",
		Category:    "app",
		Description: "Two-column settings layout: section nav + tabbed form sections. Profile, security, notifications, billing.",
		Render: func(theme, mode string) g.Node {
			th := token.Theme(theme)
			mo := token.Mode(mode)
			if th == "" {
				th = token.ThemeSwissBrut
			}
			if mo == "" {
				mo = token.ModeLight
			}

			sections := []struct {
				ic, label string
				active    bool
			}{
				{"lucide:user", "Profile", true},
				{"lucide:lock", "Security", false},
				{"lucide:bell", "Notifications", false},
				{"lucide:credit-card", "Billing", false},
				{"lucide:trash-2", "Danger Zone", false},
			}
			navItems := make([]g.Node, len(sections))
			for i, s := range sections {
				style := "display:flex;align-items:center;gap:var(--sp-2);padding:var(--sp-2) var(--sp-3);border-radius:var(--radius);font-size:var(--t-sm);font-weight:600;color:var(--muted);text-decoration:none;cursor:pointer"
				if s.active {
					style += ";background:var(--surface-2);color:var(--ink)"
				}
				navItems[i] = h.A(h.Href("#"), h.Style(style),
					icon.Icon(s.ic, icon.Props{Size: "1rem"}), g.Text(s.label))
			}

			profileSection := settingsSection("Profile",
				"Update your name, email, and avatar.",
				h.Div(
					h.Style("display:flex;align-items:flex-start;gap:var(--sp-5);margin-bottom:var(--sp-5)"),
					special.UserMenu(special.UserMenuProps{
						Initials: "AD", Size: token.SizeLG, Signal: "_sUm",
						Items: []special.UserMenuItem{{Label: "Change photo", Icon: "lucide:camera"}},
					}),
					h.Div(h.Style("font-size:var(--t-sm);color:var(--muted)"),
						h.P(g.Text("JPG, PNG or WebP")),
						h.P(g.Text("Max 2MB")),
					),
				),
				h.Div(
					h.Style("display:grid;grid-template-columns:1fr 1fr;gap:var(--sp-4)"),
					formFieldSimple("First name", h.Input(g.Attr("data-component", "input"),
						h.Type("text"), h.Value("Alex"), h.Style("width:100%"))),
					formFieldSimple("Last name", h.Input(g.Attr("data-component", "input"),
						h.Type("text"), h.Value("Developer"), h.Style("width:100%"))),
				),
				formFieldSimple("Email", h.Input(g.Attr("data-component", "input"),
					h.Type("email"), h.Value("alex@mljr.eu"), h.Style("width:100%"))),
				form.Field(form.FieldProps{Label: "Bio", Hint: "Max 160 characters"},
					form.Textarea(form.TextareaProps{
						Name:        "bio",
						Placeholder: "Tell us a bit about yourself…",
						Rows:        3,
					}),
				),
			)

			notifSection := settingsSection("Notifications",
				"Choose what updates you'd like to receive.",
				notifRow("Product updates", "Receive emails about new features and improvements.", true),
				notifRow("Security alerts", "Get notified about suspicious activity on your account.", true),
				notifRow("Marketing emails", "Promotional offers and newsletters.", false),
				notifRow("Weekly digest", "A weekly summary of your activity.", false),
			)

			return fullPage(th, mo,
				h.Div(
					layout.Navbar(layout.NavbarProps{},
						h.A(h.Href("#"), g.Text("mljr app")),
						g.Group{h.A(h.Href("#"), g.Text("Dashboard"))},
						special.ThemeToggle(),
					),
					h.Div(
						layout.Container(layout.ContainerProps{},
							h.Div(
								h.Style("display:grid;grid-template-columns:220px 1fr;gap:var(--sp-8);padding:var(--sp-8) 0"),
								// Sidebar nav
								h.Aside(
									h.Nav(h.Style("display:flex;flex-direction:column;gap:var(--sp-1)"),
										g.Group(navItems),
									),
								),
								// Content
								h.Div(
									h.Style("display:flex;flex-direction:column;gap:var(--sp-8)"),
									profileSection,
									notifSection,
									// Save
									h.Div(
										h.Style("display:flex;justify-content:flex-end;gap:var(--sp-3)"),
										primitive.Button(primitive.ButtonProps{Variant: token.Ghost}, g.Text("Cancel")),
										primitive.Button(primitive.ButtonProps{Variant: token.Primary}, g.Text("Save changes")),
									),
								),
							),
						),
					),
				),
			)
		},
	})
}

func settingsSection(title, desc string, children ...g.Node) g.Node {
	nodes := append([]g.Node{
		h.H2(h.Style("font-size:var(--t-base);font-weight:800;margin-bottom:var(--sp-1)"), g.Text(title)),
		h.P(h.Style("font-size:var(--t-sm);color:var(--muted);margin-bottom:var(--sp-5)"), g.Text(desc)),
	}, children...)
	return h.Div(
		g.Attr("data-component", "card"),
		h.Style("padding:var(--sp-6)"),
		g.Group(nodes),
	)
}

func formFieldSimple(label string, control g.Node) g.Node {
	return h.Div(
		h.Style("margin-bottom:var(--sp-3)"),
		h.Label(h.Style("display:block;font-weight:700;font-size:var(--t-sm);margin-bottom:var(--sp-1)"), g.Text(label)),
		control,
	)
}

func notifRow(title, desc string, defaultOn bool) g.Node {
	return h.Div(
		h.Style("display:flex;align-items:flex-start;justify-content:space-between;gap:var(--sp-4);padding:var(--sp-3) 0;border-bottom:var(--bw-1) solid var(--line)"),
		h.Div(
			h.Style("flex:1"),
			h.Div(h.Style("font-weight:700;font-size:var(--t-sm)"), g.Text(title)),
			h.Div(h.Style("font-size:var(--t-xs);color:var(--muted)"), g.Text(desc)),
		),
		form.Switch(form.SwitchProps{Checked: defaultOn}),
	)
}
