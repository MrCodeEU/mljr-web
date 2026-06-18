package data

import (
	"strings"
	"testing"
	"time"
)

// Regression test for the bug where the grid's start was aligned back to
// Sunday without compensating the end date, silently dropping the most
// recent days on any day but Saturday. Confirms today's date always shows
// up with its label in the rendered SVG, regardless of today's weekday.
func TestHeatmapIncludesToday(t *testing.T) {
	for weekday := time.Sunday; weekday <= time.Saturday; weekday++ {
		now := nextWeekday(time.Now(), weekday)
		cells := []HeatmapCell{
			{Date: now, Value: 7, Label: "today"},
		}
		node := Heatmap(HeatmapProps{Weeks: 52, Now: now}, cells)
		var sb strings.Builder
		if err := node.Render(&sb); err != nil {
			t.Fatalf("render: %v", err)
		}
		if !strings.Contains(sb.String(), "<title>today</title>") {
			t.Errorf("weekday %s: today's cell label missing from rendered heatmap", weekday)
		}
	}
}

func nextWeekday(from time.Time, day time.Weekday) time.Time {
	for from.Weekday() != day {
		from = from.AddDate(0, 0, 1)
	}
	return from
}
