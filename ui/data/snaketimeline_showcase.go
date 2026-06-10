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
		Slug: "snake-timeline", Name: "Snake Timeline", Category: "data",
		Summary:       "Serpentine timeline where items flow left→right in odd rows and right→left in even rows. Curved connectors bridge the row ends. Great for long event histories without excessive vertical scroll.",
		PreviewHeight: "720px",
		Code: `data.SnakeTimeline(data.SnakeTimelineProps{Cols: 2},
    data.SnakeTimelineItem{Period: "2025–Now", Title: "Software Engineer", Org: "Dynatrace", Tone: token.ToneCyan},
    data.SnakeTimelineItem{Period: "2024–2025", Title: "MSc Student", Org: "JKU", Tone: token.ToneViolet},
    data.SnakeTimelineItem{Period: "2023–2024", Title: "Tutor", Org: "JKU", Tone: token.ToneLime},
    data.SnakeTimelineItem{Period: "2022–2023", Title: "Intern", Org: "Company", Tone: token.ToneYellow},
)`,
		Render: func(p map[string]string) g.Node {
			items := []SnakeTimelineItem{
				{Period: "Nov 2025–Now · 6 mo", Title: "Software Engineer (Thesis)", Org: "Dynatrace", Desc: "Building a Prolog-based permission metamodel for abstracting IAM systems.", Tone: token.ToneCyan},
				{Period: "Oct 2024–Jun 2026", Title: "MSc Networks & IT Security", Org: "JKU Linz", Desc: "Diploma thesis + coursework in network security, cryptography, and formal methods.", Tone: token.ToneViolet},
				{Period: "Mar 2024–Nov 2025 · 1y 8mo", Title: "Software Developer", Org: "ventopay", Desc: "Go backend services, API design, SQLite, CI/CD pipelines.", Tone: token.ToneLime},
				{Period: "Oct 2023–Mar 2024 · 6mo", Title: "Digital Circuits Tutor", Org: "JKU", Desc: "Lab instructor for digital circuits course.", Tone: token.ToneSky},
				{Period: "Sep 2022–Mar 2024", Title: "BSc Computer Science", Org: "JKU Linz", Tone: token.TonePink},
				{Period: "Jun 2021–Sep 2022 · 1y 4mo", Title: "Intern", Org: "Bosch Rexroth", Desc: "Embedded Kotlin + BLE for industrial IoT.", Tone: token.ToneMint},
			}
			return h.Div(
				h.Style("display:flex;flex-direction:column;gap:var(--sp-8)"),
				h.Div(
					h.H4(h.Style("font-size:var(--t-sm);font-weight:700;color:var(--muted);margin-bottom:var(--sp-4)"), g.Text("2 columns (default):")),
					SnakeTimeline(SnakeTimelineProps{Cols: 2}, items...),
				),
				h.Div(
					h.H4(h.Style("font-size:var(--t-sm);font-weight:700;color:var(--muted);margin-bottom:var(--sp-4)"), g.Text("3 columns:")),
					SnakeTimeline(SnakeTimelineProps{Cols: 3}, items...),
				),
			)
		},
	})
}
