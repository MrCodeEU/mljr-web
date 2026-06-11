//go:build showcase

package data

import (
	"math/rand"
	"mljr-web/ui/registry"
	"time"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "heatmap", Name: "GitHub Heatmap", Category: "data",
		Summary: "GitHub-style contribution heatmap. SVG rendered server-side from date/value pairs. Zero JS, zero dependencies.",
		Code: `data.Heatmap(data.HeatmapProps{
    Weeks: 52, ShowMonthLabels: true, ShowDayLabels: true,
}, cells)`,
		Render: func(p map[string]string) g.Node {
			now := time.Now()
			cells := make([]HeatmapCell, 365)
			rng := rand.New(rand.NewSource(42))
			for i := range cells {
				d := now.AddDate(0, 0, -(364 - i))
				v := 0
				if rng.Float64() > 0.55 {
					v = rng.Intn(10) + 1
				}
				cells[i] = HeatmapCell{Date: d, Value: v}
			}
			return h.Div(
				h.Style("display:flex;flex-direction:column;gap:var(--sp-3)"),
				Heatmap(HeatmapProps{
					Weeks:           52,
					ShowMonthLabels: true,
					ShowDayLabels:   true,
				}, cells),
				h.P(h.Style("font-size:var(--t-xs);color:var(--muted);margin:0"),
					g.Text("365 days of random activity data. Color intensity scales from 0 → max value.")),
			)
		},
	})
}
