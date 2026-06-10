//go:build showcase

package primitive

import (
	"mljr-web/ui/registry"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "gradient-text", Name: "Gradient Text", Category: "primitive",
		Summary: "CSS background-clip:text gradient fill. Theme-aware by default (accent → ink). Any HTML tag, any gradient angle.",
		Code: `primitive.GradientText(primitive.GradientTextProps{
    From: "var(--accent)",
    To:   "var(--ink)",
    Tag:  "h1",
}, g.Text("Beautiful gradients"))`,
		Render: func(p map[string]string) g.Node {
			return h.Div(
				h.Style("display:flex;flex-direction:column;gap:var(--sp-5)"),
				h.Div(
					h.Style("font-size:var(--t-3xl);font-weight:900;line-height:1.1"),
					GradientText(GradientTextProps{From: "var(--accent)", To: "var(--ink)"}, g.Text("Ship faster.")),
				),
				h.Div(
					h.Style("font-size:var(--t-2xl);font-weight:800"),
					GradientText(GradientTextProps{From: "#f59e0b", To: "#ef4444", Angle: "90deg"}, g.Text("Fire gradient")),
				),
				h.Div(
					h.Style("font-size:var(--t-xl);font-weight:800"),
					GradientText(GradientTextProps{From: "#6366f1", Via: "#8b5cf6", To: "#ec4899", Angle: "135deg"}, g.Text("Purple to pink via violet")),
				),
				h.Div(
					h.Style("font-size:var(--t-lg);font-weight:700"),
					GradientText(GradientTextProps{From: "#10b981", To: "#3b82f6", Angle: "45deg"}, g.Text("Green to blue")),
				),
				h.P(
					h.Style("font-size:var(--t-base)"),
					g.Text("Inline: build "),
					GradientText(GradientTextProps{From: "var(--accent)", To: "var(--ink)"}, g.Text("amazing")),
					g.Text(" things today."),
				),
			)
		},
	})
}
