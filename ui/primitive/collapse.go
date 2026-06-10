package primitive

import (
	"fmt"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type CollapseProps struct {
	// SignalName is the Datastar signal controlling open state (default "_collapseOpen").
	// Use a unique name per instance when multiple Collapses appear on one page.
	SignalName string
	Open       bool   // initial open state
	Duration   string // CSS transition duration (default "0.25s")
}

// Collapse renders an animated show/hide region driven by a Datastar signal.
// Wrap any content — it collapses via max-height transition.
// Toggle with: data-on:click="$signalName=!$signalName"
func Collapse(p CollapseProps, children ...g.Node) g.Node {
	if p.SignalName == "" {
		p.SignalName = "_collapseOpen"
	}
	if p.Duration == "" {
		p.Duration = "0.25s"
	}
	openVal := "false"
	if p.Open {
		openVal = "true"
	}
	sig := p.SignalName

	// max-height trick: open=max-height:600px, closed=max-height:0+overflow:hidden
	heightExpr := fmt.Sprintf(`{"style":'max-height:'+($%s?'600px':'0')+';overflow:hidden;transition:max-height %s ease'}`, sig, p.Duration)

	return h.Div(
		g.Attr("data-component", "collapse"),
		g.Attr("data-signals", `{"`+sig+`":`+openVal+`}`),
		g.Attr("data-attr", heightExpr),
		h.Style("max-height:"+func() string {
			if p.Open {
				return "600px"
			}
			return "0"
		}()+";overflow:hidden;transition:max-height "+p.Duration+" ease"),
		g.Group(children),
	)
}
