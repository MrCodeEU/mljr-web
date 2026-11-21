package pages

import (
	"sort"
	"strings"
	"time"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	hpdata "mljr-web/projects/homepage/data"
	"mljr-web/ui/layout"
	"mljr-web/ui/special"
	"mljr-web/ui/swiss"
	"mljr-web/ui/token"
)

// CV is the compact profile at /cv: the same data as the homepage, laid out
// as a plain one-column résumé that scans in a minute and prints on A4.
func CV(d hpdata.SiteData, lang string, a AnalyticsConfig) g.Node {
	hero := d.ContentFor(lang).Hero
	var work, edu []hpdata.TimelineItem
	for _, t := range d.Timeline {
		switch t.Kind {
		case "work":
			work = append(work, t)
		case "education":
			edu = append(edu, t)
		} // theses are listed with their PDFs below
	}
	newestFirst := func(s []hpdata.TimelineItem) {
		sort.SliceStable(s, func(i, j int) bool { return s[i].Start > s[j].Start })
	}
	newestFirst(work)
	newestFirst(edu)

	title := tr(lang, "Kurzprofil", "Profile") + " - Michael Reinegger"
	return layout.PageShell(
		layout.PageProps{
			Title:       title,
			Description: tr(lang, "Lebenslauf von Michael Reinegger auf einer Seite: Stationen, Ausbildung, Kompetenzen und ausgewählte Projekte.", "Michael Reinegger's résumé on one page: positions, education, skills and selected projects."),
			Theme:       token.ThemeSwiss,
			AutoMode:    true,
			LockTheme:   true,
			Mode:        token.ModeLight,
			Lang:        lang,
			HeadExtra:   append([]g.Node{g.El("style", g.Raw(homepageCSS+cvCSS))}, AnalyticsHead(a)...),
		},
		special.ThemeToggleRoot(token.ThemeSwiss, token.ModeLight),
		siteNavbar(lang),
		h.Main(h.ID("main"),
			h.Div(h.Class("sub-page cv-head"),
				subHead(tr(lang, "/de/", "/"), tr(lang, "Start", "Home"), 6,
					h.H1(g.Text("Michael Reinegger")),
					h.P(h.Class("sub-lede"), g.Text(tr(lang, "IT-Security-Engineer · Linz, Österreich", "IT Security Engineer · Linz, Austria"))),
				),
				h.Div(h.Class("sub-body cv-contact"),
					h.Ul(h.Class("cv-links"),
						h.Li(h.A(h.Href("mailto:reinemic2.0@gmail.com"), g.Text("reinemic2.0@gmail.com"))),
						h.Li(h.A(h.Href("https://www.linkedin.com/in/mrcodeeu/"), g.Text("linkedin.com/in/mrcodeeu"))),
						h.Li(h.A(h.Href("https://github.com/MrCodeEU"), g.Text("github.com/MrCodeEU"))),
						h.Li(h.A(h.Href(tr(lang, "/de/", "/")), g.Text("mljr.eu"))),
					),
					h.Div(h.Class("cv-actions"),
						swiss.BtnLink(cvPath(lang), true, []g.Node{g.Attr("download", "")}, g.Text(tr(lang, "PDF-Lebenslauf ↓", "PDF résumé ↓"))),
						swiss.Btn(false, []g.Node{g.Attr("data-on:click", "window.print()")}, g.Text(tr(lang, "Drucken", "Print"))),
						swiss.BtnLink(tr(lang, "/de/", "/")+"#contact", false, nil, g.Text(tr(lang, "Kontakt", "Contact"))),
					),
				),
			),
			g.If(hero.Description != "", subRow("cv-row", tr(lang, "Profil", "Profile"), "",
				h.P(h.Class("cv-profile"), g.Text(hero.Description)),
				g.If(hero.StatusTag != "", swiss.Micro(h.Class("cv-status"), g.Text("● "+hero.StatusTag))),
			)),
			subRow("cv-row", tr(lang, "Berufserfahrung", "Experience"), "", cvEntries(work, lang)),
			subRow("cv-row", tr(lang, "Ausbildung", "Education"), "", cvEntries(edu, lang)),
			cvTheses(d.ThesisFor(lang), lang),
			subRow("cv-row", tr(lang, "Kompetenzen", "Skills"), "", cvSkills()),
			cvProjects(d.FeaturedProjects(), lang),
		),
		siteFooter(lang),
	)
}

func cvEntries(items []hpdata.TimelineItem, lang string) g.Node {
	rows := make([]g.Node, len(items))
	for i, t := range items {
		where := t.Organization
		if t.Location != "" {
			where += " · " + t.Location
		}
		rows[i] = h.Li(h.Class("cv-entry"),
			h.Span(h.Class("cv-when"), g.Text(cvPeriod(t, lang))),
			h.Div(
				h.H3(g.Text(t.TitleFor(lang))),
				h.P(h.Class("cv-where"), g.Text(where)),
				g.If(t.SummaryFor(lang) != "", h.P(g.Text(t.SummaryFor(lang)))),
				g.If(len(t.Tags) > 0, h.P(h.Class("cv-tags"), g.Text(strings.Join(t.Tags, " · ")))),
			),
		)
	}
	return h.Ol(h.Class("cv-list"), g.Group(rows))
}

func cvTheses(theses []hpdata.Thesis, lang string) g.Node {
	if len(theses) == 0 {
		return nil
	}
	theses = append([]hpdata.Thesis(nil), theses...)
	rank := map[string]int{"Master": 0, "Project": 1, "Projekt": 1, "Seminar": 2, "Bachelor": 3, "HTL": 4}
	sort.SliceStable(theses, func(i, j int) bool { return rank[theses[i].Type] < rank[theses[j].Type] })
	rows := make([]g.Node, len(theses))
	for i, t := range theses {
		rows[i] = h.Li(h.Class("cv-entry"),
			h.Span(h.Class("cv-when"), g.Text(t.Type)),
			h.Div(
				h.H3(g.If(t.PDF != "", h.A(h.Href(t.PDF), g.Text(t.Title))), g.If(t.PDF == "", g.Text(t.Title))),
				h.P(g.Text(firstSentence(t.Description, 300))),
			),
		)
	}
	return subRow("cv-row", tr(lang, "Abschlussarbeiten", "Theses & papers"), "", h.Ol(h.Class("cv-list"), g.Group(rows)))
}

func cvSkills() g.Node {
	groups := hpdata.SkillGroups()
	rows := make([]g.Node, 0, 2*len(groups))
	for _, sg := range groups {
		rows = append(rows, h.Dt(g.Text(sg.Label)), h.Dd(g.Text(strings.Join(sg.Skills, ", "))))
	}
	return h.Dl(h.Class("cv-skills"), g.Group(rows))
}

func cvProjects(featured []hpdata.Project, lang string) g.Node {
	// core picks first, then the remaining plates in their curated order
	sort.SliceStable(featured, func(i, j int) bool { return corePicks[featured[i].ID] && !corePicks[featured[j].ID] })
	if len(featured) > 3 {
		featured = featured[:3]
	}
	if len(featured) == 0 {
		return nil
	}
	rows := make([]g.Node, len(featured))
	for i, p := range featured {
		link := p.URL
		if p.HasDetailPage() {
			link = tr(lang, "/de", "") + "/projects/" + p.ID
		}
		rows[i] = h.Li(h.Class("cv-entry"),
			h.Span(h.Class("cv-when"), g.Text(strings.Join(projectTags(p, 2), " · "))),
			h.Div(
				h.H3(h.A(h.Href(link), g.Text(friendlyName(p.Name)))),
				h.P(g.Text(firstSentence(p.DescFor(lang), 240))),
			),
		)
	}
	return subRow("cv-row", tr(lang, "Ausgewählte Projekte", "Selected projects"), "", h.Ol(h.Class("cv-list"), g.Group(rows)))
}

// firstSentence keeps text up to its first full stop, so a description
// ends on a sentence; a longer first sentence is cut at the last word
// boundary before max runes.
func firstSentence(s string, max int) string {
	if i := strings.Index(s, ". "); i > 0 && i < max {
		return s[:i+1]
	}
	r := []rune(s)
	if len(r) <= max {
		return s
	}
	cut := string(r[:max])
	if i := strings.LastIndex(cut, " "); i > 0 {
		cut = cut[:i]
	}
	return strings.TrimRight(cut, " ,;:—-") + " …"
}

// cvPeriod formats "Okt 2026 – heute" in the page language.
func cvPeriod(t hpdata.TimelineItem, lang string) string {
	month := func(s string) string {
		m, err := time.Parse("2006-01", s)
		if err != nil {
			return s // "YYYY" or unknown
		}
		if lang == "de" {
			return []string{"Jän", "Feb", "Mär", "Apr", "Mai", "Jun", "Jul", "Aug", "Sep", "Okt", "Nov", "Dez"}[m.Month()-1] + " " + m.Format("2006")
		}
		return m.Format("Jan 2006")
	}
	end := tr(lang, "heute", "present")
	if t.End != nil {
		end = month(*t.End)
	}
	start := month(t.Start)
	if start == end {
		return start
	}
	return start + " – " + end
}

const cvCSS = `
.cv-head .bh { display: none; }
.cv-head .sub-head h1 { font-size: clamp(2.2rem, 1.2rem + 2.6vw, 4rem); }
.cv-contact { display: flex; flex-wrap: wrap; align-items: flex-start; justify-content: space-between; gap: 1.25rem 2rem; }
.cv-links { list-style: none; margin: 0; padding: 0; display: grid; gap: .35rem; font-family: var(--font-mono); font-size: .85rem; }
.cv-links a { color: var(--ink); text-decoration: none; border-bottom: 2px solid var(--act); }
.cv-links a:hover { background: var(--act); color: var(--act-fg); }
.cv-actions { display: flex; flex-wrap: wrap; gap: .5rem; }
.cv-profile { margin: 0; max-width: 62ch; font-size: 1.05rem; line-height: 1.55; }
.cv-status { display: block; margin-top: .75rem; color: var(--main-text); }
.cv-list { list-style: none; margin: 0; padding: 0; }
.cv-entry { display: grid; grid-template-columns: 11rem minmax(0, 1fr); gap: .3rem 1.5rem; padding: .85rem 0; border-top: 1px solid var(--rule-subtle); }
.cv-entry:first-child { border-top: 0; padding-top: 0; }
.cv-when { font-family: var(--font-mono); font-size: .75rem; font-weight: 700; text-transform: uppercase; letter-spacing: .04em; color: var(--muted); padding-top: .2rem; }
.cv-entry h3 { margin: 0; font-size: 1.05rem; font-weight: 800; letter-spacing: -.01em; }
.cv-entry h3 a { color: inherit; text-decoration: none; border-bottom: 2px solid var(--act); }
.cv-entry h3 a:hover { background: var(--act); color: var(--act-fg); }
.cv-entry p { margin: .2rem 0 0; max-width: 70ch; line-height: 1.5; font-size: .92rem; }
.cv-where { font-weight: 600; }
.cv-tags { font-family: var(--font-mono); font-size: .72rem !important; color: var(--muted); }
.cv-skills { display: grid; grid-template-columns: 11rem minmax(0, 1fr); gap: .5rem 1.5rem; margin: 0; }
.cv-skills dt { font-family: var(--font-mono); font-size: .75rem; font-weight: 700; text-transform: uppercase; letter-spacing: .04em; color: var(--muted); padding-top: .15rem; }
.cv-skills dd { margin: 0; font-size: .92rem; }
@media (max-width: 640px) {
  .cv-entry, .cv-skills { grid-template-columns: minmax(0, 1fr); }
}
@media print {
  @page { size: A4; margin: 12mm 13mm 14mm; }
  html[data-theme], html[data-theme] body { --bg: #fff; --ink: #000; --muted: #444; --rule-subtle: #ccc; --main-text: #000; background: #fff !important; color: #000; }
  body { font-size: 9pt; }
  body > [data-component="colophon"], body > footer, .cv-actions, .sub-back, .bh, [data-slot="skip"], .cv-tags { display: none !important; }
  /* one column: section titles become headings above their content */
  .sub-page { display: block; margin: 0; border: 0; padding: 0; }
  .sub-head, .sub-body { padding: 0; border: 0; }
  .sub-head-in { position: static; gap: .2rem; }
  .cv-head { display: flex; justify-content: space-between; align-items: flex-end; gap: 1rem; padding-bottom: 3mm; border-bottom: 2px solid #000; }
  .cv-head .sub-head h1 { font-size: 20pt; line-height: 1; }
  .cv-head .sub-lede { font-size: 10pt; color: #000; }
  .cv-links { font-size: 8pt; gap: .1rem; text-align: right; }
  .cv-row { margin-top: 4mm; }
  .cv-row .sub-head { border-bottom: 1px solid #000; margin-bottom: 1.5mm; }
  .sub-row-title { font-size: 9pt; letter-spacing: .06em; }
  .cv-profile { font-size: 9pt; max-width: none; }
  .cv-status { margin-top: 1mm; }
  .cv-entry { break-inside: avoid; grid-template-columns: 26mm minmax(0, 1fr); gap: 0 4mm; padding: 1.2mm 0; }
  .cv-entry h3 { font-size: 9.5pt; }
  .cv-entry p { font-size: 8.5pt; margin-top: .3mm; max-width: none; line-height: 1.35; }
  .cv-when { font-size: 7pt; padding-top: .3mm; }
  .cv-skills { grid-template-columns: 26mm minmax(0, 1fr); gap: .8mm 4mm; }
  .cv-skills dt { font-size: 7pt; }
  .cv-skills dd { font-size: 8.5pt; }
  .cv-links a, .cv-entry h3 a { border: 0; color: #000; }
}
`
