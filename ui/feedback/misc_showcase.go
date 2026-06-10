//go:build showcase

package feedback

import (
	"mljr-web/ui/icon"
	"mljr-web/ui/primitive"
	"mljr-web/ui/registry"
	"mljr-web/ui/token"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "notification-badge", Name: "Notification Badge", Category: "feedback",
		Summary: "Count badge overlaid top-right on any content — icons, buttons, avatars.",
		Code: `feedback.NotificationBadge(feedback.NotificationBadgeProps{Count: 5},
    icon.Icon("lucide:bell"),
)
// Dot mode
feedback.NotificationBadge(feedback.NotificationBadgeProps{Dot: true},
    primitive.Button(...),
)`,
		Render: func(p map[string]string) g.Node {
			return h.Div(
				h.Style("display:flex;gap:var(--sp-6);align-items:center;flex-wrap:wrap"),
				NotificationBadge(NotificationBadgeProps{Count: 3},
					icon.Icon("lucide:bell", icon.Props{Size: "1.5rem"}),
				),
				NotificationBadge(NotificationBadgeProps{Count: 99},
					icon.Icon("lucide:mail", icon.Props{Size: "1.5rem"}),
				),
				NotificationBadge(NotificationBadgeProps{Count: 150, Max: 99},
					icon.Icon("lucide:inbox", icon.Props{Size: "1.5rem"}),
				),
				NotificationBadge(NotificationBadgeProps{Dot: true},
					primitive.Button(primitive.ButtonProps{Variant: token.Outline}, g.Text("Updates")),
				),
				NotificationBadge(NotificationBadgeProps{Count: 7, Color: "var(--success)"},
					primitive.Button(primitive.ButtonProps{Variant: token.Primary}, g.Text("Messages")),
				),
			)
		},
	})

	registry.Register(&registry.Component{
		Slug: "loading-overlay", Name: "Loading Overlay", Category: "feedback",
		Summary: "Full-surface spinner overlay driven by a Datastar signal. Covers content while async operations run.",
		Code: `// Place inside a position:relative container
feedback.LoadingOverlay(feedback.LoadingOverlayProps{
    SignalName: "_loading",
    Text:       "Saving changes…",
})
// Trigger with:
// data-on:click="$_loading=true; @post('/save')"`,
		Render: func(p map[string]string) g.Node {
			return h.Div(
				h.Style("display:flex;flex-direction:column;gap:var(--sp-4)"),
				h.Div(
					h.Style("position:relative;padding:var(--sp-6);border:var(--bw-2) solid var(--line);border-radius:var(--radius);background:var(--surface);min-height:180px"),
					h.P(h.Style("font-weight:700;margin:0 0 var(--sp-2)"), g.Text("Dashboard content")),
					h.P(h.Style("color:var(--muted);font-size:var(--t-sm)"), g.Text("Click the button to simulate a loading state. The overlay covers this card.")),
					LoadingOverlay(LoadingOverlayProps{SignalName: "_lo1", Text: "Saving changes…"}),
				),
				primitive.Button(primitive.ButtonProps{Variant: token.Primary},
					g.Attr("data-on:click", "$_lo1=true;setTimeout(()=>$_lo1=false,2500)"),
					g.Text("Trigger 2.5s load"),
				),
			)
		},
	})
}
