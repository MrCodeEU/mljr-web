//go:build showcase

package primitive

import (
	"mljr-web/ui/registry"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "tooltip", Name: "Tooltip", Category: "primitive",
		Summary: "CSS-only hover tooltip. Four placement options. No JavaScript required.",
		Code: `primitive.Tooltip(
    primitive.TooltipProps{Text: "More info", Placement: primitive.TooltipTop},
    primitive.Button(primitive.ButtonProps{Variant: token.Outline}, g.Text("Hover me")),
)`,
		Controls: []registry.Control{
			{Name: "placement", Type: registry.ControlEnum, Options: []string{"top", "bottom", "left", "right"}, Default: "top"},
			{Name: "text", Type: registry.ControlText, Default: "Tooltip text"},
		},
		Render: func(p map[string]string) g.Node {
			return h.Div(h.Style("display:flex;justify-content:center;padding:var(--sp-7)"),
				Tooltip(TooltipProps{
					Text:      p["text"],
					Placement: TooltipPlacement(p["placement"]),
				},
					Button(ButtonProps{}, g.Text("Hover me")),
				),
			)
		},
	})

	registry.Register(&registry.Component{
		Slug: "kbd", Name: "Kbd", Category: "primitive",
		Summary: "Keyboard shortcut display. Monospace pill with press-depth border.",
		Code:    `h.Span(g.Text("Save with "), primitive.Kbd(g.Text("⌘")), g.Text(" + "), primitive.Kbd(g.Text("S")))`,
		Render: func(p map[string]string) g.Node {
			return h.Div(h.Style("display:flex;flex-wrap:wrap;gap:var(--sp-4);align-items:center"),
				h.Span(g.Text("Save: "), Kbd(g.Text("⌘")), g.Text("+"), Kbd(g.Text("S"))),
				h.Span(g.Text("Find: "), Kbd(g.Text("Ctrl")), g.Text("+"), Kbd(g.Text("F"))),
				h.Span(g.Text("New: "), Kbd(g.Text("Ctrl")), g.Text("+"), Kbd(g.Text("Shift")), g.Text("+"), Kbd(g.Text("N"))),
			)
		},
	})

	registry.Register(&registry.Component{
		Slug: "callout", Name: "Callout", Category: "primitive",
		Summary: "Left-bordered highlight block for tips, notes, and warnings. Variants match semantic colors.",
		Code: `primitive.Callout(primitive.CalloutProps{
    Variant: primitive.CalloutWarning,
    Title:   "Heads up",
}, g.Text("This action cannot be undone."))`,
		Controls: []registry.Control{
			{Name: "variant", Type: registry.ControlEnum, Options: []string{"", "info", "success", "warning", "danger"}, Default: "info"},
			{Name: "title", Type: registry.ControlText, Default: "Note"},
		},
		Render: func(p map[string]string) g.Node {
			return Callout(CalloutProps{
				Variant: CalloutVariant(p["variant"]),
				Title:   p["title"],
			}, h.P(g.Text("This is an important note that deserves special attention from the reader.")))
		},
	})
}
