//go:build showcase

package data

import (
	"mljr-web/ui/registry"
	"mljr-web/ui/token"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "horizontal-timeline", Name: "Horizontal Timeline", Category: "data",
		Summary: "CSS scroll-snap horizontal timeline. Dot + card per item. No JS — pure CSS scroll.",
		Code: `data.HorizontalTimeline(data.HorizontalTimelineProps{ScrollSnap: true},
    data.HorizontalTimelineItem{Label: "Founded", Period: "2020", Desc: "Company started."},
    data.HorizontalTimelineItem{Label: "Launch", Period: "2021", Tone: token.ToneAccent},
)`,
		Render: func(p map[string]string) g.Node {
			return h.Div(
				HorizontalTimeline(HorizontalTimelineProps{ScrollSnap: true},
					HorizontalTimelineItem{Label: "Founded", Period: "Jan 2019", Desc: "Company incorporated in Berlin.", Tone: token.ToneNone},
					HorizontalTimelineItem{Label: "Seed Round", Period: "Aug 2019", Desc: "$1.2M raised from angels.", Tone: token.ToneAccent},
					HorizontalTimelineItem{Label: "v1 Launch", Period: "Mar 2020", Desc: "Product shipped to 50 beta users.", Tone: token.ToneNone},
					HorizontalTimelineItem{Label: "Series A", Period: "Nov 2020", Desc: "$8M to scale the team.", Tone: token.ToneAccent},
					HorizontalTimelineItem{Label: "10k Users", Period: "Jun 2021", Desc: "Crossed 10,000 active users.", Tone: token.ToneNone},
					HorizontalTimelineItem{Label: "Today", Period: "2024", Desc: "50k users, profitable.", Tone: token.ToneAccent, Current: true},
				),
			)
		},
	})
}
