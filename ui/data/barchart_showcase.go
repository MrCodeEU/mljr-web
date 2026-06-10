//go:build showcase

package data

import (
	"mljr-web/ui/primitive"
	"mljr-web/ui/registry"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "bar-chart", Name: "Bar Chart", Category: "data",
		Summary: "Pure SVG bar chart rendered server-side. No JS, no canvas. Accepts any data slice.",
		Code: `data.BarChart(data.BarChartProps{
    Height:     180,
    ShowValues: true,
    ShowGrid:   true,
    Caption:    "Monthly deployments",
    Data: []data.BarDatum{
        {Label: "Jan", Value: 42},
        {Label: "Feb", Value: 67, Color: "var(--accent)"},
        {Label: "Mar", Value: 91},
    },
})`,
		Render: func(p map[string]string) g.Node {
			months := []BarDatum{
				{Label: "Jan", Value: 42},
				{Label: "Feb", Value: 67},
				{Label: "Mar", Value: 91},
				{Label: "Apr", Value: 55},
				{Label: "May", Value: 78},
				{Label: "Jun", Value: 110},
				{Label: "Jul", Value: 95},
				{Label: "Aug", Value: 83},
			}
			skills := []BarDatum{
				{Label: "Go", Value: 95, Color: "#00ADD8"},
				{Label: "Rust", Value: 72, Color: "#CE422B"},
				{Label: "TS", Value: 80, Color: "#3178C6"},
				{Label: "Python", Value: 65, Color: "#3776AB"},
				{Label: "SQL", Value: 88, Color: "var(--success)"},
			}
			return h.Div(
				h.Style("display:flex;flex-direction:column;gap:var(--sp-5)"),
				primitive.Card(primitive.CardProps{},
					h.P(h.Style("font-size:var(--t-xs);text-transform:uppercase;font-weight:700;letter-spacing:.06em;opacity:.5;margin:0 0 var(--sp-3)"), g.Text("Monthly deployments")),
					BarChart(BarChartProps{Height: 160, ShowValues: true, ShowGrid: true, Data: months}),
				),
				primitive.Card(primitive.CardProps{},
					h.P(h.Style("font-size:var(--t-xs);text-transform:uppercase;font-weight:700;letter-spacing:.06em;opacity:.5;margin:0 0 var(--sp-3)"), g.Text("Skill levels")),
					BarChart(BarChartProps{Height: 120, ShowValues: true, Data: skills}),
				),
			)
		},
	})
}
