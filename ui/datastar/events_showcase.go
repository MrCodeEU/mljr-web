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
		Slug: "ds-events", Name: "Events", Category: "datastar",
		Summary: "data-on:EVENT handles DOM events. data-on-interval runs on a timer. data-on-intersect fires on viewport entry.",
		Code: `// Basic click
data-on:click="$open=!$open"

// Input with 300ms debounce (fires 300ms after typing stops)
data-on:input__debounce.300ms="@get('/api/search')"

// Keyboard shortcut on the window (not just the element)
data-on:keydown__window="if(evt.key==='Escape')$open=false"

// Outside-click to close a dropdown
data-on:click__window="$ddOpen=false"
data-on:click__stop="$ddOpen=!$ddOpen"  // stop propagation on trigger

// Interval timer — runs expression every N ms
data-on-interval__duration.1s="$ticks++"

// Intersection observer — fires when element enters viewport
data-on-intersect="$visible=true"

// Once modifier — handler fires exactly once
data-on:click__once="$welcomed=true"`,
		Render: func(p map[string]string) g.Node {
			return h.Div(
				ui.Signals(`{"ticks":0,"keyLog":"","inputVal":"","debounced":"","welcomed":false,"ddOpen":false}`),
				h.Div(h.Style("display:flex;flex-direction:column;gap:var(--sp-5)"),

					// Interval timer
					primitive.Card(primitive.CardProps{},
						h.P(h.Style("font-size:var(--t-xs);text-transform:uppercase;font-weight:700;letter-spacing:.06em;opacity:.5;margin-bottom:var(--sp-3)"), g.Text("data-on-interval — timer")),
						h.Div(
							g.Attr("data-on-interval__duration.1s", "$ticks++"),
						),
						h.Div(h.Style("display:flex;align-items:center;gap:var(--sp-4)"),
							h.Div(
								h.Style("font-size:var(--t-2xl);font-weight:900;font-family:var(--font-display)"),
								g.Attr("data-text", "$ticks"),
							),
							h.Span(h.Style("opacity:.5;font-size:var(--t-sm)"), g.Text("ticks (1 s interval)")),
							primitive.Button(primitive.ButtonProps{Variant: token.Ghost, Size: token.SizeSM},
								g.Attr("data-on:click", "$ticks=0"), g.Text("Reset"),
							),
						),
					),

					// Keyboard + window events
					primitive.Card(primitive.CardProps{},
						h.P(h.Style("font-size:var(--t-xs);text-transform:uppercase;font-weight:700;letter-spacing:.06em;opacity:.5;margin-bottom:var(--sp-3)"), g.Text("data-on:keydown__window — global shortcut")),
						h.Div(
							g.Attr("data-on:keydown__window", "$keyLog=evt.key"),
						),
						h.Div(h.Style("display:flex;align-items:center;gap:var(--sp-3)"),
							h.Span(h.Style("opacity:.6;font-size:var(--t-sm)"), g.Text("Last key pressed:")),
							h.Code(
								h.Style("font-size:var(--t-base);background:var(--surface-2);padding:.15rem .4rem;border-radius:var(--radius)"),
								g.Attr("data-text", "$keyLog||'(press any key)'"),
							),
						),
					),

					// Debounce
					primitive.Card(primitive.CardProps{},
						h.P(h.Style("font-size:var(--t-xs);text-transform:uppercase;font-weight:700;letter-spacing:.06em;opacity:.5;margin-bottom:var(--sp-3)"), g.Text("data-on:input__debounce — fire after typing stops")),
						h.Input(
							g.Attr("data-component", "input"),
							h.Type("text"),
							h.Placeholder("Type fast…"),
							g.Attr("data-bind:inputVal"),
							g.Attr("data-on:input", "$inputVal=evt.target.value"),
							g.Attr("data-on:input__debounce.500ms", "$debounced=evt.target.value"),
						),
						h.Div(h.Style("display:flex;gap:var(--sp-4);margin-top:var(--sp-3);font-size:var(--t-sm)"),
							h.Div(g.Text("Live: "), h.Strong(g.Attr("data-text", "$inputVal||'—'"))),
							h.Div(g.Text("Debounced (500ms): "), h.Strong(g.Attr("data-text", "$debounced||'—'"))),
						),
					),

					// Outside click + stop propagation
					primitive.Card(primitive.CardProps{},
						h.P(h.Style("font-size:var(--t-xs);text-transform:uppercase;font-weight:700;letter-spacing:.06em;opacity:.5;margin-bottom:var(--sp-3)"), g.Text("data-on:click__window + __stop — outside-click close")),
						h.Div(
							h.Style("position:relative;display:inline-block"),
							g.Attr("data-on:click__window", "$ddOpen=false"),
							h.Div(
								g.Attr("data-on:click__stop", "$ddOpen=!$ddOpen"),
								primitive.Button(primitive.ButtonProps{Variant: token.Outline},
									g.Attr("data-text", "$ddOpen?'▲ Close':'▼ Open menu'"),
								),
							),
							h.Div(
								g.Attr("data-show", "$ddOpen"),
								h.Style("display:none;position:absolute;top:calc(100% + var(--sp-1));left:0;z-index:50;background:var(--surface);border:var(--border-w) solid var(--line);border-radius:var(--radius);box-shadow:var(--shadow);min-width:160px"),
								h.Div(h.Style("padding:var(--sp-2) var(--sp-4);font-size:var(--t-sm)"), g.Text("Option A")),
								h.Div(h.Style("padding:var(--sp-2) var(--sp-4);font-size:var(--t-sm)"), g.Text("Option B")),
								h.Div(h.Style("padding:var(--sp-2) var(--sp-4);font-size:var(--t-sm)"), g.Text("Option C")),
							),
						),
					),

					// Once modifier
					primitive.Card(primitive.CardProps{},
						h.P(h.Style("font-size:var(--t-xs);text-transform:uppercase;font-weight:700;letter-spacing:.06em;opacity:.5;margin-bottom:var(--sp-3)"), g.Text("__once — fires exactly one time")),
						h.Div(h.Style("display:flex;align-items:center;gap:var(--sp-3)"),
							primitive.Button(primitive.ButtonProps{Variant: token.Outline},
								g.Attr("data-on:click__once", "$welcomed=true"),
								g.Attr("data-attr", `{"data-variant":$welcomed?"success":""}`),
								g.Attr("data-text", "$welcomed?'✓ Already clicked':'Click me (once only)'"),
							),
						),
					),

					// Intersect
					primitive.Card(primitive.CardProps{},
						h.P(h.Style("font-size:var(--t-xs);text-transform:uppercase;font-weight:700;letter-spacing:.06em;opacity:.5;margin-bottom:var(--sp-3)"), g.Text("data-on-intersect — fires when element enters viewport")),
						h.Div(h.Style("max-height:80px;overflow-y:auto;border:var(--bw-1) solid var(--line);border-radius:var(--radius)"),
							h.Div(h.Style("height:120px;display:flex;align-items:center;justify-content:center;opacity:.4;font-size:var(--t-sm)"), g.Text("↓ scroll down ↓")),
							h.Div(
								g.Attr("data-on-intersect", "$welcomed=true"),
								h.Style("height:60px;display:flex;align-items:center;justify-content:center"),
								g.Attr("data-text", "$welcomed?'👀 I was seen!':'(not visible yet)'"),
							),
						),
					),
				),
			)
		},
	})
}
