package pages

import (
	"fmt"
	"strings"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"

	"mljr-web/internal/i18n"
	"mljr-web/ui/icon"
	"mljr-web/ui/swiss"
)

// codeShowcaseSection ("Maschinenraum") explains how this site is built:
// the request path as a numbered diagram and four headline figures, then
// real source excerpts as printed sheets — line numbers in the gutter,
// comments set in the accent colour, the caption as a margin note.
func codeShowcaseSection(p swiss.SectionProps, lang string) g.Node {
	p.Note = i18n.T(lang, "sections.codeshowcase.sub")
	examples := []struct {
		Title    string
		Caption  string
		Filename string
		Code     string
	}{
		{tr(lang, "Netzplan", "Network map"), tr(lang, "Die Gleise im Netzplan: Wegpunkte auf einem Raster, gerade Strecke plus 45°-Diagonale, als SVG-Pfad auf dem Server gerechnet.", "The Netzplan tracks: grid waypoints, a straight run plus a 45° diagonal, computed into an SVG path on the server."), "projects/homepage/pages/netzplan.go", netExcerpt},
		{tr(lang, "Abfrage", "Poller"), "The homelab panel polls VictoriaMetrics — range queries over 40 days come back empty, so the attack heatmap is fetched in 30-day chunks.", "projects/homepage/homelab/homelab.go", pollerExcerpt},
		{tr(lang, "Reiter", "Tabs"), tr(lang, "Diese Blätter, die Tafeln und der Netzplan teilen ein kleines Skript: ein Klick-Listener für die ganze Seite, Pfeiltasten, und Spiele laufen nur, solange ihr Reiter offen ist.", "These sheets, the plates and the Netzplan share one small script: one click listener for the whole page, arrow keys, and games only run while their tab is open."), "projects/homepage/pages/tabs.go · tabsJS.js", tabsExcerpt},
		{"Live", "Live updates without a framework: one attribute on the section, one SSE handler, Datastar morphs the panel in place.", "projects/homepage/handlers.go", sseExcerpt},
		{"Radar", "The skills radar — plain trigonometry rendered to an SVG polygon at request time.", "ui/data/radarchart.go", radarExcerpt},
	}

	tabs := make([]g.Node, len(examples))
	sheets := make([]g.Node, len(examples))
	for i, ex := range examples {
		n := fmt.Sprint(i + 1)
		tabs[i] = h.Button(h.Type("button"), h.Class("sheet-tab"), g.Attr("data-tab", n), g.Attr("aria-pressed", fmt.Sprint(i == 0)),
			h.Span(h.Class("sheet-tab-n"), g.Text(fmt.Sprintf("%02d", i+1))),
			h.Span(g.Text(ex.Title)),
		)
		sheets[i] = h.Article(h.Class("sheet"), g.Attr("data-tab-page", n), g.If(i > 0, g.Attr("hidden", "")),
			h.Header(h.Class("sheet-head"),
				h.Span(g.Text(fmt.Sprintf("%s %02d / %02d", tr(lang, "Blatt", "Sheet"), i+1, len(examples)))),
				h.Code(g.Text(ex.Filename)),
			),
			h.Div(h.Class("sheet-body"),
				printedCode(ex.Code, ex.Filename),
				h.Div(h.Class("sheet-note"),
					h.Span(h.Class("sheet-note-n"), g.Text(n)),
					h.P(g.Text(ex.Caption)),
				),
			),
			tabPager(lang, i, len(examples), func(j int) string { return examples[j].Title }, tr(lang, "Blatt", "Sheet")),
		)
	}

	return swiss.Section(p,
		swiss.Cell(swiss.CellProps{Span: 5},
			swiss.Kicker("A. "+tr(lang, "Anfrageweg", "Request path"), g.Text("1 → 4")),
			hoodFlow(lang),
		),
		swiss.Cell(swiss.CellProps{Span: 4, Variant: "block"},
			swiss.Kicker("B. "+tr(lang, "Kennzahlen", "Figures"), nil),
			h.Dl(h.Class("hood-figs"),
				hoodFig("175", tr(lang, "UI-Komponenten", "UI components")),
				hoodFig("0", tr(lang, "npm-Laufzeitabhängigkeiten", "npm runtime deps")),
				hoodFig("4×2", tr(lang, "Themes × Modi", "themes × modes")),
				hoodFig("1", tr(lang, "Binary zum Ausliefern", "binary to deploy")),
			),
			h.Div(h.Class("cell-foot"),
				swiss.BtnLink("https://github.com/MrCodeEU", false, []g.Node{g.Attr("target", "_blank"), g.Attr("rel", "noopener")},
					icon.Icon("simple-icons:go", icon.Props{Size: "1em"}), g.Text(tr(lang, "Quellcode ansehen ↗", "Explore the source ↗"))),
			),
		),
		swiss.Cell(swiss.CellProps{Variant: "flush", Attrs: []g.Node{h.ID("sheets"), g.Attr("data-tabs", "")}},
			h.Div(h.Class("sheet-tabs"), g.Attr("role", "group"), g.Attr("aria-label", tr(lang, "Auszüge", "Excerpts")),
				tabLead(tr(lang, "Blatt wählen", "Pick a sheet"), len(examples)), g.Group(tabs)),
			g.Group(sheets),
		),
	)
}

func hoodFig(v, label string) g.Node {
	return h.Div(h.Dt(g.Text(label)), h.Dd(g.Text(v)))
}

// hoodFlow: the request path as numbered steps on a vertical rule.
func hoodFlow(lang string) g.Node {
	steps := []struct{ title, sub string }{
		{"Browser", tr(lang, "kein Framework zur Laufzeit — nur Datastar (14 kB)", "no framework runtime — Datastar (14 kB) only")},
		{"Caddy", tr(lang, "TLS, CrowdSec und Authelia davor", "TLS, CrowdSec + Authelia in front")},
		{tr(lang, "Ein Go-Binary", "One Go binary"), tr(lang, "gomponents rendern HTML, SVG-Diagramme pro Anfrage berechnet", "gomponents render HTML, SVG charts computed per request")},
		{"SSE", tr(lang, "Server schickt Fragmente, Datastar morpht sie an Ort und Stelle", "server pushes fragments, Datastar morphs them in place")},
	}
	items := make([]g.Node, len(steps))
	for i, s := range steps {
		items[i] = h.Li(h.Span(h.Class("hf-n"), g.Text(fmt.Sprint(i+1))), h.Div(h.B(g.Text(s.title)), h.P(g.Text(s.sub))))
	}
	return h.Ol(h.Class("hf"), g.Group(items))
}

// printedCode renders Go source with gutter line numbers, tokenised by
// chroma and coloured from the swiss palette (keywords blue, strings
// magenta, numbers and types orange, functions bold, comments muted).
// printedCode highlights src with the lexer matching filename's extension.
func printedCode(src, filename string) g.Node {
	src = strings.TrimRight(src, "\n")
	lexer := lexers.Match(filename)
	if lexer == nil {
		lexer = lexers.Get("go")
	}
	it, err := lexer.Tokenise(nil, src)
	if err != nil {
		return h.Pre(h.Class("pc"), h.Code(g.Text(src)))
	}
	// Adjacent tokens of the same class share one span, and line numbers
	// come from a CSS counter, which keeps the DOM small.
	type run struct{ cls, text string }
	var lines [][]run
	cur := []run{}
	for tok := it(); tok != chroma.EOF; tok = it() {
		cls := tokenClass(tok.Type)
		parts := strings.Split(tok.Value, "\n")
		for i, part := range parts {
			if i > 0 {
				lines = append(lines, cur)
				cur = []run{}
			}
			if part == "" {
				continue
			}
			if n := len(cur); n > 0 && (cur[n-1].cls == cls || strings.TrimSpace(part) == "" && cur[n-1].cls == "") {
				cur[n-1].text += part
			} else {
				cur = append(cur, run{cls, part})
			}
		}
	}
	lines = append(lines, cur)
	rows := make([]g.Node, len(lines))
	for i, ln := range lines {
		nodes := make([]g.Node, len(ln))
		for j, r := range ln {
			if r.cls == "" {
				nodes[j] = g.Text(r.text)
			} else {
				nodes[j] = h.Span(h.Class(r.cls), g.Text(r.text))
			}
		}
		rows[i] = h.Span(h.Class("pc-l"), g.Group(nodes), g.Text("\n"))
	}
	// focusable, so keyboard users can scroll long lines sideways
	return h.Pre(h.Class("pc"), g.Attr("tabindex", "0"), g.Attr("aria-label", filename), h.Code(rows...))
}

func tokenClass(t chroma.TokenType) string {
	switch {
	case t.InCategory(chroma.Comment):
		return "pc-c"
	case t == chroma.KeywordType || t == chroma.NameBuiltin:
		return "pc-t"
	case t.InCategory(chroma.Keyword):
		return "pc-k"
	case t.InSubCategory(chroma.LiteralString):
		return "pc-s"
	case t.InSubCategory(chroma.LiteralNumber):
		return "pc-num"
	case t == chroma.NameFunction:
		return "pc-f"
	}
	return ""
}

const hoodCSS = `
.hf { list-style: none; margin: 0; padding: 0; border-left: 3px solid var(--accent); }
.hf li { display: grid; grid-template-columns: 3.2rem 1fr; gap: .4rem; padding: .7rem 0; border-bottom: 1px solid var(--rule-subtle); }
.hf li:last-child { border-bottom: 0; }
.hf-n { font-size: 2.2rem; font-weight: 900; letter-spacing: -0.05em; line-height: .9; text-align: center; color: var(--main-text); }
.hf b { font-size: 1.05rem; letter-spacing: -0.02em; }
.hf p { margin: .15rem 0 0; font-size: .85rem; color: var(--muted); }
.hood-figs { margin: 0; display: grid; grid-template-columns: 1fr 1fr; }
.hood-figs > div { padding: .6rem 0; }
.hood-figs > div:nth-child(n+3) { border-top: 1px solid currentColor; }
.hood-figs dt { font-family: var(--font-mono); font-size: .66rem; text-transform: uppercase; letter-spacing: .06em; opacity: .85; }
.hood-figs dd { margin: .2rem 0 0; font-size: clamp(2rem, 3.5vw, 3.2rem); font-weight: 900; letter-spacing: -0.05em; line-height: 1; }
.sheet-tabs { display: flex; flex-wrap: wrap; border-bottom: 1px solid var(--ink); }
.sheet-tab { all: unset; cursor: pointer; display: flex; gap: .5rem; align-items: baseline; padding: .7rem 1rem; font-weight: 700; font-size: .85rem; box-shadow: 1px 0 0 var(--ink); }
.sheet-tab-n { font-family: var(--font-mono); font-size: .7rem; color: var(--muted); }
.sheet-tab::after { content: "▸"; width: 1.1rem; height: 1.1rem; display: inline-grid; place-items: center; font-size: .65rem; background: var(--act); color: var(--act-fg); }
.sheet-tab[aria-pressed="true"]::after { content: "●"; color: inherit; background: none; }
.sheet-tab:not([aria-pressed="true"]):hover { background: var(--act); color: var(--act-fg); }
.sheet-tab:not([aria-pressed="true"]):hover .sheet-tab-n { color: inherit; }
.sheet-tabs .tab-lead { flex-direction: row; align-items: center; gap: .75rem; padding-block: 0; }
.sheet-tab[aria-pressed="true"] { background: var(--ink); color: var(--bg); }
.sheet-tab[aria-pressed="true"] .sheet-tab-n { color: inherit; }
.sheet-tab:focus-visible { outline: 3px solid var(--accent); outline-offset: -3px; }
.sheet { background: var(--surface-2); }
.sheet-head { display: flex; justify-content: space-between; gap: 1rem; padding: .5rem var(--swiss-pad); font-family: var(--font-mono); font-size: .7rem; text-transform: uppercase; letter-spacing: .06em; border-bottom: 1px solid var(--rule-subtle); }
.sheet-head code { text-transform: none; color: var(--main-text); }
.sheet-body { display: grid; grid-template-columns: minmax(0, 1fr) 16rem; }
.pc { margin: 0; padding: 1rem 0; overflow-x: auto; font-family: var(--font-mono); font-size: .78rem; line-height: 1.6; tab-size: 4; }
.pc code { display: block; min-width: max-content; counter-reset: pc; }
.pc-l { display: block; padding-right: 1rem; }
.pc-l::before { counter-increment: pc; content: counter(pc); display: inline-block; width: 3.2rem; padding-right: .9rem; margin-right: .9rem; text-align: right; color: var(--muted); border-right: 1px solid var(--rule-subtle); user-select: none; }
.pc-c { color: var(--muted); font-style: italic; }
.pc-k { color: var(--block); font-weight: 700; }
html[data-mode="dark"] .pc-k { color: #8ea6ff; }
html[data-mode="dark"] .pc-s { color: #f2679f; }
.pc-t { color: var(--second); font-weight: 600; }
.pc-s { color: var(--third); }
.pc-num { color: var(--second); }
.pc-f { color: var(--ink); font-weight: 800; }
.sheet-note { border-left: 1px solid var(--ink); padding: 1rem; display: flex; flex-direction: column; gap: .6rem; }
.sheet-note-n { width: 2rem; height: 2rem; display: grid; place-items: center; background: var(--second); color: var(--second-fg); font-weight: 900; }
.sheet-note p { margin: 0; font-size: .88rem; line-height: 1.5; }
@media (max-width: 1024px) {
  .sheet-body { grid-template-columns: minmax(0, 1fr); }
  .sheet-note { border-left: 0; border-top: 1px solid var(--ink); order: -1; flex-direction: row; align-items: flex-start; }
  .pc-l::before { width: 2.2rem; }
}
`

// The excerpts below are real source from this repo, lightly trimmed for
// display (backticks swapped for quotes where needed). Keep them in sync
// when the originals change materially.

const netExcerpt = `// Netzplan tracks: each skill line is a list of grid
// waypoints. Between two waypoints the track runs straight,
// then turns 45° into the next one, like a metro map.
for _, ln := range netLines {
	route := netRoutes[ln.key]
	var path strings.Builder
	for i, pt := range route {
		px, py := x(pt[0]), y(pt[1])
		if i == 0 {
			fmt.Fprintf(&path, "M%.1f %.1f", px, py)
			continue
		}
		lx, ly := x(route[i-1][0]), y(route[i-1][1])
		dx, dy := px-lx, py-ly
		dd := math.Min(math.Abs(dx), math.Abs(dy))
		// straight run, then the diagonal
		fmt.Fprintf(&path, " L%.1f %.1f L%.1f %.1f",
			px-math.Copysign(dd, dx), ly, px, py)
	}
	fmt.Fprintf(&b, '<path d="%s" class="nz-track" stroke="%s"/>',
		path.String(), ln.color)
}`

const pollerExcerpt = `// Attacks blocked per day, last ~12 months.
// VictoriaMetrics range queries spanning >40 days can return
// empty results (per-day index cutoff), so fetch the year in
// 30-day chunks and dedupe by day.
seenDay := map[string]bool{}
for chunk := 11; chunk >= 0; chunk-- {
	start := now.AddDate(0, 0, -30*(chunk+1))
	end := now.AddDate(0, 0, -30*chunk)
	pts, err := p.promQueryRange(ctx,
		"sum(increase(cs_bucket_overflowed_total[1d]))",
		start, end, 24*time.Hour,
	)
	if err != nil {
		continue // empty until a year of data exists
	}
	for _, pt := range pts {
		day := time.Unix(int64(pt[0]), 0)
		key := day.Format("2006-01-02")
		if pt[1] > 0 && !seenDay[key] {
			seenDay[key] = true
			s.AttackDays = append(s.AttackDays,
				DayValue{Date: day, Count: int(pt[1])})
		}
	}
}`

const tabsExcerpt = `// One delegated listener drives every tab strip on the page:
// plates, code sheets, the lab and the Netzplan stops.
function show(root, b) {
  var n = b.getAttribute('data-tab');
  root.dispatchEvent(new CustomEvent('tabs:select',
    {bubbles: true, detail: {id: n}}));   // the map listens
  root.querySelectorAll('[data-tab-page]').forEach(function(s) {
    s.hidden = s.getAttribute('data-tab-page') !== n;
    // embedded games run only while their page is visible
    s.querySelectorAll('iframe[data-src]').forEach(function(f) {
      if (s.hidden) f.removeAttribute('src');
      else if (!f.getAttribute('src'))
        f.setAttribute('src', f.getAttribute('data-src'));
    });
  });
}
document.addEventListener('click', function(e) {
  var b = e.target.closest('[data-tab]');
  if (!b) return;
  var root = b.closest('[data-tabs]');
  if (root) show(root, b);
});`

const sseExcerpt = `// One attribute on the section…
h.Section(
	h.ID("homelab"),
	g.Attr("data-on-interval__duration.60s",
		"@get('/api/homelab')"),
	// …
)

// …one SSE handler on the server.
e.GET("/api/homelab", func(c echo.Context) error {
	sse := datastar.NewSSE(c.Response().Writer, c.Request())
	return sse.PatchElements(
		web.RenderToString(pages.HomelabPanel(snapshot())),
	)
})

// Datastar morphs #homelab-panel in place every 60 s:
// no reload, no virtual DOM, no client-side state to sync.`

const radarExcerpt = `// Radar chart: axis i sits at angle 2πi/n, starting at
// 12 o'clock. Values scale along the spoke; the series
// becomes a single <polygon>.
angle := func(i int) float64 {
	return float64(i)*2*math.Pi/float64(n) - math.Pi/2
}
pt := func(radius float64, i int) (float64, float64) {
	a := angle(i)
	return cx + radius*math.Cos(a), cy + radius*math.Sin(a)
}

pts := make([]string, len(p.Axes))
for i := range p.Axes {
	scaled := (s.Values[i] / p.Max) * r
	x, y := pt(scaled, i)
	pts[i] = fmt.Sprintf("%.1f,%.1f", x, y)
}
sb.WriteString(fmt.Sprintf(
	'<polygon points="%s" fill="%s" fill-opacity="0.18"
	   stroke="%s" stroke-width="2"/>',
	strings.Join(pts, " "), color, color))`
