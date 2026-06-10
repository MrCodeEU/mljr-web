//go:build showcase

package layout

import (
	"mljr-web/ui/registry"

	g "maragu.dev/gomponents"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "breadcrumb", Name: "Breadcrumb", Category: "layout",
		Summary: "Navigation trail built from BreadcrumbItems. Last item is current page (no link).",
		Code: `layout.Breadcrumb(layout.BreadcrumbProps{},
    layout.BreadcrumbItem{Label: "Home",     Href: "/"},
    layout.BreadcrumbItem{Label: "Products", Href: "/products"},
    layout.BreadcrumbItem{Label: "Widget"},  // current page — no Href
)`,
		Render: func(p map[string]string) g.Node {
			return Breadcrumb(BreadcrumbProps{},
				BreadcrumbItem{Label: "Home", Href: "/"},
				BreadcrumbItem{Label: "Components", Href: "/components"},
				BreadcrumbItem{Label: "Breadcrumb"},
			)
		},
	})

	registry.Register(&registry.Component{
		Slug: "divider", Name: "Divider", Category: "layout",
		Summary: "Horizontal separator rule with optional centered text label.",
		Code: `layout.Divider(layout.DividerProps{})           // plain line
layout.Divider(layout.DividerProps{}, g.Text("OR")) // with label`,
		Controls: []registry.Control{
			{Name: "label", Type: registry.ControlText, Default: "OR"},
		},
		Render: func(p map[string]string) g.Node {
			if p["label"] != "" {
				return Divider(DividerProps{}, g.Text(p["label"]))
			}
			return Divider(DividerProps{})
		},
	})
}
