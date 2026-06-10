//go:build showcase

package form

import (
	"mljr-web/ui/registry"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "date-input", Name: "Date Input", Category: "form",
		Summary: "Styled native date picker. Uses browser-native calendar popup. Same visual style as other inputs.",
		Code: `form.DateInput(form.DateInputProps{Name: "birthday", Min: "1900-01-01"})
form.TimeInput(form.TimeInputProps{Name: "time", Step: 900}) // 15-min steps`,
		Render: func(p map[string]string) g.Node {
			return h.Div(
				h.Style("display:flex;flex-direction:column;gap:var(--sp-5);max-width:360px"),
				Field(FieldProps{Label: "Date of birth"},
					DateInput(DateInputProps{Name: "birthday", Min: "1900-01-01", Max: "2010-01-01"}),
				),
				Field(FieldProps{Label: "Event date"},
					DateInput(DateInputProps{Name: "event", Value: "2026-12-31"}),
				),
				Field(FieldProps{Label: "Start time"},
					TimeInput(TimeInputProps{Name: "start", Value: "09:00", Step: 900}),
				),
				h.Div(
					h.Style("display:grid;grid-template-columns:1fr 1fr;gap:var(--sp-3)"),
					Field(FieldProps{Label: "From"},
						DateInput(DateInputProps{Name: "from"}),
					),
					Field(FieldProps{Label: "To"},
						DateInput(DateInputProps{Name: "to"}),
					),
				),
			)
		},
	})

	registry.Register(&registry.Component{
		Slug: "search-input", Name: "Search Input", Category: "form",
		Summary: "Search field with icon prefix, optional Datastar debounced @get, and clear button.",
		Code: `form.SearchInput(form.SearchInputProps{
    Name:      "q",
    Target:    "/api/search",
    Debounce:  "300ms",
    Clearable: true,
})`,
		Render: func(p map[string]string) g.Node {
			return h.Div(
				h.Style("display:flex;flex-direction:column;gap:var(--sp-4)"),
				h.Div(
					h.P(h.Style("font-size:var(--t-xs);opacity:.5;font-weight:700;text-transform:uppercase;letter-spacing:.06em;margin:0 0 var(--sp-2)"), g.Text("Basic")),
					SearchInput(SearchInputProps{Placeholder: "Search components…"}),
				),
				h.Div(
					h.P(h.Style("font-size:var(--t-xs);opacity:.5;font-weight:700;text-transform:uppercase;letter-spacing:.06em;margin:0 0 var(--sp-2)"), g.Text("With clear button")),
					SearchInput(SearchInputProps{Name: "q2", Placeholder: "Search anything…", Clearable: true}),
				),
			)
		},
	})
}
