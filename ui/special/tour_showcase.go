//go:build showcase

package special

import (
	"mljr-web/ui/primitive"
	"mljr-web/ui/registry"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "onboarding-tour", Name: "Onboarding Tour", Category: "special",
		Summary: "Step-by-step spotlight tour with floating tooltip. Clip-path overlay highlights target element. No external library.",
		Code: `special.Tour(special.TourProps{Signal: "_demo"},
    special.TourStep{Target: "#step1", Title: "Welcome", Body: "This is step 1.", Placement: "bottom"},
    special.TourStep{Target: "#step2", Title: "Next", Body: "And this is step 2.", Placement: "right"},
)`,
		Render: func(p map[string]string) g.Node {
			return h.Div(
				h.Style("position:relative;padding:var(--sp-4);border:var(--bw-2) solid var(--line);border-radius:var(--radius)"),
				// Start button
				h.Button(
					g.Attr("data-component", "button"),
					g.Attr("data-variant", "primary"),
					g.Attr("data-tour-start", ""),
					g.Text("Start Tour"),
				),
				// Tour targets
				h.Div(
					h.Style("display:grid;grid-template-columns:1fr 1fr;gap:var(--sp-4);margin-top:var(--sp-4)"),
					primitive.Card(primitive.CardProps{}, h.Div(h.ID("tour-step1"),
						h.Strong(g.Text("Feature A")),
						h.P(h.Style("margin:var(--sp-1) 0 0;font-size:var(--t-sm);color:var(--muted)"), g.Text("Click 'Start Tour' above.")),
					)),
					primitive.Card(primitive.CardProps{}, h.Div(h.ID("tour-step2"),
						h.Strong(g.Text("Feature B")),
						h.P(h.Style("margin:var(--sp-1) 0 0;font-size:var(--t-sm);color:var(--muted)"), g.Text("This gets highlighted step 2.")),
					)),
				),
				Tour(TourProps{Signal: "_showcaseTour"},
					TourStep{Target: "[data-tour-start]", Title: "Start Here", Body: "Click this button to launch the onboarding tour.", Placement: "bottom"},
					TourStep{Target: "#tour-step1", Title: "Feature A", Body: "This panel highlights your first feature. Spotlight clips around the target.", Placement: "bottom"},
					TourStep{Target: "#tour-step2", Title: "Feature B", Body: "Navigate with Back/Next or click outside to dismiss.", Placement: "bottom"},
				),
			)
		},
	})
}
