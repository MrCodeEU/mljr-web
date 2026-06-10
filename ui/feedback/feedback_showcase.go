//go:build showcase

package feedback

import (
	"mljr-web/ui/layout"
	"mljr-web/ui/registry"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "alert", Name: "Alert", Category: "feedback",
		Summary: "Contextual banner: info / success / warning / danger. Optional dismiss button.",
		Code: `feedback.Alert(feedback.AlertProps{
    Variant: feedback.AlertSuccess,
    Title:   "Saved",
    Dismiss: true,
}, g.Text("Your changes have been saved successfully."))`,
		Controls: []registry.Control{
			{Name: "variant", Type: registry.ControlEnum, Options: []string{"info", "success", "warning", "danger"}, Default: "info"},
			{Name: "title", Type: registry.ControlText, Default: "Heads up"},
			{Name: "dismiss", Type: registry.ControlBool, Default: "true"},
		},
		Render: func(p map[string]string) g.Node {
			return Alert(AlertProps{
				Variant: AlertVariant(p["variant"]),
				Title:   p["title"],
				Dismiss: p["dismiss"] == "true",
			}, g.Text("This is an alert message with some additional context."))
		},
	})

	registry.Register(&registry.Component{
		Slug: "spinner", Name: "Spinner", Category: "feedback",
		Summary: "CSS-animated loading indicator. Six variants × three sizes. Swiss + Ink variants match their themes.",
		Code: `feedback.Spinner(feedback.SpinnerProps{Variant: feedback.SpinnerDots, Size: "md"})
feedback.Spinner(feedback.SpinnerProps{Variant: feedback.SpinnerSwiss})  // Swiss Brut
feedback.Spinner(feedback.SpinnerProps{Variant: feedback.SpinnerInk})    // Ink/Paper`,
		Controls: []registry.Control{
			{Name: "variant", Type: registry.ControlEnum, Options: []string{"", "dots", "pulse", "bars", "swiss", "ink"}, Default: ""},
			{Name: "size", Type: registry.ControlEnum, Options: []string{"sm", "md", "lg"}, Default: "md"},
		},
		Render: func(p map[string]string) g.Node {
			return layout.Stack(layout.StackProps{Axis: "h", Gap: "lg"},
				Spinner(SpinnerProps{Variant: SpinnerVariant(p["variant"]), Size: p["size"]}),
			)
		},
		Examples: []registry.Example{
			{Title: "All variants (md)", Node: func() g.Node {
				return h.Div(h.Style("display:flex;gap:var(--sp-8);align-items:center;flex-wrap:wrap;padding:var(--sp-4)"),
					h.Div(h.Style("display:flex;flex-direction:column;align-items:center;gap:var(--sp-3)"),
						Spinner(SpinnerProps{}),
						h.Span(h.Style("font-size:var(--t-xs);opacity:.6"), g.Text("ring")),
					),
					h.Div(h.Style("display:flex;flex-direction:column;align-items:center;gap:var(--sp-3)"),
						Spinner(SpinnerProps{Variant: SpinnerDots}),
						h.Span(h.Style("font-size:var(--t-xs);opacity:.6"), g.Text("dots")),
					),
					h.Div(h.Style("display:flex;flex-direction:column;align-items:center;gap:var(--sp-3)"),
						Spinner(SpinnerProps{Variant: SpinnerPulse}),
						h.Span(h.Style("font-size:var(--t-xs);opacity:.6"), g.Text("pulse")),
					),
					h.Div(h.Style("display:flex;flex-direction:column;align-items:center;gap:var(--sp-3)"),
						Spinner(SpinnerProps{Variant: SpinnerBars}),
						h.Span(h.Style("font-size:var(--t-xs);opacity:.6"), g.Text("bars")),
					),
					h.Div(h.Style("display:flex;flex-direction:column;align-items:center;gap:var(--sp-3)"),
						Spinner(SpinnerProps{Variant: SpinnerSwiss}),
						h.Span(h.Style("font-size:var(--t-xs);opacity:.6"), g.Text("swiss ◼")),
					),
					h.Div(h.Style("display:flex;flex-direction:column;align-items:center;gap:var(--sp-3)"),
						Spinner(SpinnerProps{Variant: SpinnerInk}),
						h.Span(h.Style("font-size:var(--t-xs);opacity:.6"), g.Text("ink ✦")),
					),
				)
			}},
		},
	})

	registry.Register(&registry.Component{
		Slug: "skeleton", Name: "Skeleton", Category: "feedback",
		Summary: "Shimmering content placeholder for loading states.",
		Code: `// text line
feedback.Skeleton(feedback.SkeletonProps{Variant: feedback.SkeletonText, Width: "60%"})

// avatar
feedback.Skeleton(feedback.SkeletonProps{Variant: feedback.SkeletonCircle, Width: "3rem", Height: "3rem"})

// card block
feedback.Skeleton(feedback.SkeletonProps{Variant: feedback.SkeletonRect, Height: "120px"})`,
		Render: func(p map[string]string) g.Node {
			return layout.Stack(layout.StackProps{},
				h.Div(
					g.Attr("style", "display:flex;gap:var(--sp-3);align-items:center"),
					Skeleton(SkeletonProps{Variant: SkeletonCircle, Width: "3rem", Height: "3rem"}),
					layout.Stack(layout.StackProps{},
						Skeleton(SkeletonProps{Variant: SkeletonText, Width: "40%"}),
						Skeleton(SkeletonProps{Variant: SkeletonText, Width: "60%"}),
					),
				),
				Skeleton(SkeletonProps{Variant: SkeletonRect, Height: "120px"}),
				Skeleton(SkeletonProps{Variant: SkeletonText, Width: "80%"}),
				Skeleton(SkeletonProps{Variant: SkeletonText, Width: "55%"}),
			)
		},
	})
}
