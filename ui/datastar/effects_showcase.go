//go:build showcase

package datastar

import (
	"mljr-web/ui"
	"mljr-web/ui/primitive"
	"mljr-web/ui/registry"
	"mljr-web/ui/token"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "ds-effects", Name: "Effects & Refs", Category: "datastar",
		Summary: "data-effect runs side-effects when signals change. data-ref captures element references. data-init fires once on mount.",
		Code: `// data-effect: runs whenever accessed signals change
// (like useEffect but reactive, no dependency array)
data-effect="document.title = 'Cart: '+$items+' items'"
data-effect="localStorage.setItem('theme', $theme)"

// data-ref: capture a DOM element reference as a signal
data-ref:canvas     // $canvas now points to this element
// use in another expression:
data-effect="$canvas.getContext('2d').clearRect(0,0,$canvas.width,$canvas.height)"

// data-init: run once when element mounts (not reactive)
data-init="$score = parseInt(localStorage.getItem('score')||'0')"`,
		Render: func(p map[string]string) g.Node {
			return h.Div(
				ui.Signals(`{"persist":0,"scrollY":0,"refW":0,"refH":0}`),
				h.Div(h.Style("display:flex;flex-direction:column;gap:var(--sp-5)"),

					// data-effect: localStorage persist
					primitive.Card(primitive.CardProps{},
						h.P(h.Style("font-size:var(--t-xs);text-transform:uppercase;font-weight:700;letter-spacing:.06em;opacity:.5;margin-bottom:var(--sp-3)"), g.Text("data-effect — localStorage persist")),
						// init restores from localStorage
						h.Div(g.Attr("data-init", "var v=localStorage.getItem('ds-demo-count'); if(v)$persist=parseInt(v)")),
						// effect writes to localStorage whenever $persist changes
						h.Div(g.Attr("data-effect", "localStorage.setItem('ds-demo-count',$persist)")),
						h.Div(h.Style("display:flex;align-items:center;gap:var(--sp-3)"),
							primitive.Button(primitive.ButtonProps{Variant: token.Outline, Size: token.SizeSM},
								g.Attr("data-on:click", "$persist=Math.max(0,$persist-1)"), g.Text("−"),
							),
							h.Span(h.Style("font-size:var(--t-xl);font-weight:900;font-family:var(--font-display);min-width:3rem;text-align:center"),
								g.Attr("data-text", "$persist"),
							),
							primitive.Button(primitive.ButtonProps{Variant: token.Primary, Size: token.SizeSM},
								g.Attr("data-on:click", "$persist++"), g.Text("+"),
							),
						),
						h.P(h.Style("font-size:var(--t-xs);opacity:.5;margin-top:var(--sp-2)"),
							g.Text("Value persists across page refreshes — stored in localStorage via data-effect."),
						),
					),

					// data-ref + data-effect: measure element
					primitive.Card(primitive.CardProps{},
						h.P(h.Style("font-size:var(--t-xs);text-transform:uppercase;font-weight:700;letter-spacing:.06em;opacity:.5;margin-bottom:var(--sp-3)"), g.Text("data-ref + data-effect — measure an element")),
						h.Div(
							g.Attr("data-ref:measured"),
							// effect reads offsetWidth/offsetHeight from the ref
							g.Attr("data-effect", "$refW=$measured?$measured.offsetWidth:0; $refH=$measured?$measured.offsetHeight:0"),
							h.Style("padding:var(--sp-5);background:var(--surface-2);border:var(--bw-1) dashed var(--line);border-radius:var(--radius);resize:both;overflow:auto;min-width:200px"),
							h.P(h.Style("font-size:var(--t-sm);text-align:center;pointer-events:none"),
								g.Text("Resize me ↘"),
							),
						),
						h.Div(h.Style("margin-top:var(--sp-3);font-size:var(--t-sm);display:flex;gap:var(--sp-4)"),
							h.Span(g.Text("W: "), h.Strong(g.Attr("data-text", "Math.round($refW)+'px'"))),
							h.Span(g.Text("H: "), h.Strong(g.Attr("data-text", "Math.round($refH)+'px'"))),
						),
					),

					// data-effect: reactive document title
					primitive.Card(primitive.CardProps{},
						h.P(h.Style("font-size:var(--t-xs);text-transform:uppercase;font-weight:700;letter-spacing:.06em;opacity:.5;margin-bottom:var(--sp-3)"), g.Text("data-effect — reactive document.title")),
						h.Div(g.Attr("data-effect", "document.title = $persist>0 ? '('+$persist+') mljr-ui showcase' : 'mljr-ui showcase'")),
						h.P(h.Style("font-size:var(--t-sm)"),
							g.Text("Increment the counter above — watch the browser tab title update reactively."),
						),
						h.P(h.Style("font-size:var(--t-xs);opacity:.5;margin-top:var(--sp-1)"),
							g.Text("data-effect runs whenever any signal it reads ($persist) changes."),
						),
					),

					// data-init: one-time setup
					primitive.Card(primitive.CardProps{},
						h.P(h.Style("font-size:var(--t-xs);text-transform:uppercase;font-weight:700;letter-spacing:.06em;opacity:.5;margin-bottom:var(--sp-3)"), g.Text("data-init vs data-effect")),
						h.Div(h.Style("display:flex;flex-direction:column;gap:var(--sp-2);font-size:var(--t-sm)"),
							h.Div(h.Style("display:flex;gap:var(--sp-2)"),
								primitive.Badge(primitive.BadgeProps{Variant: primitive.BadgeInfo}, g.Text("data-init")),
								h.Span(g.Text("Runs once on mount. Does NOT re-run when signals change.")),
							),
							h.Div(h.Style("display:flex;gap:var(--sp-2)"),
								primitive.Badge(primitive.BadgeProps{Variant: primitive.BadgeSuccess}, g.Text("data-effect")),
								h.Span(g.Text("Re-runs whenever any signal it reads changes. Reactive.")),
							),
							h.Div(h.Style("display:flex;gap:var(--sp-2)"),
								primitive.Badge(primitive.BadgeProps{}, g.Text("data-on-interval")),
								h.Span(g.Text("Runs on a timer regardless of signal changes.")),
							),
						),
					),
				),
			)
		},
	})
}
