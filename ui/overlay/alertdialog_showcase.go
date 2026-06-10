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
		Slug: "alert-dialog", Name: "Alert Dialog", Category: "overlay",
		Summary: "Confirmation dialog with cancel/confirm. Driven by Datastar signal; confirm runs an expression before closing.",
		Code: `// 1. Place AlertDialog anywhere on page
overlay.AlertDialog(overlay.AlertDialogProps{
    Title:       "Delete project?",
    Description: "This action cannot be undone.",
    ConfirmText: "Delete",
    Variant:     token.Danger,
    OnConfirm:   "deleteProject($projectId)",
})

// 2. Open via Datastar
data-on:click="$_alertOpen=true"`,
		Render: func(p map[string]string) g.Node {
			return h.Div(
				h.Style("display:flex;flex-direction:column;gap:var(--sp-4);align-items:flex-start;padding:var(--sp-4)"),
				h.P(h.Style("color:var(--muted);font-size:var(--t-sm)"), g.Text("Click a button to open the confirmation dialog.")),
				h.Div(
					h.Style("display:flex;gap:var(--sp-3);flex-wrap:wrap"),
					primitive.Button(primitive.ButtonProps{Variant: token.Danger},
						g.Attr("data-on:click", "$_alertOpen=true"),
						g.Text("Delete project"),
					),
					primitive.Button(primitive.ButtonProps{Variant: token.Outline},
						g.Attr("data-on:click", "$_alertOpen2=true"),
						g.Text("Archive item"),
					),
				),
				AlertDialog(AlertDialogProps{
					Title:       "Delete this project?",
					Description: "All data will be permanently deleted. This action cannot be undone.",
					ConfirmText: "Delete permanently",
					CancelText:  "Keep it",
					Variant:     token.Danger,
					OnConfirm:   "alert('Project deleted!')",
				}),
				AlertDialog(AlertDialogProps{
					SignalName:  "_alertOpen2",
					Title:       "Archive this item?",
					Description: "Archived items can be restored later from settings.",
					ConfirmText: "Archive",
					Variant:     token.Primary,
					OnConfirm:   "alert('Item archived!')",
				}),
			)
		},
	})
}
