//go:build showcase

package primitive

import (
	"mljr-web/ui/registry"
	"mljr-web/ui/token"

	g "maragu.dev/gomponents"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "tag", Name: "Tag", Category: "primitive",
		Summary: "Uppercase chip; colors via data-tone.",
		Code: `primitive.Tag(primitive.TagProps{Tone: token.ToneCyan}, g.Text("New"))
primitive.Tag(primitive.TagProps{Tone: token.ToneYellow}, g.Text("Draft"))
primitive.Tag(primitive.TagProps{}, g.Text("Default"))`,
		Controls: []registry.Control{
			{Name: "tone", Type: registry.ControlEnum, Options: []string{"", "yellow", "cyan", "violet", "pink", "lime", "mint", "sky", "blush"}, Default: "yellow"},
			{Name: "label", Type: registry.ControlText, Default: "New"},
		},
		Render: func(p map[string]string) g.Node {
			return Tag(TagProps{Tone: token.Tone(p["tone"])}, g.Text(p["label"]))
		},
	})
}
