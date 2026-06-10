//go:build showcase

package primitive

import (
	"mljr-web/ui/registry"
	"mljr-web/ui/token"

	g "maragu.dev/gomponents"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "button", Name: "Button", Category: "primitive",
		Summary: "Brutalist press button with variant/size/tone matrix.",
		Code: `import "mljr-web/ui/primitive"
import "mljr-web/ui/token"

// Basic
primitive.Button(primitive.ButtonProps{
    Variant: token.Primary,
    Size:    token.MD,
}, g.Text("Click me"))

// With tone
primitive.Button(primitive.ButtonProps{
    Variant: token.Tone,
    Tone:    token.ToneCyan,
}, g.Text("Cyan"))

// Icon size
primitive.Button(primitive.ButtonProps{
    Variant: token.Outline,
    Size:    token.SizeIcon,
}, icon.Icon("lucide:plus"))`,
		Controls: []registry.Control{
			{Name: "variant", Type: registry.ControlEnum, Options: []string{"primary", "secondary", "outline", "danger", "ghost", "tone"}, Default: "primary"},
			{Name: "size", Type: registry.ControlEnum, Options: []string{"sm", "md", "lg", "icon"}, Default: "md"},
			{Name: "tone", Type: registry.ControlEnum, Options: []string{"", "yellow", "cyan", "violet", "pink", "lime", "mint", "sky", "blush"}, Default: ""},
			{Name: "disabled", Type: registry.ControlBool, Default: "false"},
			{Name: "label", Type: registry.ControlText, Default: "Click me"},
		},
		Render: func(p map[string]string) g.Node {
			label := p["label"]
			if label == "" {
				label = "Click me"
			}
			return Button(ButtonProps{
				Variant:  token.Variant(p["variant"]),
				Size:     token.Size(p["size"]),
				Tone:     token.Tone(p["tone"]),
				Disabled: p["disabled"] == "true",
			}, g.Text(label))
		},
	})
}
