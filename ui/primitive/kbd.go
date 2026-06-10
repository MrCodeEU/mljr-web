package primitive

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

// Kbd renders a keyboard shortcut key in monospace pill style.
func Kbd(children ...g.Node) g.Node {
	return h.Kbd(g.Attr("data-component", "kbd"), g.Group(children))
}
