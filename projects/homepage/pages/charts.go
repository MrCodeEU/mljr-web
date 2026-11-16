package pages

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"mljr-web/ui/icon"
	"mljr-web/ui/swiss"
)

// Swiss chart primitives for the homepage: flat fills, hairline rules,
// mono labels, no gradients. All SVGs scale to their container width.

// fieldColors cycles the accent fields for categorical data.
var fieldColors = []string{"var(--block)", "var(--second)", "var(--third)", "var(--ink)", "var(--muted)"}

type dayCount struct {
	Date  time.Time
	Count int
}

// calendarSVG draws a week-column heatmap (Mon top) of the given days, last
// day rightmost. Levels: 0 = empty cell, 1–4 = tints of colour.
// unit names what is counted, for the per-day tooltip ("12 commits").
func calendarSVG(days []dayCount, color, unit, lang string) g.Node {
	if len(days) == 0 {
		return nil
	}
	sort.Slice(days, func(i, j int) bool { return days[i].Date.Before(days[j].Date) })
	first := days[0].Date
	// start on the Monday of the first week
	off := (int(first.Weekday()) + 6) % 7
	start := first.AddDate(0, 0, -off)
	byDay := map[string]int{}
	maxV := 0
	for _, d := range days {
		byDay[d.Date.Format("2006-01-02")] = d.Count
		maxV = max(maxV, d.Count)
	}
	last := days[len(days)-1].Date
	weeks := int(last.Sub(start).Hours()/24)/7 + 1
	const cell, gap, top, left = 12.0, 3.0, 18.0, 24.0
	w := left + float64(weeks)*(cell+gap)
	hgt := top + 7*(cell+gap)
	var b strings.Builder
	fmt.Fprintf(&b, `<svg viewBox="0 0 %.0f %.0f" width="100%%" style="max-width:%.0fpx" role="img" aria-label="%s" class="cal">`, w, hgt, w*1.6, tr(lang, "Aktivitätskalender", "Activity calendar"))
	months := []string{"Jan", "Feb", "Mär", "Apr", "Mai", "Jun", "Jul", "Aug", "Sep", "Okt", "Nov", "Dez"}
	if lang != "de" {
		months = []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
	}
	for i, l := range []string{"Mo", "", "Mi", "", "Fr", "", ""} {
		if l != "" {
			fmt.Fprintf(&b, `<text x="0" y="%.1f" class="cal-l">%s</text>`, top+float64(i)*(cell+gap)+cell-2, l)
		}
	}
	lastMonth := -1
	for wk := 0; wk < weeks; wk++ {
		x := left + float64(wk)*(cell+gap)
		for dow := 0; dow < 7; dow++ {
			day := start.AddDate(0, 0, wk*7+dow)
			if day.Before(first) || day.After(last) {
				continue
			}
			if dow == 0 && int(day.Month()) != lastMonth && day.Day() <= 7 {
				lastMonth = int(day.Month())
				fmt.Fprintf(&b, `<text x="%.1f" y="11" class="cal-l">%s</text>`, x, months[day.Month()-1])
			}
			v := byDay[day.Format("2006-01-02")]
			fill := "var(--rule-subtle)"
			if v > 0 && maxV > 0 {
				lvl := int(math.Ceil(float64(v) / float64(maxV) * 4))
				fill = fmt.Sprintf("color-mix(in srgb, %s %d%%, var(--bg))", color, []int{0, 35, 60, 82, 100}[lvl])
			}
			fmt.Fprintf(&b, `<rect x="%.1f" y="%.1f" width="%.0f" height="%.0f" fill="%s" data-t="%s · %d %s"/>`,
				x, top+float64(dow)*(cell+gap), cell, cell, fill, day.Format("02.01.2006"), v, unit)
		}
	}
	b.WriteString(`</svg>`)
	return g.Group{g.Raw(b.String()), h.Script(g.Raw(calTipJS))}
}

// calTipJS shows one shared tooltip for all calendar days (a <title> per
// day would double the calendar's DOM). Runs once per page.
const calTipJS = `(function(){
  if (window.__calTip) return; window.__calTip = 1;
  var tip = document.createElement('div'); tip.className = 'cal-tip'; tip.hidden = true; document.body.appendChild(tip);
  document.addEventListener('mouseover', function(e){
    var r = e.target.closest && e.target.closest('svg.cal rect[data-t]');
    if (!r) { tip.hidden = true; return; }
    var b = r.getBoundingClientRect();
    tip.textContent = r.getAttribute('data-t');
    tip.style.left = (b.left + b.width / 2) + 'px'; tip.style.top = b.top + 'px'; tip.hidden = false;
  });
  window.addEventListener('scroll', function(){ tip.hidden = true; }, { passive: true });
})();`

type segment struct {
	Label string
	Value float64
}

// propBar is one full-width bar split into proportional colour segments,
// with a numbered legend underneath.
func propBar(segs []segment, unit string) g.Node {
	total := 0.0
	for _, s := range segs {
		total += s.Value
	}
	if total <= 0 {
		return nil
	}
	parts := make([]g.Node, 0, len(segs))
	legend := make([]g.Node, 0, len(segs))
	for i, s := range segs {
		pct := s.Value / total * 100
		col := fieldColors[i%len(fieldColors)]
		parts = append(parts, h.Span(h.Class("pbar-seg"), h.Style(fmt.Sprintf("flex:%.3f 0 0;background:%s", pct, col)), h.Title(s.Label)))
		val := fmt.Sprintf("%.0f%%", pct)
		if unit != "%" {
			val = fmt.Sprintf("%.0f %s", s.Value, unit)
		}
		legend = append(legend, h.Li(h.I(h.Style("background:"+col)), h.B(g.Text(s.Label)), h.Span(g.Text(val))))
	}
	return h.Div(h.Class("pbar"),
		h.Div(h.Class("pbar-track"), g.Group(parts)),
		h.Ol(h.Class("pbar-legend"), g.Group(legend)),
	)
}

type figure struct {
	Label, Value, Sub, Variant string
	Icon                       string // optional, drawn inside a Bauhaus shape
}

// figRow lays figures out as equal ruled columns inside one cell.
func figRow(figs []figure) g.Node {
	nodes := make([]g.Node, len(figs))
	for i, f := range figs {
		cls := "figrow-item"
		if f.Variant != "" {
			cls += " mono-" + f.Variant
		}
		nodes[i] = h.Div(h.Class(cls),
			g.If(f.Icon != "", h.Span(h.Class("fig-shape"), icon.Icon(f.Icon, icon.Props{Size: "1.15rem"}))),
			swiss.Micro(g.Text(f.Label)),
			swiss.Figure(g.Text(f.Value)),
			g.If(f.Sub != "", h.Span(h.Class("figrow-sub"), g.Text(f.Sub))),
		)
	}
	return h.Div(h.Class("figrow"), h.Style(fmt.Sprintf("--n:%d", len(figs))), g.Group(nodes))
}

type colPoint struct {
	Label string // x label (shown when Tick)
	Tick  bool
	Value float64
	Title string
}

// colChart draws vertical columns on a hairline baseline.
func colChart(pts []colPoint, color string, height float64) g.Node {
	if len(pts) == 0 {
		return nil
	}
	maxV := 0.0
	for _, p := range pts {
		maxV = math.Max(maxV, p.Value)
	}
	if maxV == 0 {
		maxV = 1
	}
	const colW, gap, bottom = 10.0, 3.0, 18.0
	w := float64(len(pts)) * (colW + gap)
	var b strings.Builder
	fmt.Fprintf(&b, `<svg viewBox="0 0 %.0f %.0f" width="100%%" preserveAspectRatio="none" class="colchart" style="height:%.0fpx">`, w, height, height)
	for i, p := range pts {
		x := float64(i) * (colW + gap)
		bh := p.Value / maxV * (height - bottom - 4)
		fmt.Fprintf(&b, `<rect x="%.1f" y="%.1f" width="%.0f" height="%.1f" fill="%s"><title>%s</title></rect>`, x, height-bottom-bh, colW, bh, color, p.Title)
		if p.Tick {
			fmt.Fprintf(&b, `<line x1="%.1f" x2="%.1f" y1="0" y2="%.0f" class="colchart-tick"/>`, x-gap/2, x-gap/2, height-bottom)
		}
	}
	fmt.Fprintf(&b, `<line x1="0" x2="%.0f" y1="%.0f" y2="%.0f" class="colchart-base"/></svg>`, w, height-bottom, height-bottom)
	// x labels as HTML so the stretched SVG doesn't distort the text
	labels := make([]g.Node, 0)
	for i, p := range pts {
		if p.Tick {
			labels = append(labels, h.Span(h.Style(fmt.Sprintf("left:%.2f%%", float64(i)/float64(len(pts))*100)), g.Text(p.Label)))
		}
	}
	return h.Div(h.Class("colchart-wrap"), g.Raw(b.String()), h.Div(h.Class("colchart-x"), g.Group(labels)))
}

type lineSeries struct {
	Label  string
	Color  string
	Points []float64
}

// lineChart draws series on a shared y scale with a zero rule.
func lineChart(series []lineSeries, height float64, aria string) g.Node {
	n := 0
	lo, hi := 0.0, 0.0
	for _, s := range series {
		n = max(n, len(s.Points))
		for _, v := range s.Points {
			lo, hi = math.Min(lo, v), math.Max(hi, v)
		}
	}
	if n < 2 || hi == lo {
		return nil
	}
	const w = 600.0
	y := func(v float64) float64 { return 6 + (hi-v)/(hi-lo)*(height-12) }
	var b strings.Builder
	fmt.Fprintf(&b, `<svg viewBox="0 0 %.0f %.0f" width="100%%" preserveAspectRatio="none" role="img" aria-label="%s" class="linechart" style="height:%.0fpx">`, w, height, aria, height)
	fmt.Fprintf(&b, `<line x1="0" x2="%.0f" y1="%.1f" y2="%.1f" class="colchart-base"/>`, w, y(0), y(0))
	for _, s := range series {
		var pts []string
		for i, v := range s.Points {
			pts = append(pts, fmt.Sprintf("%.1f,%.1f", float64(i)/float64(n-1)*w, y(v)))
		}
		fmt.Fprintf(&b, `<polyline points="%s" fill="none" stroke="%s" stroke-width="2.5" vector-effect="non-scaling-stroke" stroke-linejoin="miter"/>`, strings.Join(pts, " "), s.Color)
	}
	b.WriteString(`</svg>`)
	legend := make([]g.Node, len(series))
	for i, s := range series {
		legend[i] = h.Li(h.I(h.Style("background:"+s.Color)), h.B(g.Text(s.Label)))
	}
	return h.Div(g.Raw(b.String()), g.If(len(series) > 1, h.Ol(h.Class("pbar-legend"), g.Group(legend))))
}

type barItem struct {
	Label, Value string
	N            float64
	Variant      string // "" = main colour, or a field name
}

// barList: one labelled horizontal bar per row, lengths relative to the max.
func barList(items []barItem) g.Node {
	maxV := 0.0
	for _, it := range items {
		maxV = math.Max(maxV, it.N)
	}
	if maxV == 0 {
		maxV = 1
	}
	rows := make([]g.Node, len(items))
	for i, it := range items {
		col := "var(--accent)"
		if it.Variant != "" {
			col = "var(--" + it.Variant + ")"
		}
		rows[i] = h.Li(
			h.Span(h.Class("bl-label"), g.Text(it.Label)),
			h.Span(h.Class("bl-track"), h.Span(h.Class("bl-bar"), h.Style(fmt.Sprintf("width:%.1f%%;background:%s", it.N/maxV*100, col)))),
			h.Span(h.Class("bl-val"), g.Text(it.Value)),
		)
	}
	return h.Ol(h.Class("barlist"), g.Group(rows))
}

// chartCSS: shared styles for the chart primitives above.
const chartCSS = `
.cal-tip { position: fixed; z-index: 300; transform: translate(-50%, calc(-100% - 6px)); pointer-events: none; padding: .25rem .45rem; background: var(--ink); color: var(--bg); font-family: var(--font-mono); font-size: .7rem; white-space: nowrap; }
.cal { display: block; }
.cal-l, .net-l, .net-y { font-family: var(--font-mono); font-size: 9px; font-weight: 700; fill: var(--muted); text-transform: uppercase; }
.pbar-track { display: flex; height: 28px; border: 1px solid var(--ink); }
.pbar-seg + .pbar-seg { border-left: 1px solid var(--bg); }
.pbar-legend { list-style: none; margin: .6rem 0 0; padding: 0; display: flex; flex-wrap: wrap; gap: .35rem 1.1rem; font-family: var(--font-mono); font-size: .72rem; text-transform: uppercase; }
.pbar-legend li { display: flex; align-items: center; gap: .4rem; }
.pbar-legend i { width: 12px; height: 12px; display: inline-block; }
.pbar-legend span { color: var(--muted); }
.figrow { display: grid; grid-template-columns: repeat(var(--n), minmax(0, 1fr)); }
.figrow-item { padding: var(--swiss-pad); display: flex; flex-direction: column; gap: .5rem; box-shadow: 1px 0 0 var(--ink); min-width: 0; }
.figrow-item:last-child { box-shadow: none; }
/* six figures: two rows of three until there is room for labels and icons */
@media (max-width: 1799px) {
  .figrow[style="--n:6"] { grid-template-columns: repeat(3, minmax(0, 1fr)); }
  .figrow[style="--n:6"] .figrow-item { box-shadow: 1px 0 0 var(--ink), 0 1px 0 var(--ink); }
}
.figrow-item { position: relative; }
.figrow-item:has(.fig-shape) > :nth-child(2) { padding-right: 3rem; }
/* Bauhaus shapes cycle square, circle, triangle in the three field colours */
.fig-shape { position: absolute; top: var(--swiss-pad); right: var(--swiss-pad); width: 2.4rem; height: 2.4rem; display: grid; place-items: center; background: var(--block); color: var(--block-fg); }
.figrow-item:nth-child(3n+2) .fig-shape { border-radius: 50%; background: var(--second); color: var(--second-fg); }
.figrow-item:nth-child(3n) .fig-shape { clip-path: polygon(50% 0, 100% 100%, 0 100%); background: var(--third); color: var(--third-fg); place-items: end center; padding-bottom: .3rem; }
.figrow-item:nth-child(3n) .fig-shape [data-component="icon"] { width: .9rem !important; height: .9rem !important; }
:is(.mono-block, .mono-second, .mono-third, .mono-ink) .fig-shape { background: var(--bg); color: var(--ink); }
.figrow-sub { font-size: .8rem; opacity: .75; }
.colchart, .linechart { display: block; width: 100%; }
.colchart-base { stroke: var(--ink); stroke-width: 1; vector-effect: non-scaling-stroke; }
.colchart-tick { stroke: var(--rule-subtle); stroke-width: 1; vector-effect: non-scaling-stroke; }
.colchart-wrap { position: relative; }
.colchart-x { position: relative; height: 1.2rem; margin-top: -18px; }
.colchart-x span { position: absolute; top: 2px; padding-left: 3px; font-family: var(--font-mono); font-size: .66rem; font-weight: 700; color: var(--muted); }
.barlist { list-style: none; margin: 0; padding: 0; display: grid; gap: .35rem; }
.barlist li { display: grid; grid-template-columns: minmax(0, 11rem) 1fr 4.5rem; gap: .6rem; align-items: center; font-size: .82rem; }
.bl-label { overflow: hidden; text-overflow: ellipsis; white-space: nowrap; font-weight: 600; }
.bl-track { height: 12px; background: var(--rule-subtle); }
.bl-bar { display: block; height: 100%; }
.bl-val { font-family: var(--font-mono); font-size: .72rem; font-weight: 700; text-align: right; }
@media (max-width: 1024px) {
  .figrow { grid-template-columns: repeat(2, minmax(0, 1fr)); }
  .figrow-item { box-shadow: 1px 0 0 var(--ink), 0 1px 0 var(--ink); }
  .barlist li { grid-template-columns: minmax(0, 7rem) 1fr 3.8rem; }
}
`
