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
		Slug: "tags-input", Name: "Tags Input", Category: "form",
		Summary: "Multi-value input — type and press Enter or comma to add, × to remove. Hidden input carries comma-separated values.",
		Code: `form.TagsInput(form.TagsInputProps{
    Name:    "tags",
    Default: []string{"Go", "Datastar"},
    Max:     10,
})`,
		Render: func(p map[string]string) g.Node {
			return h.Div(
				h.Style("display:flex;flex-direction:column;gap:var(--sp-5)"),
				primitive.Card(primitive.CardProps{},
					h.P(h.Style("font-size:var(--t-xs);text-transform:uppercase;font-weight:700;letter-spacing:.06em;opacity:.5;margin:0 0 var(--sp-3)"), g.Text("With pre-filled tags")),
					TagsInput(TagsInputProps{
						Name:    "tech",
						Default: []string{"Go", "Datastar", "Tailwind"},
					}),
				),
				primitive.Card(primitive.CardProps{},
					h.P(h.Style("font-size:var(--t-xs);text-transform:uppercase;font-weight:700;letter-spacing:.06em;opacity:.5;margin:0 0 var(--sp-3)"), g.Text("Empty — try typing and pressing Enter")),
					TagsInput(TagsInputProps{
						Name:        "skills",
						Placeholder: "Add a skill…",
						Max:         5,
					}),
					h.P(h.Style("font-size:var(--t-xs);opacity:.5;margin:var(--sp-2) 0 0"), g.Text("Max 5 tags")),
				),
			)
		},
	})
}
