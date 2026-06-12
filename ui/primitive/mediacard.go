package primitive

import (
	"mljr-web/ui/token"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type MediaCardProps struct {
	ImageSrc string
	ImageAlt string
	// AspectRatio: "16/9" (default) | "4/3" | "1/1" | "3/4"
	AspectRatio string
	// Badge text overlaid on the image (empty = none)
	Badge       string
	BadgeTone   token.Tone
	Title       string
	Description string
	// Href makes the title a link
	Href string
	Lazy bool
}

// MediaCard renders a card with an image on top and content below.
// Pass action nodes (buttons, links) as children.
func MediaCard(p MediaCardProps, actions ...g.Node) g.Node {
	if p.AspectRatio == "" {
		p.AspectRatio = "16/9"
	}
	if p.ImageAlt == "" {
		p.ImageAlt = p.Title
	}

	loading := "eager"
	if p.Lazy {
		loading = "lazy"
	}

	var badgeNode g.Node
	if p.Badge != "" {
		badgeNode = h.Span(
			g.Attr("data-component", "badge"),
			g.If(p.BadgeTone != "", g.Attr("data-tone", string(p.BadgeTone))),
			h.Style("position:absolute;top:var(--sp-3);left:var(--sp-3);z-index:1"),
			g.Text(p.Badge),
		)
	}

	var titleNode g.Node
	if p.Href != "" {
		titleNode = h.A(
			h.Href(p.Href),
			g.Attr("data-slot", "title"),
			g.Text(p.Title),
		)
	} else {
		titleNode = h.Div(g.Attr("data-slot", "title"), g.Text(p.Title))
	}

	return h.Div(
		g.Attr("data-component", "media-card"),
		// Image area
		h.Div(
			g.Attr("data-slot", "media"),
			h.Style("aspect-ratio:"+p.AspectRatio+";position:relative;overflow:hidden"),
			badgeNode,
			g.If(p.ImageSrc != "", h.Img(
				h.Src(p.ImageSrc),
				h.Alt(p.ImageAlt),
				g.Attr("loading", loading),
				h.Style("width:100%;height:100%;object-fit:cover"),
			)),
		),
		// Content area
		h.Div(
			g.Attr("data-slot", "body"),
			titleNode,
			g.If(p.Description != "", h.P(g.Attr("data-slot", "description"), g.Text(p.Description))),
			g.If(len(actions) > 0, h.Div(g.Attr("data-slot", "actions"), g.Group(actions))),
		),
	)
}
