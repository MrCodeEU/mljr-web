//go:build showcase

package layout

import (
	"mljr-web/ui/primitive"
	"mljr-web/ui/registry"
	"mljr-web/ui/token"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "tabs", Name: "Tabs", Category: "layout",
		Summary: "Datastar-driven tabbed panels. Signal holds the active tab slug; switching is instant.",
		Code: `layout.Tabs(layout.TabsProps{Signal: "view", Default: "overview"},
    []layout.Tab{
        {Slug: "overview", Label: "Overview", Body: h.P(g.Text("Overview content."))},
        {Slug: "details",  Label: "Details",  Body: h.P(g.Text("Details content."))},
        {Slug: "code",     Label: "Code",     Body: h.P(g.Text("Code content."))},
    },
)`,
		Render: func(p map[string]string) g.Node {
			return Tabs(TabsProps{Signal: "demoTab"},
				[]Tab{
					{Slug: "one", Label: "Overview", Body: Stack(StackProps{},
						primitive.Heading(primitive.HeadingProps{Level: 3}, g.Text("Overview")),
						h.P(g.Text("High-level summary of the feature or content area.")),
					)},
					{Slug: "two", Label: "Details", Body: Stack(StackProps{},
						primitive.Heading(primitive.HeadingProps{Level: 3}, g.Text("Details")),
						h.P(g.Text("Technical specifications, parameters, and edge cases.")),
						primitive.Badge(primitive.BadgeProps{Variant: primitive.BadgeSuccess}, g.Text("Stable")),
					)},
					{Slug: "three", Label: "Code", Body: Stack(StackProps{},
						primitive.Heading(primitive.HeadingProps{Level: 3}, g.Text("Usage")),
						h.Pre(h.Style("font-family:var(--font-mono);font-size:var(--t-sm)"),
							h.Code(g.Text("layout.Tabs(TabsProps{}, []Tab{...})")),
						),
					)},
					{Slug: "four", Label: "History", Body: Stack(StackProps{},
						primitive.Card(primitive.CardProps{Tone: token.ToneMint},
							h.P(g.Text("v1.0 — initial release")),
						),
					)},
				},
			)
		},
	})
}
