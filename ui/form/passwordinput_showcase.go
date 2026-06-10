//go:build showcase

package form

import (
	"mljr-web/ui/primitive"
	"mljr-web/ui/registry"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "password-input", Name: "Password Input", Category: "form",
		Summary: "Password field with show/hide toggle. Datastar signal tracks visibility; data-attr swaps input type.",
		Code: `form.PasswordInput(form.PasswordInputProps{
    Name:         "password",
    Placeholder:  "Enter password",
    Autocomplete: "current-password",
})`,
		Render: func(p map[string]string) g.Node {
			return h.Div(
				h.Style("display:flex;flex-direction:column;gap:var(--sp-5);max-width:360px"),
				Field(FieldProps{Label: "Current password"},
					PasswordInput(PasswordInputProps{
						Name:         "current",
						Placeholder:  "Enter current password",
						Autocomplete: "current-password",
					}),
				),
				Field(FieldProps{Label: "New password", Hint: "Min 8 characters"},
					PasswordInput(PasswordInputProps{
						Name:         "newpw",
						Placeholder:  "Choose a strong password",
						Autocomplete: "new-password",
					}),
				),
				primitive.Button(primitive.ButtonProps{},
					g.Attr("data-component", "button"),
					g.Text("Update password"),
				),
			)
		},
	})
}
