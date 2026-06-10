//go:build showcase

package primitive

import (
	"mljr-web/ui/icon"
	"mljr-web/ui/registry"
	"mljr-web/ui/token"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "flip-card", Name: "Flip Card", Category: "primitive",
		Summary: "3D CSS flip on hover or Datastar-signal click. Front and back faces are fully custom content.",
		Code: `overlay.FlipCard(overlay.FlipCardProps{Height: "200px"},
    // front
    h.Div(g.Text("Front face")),
    // back
    h.Div(g.Text("Back face")),
)`,
		Render: func(p map[string]string) g.Node {
			face := func(bg, text, label string) g.Node {
				return h.Div(
					h.Style("height:100%;background:"+bg+";border-radius:var(--radius);display:flex;flex-direction:column;align-items:center;justify-content:center;gap:var(--sp-3);padding:var(--sp-5)"),
					icon.Icon(text, icon.Props{Size: "2rem"}),
					h.Div(h.Style("font-weight:800;font-size:var(--t-sm)"), g.Text(label)),
				)
			}
			return h.Div(
				h.Style("display:flex;gap:var(--sp-6);flex-wrap:wrap"),
				h.Div(
					h.P(h.Style("font-size:var(--t-sm);color:var(--muted);margin-bottom:var(--sp-2)"), g.Text("Hover:")),
					FlipCard(FlipCardProps{Height: "200px"},
						face("var(--accent)", "lucide:star", "Hover me"),
						face("var(--surface-2)", "lucide:zap", "Back face!"),
					),
				),
				h.Div(
					h.P(h.Style("font-size:var(--t-sm);color:var(--muted);margin-bottom:var(--sp-2)"), g.Text("Click:")),
					FlipCard(FlipCardProps{Height: "200px", Trigger: "click", Signal: "_fc2"},
						face("var(--surface-2)", "lucide:rotate-cw", "Click me"),
						face("var(--accent)", "lucide:check", "Flipped!"),
					),
				),
				h.Div(
					h.P(h.Style("font-size:var(--t-sm);color:var(--muted);margin-bottom:var(--sp-2)"), g.Text("Profile card:")),
					FlipCard(FlipCardProps{Height: "200px"},
						h.Div(
							h.Style("height:100%;background:var(--surface-2);border-radius:var(--radius);display:flex;flex-direction:column;align-items:center;justify-content:center;gap:var(--sp-3);padding:var(--sp-5)"),
							Avatar(AvatarProps{Initials: "JS", Size: token.SizeLG, Tone: token.ToneAccent}),
							h.Div(h.Style("font-weight:800"), g.Text("Jane Smith")),
							h.Div(h.Style("font-size:var(--t-xs);color:var(--muted)"), g.Text("Senior Engineer")),
						),
						h.Div(
							h.Style("height:100%;background:var(--ink);color:var(--surface);border-radius:var(--radius);display:flex;flex-direction:column;justify-content:center;gap:var(--sp-2);padding:var(--sp-5)"),
							h.Div(h.Style("font-size:var(--t-sm)"), icon.Icon("lucide:mail", icon.Props{Size: "1rem"}), g.Text(" jane@example.com")),
							h.Div(h.Style("font-size:var(--t-sm)"), icon.Icon("lucide:github", icon.Props{Size: "1rem"}), g.Text(" @janesmith")),
							h.Div(h.Style("font-size:var(--t-sm)"), icon.Icon("lucide:map-pin", icon.Props{Size: "1rem"}), g.Text(" Berlin, DE")),
						),
					),
				),
			)
		},
	})
}
