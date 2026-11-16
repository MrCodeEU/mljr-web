package pages

import (
	"fmt"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	hpdata "mljr-web/projects/homepage/data"
	"mljr-web/projects/homepage/homelab"
	"mljr-web/ui/layout"
	"mljr-web/ui/overlay"
	"mljr-web/ui/special"
	"mljr-web/ui/swiss"
	"mljr-web/ui/token"
)

// homeSection is one row of the homepage sheet. Num is assigned in order so
// optional sections (live tools, Strava, papers) never leave gaps.
type homeSection struct {
	id, de, en string
	render     func(p swiss.SectionProps) g.Node
}

// sectionField gives some anchors a solid colour field, so the sheet carries
// large flat colour areas down the page instead of only at the top.
var sectionField = map[string]string{
	"projects": "second", "playground": "third", "experience": "block",
	"homelab": "second", "skills": "third", "contact": "ink",
}

func Home(d hpdata.SiteData, lang string, a AnalyticsConfig, hl homelab.Snapshot) g.Node {
	featured := d.FeaturedProjects()
	rest := d.AllProjects()
	if len(featured) == 0 {
		featured, rest = rest, nil
	}
	liveTools := d.LiveToolProjects()
	title := tr(lang, "Michael Reinegger — IT-Security-Engineer in Linz", "Michael Reinegger — IT Security Engineer, Linz")
	ogDescription := tr(lang,
		"Portfolio von Michael Reinegger: Dipl.-Ing. Netzwerke & IT-Sicherheit (JKU Linz), Go, selbst gehostete Infrastruktur, Projekte, Laufbahn und Homelab live.",
		"Portfolio of Michael Reinegger: Dipl.-Ing. in Networks & IT Security (JKU Linz), Go, self-hosted infrastructure, projects, career and a live homelab.")
	ogURL := tr(lang, "https://mljr.eu/de/", "https://mljr.eu/")
	headExtra := append([]g.Node{
		g.El("style", g.Raw(homepageCSS+homeCSS+chartCSS+netCSS+careerMapCSS+leitstandCSS+hoodCSS+contactCSS+bewegungCSS+catalogueCSS+tabsCSS+labCSS+skillsCSS)),
		h.Script(h.Type("application/ld+json"), g.Raw(personJSONLD)),
		g.El("meta", g.Attr("property", "og:title"), g.Attr("content", title)),
		g.El("meta", g.Attr("property", "og:description"), g.Attr("content", ogDescription)),
		g.El("meta", g.Attr("property", "og:type"), g.Attr("content", "website")),
		g.El("meta", g.Attr("property", "og:url"), g.Attr("content", ogURL)),
		g.El("meta", g.Attr("property", "og:locale"), g.Attr("content", tr(lang, "de_AT", "en_US"))),
		g.El("meta", g.Attr("property", "og:image"), g.Attr("content", "https://mljr.eu/assets/portfolio/mljr-web/0.webp")),
		g.El("meta", g.Attr("name", "twitter:card"), g.Attr("content", "summary_large_image")),
		g.El("meta", g.Attr("name", "twitter:title"), g.Attr("content", title)),
		g.El("meta", g.Attr("name", "twitter:description"), g.Attr("content", ogDescription)),
		g.El("meta", g.Attr("name", "twitter:image"), g.Attr("content", "https://mljr.eu/assets/portfolio/mljr-web/0.webp")),
	}, AnalyticsHead(a)...)

	all := []homeSection{
		{"hero", "Identität", "Identity", func(p swiss.SectionProps) g.Node { return heroSection(p, d, lang) }},
		{"projects", "Werkkatalog", "Catalogue", func(p swiss.SectionProps) g.Node { return catalogueSection(p, featured, rest, liveTools, lang) }},
		{"playground", "Labor", "Kinetic Lab", func(p swiss.SectionProps) g.Node { return labSection(p, lang) }},
		{"opensource", "Quelltext", "Open Source", func(p swiss.SectionProps) g.Node { return githubSection(p, d, lang) }},
		{"experience", "Laufbahn", "Career Network", func(p swiss.SectionProps) g.Node { return careerSection(p, d, lang) }},
		{"papers", "Schriften", "Papers", func(p swiss.SectionProps) g.Node { return papersSection(p, d, lang) }},
		{"homelab", "Heimnetz", "Homelab", func(p swiss.SectionProps) g.Node { return homelabSection(p, hl, lang) }},
		{"activity", "Bewegung", "Training", func(p swiss.SectionProps) g.Node { return stravaSection(p, d, lang) }},
		{"skills", "Kompetenzen", "Skills", func(p swiss.SectionProps) g.Node { return skillsSection(p, lang) }},
		{"under-the-hood", "Maschinenraum", "Under the Hood", func(p swiss.SectionProps) g.Node { return codeShowcaseSection(p, lang) }},
		{"contact", "Kontakt", "Contact", func(p swiss.SectionProps) g.Node { return contactSection(p, d.ContentFor(lang).Contact, lang) }},
	}
	var nodes []g.Node
	var items []swiss.NavItem
	for _, s := range all {
		num := fmt.Sprintf("%02d", len(items))
		n := s.render(swiss.SectionProps{ID: s.id, Num: num, Title: s.de, Sub: s.en, Field: sectionField[s.id], AnchorExtra: bauhausMark(len(items))})
		if n == nil {
			continue
		}
		nodes = append(nodes, n)
		items = append(items, swiss.NavItem{ID: s.id, Num: num, Label: s.de, Sub: s.en})
	}
	nav := swiss.ColophonProps{
		Brand:     "mljr.eu",
		BrandSub:  "— Michael Reinegger · Linz (AT)",
		Items:     items,
		Actions:   siteActions(lang),
		PageID:    "page",
		SkipLabel: tr(lang, "Zum Inhalt springen", "Skip to content"),
	}

	return layout.PageShell(
		layout.PageProps{
			Title:       title,
			Description: ogDescription,
			Theme:       token.ThemeSwiss,
			AutoMode:    true,
			LockTheme:   true,
			Mode:        token.ModeLight,
			HeadExtra:   headExtra,
			Lang:        lang,
		},
		special.ThemeToggleRoot(token.ThemeSwiss, token.ModeLight),
		swiss.Colophon(nav),
		h.Div(h.ID("page"),
			swiss.Page("sheet",
				special.LogoGrid(special.LogoGridProps{SlotID: "logo-slot", SheetID: "sheet", HeroID: "hero"}),
				h.Main(append([]g.Node{h.ID("main")}, nodes...)...),
			),
			siteFooter(lang),
		),
		swiss.Register(nav),
		overlay.Toaster(overlay.ToasterProps{}),
		overlay.Portal("portal"),
		h.Script(g.Raw(tabsJS)),
	)
}

// personJSONLD tells search engines who the site is about (schema.org Person).
const personJSONLD = `{"@context":"https://schema.org","@type":"Person","name":"Michael Reinegger","url":"https://mljr.eu/","jobTitle":"IT Security Engineer","address":{"@type":"PostalAddress","addressLocality":"Linz","addressCountry":"AT"},"alumniOf":[{"@type":"CollegeOrUniversity","name":"Johannes Kepler University Linz"},{"@type":"HighSchool","name":"HTL Steyr"}],"knowsAbout":["IT Security","Go","Networks","Self-hosting"],"sameAs":["https://github.com/MrCodeEU","https://www.linkedin.com/in/mrcodeeu/"]}`
