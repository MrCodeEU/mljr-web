//go:build showcase

package feedback

import (
	"mljr-web/ui/primitive"
	"mljr-web/ui/registry"
	"mljr-web/ui/token"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "notification-stack", Name: "Notification Stack", Category: "feedback",
		Summary: "Fixed-position persistent notification feed. Add via window._pushNotification(). Auto-dismisses with duration.",
		Code: `// Mount once per page
feedback.NotificationStack(feedback.NotificationStackProps{
    Position: "top-right",
    Max:      5,
})

// Push from anywhere (JS)
window._pushNotification({
    title:    "Deploy succeeded",
    body:     "v2.4.1 is live.",
    variant:  "success",  // "info" | "success" | "warning" | "error"
    duration: 5000,       // ms; 0 = manual dismiss only
})`,
		Render: func(p map[string]string) g.Node {
			return h.Div(
				h.Style("display:flex;flex-direction:column;gap:var(--sp-4);padding:var(--sp-4)"),
				h.P(h.Style("color:var(--muted);font-size:var(--t-sm)"), g.Text("Click buttons to push notifications. They stack top-right and auto-dismiss after 4 s.")),
				h.Div(h.Style("display:flex;flex-wrap:wrap;gap:var(--sp-2)"),
					primitive.Button(primitive.ButtonProps{Variant: token.Primary, Size: token.SizeSM},
						g.Attr("onclick", `window._pushNotification({title:'Deploy succeeded',body:'v2.4.1 is now live on production.',variant:'success',duration:4000})`),
						g.Text("Success"),
					),
					primitive.Button(primitive.ButtonProps{Variant: token.Outline, Size: token.SizeSM},
						g.Attr("onclick", `window._pushNotification({title:'New message',body:'Jordan sent you a file.',variant:'info',duration:4000})`),
						g.Text("Info"),
					),
					primitive.Button(primitive.ButtonProps{Variant: token.Ghost, Size: token.SizeSM},
						g.Attr("onclick", `window._pushNotification({title:'Storage warning',body:'You are at 90% of your quota.',variant:'warning',duration:4000})`),
						g.Text("Warning"),
					),
					primitive.Button(primitive.ButtonProps{Variant: token.Danger, Size: token.SizeSM},
						g.Attr("onclick", `window._pushNotification({title:'Build failed',body:'3 errors in ui/form/input.go.',variant:'error',duration:0})`),
						g.Text("Error (manual dismiss)"),
					),
				),
				NotificationStack(NotificationStackProps{Position: "top-right", Max: 5}),
			)
		},
	})
}
