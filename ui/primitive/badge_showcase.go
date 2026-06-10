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
		Slug: "badge", Name: "Badge", Category: "primitive",
		Summary: "Small pill label for counts, statuses, or tags. Variants match semantic colors.",
		Code: `// status badge
primitive.Badge(primitive.BadgeProps{Variant: primitive.BadgeSuccess}, g.Text("Active"))

// count badge
primitive.Badge(primitive.BadgeProps{Variant: primitive.BadgeDanger}, g.Text("12"))

// toned badge
primitive.Badge(primitive.BadgeProps{Variant: primitive.BadgeTone, Tone: token.ToneCyan}, g.Text("Beta"))`,
		Controls: []registry.Control{
			{Name: "variant", Type: registry.ControlEnum, Options: []string{"", "success", "danger", "warning", "info", "outline", "tone"}, Default: ""},
			{Name: "tone", Type: registry.ControlEnum, Options: []string{"", "yellow", "cyan", "violet", "pink", "lime", "mint"}, Default: ""},
			{Name: "size", Type: registry.ControlEnum, Options: []string{"", "sm", "md", "lg"}, Default: ""},
			{Name: "label", Type: registry.ControlText, Default: "New"},
		},
		Render: func(p map[string]string) g.Node {
			return Badge(BadgeProps{
				Variant: BadgeVariant(p["variant"]),
				Tone:    token.Tone(p["tone"]),
				Size:    token.Size(p["size"]),
			}, g.Text(p["label"]))
		},
		Examples: []registry.Example{
			{Title: "All variants", Node: func() g.Node {
				return h.Div(h.Style("display:flex;gap:var(--sp-2);flex-wrap:wrap;align-items:center"),
					Badge(BadgeProps{}, g.Text("Default")),
					Badge(BadgeProps{Variant: BadgeSuccess}, g.Text("Active")),
					Badge(BadgeProps{Variant: BadgeDanger}, g.Text("99+")),
					Badge(BadgeProps{Variant: BadgeWarning}, g.Text("Warn")),
					Badge(BadgeProps{Variant: BadgeInfo}, g.Text("Info")),
					Badge(BadgeProps{Variant: BadgeOutline}, g.Text("Draft")),
				)
			}},
		},
	})
}
