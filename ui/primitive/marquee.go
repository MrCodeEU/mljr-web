package primitive

import (
	"fmt"
	"strings"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type MarqueeProps struct {
	Speed        string // CSS animation duration, default "25s"
	Direction    string // "left" (default) | "right"
	PauseOnHover bool
	Gap          string // gap between items, default "var(--sp-6)"
}

// Marquee renders items in an infinite horizontal scroll loop.
// Items are duplicated server-side to create a seamless CSS-only loop.
// No JS required; gap/speed/direction all CSS-driven.
func Marquee(p MarqueeProps, items ...g.Node) g.Node {
	if p.Speed == "" {
		p.Speed = "25s"
	}
	if p.Gap == "" {
		p.Gap = "var(--sp-6)"
	}
	if p.Direction == "" {
		p.Direction = "left"
	}

	attrs := []g.Node{
		g.Attr("data-component", "marquee"),
		g.If(p.Direction == "right", g.Attr("data-direction", "right")),
		g.If(p.PauseOnHover, g.Attr("data-pause-on-hover", "true")),
	}

	style := fmt.Sprintf("--marquee-speed:%s;--marquee-gap:%s", p.Speed, p.Gap)
	attrs = append(attrs, h.Style(style))

	// Build two copies of items for seamless loop
	var sb strings.Builder
	_ = sb // silence unused warning

	setA := make([]g.Node, len(items))
	setB := make([]g.Node, len(items))
	for i, item := range items {
		setA[i] = h.Div(g.Attr("data-slot", "item"), item)
		setB[i] = h.Div(g.Attr("aria-hidden", "true"), g.Attr("data-slot", "item"), item)
	}

	track := h.Div(
		g.Attr("data-slot", "track"),
		g.Group(setA),
		g.Group(setB),
	)

	return h.Div(append(attrs, track)...)
}
