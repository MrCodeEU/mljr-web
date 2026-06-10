//go:build showcase

package layout

import (
	"mljr-web/ui/registry"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "background", Name: "Background", Category: "layout",
		Summary: "Decorative CSS pattern layer. Dots, grid, lines, diagonal, cross, or gradient. Position absolute, pointer-events none.",
		Code: `// Inside a position:relative container
h.Div(h.Style("position:relative;height:200px;overflow:hidden"),
    layout.Background(layout.BackgroundProps{
        Pattern: layout.BGDots,
        Color:   "var(--line)",
        Opacity: 0.5,
    }),
    // ... your content above the background ...
)`,
		Render: func(p map[string]string) g.Node {
			patterns := []struct {
				name    string
				pattern BackgroundPattern
				color   string
			}{
				{"Dots", BGDots, "var(--primary)"},
				{"Grid", BGGrid, "var(--line)"},
				{"Lines", BGLines, "var(--accent)"},
				{"Diagonal", BGDiagonal, "var(--line)"},
				{"Cross", BGCross, "var(--success)"},
				{"Gradient", BGGradient, "var(--accent)"},
			}
			cells := make([]g.Node, len(patterns))
			for i, pt := range patterns {
				cells[i] = h.Div(
					h.Style("position:relative;height:100px;border:var(--bw-1) solid var(--line);border-radius:var(--radius);overflow:hidden;display:flex;align-items:center;justify-content:center"),
					Background(BackgroundProps{Pattern: pt.pattern, Color: pt.color, Opacity: 0.6}),
					h.Span(h.Style("position:relative;z-index:1;font-size:var(--t-xs);font-weight:700;text-transform:uppercase;letter-spacing:.06em"), g.Text(pt.name)),
				)
			}
			return h.Div(
				h.Style("display:grid;grid-template-columns:repeat(3,1fr);gap:var(--sp-3)"),
				g.Group(cells),
			)
		},
	})
}
