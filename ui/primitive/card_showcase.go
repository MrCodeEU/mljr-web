//go:build showcase

package primitive

import (
	"mljr-web/ui/registry"
	"mljr-web/ui/token"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "card", Name: "Card", Category: "primitive",
		Summary: "Bordered, shadow-dropped container. Tone controls fill.",
		Code: `primitive.Card(primitive.CardProps{
    Tone: token.ToneCyan,
},
    primitive.Heading(primitive.HeadingProps{Level: 3}, g.Text("Title")),
    h.P(g.Text("Body content")),
)`,
		Controls: []registry.Control{
			{Name: "tone", Type: registry.ControlEnum, Options: []string{"", "yellow", "cyan", "violet", "pink", "lime", "mint", "sky", "blush", "accent", "accent-2"}, Default: ""},
			{Name: "interactive", Type: registry.ControlBool, Default: "false"},
		},
		Render: func(p map[string]string) g.Node {
			return Card(CardProps{
				Tone:        token.Tone(p["tone"]),
				Interactive: p["interactive"] == "true",
			},
				Heading(HeadingProps{Level: 3}, g.Text("Card title")),
				h.P(g.Text("Tonal surface, hard border, brutalist shadow.")),
			)
		},
	})
}
