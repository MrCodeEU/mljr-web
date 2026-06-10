//go:build showcase

package form

import (
	"mljr-web/ui"
	"mljr-web/ui/layout"
	"mljr-web/ui/primitive"
	"mljr-web/ui/registry"
	"mljr-web/ui/token"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "input", Name: "Input", Category: "form",
		Summary: "Text input bound to a Datastar signal. Supports type, placeholder, required.",
		Code: `form.Field(form.FieldProps{Label: "Email", Hint: "We won't spam you", ErrorFor: "email"},
    form.Input(form.InputProps{
        Type:        "email",
        Placeholder: "you@example.com",
        Signal:      "email",
        Required:    true,
    }),
)`,
		Render: func(p map[string]string) g.Node {
			return h.Div(
				ui.Signals(`{demoVal:""}`),
				Field(FieldProps{Label: "Your name", Hint: "Will be greeted"},
					Input(InputProps{
						Type:        "text",
						Placeholder: "e.g. Michael",
						Signal:      "demoVal",
					}),
				),
				h.P(h.Style("margin-top:8px;font-size:.85rem"), g.Text("Value: "), h.Span(g.Attr("data-text", `$demoVal || "(empty)"`))),
			)
		},
	})

	registry.Register(&registry.Component{
		Slug: "textarea", Name: "Textarea", Category: "form",
		Summary: "Multiline text input with Datastar signal binding.",
		Code: `form.Field(form.FieldProps{Label: "Message", Hint: "At least 10 characters"},
    form.Textarea(form.TextareaProps{
        Placeholder: "Write something…",
        Signal:      "message",
        Rows:        4,
        Required:    true,
    }),
)`,
		Render: func(p map[string]string) g.Node {
			return Field(FieldProps{Label: "Message", Hint: "At least 10 characters"},
				Textarea(TextareaProps{
					Placeholder: "Write something…",
					Signal:      "msg",
					Rows:        4,
				}),
			)
		},
	})

	registry.Register(&registry.Component{
		Slug: "select", Name: "Select", Category: "form",
		Summary: "Dropdown select with custom arrow, bound to a Datastar signal.",
		Code: `form.Field(form.FieldProps{Label: "Country"},
    form.Select(form.SelectProps{
        Signal: "country",
        Options: []form.SelectOption{
            {Value: "", Label: "Pick one…"},
            {Value: "at", Label: "Austria"},
            {Value: "de", Label: "Germany"},
        },
    }),
)`,
		Render: func(p map[string]string) g.Node {
			return h.Div(
				ui.Signals(`{fruit:""}`),
				Field(FieldProps{Label: "Favourite fruit"},
					Select(SelectProps{
						Signal: "fruit",
						Options: []SelectOption{
							{Value: "", Label: "Pick one…"},
							{Value: "apple", Label: "Apple"},
							{Value: "banana", Label: "Banana"},
							{Value: "cherry", Label: "Cherry"},
						},
					}),
				),
				h.P(h.Style("margin-top:8px;font-size:.85rem"), g.Text("Value: "), h.Span(g.Attr("data-text", `$fruit || "(none)"`))),
			)
		},
	})

	registry.Register(&registry.Component{
		Slug: "checkbox", Name: "Checkbox", Category: "form",
		Summary: "Styled checkbox bound to a boolean Datastar signal.",
		Code: `// declare signal in ancestor: ui.Signals("{agreed:false}")
form.Checkbox(form.CheckboxProps{
    Label:  "I accept the terms",
    Signal: "agreed",
})`,
		Render: func(p map[string]string) g.Node {
			return h.Div(
				ui.Signals(`{agreed:false}`),
				Checkbox(CheckboxProps{
					Label:  "I agree to the terms",
					Signal: "agreed",
				}),
				h.P(h.Style("margin-top:8px;font-size:.85rem"), g.Text("Agreed: "), h.Span(g.Attr("data-text", `$agreed ? "yes" : "no"`))),
			)
		},
	})

	registry.Register(&registry.Component{
		Slug: "radio", Name: "Radio Group", Category: "form",
		Summary: "Styled radio buttons via RadioGroup, bound to a shared signal.",
		Code: `form.RadioGroup(form.RadioGroupProps{
    Signal: "plan",
    Name:   "plan",
    Options: []form.RadioOption{
        {Value: "free",       Label: "Free"},
        {Value: "pro",        Label: "Pro"},
        {Value: "enterprise", Label: "Enterprise"},
    },
})`,
		Render: func(p map[string]string) g.Node {
			return h.Div(
				ui.Signals(`{plan:""}`),
				Field(FieldProps{Label: "Choose a plan"},
					RadioGroup(RadioGroupProps{
						Signal: "plan",
						Name:   "plan",
						Options: []RadioOption{
							{Value: "free", Label: "Free"},
							{Value: "pro", Label: "Pro"},
							{Value: "enterprise", Label: "Enterprise"},
						},
					}),
				),
				h.P(h.Style("margin-top:8px;font-size:.85rem"), g.Text("Plan: "), h.Span(g.Attr("data-text", `$plan || "(none)"`))),
			)
		},
	})

	registry.Register(&registry.Component{
		Slug: "contact-form", Name: "Contact Form", Category: "form",
		Summary: "Full contact form: fields, validation error display, altcha captcha, honeypot, SSE submit.",
		Code: `// See projects/homepage/handlers.go for the full SSE handler.
// Requires /api/altcha + /api/contact endpoints.`,
		Render: func(p map[string]string) g.Node {
			return h.Div(
				ui.Signals(`{name:'',email:'',message:'',sending:false,nameError:'',emailError:'',msgError:''}`),
				layout.Stack(layout.StackProps{},
					Field(FieldProps{Label: "Name", ErrorFor: "name"},
						Input(InputProps{Type: "text", Placeholder: "Your name", Signal: "name"}),
					),
					Field(FieldProps{Label: "Email", ErrorFor: "email"},
						Input(InputProps{Type: "email", Placeholder: "your@email.com", Signal: "email"}),
					),
					Field(FieldProps{Label: "Message", ErrorFor: "msg"},
						Textarea(TextareaProps{Placeholder: "What's on your mind?", Signal: "message", Rows: 4}),
					),
					primitive.Button(primitive.ButtonProps{
						Variant: token.Primary,
						Type:    "submit",
					}, g.Text("Send message")),
				),
			)
		},
	})

	registry.Register(&registry.Component{
		Slug: "switch", Name: "Switch", Category: "form",
		Summary: "Toggle switch bound to a boolean Datastar signal. Companion to Checkbox.",
		Code: `form.Switch(form.SwitchProps{
    Label:  "Enable notifications",
    Signal: "notifs",
})`,
		Render: func(p map[string]string) g.Node {
			return h.Div(
				ui.Signals(`{darkMode:false,autoSave:true,newsletter:false}`),
				layout.Stack(layout.StackProps{},
					Switch(SwitchProps{Label: "Dark mode", Signal: "darkMode"}),
					Switch(SwitchProps{Label: "Auto-save", Signal: "autoSave", Checked: true}),
					Switch(SwitchProps{Label: "Newsletter", Signal: "newsletter"}),
				),
			)
		},
	})
}
