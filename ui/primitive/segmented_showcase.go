//go:build showcase

package primitive

import (
	"mljr-web/ui/registry"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "segmented", Name: "Segmented Control", Category: "primitive",
		Summary: "Radio group styled as a connected button bar. Keyboard accessible via native radio inputs.",
		Code: `primitive.Segmented(primitive.SegmentedProps{
    Name:    "view",
    Default: "week",
    Options: []primitive.SegmentedOption{
        {Value: "day",   Label: "Day"},
        {Value: "week",  Label: "Week"},
        {Value: "month", Label: "Month"},
    },
})`,
		Render: func(p map[string]string) g.Node {
			return h.Div(
				h.Style("display:flex;flex-direction:column;gap:var(--sp-5);align-items:flex-start"),
				Segmented(SegmentedProps{
					Name:    "view",
					Default: "week",
					Options: []SegmentedOption{
						{Value: "day", Label: "Day"},
						{Value: "week", Label: "Week"},
						{Value: "month", Label: "Month"},
					},
				}),
				Segmented(SegmentedProps{
					Name:    "size",
					Default: "md",
					Options: []SegmentedOption{
						{Value: "sm", Label: "SM"},
						{Value: "md", Label: "MD"},
						{Value: "lg", Label: "LG"},
						{Value: "xl", Label: "XL"},
					},
				}),
				Segmented(SegmentedProps{
					Name:    "align",
					Default: "left",
					Options: []SegmentedOption{
						{Value: "left", Label: "Left"},
						{Value: "center", Label: "Center"},
						{Value: "right", Label: "Right"},
					},
				}),
			)
		},
	})
}
