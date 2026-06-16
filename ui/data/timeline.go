package data

import (
	"mljr-web/ui/primitive"
	"mljr-web/ui/token"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type TimelineItem struct {
	Period   string
	Title    string
	Org      string
	OrgLogo  string // optional img src — rendered absolute top-right of card
	Desc     string
	Tags     []string
	TagNodes []g.Node // pre-built tag nodes with icons/tones; overrides Tags when set
	Tone     token.Tone
	Attrs    []g.Node
}

type TimelineProps struct {
	Attrs []g.Node
}

func Timeline(p TimelineProps, items ...TimelineItem) g.Node {
	nodes := make([]g.Node, 0, len(items)+2)
	nodes = append(nodes, g.Attr("data-component", "timeline"))
	nodes = append(nodes, g.Group(p.Attrs))
	for _, item := range items {
		nodes = append(nodes, timelineItem(item))
	}
	return h.Div(nodes...)
}

func timelineItem(item TimelineItem) g.Node {
	tagNodes := item.TagNodes
	if len(tagNodes) == 0 {
		tagNodes = make([]g.Node, 0, len(item.Tags))
		for _, t := range item.Tags {
			tagNodes = append(tagNodes, primitive.Tag(primitive.TagProps{}, g.Text(t)))
		}
	}

	cardChildren := []g.Node{
		// Logo: absolute top-right, rendered via data-slot="logo"
		g.If(item.OrgLogo != "",
			h.Img(
				g.Attr("data-slot", "logo"),
				h.Src(item.OrgLogo),
				h.Alt(item.Org),
			),
		),
		h.Span(g.Attr("data-slot", "period"), g.Text(item.Period)),
		h.Div(g.Attr("data-slot", "title"), g.Text(item.Title)),
		h.Div(g.Attr("data-slot", "org"), g.Text(item.Org)),
		g.If(item.Desc != "", h.P(g.Attr("data-slot", "desc"), g.Text(item.Desc))),
		g.If(len(item.Tags) > 0, h.Div(g.Attr("data-slot", "tags"), g.Group(tagNodes))),
	}

	return h.Div(
		g.Attr("data-component", "timeline-item"),
		g.Group(item.Attrs),
		h.Span(g.Attr("data-slot", "dot")),
		h.Span(g.Attr("data-slot", "line")),
		h.Div(
			g.Attr("data-slot", "card"),
			primitive.Card(primitive.CardProps{Tone: item.Tone}, cardChildren...),
		),
	)
}
