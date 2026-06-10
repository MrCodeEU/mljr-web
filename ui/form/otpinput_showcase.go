//go:build showcase

package form

import (
	"mljr-web/ui/primitive"
	"mljr-web/ui/registry"
	"mljr-web/ui/token"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "otp-input", Name: "OTP Input", Category: "form",
		Summary: "One-time-password entry — single-digit boxes with auto-advance, backspace navigation, and paste support.",
		Code: `form.OTPInput(form.OTPInputProps{
    Name:   "code",
    Length: 6,
})`,
		Render: func(p map[string]string) g.Node {
			return h.Div(
				h.Style("display:flex;flex-direction:column;gap:var(--sp-6);align-items:flex-start"),
				h.Div(
					h.Style("display:flex;flex-direction:column;gap:var(--sp-3)"),
					h.P(h.Style("font-size:var(--t-xs);text-transform:uppercase;font-weight:700;letter-spacing:.06em;opacity:.5;margin:0"), g.Text("6-digit code")),
					OTPInput(OTPInputProps{Name: "code6", Length: 6}),
				),
				h.Div(
					h.Style("display:flex;flex-direction:column;gap:var(--sp-3)"),
					h.P(h.Style("font-size:var(--t-xs);text-transform:uppercase;font-weight:700;letter-spacing:.06em;opacity:.5;margin:0"), g.Text("4-digit PIN")),
					OTPInput(OTPInputProps{Name: "pin", Length: 4, Label: "PIN"}),
				),
				primitive.Button(primitive.ButtonProps{Variant: token.Primary},
					g.Text("Verify code"),
				),
			)
		},
	})
}
