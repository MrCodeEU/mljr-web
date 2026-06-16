package pages

import (
	"fmt"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"mljr-web/internal/i18n"
	hpdata "mljr-web/projects/homepage/data"
	uidata "mljr-web/ui/data"
	"mljr-web/ui/icon"
	"mljr-web/ui/layout"
	"mljr-web/ui/primitive"
	"mljr-web/ui/special"
	"mljr-web/ui/token"
)

func experienceSection(d hpdata.SiteData, lang string) g.Node {
	jobs := d.WorkItems()
	edu := d.EduItems()

	tones := []token.Tone{token.ToneCyan, token.ToneViolet, token.ToneLime, token.ToneSky, token.TonePink, token.ToneMint, token.ToneBlush, token.ToneYellow}

	snakeItems := make([]uidata.SnakeTimelineItem, 0, len(jobs))
	mobileItems := make([]uidata.TimelineItem, 0, len(jobs))
	for i, j := range jobs {
		details := j.DetailsFor(lang)
		desc := ""
		if len(details) > 0 {
			desc = details[0]
		} else if j.Summary != "" {
			desc = j.Summary
		}
		tagNodes := buildTagNodes(j.Tags)
		snakeItems = append(snakeItems, uidata.SnakeTimelineItem{
			Period:   j.FormatPeriod() + " · " + j.FormatDuration(),
			Title:    j.TitleFor(lang),
			Org:      j.Organization,
			OrgLogo:  j.Logo,
			Desc:     desc,
			TagNodes: tagNodes,
			Tone:     tones[i%len(tones)],
		})
		mobileItems = append(mobileItems, uidata.TimelineItem{
			Period:   j.FormatPeriod() + " · " + j.FormatDuration(),
			Title:    j.TitleFor(lang),
			Org:      j.Organization,
			OrgLogo:  j.Logo,
			Desc:     desc,
			TagNodes: tagNodes,
			Tone:     tones[i%len(tones)],
		})
	}

	eduTones := []token.Tone{token.ToneAccent, token.ToneSky, token.ToneMint}
	eduItems := make([]uidata.TimelineItem, 0, len(edu))
	for i, e := range edu {
		eduItems = append(eduItems, uidata.TimelineItem{
			Period:  e.FormatPeriod(),
			Title:   e.TitleFor(lang),
			Org:     e.Organization,
			OrgLogo: e.Logo,
			Desc:    e.SummaryFor(lang),
			Tone:    eduTones[i%len(eduTones)],
		})
	}

	return h.Section(
		h.ID("experience"),
		h.Style("padding:var(--sp-12) 0"),
		layout.Container(layout.ContainerProps{},
			sectionHeader("05", i18n.T(lang, "sections.experience.title"), fmt.Sprintf(i18n.T(lang, "sections.experience.positions"), len(jobs)), token.ToneCyan),
			// Snake timeline spans full width
			h.Div(h.Class("experience-snake"), uidata.SnakeTimeline(uidata.SnakeTimelineProps{Cols: 3}, snakeItems...)),
			h.Div(h.Class("experience-mobile-timeline"), uidata.Timeline(uidata.TimelineProps{}, mobileItems...)),
			// Education row below
			h.Div(
				h.Style("margin-top:var(--sp-12)"),
				sectionHeader("", i18n.T(lang, "sections.education.title"), fmt.Sprintf(i18n.T(lang, "sections.education.degrees"), len(edu)), token.ToneSky),
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
							g.If(edu.Desc != "",
								h.P(h.Style("margin:var(--sp-2) 0 0;font-size:var(--t-sm);line-height:1.5;opacity:.8"), g.Text(edu.Desc)),
							),
						)
					}
					return nodes
				}()),
				),
				// Thesis section
				g.If(len(d.ThesisFor(lang)) > 0,
					h.Div(
						h.Style("margin-top:var(--sp-10)"),
						sectionHeader("", i18n.T(lang, "sections.papers.title"), i18n.T(lang, "sections.papers.sub"), token.ToneViolet),
						h.Div(
							h.Style("display:grid;grid-template-columns:repeat(auto-fill,minmax(260px,1fr));gap:var(--sp-4)"),
							g.Group(func() []g.Node {
								thesis := d.ThesisFor(lang)
								nodes := make([]g.Node, 0, len(thesis))
								for _, t := range thesis {
									nodes = append(nodes, primitive.Card(primitive.CardProps{Tone: thesisTone(t.Type), Attrs: []g.Node{h.Style("position:relative")}},
										// Corner badge: center of tag sits at top-right card corner
										h.Div(h.Style("position:absolute;top:0;right:0;transform:translate(50%,-50%);z-index:2"),
											primitive.Tag(primitive.TagProps{}, g.Text(t.Type)),
										),
										h.Div(h.Style("display:flex;align-items:center;gap:var(--sp-2);margin-bottom:var(--sp-2)"),
											g.If(thesisLogo(t.Type) != "",
												uidata.OrgLogoChip(thesisLogo(t.Type), ""),
											),
											h.H4(h.Style("font-weight:900;font-size:var(--t-base);margin:0;line-height:1.35"), g.Text(t.Title)),
										),
										h.P(h.Style("margin:0;font-size:var(--t-sm);line-height:1.55;opacity:.85"), g.Text(t.Description)),
										g.If(t.PDF != "",
											h.A(h.Href(t.PDF), g.Attr("target", "_blank"), g.Attr("rel", "noopener"), h.Style("margin-top:var(--sp-3);display:block"),
												primitive.Button(primitive.ButtonProps{Variant: token.Outline, Size: token.SizeSM},
													icon.Icon("lucide:file-text"),
													g.Text(i18n.T(lang, "sections.experience.thesis_view_pdf")),
												),
											),
										),
										g.If(t.PDF == "",
											h.Div(h.Style("margin-top:var(--sp-3)"),
												primitive.Tag(primitive.TagProps{}, icon.Icon("lucide:clock"), g.Text(i18n.T(lang, "sections.experience.thesis_coming_soon"))),
											),
										),
									))
								}
								return nodes
							}()),
						),
					),
				),
			),
			experienceLocationMap(d, lang),
		),
	)
}

type orgLocation struct {
	lat float64
	lng float64
}

func experienceLocationMap(d hpdata.SiteData, lang string) g.Node {
	pins := []special.MapPin{
		{Lat: 48.143, Lng: 14.461, Label: "Home · Thaling, Upper Austria", Popup: "<strong>Home · Thaling, Upper Austria</strong><br>Current base"},
	}
	// Count occurrences per anchor first: locations with >1 pin are excluded
	// from clustering entirely and instead spread around their anchor on
	// fixed pixel-radius spokes (so they never get hidden behind a cluster
	// bubble for their own siblings).
	counts := map[orgLocation]int{}
	for _, j := range d.WorkItems() {
		if loc, ok := orgLocation2(j.Organization); ok {
			counts[loc]++
		}
	}
	for _, e := range d.EduItems() {
		if loc, ok := orgLocation2(e.Organization); ok {
			counts[loc]++
		}
	}
	seen := map[orgLocation]int{}
	spread := func(loc orgLocation) (lineLat, lineLng, angle, radius float64) {
		n := seen[loc]
		seen[loc] = n + 1
		if counts[loc] <= 1 {
			return 0, 0, 0, 0
		}
		return loc.lat, loc.lng, float64(n) * 2.399963, 22 + 18*float64(n) // golden angle, growing radius
	}
	idx := 1
	for _, j := range d.WorkItems() {
		loc, ok := orgLocation2(j.Organization)
		if !ok {
			continue
		}
		lineLat, lineLng, angle, radius := spread(loc)
		pins = append(pins, special.MapPin{
			AnchorLat:    loc.lat,
			AnchorLng:    loc.lng,
			LineLat:      lineLat,
			LineLng:      lineLng,
			SpreadAngle:  angle,
			SpreadRadius: radius,
			Label:        j.Organization,
			Popup:        orgPopup(idx, j.Organization, j.Logo, j.Title, j.FormatPeriod()),
			Icon:         j.Logo,
		})
		idx++
	}
	for _, e := range d.EduItems() {
		loc, ok := orgLocation2(e.Organization)
		if !ok {
			continue
		}
		lineLat, lineLng, angle, radius := spread(loc)
		pins = append(pins, special.MapPin{
			AnchorLat:    loc.lat,
			AnchorLng:    loc.lng,
			LineLat:      lineLat,
			LineLng:      lineLng,
			SpreadAngle:  angle,
			SpreadRadius: radius,
			Label:        e.Organization,
			Popup:        orgPopup(idx, e.Organization, e.Logo, e.Title, e.FormatPeriod()),
			Icon:         e.Logo,
		})
		idx++
	}

	return h.Div(
		h.Style("margin-top:var(--sp-12)"),
		sectionHeader("", i18n.T(lang, "sections.places.title"), i18n.T(lang, "sections.places.sub"), token.ToneMint),
		special.OpenMap(special.OpenMapProps{
			CenterLat: 48.22,
			CenterLng: 14.34,
			Zoom:      9,
			Height:    "520px",
			ID:        "experience-map",
		}, pins...),
	)
}

// orgLocation2 maps organization names (EN or DE) to map coordinates.
func orgLocation2(org string) (orgLocation, bool) {
	m := map[string]orgLocation{
		"Dynatrace":                        {48.3069, 14.2858},
		"Johannes Kepler University Linz":  {48.3371, 14.3196},
		"Johannes Kepler Universität Linz": {48.3371, 14.3196},
		"ventopay gmbh":                    {48.3678, 14.5165},
		"Bosch":                            {48.3069, 14.2858},
		"Bosch Rexroth":                    {48.2462, 14.2348},
		"ENGEL":                            {48.2735, 14.5861},
		"HerzReha Bad Ischl":               {47.7111, 13.6239},
		"HTL Steyr":                        {48.0427, 14.4213},
	}
	loc, ok := m[org]
	return loc, ok
}

func buildTagNodes(tags []string) []g.Node {
	nodes := make([]g.Node, 0, len(tags))
	for _, t := range tags {
		if t == "" {
			continue
		}
		nodes = append(nodes, primitive.Tag(
			primitive.TagProps{Icon: hpdata.TechIcon(t), Tone: hpdata.TagTone(t)},
			g.Text(t),
		))
	}
	return nodes
}

func thesisLogo(typ string) string {
	switch typ {
	case "HTL":
		return "/static/logos/htl-steyr.png"
	case "Bachelor", "Master", "Project", "Projekt":
		return "/static/logos/jku.jpg"
	}
	return ""
}

func thesisTone(typ string) token.Tone {
	switch typ {
	case "BACHELOR", "Bachelor":
		return token.ToneSky
	case "MASTER", "Master":
		return token.ToneViolet
	case "PROJEKT", "PROJECT", "Projekt", "Project":
		return token.ToneLime
	default:
		return token.ToneAccent
	}
}

func orgPopup(n int, org, logo, title, period string) string {
	logoHTML := ""
	if logo != "" {
		logoHTML = fmt.Sprintf(`<img class="map-tooltip-logo" src="%s" alt="">`, logo)
	}
	return fmt.Sprintf(`<div class="map-tooltip-body">%s<div><div class="map-tooltip-title">%02d · %s</div><div class="map-tooltip-meta">%s<br>%s</div></div></div>`, logoHTML, n, org, title, period)
}
