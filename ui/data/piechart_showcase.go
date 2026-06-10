//go:build showcase

package data

import (
	"mljr-web/ui/primitive"
	"mljr-web/ui/registry"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "pie-chart", Name: "Pie Chart", Category: "data",
		Summary: "Pure SVG solid pie chart using arc paths. Legend with percentage labels. No JS.",
		Code: `data.PieChart(data.PieChartProps{
    Slices: []data.DonutSlice{
        {Label: "Go",    Value: 45, Color: "#00ADD8"},
        {Label: "Other", Value: 55},
    },
})`,
		Render: func(p map[string]string) g.Node {
			return h.Div(
				h.Style("display:flex;flex-wrap:wrap;gap:var(--sp-5)"),
				primitive.Card(primitive.CardProps{},
					h.P(h.Style("font-size:var(--t-xs);text-transform:uppercase;font-weight:700;letter-spacing:.06em;opacity:.5;margin:0 0 var(--sp-3)"), g.Text("Traffic sources")),
					PieChart(PieChartProps{
						Slices: []DonutSlice{
							{Label: "Organic", Value: 42},
							{Label: "Direct", Value: 28},
							{Label: "Social", Value: 18},
							{Label: "Email", Value: 12},
						},
					}),
				),
				primitive.Card(primitive.CardProps{},
					h.P(h.Style("font-size:var(--t-xs);text-transform:uppercase;font-weight:700;letter-spacing:.06em;opacity:.5;margin:0 0 var(--sp-3)"), g.Text("OS breakdown")),
					PieChart(PieChartProps{
						Size: 140,
						Slices: []DonutSlice{
							{Label: "Linux", Value: 55, Color: "#CE422B"},
							{Label: "macOS", Value: 30, Color: "#555"},
							{Label: "Windows", Value: 15, Color: "#0078D4"},
						},
					}),
				),
			)
		},
	})

	registry.Register(&registry.Component{
		Slug: "tree-view", Name: "Tree View", Category: "data",
		Summary: "Hierarchical collapsible tree using native <details> elements. No JS, keyboard accessible.",
		Code: `data.TreeView(data.TreeViewProps{
    Nodes: []data.TreeNode{
        {Label: "src", Icon: "lucide:folder", Open: true, Children: []data.TreeNode{
            {Label: "main.go", Icon: "lucide:file-text"},
            {Label: "ui", Icon: "lucide:folder", Children: []data.TreeNode{
                {Label: "button.go"},
            }},
        }},
    },
})`,
		Render: func(p map[string]string) g.Node {
			return TreeView(TreeViewProps{
				Nodes: []TreeNode{
					{Label: "mljr-web", Icon: "lucide:folder", Open: true, Children: []TreeNode{
						{Label: "projects", Icon: "lucide:folder", Open: true, Children: []TreeNode{
							{Label: "homepage", Icon: "lucide:folder", Children: []TreeNode{
								{Label: "main.go", Icon: "lucide:file-text"},
								{Label: "pages", Icon: "lucide:folder", Children: []TreeNode{
									{Label: "home.go", Icon: "lucide:file-text"},
									{Label: "contact.go", Icon: "lucide:file-text"},
								}},
							}},
							{Label: "showcase", Icon: "lucide:folder", Children: []TreeNode{
								{Label: "main.go", Icon: "lucide:file-text"},
							}},
						}},
						{Label: "ui", Icon: "lucide:folder", Open: true, Children: []TreeNode{
							{Label: "primitive", Icon: "lucide:folder", Children: []TreeNode{
								{Label: "button.go", Icon: "lucide:file-text"},
								{Label: "card.go", Icon: "lucide:file-text"},
								{Label: "badge.go", Icon: "lucide:file-text"},
							}},
							{Label: "form", Icon: "lucide:folder", Children: []TreeNode{
								{Label: "input.go", Icon: "lucide:file-text"},
								{Label: "select.go", Icon: "lucide:file-text"},
							}},
							{Label: "css", Icon: "lucide:folder", Children: []TreeNode{
								{Label: "_primitive.css", Icon: "lucide:file-text"},
								{Label: "_form.css", Icon: "lucide:file-text"},
							}},
						}},
						{Label: "go.mod", Icon: "lucide:file-text"},
						{Label: "Makefile", Icon: "lucide:file-text"},
					}},
				},
			})
		},
	})
}
