//go:build showcase

package datastar

import (
	"mljr-web/ui"
	"mljr-web/ui/primitive"
	"mljr-web/ui/registry"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "ds-modifiers", Name: "Modifiers", Category: "datastar",
		Summary: "Event modifiers after __ control how handlers fire: debounce, throttle, window, document, stop, prevent, once, passive.",
		Code: `// Format: data-on:EVENT__mod1.val__mod2="expression"

data-on:input__debounce.300ms="@get('/search')"  // wait 300ms after last input
data-on:scroll__throttle.100ms="$y=window.scrollY" // max once per 100ms
data-on:click__window="$open=false"               // listen on window, not element
data-on:keydown__document="..."                   // listen on document
data-on:click__stop="$open=!$open"               // stopPropagation()
data-on:submit__prevent="@post('/api')"           // preventDefault()
data-on:click__once="$seen=true"                  // fire exactly once
data-on:scroll__passive="$y=scrollY"             // passive: true (perf)
data-on:keydown__window="if(evt.key==='/'){ evt.preventDefault(); $searchOpen=true }"`,
		Render: func(p map[string]string) g.Node {
			type mod struct {
				mod, syntax, when, demo string
			}
			mods := []struct {
				name, syntax, desc string
			}{
				{"debounce", `data-on:input__debounce.300ms`, "Delays firing until N ms after the last event. Essential for search-as-you-type."},
				{"throttle", `data-on:scroll__throttle.100ms`, "Fires at most once per N ms. Good for scroll/resize handlers."},
				{"window", `data-on:click__window`, "Listens on window instead of the element. Use for outside-click close."},
				{"document", `data-on:keydown__document`, "Listens on document. Slightly narrower than window."},
				{"stop", `data-on:click__stop`, "Calls event.stopPropagation(). Prevents event from bubbling."},
				{"prevent", `data-on:submit__prevent`, "Calls event.preventDefault(). Prevents form submit / link follow."},
				{"once", `data-on:click__once`, "Handler fires exactly one time, then unregisters itself."},
				{"passive", `data-on:scroll__passive`, "Adds {passive:true} to addEventListener. Improves scroll perf."},
			}

			rows := make([]g.Node, len(mods))
			for i, m := range mods {
				rows[i] = h.Tr(
					h.Td(h.Style("padding:var(--sp-2) var(--sp-3);border-bottom:var(--bw-1) solid color-mix(in srgb,var(--line) 20%,transparent)"),
						primitive.Badge(primitive.BadgeProps{Variant: primitive.BadgeInfo}, g.Text(m.name)),
					),
					h.Td(h.Style("padding:var(--sp-2) var(--sp-3);border-bottom:var(--bw-1) solid color-mix(in srgb,var(--line) 20%,transparent)"),
						h.Code(h.Style("font-size:var(--t-xs)"), g.Text(m.syntax)),
					),
					h.Td(h.Style("padding:var(--sp-2) var(--sp-3);border-bottom:var(--bw-1) solid color-mix(in srgb,var(--line) 20%,transparent);font-size:var(--t-xs);opacity:.7"),
						g.Text(m.desc),
					),
				)
			}

			return h.Div(
				ui.Signals(`{"_debLog":"","_throttLog":0,"_scrollY":0,"_stopOuter":0,"_stopInner":0,"_onceCount":0}`),

				// Reference table
				primitive.Card(primitive.CardProps{},
					h.P(h.Style("font-size:var(--t-xs);text-transform:uppercase;font-weight:700;letter-spacing:.06em;opacity:.5;margin-bottom:var(--sp-3)"), g.Text("All modifiers")),
					h.Div(h.Style("overflow-x:auto"),
						h.Table(
							g.Attr("data-component", "table"),
							h.THead(h.Tr(
								h.Th(g.Text("Modifier")),
								h.Th(g.Text("Syntax")),
								h.Th(g.Text("Effect")),
							)),
							h.TBody(g.Group(rows)),
						),
					),
				),

				h.Div(h.Style("display:flex;flex-direction:column;gap:var(--sp-4);margin-top:var(--sp-5)"),

					// Debounce live demo
					primitive.Card(primitive.CardProps{},
						h.P(h.Style("font-size:var(--t-xs);text-transform:uppercase;font-weight:700;letter-spacing:.06em;opacity:.5;margin-bottom:var(--sp-3)"), g.Text("debounce live demo")),
						h.Input(
							g.Attr("data-component", "input"),
							h.Type("text"),
							h.Placeholder("Type fast…"),
							g.Attr("data-on:input", "$_debLog=evt.target.value"),
							g.Attr("data-on:input__debounce.500ms", `$_debLog='[debounced 500ms] '+evt.target.value`),
						),
						h.Div(h.Style("margin-top:var(--sp-2);font-size:var(--t-sm)"),
							h.Span(h.Style("opacity:.6"), g.Text("Debounced: ")),
							h.Strong(g.Attr("data-text", "$_debLog||'(type to see)'")),
						),
					),

					// stop propagation demo
					primitive.Card(primitive.CardProps{},
						h.P(h.Style("font-size:var(--t-xs);text-transform:uppercase;font-weight:700;letter-spacing:.06em;opacity:.5;margin-bottom:var(--sp-3)"), g.Text("__stop vs no stop")),
						h.Div(h.Style("display:flex;gap:var(--sp-4)"),
							h.Div(
								h.Style("padding:var(--sp-4);border:var(--bw-1) solid var(--line);border-radius:var(--radius);cursor:pointer"),
								g.Attr("data-on:click", "$_stopOuter++"),
								h.P(h.Style("font-size:var(--t-xs);opacity:.5;margin-bottom:var(--sp-2)"), g.Text("Outer (clicks: "), h.Span(g.Attr("data-text", "$_stopOuter")), g.Text(")")),
								h.Button(
									g.Attr("data-component", "button"),
									g.Attr("data-on:click__stop", "$_stopInner++"),
									g.Attr("data-text", "`Inner __stop (${$_stopInner})`"),
								),
							),
						),
						h.P(h.Style("font-size:var(--t-xs);opacity:.5;margin-top:var(--sp-2)"),
							g.Text("Inner button has __stop — clicks don't propagate to outer div."),
						),
					),

					// once demo
					primitive.Card(primitive.CardProps{},
						h.P(h.Style("font-size:var(--t-xs);text-transform:uppercase;font-weight:700;letter-spacing:.06em;opacity:.5;margin-bottom:var(--sp-3)"), g.Text("__once live demo")),
						h.Button(
							g.Attr("data-component", "button"),
							g.Attr("data-on:click__once", "$_onceCount++"),
							g.Attr("data-attr", `{"data-variant":$_onceCount>0?"success":"outline"}`),
							g.Attr("data-text", `$_onceCount>0?'Fired ('+$_onceCount+'x) — no more':'Click me (fires once)'`),
						),
					),
				),
			)
		},
	})
}
