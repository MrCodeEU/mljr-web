//go:build showcase

package primitive

import (
	"mljr-web/ui/registry"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "video-player", Name: "Video Player", Category: "primitive",
		Summary: "Native HTML5 video with custom-styled controls: play/pause, scrub, volume, mute, speed selector, fullscreen. Zero external JS.",
		Code: `primitive.VideoPlayer(primitive.VideoPlayerProps{
    Src:    "/static/demo.mp4",
    Poster: "/static/demo-poster.jpg",
})`,
		Render: func(p map[string]string) g.Node {
			return h.Div(
				h.Style("display:flex;flex-direction:column;gap:var(--sp-4)"),
				VideoPlayer(VideoPlayerProps{
					Src: "https://www.w3schools.com/html/mov_bbb.mp4",
				}),
				h.P(h.Style("font-size:var(--t-xs);color:var(--muted);margin:0"),
					g.Text("Custom controls overlay native <video>. No external library. Uses pointer events to sync playback state with buttons.")),
			)
		},
	})
}
