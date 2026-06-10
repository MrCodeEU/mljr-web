//go:build showcase

package special

import (
	"mljr-web/ui/registry"
	"mljr-web/ui/token"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "user-menu", Name: "User Menu", Category: "special",
		Summary: "Composite: clickable Avatar opens a dropdown with identity header + action items. Composed from primitive.Avatar + custom Datastar dropdown.",
		Code: `// import "mljr-web/ui/special"
special.UserMenu(special.UserMenuProps{
    Name:     "Jane Smith",
    Email:    "jane@example.com",
    Initials: "JS",
    Signal:   "_um1",
    Items: []special.UserMenuItem{
        {Label: "Profile", Icon: "lucide:user",     Href: "/profile"},
        {Label: "Settings", Icon: "lucide:settings", Href: "/settings"},
        {Divider: true},
        {Label: "Sign out", Icon: "lucide:log-out",  Danger: true, OnClick: "alert('bye')"},
    },
})`,
		Render: func(p map[string]string) g.Node {
			return h.Div(
				h.Style("display:flex;gap:var(--sp-8);flex-wrap:wrap;align-items:flex-start;padding:var(--sp-4)"),
				h.Div(
					h.P(h.Style("font-size:var(--t-xs);color:var(--muted);margin-bottom:var(--sp-2)"), g.Text("Initials avatar:")),
					UserMenu(UserMenuProps{
						Name:     "Jane Smith",
						Email:    "jane@example.com",
						Initials: "JS",
						Signal:   "_um1",
						Align:    "left",
						Items: []UserMenuItem{
							{Label: "Profile", Icon: "lucide:user", Href: "#"},
							{Label: "Settings", Icon: "lucide:settings", Href: "#"},
							{Label: "Billing", Icon: "lucide:credit-card", Href: "#"},
							{Divider: true},
							{Label: "Sign out", Icon: "lucide:log-out", Danger: true, OnClick: "alert('Signed out')"},
						},
					}),
				),
				h.Div(
					h.P(h.Style("font-size:var(--t-xs);color:var(--muted);margin-bottom:var(--sp-2)"), g.Text("Large, right-aligned:")),
					UserMenu(UserMenuProps{
						Name:     "Alex Developer",
						Email:    "alex@mljr.eu",
						Initials: "AD",
						Size:     token.SizeLG,
						Signal:   "_um2",
						Align:    "right",
						Items: []UserMenuItem{
							{Label: "Dashboard", Icon: "lucide:layout-dashboard", Href: "#"},
							{Label: "API Keys", Icon: "lucide:key", Href: "#"},
							{Divider: true},
							{Label: "Delete Account", Icon: "lucide:trash-2", Danger: true, OnClick: "alert('Really?')"},
						},
					}),
				),
			)
		},
	})
}
