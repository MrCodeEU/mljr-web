//go:build showcase

package feedback

import (
	"mljr-web/ui/registry"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "shimmer", Name: "Shimmer", Category: "feedback",
		Summary: "Animated gradient loading placeholder. Lighter than Skeleton — just the shimmer effect with configurable size and shape.",
		Code: `feedback.Shimmer(feedback.ShimmerProps{
    Width:  "100%",
    Height: "1.2em",
})

// Multiple lines
feedback.Shimmer(feedback.ShimmerProps{
    Width:  "280px",
    Height: "0.9em",
    Lines:  4,
})

// Circle avatar
feedback.Shimmer(feedback.ShimmerProps{
    Width:  "48px",
    Height: "48px",
    Circle: true,
})`,
		Render: func(p map[string]string) g.Node {
			return h.Div(
				h.Style("display:flex;flex-direction:column;gap:var(--sp-6)"),
				// Card skeleton
				h.Div(
					h.Style("border:var(--bw-1) solid var(--line);border-radius:var(--radius);padding:var(--sp-4);display:flex;flex-direction:column;gap:var(--sp-3)"),
					Shimmer(ShimmerProps{Width: "100%", Height: "180px"}),
					Shimmer(ShimmerProps{Width: "60%", Height: "1.2em"}),
					Shimmer(ShimmerProps{Width: "100%", Height: "0.9em", Lines: 3}),
					h.Div(h.Style("display:flex;gap:var(--sp-2)"),
						Shimmer(ShimmerProps{Width: "80px", Height: "32px"}),
						Shimmer(ShimmerProps{Width: "80px", Height: "32px"}),
					),
				),
				// Profile list skeleton
				h.Div(
					h.Style("display:flex;flex-direction:column;gap:var(--sp-3)"),
					h.P(h.Style("font-size:var(--t-sm);font-weight:700;color:var(--muted)"), g.Text("List skeleton:")),
					profileSkeleton(), profileSkeleton(), profileSkeleton(),
				),
			)
		},
	})
}

func profileSkeleton() g.Node {
	return h.Div(
		h.Style("display:flex;align-items:center;gap:var(--sp-3)"),
		Shimmer(ShimmerProps{Width: "40px", Height: "40px", Circle: true}),
		h.Div(h.Style("flex:1;display:flex;flex-direction:column;gap:var(--sp-2)"),
			Shimmer(ShimmerProps{Width: "40%", Height: "0.85em"}),
			Shimmer(ShimmerProps{Width: "60%", Height: "0.75em"}),
		),
	)
}
