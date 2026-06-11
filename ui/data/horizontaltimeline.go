package data

import (
	"mljr-web/ui/primitive"
	"mljr-web/ui/token"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type HorizontalTimelineItem struct {
	Label   string
	Period  string
	Desc    string
	Tone    token.Tone
	Current bool
}

type HorizontalTimelineProps struct {
	// ScrollSnap enables CSS scroll-snap for swipe-friendly horizontal scrolling.
	ScrollSnap bool
}

// HorizontalTimeline renders events in a horizontal scrollable strip.
// Each event has a dot, vertical connector to a card, and label below.
// No JS — pure CSS scroll.
func HorizontalTimeline(p HorizontalTimelineProps, items ...HorizontalTimelineItem) g.Node {
	containerStyle := "display:flex;overflow-x:auto;padding-bottom:var(--sp-4);gap:0;align-items:flex-start;scrollbar-width:thin;scrollbar-color:var(--muted) transparent"
	if p.ScrollSnap {
		containerStyle += ";scroll-snap-type:x mandatory"
	}

	nodes := make([]g.Node, len(items))
	for i, item := range items {
		isLast := i == len(items)-1

		dotTone := item.Tone
		if dotTone == "" {
			dotTone = token.ToneAccent
		}

		snapStyle := ""
		if p.ScrollSnap {
			snapStyle = "scroll-snap-align:start;"
		}

		nodes[i] = h.Div(
			g.Attr("data-slot", "ht-item"),
			g.If(item.Current, g.Attr("data-current", "")),
			h.Style(snapStyle+"display:flex;flex-direction:column;align-items:center;min-width:180px;max-width:240px;flex-shrink:0;padding:0 var(--sp-2)"),
			// Top: horizontal line + dot
			h.Div(
				h.Style("display:flex;align-items:center;width:100%;margin-bottom:var(--sp-3)"),
				// Left line (hidden for first item)
				g.If(i > 0, h.Div(h.Style("flex:1;height:3px;background:var(--ink);opacity:.25"))),
				g.If(i == 0, h.Div(h.Style("flex:1"))),
				// Dot
				h.Div(
					g.Attr("data-slot", "ht-dot"),
					primitive.Tag(primitive.TagProps{Tone: dotTone},
						g.If(item.Current,
							g.Raw(`<svg width="8" height="8" viewBox="0 0 8 8"><circle cx="4" cy="4" r="4" fill="currentColor"/></svg>`),
						),
					),
				),
				// Right line (hidden for last item)
				g.If(!isLast, h.Div(h.Style("flex:1;height:3px;background:var(--ink);opacity:.25"))),
				g.If(isLast, h.Div(h.Style("flex:1"))),
			),
			// Card below dot
			h.Div(
				h.Style("width:100%"),
				primitive.Card(primitive.CardProps{Tone: item.Tone},
					g.If(item.Period != "", h.Div(h.Style("font-size:var(--t-xs);font-weight:800;opacity:.6;text-transform:uppercase;letter-spacing:.06em;margin-bottom:var(--sp-1)"), g.Text(item.Period))),
					h.Div(h.Style("font-weight:900;font-size:var(--t-sm);margin-bottom:var(--sp-1)"), g.Text(item.Label)),
					g.If(item.Desc != "", h.P(h.Style("font-size:var(--t-xs);opacity:.75;margin:0;line-height:1.4"), g.Text(item.Desc))),
				),
			),
		)
	}

	return h.Div(
		g.Attr("data-component", "horizontal-timeline"),
		h.Style(containerStyle),
		g.Group(nodes),
	)
}
