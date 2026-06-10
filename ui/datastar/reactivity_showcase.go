//go:build showcase

package datastar

import (
	"mljr-web/ui"
	"mljr-web/ui/form"
	"mljr-web/ui/primitive"
	"mljr-web/ui/registry"
	"mljr-web/ui/token"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "ds-reactivity", Name: "Reactivity", Category: "datastar",
		Summary: "data-text, data-show, data-bind, data-attr, data-style, data-class — six ways to reflect signal state in the DOM.",
		Code: `// Reactive text content
data-text="$name || 'anonymous'"

// Conditional display (CSS display:none toggle)
data-show="$open"

// Two-way input binding
data-bind:search   // input value ↔ $search signal

// Reactive HTML attribute
data-attr='{"data-variant": $status, "aria-label": "Status: "+$status}'

// Reactive inline style property
data-style='{"--accent": $dark ? "#000" : "#e23"}'

// Reactive CSS class (adds/removes individual classes)
data-class='{"is-active": $active, "has-error": $error}'`,
		Render: func(p map[string]string) g.Node {
			return h.Div(
				ui.Signals(`{"name":"","open":true,"theme":"primary","active":false,"progress":40}`),
				h.Div(h.Style("display:flex;flex-direction:column;gap:var(--sp-5)"),

					// data-text
					primitive.Card(primitive.CardProps{},
						h.P(h.Style("font-size:var(--t-xs);text-transform:uppercase;font-weight:700;letter-spacing:.06em;opacity:.5;margin-bottom:var(--sp-3)"), g.Text("data-text — reactive content")),
						form.Field(form.FieldProps{Label: "Name"},
							form.Input(form.InputProps{Signal: "name", Placeholder: "Type something…"}),
						),
						h.Div(h.Style("margin-top:var(--sp-3);font-size:var(--t-base)"),
							g.Attr("data-text", "`Hello, ${$name || 'stranger'} 👋`"),
						),
					),

					// data-show
					primitive.Card(primitive.CardProps{},
						h.P(h.Style("font-size:var(--t-xs);text-transform:uppercase;font-weight:700;letter-spacing:.06em;opacity:.5;margin-bottom:var(--sp-3)"), g.Text("data-show — conditional display")),
						h.Div(h.Style("display:flex;gap:var(--sp-2);margin-bottom:var(--sp-3)"),
							primitive.Button(primitive.ButtonProps{Variant: token.Outline, Size: token.SizeSM},
								g.Attr("data-on:click", "$open=!$open"),
								g.Attr("data-text", "$open?'Hide panel':'Show panel'"),
							),
						),
						h.Div(
							g.Attr("data-show", "$open"),
							h.Style("display:none"),
							h.Div(h.Style("padding:var(--sp-4);background:var(--surface-2);border-radius:var(--radius)"),
								g.Text("Panel content — toggled with "),
								h.Code(g.Text("data-show")),
								g.Text(". No JS framework, just CSS "),
								h.Code(g.Text("display:none")),
								g.Text(" flipped by Datastar."),
							),
						),
					),

					// data-attr + data-style
					primitive.Card(primitive.CardProps{},
						h.P(h.Style("font-size:var(--t-xs);text-transform:uppercase;font-weight:700;letter-spacing:.06em;opacity:.5;margin-bottom:var(--sp-3)"), g.Text("data-attr + data-style — reactive attributes")),
						h.Div(h.Style("display:flex;gap:var(--sp-2);flex-wrap:wrap;margin-bottom:var(--sp-3)"),
							primitive.Button(primitive.ButtonProps{Variant: token.Outline, Size: token.SizeSM},
								g.Attr("data-on:click", "$theme='primary'"), g.Text("Primary"),
							),
							primitive.Button(primitive.ButtonProps{Variant: token.Outline, Size: token.SizeSM},
								g.Attr("data-on:click", "$theme='danger'"), g.Text("Danger"),
							),
							primitive.Button(primitive.ButtonProps{Variant: token.Outline, Size: token.SizeSM},
								g.Attr("data-on:click", "$theme='success'"), g.Text("Success"),
							),
						),
						h.Div(
							g.Attr("data-component", "badge"),
							g.Attr("data-attr", `{"data-variant":$theme}`),
							h.Style("font-size:var(--t-base);padding:.5rem 1rem"),
							g.Attr("data-text", "`Variant: ${$theme}`"),
						),
					),

					// data-bind
					primitive.Card(primitive.CardProps{},
						h.P(h.Style("font-size:var(--t-xs);text-transform:uppercase;font-weight:700;letter-spacing:.06em;opacity:.5;margin-bottom:var(--sp-3)"), g.Text("data-bind — two-way binding")),
						h.Div(h.Style("display:flex;flex-direction:column;gap:var(--sp-3)"),
							h.Div(h.Style("display:flex;align-items:center;gap:var(--sp-3)"),
								h.Span(h.Style("min-width:5rem;font-size:var(--t-sm);opacity:.6"), g.Text("Range:")),
								h.Input(
									h.Type("range"), h.Min("0"), h.Max("100"),
									g.Attr("data-component", "slider"),
									g.Attr("data-bind:progress"),
									g.Attr("data-on:input", "$progress=Number(evt.target.value)"),
									h.Style("flex:1"),
								),
								h.Span(
									h.Style("min-width:3rem;text-align:right;font-weight:700;font-family:var(--font-display)"),
									g.Attr("data-text", "$progress+'%'"),
								),
							),
							h.Div(h.Style("height:8px;background:var(--surface-2);border-radius:4px;overflow:hidden"),
								h.Div(
									g.Attr("data-style", `{"width":$progress+"%","background":"var(--accent)"}`),
									h.Style("height:100%;transition:width .1s"),
								),
							),
						),
					),
				),
			)
		},
	})
}
