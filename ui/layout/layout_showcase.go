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
		Slug: "container", Name: "Container", Category: "layout",
		Summary: "Max-width content wrapper with responsive horizontal padding.",
		Code: `layout.Container(layout.ContainerProps{},
    // page content here
)`,
		Render: func(p map[string]string) g.Node {
			return h.Div(h.Style("width:100%;background:var(--surface-2,var(--surface));border:var(--bw-1) dashed var(--line)"),
				Container(ContainerProps{},
					h.P(g.Text("Content lives here — constrained and padded.")),
				),
			)
		},
	})

	registry.Register(&registry.Component{
		Slug: "stack", Name: "Stack", Category: "layout",
		Summary: "Flex column (v) or row (h) with configurable gap, align, and justify.",
		Code: `layout.Stack(layout.StackProps{Axis: "h", Gap: "md"},
    primitive.Button(primitive.ButtonProps{Variant: token.Primary}, g.Text("A")),
    primitive.Button(primitive.ButtonProps{Variant: token.Outline}, g.Text("B")),
)`,
		Controls: []registry.Control{
			{Name: "axis", Type: registry.ControlEnum, Options: []string{"v", "h"}, Default: "h"},
			{Name: "gap", Type: registry.ControlEnum, Options: []string{"sm", "md", "lg"}, Default: "md"},
		},
		Render: func(p map[string]string) g.Node {
			return Stack(StackProps{Axis: p["axis"], Gap: p["gap"]},
				primitive.Button(primitive.ButtonProps{Variant: token.Primary}, g.Text("Alpha")),
				primitive.Button(primitive.ButtonProps{Variant: token.Outline}, g.Text("Beta")),
				primitive.Button(primitive.ButtonProps{Variant: token.Ghost}, g.Text("Gamma")),
			)
		},
	})

	registry.Register(&registry.Component{
		Slug: "grid", Name: "Grid / Col", Category: "layout",
		Summary: "12-column CSS grid. Col data-span controls width.",
		Code: `layout.Grid(layout.GridProps{},
    layout.Col(layout.ColProps{Span: 8}, mainContent),
    layout.Col(layout.ColProps{Span: 4}, sidebar),
)`,
		Render: func(p map[string]string) g.Node {
			return h.Div(h.Style("width:100%"),
				Grid(GridProps{},
					Col(ColProps{Span: 4},
						primitive.Card(primitive.CardProps{Tone: token.ToneCyan},
							h.P(g.Text("span 4")),
						),
					),
					Col(ColProps{Span: 4},
						primitive.Card(primitive.CardProps{Tone: token.ToneLime},
							h.P(g.Text("span 4")),
						),
					),
					Col(ColProps{Span: 4},
						primitive.Card(primitive.CardProps{Tone: token.ToneViolet},
							h.P(g.Text("span 4")),
						),
					),
					Col(ColProps{Span: 8},
						primitive.Card(primitive.CardProps{Tone: token.TonePink},
							h.P(g.Text("span 8")),
						),
					),
					Col(ColProps{Span: 4},
						primitive.Card(primitive.CardProps{Tone: token.ToneYellow},
							h.P(g.Text("span 4")),
						),
					),
				),
			)
		},
	})

	registry.Register(&registry.Component{
		Slug: "navbar", Name: "Navbar", Category: "layout",
		Summary: "Page header with brand, nav links, and actions slots.",
		Code: `layout.Navbar(layout.NavbarProps{},
    g.Text("MyApp"),
    g.Group{h.A(h.Href("/about"), g.Text("About"))},
    g.Group{primitive.Button(...)},
)`,
		Render: func(p map[string]string) g.Node {
			return Navbar(NavbarProps{},
				h.Strong(g.Text("MyApp")),
				g.Group{
					h.A(h.Href("#"), g.Text("Home")),
					h.A(h.Href("#"), g.Text("About")),
					h.A(h.Href("#"), g.Text("Work")),
				},
				g.Group{
					primitive.Button(primitive.ButtonProps{Variant: token.Primary, Size: token.SizeSM}, g.Text("Contact")),
				},
			)
		},
	})

	registry.Register(&registry.Component{
		Slug: "footer", Name: "Footer", Category: "layout",
		Summary: "Page footer container. Add columns, links, copyright via children.",
		Code: `layout.Footer(layout.FooterProps{},
    h.Div(g.Text("© 2025 MyApp")),
)`,
		Render: func(p map[string]string) g.Node {
			return Footer(FooterProps{},
				h.Div(
					g.Attr("style", "display:flex;justify-content:space-between;align-items:center;flex-wrap:wrap;gap:var(--sp-3)"),
					h.Div(g.Text("© 2025 mljr-web")),
					h.Div(
						g.Attr("style", "display:flex;gap:var(--sp-4)"),
						h.A(h.Href("#"), g.Text("Privacy")),
						h.A(h.Href("#"), g.Text("Terms")),
						h.A(h.Href("#"), g.Text("Contact")),
					),
				),
			)
		},
	})
}
