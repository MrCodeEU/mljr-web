//go:build showcase

package data

import (
	"mljr-web/ui/registry"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "gauge", Name: "Gauge / Meter", Category: "data",
		Summary: "SVG circular gauge with 270° arc. Optional tick marks. Server-side rendered — zero JS.",
		Code: `data.Gauge(data.GaugeProps{
    Value: 72, Max: 100, Label: "CPU", Unit: "%",
    Ticks: true, Size: 200,
})`,
		Render: func(p map[string]string) g.Node {
			return h.Div(
				h.Style("display:grid;grid-template-columns:repeat(3,1fr);gap:var(--sp-4);justify-items:center"),
				Gauge(GaugeProps{Value: 72, Max: 100, Label: "CPU", Unit: "%", Ticks: true, Size: 180}),
				Gauge(GaugeProps{Value: 3.8, Max: 5, Label: "Rating", Unit: "★", Color: "var(--warning)", Size: 180}),
				Gauge(GaugeProps{Value: 24, Max: 100, Label: "Disk", Unit: "%", Color: "var(--success)", Ticks: true, Size: 180}),
			)
		},
	})
}
