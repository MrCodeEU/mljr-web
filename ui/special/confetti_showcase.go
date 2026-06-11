//go:build showcase

package special

import (
	"mljr-web/ui/registry"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "confetti", Name: "Confetti", Category: "special",
		Summary: "Canvas confetti burst. Trigger on click, page load, or manually via window._confetti(). rAF animated, self-cleaning canvas.",
		Code: `special.Confetti(special.ConfettiProps{
    ParticleCount: 120,
    Duration:      3000,
    Trigger:       "click",
    ButtonLabel:   "🎉 Celebrate",
})`,
		Render: func(p map[string]string) g.Node {
			return h.Div(
				h.Style("display:flex;gap:var(--sp-4);flex-wrap:wrap;align-items:center"),
				Confetti(ConfettiProps{
					ParticleCount: 120,
					Duration:      3000,
					Trigger:       "click",
					ButtonLabel:   "🎉 Celebrate",
				}),
				Confetti(ConfettiProps{
					ParticleCount: 60,
					Duration:      2000,
					Trigger:       "click",
					ButtonLabel:   "Mini burst",
					Colors:        []string{"#f59e0b", "#10b981", "#6366f1"},
				}),
			)
		},
	})
}
