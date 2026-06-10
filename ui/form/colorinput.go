package form

import (
	"fmt"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type ColorInputProps struct {
	Name    string
	ID      string
	Value   string // hex color default (e.g. "#f28d1d")
	ShowHex bool   // show hex value label next to swatch
}

// ColorInput renders a styled native color picker with optional hex label.
func ColorInput(p ColorInputProps) g.Node {
	if p.Value == "" {
		p.Value = "#000000"
	}
	sig := "_cv"
	if p.Name != "" {
		sig += p.Name
	}

	inputAttrs := []g.Node{
		g.Attr("data-component", "color-input-native"),
		h.Type("color"),
		h.Name(p.Name),
		h.Value(p.Value),
		g.Attr("data-bind:"+sig),
	}
	if p.ID != "" {
		inputAttrs = append(inputAttrs, h.ID(p.ID))
	}

	nodes := []g.Node{
		g.Attr("data-component", "color-input"),
		g.Attr("data-signals", `{"`+sig+`":"`+p.Value+`"}`),
		h.Input(inputAttrs...),
	}

	if p.ShowHex {
		nodes = append(nodes,
			h.Span(
				g.Attr("data-slot", "hex"),
				g.Attr("data-text", "$"+sig),
				h.Style("font-family:var(--font-mono);font-size:var(--t-sm)"),
				g.Text(p.Value),
			),
		)
	}

	return h.Label(nodes...)
}

type RangePairProps struct {
	Name    string // base name; creates name_min and name_max
	Min     float64
	Max     float64
	Step    float64
	LowVal  float64 // initial low value
	HighVal float64 // initial high value
	Format  string  // JS format fn name, empty = raw number
}

// RangePair renders a min/max dual-handle range selector.
// Two range inputs overlap on the same track; JS ensures low ≤ high.
func RangePair(p RangePairProps) g.Node {
	if p.Max == 0 {
		p.Max = 100
	}
	if p.Step == 0 {
		p.Step = 1
	}
	if p.HighVal == 0 {
		p.HighVal = p.Max
	}

	minSig := "_rp" + p.Name + "lo"
	maxSig := "_rp" + p.Name + "hi"

	return h.Div(
		g.Attr("data-component", "range-pair"),
		g.Attr("data-signals", fmt.Sprintf(`{"%s":%v,"%s":%v}`, minSig, p.LowVal, maxSig, p.HighVal)),
		h.Div(
			g.Attr("data-slot", "track"),
			h.Input(
				h.Type("range"),
				h.Name(p.Name+"_min"),
				g.Attr("min", fmt.Sprintf("%v", p.Min)),
				g.Attr("max", fmt.Sprintf("%v", p.Max)),
				g.Attr("step", fmt.Sprintf("%v", p.Step)),
				g.Attr("data-bind:"+minSig),
				g.Attr("data-on:input__debounce.0ms", fmt.Sprintf("$%s=Math.min(Number(evt.target.value),$%s-$%s)", minSig, maxSig, "1")),
				h.Value(fmt.Sprintf("%v", p.LowVal)),
			),
			h.Input(
				h.Type("range"),
				h.Name(p.Name+"_max"),
				g.Attr("min", fmt.Sprintf("%v", p.Min)),
				g.Attr("max", fmt.Sprintf("%v", p.Max)),
				g.Attr("step", fmt.Sprintf("%v", p.Step)),
				g.Attr("data-bind:"+maxSig),
				g.Attr("data-on:input__debounce.0ms", fmt.Sprintf("$%s=Math.max(Number(evt.target.value),$%s+$%s)", maxSig, minSig, "1")),
				h.Value(fmt.Sprintf("%v", p.HighVal)),
			),
		),
		h.Div(
			g.Attr("data-slot", "values"),
			h.Span(g.Attr("data-text", "$"+minSig), g.Text(fmt.Sprintf("%v", p.LowVal))),
			h.Span(g.Text("–")),
			h.Span(g.Attr("data-text", "$"+maxSig), g.Text(fmt.Sprintf("%v", p.HighVal))),
		),
	)
}
