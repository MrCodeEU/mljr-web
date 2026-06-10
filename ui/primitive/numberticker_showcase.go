//go:build showcase

package primitive

import (
	"mljr-web/ui/registry"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "number-ticker", Name: "Number Ticker", Category: "primitive",
		Summary: "Animated number counter using requestAnimationFrame with ease-out cubic easing. Optional IntersectionObserver trigger.",
		Code: `// import "mljr-web/ui/primitive"
primitive.NumberTicker(primitive.NumberTickerProps{
    Value:         12500,
    Prefix:        "$",
    TriggerOnView: true,
    ID:            "revenue",
})

// Percentage
primitive.NumberTicker(primitive.NumberTickerProps{
    Value:    98.6,
    Suffix:   "%",
    Decimals: 1,
    Duration: 2000,
    ID:       "uptime",
})`,
		Render: func(p map[string]string) g.Node {
			return h.Div(
				h.Style("display:flex;flex-direction:column;gap:var(--sp-6)"),
				h.Div(
					h.Style("display:grid;grid-template-columns:repeat(auto-fit,minmax(140px,1fr));gap:var(--sp-4)"),
					statTicker("Revenue", NumberTickerProps{Value: 128500, Prefix: "$", TriggerOnView: true, ID: "nt-rev", Duration: 1500}),
					statTicker("Users", NumberTickerProps{Value: 42317, TriggerOnView: true, ID: "nt-usr", Duration: 1800}),
					statTicker("Uptime", NumberTickerProps{Value: 99.97, Suffix: "%", Decimals: 2, TriggerOnView: true, ID: "nt-up", Duration: 1200}),
					statTicker("Speed", NumberTickerProps{Value: 248, Suffix: " ms", TriggerOnView: true, ID: "nt-spd", Duration: 1000}),
				),
				h.P(h.Style("font-size:var(--t-xs);color:var(--muted)"), g.Text("Counters animate when scrolled into view. Reload to replay.")),
			)
		},
	})
}

func statTicker(label string, p NumberTickerProps) g.Node {
	return h.Div(
		h.Style("padding:var(--sp-4);background:var(--surface-2);border:var(--bw-1) solid var(--line);border-radius:var(--radius);text-align:center"),
		h.Div(
			h.Style("font-size:var(--t-2xl);font-weight:800;line-height:1.1"),
			NumberTicker(p),
		),
		h.Div(h.Style("font-size:var(--t-xs);color:var(--muted);margin-top:var(--sp-1)"), g.Text(label)),
	)
}
