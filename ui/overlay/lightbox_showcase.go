//go:build showcase

package overlay

import (
	"mljr-web/ui/registry"

	g "maragu.dev/gomponents"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "lightbox", Name: "Lightbox", Category: "overlay",
		PreviewHeight: "520px",
		Summary: "Full-screen image viewer. Click thumbnail to open, ←/→ keys navigate, click backdrop or × to close.",
		Code: `overlay.Lightbox(overlay.LightboxProps{
    ID:      "gallery",
    Columns: 3,
    Images: []overlay.LightboxImage{
        {Src: "/img/photo1.jpg", Thumb: "/img/photo1-sm.jpg", Alt: "Photo 1", Caption: "Caption 1"},
    },
})`,
		Render: func(p map[string]string) g.Node {
			imgs := []LightboxImage{
				{Src: "https://picsum.photos/seed/lb1/800/600", Thumb: "https://picsum.photos/seed/lb1/240/180", Alt: "Mountain lake", Caption: "Mountain lake at dawn"},
				{Src: "https://picsum.photos/seed/lb2/800/600", Thumb: "https://picsum.photos/seed/lb2/240/180", Alt: "Forest path", Caption: "Forest path in autumn"},
				{Src: "https://picsum.photos/seed/lb3/800/600", Thumb: "https://picsum.photos/seed/lb3/240/180", Alt: "City skyline", Caption: "City skyline at night"},
				{Src: "https://picsum.photos/seed/lb4/800/600", Thumb: "https://picsum.photos/seed/lb4/240/180", Alt: "Abstract art", Caption: "Abstract composition"},
				{Src: "https://picsum.photos/seed/lb5/800/600", Thumb: "https://picsum.photos/seed/lb5/240/180", Alt: "Coastline", Caption: "Coastline at sunset"},
				{Src: "https://picsum.photos/seed/lb6/800/600", Thumb: "https://picsum.photos/seed/lb6/240/180", Alt: "Architecture", Caption: "Modern architecture"},
			}
			return Lightbox(LightboxProps{
				ID:       "demo-lb",
				Images:   imgs,
				Columns:  3,
				ThumbSize: "140px",
			})
		},
	})
}
