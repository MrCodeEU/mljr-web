//go:build showcase

package primitive

import (
	"mljr-web/ui/registry"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "audio-player", Name: "Audio Player", Category: "primitive",
		Summary: "Native HTML5 audio with animated Canvas waveform (40 bars), play/pause, scrub, time display. Zero external JS.",
		Code: `primitive.AudioPlayer(primitive.AudioPlayerProps{
    Src:   "/static/demo.mp3",
    Title: "Track Name",
})`,
		Render: func(p map[string]string) g.Node {
			return h.Div(
				h.Style("display:flex;flex-direction:column;gap:var(--sp-4);max-width:480px"),
				AudioPlayer(AudioPlayerProps{
					Src:   "https://www.w3schools.com/html/horse.mp3",
					Title: "Demo Audio",
				}),
				h.P(h.Style("font-size:var(--t-xs);color:var(--muted);margin:0"),
					g.Text("Canvas waveform animates on playback. Bars reflect real audio frequency data via Web Audio API analyser.")),
			)
		},
	})
}
