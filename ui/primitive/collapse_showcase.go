//go:build showcase

package primitive

import (
	"mljr-web/ui/registry"
	"mljr-web/ui/token"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "collapse", Name: "Collapse", Category: "primitive",
		Summary: "Animated show/hide region driven by a Datastar signal. Max-height CSS transition, no JS.",
		Code: `// Toggle button
primitive.Button(primitive.ButtonProps{},
    g.Attr("data-on:click", "$_collapseOpen=!$_collapseOpen"),
    g.Text("Toggle"),
)

// Collapse region (signals declared inside)
primitive.Collapse(primitive.CollapseProps{SignalName: "_collapseOpen"},
    h.P(g.Text("Hidden content...")),
)`,
		Render: func(p map[string]string) g.Node {
			return h.Div(
				h.Style("display:flex;flex-direction:column;gap:var(--sp-4)"),
				h.Div(
					h.Style("display:flex;flex-direction:column;gap:var(--sp-2)"),
					h.Div(
						h.Style("display:flex;align-items:center;justify-content:space-between;padding:var(--sp-3);background:var(--surface-2);border:var(--bw-1) solid var(--line);border-radius:var(--radius);cursor:pointer"),
						g.Attr("data-on:click", "$_c1=!$_c1"),
						g.Attr("data-signals", `{"_c1":false}`),
						h.Span(h.Style("font-weight:700"), g.Text("FAQ: How does this work?")),
						h.Span(g.Attr("data-text", `$_c1?"▲":"▼"`), g.Text("▼")),
					),
					Collapse(CollapseProps{SignalName: "_c1"},
						h.Div(
							h.Style("padding:var(--sp-4);background:var(--surface);border:var(--bw-1) solid var(--line);border-top:none;border-radius:0 0 var(--radius) var(--radius)"),
							h.P(h.Style("margin:0;color:var(--muted)"), g.Text("Collapse uses a max-height CSS transition driven by a Datastar signal. No JavaScript needed for the animation — just CSS and reactive attributes.")),
						),
					),
				),
				h.Div(
					h.Style("display:flex;flex-direction:column;gap:var(--sp-2)"),
					Button(ButtonProps{Variant: token.Outline},
						g.Attr("data-on:click", "$_c2=!$_c2"),
						g.Text("Toggle advanced settings"),
					),
					Collapse(CollapseProps{SignalName: "_c2"},
						h.Div(
							h.Style("padding:var(--sp-4);background:var(--surface-2);border:var(--bw-2) dashed var(--line);border-radius:var(--radius);display:flex;flex-direction:column;gap:var(--sp-3)"),
							h.P(h.Style("font-weight:700;margin:0"), g.Text("Advanced settings")),
							h.P(h.Style("color:var(--muted);font-size:var(--t-sm);margin:0"), g.Text("These settings are hidden by default and revealed with a smooth animation. Works with any content — forms, images, long text.")),
						),
					),
				),
			)
		},
	})
}
