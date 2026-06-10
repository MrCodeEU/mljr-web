package feedback

import (
	"mljr-web/ui/icon"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type EmptyStateProps struct {
	Icon    string // icon name (e.g. "lucide:inbox")
	Title   string
	Message string
	Attrs   []g.Node
}

// EmptyState renders a centred placeholder for empty data views.
// Pass action buttons (e.g. Button) as children.
func EmptyState(p EmptyStateProps, actions ...g.Node) g.Node {
	var iconNode g.Node
	if p.Icon != "" {
		iconNode = h.Div(g.Attr("data-slot", "icon"), icon.Icon(p.Icon))
	}

	var actionsNode g.Node
	if len(actions) > 0 {
		actionsNode = h.Div(g.Attr("data-slot", "actions"), g.Group(actions))
	}

	return h.Div(
		g.Attr("data-component", "empty-state"),
		g.Group(p.Attrs),
		iconNode,
		g.If(p.Title != "", h.Div(g.Attr("data-slot", "title"), g.Text(p.Title))),
		g.If(p.Message != "", h.P(g.Attr("data-slot", "message"), g.Text(p.Message))),
		actionsNode,
	)
}
