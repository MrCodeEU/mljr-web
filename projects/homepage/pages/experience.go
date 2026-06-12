package pages

import (
	"fmt"
	"math"
	"strings"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	hpdata "mljr-web/projects/homepage/data"
	uidata "mljr-web/ui/data"
	"mljr-web/ui/layout"
	"mljr-web/ui/primitive"
	"mljr-web/ui/special"
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
								// Same logo-chip header as the experience snake cards
								h.Div(
									h.Style("display:flex;align-items:center;gap:var(--sp-3);margin-bottom:var(--sp-3)"),
									uidata.OrgLogoChip(edu.OrgLogo, edu.Org),
									h.Div(
										h.Style("min-width:0"),
										h.Div(h.Style("font-size:var(--t-sm);font-weight:800;line-height:1.25"), g.Text(edu.Org)),
										h.Div(h.Style("font-size:var(--t-xs);font-family:var(--font-mono,monospace);font-weight:600;opacity:.65;margin-top:2px"), g.Text(edu.Period)),
									),
								),
								h.H4(h.Style("font-weight:900;font-size:var(--t-base);margin:0;line-height:1.35"), g.Text(edu.Title)),
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
			experienceLocationMap(li),
		),
	)
}

type orgLocation struct {
	lat float64
	lng float64
}

func experienceLocationMap(li hpdata.LinkedInData) g.Node {
	pins := []special.MapPin{
		{Lat: 48.143, Lng: 14.461, Label: "Home · Thaling, Upper Austria", Popup: "<strong>Home · Thaling, Upper Austria</strong><br>Current base"},
	}
	seen := map[string]int{}
	idx := 1
	for _, j := range li.RelevantExperience(100) {
		loc, ok := companyLocation(j.Company)
		if !ok {
			continue
		}
		n := seen[j.Company]
		seen[j.Company] = n + 1
		lat, lng := offsetLocation(loc, n)
		pins = append(pins, special.MapPin{
			Lat:       lat,
			Lng:       lng,
			AnchorLat: loc.lat,
			AnchorLng: loc.lng,
			Label:     j.Company,
			Popup:     orgPopup(idx, j.Company, hpdata.LogoForCompany(j.Company), j.Title, j.Period),
			Icon:      hpdata.LogoForCompany(j.Company),
		})
		idx++
	}
	for _, e := range li.Education {
		loc, ok := schoolLocation(e.School)
		if !ok {
			continue
		}
		n := seen[e.School]
		seen[e.School] = n + 1
		lat, lng := offsetLocation(loc, n)
		pins = append(pins, special.MapPin{
			Lat:       lat,
			Lng:       lng,
			AnchorLat: loc.lat,
			AnchorLng: loc.lng,
			Label:     e.School,
			Popup:     orgPopup(idx, e.School, hpdata.LogoForSchool(e.School), e.Degree, e.Period),
			Icon:      hpdata.LogoForSchool(e.School),
		})
		idx++
	}

	return h.Div(
		h.Style("margin-top:var(--sp-12)"),
		sectionHeader("", "Places", "work · education · home", token.ToneMint),
		special.OpenMap(special.OpenMapProps{
			CenterLat: 48.22,
			CenterLng: 14.34,
			Zoom:      9,
			Height:    "360px",
			ID:        "experience-map",
		}, pins...),
	)
}

func companyLocation(company string) (orgLocation, bool) {
	m := map[string]orgLocation{
		"Dynatrace":                        {48.3069, 14.2858},
		"Johannes Kepler Universität Linz": {48.3371, 14.3196},
		"ventopay gmbh":                    {48.3678, 14.5165},
		"Bosch":                            {48.3069, 14.2858},
		"Bosch Rexroth":                    {48.2462, 14.2348},
		"ENGEL":                            {48.2735, 14.5861},
		"HerzReha Bad Ischl":               {47.7111, 13.6239},
	}
	loc, ok := m[company]
	return loc, ok
}

func schoolLocation(school string) (orgLocation, bool) {
	m := map[string]orgLocation{
		"Johannes Kepler Universität Linz": {48.3371, 14.3196},
		"HTL Steyr":                        {48.0427, 14.4213},
	}
	loc, ok := m[school]
	return loc, ok
}

func offsetLocation(loc orgLocation, n int) (float64, float64) {
	if n == 0 {
		return loc.lat, loc.lng
	}
	angle := float64(n) * (math.Pi * 2 / 7)
	radius := 0.013 + float64(n/7)*0.004
	return loc.lat + math.Sin(angle)*radius, loc.lng + math.Cos(angle)*radius*1.35
}

func orgPopup(n int, org, logo, title, period string) string {
	logoHTML := ""
	if logo != "" {
		logoHTML = fmt.Sprintf(`<img src="%s" alt="" style="width:34px;height:34px;object-fit:contain;border:2px solid #111;background:#fff;margin-right:8px;vertical-align:middle">`, logo)
	}
	return fmt.Sprintf(`<div style="display:flex;align-items:center">%s<div><strong>%02d · %s</strong><br>%s<br>%s</div></div>`, logoHTML, n, org, title, period)
}
