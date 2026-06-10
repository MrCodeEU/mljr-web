//go:build showcase

package data

import (
	"mljr-web/ui/registry"
	"mljr-web/ui/token"

	g "maragu.dev/gomponents"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "timeline", Name: "Timeline", Category: "data",
		Summary: "Vertical timeline of dated entries with card cards and tags.",
		Code: `data.Timeline(data.TimelineProps{},
    data.TimelineItem(data.TimelineItemProps{
        Period: "Jan. 2024–Today",
        Title:  "Senior Engineer",
        Org:    "Acme Corp",
        Tags:   []string{"Go", "Kubernetes"},
    }, h.P(g.Text("Led platform migration."))),
)`,
		Controls: []registry.Control{
			{Name: "tone", Type: registry.ControlEnum, Options: []string{"", "yellow", "cyan", "violet", "pink", "lime", "mint", "sky", "blush", "accent"}, Default: ""},
		},
		Render: func(p map[string]string) g.Node {
			return Timeline(TimelineProps{},
				TimelineItem{
					Period: "Jan. 2024–Heute",
					Title:  "Senior Engineer",
					Org:    "Acme Corp",
					Desc:   "Led platform migration to cloud-native architecture.",
					Tags:   []string{"Go", "Kubernetes", "Datastar"},
					Tone:   token.Tone(p["tone"]),
				},
				TimelineItem{
					Period: "Mar. 2021–Dec. 2023",
					Title:  "Software Engineer",
					Org:    "Beta GmbH",
					Desc:   "Built internal tooling and APIs.",
					Tags:   []string{"Go", "Postgres"},
					Tone:   token.Tone(p["tone"]),
				},
			)
		},
	})
}
