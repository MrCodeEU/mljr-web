//go:build showcase

package primitive

import (
	"fmt"
	"mljr-web/ui/registry"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "scroll-area", Name: "Scroll Area", Category: "primitive",
		Summary: "Styled scroll container with thin themed scrollbars via CSS scrollbar-width + scrollbar-color.",
		Code: `// import "mljr-web/ui/primitive"
primitive.ScrollArea(primitive.ScrollAreaProps{
    Height: "240px",
}, content...)

// Horizontal
primitive.ScrollArea(primitive.ScrollAreaProps{
    Direction: "horizontal",
    Height:    "80px",
}, wideContent...)`,
		Render: func(p map[string]string) g.Node {
			// Vertical content
			vertItems := make([]g.Node, 20)
			for i := range vertItems {
				vertItems[i] = h.Div(
					h.Style("padding:var(--sp-3) var(--sp-4);border-bottom:var(--bw-1) solid var(--line);font-size:var(--t-sm)"),
					g.Text(fmt.Sprintf("Row %d — scroll to see more items below", i+1)),
				)
			}

			// Horizontal content
			horzItems := make([]g.Node, 12)
			for i := range horzItems {
				horzItems[i] = h.Div(
					h.Style("flex-shrink:0;width:120px;height:60px;display:flex;align-items:center;justify-content:center;background:var(--surface-2);border:var(--bw-1) solid var(--line);border-radius:var(--radius);font-size:var(--t-sm);font-weight:600"),
					g.Text(fmt.Sprintf("Card %d", i+1)),
				)
			}

			return h.Div(
				h.Style("display:flex;flex-direction:column;gap:var(--sp-6)"),
				h.Div(
					h.H4(h.Style("font-size:var(--t-sm);font-weight:700;margin-bottom:var(--sp-2)"), g.Text("Vertical")),
					ScrollArea(ScrollAreaProps{Height: "200px"},
						h.Div(g.Group(vertItems)),
					),
				),
				h.Div(
					h.H4(h.Style("font-size:var(--t-sm);font-weight:700;margin-bottom:var(--sp-2)"), g.Text("Horizontal")),
					ScrollArea(ScrollAreaProps{Direction: "horizontal", Height: "80px"},
						h.Div(
							h.Style("display:flex;gap:var(--sp-3);padding:var(--sp-2)"),
							g.Group(horzItems),
						),
					),
				),
			)
		},
	})
}
