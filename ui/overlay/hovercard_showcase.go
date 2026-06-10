//go:build showcase

package overlay

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
		Slug: "hover-card", Name: "Hover Card", Category: "overlay",
		Summary: "Rich content card triggered by hovering. Pure CSS — no JS, no Datastar. Ideal for user profile previews, link previews.",
		Code: `// import "mljr-web/ui/overlay"
overlay.HoverCard(
    overlay.HoverCardProps{Placement: "top"},
    h.Span(g.Text("@username")),
    h.Div( /* card content */ ),
)`,
		Render: func(p map[string]string) g.Node {
			profileCard := h.Div(
				primitive.Avatar(primitive.AvatarProps{
					Initials: "JS",
					Tone:     token.ToneAccent,
					Size:     token.SizeLG,
				}),
				h.Div(
					h.Style("margin-top:var(--sp-3)"),
					h.Div(h.Style("font-weight:800;font-size:var(--t-base)"), g.Text("Jane Smith")),
					h.Div(h.Style("color:var(--muted);font-size:var(--t-sm)"), g.Text("@janesmith · Senior Engineer")),
					h.Div(
						h.Style("margin-top:var(--sp-2);font-size:var(--t-sm)"),
						g.Text("Building UI components for the web. Go + HTMX enthusiast."),
					),
					h.Div(
						h.Style("display:flex;gap:var(--sp-4);margin-top:var(--sp-3);font-size:var(--t-sm)"),
						h.Span(h.Strong(g.Text("142")), g.Text(" Following")),
						h.Span(h.Strong(g.Text("1.2K")), g.Text(" Followers")),
					),
				),
			)

			linkCard := h.Div(
				h.Div(h.Style("font-weight:700;font-size:var(--t-sm);margin-bottom:var(--sp-1)"), g.Text("mljr-web — Component Library")),
				h.Div(h.Style("font-size:var(--t-xs);color:var(--muted);margin-bottom:var(--sp-2)"), g.Text("Go + gomponents + Datastar 1.0.2 + Tailwind v4")),
				h.Div(
					h.Style("display:flex;gap:var(--sp-3);font-size:var(--t-xs)"),
					h.Span(icon.Icon("lucide:star", icon.Props{Size: "0.85rem"}), g.Text(" 247")),
					h.Span(icon.Icon("lucide:git-branch", icon.Props{Size: "0.85rem"}), g.Text(" 18")),
				),
			)

			return h.Div(
				h.Style("display:flex;gap:var(--sp-8);flex-wrap:wrap;padding:var(--sp-8) var(--sp-4)"),
				h.Div(
					h.P(h.Style("font-size:var(--t-sm);color:var(--muted);margin-bottom:var(--sp-2)"), g.Text("Profile card:")),
					HoverCard(
						HoverCardProps{Placement: "top"},
						h.A(h.Href("#"), h.Style("font-weight:700;color:var(--accent);text-decoration:underline"), g.Text("@janesmith")),
						profileCard,
					),
				),
				h.Div(
					h.P(h.Style("font-size:var(--t-sm);color:var(--muted);margin-bottom:var(--sp-2)"), g.Text("Link preview:")),
					HoverCard(
						HoverCardProps{Placement: "bottom", Width: "300px"},
						h.A(h.Href("#"), h.Style("font-weight:700;color:var(--accent);text-decoration:underline"), g.Text("mljr-web")),
						linkCard,
					),
				),
			)
		},
	})
}
