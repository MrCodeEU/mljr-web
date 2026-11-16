package pages

import (
	"fmt"
	"math"
	"slices"
	"strings"
	"time"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	hpdata "mljr-web/projects/homepage/data"
	"mljr-web/ui/swiss"
)

// The skill network ("Netzplan"): each line is a field of work, each stop a
// job, a degree or a project. Where fields meet, the stop is an interchange
// drawn as a capsule across both tracks. Positions are hand-laid on a
// column/row grid so every track runs straight and labels never collide.

type netLine struct {
	key, de, en, color string
}

var netLines = []netLine{
	{"h", "Hardware & Embedded", "Hardware & embedded", "var(--nz-h)"},
	{"g", "Software & Go", "Software & Go", "var(--nz-g)"},
	{"s", "Security", "Security", "var(--nz-s)"},
	{"w", "Web & UI", "Web & UI", "var(--nz-w)"},
}

type netStop struct {
	key   string
	label [2]string // two short label lines
	col   int
	row   int    // top row of the stop
	row2  int    // bottom row for interchanges (== row otherwise)
	lines string // line keys, e.g. "gs"
	tl    string // timeline id, or
	proj  string // project id
}

var netStops = []netStop{
	{"htl", [2]string{"HTL", "Steyr"}, 0, 0, 0, "h", "htl-steyr", ""},
	{"engel", [2]string{"ENGEL", ""}, 1, 0, 0, "h", "engel-intern-2017", ""},
	{"rexroth", [2]string{"Bosch", "Rexroth"}, 2, 0, 0, "h", "bosch-rexroth-intern", ""},
	{"tutor", [2]string{"JKU", "Tutor"}, 4, 0, 0, "h", "jku-tutor", ""},
	{"bsc", [2]string{"JKU", "Informatik"}, 3, 1, 1, "g", "jku-bachelor", ""},
	{"bthesis", [2]string{"Bachelor-", "arbeit"}, 4, 1, 1, "g", "bachelor-thesis", ""},
	{"prj", [2]string{"Master-", "projekt"}, 5, 1, 1, "hg", "", "jku-master-project-network-devices-power-consumption-analysis"},
	{"ntray", [2]string{"Nightscout", "Tray"}, 6, 1, 1, "g", "", "nightscout-tray"},
	{"godrive", [2]string{"goDrive", ""}, 7, 1, 1, "g", "", "godrive"},
	{"dyni", [2]string{"Dynatrace", "Praktikum"}, 8, 1, 2, "gs", "dynatrace-intern", ""},
	{"dynt", [2]string{"Dynatrace", "Thesis"}, 9, 1, 2, "gs", "dynatrace-thesis", ""},
	{"homelab", [2]string{"Homelab", "Automation"}, 10, 1, 2, "gs", "", "homelab-automation"},
	{"mljr", [2]string{"mljr.eu", ""}, 11, 0, 1, "gw", "", "mljr-web"},
	{"glucava", [2]string{"Glucava", ""}, 12, 1, 1, "g", "", "glucava"},
	{"msc", [2]string{"JKU", "Security MSc"}, 6, 2, 2, "s", "jku-master", ""},
	{"faceauth", [2]string{"Face Auth", ""}, 7, 2, 2, "s", "", "face-auth"},
	{"landooe", [2]string{"Land OÖ", ""}, 12, 2, 2, "s", "land-ooe-security", ""},
	{"bosch", [2]string{"Bosch", ""}, 6, 0, 0, "w", "bosch-intern-2023", ""},
	{"ventopay", [2]string{"ventopay", ""}, 7, 0, 0, "w", "ventopay-intern", ""},
	{"wordwiz", [2]string{"WordWiz", ""}, 8, 0, 0, "w", "", "WordWiz"},
	{"sudoku", [2]string{"Sudoku", "Gen"}, 9, 0, 0, "w", "", "sudoku-gen"},
}

// netRoutes: track waypoints per line as (col,row), in order.
var netRoutes = map[string][][2]int{
	"h": {{0, 0}, {4, 0}, {5, 1}},
	"g": {{3, 1}, {12, 1}},
	"s": {{6, 2}, {12, 2}},
	"w": {{6, 0}, {11, 0}},
}

// yearSpan: "2021", "2021–24" or "2026–" from a timeline entry.
func yearSpan(t hpdata.TimelineItem) string {
	if len(t.Start) < 4 {
		return ""
	}
	y := t.Start[:4]
	switch {
	case t.End == nil:
		return y + "–"
	case len(*t.End) >= 4 && (*t.End)[:4] != y:
		return y + "–" + (*t.End)[2:4]
	}
	return y
}

// stopInfo is what a panel shows for one stop or timeline entry.
type stopInfo struct {
	key, lines                                 string
	title, org, period, place, kind, desc, url string
	logo                                       string
	images                                     []string
	details, tags                              []string
}

func timelineInfo(key, lines string, t hpdata.TimelineItem, lang string) stopInfo {
	return stopInfo{
		key: key, lines: lines,
		title: t.TitleFor(lang), org: t.Organization, period: t.FormatPeriod() + " · " + t.FormatDuration(), place: t.Location,
		kind: map[string]string{"work": "Job", "education": tr(lang, "Studium", "Study"), "thesis": tr(lang, "Abschlussarbeit", "Thesis")}[t.Kind],
		desc: t.SummaryFor(lang), details: t.DetailsFor(lang), tags: t.Tags, logo: t.Logo,
	}
}

// careerSection ("Laufbahn"): the skill network on top, a year strip with
// every timeline entry below it, and one detail panel shared by both.
// Selecting a stop or a bar shows its panel (logo or screenshots, period,
// place, summary, details).
func careerSection(p swiss.SectionProps, d hpdata.SiteData, lang string) g.Node {
	if len(d.Timeline) == 0 {
		return nil
	}
	p.NoScatter = true
	p.Note = tr(lang,
		"Netzplan: jede Linie ein Fachgebiet, jede Station ein Job, Studium oder Projekt. Umstiege zeigen, wo Gebiete zusammenkommen. Darunter die Zeitleiste.",
		"Network map: each line a field, each stop a job, degree or project. Interchanges show where fields meet. The timeline sits below.")

	tl := map[string]hpdata.TimelineItem{}
	for _, t := range d.Timeline {
		tl[t.ID] = t
	}
	projects := map[string]hpdata.Project{}
	for _, pr := range d.GitHub {
		projects[pr.ID] = pr
	}

	// panel data: every network stop, then timeline entries without a stop
	var infos []stopInfo
	stopOf := map[string]string{} // timeline id -> panel key
	lineOf := map[string]string{} // timeline id -> line key for the strip
	for _, st := range netStops {
		if st.tl != "" {
			t, ok := tl[st.tl]
			if !ok {
				continue
			}
			in := timelineInfo(st.key, st.lines, t, lang)
			infos = append(infos, in)
			stopOf[st.tl] = st.key
			lineOf[st.tl] = st.lines[:1]
			continue
		}
		pr, ok := projects[st.proj]
		if !ok {
			continue
		}
		in := stopInfo{key: st.key, lines: st.lines, title: friendlyName(pr.Name), kind: tr(lang, "Projekt", "Project"),
			desc: pr.DescFor(lang), tags: projectTags(pr, 6), images: pr.LocalImages()}
		if HasProjectDetail(pr) {
			in.url = "/projects/" + pr.ID
		}
		infos = append(infos, in)
	}
	// ENGEL's second internship shares the stop and its panel; anything else
	// gets its own panel
	if t2, ok := tl["engel-intern-2018"]; ok {
		stopOf["engel-intern-2018"], lineOf["engel-intern-2018"] = "engel", "h"
		for i := range infos {
			if infos[i].key == "engel" {
				infos[i].period += " / " + t2.FormatPeriod() + " · " + t2.FormatDuration()
				if t2.Location != "" && t2.Location != infos[i].place {
					infos[i].place += " / " + t2.Location
				}
				ds := slices.Clone(infos[i].details)
				for _, dt := range t2.DetailsFor(lang) {
					if !slices.Contains(ds, dt) {
						ds = append(ds, dt)
					}
				}
				infos[i].details = ds
			}
		}
	}
	for _, t := range d.Timeline {
		if _, ok := stopOf[t.ID]; ok {
			continue
		}
		infos = append(infos, timelineInfo(t.ID, "", t, lang))
		stopOf[t.ID], lineOf[t.ID] = t.ID, "o"
	}
	first := ""
	if len(infos) > 0 {
		first = infos[0].key
	}

	return swiss.Section(p,
		swiss.Cell(swiss.CellProps{Variant: "flush", Attrs: []g.Node{h.ID("nz"), g.Attr("data-tabs", ""), g.Attr("data-noscatter", "")}},
			netLegend(lang),
			h.Div(h.Class("nz-wrap"), g.Raw(netSVG(tl, projects, first, lang))),
			swiss.Micro(h.Class("nz-swipe"), g.Attr("aria-hidden", "true"), g.Text(tr(lang, "← seitlich wischen →", "← swipe sideways →"))),
			yearStrip(d.Timeline, stopOf, lineOf, first, lang),
			g.Group(netPanels(infos, first, lang)),
			careerMap(d, stopOf, lang),
			swiss.GalleryScript(),
			h.Script(g.Raw(netJS)),
		),
	)
}

func netLegend(lang string) g.Node {
	legend := make([]g.Node, len(netLines))
	for i, ln := range netLines {
		legend[i] = h.Button(h.Type("button"), h.Class("nz-line"), g.Attr("data-netline", ln.key), g.Attr("aria-pressed", "false"),
			h.Span(h.Class("nz-line-badge"), g.Attr("data-line", ln.key), h.Style("background:"+ln.color), g.Text(strings.ToUpper(ln.key))),
			h.Span(g.Text(tr(lang, ln.de, ln.en))),
		)
	}
	return h.Div(h.Class("nz-legend"),
		swiss.Micro(h.Class("nz-hint"), g.Text(tr(lang, "Linie wählen zum Hervorheben · Station oder Balken anklicken für Details", "Pick a line to highlight · click a stop or bar for details"))),
		h.Div(h.Class("nz-lines"), g.Group(legend)),
	)
}

// netSVG draws tracks, badges and stops.
func netSVG(tl map[string]hpdata.TimelineItem, projects map[string]hpdata.Project, first, lang string) string {
	const colW, rowH, padX, top = 84.0, 120.0, 70.0, 130.0
	x := func(c int) float64 { return padX + float64(c)*colW }
	y := func(r int) float64 { return top + float64(r)*rowH }
	W, H := padX*2+12*colW+20, top+2*rowH+90

	var b strings.Builder
	fmt.Fprintf(&b, `<svg viewBox="0 0 %.0f %.0f" class="nz" role="group" aria-label="%s">`, W, H,
		tr(lang, "Netzplan der Fachgebiete", "Network map of fields"))
	for _, ln := range netLines {
		route := netRoutes[ln.key]
		var path strings.Builder
		for i, pt := range route {
			px, py := x(pt[0]), y(pt[1])
			if i == 0 {
				fmt.Fprintf(&path, "M%.1f %.1f", px, py)
				continue
			}
			// straight run, then a 45° diagonal into the waypoint
			lx, ly := x(route[i-1][0]), y(route[i-1][1])
			dx, dy := px-lx, py-ly
			dd := math.Min(math.Abs(dx), math.Abs(dy))
			fmt.Fprintf(&path, " L%.1f %.1f L%.1f %.1f", px-math.Copysign(dd, dx), ly, px, py)
		}
		fmt.Fprintf(&b, `<path d="%s" class="nz-track" data-lines="%s" stroke="%s"/>`, path.String(), ln.key, ln.color)
		sx, sy := x(route[0][0])-40, y(route[0][1])
		fmt.Fprintf(&b, `<g class="nz-badge" data-lines="%s"><rect x="%.1f" y="%.1f" width="26" height="26" fill="%s"/><text x="%.1f" y="%.1f">%s</text></g>`,
			ln.key, sx-13, sy-13, ln.color, sx, sy+5, strings.ToUpper(ln.key))
	}
	for _, st := range netStops {
		title, meta := "", ""
		if st.tl != "" {
			t, ok := tl[st.tl]
			if !ok {
				continue
			}
			title, meta = t.TitleFor(lang), yearSpan(t)
			if st.key == "engel" {
				meta = "2017 · 18"
			}
		} else {
			pr, ok := projects[st.proj]
			if !ok {
				continue
			}
			title, meta = friendlyName(pr.Name), strings.ToUpper(tr(lang, "Projekt", "Project"))
		}
		cx, y1, y2 := x(st.col), y(st.row), y(st.row2)
		cls := "nz-stop"
		if len(st.lines) > 1 {
			cls += " nz-x"
		}
		fmt.Fprintf(&b, `<g class="%s" data-lines="%s" data-tab="%s" role="button" tabindex="0" aria-pressed="%v" aria-label="%s">`,
			cls, strings.Join(strings.Split(st.lines, ""), " "), st.key, st.key == first, title)
		labels := []string{st.label[0]}
		if st.label[1] != "" {
			labels = append(labels, st.label[1])
		}
		ly := y2 + 32 // below the lowest track …
		if st.row == 0 {
			ly = y1 - 26 - float64(len(labels))*14 // … or above the top row
		}
		t0 := math.Min(y1-12, ly-14)
		t1 := math.Max(y2+12, ly+float64(len(labels))*14+4)
		fmt.Fprintf(&b, `<rect class="nz-hit" x="%.1f" y="%.1f" width="%.0f" height="%.1f"/>`, cx-colW/2, t0, colW, t1-t0)
		if y2 > y1 {
			fmt.Fprintf(&b, `<rect x="%.1f" y="%.1f" width="24" height="%.1f" rx="12" class="nz-mark"/>`, cx-12, y1-12, y2-y1+24)
		} else {
			fmt.Fprintf(&b, `<rect x="%.1f" y="%.1f" width="20" height="20" class="nz-mark"/>`, cx-10, y1-10)
		}
		fmt.Fprintf(&b, `<text x="%.1f" y="%.1f" class="nz-n">`, cx, ly)
		for i, l := range labels {
			dy := 0
			if i > 0 {
				dy = 14
			}
			fmt.Fprintf(&b, `<tspan x="%.1f" dy="%d">%s</tspan>`, cx, dy, l)
		}
		fmt.Fprintf(&b, `</text><text x="%.1f" y="%.1f" class="nz-y">%s</text></g>`, cx, ly+float64(len(labels)-1)*14+14, meta)
	}
	b.WriteString(`</svg>`)
	return b.String()
}

// yearStrip: one row per line, each timeline entry a clickable bar on a
// shared year axis (projects have no dates and stay on the map only).
func yearStrip(items []hpdata.TimelineItem, stopOf, lineOf map[string]string, first, lang string) g.Node {
	const startYear = 2015
	now := time.Now()
	span := float64((now.Year()-startYear)*12 + int(now.Month()))
	pos := func(ym string) float64 {
		t, err := time.Parse("2006-01", ym)
		if err != nil {
			return 0
		}
		return float64((t.Year()-startYear)*12+int(t.Month())-1) / span * 100
	}
	rows := append(append([]netLine(nil), netLines...), netLine{"o", "Sonstiges", "Other", "var(--muted)"})
	var rowNodes []g.Node
	for _, ln := range rows {
		var bars []g.Node
		for _, t := range items {
			if lineOf[t.ID] != ln.key {
				continue
			}
			l := pos(t.Start)
			r := 100.0
			if t.End != nil {
				r = pos(*t.End) + 100/span
			}
			label := orgShort[t.Organization]
			if label == "" {
				label = t.Organization
			}
			bars = append(bars, h.Button(h.Type("button"), h.Class("ys-bar"), g.Attr("data-tab", stopOf[t.ID]),
				g.Attr("aria-pressed", fmt.Sprint(stopOf[t.ID] == first)), h.Title(t.TitleFor(lang)+" · "+t.FormatPeriod()),
				h.Style(fmt.Sprintf("left:%.2f%%;width:max(10px,%.2f%%);--c:%s", l, r-l, ln.color)),
				h.Span(g.Text(label)),
			))
		}
		if len(bars) == 0 {
			continue
		}
		rowNodes = append(rowNodes, h.Div(h.Class("ys-row"),
			h.Span(h.Class("ys-key"), h.Style("background:"+ln.color), g.Attr("data-line", ln.key), g.Text(strings.ToUpper(ln.key))),
			h.Div(h.Class("ys-track"), g.Group(bars)),
		))
	}
	ticks := make([]g.Node, 0, now.Year()-startYear+1)
	for y := startYear; y <= now.Year(); y++ {
		ticks = append(ticks, h.Span(h.Style(fmt.Sprintf("left:%.2f%%", float64((y-startYear)*12)/span*100)), g.Text(fmt.Sprint(y))))
	}
	return h.Div(h.Class("ys"),
		swiss.Kicker(tr(lang, "Zeitleiste", "Timeline"), g.Text(fmt.Sprintf("%d–%s", startYear, tr(lang, "heute", "today")))),
		h.Div(h.Class("ys-row ys-axis"), h.Span(), h.Div(h.Class("ys-track"), g.Group(ticks))),
		g.Group(rowNodes),
	)
}

func netPanels(infos []stopInfo, first, lang string) []g.Node {
	out := make([]g.Node, len(infos))
	for i, in := range infos {
		chips := make([]g.Node, 0, 2)
		for _, ln := range netLines {
			if strings.Contains(in.lines, ln.key) {
				chips = append(chips, h.Span(h.Class("nz-chip"), g.Attr("data-line", ln.key), h.Style("background:"+ln.color), g.Text(strings.ToUpper(ln.key)+" · "+tr(lang, ln.de, ln.en))))
			}
		}
		var media g.Node
		switch {
		case len(in.images) > 0:
			media = swiss.Gallery(swiss.GalleryProps{Images: in.images, Alt: in.title, Aspect: "16/10",
				Label: tr(lang, "Abb.", "Fig."), Prev: tr(lang, "Vorheriges Bild", "Previous image"), Next: tr(lang, "Nächstes Bild", "Next image")})
		case in.logo != "":
			media = h.Div(h.Class("nz-logo"), h.Img(h.Src(in.logo), h.Alt(in.org), g.Attr("loading", "lazy")))
		}
		detail := make([]g.Node, len(in.details))
		for j, dt := range in.details {
			detail[j] = h.Li(g.Text(dt))
		}
		cls := "nz-panel"
		if len(in.images) == 0 && in.logo != "" {
			cls += " nz-panel-logo"
		}
		out[i] = h.Div(h.Class(cls), g.Attr("data-tab-page", in.key), g.If(in.key != first, g.Attr("hidden", "")),
			g.If(media != nil, h.Div(h.Class("nz-media"), media)),
			h.Div(h.Class("nz-body"),
				h.Div(h.Class("nz-panel-head"),
					swiss.Micro(g.Text(strings.ToUpper(in.kind)+g2(in.org)+g2(in.period)+g2(in.place))),
					h.Div(h.Class("nz-chips"), g.Group(chips)),
				),
				h.H3(h.Class("nz-panel-title"), g.Text(in.title)),
				g.If(in.desc != "", h.P(h.Class("cell-text"), g.Text(truncate(in.desc, 320)))),
				g.If(len(detail) > 0, h.Ul(h.Class("nz-details"), g.Group(detail))),
				g.If(len(in.tags) > 0, tagPills(in.tags)),
				g.If(in.url != "", h.Div(swiss.BtnLink(in.url, true, nil, g.Text(tr(lang, "Projekt ansehen →", "View project →"))))),
			),
		)
	}
	return out
}

const netJS = `(function(){
  var root=document.getElementById('nz'); if(!root) return;
  root.addEventListener('click',function(e){
    var b=e.target.closest('[data-netline]'); if(!b) return;
    var on=b.getAttribute('aria-pressed')!=='true';
    root.querySelectorAll('[data-netline]').forEach(function(x){x.setAttribute('aria-pressed',String(on&&x===b));});
    if(on) root.setAttribute('data-focus',b.getAttribute('data-netline')); else root.removeAttribute('data-focus');
  });
  root.addEventListener('keydown',function(e){
    if(e.key!=='Enter'&&e.key!==' ') return;
    var s=e.target.closest('g[data-tab]'); if(!s) return;
    e.preventDefault(); s.dispatchEvent(new MouseEvent('click',{bubbles:true}));
  });
})();`

const netCSS = `
#nz { --nz-h: var(--second); --nz-g: var(--block); --nz-s: var(--third); --nz-w: var(--ink); }
html[data-mode="dark"] #nz { --nz-h: #ff8a2a; --nz-g: #5f80ff; }
.nz-legend { display: flex; flex-wrap: wrap; justify-content: space-between; align-items: center; gap: .75rem 1.5rem; padding: .9rem var(--swiss-pad); border-bottom: 1px solid var(--ink); }
.nz-hint { color: var(--muted); }
.nz-lines { display: flex; flex-wrap: wrap; gap: .4rem; }
.nz-line { all: unset; box-sizing: border-box; cursor: pointer; display: inline-flex; align-items: center; gap: .5rem; padding: .3rem .7rem .3rem .3rem; border: 1px solid var(--ink); font-weight: 700; font-size: .82rem; }
.nz-line:hover { background: var(--act); color: var(--act-fg); }
.nz-line[aria-pressed="true"] { background: var(--ink); color: var(--bg); }
.nz-line:focus-visible { outline: 3px solid var(--accent); outline-offset: 2px; }
.nz-line-badge, .nz-chip { color: #fff; font-family: var(--font-mono); font-weight: 800; }
.nz-line-badge[data-line="w"], .nz-chip[data-line="w"] { color: var(--bg); }
.nz-line-badge[data-line="h"], .nz-chip[data-line="h"] { color: #0e0f14; }
.nz-line-badge { width: 1.6rem; height: 1.6rem; display: grid; place-items: center; font-size: .8rem; }
.nz-wrap { overflow-x: auto; padding: 0 var(--swiss-pad); }
.nz { display: block; width: 100%; min-width: 860px; height: auto; }
.nz-track { fill: none; stroke-width: 10; stroke-linejoin: round; transition: opacity .2s; }
.nz-badge text { font-family: var(--font-mono); font-size: 13px; font-weight: 800; fill: #fff; text-anchor: middle; }
.nz-badge[data-lines="w"] text { fill: var(--bg); }
.nz-stop { cursor: pointer; transition: opacity .2s; outline: none; }
.nz-mark { fill: var(--bg); stroke: var(--ink); stroke-width: 3; transition: fill .15s; }
.nz-x .nz-mark { stroke-width: 4; }
.nz-hit { fill: transparent; }
.nz-stop:hover .nz-hit { fill: var(--act); }
.nz-stop:hover text { fill: var(--act-fg); }
.nz-stop:hover .nz-mark, .nz-stop:focus-visible .nz-mark { stroke-width: 5; }
.nz-stop[aria-pressed="true"] .nz-mark { fill: var(--accent); }
.nz-stop:focus-visible .nz-hit { fill: var(--act); stroke: var(--ink); stroke-width: 2; }
.nz-stop:focus-visible text { fill: var(--act-fg); }
.nz-n { font-family: var(--font-body); font-size: 13px; font-weight: 800; fill: var(--ink); text-anchor: middle; letter-spacing: -0.01em; }
.nz-stop[aria-pressed="true"] .nz-n { fill: var(--main-text); }
.nz-y { font-family: var(--font-mono); font-size: 9.5px; font-weight: 700; fill: var(--muted); text-anchor: middle; letter-spacing: .04em; }
#nz[data-focus="h"] [data-lines]:not([data-lines~="h"]),
#nz[data-focus="g"] [data-lines]:not([data-lines~="g"]),
#nz[data-focus="s"] [data-lines]:not([data-lines~="s"]),
#nz[data-focus="w"] [data-lines]:not([data-lines~="w"]) { opacity: .12; }
.nz-swipe { display: none; }
@media (max-width: 1024px) {
  .nz-wrap { mask-image: linear-gradient(to right, #000 85%, transparent); -webkit-mask-image: linear-gradient(to right, #000 85%, transparent); }
  .nz-swipe { display: block; padding: .4rem var(--swiss-pad) .8rem; color: var(--main-text); text-align: right; }
}
.nz-panel { display: grid; grid-template-columns: minmax(0, 4fr) minmax(0, 5fr); border-top: 3px solid var(--ink); background: var(--surface-2); }
.nz-panel:not(:has(.nz-media)) { grid-template-columns: minmax(0, 1fr); }
/* a logo is a small mark: one narrow column, the text gets the rest */
.nz-panel-logo { grid-template-columns: minmax(150px, 1fr) minmax(0, 6fr); }
.nz-panel-logo .nz-logo { min-height: 0; padding: 1.25rem; place-items: start center; }
.nz-panel-logo .nz-logo img { max-width: 100%; max-height: 96px; }
.nz-media { box-shadow: 1px 0 0 var(--ink); min-width: 0; display: flex; flex-direction: column; }
.nz-media [data-component="gallery"], .nz-media [data-slot="stage"] { flex: 1; }
.nz-logo { flex: 1; display: grid; place-items: center; padding: 2rem; background: #fff; min-height: 220px; }
.nz-logo img { max-width: 70%; max-height: 160px; object-fit: contain; }
.nz-body { display: flex; flex-direction: column; gap: .6rem; padding: var(--swiss-pad); }
.nz-details { margin: 0; padding-left: 1.1rem; display: grid; gap: .25rem; font-size: .88rem; }
.ys { padding: .9rem var(--swiss-pad) 1rem; border-top: 1px solid var(--ink); }
.ys-row { display: grid; grid-template-columns: 1.6rem 1fr; gap: .6rem; align-items: center; margin-top: 4px; }
.ys-key { height: 1.6rem; display: grid; place-items: center; color: #fff; font-family: var(--font-mono); font-size: .72rem; font-weight: 800; }
.ys-key[data-line="w"] { color: var(--bg); }
.ys-track { position: relative; height: 1.6rem; background: repeating-linear-gradient(to right, var(--rule-subtle) 0 1px, transparent 1px calc(100% / 12)); }
.ys-axis .ys-track { height: 1.1rem; background: none; }
.ys-axis .ys-track span { position: absolute; top: 0; font-family: var(--font-mono); font-size: .62rem; font-weight: 700; color: var(--muted); padding-left: 2px; border-left: 1px solid var(--rule-subtle); }
.ys-bar { all: unset; box-sizing: border-box; position: absolute; top: 2px; bottom: 2px; cursor: pointer; background: var(--c); color: #fff; overflow: hidden; display: flex; align-items: center; padding: 0 .35rem; box-shadow: 0 0 0 1px var(--bg); }
.ys-bar span { font-family: var(--font-mono); font-size: .6rem; font-weight: 700; text-transform: uppercase; letter-spacing: .04em; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.ys-bar:hover { background: var(--act); color: var(--act-fg); outline: 1px solid var(--act-fg); z-index: 1; }
.ys-bar[aria-pressed="true"] { outline: 3px solid var(--accent); outline-offset: 1px; z-index: 2; }
.ys-bar:focus-visible { outline: 3px solid var(--accent); outline-offset: 2px; }
html[data-mode="dark"] .ys-bar { color: #0b0d16; }
.ys-bar[style*="--nz-h"] { color: #0e0f14; }
.nz-panel-head { display: flex; flex-wrap: wrap; justify-content: space-between; gap: .5rem; align-items: center; }
.nz-chips { display: flex; flex-wrap: wrap; gap: .3rem; }
.nz-chip { font-size: .66rem; padding: .2rem .45rem; text-transform: uppercase; letter-spacing: .05em; }
.nz-panel-title { margin: 0; font-size: clamp(1.4rem, 1.1rem + 1.2vw, 2.2rem); font-weight: 800; letter-spacing: -0.03em; line-height: 1.05; }
.nz-panel .cell-text { max-width: 75ch; margin: 0; }
@media (max-width: 1024px) {
  .nz-panel, .nz-panel:not(:has(.nz-media)) { grid-template-columns: minmax(0, 1fr); }
  .nz-media { box-shadow: 0 1px 0 var(--ink); }
  .nz-logo { min-height: 140px; }
  .nz-panel-logo { grid-template-columns: minmax(0, 1fr); }
  .nz-panel-logo .nz-logo { min-height: 0; place-items: center start; }
  .nz-panel-logo .nz-logo img { max-height: 56px; }
  .ys-axis .ys-track span:nth-child(even) { display: none; }
}
`
