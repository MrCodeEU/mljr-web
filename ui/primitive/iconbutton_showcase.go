//go:build showcase

package primitive

import (
	"mljr-web/ui/registry"
	"mljr-web/ui/token"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "icon-button", Name: "Icon Button", Category: "primitive",
		Summary: "Square icon-only button. Thin wrapper over Button with size=icon — adds aria-label and pass-through Attrs.",
		Code: `primitive.IconButton(primitive.IconButtonProps{
    Icon:    "lucide:settings",
    Label:   "Open settings",
    Variant: token.Outline,
})`,
		Render: func(p map[string]string) g.Node {
			return h.Div(
				h.Style("display:flex;flex-direction:column;gap:var(--sp-5)"),
				h.Div(
					h.P(h.Style("font-size:var(--t-xs);opacity:.5;font-weight:700;text-transform:uppercase;letter-spacing:.06em;margin:0 0 var(--sp-2)"), g.Text("Variants")),
					h.Div(h.Style("display:flex;gap:var(--sp-2);flex-wrap:wrap"),
						IconButton(IconButtonProps{Icon: "lucide:settings", Label: "Settings", Variant: token.Outline}),
						IconButton(IconButtonProps{Icon: "lucide:search", Label: "Search", Variant: token.Primary}),
						IconButton(IconButtonProps{Icon: "lucide:trash-2", Label: "Delete", Variant: token.Danger}),
						IconButton(IconButtonProps{Icon: "lucide:heart", Label: "Like", Variant: token.Ghost}),
						IconButton(IconButtonProps{Icon: "lucide:share", Label: "Share", Variant: token.Secondary}),
					),
				),
				h.Div(
					h.P(h.Style("font-size:var(--t-xs);opacity:.5;font-weight:700;text-transform:uppercase;letter-spacing:.06em;margin:0 0 var(--sp-2)"), g.Text("Inline with text buttons")),
					h.Div(h.Style("display:flex;gap:var(--sp-2);align-items:center;flex-wrap:wrap"),
						Button(ButtonProps{Variant: token.Primary}, g.Text("Save")),
						IconButton(IconButtonProps{Icon: "lucide:copy", Label: "Copy", Variant: token.Outline}),
						IconButton(IconButtonProps{Icon: "lucide:external-link", Label: "Open", Variant: token.Outline}),
						IconButton(IconButtonProps{Icon: "lucide:trash-2", Label: "Delete", Variant: token.Danger}),
					),
				),
			)
		},
	})
}
