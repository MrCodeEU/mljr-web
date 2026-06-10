//go:build showcase

package primitive

import (
	"mljr-web/ui/registry"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "progress", Name: "Progress", Category: "primitive",
		Summary: "Horizontal progress bar with value, optional label, semantic variant colors, and three sizes.",
		Code: `primitive.Progress(primitive.ProgressProps{
    Value:     72,
    Label:     "Uploading…",
    ShowLabel: true,
    Variant:   primitive.ProgressSuccess,
})`,
		Controls: []registry.Control{
			{Name: "value", Type: registry.ControlEnum, Options: []string{"0", "25", "50", "72", "100"}, Default: "72"},
			{Name: "variant", Type: registry.ControlEnum, Options: []string{"", "success", "warning", "danger"}, Default: ""},
			{Name: "size", Type: registry.ControlEnum, Options: []string{"sm", "md", "lg"}, Default: "md"},
			{Name: "label", Type: registry.ControlBool, Default: "true"},
		},
		Render: func(p map[string]string) g.Node {
			val := 72
			switch p["value"] {
			case "0":
				val = 0
			case "25":
				val = 25
			case "50":
				val = 50
			case "100":
				val = 100
			}
			return Progress(ProgressProps{
				Value:     val,
				Label:     "Progress",
				ShowLabel: p["label"] == "true",
				Variant:   ProgressVariant(p["variant"]),
				Size:      p["size"],
			})
		},
		Examples: []registry.Example{
			{Title: "All variants", Node: func() g.Node {
				return h.Div(h.Style("display:flex;flex-direction:column;gap:var(--sp-3);width:100%"),
					Progress(ProgressProps{Value: 80, Label: "Default", ShowLabel: true}),
					Progress(ProgressProps{Value: 65, Label: "Success", ShowLabel: true, Variant: ProgressSuccess}),
					Progress(ProgressProps{Value: 45, Label: "Warning", ShowLabel: true, Variant: ProgressWarning}),
					Progress(ProgressProps{Value: 30, Label: "Danger", ShowLabel: true, Variant: ProgressDanger}),
				)
			}},
			{Title: "Sizes", Node: func() g.Node {
				return h.Div(h.Style("display:flex;flex-direction:column;gap:var(--sp-4);width:100%"),
					Progress(ProgressProps{Value: 60, Size: "sm"}),
					Progress(ProgressProps{Value: 60}),
					Progress(ProgressProps{Value: 60, Size: "lg"}),
				)
			}},
		},
	})
}
