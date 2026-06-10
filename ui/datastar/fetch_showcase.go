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
		Slug: "ds-fetch", Name: "Fetch (SSE)", Category: "datastar",
		Summary: "@get and @post send signals to the server via SSE. data-indicator shows loading state. The connection is request-response (not persistent).",
		Code: `// GET — sends all signals as query params, streams response events
data-on:click="@get('/api/items')"

// POST — sends signals as JSON body
data-on:click="@post('/api/save')"

// Filter signals sent (default: exclude _ prefixed)
data-on:click="@get('/api/search', {filterSignals: {include: /^q/}})"

// Loading indicator — sets aria-busy + data-loading attribute
data-indicator:fetching
// use in CSS: [data-loading] { opacity:.5 }
// or show a spinner:
data-show="$fetching"   // $fetching is true while SSE open

// Server side (Go):
sse := datastar.NewSSE(c.Response().Writer, c.Request())
sse.PatchElements(web.RenderToString(myFragment))
datastar.MarshalAndPatchSignals(sse, map[string]any{"loading": false})`,
		Render: func(p map[string]string) g.Node {
			return h.Div(
				ui.Signals(`{"_q":"","echoResult":"","_fetching":false}`),
				h.Div(h.Style("display:flex;flex-direction:column;gap:var(--sp-5)"),

					// Echo demo
					primitive.Card(primitive.CardProps{},
						h.P(h.Style("font-size:var(--t-xs);text-transform:uppercase;font-weight:700;letter-spacing:.06em;opacity:.5;margin-bottom:var(--sp-3)"), g.Text("@post — echo server round-trip")),
						h.Div(h.Style("display:flex;gap:var(--sp-2)"),
							form.Input(form.InputProps{Signal: "_q", Placeholder: "Type a message…"}),
							primitive.Button(primitive.ButtonProps{Variant: token.Primary},
								g.Attr("data-indicator:_fetching"),
								g.Attr("data-on:click", "@post('/demo/echo')"),
								g.Attr("data-attr", `{"aria-busy":$_fetching}`),
								g.Attr("data-text", "$_fetching?'Sending…':'Send to server'"),
							),
						),
						h.Div(
							h.ID("echo-result"),
							h.Style("margin-top:var(--sp-3);padding:var(--sp-3);background:var(--surface-2);border-radius:var(--radius);font-size:var(--t-sm);min-height:2.5rem"),
							g.Attr("data-show", `$echoResult!==''`),
							h.Style("display:none"),
							g.Attr("data-text", "$echoResult"),
						),
						h.P(h.Style("font-size:var(--t-xs);opacity:.5;margin-top:var(--sp-2)"),
							g.Text("Server reads $q signal, patches #echo-result fragment + updates $echoResult signal."),
						),
					),

					// Search demo
					primitive.Card(primitive.CardProps{},
						h.P(h.Style("font-size:var(--t-xs);text-transform:uppercase;font-weight:700;letter-spacing:.06em;opacity:.5;margin-bottom:var(--sp-3)"), g.Text("@get + debounce — search-as-you-type")),
						form.Input(form.InputProps{
							Signal:      "_q",
							Placeholder: "Search fruits…",
						},
							g.Attr("data-on:input__debounce.300ms", "@get('/demo/search')"),
						),
						h.Div(
							h.ID("search-results"),
							h.Style("margin-top:var(--sp-3)"),
							h.P(h.Style("opacity:.4;font-size:var(--t-sm)"), g.Text("Results appear here…")),
						),
					),

					// Indicator explanation
					primitive.Card(primitive.CardProps{},
						h.P(h.Style("font-size:var(--t-xs);text-transform:uppercase;font-weight:700;letter-spacing:.06em;opacity:.5;margin-bottom:var(--sp-3)"), g.Text("data-indicator — loading state")),
						h.Div(h.Style("display:flex;flex-direction:column;gap:var(--sp-2);font-size:var(--t-sm)"),
							h.Code(g.Text(`data-indicator:myLoading`)),
							h.P(g.Text("Adds this attribute to an element that triggers fetch. Sets signal "), h.Code(g.Text("$myLoading=true")), g.Text(" while the SSE stream is open, "), h.Code(g.Text("false")), g.Text(" after.")),
							h.P(g.Text("Also sets "), h.Code(g.Text("aria-busy")), g.Text(" and "), h.Code(g.Text("data-loading")), g.Text(" on the element. Use in CSS or with "), h.Code(g.Text("data-show")), g.Text(" to render a spinner.")),
						),
					),
				),
			)
		},
	})
}
