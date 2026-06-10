package layout

import (
	"fmt"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

// BackgroundPattern selects the decorative pattern style.
type BackgroundPattern string

const (
	BGDots     BackgroundPattern = "dots"
	BGGrid     BackgroundPattern = "grid"
	BGLines    BackgroundPattern = "lines"
	BGDiagonal BackgroundPattern = "diagonal"
	BGCross    BackgroundPattern = "cross"
	BGGradient BackgroundPattern = "gradient"
)

type BackgroundProps struct {
	Pattern BackgroundPattern
	Color   string  // CSS color for the pattern (default "var(--line)")
	Size    string  // CSS background-size (default varies per pattern)
	Opacity float64 // 0–1 (default 0.5)
	Fixed   bool    // position:fixed vs absolute (default absolute)
}

// Background renders a decorative full-coverage pattern layer.
// Place as the first child of a position:relative container.
// Pointer-events are disabled so it doesn't block interaction.
func Background(p BackgroundProps) g.Node {
	if p.Color == "" {
		p.Color = "var(--line)"
	}
	if p.Opacity == 0 {
		p.Opacity = 0.5
	}

	pos := "absolute"
	if p.Fixed {
		pos = "fixed"
	}

	bgStyle := bgCSS(p)
	return h.Div(
		g.Attr("data-component", "background"),
		g.Attr("data-pattern", string(p.Pattern)),
		h.Style(pos+":inset(0);inset:0;pointer-events:none;z-index:0;opacity:"+fmt.Sprintf("%.2f", p.Opacity)+";"+bgStyle),
		g.Attr("aria-hidden", "true"),
	)
}

func bgCSS(p BackgroundProps) string {
	color := p.Color
	size := p.Size

	switch p.Pattern {
	case BGDots:
		if size == "" {
			size = "24px 24px"
		}
		return "background-image:radial-gradient(" + color + " 1px,transparent 1px);background-size:" + size
	case BGGrid:
		if size == "" {
			size = "32px 32px"
		}
		return "background-image:linear-gradient(" + color + " 1px,transparent 1px),linear-gradient(90deg," + color + " 1px,transparent 1px);background-size:" + size
	case BGLines:
		if size == "" {
			size = "24px 24px"
		}
		return "background-image:linear-gradient(" + color + " 1px,transparent 1px);background-size:" + size
	case BGDiagonal:
		if size == "" {
			size = "20px 20px"
		}
		return "background-image:repeating-linear-gradient(45deg," + color + " 0," + color + " 1px,transparent 0,transparent 50%);background-size:" + size
	case BGCross:
		if size == "" {
			size = "32px 32px"
		}
		return "background-image:radial-gradient(circle," + color + " 2px,transparent 2px);background-size:" + size
	case BGGradient:
		return "background:radial-gradient(ellipse at top left," + color + " 0%,transparent 60%)"
	}
	return ""
}
