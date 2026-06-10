//go:build showcase

package data

import (
	"mljr-web/ui/layout"
	"mljr-web/ui/registry"

	g "maragu.dev/gomponents"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "stat-card", Name: "Stat Card", Category: "data",
		Summary: "Key metric display: label, large value, optional delta with up/down coloring.",
		Code: `data.StatCard(data.StatCardProps{
    Label:   "Monthly revenue",
    Value:   "$48,200",
    Delta:   "+12.4%",
    DeltaUp: true,
})`,
		Controls: []registry.Control{
			{Name: "delta", Type: registry.ControlBool, Default: "true"},
			{Name: "deltaup", Type: registry.ControlBool, Default: "true"},
		},
		Render: func(p map[string]string) g.Node {
			delta := ""
			if p["delta"] == "true" {
				if p["deltaup"] == "true" {
					delta = "+12.4%"
				} else {
					delta = "−8.3%"
				}
			}
			return layout.Grid(layout.GridProps{},
				layout.Col(layout.ColProps{Span: 4},
					StatCard(StatCardProps{Label: "Revenue", Value: "$48,200", Delta: delta, DeltaUp: p["deltaup"] == "true"}),
				),
				layout.Col(layout.ColProps{Span: 4},
					StatCard(StatCardProps{Label: "Users", Value: "12,841", Delta: delta, DeltaUp: p["deltaup"] == "true"}),
				),
				layout.Col(layout.ColProps{Span: 4},
					StatCard(StatCardProps{Label: "Conversion", Value: "3.6%", Delta: delta, DeltaUp: p["deltaup"] == "true"}),
				),
			)
		},
	})
}
