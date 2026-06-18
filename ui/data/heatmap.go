package data

import (
	"fmt"
	stdhtml "html"
	"math"
	"strings"
	"time"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type HeatmapCell struct {
	Date  time.Time
	Value int    // contribution/activity count
	Label string // tooltip label (optional)
}

type HeatmapProps struct {
	// Weeks is the number of weeks to display (default 52).
	Weeks int
	// CellSize in px (default 13).
	CellSize int
	// Gap between cells in px (default 3).
	Gap int
	// MaxValue for color intensity scaling (0 = auto from data).
	MaxValue int
	// ShowMonthLabels renders month abbreviations above grid.
	ShowMonthLabels bool
	// ShowDayLabels renders Mon/Wed/Fri on the left.
	ShowDayLabels bool
	// Now overrides the reference "today" (default time.Now()); exposed
	// for tests so grid-alignment behavior can be checked across weekdays.
	Now time.Time
}

// Heatmap renders a GitHub-style contribution heatmap as an SVG.
// All rendering is server-side — zero JS. Pass cells for each day.
func Heatmap(p HeatmapProps, cells []HeatmapCell) g.Node {
	if p.Weeks == 0 {
		p.Weeks = 52
	}
	if p.CellSize == 0 {
		p.CellSize = 13
	}
	if p.Gap == 0 {
		p.Gap = 3
	}

	// Build lookup map date→value
	lookup := make(map[string]HeatmapCell, len(cells))
	maxVal := p.MaxValue
	for _, c := range cells {
		key := c.Date.Format("2006-01-02")
		lookup[key] = c
		if p.MaxValue == 0 && c.Value > maxVal {
			maxVal = c.Value
		}
	}
	if maxVal == 0 {
		maxVal = 1
	}

	// Anchor the grid on the Saturday on/after today so the current
	// (possibly incomplete) week is always included, then walk back
	// Weeks full Sun–Sat weeks. Aligning only the start (and not the
	// end) would shift the whole window earlier on any day but
	// Saturday, silently dropping the most recent days.
	now := p.Now
	if now.IsZero() {
		now = time.Now()
	}
	end := now
	for end.Weekday() != time.Saturday {
		end = end.AddDate(0, 0, 1)
	}
	start := end.AddDate(0, 0, -(p.Weeks*7 - 1))

	step := p.CellSize + p.Gap
	leftPad := 0
	if p.ShowDayLabels {
		leftPad = 28
	}
	topPad := 0
	if p.ShowMonthLabels {
		topPad = 20
	}

	svgW := p.Weeks*step + leftPad
	svgH := 7*step + topPad

	var sb strings.Builder
	fmt.Fprintf(&sb, `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %d %d" style="width:100%%;height:auto;max-width:%dpx" data-component="heatmap">`,
		svgW, svgH, svgW)

	// Day labels
	if p.ShowDayLabels {
		days := []struct{ d, row int }{{1, 1}, {3, 3}, {5, 5}} // Mon, Wed, Fri
		for _, dl := range days {
			y := topPad + dl.row*step + p.CellSize/2 + 4
			fmt.Fprintf(&sb, `<text x="0" y="%d" fill="var(--muted)" font-size="9" font-family="var(--font-mono)">%s</text>`,
				y, time.Weekday(dl.d).String()[:3])
		}
	}

	// Month labels and cells
	lastMonth := -1
	for w := 0; w < p.Weeks; w++ {
		for d := 0; d < 7; d++ {
			date := start.AddDate(0, 0, w*7+d)
			key := date.Format("2006-01-02")
			cell := lookup[key]

			x := leftPad + w*step
			y := topPad + d*step

			// Month label
			if p.ShowMonthLabels && date.Day() <= 7 && date.Month() != time.Month(lastMonth) {
				lastMonth = int(date.Month())
				fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="var(--muted)" font-size="10" font-family="var(--font-mono)">%s</text>`,
					x, topPad-6, date.Month().String()[:3])
			}

			intensity := 0.0
			if cell.Value > 0 {
				// sqrt compresses the high end so a few outlier days don't
				// flatten every other day into the same top bucket.
				intensity = math.Sqrt(float64(cell.Value) / float64(maxVal))
			}
			color := heatmapColor(intensity)

			label := cell.Label
			if label == "" && cell.Value > 0 {
				label = fmt.Sprintf("%d on %s", cell.Value, key)
			}

			fmt.Fprintf(&sb, `<rect x="%d" y="%d" width="%d" height="%d" rx="2" fill="%s">`,
				x, y, p.CellSize, p.CellSize, color)
			if label != "" {
				fmt.Fprintf(&sb, `<title>%s</title>`, stdhtml.EscapeString(label))
			}
			sb.WriteString(`</rect>`)
		}
	}
	sb.WriteString(`</svg>`)

	return h.Div(
		g.Attr("data-component", "heatmap-wrap"),
		h.Style("overflow-x:auto"),
		g.Raw(sb.String()),
	)
}

// heatmapColor returns a CSS color for the given 0–1 intensity.
// Uses CSS vars so it themes correctly.
func heatmapColor(intensity float64) string {
	if intensity <= 0 {
		return "var(--surface-2)"
	}
	// 4 levels
	switch {
	case intensity < 0.25:
		return "color-mix(in srgb, var(--accent) 30%, var(--surface-2))"
	case intensity < 0.5:
		return "color-mix(in srgb, var(--accent) 55%, var(--surface-2))"
	case intensity < 0.75:
		return "color-mix(in srgb, var(--accent) 80%, var(--surface-2))"
	default:
		return "var(--accent)"
	}
}
