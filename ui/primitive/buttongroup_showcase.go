//go:build showcase

package primitive

import (
	"mljr-web/ui/icon"
	"mljr-web/ui/registry"
	"mljr-web/ui/token"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "button-group", Name: "Button Group", Category: "primitive",
		Summary: "Horizontal group of related buttons. Attached mode collapses shared borders into a connected strip.",
		Code: `// Attached (connected border)
primitive.ButtonGroup(primitive.ButtonGroupProps{Attached: true, Label: "Text alignment"},
    primitive.Button(primitive.ButtonProps{Variant: token.Outline}, icon.Icon("lucide:align-left")),
    primitive.Button(primitive.ButtonProps{Variant: token.Primary}, icon.Icon("lucide:align-center")),
    primitive.Button(primitive.ButtonProps{Variant: token.Outline}, icon.Icon("lucide:align-right")),
)`,
		Render: func(p map[string]string) g.Node {
			return h.Div(
				h.Style("display:flex;flex-direction:column;gap:var(--sp-5)"),
				h.Div(
					h.P(h.Style("font-size:var(--t-xs);opacity:.5;font-weight:700;text-transform:uppercase;letter-spacing:.06em;margin:0 0 var(--sp-2)"), g.Text("Attached")),
					ButtonGroup(ButtonGroupProps{Attached: true, Label: "Text alignment"},
						Button(ButtonProps{Variant: token.Primary, Size: token.SizeIcon}, icon.Icon("lucide:align-left")),
						Button(ButtonProps{Variant: token.Outline, Size: token.SizeIcon}, icon.Icon("lucide:align-center")),
						Button(ButtonProps{Variant: token.Outline, Size: token.SizeIcon}, icon.Icon("lucide:align-right")),
					),
				),
				h.Div(
					h.P(h.Style("font-size:var(--t-xs);opacity:.5;font-weight:700;text-transform:uppercase;letter-spacing:.06em;margin:0 0 var(--sp-2)"), g.Text("Spaced")),
					ButtonGroup(ButtonGroupProps{Label: "Actions"},
						Button(ButtonProps{Variant: token.Outline}, g.Text("Export")),
						Button(ButtonProps{Variant: token.Outline}, g.Text("Share")),
						Button(ButtonProps{Variant: token.Danger, Size: token.SizeSM}, g.Text("Delete")),
					),
				),
				h.Div(
					h.P(h.Style("font-size:var(--t-xs);opacity:.5;font-weight:700;text-transform:uppercase;letter-spacing:.06em;margin:0 0 var(--sp-2)"), g.Text("Pagination-style attached")),
					ButtonGroup(ButtonGroupProps{Attached: true},
						Button(ButtonProps{Variant: token.Outline, Size: token.SizeIcon}, icon.Icon("lucide:chevron-left")),
						Button(ButtonProps{Variant: token.Primary}, g.Text("1")),
						Button(ButtonProps{Variant: token.Outline}, g.Text("2")),
						Button(ButtonProps{Variant: token.Outline}, g.Text("3")),
						Button(ButtonProps{Variant: token.Outline, Size: token.SizeIcon}, icon.Icon("lucide:chevron-right")),
					),
				),
			)
		},
	})
}
