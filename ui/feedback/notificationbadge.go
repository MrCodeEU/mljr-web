package feedback

import (
	"fmt"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type NotificationBadgeProps struct {
	Count int    // 0 = hide badge unless Dot is true
	Max   int    // cap display (default 99); shows "99+" when exceeded
	Dot   bool   // show as dot regardless of count
	Color string // background color (default "var(--danger)")
}

// NotificationBadge wraps any content and overlays a count badge top-right.
// Use inside a position:relative container (or it creates its own).
func NotificationBadge(p NotificationBadgeProps, child g.Node) g.Node {
	if p.Color == "" {
		p.Color = "var(--danger)"
	}
	if p.Max == 0 {
		p.Max = 99
	}

	show := p.Dot || p.Count > 0
	if !show {
		return child
	}

	var label string
	if !p.Dot {
		if p.Count > p.Max {
			label = fmt.Sprintf("%d+", p.Max)
		} else {
			label = fmt.Sprintf("%d", p.Count)
		}
	}

	badgeStyle := "position:absolute;top:-4px;right:-4px;min-width:18px;height:18px;" +
		"background:" + p.Color + ";color:#fff;border-radius:9px;font-size:10px;" +
		"font-weight:800;display:flex;align-items:center;justify-content:center;" +
		"padding:0 4px;pointer-events:none;z-index:1"
	if p.Dot {
		badgeStyle = "position:absolute;top:-2px;right:-2px;width:10px;height:10px;" +
			"background:" + p.Color + ";border-radius:50%;pointer-events:none;z-index:1"
	}

	return h.Div(
		g.Attr("data-component", "notification-badge"),
		h.Style("position:relative;display:inline-flex"),
		child,
		h.Span(
			g.Attr("data-slot", "badge"),
			h.Style(badgeStyle),
			g.Attr("aria-label", label+" notifications"),
			g.If(!p.Dot, g.Text(label)),
		),
	)
}
