package overlay

import (
	"fmt"

	"mljr-web/ui/icon"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type DropdownItem struct {
	Label   string
	Href    string // if set, renders <a>; otherwise <button>
	OnClick string // Datastar expression
	Icon    string // icon name (e.g. "lucide:edit")
	Variant string // "danger" for destructive items
	Divider bool   // render a divider before this item
}

type DropdownProps struct {
	Signal string // signal name for open state (default "dd")
	Align  string // "left" (default) | "right"
	Attrs  []g.Node
}

// Dropdown wraps a trigger + Datastar-driven menu. Signal toggles open/close.
// Close on outside-click handled via data-on:click__window.
func Dropdown(p DropdownProps, trigger g.Node, items ...DropdownItem) g.Node {
	sig := p.Signal
	if sig == "" {
		sig = "dd"
	}
	openExpr := "$" + sig
	toggleExpr := fmt.Sprintf("$%s=!$%s", sig, sig)
	closeExpr := fmt.Sprintf("$%s=false", sig)

	menuItems := make([]g.Node, 0, len(items))
	for _, item := range items {
		if item.Divider {
			menuItems = append(menuItems, h.Div(g.Attr("data-component", "dropdown-divider")))
		}
		var iconNode g.Node
		if item.Icon != "" {
			iconNode = icon.Icon(item.Icon)
		}
		attrs := []g.Node{
			g.Attr("data-component", "dropdown-item"),
			g.If(item.Variant != "", g.Attr("data-variant", item.Variant)),
		}
		if item.Href != "" {
			attrs = append(attrs, h.Href(item.Href))
			if item.OnClick != "" {
				attrs = append(attrs, g.Attr("data-on:click", item.OnClick+";"+closeExpr))
			}
			attrs = append(attrs, iconNode, g.Text(item.Label))
			menuItems = append(menuItems, h.A(attrs...))
		} else {
			expr := item.OnClick
			if expr != "" {
				expr += ";" + closeExpr
			} else {
				expr = closeExpr
			}
			attrs = append(attrs, g.Attr("data-on:click", expr), h.Type("button"))
			attrs = append(attrs, iconNode, g.Text(item.Label))
			menuItems = append(menuItems, h.Button(attrs...))
		}
	}

	align := p.Align
	if align == "" {
		align = "left"
	}

	return h.Div(
		g.Attr("data-component", "dropdown"),
		g.Attr("data-signals", fmt.Sprintf(`{%s:false}`, sig)),
		g.Attr("data-on:click__window", closeExpr),
		g.Group(p.Attrs),
		h.Div(
			g.Attr("data-on:click__stop", toggleExpr),
			trigger,
		),
		h.Div(
			g.Attr("data-component", "dropdown-menu"),
			g.If(align == "right", g.Attr("data-align", "right")),
			g.Attr("data-show", openExpr),
			h.Style("display:none"),
			g.Attr("data-on:click__stop", "void 0"),
			g.Group(menuItems),
		),
	)
}
