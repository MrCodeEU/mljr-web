package primitive

import (
	"mljr-web/ui/icon"
	"mljr-web/ui/token"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type FABProps struct {
	Icon     string       // lucide icon name (default "lucide:plus")
	Label    string       // aria-label
	Variant  token.Variant // default Primary
	Size     string       // "sm" | "md" (default) | "lg"
	Position string       // "bottom-right" (default) | "bottom-left" | "top-right" | "top-left"
	Href     string       // if set, renders as <a>
}

// FAB renders a floating action button fixed to the viewport corner.
func FAB(p FABProps, extra ...g.Node) g.Node {
	if p.Icon == "" {
		p.Icon = "lucide:plus"
	}
	if p.Label == "" {
		p.Label = "Action"
	}
	if p.Variant == "" {
		p.Variant = token.Primary
	}
	if p.Position == "" {
		p.Position = "bottom-right"
	}

	attrs := []g.Node{
		g.Attr("data-component", "fab"),
		g.Attr("data-variant", string(p.Variant)),
		g.Attr("data-position", p.Position),
		g.Attr("data-size", p.Size),
		g.Attr("aria-label", p.Label),
		icon.Icon(p.Icon),
		g.Group(extra),
	}

	if p.Href != "" {
		return h.A(append([]g.Node{h.Href(p.Href)}, attrs...)...)
	}
	return h.Button(append([]g.Node{h.Type("button")}, attrs...)...)
}

// SpeedDialItem is one action in the speed dial menu.
type SpeedDialItem struct {
	Icon    string
	Label   string
	OnClick string
	Href    string
}

// SpeedDial renders a FAB with an expandable list of mini action buttons.
func SpeedDial(p FABProps, items []SpeedDialItem) g.Node {
	if p.Icon == "" {
		p.Icon = "lucide:plus"
	}
	if p.Variant == "" {
		p.Variant = token.Primary
	}
	if p.Position == "" {
		p.Position = "bottom-right"
	}
	sig := "_sdOpen"

	miniItems := make([]g.Node, len(items))
	for i, item := range items {
		inner := []g.Node{
			icon.Icon(item.Icon, icon.Props{Size: "1.1rem"}),
		}
		var el g.Node
		if item.Href != "" {
			el = h.A(
				h.Href(item.Href),
				g.Attr("data-component", "fab"),
				g.Attr("data-variant", string(p.Variant)),
				g.Attr("data-size", "sm"),
				g.Attr("aria-label", item.Label),
				g.Group(inner),
			)
		} else {
			el = h.Button(
				h.Type("button"),
				g.Attr("data-component", "fab"),
				g.Attr("data-variant", string(p.Variant)),
				g.Attr("data-size", "sm"),
				g.Attr("aria-label", item.Label),
				g.If(item.OnClick != "", g.Attr("data-on:click", item.OnClick+";$"+sig+"=false")),
				g.Group(inner),
			)
		}
		miniItems[i] = h.Div(
			h.Style("display:flex;align-items:center;gap:var(--sp-2)"),
			h.Span(
				h.Style("background:var(--surface);border:var(--bw-1) solid var(--line);border-radius:var(--radius);padding:var(--sp-1) var(--sp-2);font-size:var(--t-xs);font-weight:700;white-space:nowrap;box-shadow:var(--shadow)"),
				g.Text(item.Label),
			),
			el,
		)
	}

	return h.Div(
		g.Attr("data-component", "speed-dial"),
		g.Attr("data-position", p.Position),
		g.Attr("data-signals", `{"`+sig+`":false}`),
		h.Div(
			g.Attr("data-slot", "items"),
			g.Attr("data-show", "$"+sig),
			h.Style("display:none;flex-direction:column;gap:var(--sp-2);align-items:flex-end"),
			g.Group(miniItems),
		),
		h.Button(
			h.Type("button"),
			g.Attr("data-component", "fab"),
			g.Attr("data-variant", string(p.Variant)),
			g.Attr("aria-label", p.Label),
			g.Attr("aria-expanded", "false"),
			g.Attr("data-attr", `{"aria-expanded":$`+sig+`}`),
			g.Attr("data-on:click", "$"+sig+"=!$"+sig),
			icon.Icon(p.Icon),
		),
	)
}
