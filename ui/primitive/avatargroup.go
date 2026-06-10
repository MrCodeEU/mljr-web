package primitive

import (
	"fmt"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type AvatarGroupProps struct {
	Max   int // max avatars shown before +N overflow badge (0 = no limit)
	Attrs []g.Node
}

// AvatarGroup stacks Avatar components with overlap and shows +N for extras.
func AvatarGroup(p AvatarGroupProps, avatars ...g.Node) g.Node {
	visible := avatars
	var overflow g.Node

	if p.Max > 0 && len(avatars) > p.Max {
		visible = avatars[:p.Max]
		extra := len(avatars) - p.Max
		overflow = h.Div(
			g.Attr("data-slot", "overflow"),
			g.Text(fmt.Sprintf("+%d", extra)),
		)
	}

	nodes := []g.Node{
		g.Attr("data-component", "avatar-group"),
		g.Group(p.Attrs),
	}
	nodes = append(nodes, visible...)
	if overflow != nil {
		nodes = append(nodes, overflow)
	}
	return h.Div(nodes...)
}
