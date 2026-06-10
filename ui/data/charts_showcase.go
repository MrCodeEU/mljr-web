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
		Slug: "line-chart", Name: "Line Chart", Category: "data",
		Summary: "Pure-SVG multi-series line/area chart rendered server-side. Smooth bezier curves, optional grid, dots, and fill.",
		Code: `data.LineChart(data.LineChartProps{
    Height:   160,
    ShowDots: true,
    ShowGrid: true,
    Labels:   []string{"Jan","Feb","Mar","Apr","May"},
    Series: []data.LineChartSeries{
        {Label: "Revenue", Points: []float64{42,67,91,55,78}, Fill: true},
        {Label: "Costs",   Points: []float64{30,35,40,38,42}, Color: "var(--danger)"},
    },
})`,
		Render: func(p map[string]string) g.Node {
			months := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug"}
			return h.Div(
				h.Style("display:flex;flex-direction:column;gap:var(--sp-5)"),
				primitive.Card(primitive.CardProps{},
					h.P(h.Style("font-size:var(--t-xs);text-transform:uppercase;font-weight:700;letter-spacing:.06em;opacity:.5;margin:0 0 var(--sp-3)"), g.Text("Revenue vs Costs")),
					LineChart(LineChartProps{
						Height: 160, ShowDots: true, ShowGrid: true, Labels: months,
						Series: []LineChartSeries{
							{Label: "Revenue", Points: []float64{42, 67, 91, 55, 78, 110, 95, 83}, Fill: true},
							{Label: "Costs", Points: []float64{30, 35, 40, 38, 42, 48, 45, 44}, Color: "var(--danger)"},
						},
					}),
				),
				primitive.Card(primitive.CardProps{},
					h.P(h.Style("font-size:var(--t-xs);text-transform:uppercase;font-weight:700;letter-spacing:.06em;opacity:.5;margin:0 0 var(--sp-3)"), g.Text("Single series with fill")),
					LineChart(LineChartProps{
						Height: 100, ShowGrid: true, Labels: months,
						Series: []LineChartSeries{
							{Label: "Users", Points: []float64{120, 145, 160, 138, 175, 210, 195, 220}, Fill: true, Color: "var(--accent)"},
						},
					}),
				),
			)
		},
	})

	registry.Register(&registry.Component{
		Slug: "donut-chart", Name: "Donut Chart", Category: "data",
		Summary: "Pure-SVG donut chart using stroke-dasharray. Includes legend, center label, and hover tooltips.",
		Code: `data.DonutChart(data.DonutChartProps{
    Label:    "92%",
    Sublabel: "uptime",
    Slices: []data.DonutSlice{
        {Label: "Go",         Value: 45},
        {Label: "TypeScript", Value: 30},
        {Label: "Python",     Value: 25},
    },
})`,
		Render: func(p map[string]string) g.Node {
			return h.Div(
				h.Style("display:flex;flex-wrap:wrap;gap:var(--sp-5)"),
				primitive.Card(primitive.CardProps{},
					h.P(h.Style("font-size:var(--t-xs);text-transform:uppercase;font-weight:700;letter-spacing:.06em;opacity:.5;margin:0 0 var(--sp-3)"), g.Text("Tech stack")),
					DonutChart(DonutChartProps{
						Label: "5", Sublabel: "languages",
						Slices: []DonutSlice{
							{Label: "Go", Value: 45, Color: "#00ADD8"},
							{Label: "TypeScript", Value: 30, Color: "#3178C6"},
							{Label: "Python", Value: 15, Color: "#3776AB"},
							{Label: "Rust", Value: 8, Color: "#CE422B"},
							{Label: "Other", Value: 2},
						},
					}),
				),
				primitive.Card(primitive.CardProps{},
					h.P(h.Style("font-size:var(--t-xs);text-transform:uppercase;font-weight:700;letter-spacing:.06em;opacity:.5;margin:0 0 var(--sp-3)"), g.Text("Status")),
					DonutChart(DonutChartProps{
						Label: "92%", Sublabel: "uptime", Size: 140, Thickness: 28,
						Slices: []DonutSlice{
							{Label: "Healthy", Value: 92, Color: "var(--success)"},
							{Label: "Degraded", Value: 8, Color: "var(--danger)"},
						},
					}),
				),
			)
		},
	})

	registry.Register(&registry.Component{
		Slug: "sparkline", Name: "Sparkline", Category: "data",
		Summary: "Tiny inline SVG trend line for use inside stat cards, tables, or dashboard tiles.",
		Code: `data.Sparkline(data.SparklineProps{
    Points: []float64{12, 18, 9, 24, 17, 30, 28},
    Fill:   true,
    Color:  "var(--success)",
})`,
		Render: func(p map[string]string) g.Node {
			rows := []struct {
				label  string
				val    string
				delta  string
				color  string
				points []float64
			}{
				{"Revenue", "€24,580", "+12.4%", "var(--success)", []float64{18, 22, 19, 28, 24, 30, 28}},
				{"Users", "9,420", "+8.1%", "var(--primary)", []float64{60, 72, 65, 80, 74, 88, 92}},
				{"Churn", "2.1%", "-0.4%", "var(--danger)", []float64{3.1, 2.8, 3.2, 2.6, 2.4, 2.2, 2.1}},
				{"MRR", "€8,200", "+5.3%", "var(--accent)", []float64{72, 74, 71, 76, 77, 80, 82}},
			}
			rowNodes := make([]g.Node, len(rows))
			for i, r := range rows {
				rowNodes[i] = h.Tr(
					h.Td(h.Style("padding:var(--sp-2) var(--sp-3);font-weight:600;font-size:var(--t-sm)"), g.Text(r.label)),
					h.Td(h.Style("padding:var(--sp-2) var(--sp-3);font-family:var(--font-display);font-weight:900"), g.Text(r.val)),
					h.Td(h.Style("padding:var(--sp-2) var(--sp-3);color:"+r.color+";font-size:var(--t-sm);font-weight:700"), g.Text(r.delta)),
					h.Td(h.Style("padding:var(--sp-2) var(--sp-3)"),
						Sparkline(SparklineProps{Points: r.points, Fill: true, Color: r.color}),
					),
				)
			}
			return h.Table(
				h.Style("width:100%;border-collapse:collapse"),
				h.THead(h.Tr(
					h.Th(h.Style("text-align:left;padding:var(--sp-2) var(--sp-3);font-size:var(--t-xs);opacity:.5;font-weight:700;border-bottom:var(--bw-1) solid var(--line)"), g.Text("Metric")),
					h.Th(h.Style("text-align:left;padding:var(--sp-2) var(--sp-3);font-size:var(--t-xs);opacity:.5;font-weight:700;border-bottom:var(--bw-1) solid var(--line)"), g.Text("Value")),
					h.Th(h.Style("text-align:left;padding:var(--sp-2) var(--sp-3);font-size:var(--t-xs);opacity:.5;font-weight:700;border-bottom:var(--bw-1) solid var(--line)"), g.Text("Δ")),
					h.Th(h.Style("text-align:left;padding:var(--sp-2) var(--sp-3);font-size:var(--t-xs);opacity:.5;font-weight:700;border-bottom:var(--bw-1) solid var(--line)"), g.Text("Trend")),
				)),
				h.TBody(g.Group(rowNodes)),
			)
		},
	})
}
