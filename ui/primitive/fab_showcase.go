//go:build showcase

package primitive

import (
	"mljr-web/ui/registry"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "fab", Name: "FAB", Category: "primitive",
		Summary: "Floating action button fixed to a viewport corner. Speed Dial variant expands a mini-action list.",
		Code: `// Simple FAB
primitive.FAB(primitive.FABProps{Icon: "lucide:plus", Label: "New item"})

// Speed Dial
primitive.SpeedDial(
    primitive.FABProps{Icon: "lucide:plus", Label: "Actions"},
    []primitive.SpeedDialItem{
        {Icon: "lucide:edit",     Label: "Edit",    OnClick: "edit()"},
        {Icon: "lucide:share",    Label: "Share",   OnClick: "share()"},
        {Icon: "lucide:trash-2",  Label: "Delete",  OnClick: "del()"},
    },
)`,
		Render: func(p map[string]string) g.Node {
			return h.Div(
				h.Style("position:relative;min-height:300px;border:var(--bw-1) dashed var(--line);border-radius:var(--radius);padding:var(--sp-4)"),
				h.P(h.Style("color:var(--muted);font-size:var(--t-sm)"), g.Text("FAB components are positioned inside this container (normally they float on the viewport).")),
				h.Div(h.Style("position:absolute;bottom:var(--sp-4);right:var(--sp-4);display:flex;flex-direction:column;align-items:flex-end;gap:var(--sp-3)"),
					SpeedDial(
						FABProps{Icon: "lucide:plus", Label: "New"},
						[]SpeedDialItem{
							{Icon: "lucide:file-text", Label: "New document", OnClick: "alert('doc')"},
							{Icon: "lucide:folder", Label: "New folder", OnClick: "alert('folder')"},
							{Icon: "lucide:upload", Label: "Upload", OnClick: "alert('upload')"},
						},
					),
				),
				h.Div(h.Style("position:absolute;bottom:var(--sp-4);left:var(--sp-4)"),
					FAB(FABProps{Icon: "lucide:heart", Label: "Like", Position: "bottom-left"}),
				),
			)
		},
	})
}
