//go:build showcase

package special

import (
	"mljr-web/ui/layout"
	"mljr-web/ui/primitive"
	"mljr-web/ui/registry"
	"mljr-web/ui/token"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "theme-toggle", Name: "Theme / Mode Toggle", Category: "special",
		Summary: "Cycles $theme (swissbrut↔ink) and $mode (light↔dark). ThemeToggleRoot must appear once per page.",
		Code: `// in PageShell body — one per page
special.ThemeToggleRoot(token.ThemeSwissBrut, token.ModeLight)

// trigger buttons (place in navbar)
special.ThemeToggle() // cycles swissbrut ↔ ink
special.ModeToggle()  // flips light ↔ dark`,
		Render: func(p map[string]string) g.Node {
			return layout.Stack(layout.StackProps{Axis: "h", Gap: "sm"},
				ThemeToggle(),
				ModeToggle(),
				h.Span(h.Style("font-size:var(--t-sm);opacity:.6"), g.Text("Toggle theme or mode above")),
			)
		},
	})

	registry.Register(&registry.Component{
		Slug: "captcha", Name: "Captcha (Altcha)", Category: "special",
		Summary: "Self-hosted proof-of-work captcha. No CDN, no telemetry. challenge= sets the JSON challenge endpoint.",
		Code: `// import "mljr-web/ui/special"
special.Captcha(special.CaptchaProps{ChallengeURL: "/api/altcha"})`,
		Render: func(p map[string]string) g.Node {
			return h.Div(h.Style("width:100%;max-width:400px"),
				primitive.Card(primitive.CardProps{Tone: token.ToneNone},
					h.P(h.Style("font-size:var(--t-sm);opacity:.7;margin-bottom:var(--sp-3)"),
						g.Text("Uses the showcase /api/altcha endpoint and writes the verified payload into the widget's hidden input."),
					),
					Captcha(CaptchaProps{ChallengeURL: "/api/altcha"}),
				),
			)
		},
	})
}
