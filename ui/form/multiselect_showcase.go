//go:build showcase

package form

import (
	"mljr-web/ui/registry"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "multi-select", Name: "Multi Select", Category: "form",
		Summary: "Chip-based multi-value selector with dropdown. Selected values displayed as removable chips; hidden inputs carry values for form submission.",
		Code: `// import "mljr-web/ui/form"
form.MultiSelect(form.MultiSelectProps{
    Name: "tags",
    Options: []form.MultiSelectOption{
        {Value: "go", Label: "Go"},
        {Value: "ts", Label: "TypeScript"},
    },
    Default: []string{"go"},
    Max:     3,
    Placeholder: "Select languages…",
})`,
		Render: func(p map[string]string) g.Node {
			return h.Div(
				h.Style("display:flex;flex-direction:column;gap:var(--sp-6)"),
				h.Div(
					h.H4(h.Style("font-size:var(--t-sm);font-weight:700;margin-bottom:var(--sp-2)"), g.Text("Languages (max 3)")),
					MultiSelect(MultiSelectProps{
						Name: "languages",
						Options: []MultiSelectOption{
							{Value: "go", Label: "Go"},
							{Value: "ts", Label: "TypeScript"},
							{Value: "py", Label: "Python"},
							{Value: "rs", Label: "Rust"},
							{Value: "kt", Label: "Kotlin"},
							{Value: "swift", Label: "Swift"},
						},
						Default:     []string{"go"},
						Max:         3,
						Placeholder: "Select languages…",
					}),
				),
				h.Div(
					h.H4(h.Style("font-size:var(--t-sm);font-weight:700;margin-bottom:var(--sp-2)"), g.Text("Topics (unlimited)")),
					MultiSelect(MultiSelectProps{
						Name: "topics",
						Options: []MultiSelectOption{
							{Value: "frontend", Label: "Frontend"},
							{Value: "backend", Label: "Backend"},
							{Value: "devops", Label: "DevOps"},
							{Value: "ml", Label: "Machine Learning"},
							{Value: "mobile", Label: "Mobile"},
							{Value: "security", Label: "Security"},
						},
						Placeholder: "Pick topics…",
					}),
				),
			)
		},
	})
}
