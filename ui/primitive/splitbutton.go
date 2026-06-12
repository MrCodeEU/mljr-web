package primitive

import (
	"mljr-web/ui/icon"
	"mljr-web/ui/token"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type SplitButtonItem struct {
	Label   string
	OnClick string // Datastar or JS expression
	Href    string
}

type SplitButtonProps struct {
	Label      string
	Variant    token.Variant
	Size       token.Size
	OnClick    string            // main button action
	Href       string            // main button href (if set, renders as link)
	Items      []SplitButtonItem // dropdown items
	SignalName string            // default "_sbOpen"
}

// SplitButton renders a primary action button + chevron dropdown for alternate actions.
func SplitButton(p SplitButtonProps) g.Node {
	if p.Variant == "" {
		p.Variant = token.Primary
	}
	if p.SignalName == "" {
		p.SignalName = "_sbOpen"
	}
	sig := p.SignalName

	// Main button
	var mainBtn g.Node
	mainAttrs := []g.Node{
		g.Attr("data-component", "button"),
		g.Attr("data-variant", string(p.Variant)),
		g.Attr("data-size", string(p.Size)),
		g.Attr("data-slot", "main"),
	}
	if p.OnClick != "" {
		mainAttrs = append(mainAttrs, g.Attr("data-on:click", p.OnClick))
	}
	if p.Href != "" {
		mainBtn = h.A(append([]g.Node{h.Href(p.Href)}, append(mainAttrs, g.Text(p.Label))...)...)
	} else {
		mainBtn = h.Button(append(mainAttrs, h.Type("button"), g.Text(p.Label))...)
	}

	// Dropdown items
	itemNodes := make([]g.Node, len(p.Items))
	for i, item := range p.Items {
		var el g.Node
		if item.Href != "" {
			el = h.A(
				h.Href(item.Href),
				g.Attr("data-slot", "item"),
				g.Attr("data-on:click", "$"+sig+"=false"),
				g.Text(item.Label),
			)
		} else {
			clickExpr := "$" + sig + "=false"
			if item.OnClick != "" {
				clickExpr = item.OnClick + ";" + clickExpr
			}
			el = h.Button(
				h.Type("button"),
				g.Attr("data-slot", "item"),
				g.Attr("data-on:click", clickExpr),
				g.Text(item.Label),
			)
		}
		itemNodes[i] = h.Li(el)
	}

	return h.Div(
		g.Attr("data-component", "split-button"),
		g.Attr("data-signals", `{"`+sig+`":false}`),
		mainBtn,
		h.Button(
			h.Type("button"),
			g.Attr("data-component", "button"),
			g.Attr("data-variant", string(p.Variant)),
			g.Attr("data-size", string(p.Size)),
			g.Attr("data-slot", "chevron"),
			g.Attr("data-on:click", "$"+sig+"=!$"+sig),
			g.Attr("aria-label", "More actions"),
			icon.Icon("lucide:chevron-down", icon.Props{Size: "1rem"}),
		),
		h.Ul(
			g.Attr("data-slot", "dropdown"),
			g.Attr("data-show", "$"+sig),
			h.Style("display:none"),
			h.Role("menu"),
			g.Group(itemNodes),
		),
	)
}
