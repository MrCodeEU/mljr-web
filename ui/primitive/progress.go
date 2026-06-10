package primitive

import (
	"fmt"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type ProgressVariant string

const (
	ProgressDefault ProgressVariant = ""
	ProgressSuccess ProgressVariant = "success"
	ProgressWarning ProgressVariant = "warning"
	ProgressDanger  ProgressVariant = "danger"
)

type ProgressProps struct {
	Value      int    // 0–100
	Label      string // optional left label
	ValueLabel string // optional right label (defaults to "N%")
	Variant    ProgressVariant
	Size       string // sm | md (default) | lg
	ShowLabel  bool
	Attrs      []g.Node
}

// Progress renders a horizontal progress bar.
func Progress(p ProgressProps) g.Node {
	if p.Value < 0 {
		p.Value = 0
	}
	if p.Value > 100 {
		p.Value = 100
	}
	valLabel := p.ValueLabel
	if valLabel == "" {
		valLabel = fmt.Sprintf("%d%%", p.Value)
	}

	var labelRow g.Node
	if p.ShowLabel || p.Label != "" {
		labelRow = h.Div(
			g.Attr("data-slot", "label"),
			g.If(p.Label != "", h.Span(g.Text(p.Label))),
			h.Span(g.Text(valLabel)),
		)
	}

	return h.Div(
		g.Attr("data-component", "progress-wrap"),
		g.Group(p.Attrs),
		labelRow,
		h.Div(
			g.Attr("data-component", "progress"),
			g.If(p.Variant != "", g.Attr("data-variant", string(p.Variant))),
			g.If(p.Size != "", g.Attr("data-size", p.Size)),
			g.Attr("role", "progressbar"),
			g.Attr("aria-valuenow", fmt.Sprintf("%d", p.Value)),
			g.Attr("aria-valuemin", "0"),
			g.Attr("aria-valuemax", "100"),
			h.Div(
				g.Attr("data-slot", "fill"),
				h.Style(fmt.Sprintf("width:%d%%", p.Value)),
			),
		),
	)
}
