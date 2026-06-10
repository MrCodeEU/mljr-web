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
		Slug: "share-button", Name: "Share Button", Category: "primitive",
		Summary: "Web Share API with clipboard fallback. Uses navigator.share() on supported devices; copies link on desktop.",
		Code: `// import "mljr-web/ui/primitive"
primitive.ShareButton(primitive.ShareButtonProps{
    Title: "Check this out",
    Text:  "Found something interesting",
})

// Custom URL
primitive.ShareButton(primitive.ShareButtonProps{
    URL:     "https://mljr.eu",
    Title:   "mljr.eu",
    Variant: token.Primary,
    Label:   "Share Page",
})`,
		Render: func(p map[string]string) g.Node {
			return h.Div(
				h.Style("display:flex;flex-direction:column;gap:var(--sp-5)"),
				h.P(h.Style("font-size:var(--t-sm);color:var(--muted)"), g.Text("On mobile: opens native share sheet. On desktop: copies link to clipboard.")),
				h.Div(
					h.Style("display:flex;gap:var(--sp-3);flex-wrap:wrap;align-items:center"),
					ShareButton(ShareButtonProps{
						Title: "Check this out",
						Text:  "Found something interesting on mljr.eu",
					}),
					ShareButton(ShareButtonProps{
						URL:     "https://mljr.eu",
						Title:   "mljr.eu",
						Variant: token.Primary,
						Label:   "Share Page",
					}),
					ShareButton(ShareButtonProps{
						URL:     "https://mljr.eu",
						Variant: token.Ghost,
						Label:   "Share",
					}),
					ShareButton(ShareButtonProps{
						URL:     "https://mljr.eu",
						Variant: token.Outline,
						Size:    token.SizeSM,
						Label:   "Copy link",
					}),
				),
			)
		},
	})
}
