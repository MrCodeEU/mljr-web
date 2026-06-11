package pages

import (
	"fmt"
	"strings"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	hpdata "mljr-web/projects/homepage/data"
	uidata "mljr-web/ui/data"
	"mljr-web/ui/layout"
	"mljr-web/ui/primitive"
	"mljr-web/ui/token"
)

func experienceSection(li hpdata.LinkedInData) g.Node {
	jobs := li.RelevantExperience(100) // all positions

	tones := []token.Tone{token.ToneCyan, token.ToneViolet, token.ToneLime, token.ToneSky, token.TonePink, token.ToneMint, token.ToneBlush, token.ToneYellow}

	snakeItems := make([]uidata.SnakeTimelineItem, 0, len(jobs))
	mobileItems := make([]uidata.TimelineItem, 0, len(jobs))
	for i, j := range jobs {
		desc := j.Desc
		if len(desc) < 30 || strings.Contains(desc, "Österreich") || strings.Contains(desc, "sterreich") {
			desc = ""
		}
		tags := []string{j.Type}
		if j.Type == "" {
			tags = nil
		}
		snakeItems = append(snakeItems, uidata.SnakeTimelineItem{
			Period:  j.Period + " · " + j.Duration,
			Title:   j.Title,
			Org:     j.Company,
			OrgLogo: hpdata.LogoForCompany(j.Company),
			Desc:    desc,
			Tags:    tags,
			Tone:    tones[i%len(tones)],
		})
		mobileItems = append(mobileItems, uidata.TimelineItem{
			Period:  j.Period + " · " + j.Duration,
			Title:   j.Title,
			Org:     j.Company,
			OrgLogo: hpdata.LogoForCompany(j.Company),
			Desc:    desc,
			Tags:    tags,
			Tone:    tones[i%len(tones)],
		})
	}

	eduItems := make([]uidata.TimelineItem, 0, len(li.Education))
	eduTones := []token.Tone{token.ToneAccent, token.ToneSky, token.ToneMint}
	for i, e := range li.Education {
		eduItems = append(eduItems, uidata.TimelineItem{
			Period:  e.Period,
			Title:   e.Degree,
			Org:     e.School,
			OrgLogo: hpdata.LogoForSchool(e.School),
			Tone:    eduTones[i%len(eduTones)],
		})
	}

	return h.Section(
		h.ID("experience"),
		h.Style("padding:var(--sp-12) 0"),
		layout.Container(layout.ContainerProps{},
			sectionHeader("01", "Experience", fmt.Sprintf("%d positions", len(jobs)), token.ToneCyan),
			// Snake timeline spans full width
			h.Div(h.Class("experience-snake"), uidata.SnakeTimeline(uidata.SnakeTimelineProps{Cols: 3}, snakeItems...)),
			h.Div(h.Class("experience-mobile-timeline"), uidata.Timeline(uidata.TimelineProps{}, mobileItems...)),
			// Education row below
			h.Div(
				h.Style("margin-top:var(--sp-12)"),
				sectionHeader("", "Education", fmt.Sprintf("%d degrees", len(li.Education)), token.ToneSky),
				h.Div(
					h.Style("display:grid;grid-template-columns:repeat(auto-fill,minmax(260px,1fr));gap:var(--sp-4)"),
					g.Group(func() []g.Node {
						nodes := make([]g.Node, len(eduItems))
						for i, edu := range eduItems {
							nodes[i] = primitive.Card(primitive.CardProps{Tone: edu.Tone},
								h.Div(
									h.Style("display:flex;align-items:flex-start;justify-content:space-between;gap:var(--sp-3);margin-bottom:var(--sp-2)"),
									h.Span(h.Style("font-size:var(--t-xs);font-weight:700;opacity:.6;text-transform:uppercase;letter-spacing:.06em"), g.Text(edu.Period)),
									g.If(edu.OrgLogo != "", h.Img(h.Src(edu.OrgLogo), h.Alt(edu.Org), h.Style("width:32px;height:32px;object-fit:contain;border-radius:4px"))),
								),
								h.H4(h.Style("font-weight:800;font-size:var(--t-base);margin:0 0 var(--sp-1)"), g.Text(edu.Title)),
								h.Div(h.Style("font-size:var(--t-sm);opacity:.7"), g.Text(edu.Org)),
							)
						}
						return nodes
					}()),
				),
				// Thesis callout
				h.Div(
					h.Style("margin-top:var(--sp-4)"),
					primitive.Callout(primitive.CalloutProps{Variant: primitive.CalloutInfo},
						h.Strong(g.Text("Current thesis: ")),
						g.Text("Abstracting Permission Systems → Prolog-based metamodel @ Dynatrace"),
					),
				),
			),
		),
	)
}
