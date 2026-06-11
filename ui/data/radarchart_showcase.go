//go:build showcase

package data

import (
	"mljr-web/ui/registry"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "radar-chart", Name: "Radar Chart", Category: "data",
		Summary: "Multi-series radar/spider chart rendered as inline SVG. Server-side Go math — zero JS, zero canvas.",
		Code: `data.RadarChart(data.RadarChartProps{
    Axes:       []string{"Go", "Rust", "TypeScript", "Python", "SQL"},
    ShowGrid:   true,
    GridLevels: 5,
},
    data.RadarSeries{Label: "Alice", Values: []float64{90, 60, 75, 50, 85}},
    data.RadarSeries{Label: "Bob",   Values: []float64{55, 80, 70, 90, 60}},
)`,
		Render: func(p map[string]string) g.Node {
			return h.Div(
				h.Style("display:grid;grid-template-columns:1fr 1fr;gap:var(--sp-6)"),
				h.Div(
					h.P(h.Style("font-weight:800;margin:0 0 var(--sp-3)"), g.Text("Tech Skills")),
					RadarChart(RadarChartProps{
						Axes:       []string{"Go", "Rust", "TypeScript", "Python", "SQL"},
						ShowGrid:   true,
						GridLevels: 5,
						Size:       260,
					},
						RadarSeries{Label: "Alice", Values: []float64{90, 60, 75, 50, 85}},
						RadarSeries{Label: "Bob", Values: []float64{55, 80, 70, 90, 60}},
					),
				),
				h.Div(
					h.P(h.Style("font-weight:800;margin:0 0 var(--sp-3)"), g.Text("Soft Skills")),
					RadarChart(RadarChartProps{
						Axes:       []string{"Leadership", "Comms", "Teamwork", "Initiative", "Delivery"},
						ShowGrid:   true,
						GridLevels: 4,
						Size:       260,
					},
						RadarSeries{Label: "Score", Values: []float64{80, 70, 90, 75, 85}},
					),
				),
			)
		},
	})
}
