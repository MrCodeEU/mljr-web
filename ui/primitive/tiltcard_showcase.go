//go:build showcase

package primitive

import (
	"mljr-web/ui/registry"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "tilt-card", Name: "Tilt Card", Category: "primitive",
		Summary: "3D perspective tilt driven by pointer position. Optional shine gradient overlay. Zero JS libraries.",
		Code: `primitive.TiltCard(primitive.TiltCardProps{
    MaxTilt:     15,
    Scale:       1.04,
    Perspective: 800,
    Shine:       true,
}, content...)`,
		Render: func(p map[string]string) g.Node {
			card := func(label, color string) g.Node {
				return TiltCard(TiltCardProps{MaxTilt: 15, Scale: 1.04, Perspective: 800, Shine: true},
					h.Div(
						h.Style("background:"+color+";border:var(--bw-2) solid var(--ink);border-radius:var(--radius);padding:var(--sp-6) var(--sp-5);box-shadow:var(--shadow-md);min-height:140px;display:flex;align-items:center;justify-content:center"),
						h.Strong(h.Style("font-size:var(--t-xl);font-family:var(--font-display);font-weight:900"), g.Text(label)),
					),
				)
			}
			return h.Div(
				h.Style("display:grid;grid-template-columns:repeat(3,1fr);gap:var(--sp-5);max-width:600px"),
				card("Hover me", "var(--accent)"),
				card("Tilt!", "var(--surface-2)"),
				card("3D ✨", "var(--ink);color:var(--bg)"),
			)
		},
	})
}
