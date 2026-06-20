package layout

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type NavPillItem struct {
	Label  string
	Href   string
	Active bool
}

// NavPills renders a row of link pills for multi-page (non-Datastar) nav,
// styled like the tab-btn/tab-list pair in _layout.css but using plain <a>
// tags since each item navigates to a different page rather than switching
// a client-side panel.
func NavPills(items []NavPillItem) g.Node {
	links := make([]g.Node, len(items))
	for i, it := range items {
		state := ""
		if it.Active {
			state = "active"
		}
		links[i] = h.A(
			g.Attr("data-component", "nav-pill"),
			g.Attr("data-state", state),
			h.Href(it.Href),
			g.Text(it.Label),
		)
	}
	return h.Nav(g.Attr("data-component", "nav-pills"), g.Group(links))
}
