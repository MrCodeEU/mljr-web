package special

import (
	"fmt"

	"mljr-web/ui/icon"
	"mljr-web/ui/primitive"
	"mljr-web/ui/token"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type UserMenuItem struct {
	Label   string
	Href    string
	OnClick string
	Icon    string
	Divider bool
	Danger  bool
}

type UserMenuProps struct {
	Name      string
	Email     string
	AvatarSrc string
	Initials  string
	Size      token.Size
	// Signal is the Datastar signal name (default "_um").
	Signal string
	// Align: "left" | "right" (default "right").
	Align string
	Items []UserMenuItem
}

// UserMenu renders a clickable avatar that opens a dropdown with user identity + actions.
// Composite of primitive.Avatar + a custom-built Datastar dropdown.
func UserMenu(p UserMenuProps) g.Node {
	sig := p.Signal
	if sig == "" {
		sig = "_um"
	}
	align := p.Align
	if align == "" {
		align = "right"
	}

	toggleExpr := fmt.Sprintf("$%s=!$%s", sig, sig)
	closeExpr := fmt.Sprintf("$%s=false", sig)
	openExpr := "$" + sig

	// Build menu items
	var itemNodes []g.Node
	for _, item := range p.Items {
		if item.Divider {
			itemNodes = append(itemNodes, h.Div(h.Style("height:var(--bw-1);background:var(--line);margin:var(--sp-1) 0")))
		}
		var iconNode g.Node
		if item.Icon != "" {
			iconNode = icon.Icon(item.Icon)
		}
		expr := closeExpr
		if item.OnClick != "" {
			expr = item.OnClick + ";" + closeExpr
		}
		variant := ""
		if item.Danger {
			variant = "danger"
		}

		var el g.Node
		if item.Href != "" {
			el = h.A(
				h.Href(item.Href),
				g.Attr("data-slot", "item"),
				g.If(variant != "", g.Attr("data-variant", variant)),
				g.Attr("data-on:click", closeExpr),
				iconNode,
				g.Text(item.Label),
			)
		} else {
			el = h.Button(
				h.Type("button"),
				g.Attr("data-slot", "item"),
				g.If(variant != "", g.Attr("data-variant", variant)),
				g.Attr("data-on:click", expr),
				iconNode,
				g.Text(item.Label),
			)
		}
		itemNodes = append(itemNodes, el)
	}

	return h.Div(
		g.Attr("data-component", "user-menu"),
		g.Attr("data-signals", fmt.Sprintf(`{"%s":false}`, sig)),
		g.Attr("data-on:click__window", closeExpr),
		g.Attr("style", "position:relative;display:inline-flex"),

		// Trigger
		h.Div(
			g.Attr("data-on:click__stop", toggleExpr),
			h.Style("cursor:pointer"),
			primitive.Avatar(primitive.AvatarProps{
				Src:      p.AvatarSrc,
				Initials: p.Initials,
				Size:     p.Size,
			}),
		),

		// Menu
		h.Div(
			g.Attr("data-component", "user-menu-panel"),
			g.If(align == "right", g.Attr("data-align", "right")),
			g.Attr("data-show", openExpr),
			h.Style("display:none"),
			g.Attr("data-on:click__stop", "void 0"),

			// Identity header
			h.Div(
				g.Attr("data-slot", "header"),
				primitive.Avatar(primitive.AvatarProps{
					Src:      p.AvatarSrc,
					Initials: p.Initials,
					Size:     token.SizeSM,
				}),
				h.Div(
					h.Style("min-width:0;flex:1"),
					h.Div(h.Style("font-weight:700;font-size:var(--t-sm);overflow:hidden;text-overflow:ellipsis;white-space:nowrap"), g.Text(p.Name)),
					g.If(p.Email != "", h.Div(h.Style("font-size:var(--t-xs);color:var(--muted);overflow:hidden;text-overflow:ellipsis;white-space:nowrap"), g.Text(p.Email))),
				),
			),

			h.Div(h.Style("height:var(--bw-1);background:var(--line);margin:var(--sp-1) 0")),
			g.Group(itemNodes),
		),
	)
}
