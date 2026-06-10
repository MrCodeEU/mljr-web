package data

import (
	"mljr-web/ui/icon"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type TreeNode struct {
	Label    string
	Icon     string // optional lucide icon
	Href     string // if set, label is a link
	Children []TreeNode
	Open     bool // default open state for branch nodes
}

type TreeViewProps struct {
	Nodes []TreeNode
}

// TreeView renders a hierarchical collapsible tree using nested <details> elements.
// No JS required — native browser disclosure behavior.
func TreeView(p TreeViewProps) g.Node {
	return h.Div(
		g.Attr("data-component", "tree-view"),
		h.Ul(g.Attr("data-slot", "root"), treeNodes(p.Nodes, 0)),
	)
}

func treeNodes(nodes []TreeNode, depth int) g.Node {
	items := make([]g.Node, len(nodes))
	for i, n := range nodes {
		items[i] = treeNode(n, depth)
	}
	return g.Group(items)
}

func treeNode(n TreeNode, depth int) g.Node {
	if len(n.Children) == 0 {
		// Leaf node
		label := labelEl(n)
		return h.Li(g.Attr("data-slot", "leaf"), label)
	}

	// Branch node using <details>
	return h.Li(
		g.El("details",
			g.If(n.Open, g.Attr("open", "")),
			g.El("summary",
				g.Attr("data-slot", "branch"),
				icon.Icon("lucide:chevron-right", icon.Props{Size: "1rem"}),
				g.If(n.Icon != "", icon.Icon(n.Icon, icon.Props{Size: "1rem"})),
				g.Text(n.Label),
			),
			h.Ul(g.Attr("data-slot", "subtree"), treeNodes(n.Children, depth+1)),
		),
	)
}

func labelEl(n TreeNode) g.Node {
	inner := []g.Node{}
	if n.Icon != "" {
		inner = append(inner, icon.Icon(n.Icon, icon.Props{Size: "1rem"}))
	}
	inner = append(inner, g.Text(n.Label))

	if n.Href != "" {
		return h.A(
			g.Attr("data-slot", "label"),
			h.Href(n.Href),
			g.Group(inner),
		)
	}
	return h.Span(g.Attr("data-slot", "label"), g.Group(inner))
}
