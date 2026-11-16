package pages

import (
	"fmt"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

// bhShape is one form in a section's Bauhaus mark: a grid area on a 4×4
// raster, a kind (see bauhausCSS) and a colour role (ink, bg, blk, sec, thr).
type bhShape struct {
	area, kind, color string
}

// bhMarks are the compositions for the section anchors, one per section, in
// the footer's language: flat primaries, circles, triangles, hard edges.
var bhMarks = [][]bhShape{
	{{"1/1/4/4", "circle", "sec"}, {"3/3/5/5", "sq", "ink"}, {"1/4/2/5", "dot", "thr"}},
	{{"1/1/5/5", "tri", "thr"}, {"3/1/5/3", "half", "bg"}, {"1/4/2/5", "circle", "ink"}},
	{{"1/1/5/3", "stripes", "ink"}, {"2/2/5/5", "circle", "sec"}},
	{{"1/1/3/3", "quarter", "ink"}, {"1/3/3/5", "sq", "sec"}, {"3/1/5/5", "half", "blk"}, {"3/2/5/4", "dot", "bg"}},
	{{"1/1/5/5", "ring", "ink"}, {"2/2/4/4", "sq", "thr"}},
	{{"2/1/5/4", "tri", "blk"}, {"1/3/3/5", "circle", "thr"}, {"4/4/5/5", "sq", "ink"}},
	{{"1/1/5/3", "halfr", "thr"}, {"1/3/5/5", "halfl", "blk"}, {"2/2/4/4", "dot", "bg"}},
	{{"1/1/3/5", "half", "sec"}, {"3/1/5/3", "sq", "ink"}, {"3/3/5/5", "quarter2", "thr"}},
	{{"1/1/3/3", "sq", "blk"}, {"3/3/5/5", "sq", "ink"}, {"1/3/3/5", "rtri", "sec"}, {"3/1/5/3", "circle", "thr"}},
	{{"1/1/5/5", "quarter2", "sec"}, {"1/1/3/3", "circle", "ink"}},
	{{"1/1/5/2", "sq", "thr"}, {"1/2/5/3", "sq", "sec"}, {"1/3/5/4", "sq", "blk"}, {"2/3/4/5", "circle", "ink"}},
	{{"1/1/5/5", "rtri", "ink"}, {"3/3/5/5", "circle", "sec"}, {"1/4/2/5", "dot", "thr"}},
}

// bauhausMark renders composition n (cycled) for a section anchor.
// Decorative only: hidden from assistive technology.
func bauhausMark(n int) g.Node {
	shapes := bhMarks[n%len(bhMarks)]
	nodes := make([]g.Node, len(shapes))
	for i, s := range shapes {
		nodes[i] = h.I(h.Class("bh-"+s.kind), h.Style(fmt.Sprintf("grid-area:%s;background:var(--bh-%s);z-index:%d", s.area, s.color, i+1)))
	}
	return h.Div(h.Class("bh"), g.Attr("aria-hidden", "true"), g.Group(nodes))
}

const bauhausCSS = `
/* colours: a role that matches the section's own field colour is swapped
   for one that contrasts with it */
.bh { --bh-ink: var(--ink); --bh-bg: var(--bg); --bh-blk: var(--block); --bh-sec: var(--second); --bh-thr: var(--third);
  display: grid; grid-template: repeat(4, minmax(0, 1fr)) / repeat(4, minmax(0, 1fr));
  width: min(100%, 10.5rem); aspect-ratio: 1; margin-top: auto;
  /* own stacking context below the sticky anchor head, so the head always
     scrolls over the mark and the shapes keep their order */
  position: relative; z-index: 0; isolation: isolate; }
/* the sticky head slides over the mark when scrolling: opaque, in the
   anchor's own colour, so the title stays readable on top */
[data-component="swiss-section"] > [data-slot="anchor"] [data-slot="head"] { z-index: 1; background: var(--bg); padding-bottom: .5rem; }
[data-component="swiss-section"][data-field] > [data-slot="anchor"] [data-slot="head"] { background: inherit; }
.bh i { position: relative; }
[data-field="block"] .bh { --bh-blk: var(--bg); }
[data-field="second"] .bh { --bh-sec: var(--block); }
[data-field="third"] .bh { --bh-thr: var(--block); }
[data-field="ink"] .bh { --bh-ink: var(--bg); --bh-bg: var(--block); }
.bh-circle { border-radius: 50%; }
.bh-dot { border-radius: 50%; clip-path: circle(32%); }
.bh-ring { border-radius: 50%; -webkit-mask: radial-gradient(circle, transparent 46%, #000 47%); mask: radial-gradient(circle, transparent 46%, #000 47%); }
.bh-half { border-radius: 999px 999px 0 0; }
.bh-halfl { border-radius: 999px 0 0 999px; }
.bh-halfr { border-radius: 0 999px 999px 0; }
.bh-quarter { border-top-left-radius: 100%; }
.bh-quarter2 { border-bottom-right-radius: 100%; }
.bh-tri { clip-path: polygon(50% 0, 100% 100%, 0 100%); }
.bh-rtri { clip-path: polygon(0 0, 100% 0, 0 100%); }
.bh-stripes { -webkit-mask: repeating-linear-gradient(90deg, #000 0 22%, transparent 22% 39%); mask: repeating-linear-gradient(90deg, #000 0 22%, transparent 22% 39%); }
[data-component="swiss-section"][data-wide] .bh { display: none; }
@media (max-width: 1024px) { .bh { display: none; } }
@media (prefers-reduced-motion: no-preference) {
  .bh i { transition: transform .5s cubic-bezier(.2, .8, .2, 1); }
  [data-component="swiss-section"]:hover .bh i:nth-child(2) { transform: rotate(90deg); }
}
`
