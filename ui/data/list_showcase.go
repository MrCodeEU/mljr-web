//go:build showcase

package data

import (
	"mljr-web/ui/primitive"
	"mljr-web/ui/registry"
	"mljr-web/ui/token"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "list", Name: "List", Category: "data",
		Summary: "Styled unordered or ordered list. Divided variant adds row separators.",
		Code: `data.List(data.ListProps{Variant: data.ListDivided},
    data.ListItem(g.Text("First item")),
    data.ListItem(g.Text("Second item")),
    data.ListItem(g.Text("Third item")),
)`,
		Controls: []registry.Control{
			{Name: "variant", Type: registry.ControlEnum, Options: []string{"", "divided", "ordered"}, Default: "divided"},
		},
		Render: func(p map[string]string) g.Node {
			return List(ListProps{Variant: ListVariant(p["variant"])},
				ListItem(h.Strong(g.Text("Alice Müller")), h.Span(g.Attr("style", "opacity:.6;margin-left:var(--sp-2)"), g.Text("Engineer"))),
				ListItem(h.Strong(g.Text("Bob Chen")), h.Span(g.Attr("style", "opacity:.6;margin-left:var(--sp-2)"), g.Text("Designer"))),
				ListItem(h.Strong(g.Text("Carol Singh")), h.Span(g.Attr("style", "opacity:.6;margin-left:var(--sp-2)"), g.Text("PM"))),
				ListItem(h.Strong(g.Text("Dave Park")), h.Span(g.Attr("style", "opacity:.6;margin-left:var(--sp-2)"), g.Text("DevOps"))),
			)
		},
	})

	registry.Register(&registry.Component{
		Slug: "description-list", Name: "Description List", Category: "data",
		Summary: "Term/definition pairs in an aligned grid layout.",
		Code: `data.DescriptionList(data.DescriptionListProps{},
    data.DescriptionItem{Term: "Version",  Desc: g.Text("v1.4.2")},
    data.DescriptionItem{Term: "Released", Desc: g.Text("2025-03-14")},
    data.DescriptionItem{Term: "License",  Desc: g.Text("MIT")},
)`,
		Render: func(p map[string]string) g.Node {
			return DescriptionList(DescriptionListProps{},
				DescriptionItem{Term: "Version", Desc: primitive.Badge(primitive.BadgeProps{Variant: primitive.BadgeSuccess}, g.Text("v1.4.2"))},
				DescriptionItem{Term: "Released", Desc: g.Text("2025-03-14")},
				DescriptionItem{Term: "License", Desc: g.Text("MIT")},
				DescriptionItem{Term: "Maintainer", Desc: g.Text("Michael Reinegger")},
				DescriptionItem{Term: "Status", Desc: primitive.Badge(primitive.BadgeProps{Variant: primitive.BadgeInfo}, g.Text("Active"))},
				DescriptionItem{Term: "Theme", Desc: primitive.Tag(primitive.TagProps{Tone: token.ToneCyan}, g.Text("swissbrut"))},
			)
		},
	})
}
