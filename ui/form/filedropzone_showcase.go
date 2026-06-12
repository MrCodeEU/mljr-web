//go:build showcase

package form

import (
	"mljr-web/ui/registry"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "file-drop-zone", Name: "File Drop Zone", Category: "form",
		Summary: "Drag-and-drop file upload area. Shows file names on selection, drag-over highlight state.",
		Code: `form.FileDropZone(form.FileDropZoneProps{
    Name:      "attachment",
    Accept:    "image/*,.pdf",
    Multiple:  true,
    MaxSizeMB: 10,
})`,
		Render: func(p map[string]string) g.Node {
			return h.Div(
				h.Style("display:flex;flex-direction:column;gap:var(--sp-5)"),
				h.Div(
					h.P(h.Style("font-size:var(--t-xs);text-transform:uppercase;font-weight:700;letter-spacing:.06em;opacity:.5;margin:0 0 var(--sp-3)"), g.Text("Single file")),
					FileDropZone(FileDropZoneProps{
						Name:   "doc",
						Accept: ".pdf,.doc,.docx",
					}),
				),
				h.Div(
					h.P(h.Style("font-size:var(--t-xs);text-transform:uppercase;font-weight:700;letter-spacing:.06em;opacity:.5;margin:0 0 var(--sp-3)"), g.Text("Multiple images")),
					FileDropZone(FileDropZoneProps{
						Name:      "images",
						Accept:    "image/*",
						Multiple:  true,
						MaxSizeMB: 5,
					}),
				),
			)
		},
	})
}
