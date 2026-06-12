//go:build showcase

package overlay

import (
	"mljr-web/ui/primitive"
	"mljr-web/ui/registry"
	"mljr-web/ui/token"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "sheet", Name: "Sheet", Category: "overlay",
		PreviewHeight: "460px",
		Summary:       "Full-edge slide-in panel. Bottom placement suits mobile; right placement suits desktop side panels.",
		Code: `// Open trigger
primitive.Button(..., g.Attr("data-on:click", "$_sheetOpen=true"), ...)

// Sheet (place anywhere on page)
overlay.Sheet(overlay.SheetProps{
    Title:     "Filter options",
    Placement: "bottom",
},
    h.P(g.Text("Sheet content here")),
)`,
		Render: func(p map[string]string) g.Node {
			return h.Div(
				h.Style("display:flex;flex-direction:column;gap:var(--sp-4);align-items:flex-start;padding:var(--sp-4)"),
				h.P(h.Style("color:var(--muted);font-size:var(--t-sm)"), g.Text("Click buttons to open sheets from different edges.")),
				h.Div(h.Style("display:flex;gap:var(--sp-3);flex-wrap:wrap"),
					primitive.Button(primitive.ButtonProps{Variant: token.Primary},
						g.Attr("data-on:click", "$_sh1=true"),
						g.Text("Open bottom sheet"),
					),
					primitive.Button(primitive.ButtonProps{Variant: token.Outline},
						g.Attr("data-on:click", "$_sh2=true"),
						g.Text("Open right sheet"),
					),
				),
				Sheet(SheetProps{SignalName: "_sh1", Title: "Filter options", Placement: "bottom"},
					h.Div(h.Style("display:flex;flex-direction:column;gap:var(--sp-4)"),
						h.P(h.Style("color:var(--muted);font-size:var(--t-sm)"), g.Text("Bottom sheets work well for mobile filter panels, share menus, and action lists.")),
						h.Div(h.Style("display:flex;flex-direction:column;gap:var(--sp-3)"),
							h.Label(g.Attr("data-component", "checkbox"),
								h.Input(h.Type("checkbox")), h.Span(g.Attr("data-slot", "box")),
								h.Span(g.Attr("data-slot", "label"), g.Text("In stock only")),
							),
							h.Label(g.Attr("data-component", "checkbox"),
								h.Input(h.Type("checkbox"), g.Attr("checked", "")), h.Span(g.Attr("data-slot", "box")),
								h.Span(g.Attr("data-slot", "label"), g.Text("Free shipping")),
							),
						),
						primitive.Button(primitive.ButtonProps{Variant: token.Primary},
							g.Attr("data-on:click", "$_sh1=false"),
							h.Style("width:100%"),
							g.Text("Apply filters"),
						),
					),
				),
				Sheet(SheetProps{SignalName: "_sh2", Title: "Notifications", Placement: "right"},
					h.Div(h.Style("display:flex;flex-direction:column;gap:var(--sp-3)"),
						func() g.Node {
							items := []struct{ title, time string }{
								{"New comment on your post", "2 min ago"},
								{"Project build succeeded", "12 min ago"},
								{"Your trial expires in 3 days", "1 hr ago"},
							}
							nodes := make([]g.Node, len(items))
							for i, it := range items {
								nodes[i] = h.Div(
									h.Style("padding:var(--sp-3);border-bottom:var(--bw-1) solid var(--line)"),
									h.P(h.Style("font-weight:600;font-size:var(--t-sm);margin:0"), g.Text(it.title)),
									h.Span(h.Style("font-size:var(--t-xs);color:var(--muted)"), g.Text(it.time)),
								)
							}
							return g.Group(nodes)
						}(),
					),
				),
			)
		},
	})
}
