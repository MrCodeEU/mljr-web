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
		Slug: "avatar", Name: "Avatar", Category: "primitive",
		Summary: "User avatar with image or initials fallback, optional status dot, and size/shape variants.",
		Code: `// with image
primitive.Avatar(primitive.AvatarProps{Src: "/img/user.jpg", Alt: "Alice", Size: token.MD})

// initials fallback with tone
primitive.Avatar(primitive.AvatarProps{Initials: "AB", Tone: token.ToneCyan, Status: primitive.AvatarOnline})

// square shape
primitive.Avatar(primitive.AvatarProps{Initials: "BC", Shape: primitive.AvatarSquare, Tone: token.ToneViolet})`,
		Controls: []registry.Control{
			{Name: "initials", Type: registry.ControlText, Default: "MR"},
			{Name: "size", Type: registry.ControlEnum, Options: []string{"sm", "md", "lg", "xl"}, Default: "md"},
			{Name: "shape", Type: registry.ControlEnum, Options: []string{"", "square"}, Default: ""},
			{Name: "status", Type: registry.ControlEnum, Options: []string{"", "online", "away", "offline"}, Default: "online"},
			{Name: "tone", Type: registry.ControlEnum, Options: []string{"", "cyan", "violet", "lime", "pink", "yellow"}, Default: "cyan"},
		},
		Render: func(p map[string]string) g.Node {
			return Avatar(AvatarProps{
				Initials: p["initials"],
				Size:     token.Size(p["size"]),
				Shape:    AvatarShape(p["shape"]),
				Status:   AvatarStatus(p["status"]),
				Tone:     token.Tone(p["tone"]),
			})
		},
		Examples: []registry.Example{
			{Title: "Sizes", Node: func() g.Node {
				return h.Div(h.Style("display:flex;gap:var(--sp-3);align-items:center;flex-wrap:wrap"),
					Avatar(AvatarProps{Initials: "SM", Size: "sm", Tone: token.ToneCyan}),
					Avatar(AvatarProps{Initials: "MD", Tone: token.ToneViolet}),
					Avatar(AvatarProps{Initials: "LG", Size: "lg", Tone: token.ToneLime}),
					Avatar(AvatarProps{Initials: "XL", Size: "xl", Tone: token.TonePink}),
				)
			}},
			{Title: "With status", Node: func() g.Node {
				return h.Div(h.Style("display:flex;gap:var(--sp-3);align-items:center;flex-wrap:wrap"),
					Avatar(AvatarProps{Initials: "ON", Tone: token.ToneCyan, Status: AvatarOnline}),
					Avatar(AvatarProps{Initials: "AW", Tone: token.ToneYellow, Status: AvatarAway}),
					Avatar(AvatarProps{Initials: "OF", Status: AvatarOffline}),
				)
			}},
		},
	})
}
