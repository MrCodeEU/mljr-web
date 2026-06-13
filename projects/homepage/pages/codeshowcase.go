package pages

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	uidata "mljr-web/ui/data"
	"mljr-web/ui/icon"
	"mljr-web/ui/layout"
	"mljr-web/ui/primitive"
	"mljr-web/ui/token"
)

// codeShowcaseSection explains how this site is built: an architecture
// sketch on the left, a paged carousel of real source excerpts on the right
// (highlighted server-side with chroma, paged by the same PaginatedPages
// component example #3 shows).
func codeShowcaseSection() g.Node {
	examples := []struct {
		Caption  string
		Filename string
		Code     string
	}{
		{"The arc math behind the gauges in the open-source section.", "ui/data/gauge.go", gaugeExcerpt},
		{"The homelab panel polls VictoriaMetrics — range queries over 40 days come back empty, so the attack heatmap is fetched in 30-day chunks.", "projects/homepage/homelab/homelab.go", pollerExcerpt},
		{"This very carousel. Pages flip on a Datastar signal; the observer only reacts to display:none→visible, so the animation can never re-trigger itself.", "ui/data/pagination.go", paginatedExcerpt},
		{"Live updates without a framework: one attribute on the section, one SSE handler, Datastar morphs the panel in place.", "projects/homepage/handlers.go", sseExcerpt},
		{"The skills radar — plain trigonometry rendered to an SVG polygon at request time.", "ui/data/radarchart.go", radarExcerpt},
	}

	pages := make([]g.Node, len(examples))
	for i, ex := range examples {
		pages[i] = h.Div(
			h.P(h.Style("margin:0 0 var(--sp-3);font-size:var(--t-sm);font-weight:700;color:var(--muted);line-height:1.5"),
				g.Text(ex.Caption)),
			uidata.SyntaxHighlighter(uidata.SyntaxHighlighterProps{
				Language: "go",
				Theme:    "monokai",
				Filename: ex.Filename,
			}, ex.Code),
		)
	}

	return h.Section(
		h.ID("under-the-hood"),
		h.Style("padding:var(--sp-12) 0;border-top:var(--bw-2) solid var(--ink)"),
		uidata.PaginationSignals("hood", 1),
		layout.Container(layout.ContainerProps{},
			sectionHeader("08", "Under the hood", "this site, in code", token.ToneSky),
			h.Div(
				h.Class("hood-grid"),
				h.Style("display:grid;grid-template-columns:1fr 1.2fr;gap:var(--sp-6);align-items:start"),
				// Left: narrative + architecture sketch + principle chips
				h.Div(
					h.Style("display:flex;flex-direction:column;gap:var(--sp-4);position:sticky;top:var(--sp-8)"),
					h.P(h.Style("font-size:var(--t-lg);font-weight:700;line-height:1.5;margin:0"),
						g.Text("Every component on this page is a Go function. No React, no npm, no build pipeline for the markup — the server renders HTML, Datastar adds reactivity, and SVG charts are computed at request time."),
					),
					archDiagram(),
					h.Div(
						h.Style("display:grid;grid-template-columns:1fr 1fr;gap:var(--sp-3)"),
						hoodFact("lucide:boxes", "175", "UI components"),
						hoodFact("lucide:package-x", "0", "npm runtime deps"),
						hoodFact("lucide:palette", "4×2", "themes × modes"),
						hoodFact("lucide:server", "1", "binary to deploy"),
					),
					h.A(h.Href("https://github.com/MrCodeEU"), g.Attr("target", "_blank"), g.Attr("rel", "noopener"),
						primitive.Button(primitive.ButtonProps{Variant: token.Primary},
							icon.Icon("simple-icons:go"),
							g.Text("Explore the source"),
							icon.Icon("lucide:arrow-up-right"),
						),
					),
				),
				// Right: paged source excerpts — what you see is what runs
				h.Div(
					h.Style("min-width:0"),
					h.Div(h.Style("display:flex;align-items:center;justify-content:space-between;gap:var(--sp-3);margin-bottom:var(--sp-4);flex-wrap:wrap"),
						h.Div(h.Style("font-size:var(--t-xs);font-weight:900;text-transform:uppercase;letter-spacing:.1em;color:var(--muted)"),
							g.Text("5 real excerpts from this site")),
						uidata.Pagination(uidata.PaginationProps{ID: "hood", Total: len(examples), PerPage: 1}),
					),
					uidata.PaginatedPages(uidata.PaginatedPagesProps{ID: "hood", Animation: uidata.PageAnimSlideLeft}, pages...),
				),
			),
		),
	)
}

// archDiagram sketches the request path as stacked neo-brutalist boxes.
func archDiagram() g.Node {
	box := func(bg, title, sub string, ic string) g.Node {
		return h.Div(
			h.Style("border:var(--bw-2) solid var(--ink);box-shadow:var(--shadow-sm);background:"+bg+";padding:var(--sp-3) var(--sp-4);display:flex;align-items:center;gap:var(--sp-3)"),
			icon.Icon(ic, icon.Props{Size: "1.3rem"}),
			h.Div(
				h.Div(h.Style("font-weight:900;font-size:var(--t-sm);line-height:1.2"), g.Text(title)),
				h.Div(h.Style("font-size:var(--t-xs);font-weight:700;color:var(--muted)"), g.Text(sub)),
			),
		)
	}
	arrow := func(label string) g.Node {
		return h.Div(
			h.Style("display:flex;align-items:center;gap:var(--sp-2);padding:2px 0 2px var(--sp-5)"),
			h.Span(h.Style("font-size:1rem;font-weight:900;line-height:1"), g.Text("↓")),
			h.Span(h.Style("font-size:var(--t-xs);font-family:var(--font-mono,monospace);font-weight:700;color:var(--muted)"), g.Text(label)),
		)
	}
	return h.Div(
		h.Style("display:flex;flex-direction:column"),
		box("var(--surface)", "Browser", "no framework runtime — Datastar (14 kB) only", "lucide:globe"),
		arrow("HTTPS · one round trip"),
		box("var(--yellow-bg,#fef08a)", "Caddy ingress", "TLS, CrowdSec + Authelia in front", "lucide:shield"),
		arrow("reverse proxy"),
		box("var(--lime-bg,#d9f99d)", "One Go binary", "gomponents render HTML · SSE patches fragments · SVG charts computed per request", "simple-icons:go"),
	)
}

func hoodFact(ic, value, label string) g.Node {
	return h.Div(
		h.Style("border:var(--bw-2) solid var(--ink);background:var(--surface);box-shadow:var(--shadow-sm);padding:var(--sp-3) var(--sp-4);display:flex;align-items:center;gap:var(--sp-3)"),
		icon.Icon(ic, icon.Props{Size: "1.4rem"}),
		h.Div(
			h.Div(h.Style("font-weight:900;font-size:var(--t-xl);line-height:1"), g.Text(value)),
			h.Div(h.Style("font-size:var(--t-xs);font-weight:700;color:var(--muted);text-transform:uppercase;letter-spacing:.08em;margin-top:2px"), g.Text(label)),
		),
	)
}

// The excerpts below are real source from this repo, lightly trimmed for
// display (backticks swapped for quotes where needed). Keep them in sync
// when the originals change materially.

const gaugeExcerpt = `// Gauge renders a circular gauge as inline SVG.
// Zero JS, zero dependencies — pure server-side Go.
func Gauge(p GaugeProps) g.Node {
	cx, cy := float64(p.Size)/2, float64(p.Size)/2
	radius := float64(p.Size) * 0.38
	strokeW := float64(p.Size) * 0.1

	// Arc spans 270° (from 135° to 45°)
	startAngle := 135.0
	totalAngle := 270.0
	pct := (p.Value - p.Min) / (p.Max - p.Min)

	toRad := func(deg float64) float64 {
		return deg * math.Pi / 180
	}
	arcPt := func(deg float64) (float64, float64) {
		a := toRad(deg)
		return cx + radius*math.Cos(a),
			cy + radius*math.Sin(a)
	}
	arcPath := func(start, end float64, color string) string {
		x1, y1 := arcPt(start)
		x2, y2 := arcPt(end)
		large := 0
		if math.Abs(end-start) > 180 {
			large = 1
		}
		return fmt.Sprintf(
			'<path d="M %.1f %.1f A %.1f %.1f 0 %d 1 %.1f %.1f"
			   fill="none" stroke="%s" stroke-width="%.1f"
			   stroke-linecap="round"/>',
			x1, y1, radius, radius, large, x2, y2,
			color, strokeW)
	}

	// Track (full arc), then value arc on top
	sb.WriteString(arcPath(startAngle, startAngle+totalAngle, p.TrackColor))
	sb.WriteString(arcPath(startAngle, startAngle+totalAngle*pct, p.Color))
	// …
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

const paginatedExcerpt = `// PaginatedPages: each page shown while the shared
// Datastar signal matches its index.
for i, page := range pages {
	wrapped[i] = h.Div(
		g.Attr("data-slot", "page"),
		ui.Show(fmt.Sprintf("$%s === %d", sig, i)),
		page,
	)
}

// Entrance animation replays exactly once per page switch:
// the observer only reacts to style mutations whose OLD value
// was display:none — the CSS animation itself can never
// re-trigger it.
var obs=new MutationObserver(function(muts){
  muts.forEach(function(m){
    var was=(m.oldValue||'');
    var wasHidden=was.indexOf('display:none')>-1;
    if(wasHidden&&m.target.style.display!=='none'){
      m.target.removeAttribute('data-anim');
      void m.target.offsetWidth;      // restart keyframes
      m.target.setAttribute('data-anim','');
    }
  });
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
