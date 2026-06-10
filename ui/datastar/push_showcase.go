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
		Slug: "ds-push", Name: "Server Push (SSE)", Category: "datastar",
		Summary: "PatchElements morphs DOM fragments. PatchSignals updates client state. data-on-signal-patch reacts to incoming patches.",
		Code: `// Server — Go SSE handler
sse := datastar.NewSSE(c.Response().Writer, c.Request())

// Patch a DOM fragment (morphs by id)
sse.PatchElements(web.RenderToString(
    h.Div(h.ID("counter"), g.Text(fmt.Sprintf("%d", count))),
))

// Patch signals (updates client state)
datastar.MarshalAndPatchSignals(sse, map[string]any{
    "serverTime": time.Now().Format("15:04:05"),
    "online":     true,
})

// For persistent streams: loop + flush, check ctx
for {
    select {
    case <-c.Request().Context().Done():
        return nil
    case <-ticker.C:
        datastar.MarshalAndPatchSignals(sse, ...)
    }
}

// Client — react to signal patches
data-on-signal-patch="if(patch.serverTime) console.log('new time:', $serverTime)"`,
		Render: func(p map[string]string) g.Node {
			return h.Div(
				ui.Signals(`{"serverTime":"—","online":false,"_streamActive":false}`),
				h.Div(h.Style("display:flex;flex-direction:column;gap:var(--sp-5)"),

					// Live clock via persistent SSE
					primitive.Card(primitive.CardProps{},
						h.P(h.Style("font-size:var(--t-xs);text-transform:uppercase;font-weight:700;letter-spacing:.06em;opacity:.5;margin-bottom:var(--sp-3)"), g.Text("Persistent SSE — live server clock")),
						h.Div(h.Style("font-size:var(--t-2xl);font-weight:900;font-family:var(--font-display);letter-spacing:-.02em"),
							g.Attr("data-text", "$serverTime"),
						),
						h.Div(h.Style("display:flex;gap:var(--sp-2);margin-top:var(--sp-3)"),
							primitive.Button(primitive.ButtonProps{Variant: token.Primary, Size: token.SizeSM},
								g.Attr("data-show", "!$_streamActive"),
								g.Attr("data-on:click", "$_streamActive=true;@get('/demo/time')"),
								g.Text("Start stream"),
							),
							h.Div(
								g.Attr("data-show", "$_streamActive"),
								h.Style("display:none;display:flex;align-items:center;gap:var(--sp-2)"),
								h.Span(h.Style("width:8px;height:8px;border-radius:50%;background:var(--success)")),
								h.Span(h.Style("font-size:var(--t-sm)"), g.Text("Streaming — auto-closes after 15 s")),
							),
						),
					),

					// PatchElements vs PatchSignals
					primitive.Card(primitive.CardProps{},
						h.P(h.Style("font-size:var(--t-xs);text-transform:uppercase;font-weight:700;letter-spacing:.06em;opacity:.5;margin-bottom:var(--sp-3)"), g.Text("PatchElements vs PatchSignals")),
						h.Div(h.Style("display:flex;flex-direction:column;gap:var(--sp-3);font-size:var(--t-sm)"),
							h.Div(h.Style("display:flex;gap:var(--sp-2)"),
								primitive.Badge(primitive.BadgeProps{Variant: primitive.BadgeInfo}, g.Text("PatchElements")),
								h.Span(g.Text("Morphs DOM fragment by id. Best for server-rendered HTML (lists, tables, paginated content).")),
							),
							h.Div(h.Style("display:flex;gap:var(--sp-2)"),
								primitive.Badge(primitive.BadgeProps{Variant: primitive.BadgeSuccess}, g.Text("PatchSignals")),
								h.Span(g.Text("Updates client signal values. Best for status, counters, flags. Client renders from signals.")),
							),
							h.Div(h.Style("display:flex;gap:var(--sp-2)"),
								primitive.Badge(primitive.BadgeProps{}, g.Text("Both")),
								h.Span(g.Text("Can be mixed in one SSE response. PatchSignals first, then PatchElements — signals are available when fragment renders.")),
							),
						),
					),

					// on-signal-patch
					primitive.Card(primitive.CardProps{},
						h.P(h.Style("font-size:var(--t-xs);text-transform:uppercase;font-weight:700;letter-spacing:.06em;opacity:.5;margin-bottom:var(--sp-3)"), g.Text("data-on-signal-patch — react to incoming patches")),
						h.Div(
							g.Attr("data-on-signal-patch", "patch.serverTime&&console.log('Server time updated:',patch.serverTime)"),
						),
						h.P(h.Style("font-size:var(--t-sm)"),
							g.Text("Fires whenever "),
							h.Code(g.Text("PatchSignals")),
							g.Text(" arrives from any SSE. Receives the "),
							h.Code(g.Text("patch")),
							g.Text(" object. Useful for analytics, notifications, or triggering client-side animations when new data arrives."),
						),
					),
				),
			)
		},
	})
}
