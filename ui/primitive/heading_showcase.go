//go:build showcase

package primitive

import (
	"strconv"

	"mljr-web/ui/registry"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "heading", Name: "Heading", Category: "primitive",
		Summary: "h1..h5 with Swiss display type and tracking.",
		Code: `// import "mljr-web/ui/primitive"
primitive.Heading(primitive.HeadingProps{Level: 1}, g.Text("Page title"))
primitive.Heading(primitive.HeadingProps{Level: 2}, g.Text("Section"))
primitive.Heading(primitive.HeadingProps{Level: 3}, g.Text("Subsection"))`,
		Controls: []registry.Control{
			{Name: "level", Type: registry.ControlEnum, Options: []string{"1", "2", "3", "4", "5"}, Default: "2"},
			{Name: "text", Type: registry.ControlText, Default: "The quick brown fox"},
		},
		Render: func(p map[string]string) g.Node {
			lvl, _ := strconv.Atoi(p["level"])
			return Heading(HeadingProps{Level: lvl}, g.Text(p["text"]))
		},
	})
	registry.Register(&registry.Component{
		Slug: "display", Name: "Display", Category: "primitive",
		Summary: "Hero display headline with highlighter <em> and accent <mark>.",
		Code: `primitive.Display(primitive.DisplayProps{},
    g.Text("Build something "),
    h.Em(g.Text("bold")),   // highlight strip
)`,
		Render: func(p map[string]string) g.Node {
			return Display(DisplayProps{},
				g.Text("Building a "),
				h.Em(g.Text("brutalist")),
				g.Text(" Go web stack."),
			)
		},
	})
}
