package feedback

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type LoadingOverlayProps struct {
	// SignalName is the Datastar signal (default "_loading").
	// Set to true to show: data-on:click="$_loading=true"
	SignalName string
	Text       string // optional loading label
}

// LoadingOverlay renders a full-surface spinner overlay driven by a Datastar signal.
// Place inside a position:relative container or set Absolute to cover the whole viewport.
func LoadingOverlay(p LoadingOverlayProps) g.Node {
	if p.SignalName == "" {
		p.SignalName = "_loading"
	}
	sig := p.SignalName

	inner := []g.Node{
		Spinner(SpinnerProps{}),
	}
	if p.Text != "" {
		inner = append(inner, h.Span(
			h.Style("font-size:var(--t-sm);font-weight:600;color:var(--fg)"),
			g.Text(p.Text),
		))
	}

	return h.Div(
		g.Attr("data-component", "loading-overlay"),
		g.Attr("data-signals", `{"`+sig+`":false}`),
		g.Attr("data-show", "$"+sig),
		h.Style("display:none"),
		h.Div(g.Attr("data-slot", "backdrop")),
		h.Div(
			g.Attr("data-slot", "content"),
			g.Group(inner),
		),
	)
}
