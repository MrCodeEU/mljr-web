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
		Slug: "combobox", Name: "Combobox", Category: "form",
		Summary: "Filterable select — type to search, arrow keys to navigate, Enter to select. Hidden input carries the value.",
		Code: `form.Combobox(form.ComboboxProps{
    Name:    "language",
    Default: "go",
    Options: []form.ComboboxOption{
        {Value: "go",   Label: "Go"},
        {Value: "rust", Label: "Rust"},
        {Value: "ts",   Label: "TypeScript"},
    },
})`,
		Render: func(p map[string]string) g.Node {
			langs := []ComboboxOption{
				{Value: "go", Label: "Go"},
				{Value: "rust", Label: "Rust"},
				{Value: "ts", Label: "TypeScript"},
				{Value: "python", Label: "Python"},
				{Value: "kotlin", Label: "Kotlin"},
				{Value: "swift", Label: "Swift"},
				{Value: "csharp", Label: "C#"},
				{Value: "java", Label: "Java"},
				{Value: "elixir", Label: "Elixir"},
				{Value: "zig", Label: "Zig"},
			}
			countries := []ComboboxOption{
				{Value: "at", Label: "Austria"},
				{Value: "de", Label: "Germany"},
				{Value: "ch", Label: "Switzerland"},
				{Value: "fr", Label: "France"},
				{Value: "us", Label: "United States"},
				{Value: "jp", Label: "Japan"},
				{Value: "gb", Label: "United Kingdom"},
			}
			return h.Div(
				h.Style("display:flex;flex-direction:column;gap:var(--sp-5)"),
				primitive.Card(primitive.CardProps{},
					h.P(h.Style("font-size:var(--t-xs);text-transform:uppercase;font-weight:700;letter-spacing:.06em;opacity:.5;margin:0 0 var(--sp-3)"), g.Text("Language")),
					Combobox(ComboboxProps{Name: "language", Default: "go", Options: langs}),
				),
				primitive.Card(primitive.CardProps{},
					h.P(h.Style("font-size:var(--t-xs);text-transform:uppercase;font-weight:700;letter-spacing:.06em;opacity:.5;margin:0 0 var(--sp-3)"), g.Text("Country")),
					Combobox(ComboboxProps{Name: "country", Placeholder: "Search country…", Options: countries}),
				),
			)
		},
	})
}
