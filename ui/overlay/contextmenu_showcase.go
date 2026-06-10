//go:build showcase

package overlay

import (
	"mljr-web/ui/primitive"
	"mljr-web/ui/registry"
	"mljr-web/ui/token"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "context-menu", Name: "Context Menu", Category: "overlay",
		Summary: "Right-click context menu anchored to the cursor. Closes on click-outside, Escape, or item selection.",
		Code: `// 1. Attach to any element via data-ctx
h.Div(g.Attr("data-ctx", "my-menu"), g.Text("Right-click me"))

// 2. Render the menu (place anywhere on page)
overlay.ContextMenu(overlay.ContextMenuProps{
    ID: "my-menu",
    Items: []overlay.ContextMenuItem{
        {Label: "Edit",   Icon: "lucide:edit",   OnClick: "editItem()"},
        {Label: "Share",  Icon: "lucide:share",  OnClick: "shareItem()"},
        {Separator: true},
        {Label: "Delete", Icon: "lucide:trash-2", OnClick: "deleteItem()"},
    },
})`,
		Render: func(p map[string]string) g.Node {
			return h.Div(
				h.Style("display:flex;flex-direction:column;gap:var(--sp-5)"),
				h.P(h.Style("color:var(--muted);font-size:var(--t-sm)"), g.Text("Right-click the card below to open the context menu.")),
				h.Div(
					g.Attr("data-ctx", "demo-ctx"),
					h.Style("padding:var(--sp-6);border:var(--bw-2) dashed var(--line);border-radius:var(--radius);text-align:center;cursor:context-menu;background:var(--surface-2)"),
					h.P(h.Style("font-weight:700;margin:0"), g.Text("Right-click me")),
					h.P(h.Style("color:var(--muted);font-size:var(--t-sm);margin:var(--sp-1) 0 0"), g.Text("file-name.go · 2.4 KB")),
				),
				primitive.Button(primitive.ButtonProps{Variant: token.Outline, Size: token.SizeSM},
					g.Attr("data-ctx", "demo-ctx"),
					g.Text("Or right-click this button"),
				),
				ContextMenu(ContextMenuProps{
					ID: "demo-ctx",
					Items: []ContextMenuItem{
						{Label: "Open", Icon: "lucide:eye", OnClick: "alert('Opening...')"},
						{Label: "Edit", Icon: "lucide:edit", OnClick: "alert('Editing...')"},
						{Label: "Copy path", Icon: "lucide:copy", OnClick: "navigator.clipboard.writeText('demo/file.go')"},
						{Separator: true},
						{Label: "Share", Icon: "lucide:share", OnClick: "alert('Sharing...')"},
						{Separator: true},
						{Label: "Delete", Icon: "lucide:trash-2", OnClick: "alert('Deleting...')"},
					},
				}),
			)
		},
	})
}
