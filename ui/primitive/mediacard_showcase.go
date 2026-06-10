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
		Slug: "media-card", Name: "Media Card", Category: "primitive",
		Summary: "Card with an image/media area on top and content below. Badge overlay, lazy loading, optional title link, action slot.",
		Code: `primitive.MediaCard(primitive.MediaCardProps{
    ImageSrc:    "/static/hero.jpg",
    AspectRatio: "16/9",
    Badge:       "New",
    BadgeTone:   token.ToneAccent,
    Title:       "Article title",
    Description: "Short description…",
    Lazy:        true,
},
    primitive.Button(primitive.ButtonProps{Variant: token.Outline, Size: token.SizeSM}, g.Text("Read more")),
)`,
		Render: func(p map[string]string) g.Node {
			placeholder := func(ratio, color string) g.Node {
				return h.Div(h.Style("width:100%;height:100%;background:"+color+";display:flex;align-items:center;justify-content:center;color:white;font-weight:700;font-size:var(--t-sm)"), g.Text("Image"))
			}
			return h.Div(
				h.Style("display:grid;grid-template-columns:repeat(auto-fill,minmax(260px,1fr));gap:var(--sp-5)"),
				MediaCard(MediaCardProps{
					AspectRatio: "16/9",
					Badge:       "Featured",
					BadgeTone:   token.ToneAccent,
					Title:       "Getting started with Go",
					Description: "A practical guide to building web services in Go.",
				},
					placeholder("16/9", "var(--accent)"),
					Button(ButtonProps{Variant: token.Outline, Size: token.SizeSM}, g.Text("Read more")),
				),
				MediaCard(MediaCardProps{
					AspectRatio: "4/3",
					Badge:       "New",
					Title:       "Tailwind v4 deep dive",
					Description: "Everything that changed in Tailwind CSS v4.",
				},
					placeholder("4/3", "#6366f1"),
				),
				MediaCard(MediaCardProps{
					AspectRatio: "1/1",
					Title:       "Motion v10",
					Description: "Animation library for the web.",
					Href:        "#",
				},
					placeholder("1/1", "#10b981"),
				),
			)
		},
	})
}
