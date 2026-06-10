//go:build showcase

package overlay

import (
	"mljr-web/ui/primitive"
	"mljr-web/ui/registry"
	"mljr-web/ui/token"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "popover", Name: "Popover", Category: "overlay",
		Summary: "Positioned floating panel gated by a Datastar signal. Closes on outside-click.",
		Code: `overlay.Popover(
    overlay.PopoverProps{Signal: "info", Placement: overlay.PopoverBottom},
    primitive.Button(primitive.ButtonProps{Variant: token.Outline}, g.Text("More info")),
    h.Div(
        h.P(g.Text("Popover content goes here.")),
    ),
)`,
		Controls: []registry.Control{
			{Name: "placement", Type: registry.ControlEnum, Options: []string{"bottom", "top", "left", "right"}, Default: "bottom"},
		},
		Render: func(p map[string]string) g.Node {
			return h.Div(h.Style("display:flex;justify-content:center;padding:var(--sp-9)"),
				Popover(
					PopoverProps{Signal: "pop1", Placement: PopoverPlacement(p["placement"])},
					primitive.Button(primitive.ButtonProps{Variant: token.Outline}, g.Text("Open popover")),
					h.Div(
						primitive.Heading(primitive.HeadingProps{Level: 4}, g.Text("Popover title")),
						h.P(g.Text("Contextual content, a mini form, or extra details.")),
						primitive.Badge(primitive.BadgeProps{Variant: primitive.BadgeSuccess}, g.Text("New")),
					),
				),
			)
		},
	})
}
