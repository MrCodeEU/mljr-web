//go:build showcase

package primitive

import (
	"mljr-web/ui/icon"
	"mljr-web/ui/registry"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "marquee", Name: "Marquee", Category: "primitive",
		Summary: "Infinite horizontal scroll strip. CSS-only animation — items duplicated server-side for seamless loop. Optional pause-on-hover.",
		Code: `// import "mljr-web/ui/primitive"
primitive.Marquee(primitive.MarqueeProps{
    Speed:        "20s",
    PauseOnHover: true,
},
    h.Div(g.Text("Item 1")),
    h.Div(g.Text("Item 2")),
)`,
		Render: func(p map[string]string) g.Node {
			logos := []string{
				"simple-icons:go", "simple-icons:typescript", "simple-icons:rust",
				"simple-icons:python", "simple-icons:docker", "simple-icons:kubernetes",
				"simple-icons:postgresql", "simple-icons:redis", "simple-icons:nginx",
				"simple-icons:linux",
			}
			logoItems := make([]g.Node, len(logos))
			for i, ic := range logos {
				logoItems[i] = h.Div(
					h.Style("display:flex;align-items:center;gap:var(--sp-2);padding:var(--sp-2) var(--sp-4);background:var(--surface-2);border:var(--bw-1) solid var(--line);border-radius:var(--radius);white-space:nowrap;font-size:var(--t-sm);font-weight:600"),
					icon.Icon(ic, icon.Props{Size: "1.2rem"}),
				)
			}

			words := []string{"Fast", "Typesafe", "Composable", "Themeable", "Accessible", "Minimal"}
			wordItems := make([]g.Node, len(words))
			for i, w := range words {
				wordItems[i] = h.Div(
					h.Style("padding:var(--sp-2) var(--sp-5);background:var(--accent);color:var(--accent-ink);border-radius:var(--radius);font-size:var(--t-sm);font-weight:800;white-space:nowrap"),
					g.Text(w),
				)
			}

			return h.Div(
				h.Style("display:flex;flex-direction:column;gap:var(--sp-5)"),
				h.P(h.Style("font-size:var(--t-sm);color:var(--muted)"), g.Text("Hover to pause.")),
				Marquee(MarqueeProps{Speed: "22s", PauseOnHover: true}, logoItems...),
				Marquee(MarqueeProps{Speed: "18s", Direction: "right", PauseOnHover: true}, wordItems...),
			)
		},
	})
}
