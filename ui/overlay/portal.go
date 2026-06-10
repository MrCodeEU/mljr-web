package overlay

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

// Portal is the fixed-position mount region for overlays. Drop one in PageShell.
// Modals/toasts patch into #portal via SSE PatchElements or render inline with data-show.
func Portal(id string) g.Node {
	if id == "" {
		id = "portal"
	}
	return h.Div(
		h.ID(id),
		g.Attr("data-component", "portal"),
	)
}
