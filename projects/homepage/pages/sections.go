package pages

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"mljr-web/internal/i18n"
	hpdata "mljr-web/projects/homepage/data"
	"mljr-web/ui/icon"
	"mljr-web/ui/primitive"
	"mljr-web/ui/swiss"
)

// ── 00 Identität ────────────────────────────────────────────────────────────

func heroSection(p swiss.SectionProps, d hpdata.SiteData, lang string) g.Node {
	li := d.LinkedIn
	hero := d.ContentFor(lang).Hero
	b := hero.Bento
	p.Title = "Identität & System"
	p.Sub = "Identity & System"
	p.Note = b.Education.Value
	p.AnchorExtra = h.Div(h.ID("logo-slot"), h.Class("hero-logo-slot"))
	p.Attrs = []g.Node{h.Class("hero")}

	tw := hero.TaglineLines
	fact := func(c hpdata.BentoCell, n string) g.Node {
		return swiss.Cell(swiss.CellProps{Span: 3},
			swiss.Kicker(n+". "+c.Label, nil),
			h.P(h.Class("hero-fact"), g.Text(c.Value)),
			g.If(c.Sub != "", h.P(h.Class("hero-fact-sub"), g.Text(c.Sub))),
		)
	}
	return swiss.Section(p,
		swiss.Cell(swiss.CellProps{Span: 6, Attrs: []g.Node{h.Class("hero-head")}},
			h.Div(
				swiss.Micro(h.Class("hero-status"), g.Text("● "+hero.StatusTag)),
				h.H1(h.Class("hero-name"), g.Text("Michael "), h.Br(), g.Text("Reinegger")),
				g.If(len(tw) > 0, h.P(h.Class("hero-type"), swiss.SROnly(g.Text(strings.Join(tw, ", "))),
					h.Span(g.Attr("aria-hidden", "true"), g.Text("› "),
						primitive.Typewriter(primitive.TypewriterProps{Lines: tw, Speed: 70, DeleteSpeed: 35, Pause: 1800, ID: "hero-tw"}),
					),
				)),
			),
			h.Div(
				h.P(h.Class("hero-lead"), g.Text(hero.Description)),
				h.Div(h.Class("hero-ctas"),
					swiss.BtnLink("#projects", true, nil, g.Text(i18n.T(lang, "hero.cta_projects")+" ↓")),
					swiss.BtnLink(tr(lang, "/de/cv", "/cv"), false, nil, g.Text(tr(lang, "Kurzprofil · 60 s", "Profile · 60 s"))),
					swiss.BtnLink("https://github.com/MrCodeEU", false, []g.Node{g.Attr("target", "_blank"), g.Attr("rel", "noopener noreferrer")},
						icon.Icon("simple-icons:github", icon.Props{Size: "1em"}), g.Text(i18n.T(lang, "hero.cta_github"))),
					swiss.BtnLink(cvPath(lang), false, []g.Node{g.Attr("download", "")},
						icon.Icon("lucide:download", icon.Props{Size: "1em"}), g.Text(i18n.T(lang, "hero.cta_download_cv"))),
				),
			),
		),
		// one column beside the name: the blue status block stretches, the
		// photo below it stays square
		swiss.Cell(swiss.CellProps{Span: 3, Variant: "flush", Attrs: []g.Node{h.Class("hero-side")}},
			h.Div(h.Class("hero-block"),
				swiss.Kicker(b.Status.Label, g.Text("Live")),
				h.P(h.Class("hero-credo"), g.Text(b.Status.Value)),
				h.P(h.Class("hero-credo-sub"), g.Text(b.Status.Sub)),
				swiss.Micro(h.Class("hero-scroll"), g.Text(tr(lang, "Scrollen ↓ zum Zerlegen", "Scroll ↓ to disassemble"))),
			),
			g.If(li.Profile.PhotoURL != "", h.Div(h.Class("hero-photo"),
				// above the fold: load eagerly, it is the likely LCP image
				h.Img(h.Src(li.Profile.PhotoURL), h.Alt(i18n.T(lang, "hero.photo_name")), g.Attr("fetchpriority", "high")),
				h.Div(h.Class("hero-photo-cap"),
					swiss.Micro(g.Text("Abb. 1 — "+i18n.T(lang, "hero.photo_name"))),
					swiss.Micro(g.Text(i18n.T(lang, "hero.photo_location"))),
				),
			)),
		),
		fact(b.Focus, "A"),
		fact(b.Education, "B"),
		fact(b.Homelab, "C"),
	)
}

// ── 02 Werkkatalog ──────────────────────────────────────────────────────────

func projectTags(p hpdata.Project, limit int) []string {
	seen := map[string]bool{}
	var out []string
	add := func(t string) {
		k := strings.ToLower(t)
		if t == "" || seen[k] || k == "live-tool" {
			return
		}
		seen[k] = true
		out = append(out, t)
	}
	add(p.Language)
	for _, t := range p.Topics {
		add(t)
	}
	if limit > 0 && len(out) > limit {
		out = out[:limit]
	}
	return out
}

func tagPills(tags []string) g.Node {
	nodes := make([]g.Node, len(tags))
	for i, t := range tags {
		ic := hpdata.TechIcon(t)
		nodes[i] = swiss.Pill(g.If(ic != "", icon.Icon(ic, icon.Props{Size: "1em"})), g.Text(t))
	}
	return swiss.Pills(nodes...)
}

func projectLinks(p hpdata.Project, lang string) g.Node {
	var nodes []g.Node
	ext := []g.Node{g.Attr("target", "_blank"), g.Attr("rel", "noopener noreferrer")}
	if HasProjectDetail(p) {
		nodes = append(nodes, swiss.BtnLink("/projects/"+p.ID, true, nil, g.Text(i18n.T(lang, "projects.details")+" →")))
	}
	if p.URL != "" {
		nodes = append(nodes, swiss.BtnLink(p.URL, false, ext, icon.Icon("simple-icons:github", icon.Props{Size: "1em"}), g.Text(i18n.T(lang, "projects.source")+" ↗")))
	}
	for _, l := range p.Links {
		if l.URL != "" {
			nodes = append(nodes, swiss.BtnLink(l.URL, false, ext, g.Text(l.Name+" ↗")))
		}
	}
	return h.Div(h.Class("cat-links"), g.Group(nodes))
}

// corePicks are the projects a visitor short on time should open first;
// they get a "Kernprojekt" mark on their plate and plate tab.
var corePicks = map[string]bool{"homelab-automation": true, "mljr-web": true}

// corePick is the mark shown on core projects and the master's thesis.
func corePick(lang string) g.Node {
	return h.Span(h.Class("core-pick"), g.Text(tr(lang, "★ Kernprojekt", "★ Core pick")))
}

func catalogueSection(p swiss.SectionProps, featured, rest, tools []hpdata.Project, lang string) g.Node {
	p.Note = tr(lang, "Tafeln für Hauptprojekte, Kleintafeln und Verzeichnis für den Rest.", "Plates for the main projects, small plates and an index for the rest.")
	gal := func(pr hpdata.Project, aspect string) g.Node {
		return swiss.Gallery(swiss.GalleryProps{
			Images: pr.LocalImages(), Alt: friendlyName(pr.Name), Aspect: aspect,
			Label: tr(lang, "Abb.", "Fig."), Prev: tr(lang, "Vorheriges Bild", "Previous image"), Next: tr(lang, "Nächstes Bild", "Next image"),
		})
	}
	nodes := []g.Node{swiss.GalleryScript()}

	// Featured: one plate at a time, picked from a numbered plate index.
	roman := []string{"I", "II", "III", "IV", "V", "VI", "VII", "VIII", "IX", "X"}
	tabs := make([]g.Node, len(featured))
	plates := make([]g.Node, len(featured))
	for i, pr := range featured {
		num := strconv.Itoa(i + 1)
		if i < len(roman) {
			num = roman[i]
		}
		n := strconv.Itoa(i + 1)
		tabs[i] = h.Button(h.Type("button"), h.Class("plate-tab"), g.Attr("data-tab", n), g.Attr("aria-pressed", strconv.FormatBool(i == 0)),
			h.Span(h.Class("plate-tab-n"), g.Text(num)),
			h.Span(h.Class("plate-tab-t"), g.Text(friendlyName(pr.Name))),
			g.If(corePicks[pr.ID], h.Span(h.Class("plate-tab-core"), g.Text(tr(lang, "★ Kern", "★ Core")))),
		)
		media := gal(pr, "4/3")
		if media == nil {
			media = monogram(pr, i)
		}
		plates[i] = h.Article(h.Class("plate"), g.Attr("data-tab-page", n), g.If(i > 0, g.Attr("hidden", "")),
			h.Div(h.Class("plate-img"), media),
			h.Div(h.Class("plate-body"),
				swiss.Kicker(tr(lang, "Tafel ", "Plate ")+num, g.Text(fmt.Sprintf("%d / %d", i+1, len(featured)))),
				h.H3(h.Class("plate-title"), g.Text(friendlyName(pr.Name))),
				g.If(corePicks[pr.ID], corePick(lang)),
				g.If(pr.Stars > 0, swiss.Micro(g.Text(fmt.Sprintf("★ %d", pr.Stars)))),
				h.P(g.Text(truncate(pr.DescFor(lang), 350))),
				tagPills(projectTags(pr, 6)),
				projectLinks(pr, lang),
			),
			h.Div(h.Class("plate-pager"), tabPager(lang, i, len(featured), func(j int) string { return friendlyName(featured[j].Name) }, tr(lang, "Tafel", "Plate"))),
		)
	}
	if len(featured) > 0 {
		nodes = append(nodes, h.Div(h.Class("plates"), g.Attr("data-tabs", ""), g.Attr("data-noscatter", ""),
			h.Div(h.Class("plate-tabs"), g.Attr("role", "group"), g.Attr("aria-label", tr(lang, "Tafeln", "Plates")),
				tabLead(tr(lang, "Tafel wählen", "Pick a plate"), len(featured)), g.Group(tabs)),
			g.Group(plates),
		))
	}

	// Live tools: small cells that open the running tool.
	if len(tools) > 0 {
		nodes = append(nodes, swiss.Cell(swiss.CellProps{Variant: "ink", Attrs: []g.Node{h.Class("cat-head")}},
			swiss.Kicker(tr(lang, "Werkzeuge · im Browser nutzbar", "Tools · usable in the browser"), g.Text("● Live")),
		))
		for i, t := range tools {
			nodes = append(nodes, swiss.Cell(swiss.CellProps{Span: 9 / min(len(tools), 3)},
				swiss.Kicker(fmt.Sprintf("%s. %s", string(rune('A'+i)), tr(lang, "Werkzeug", "Tool")), g.Text("● Live")),
				h.H3(h.Class("cell-title"), g.Text(friendlyName(t.Name))),
				h.P(h.Class("cell-text"), g.Text(truncate(t.DescFor(lang), 160))),
				h.Div(h.Class("cat-links cell-foot"),
					swiss.BtnLink(t.LiveURL(), true, []g.Node{g.Attr("target", "_blank"), g.Attr("rel", "noopener")}, g.Text(tr(lang, "Öffnen ↗", "Open ↗"))),
					g.If(t.URL != "", swiss.BtnLink(t.URL, false, []g.Node{g.Attr("target", "_blank"), g.Attr("rel", "noopener noreferrer")}, g.Text(i18n.T(lang, "projects.source")))),
				),
			))
		}
	}

	// Rest: the first three as picture cards, the others as a folding index.
	if len(rest) > 0 {
		nodes = append(nodes, swiss.Cell(swiss.CellProps{Variant: "ink", Attrs: []g.Node{h.Class("cat-head")}},
			swiss.Kicker(tr(lang, "Kleintafeln & Verzeichnis", "Small plates & index"), g.Text(fmt.Sprintf("%d %s", len(rest), tr(lang, "weitere Werke", "more works")))),
		))
		cards := min(3, len(rest))
		for i, pr := range rest[:cards] {
			media := gal(pr, "16/10")
			if media == nil {
				media = monogram(pr, i)
			}
			nodes = append(nodes, swiss.Cell(swiss.CellProps{Span: 3, Variant: "flush", Attrs: []g.Node{h.Class("mini")}},
				media,
				h.Div(h.Class("mini-body"),
					swiss.Kicker(fmt.Sprintf("Nr. %02d", len(featured)+i+1), g.If(pr.Stars > 0, g.Text(fmt.Sprintf("★ %d", pr.Stars)))),
					h.H3(h.Class("cell-title"), g.Text(friendlyName(pr.Name))),
					h.P(h.Class("cell-text"), g.Text(truncate(pr.DescFor(lang), 170))),
					tagPills(projectTags(pr, 4)),
					h.Div(h.Class("cell-foot"), projectLinks(pr, lang)),
				),
			))
		}
		if len(rest) > cards {
			rows := make([]g.Node, 0, len(rest)-cards)
			for i, pr := range rest[cards:] {
				idx := cards + i
				media := gal(pr, "16/10")
				if media == nil {
					media = monogram(pr, idx)
				}
				rows = append(rows, h.Details(h.Class("idx-row"), g.Attr("name", "idx"), g.If(i == 0, g.Attr("open", "")),
					h.Summary(
						h.Span(h.Class("idx-n"), g.Text(fmt.Sprintf("%02d", len(featured)+idx+1))),
						h.Span(h.Class("idx-t"), g.Text(friendlyName(pr.Name))),
						h.Span(h.Class("idx-d"), g.Text(truncate(pr.DescFor(lang), 90))),
						h.Span(h.Class("idx-l"), g.Text(pr.Language)),
						h.Span(h.Class("idx-x"), g.Attr("aria-hidden", "true"), h.Span(h.Class("idx-x-l idx-open"), g.Text(tr(lang, "Öffnen", "Open"))), h.Span(h.Class("idx-x-l idx-close"), g.Text(tr(lang, "Schließen", "Close"))), h.B()),
					),
					h.Div(h.Class("idx-body"),
						h.Div(h.Class("idx-media"), media),
						h.Div(h.Class("idx-text"),
							g.If(pr.Stars > 0, swiss.Micro(g.Text(fmt.Sprintf("★ %d", pr.Stars)))),
							h.P(h.Class("cell-text"), g.Text(pr.DescFor(lang))),
							tagPills(projectTags(pr, 6)),
							projectLinks(pr, lang),
						),
					),
				))
			}
			nodes = append(nodes, swiss.Cell(swiss.CellProps{Variant: "flush", Attrs: []g.Node{h.Class("idx"), g.Attr("data-noscatter", "")}},
				h.Div(h.Class("idx-head"),
					h.Span(g.Text("Nr.")), h.Span(g.Text(tr(lang, "Werk", "Work"))), h.Span(h.Class("idx-d"), g.Text(tr(lang, "Kurzbeschreibung", "Summary"))), h.Span(h.Class("idx-l"), g.Text(tr(lang, "Sprache", "Language"))), h.Span(h.Class("idx-hint"), g.Text(tr(lang, "Zeile anklicken ▾", "Click a row ▾"))),
				),
				g.Group(rows),
			))
		}
	}
	return swiss.Section(p, nodes...)
}

// monogram stands in for a project without screenshots: a solid colour field
// with the project's initials, cycling through the three accent colours.
func monogram(pr hpdata.Project, i int) g.Node {
	fields := []string{"block", "second", "third"}
	words := strings.FieldsFunc(friendlyName(pr.Name), func(r rune) bool { return r == ' ' || r == '-' || r == '_' })
	mono := ""
	for _, w := range words {
		if len(mono) < 2 && w != "" {
			mono += strings.ToUpper(string([]rune(w)[0]))
		}
	}
	return h.Div(h.Class("mono-field mono-"+fields[i%len(fields)]), g.Attr("aria-hidden", "true"),
		h.Span(g.Text(mono)),
		swiss.Micro(g.Text(pr.Language)),
	)
}

// ── 04 Labor ────────────────────────────────────────────────────────────────

func labSection(p swiss.SectionProps, lang string) g.Node {
	p.Title = "Kinetisches Labor"
	p.Sub = "Kinetic Lab"
	p.Note = i18n.T(lang, "sections.playground.sub")
	n := len(playgroundItems)
	tabs := make([]g.Node, n)
	pages := make([]g.Node, n)
	for i, it := range playgroundItems {
		k := strconv.Itoa(i + 1)
		tabs[i] = h.Button(h.Type("button"), h.Class("lab-tab"), g.Attr("data-tab", k), g.Attr("aria-pressed", strconv.FormatBool(i == 0)),
			icon.Icon(it.Icon, icon.Props{Size: "1.1em"}), h.Span(g.Text(it.Label)),
		)
		// Only the visible game runs: labJS gives the first its src once the
		// lab nears the screen, the others get it when their tab opens (and
		// lose it again when it closes).
		frame := h.IFrame(h.Class("lab-frame"), h.Title(it.Label), g.Attr("data-src", it.Href))
		pages[i] = h.Div(h.Class("lab-page"), g.Attr("data-tab-page", k), g.If(i > 0, g.Attr("hidden", "")),
			h.Div(h.Class("lab-stage"), frame),
			h.Div(h.Class("lab-info"),
				swiss.Kicker(fmt.Sprintf("Nr. %02d / %02d", i+1, n), icon.Icon(it.Icon, icon.Props{Size: "1.1em"})),
				h.H3(h.Class("cell-title"), g.Text(it.Label)),
				h.P(h.Class("cell-text"), g.Text(i18n.T(lang, it.DescKey))),
				h.Div(h.Class("cell-foot"), swiss.BtnLink(it.Href, false, nil, g.Text(tr(lang, "Vollbild & Erklärung ↗", "Full page & explainer ↗")))),
			),
			h.Div(h.Class("lab-pager"), tabPager(lang, i, n, func(j int) string { return playgroundItems[j].Label }, tr(lang, "Versuch", "Experiment"))),
		)
	}
	return swiss.Section(p,
		h.Div(h.Class("lab"), g.Attr("data-tabs", ""), g.Attr("data-noscatter", ""),
			h.Div(h.Class("lab-tabs"), g.Attr("role", "group"), g.Attr("aria-label", tr(lang, "Versuche", "Experiments")),
				tabLead(tr(lang, "Versuch wählen", "Pick an experiment"), n), g.Group(tabs)),
			g.Group(pages),
			h.Script(g.Raw(labJS)),
		),
	)
}

// labJS sizes each game iframe to its content, so the whole game (canvas
// and controls) is visible without scrolling inside the frame.
const labJS = `(function(){
  function fit(f){
    var d; try { d = f.contentDocument; } catch(e) { return; }
    if (!d || !d.body) return;
    var m = d.querySelector('main') || d.body;
    var size = function(){ f.style.height = Math.ceil(m.getBoundingClientRect().height) + 'px'; };
    size();
    if (window.ResizeObserver) new ResizeObserver(size).observe(m);
  }
  document.querySelectorAll('.lab-frame').forEach(function(f){
    f.addEventListener('load', function(){ if (f.getAttribute('src')) fit(f); });
  });
  // start the visible game only when the lab is about to scroll into view
  var lab = document.querySelector('.lab');
  if (!lab) return;
  var start = function(){
    var f = lab.querySelector('[data-tab-page]:not([hidden]) iframe[data-src]');
    if (f && !f.getAttribute('src')) f.setAttribute('src', f.getAttribute('data-src'));
  };
  if (!window.IntersectionObserver) return start();
  var io = new IntersectionObserver(function(es){
    if (es.some(function(e){ return e.isIntersecting; })) { io.disconnect(); start(); }
  }, { rootMargin: '300px 0px' });
  io.observe(lab);
})();`

const labCSS = `
.lab { grid-column: 1 / -1; }
.lab-tabs { display: grid; grid-template-columns: repeat(6, minmax(0, 1fr)); box-shadow: 0 1px 0 var(--ink); }
.lab-tabs .tab-lead { grid-row: span 2; }
.lab-tab { all: unset; box-sizing: border-box; cursor: pointer; display: flex; flex-direction: column; gap: .45rem; padding: .7rem .8rem; box-shadow: 1px 0 0 var(--ink), 0 1px 0 var(--ink); font-size: .8rem; font-weight: 700; line-height: 1.15; }
.lab-tab:hover { background: var(--act); color: var(--act-fg); }
.lab-tab[aria-pressed="true"] { background: var(--third); color: var(--third-fg); }
.lab-tab:focus-visible { outline: 3px solid var(--accent); outline-offset: -3px; }
.lab-page { display: grid; grid-template-columns: minmax(0, 1fr); }
.lab-page[hidden] { display: none; }
.lab-stage { background: var(--surface-2); }
.lab-frame:focus { outline: 3px solid var(--act); outline-offset: -3px; }
.lab-frame { display: block; width: 100%; height: 640px; border: 0; background: var(--bg); }
/* info strip above the game, so the game gets the full width */
.lab-info { order: -1; padding: .9rem var(--swiss-pad); display: grid; grid-template-columns: minmax(12rem, 16rem) minmax(0, 1fr) auto; column-gap: 1.5rem; row-gap: .3rem; align-items: center; box-shadow: 0 1px 0 var(--ink); }
.lab-info > :nth-child(1) { grid-column: 1; grid-row: 1; }
.lab-info > :nth-child(2) { grid-column: 1; grid-row: 2; margin: 0; }
.lab-info > :nth-child(3) { grid-column: 2; grid-row: 1 / span 2; margin: 0; }
.lab-info > :nth-child(4) { grid-column: 3; grid-row: 1 / span 2; margin: 0; }
.lab-pager { grid-column: 1 / -1; }
@media (max-width: 1024px) {
  .lab-tabs { grid-template-columns: repeat(2, minmax(0, 1fr)); }
  .lab-tabs .tab-lead { grid-row: auto; grid-column: 1 / -1; }
  .lab-page { grid-template-columns: minmax(0, 1fr); }
  .lab-info { grid-template-columns: minmax(0, 1fr); }
  .lab-info > * { grid-column: 1 !important; grid-row: auto !important; }
  .lab-frame { height: 70vh; }
}
`

// ── Laufbahn / Schriften ─────────────────────────────────────────────

var orgShort = map[string]string{
	"Johannes Kepler University Linz":  "JKU Linz",
	"Johannes Kepler Universität Linz": "JKU Linz",
	"Land Oberösterreich":              "Land OÖ",
	"HerzReha Bad Ischl":               "HerzReha",
}

// thesisEdu maps a thesis type to the timeline id of its degree.
// degree describes how a thesis type is shown. Degree theses get a colour
// poster; everything else (project reports, seminar papers) is a "Schrift":
// a slim paper cell with less visual weight.
type degree struct {
	field, mark string
	sub         [2]string // DE, EN
	org         string    // timeline id of the programme (logo, organisation, title)
	when        string    // timeline id whose period dates the work itself
	order       int
}

var degrees = map[string]degree{
	"HTL":      {"second", "HTL", [2]string{"Reife- & Diplomprüfung", "Matura & diploma"}, "htl-steyr", "", 1},
	"Bachelor": {"block", "BSc", [2]string{"Bachelor of Science", "Bachelor of Science"}, "jku-bachelor", "bachelor-thesis", 2},
	"Master":   {"third", "DI", [2]string{"Diplom-Ingenieur", "Diplom-Ingenieur · Master's degree"}, "jku-master", "dynatrace-thesis", 4},
}

// paperOrder places papers between degrees: a master project sits after
// the bachelor and before the master thesis.
var paperOrder = map[string]int{"Project": 3, "Projekt": 3, "Seminar": 3}

func papersSection(p swiss.SectionProps, d hpdata.SiteData, lang string) g.Node {
	thesis := append([]hpdata.Thesis(nil), d.ThesisFor(lang)...)
	if len(thesis) == 0 {
		return nil
	}
	p.Title = "Abschlüsse & Schriften"
	p.Sub = "Degrees & Papers"
	p.Wide = true
	p.Note = i18n.T(lang, "sections.papers.sub")
	tl := map[string]hpdata.TimelineItem{}
	for _, t := range d.Timeline {
		tl[t.ID] = t
	}
	rank := func(t hpdata.Thesis) int {
		if dg, ok := degrees[t.Type]; ok {
			return dg.order
		}
		if o, ok := paperOrder[t.Type]; ok {
			return o
		}
		return 3
	}
	sort.SliceStable(thesis, func(i, j int) bool { return rank(thesis[i]) < rank(thesis[j]) })

	pdf := func(t hpdata.Thesis, primary bool) g.Node {
		if t.PDF == "" {
			return swiss.Micro(h.Style("color:var(--muted)"), g.Text("◷ "+i18n.T(lang, "sections.experience.thesis_coming_soon")))
		}
		return swiss.BtnLink(t.PDF, primary, []g.Node{g.Attr("target", "_blank"), g.Attr("rel", "noopener")},
			g.Text("PDF ↓ "+i18n.T(lang, "sections.experience.thesis_view_pdf")))
	}
	// Degree posters take 3 columns (the highest degree 4); the papers share
	// one 2-column stack, each a horizontal sheet, between bachelor and master.
	var cells []g.Node
	var stack []g.Node
	stackAt := -1
	for i, t := range thesis {
		letter := string(rune('A' + i))
		dg, isDegree := degrees[t.Type]
		if !isDegree {
			if stackAt < 0 {
				stackAt = len(cells)
				cells = append(cells, nil)
			}
			meta := paperPeriod(t, lang)
			if len(t.CoAuthors) > 0 {
				meta = append(meta, tr(lang, "mit ", "with ")+strings.Join(t.CoAuthors, ", "))
			}
			stack = append(stack, h.Div(h.Class("paper"),
				swiss.Kicker(letter+". "+tr(lang, "Schrift", "Paper"), g.Text(t.Type)),
				h.Span(h.Class("paper-mark"), g.Text(paperMark(t.Type))),
				h.H3(h.Class("paper-title"), g.Text(t.Title)),
				g.If(len(meta) > 0, swiss.Micro(h.Class("paper-meta"), g.Text(strings.Join(meta, " · ")))),
				h.P(h.Class("cell-text"), g.Text(truncate(t.Description, 150))),
				h.Div(h.Class("cell-foot"), pdf(t, false)),
			))
			continue
		}
		e, hasOrg := tl[dg.org]
		period := ""
		if w, ok := tl[dg.when]; ok {
			period = w.FormatPeriod()
		} else if hasOrg && e.End != nil && len(*e.End) >= 4 {
			period = tr(lang, "Abschluss ", "Completed ") + (*e.End)[:4]
		}
		sub := dg.sub[0]
		if lang != "de" {
			sub = dg.sub[1]
		}
		span := 3
		if dg.order == 4 {
			span = 4
		}
		logo := thesisLogo(t.Type)
		cells = append(cells, swiss.Cell(swiss.CellProps{Span: span, Variant: "flush", Attrs: []g.Node{h.Class("poster")}},
			h.Div(h.Class("poster-field mono-"+dg.field),
				h.Div(h.Class("poster-top"),
					swiss.Micro(g.Text(letter+". "+tr(lang, "Abschlussarbeit", "Thesis"))),
					g.If(period != "", swiss.Micro(g.Text(period))),
				),
				h.Div(
					h.Span(h.Class("poster-mark"), g.Text(dg.mark)),
					h.Span(h.Class("poster-sub"), g.Text(sub)),
				),
				h.Div(h.Class("poster-org"),
					g.If(logo != "", h.Img(h.Src(logo), h.Alt(""), g.Attr("loading", "lazy"))),
					g.If(hasOrg, swiss.Micro(g.Text(e.Organization))),
				),
			),
			h.Div(h.Class("poster-body"),
				g.If(hasOrg, swiss.Kicker(e.TitleFor(lang), nil)),
				h.H3(h.Class("cell-title"), g.Text(t.Title)),
				g.If(t.Type == "Master", corePick(lang)),
				h.P(h.Class("cell-text"), g.Text(truncate(t.Description, 360))),
				h.Div(h.Class("cell-foot"), pdf(t, true)),
			),
		))
	}
	if stackAt >= 0 {
		cells[stackAt] = swiss.Cell(swiss.CellProps{Span: 2, Variant: "flush", Attrs: []g.Node{h.Class("paper-stack")}}, g.Group(stack))
	}
	return swiss.Section(p, cells...)
}

// paperPeriod formats when a paper was written: "Aug 2025 – Mai 2026 ·
// 9 Monate" with a known duration, else just the submission month.
func paperPeriod(t hpdata.Thesis, lang string) []string {
	end, err := time.Parse("2006-01", t.Submitted)
	if err != nil {
		return nil
	}
	month := func(m time.Time) string {
		if lang == "de" {
			return []string{"Jän", "Feb", "Mär", "Apr", "Mai", "Jun", "Jul", "Aug", "Sep", "Okt", "Nov", "Dez"}[m.Month()-1] + " " + m.Format("2006")
		}
		return m.Format("Jan 2006")
	}
	if t.Months <= 1 {
		return []string{tr(lang, "Abgabe ", "Submitted ") + month(end)}
	}
	start := end.AddDate(0, -(t.Months - 1), 0)
	return []string{month(start) + " – " + month(end), fmt.Sprintf("%d %s", t.Months, tr(lang, "Monate", "months"))}
}

func paperMark(kind string) string {
	switch kind {
	case "Project", "Projekt":
		return "PRJ"
	case "Seminar":
		return "SE"
	}
	return strings.ToUpper(kind)
}

// ── Kompetenzen ─────────────────────────────────────────────────────────────

var gateShort = map[string]string{
	"Languages": "SPR", "Web": "WEB", "Infra / Homelab": "INF",
	"Security": "SEC", "Embedded / BLE": "EMB", "Ops / Data": "OPS",
}

func skillsSection(p swiss.SectionProps, d hpdata.SiteData, lang string) g.Node {
	var rows []swiss.FlapRow
	groups := hpdata.SkillGroups()
	for _, sg := range groups {
		for _, s := range sg.Skills {
			v := "—"
			if n := d.RepoCount(s); n > 0 {
				v = strconv.Itoa(n)
			}
			rows = append(rows, swiss.FlapRow{Icon: hpdata.TechIcon(s), Name: s, Gate: gateShort[sg.Label], Value: v, Group: sg.Label})
		}
	}
	p.Note = fmt.Sprintf(tr(lang,
		"Fallblatt-Anzeige. %d Kompetenzen, %d Gleise. Repos = öffentliche GitHub-Projekte mit dieser Technik; Berufserfahrung steht in der Laufbahn.",
		"Split-flap display. %d skills, %d groups. Repos = public GitHub projects using it; work experience is in the career section."), len(rows), len(groups))
	p.Attrs = append(p.Attrs, h.Class("sec-inv"))
	fb := swiss.FlapBoardProps{ID: "skills-board", Rows: rows, ColValue: "Repos"}
	if lang != "de" {
		fb.PageLabel, fb.Pause, fb.Resume, fb.Prev, fb.Next = "Page", "Pause", "Resume", "Previous page", "Next page"
		fb.Caption, fb.ColGroup = "Skills by group", "Group"
	}
	counts := make([]int, len(groups))
	most := 1
	for i, sg := range groups {
		counts[i] = d.RepoCount(sg.Skills...)
		most = max(most, counts[i])
	}
	legend := make([]g.Node, len(groups))
	for i, sg := range groups {
		legend[i] = h.Li(
			h.Span(h.Class("sk-gate"), g.Text(gateShort[sg.Label])),
			h.Span(h.Class("sk-name"), g.Text(sg.Label), h.Small(g.Text(fmt.Sprintf(" · %d", len(sg.Skills))))),
			h.Span(h.Class("sk-bar"), h.I(h.Style(fmt.Sprintf("width:%d%%", counts[i]*100/most)))),
			h.Span(h.Class("sk-val"), g.Text(strconv.Itoa(counts[i]))),
		)
	}
	return swiss.Section(p,
		swiss.Cell(swiss.CellProps{},
			swiss.Kicker(tr(lang, "Gleise", "Groups"), g.Text(tr(lang, "Kompetenzen · Repos", "skills · repos"))),
			h.Ul(h.Class("sk-legend"), g.Group(legend)),
		),
		swiss.Cell(swiss.CellProps{}, swiss.FlapBoard(fb)),
	)
}

const skillsCSS = `
.sk-legend { list-style: none; margin: 0; padding: 0; display: grid; grid-template-columns: repeat(auto-fit, minmax(min(100%, 26rem), 1fr)); gap: .55rem 2.5rem; }
.sk-legend li { display: grid; grid-template-columns: 3.2rem minmax(0, 12rem) minmax(0, 1fr) 2.5rem; gap: .8rem; align-items: center; }
.sk-gate { font-family: var(--font-mono); font-weight: 800; font-size: .75rem; color: var(--tile-hi); }
.sk-name { font-weight: 700; font-size: .9rem; }
.sk-name small { font-family: var(--font-mono); font-weight: 400; opacity: .7; }
.sk-bar { height: 10px; background: var(--tile); }
.sk-bar i { display: block; height: 100%; background: var(--tile-fg); }
.sk-val { font-family: var(--font-mono); font-size: .8rem; text-align: right; }
/* the board is dark in both modes: give the radar grid and kickers board colours */
.sec-inv > [data-slot="body"] { --line: #5a5f70; --muted: #a6a8b0; }
.sec-inv [data-component="swiss-kicker"] [data-slot="aside"] { color: var(--tile-hi); }
@media (max-width: 640px) { .sk-legend li { grid-template-columns: 2.8rem minmax(0, 1fr) 2.2rem; } .sk-bar { display: none; } }
`

const catalogueCSS = `
.plates { grid-column: 1 / -1; }
.plate[hidden] { display: none; }
.plate-tabs { display: grid; grid-template-columns: repeat(auto-fit, minmax(130px, 1fr)); box-shadow: 0 1px 0 var(--ink); }
.plate-pager { grid-column: 1 / -1; }
.plate-tab::after { content: "▸"; position: absolute; top: .6rem; right: .7rem; width: 1.2rem; height: 1.2rem; display: grid; place-items: center; font-size: .7rem; background: var(--act); color: var(--act-fg); }
.plate-tab[aria-pressed="true"]::after { content: "●"; color: inherit; background: none; }
.plate-tab:not([aria-pressed="true"]) { box-shadow: 1px 0 0 var(--ink), inset 0 -4px 0 var(--rule-subtle); }
.plate-tab:not([aria-pressed="true"]):hover { box-shadow: 1px 0 0 var(--ink), inset 0 -4px 0 var(--act-fg); }
.plate-tab { all: unset; position: relative; box-sizing: border-box; cursor: pointer; display: flex; flex-direction: column; gap: .35rem; padding: .8rem var(--swiss-pad); box-shadow: 1px 0 0 var(--ink); min-width: 0; }
.core-pick { align-self: flex-start; display: inline-block; font-family: var(--font-mono); font-size: .7rem; font-weight: 800; text-transform: uppercase; letter-spacing: .06em; padding: .2rem .45rem; background: var(--block); color: var(--block-fg); }
.plate-tab-core { font-family: var(--font-mono); font-size: .62rem; font-weight: 800; text-transform: uppercase; letter-spacing: .06em; color: var(--main-text); }
.plate-tab[aria-pressed="true"] .plate-tab-core { color: inherit; }
.plate-tab-n { font-size: 1.9rem; font-weight: 900; letter-spacing: -0.05em; line-height: .85; color: var(--main-text); }
.plate-tab-t { font-family: var(--font-mono); font-size: .68rem; font-weight: 700; text-transform: uppercase; letter-spacing: .05em; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.plate-tab:hover { background: var(--act); color: var(--act-fg); }
.plate-tab:not([aria-pressed="true"]):hover .plate-tab-n { color: inherit; }
.plate-tab[aria-pressed="true"] { background: var(--block); color: var(--block-fg); }
.plate-tab[aria-pressed="true"] .plate-tab-n { color: inherit; }
.plate-tab:focus-visible { outline: 3px solid var(--accent); outline-offset: -3px; }
.idx-head, .idx-row > summary { display: grid; grid-template-columns: 3.5rem minmax(0, 14rem) minmax(0, 1fr) 7rem 7.5rem; gap: 1rem; align-items: baseline; padding: .5rem var(--swiss-pad); }
.idx-head .idx-d { font-size: inherit; color: inherit; }
.idx-head { font-family: var(--font-mono); font-size: .64rem; font-weight: 700; text-transform: uppercase; letter-spacing: .08em; color: var(--muted); border-bottom: 1px solid var(--ink); padding-block: .5rem; }
.idx-row { border-bottom: 1px solid var(--rule-subtle); }
.idx-row:last-child { border-bottom: 0; }
.idx-row > summary { list-style: none; cursor: pointer; }
.idx-row > summary::-webkit-details-marker { display: none; }
.idx-row > summary:hover { background: var(--act); color: var(--act-fg); }
.idx-row > summary:hover * { color: inherit; }
.idx-row > summary:focus-visible { outline: 3px solid var(--accent); outline-offset: -3px; }
.idx-n { font-weight: 900; font-size: 1.25rem; letter-spacing: -0.04em; color: var(--main-text); }
.idx-t { font-weight: 800; letter-spacing: -0.01em; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.idx-d { font-size: .85rem; color: var(--muted); overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.idx-l { font-family: var(--font-mono); font-size: .68rem; text-transform: uppercase; letter-spacing: .05em; }
.idx-x { display: flex; justify-content: flex-end; align-items: center; gap: .5rem; }
.idx-x-l { font-family: var(--font-mono); font-size: .66rem; font-weight: 700; text-transform: uppercase; letter-spacing: .08em; color: var(--main-text); }
.idx-x b { width: 1.5rem; height: 1.5rem; display: grid; place-items: center; border: 1px solid var(--act-fg); background: var(--act); color: var(--act-fg); font-size: 1rem; line-height: 1; }
.idx-x b::before { content: "+"; }
.idx-row[open] .idx-x b::before { content: "–"; }
.idx-hint { text-align: right; }
.idx-row[open] .idx-x-l { color: inherit; }
.idx-close, .idx-row[open] .idx-open { display: none; }
.idx-row[open] .idx-close { display: inline; }
.idx-row[open] > summary { background: var(--block); color: var(--block-fg); }
.idx-row[open] > summary .idx-n, .idx-row[open] > summary .idx-d { color: inherit; }
.idx-body { display: grid; grid-template-columns: minmax(0, 4fr) minmax(0, 5fr); border-top: 1px solid var(--ink); }
.idx-media { box-shadow: 1px 0 0 var(--ink); min-width: 0; }
.idx-text { padding: var(--swiss-pad); display: flex; flex-direction: column; gap: .75rem; }
@media (max-width: 1024px) {
  .plate-tabs { grid-template-columns: repeat(3, minmax(0, 1fr)); }
  .idx-head { display: none; }
  .idx-row > summary { grid-template-columns: 2.5rem minmax(0, 1fr) auto; padding-inline: 0; }
  .idx-row > summary .idx-d, .idx-row > summary .idx-l { display: none; }
  .idx-body { grid-template-columns: minmax(0, 1fr); }
  .idx-media { box-shadow: 0 1px 0 var(--ink); }
  .idx-text { padding-inline: 0; }
}
`
