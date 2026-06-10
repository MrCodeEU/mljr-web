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
		Slug: "copy-button", Name: "Copy Button", Category: "primitive",
		Summary: "Copies text to clipboard. Datastar signal flips to show checkmark feedback, resets after 2 s.",
		Code: `// import "mljr-web/ui/primitive"
primitive.CopyButton(primitive.CopyButtonProps{
    Text: "npm install mljr-ui",
})

// With label
primitive.CopyButton(primitive.CopyButtonProps{
    Text:    someText,
    Variant: token.Outline,
    Size:    token.SizeSM,
    Label:   "Copy code",
})`,
		Render: func(p map[string]string) g.Node {
			return h.Div(
				h.Style("display:flex;flex-direction:column;gap:var(--sp-5)"),
				h.Div(
					h.Style("display:flex;gap:var(--sp-3);align-items:center;flex-wrap:wrap"),
					CopyButton(CopyButtonProps{Text: "npm install mljr-ui"}),
					CopyButton(CopyButtonProps{Text: "go get mljr-web/ui", Variant: token.Primary}),
					CopyButton(CopyButtonProps{Text: "curl https://mljr.eu", Variant: token.Outline, Size: token.SizeSM, Label: "Copy URL"}),
				),
				h.Div(
					h.Style("display:flex;align-items:center;gap:var(--sp-3);padding:var(--sp-3);background:var(--surface-2);border-radius:var(--radius);border:var(--bw-1) solid var(--line)"),
					h.Code(h.Style("flex:1;font-family:var(--font-mono);font-size:var(--t-sm)"), g.Text("go get mljr-web/ui")),
					CopyButton(CopyButtonProps{Text: "go get mljr-web/ui"}),
				),
			)
		},
	})
}
