//go:build showcase

package layout

import (
	"mljr-web/ui/registry"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "banner", Name: "Banner", Category: "layout",
		Summary: "Full-width announcement strip. Variants match semantic colors. Optional dismiss button.",
		Code: `layout.Banner(layout.BannerProps{Variant: layout.BannerWarning, Dismiss: true},
    g.Text("Scheduled maintenance on Sunday 02:00–04:00 UTC."),
)`,
		Controls: []registry.Control{
			{Name: "variant", Type: registry.ControlEnum, Options: []string{"", "info", "success", "warning", "danger"}, Default: "info"},
			{Name: "dismiss", Type: registry.ControlBool, Default: "true"},
		},
		Render: func(p map[string]string) g.Node {
			return Banner(BannerProps{
				Variant: BannerVariant(p["variant"]),
				Dismiss: p["dismiss"] == "true",
			}, g.Text("Scheduled maintenance on Sunday 02:00–04:00 UTC. Expect brief downtime."))
		},
	})

	registry.Register(&registry.Component{
		Slug: "stepper", Name: "Stepper", Category: "layout",
		Summary: "Horizontal multi-step progress indicator. States: complete / active / upcoming.",
		Code: `layout.Stepper(layout.StepperProps{},
    layout.Step{Label: "Account",  State: layout.StepComplete},
    layout.Step{Label: "Profile",  State: layout.StepActive},
    layout.Step{Label: "Payment"},
    layout.Step{Label: "Confirm"},
)`,
		Controls: []registry.Control{
			{Name: "active", Type: registry.ControlEnum, Options: []string{"1", "2", "3", "4"}, Default: "2"},
		},
		Render: func(p map[string]string) g.Node {
			labels := []string{"Account", "Profile", "Payment", "Confirm"}
			steps := make([]Step, len(labels))
			active := 1
			switch p["active"] {
			case "1":
				active = 0
			case "3":
				active = 2
			case "4":
				active = 3
			default:
				active = 1
			}
			for i, l := range labels {
				switch {
				case i < active:
					steps[i] = Step{Label: l, State: StepComplete}
				case i == active:
					steps[i] = Step{Label: l, State: StepActive}
				default:
					steps[i] = Step{Label: l}
				}
			}
			return h.Div(h.Style("padding:var(--sp-4)"), Stepper(StepperProps{}, steps...))
		},
	})
}
