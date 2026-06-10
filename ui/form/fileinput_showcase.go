//go:build showcase

package form

import (
	"mljr-web/ui"
	"mljr-web/ui/registry"

	g "maragu.dev/gomponents"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "file-input", Name: "File Input", Category: "form",
		Summary: "Styled file upload trigger with filename feedback. Wraps a hidden native input.",
		Code: `form.FileInput(form.FileInputProps{
    Signal: "uploadFile",
    Name:   "file",
    Accept: "image/*,.pdf",
    Label:  "Choose file or drag & drop",
})`,
		Render: func(p map[string]string) g.Node {
			return g.Group{
				ui.Signals(`{"uploadFile":""}`),
				Field(FieldProps{Label: "Upload file", Hint: "PNG, JPG or PDF up to 10 MB"},
					FileInput(FileInputProps{
						Signal: "uploadFile",
						Name:   "file",
						Accept: "image/*,.pdf",
					}),
				),
			}
		},
	})
}
