//go:build showcase

package primitive

import (
	"fmt"
	"mljr-web/ui"
	"mljr-web/ui/registry"
	"mljr-web/ui/token"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "rating", Name: "Rating", Category: "primitive",
		Summary: "Star rating widget. Datastar signal holds current value. Read-only variant available.",
		Code: `// declare signal in ancestor
ui.Signals("{stars:0}")

primitive.Rating(primitive.RatingProps{Signal: "stars", Max: 5})`,
		Controls: []registry.Control{
			{Name: "max", Type: registry.ControlEnum, Options: []string{"3", "5", "10"}, Default: "5"},
			{Name: "readonly", Type: registry.ControlBool, Default: "false"},
		},
		Render: func(p map[string]string) g.Node {
			max := 5
			if p["max"] == "3" {
				max = 3
			} else if p["max"] == "10" {
				max = 10
			}
			return h.Div(
				ui.Signals(`{"ratingVal":3}`),
				Rating(RatingProps{Signal: "ratingVal", Max: max, ReadOnly: p["readonly"] == "true"}),
				h.P(h.Style("font-size:var(--t-sm);margin-top:var(--sp-2)"),
					g.Text("Rating: "), h.Span(g.Attr("data-text", `$ratingVal+" / `+fmt.Sprintf("%d", max)+`"`)),
				),
			)
		},
	})

	registry.Register(&registry.Component{
		Slug: "toggle-group", Name: "Toggle Group", Category: "primitive",
		Summary: "Segmented button control. Only one option active at a time, driven by a Datastar signal.",
		Code: `primitive.ToggleGroup(
    primitive.ToggleGroupProps{Signal: "view", Default: "grid"},
    primitive.ToggleOption{Value: "grid",  Label: "Grid",  Icon: "lucide:layout-grid"},
    primitive.ToggleOption{Value: "list",  Label: "List",  Icon: "lucide:list"},
    primitive.ToggleOption{Value: "table", Label: "Table", Icon: "lucide:table"},
)`,
		Render: func(p map[string]string) g.Node {
			return h.Div(h.Style("display:flex;flex-direction:column;gap:var(--sp-5)"),
				ToggleGroup(ToggleGroupProps{Signal: "view1", Default: "grid"},
					ToggleOption{Value: "grid", Label: "Grid", Icon: "lucide:layout-grid"},
					ToggleOption{Value: "list", Label: "List", Icon: "lucide:list"},
					ToggleOption{Value: "table", Label: "Table", Icon: "lucide:table"},
				),
				ToggleGroup(ToggleGroupProps{Signal: "align1", Default: "left"},
					ToggleOption{Value: "left", Icon: "lucide:align-left"},
					ToggleOption{Value: "center", Icon: "lucide:align-center"},
					ToggleOption{Value: "right", Icon: "lucide:align-right"},
				),
			)
		},
	})

	registry.Register(&registry.Component{
		Slug: "code-block", Name: "Code Block", Category: "primitive",
		Summary: "Styled pre/code display with language label and copy-to-clipboard button.",
		Code: `primitive.CodeBlock(primitive.CodeBlockProps{Language: "go", ID: "my-snippet"},
    "func Hello(name string) string {\n    return \"Hello, \"+name\n}",
)`,
		Render: func(p map[string]string) g.Node {
			return CodeBlock(CodeBlockProps{Language: "go", ID: "demo-cb"},
				`package main

import "fmt"

func main() {
    fmt.Println("Hello, mljr-web!")
}`)
		},
	})

	registry.Register(&registry.Component{
		Slug: "chip", Name: "Chip", Category: "primitive",
		Summary: "Dismissible pill tag. Click × removes from DOM or calls a custom Datastar expression.",
		Code: `primitive.Chip(primitive.ChipProps{Tone: token.ToneCyan, Dismiss: true}, g.Text("React"))
primitive.Chip(primitive.ChipProps{Tone: token.ToneViolet, Dismiss: true}, g.Text("TypeScript"))`,
		Controls: []registry.Control{
			{Name: "tone", Type: registry.ControlEnum, Options: []string{"", "cyan", "violet", "lime", "pink", "yellow"}, Default: "cyan"},
			{Name: "dismiss", Type: registry.ControlBool, Default: "true"},
		},
		Render: func(p map[string]string) g.Node {
			return h.Div(h.Style("display:flex;flex-wrap:wrap;gap:var(--sp-2)"),
				Chip(ChipProps{Tone: token.ToneCyan, Dismiss: p["dismiss"] == "true"}, g.Text("React")),
				Chip(ChipProps{Tone: token.ToneViolet, Dismiss: p["dismiss"] == "true"}, g.Text("TypeScript")),
				Chip(ChipProps{Tone: token.ToneLime, Dismiss: p["dismiss"] == "true"}, g.Text("Go")),
				Chip(ChipProps{Tone: token.ToneYellow, Dismiss: p["dismiss"] == "true"}, g.Text("Tailwind")),
				Chip(ChipProps{Tone: token.TonePink, Dismiss: p["dismiss"] == "true"}, g.Text("Datastar")),
			)
		},
	})

	registry.Register(&registry.Component{
		Slug: "progress-ring", Name: "Progress Ring", Category: "primitive",
		Summary: "SVG circular progress indicator. Server-rendered with correct stroke-dashoffset.",
		Code: `primitive.ProgressRing(primitive.ProgressRingProps{
    Value:     72,
    Size:      80,
    Thickness: 8,
    Label:     "72%",
    Variant:   primitive.ProgressSuccess,
})`,
		Controls: []registry.Control{
			{Name: "value", Type: registry.ControlEnum, Options: []string{"0", "25", "50", "72", "100"}, Default: "72"},
			{Name: "variant", Type: registry.ControlEnum, Options: []string{"", "success", "warning", "danger"}, Default: ""},
			{Name: "size", Type: registry.ControlEnum, Options: []string{"48", "72", "96"}, Default: "72"},
		},
		Render: func(p map[string]string) g.Node {
			val := map[string]int{"0": 0, "25": 25, "50": 50, "72": 72, "100": 100}[p["value"]]
			size := map[string]int{"48": 48, "72": 72, "96": 96}[p["size"]]
			label := fmt.Sprintf("%d%%", val)
			return ProgressRing(ProgressRingProps{
				Value:     val,
				Size:      size,
				Thickness: 7,
				Label:     label,
				Variant:   ProgressVariant(p["variant"]),
			})
		},
	})

	registry.Register(&registry.Component{
		Slug: "avatar-group", Name: "Avatar Group", Category: "primitive",
		Summary: "Stacked overlapping avatars with +N overflow badge for large sets.",
		Code: `primitive.AvatarGroup(primitive.AvatarGroupProps{Max: 4},
    primitive.Avatar(primitive.AvatarProps{Initials: "AB", Tone: token.ToneCyan}),
    primitive.Avatar(primitive.AvatarProps{Initials: "CD", Tone: token.ToneViolet}),
    primitive.Avatar(primitive.AvatarProps{Initials: "EF", Tone: token.ToneLime}),
    primitive.Avatar(primitive.AvatarProps{Initials: "GH", Tone: token.TonePink}),
    primitive.Avatar(primitive.AvatarProps{Initials: "IJ", Tone: token.ToneYellow}),
)`,
		Controls: []registry.Control{
			{Name: "max", Type: registry.ControlEnum, Options: []string{"2", "3", "4", "99"}, Default: "4"},
		},
		Render: func(p map[string]string) g.Node {
			max := map[string]int{"2": 2, "3": 3, "4": 4, "99": 99}[p["max"]]
			return AvatarGroup(AvatarGroupProps{Max: max},
				Avatar(AvatarProps{Initials: "AB", Tone: token.ToneCyan}),
				Avatar(AvatarProps{Initials: "CD", Tone: token.ToneViolet}),
				Avatar(AvatarProps{Initials: "EF", Tone: token.ToneLime}),
				Avatar(AvatarProps{Initials: "GH", Tone: token.TonePink}),
				Avatar(AvatarProps{Initials: "IJ", Tone: token.ToneYellow}),
				Avatar(AvatarProps{Initials: "KL"}),
			)
		},
	})
}
