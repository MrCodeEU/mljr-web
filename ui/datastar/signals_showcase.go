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
		Slug: "ds-signals", Name: "Signals", Category: "datastar",
		Summary: "data-signals declares reactive state. data-computed derives new signals from expressions. Signals scope to their subtree.",
		Code: `// Declare
data-signals='{"count":0,"user":{"name":"Alice","score":0}}'

// Read in expression
data-text="$count"
data-text="'Hi '+$user.name+' — '+$user.score+' pts'"

// Write from event
data-on:click="$count++"
data-on:click="$user.score+= 10"

// Derived signal (data-computed)
data-computed:doubled="$count * 2"
data-text="$doubled"   // reactive, updates automatically

// Signals prefixed _ are excluded from SSE submissions
data-signals='{"_localOnly": true}'`,
		Render: func(p map[string]string) g.Node {
			return h.Div(
				ui.Signals(`{"count":0,"user":{"name":"Alice","score":0}}`),
				h.Div(h.Style("display:flex;flex-direction:column;gap:var(--sp-5)"),

					// Counter + computed
					primitive.Card(primitive.CardProps{},
						h.Div(
							g.Attr("data-computed:doubled", "$count * 2"),
							h.P(h.Style("font-size:var(--t-xs);text-transform:uppercase;font-weight:700;letter-spacing:.06em;opacity:.5;margin-bottom:var(--sp-3)"), g.Text("Simple signal + computed")),
							h.Div(h.Style("display:flex;align-items:center;gap:var(--sp-3)"),
								primitive.Button(primitive.ButtonProps{Variant: token.Outline, Size: token.SizeSM},
									g.Attr("data-on:click", "$count=Math.max(0,$count-1)"),
									g.Text("−"),
								),
								h.Span(h.Style("font-size:var(--t-2xl);font-weight:900;font-family:var(--font-display);min-width:3rem;text-align:center"),
									g.Attr("data-text", "$count"),
								),
								primitive.Button(primitive.ButtonProps{Variant: token.Primary, Size: token.SizeSM},
									g.Attr("data-on:click", "$count++"),
									g.Text("+"),
								),
								primitive.Button(primitive.ButtonProps{Variant: token.Ghost, Size: token.SizeSM},
									g.Attr("data-on:click", "$count=0"),
									g.Text("Reset"),
								),
							),
							h.Div(h.Style("display:flex;gap:var(--sp-4);margin-top:var(--sp-3);font-size:var(--t-sm)"),
								h.Span(g.Text("doubled: "), h.Strong(g.Attr("data-text", "$doubled"))),
								h.Span(g.Text("parity: "), h.Strong(g.Attr("data-text", "$count%2===0?'even':'odd'"))),
								h.Span(g.Text("level: "), h.Strong(g.Attr("data-text", "$count<5?'newbie':$count<15?'pro':'legend'"))),
							),
						),
					),

					// Nested object signal
					primitive.Card(primitive.CardProps{},
						h.P(h.Style("font-size:var(--t-xs);text-transform:uppercase;font-weight:700;letter-spacing:.06em;opacity:.5;margin-bottom:var(--sp-3)"), g.Text("Nested signal object")),
						h.Div(h.Style("font-size:var(--t-base);margin-bottom:var(--sp-3)"),
							g.Attr("data-text", "`${$user.name} has ${$user.score} pts`"),
						),
						h.Div(h.Style("display:flex;gap:var(--sp-2);flex-wrap:wrap"),
							primitive.Button(primitive.ButtonProps{Variant: token.Outline, Size: token.SizeSM},
								g.Attr("data-on:click", "$user.score+=10"),
								g.Text("+10 points"),
							),
							primitive.Button(primitive.ButtonProps{Variant: token.Outline, Size: token.SizeSM},
								g.Attr("data-on:click", "$user.name=$user.name==='Alice'?'Bob':'Alice'"),
								g.Text("Swap name"),
							),
							primitive.Button(primitive.ButtonProps{Variant: token.Ghost, Size: token.SizeSM},
								g.Attr("data-on:click", "$user={name:'Alice',score:0}"),
								g.Text("Reset"),
							),
						),
					),

					// Signal scoping note
					primitive.Card(primitive.CardProps{},
						h.P(h.Style("font-size:var(--t-xs);text-transform:uppercase;font-weight:700;letter-spacing:.06em;opacity:.5;margin-bottom:var(--sp-3)"), g.Text("Signal scope")),
						h.Div(h.Style("display:flex;gap:var(--sp-4)"),
							h.Div(
								ui.Signals(`{"local":0}`),
								h.Style("flex:1;padding:var(--sp-3);border:var(--bw-1) solid var(--line);border-radius:var(--radius)"),
								h.P(h.Style("font-size:var(--t-xs);opacity:.6;margin-bottom:var(--sp-2)"), g.Text("Subtree A")),
								h.Div(h.Style("display:flex;align-items:center;gap:var(--sp-2)"),
									primitive.Button(primitive.ButtonProps{Variant: token.Outline, Size: token.SizeSM},
										g.Attr("data-on:click", "$local++"), g.Text("+"),
									),
									h.Span(g.Attr("data-text", "$local")),
								),
							),
							h.Div(
								ui.Signals(`{"local":0}`),
								h.Style("flex:1;padding:var(--sp-3);border:var(--bw-1) solid var(--line);border-radius:var(--radius)"),
								h.P(h.Style("font-size:var(--t-xs);opacity:.6;margin-bottom:var(--sp-2)"), g.Text("Subtree B (own $local)")),
								h.Div(h.Style("display:flex;align-items:center;gap:var(--sp-2)"),
									primitive.Button(primitive.ButtonProps{Variant: token.Outline, Size: token.SizeSM},
										g.Attr("data-on:click", "$local++"), g.Text("+"),
									),
									h.Span(g.Attr("data-text", "$local")),
								),
							),
						),
					),
				),
			)
		},
	})
}
