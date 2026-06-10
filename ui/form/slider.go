package form

import (
	"fmt"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type SliderProps struct {
	Signal    string
	Name      string
	Min       int
	Max       int
	Step      int
	Label     string
	ShowValue bool
	Attrs     []g.Node
}

// Slider renders a styled HTML range input bound to a Datastar signal.
// ShowValue adds a floating label above the thumb tracking the current value.
func Slider(p SliderProps, attrs ...g.Node) g.Node {
	if p.Max == 0 {
		p.Max = 100
	}
	if p.Step == 0 {
		p.Step = 1
	}
	sig := p.Signal

	var labelNode g.Node
	if p.Label != "" {
		labelNode = h.Div(g.Attr("data-slot", "label"), g.Text(p.Label))
	}

	var thumbLabel g.Node
	if p.ShowValue && sig != "" {
		// pct = 0..100 proportion along the track
		expr := fmt.Sprintf("($%s-%d)/(%d-%d)*100", sig, p.Min, p.Max, p.Min)
		// clamp so label doesn't overflow at track edges
		dataAttr := fmt.Sprintf(`{"style":"left:clamp(1rem,"+(%s).toFixed(1)+"%%,calc(100%% - 1rem))"}`, expr)
		thumbLabel = h.Div(
			g.Attr("data-slot", "thumb-label"),
			g.Attr("data-text", "$"+sig),
			g.Attr("data-attr", dataAttr),
		)
	}

	return h.Div(
		g.Attr("data-component", "slider-wrap"),
		g.Group(p.Attrs),
		labelNode,
		h.Div(
			g.Attr("data-component", "slider-track"),
			thumbLabel,
			h.Input(
				h.Type("range"),
				g.Attr("data-component", "slider"),
				g.If(sig != "", g.Attr("data-bind:"+sig)),
				g.If(sig != "", g.Attr("data-on:input", "$"+sig+"=Number(evt.target.value)")),
				g.If(p.Name != "", h.Name(p.Name)),
				h.Min(fmt.Sprintf("%d", p.Min)),
				h.Max(fmt.Sprintf("%d", p.Max)),
				h.Step(fmt.Sprintf("%d", p.Step)),
				g.Group(attrs),
			),
		),
	)
}
