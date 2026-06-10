package form

import (
	"fmt"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type DateInputProps struct {
	Name  string
	ID    string
	Value string // ISO date e.g. "2025-01-15"
	Min   string
	Max   string
}

// DateInput renders a styled native date picker.
func DateInput(p DateInputProps) g.Node {
	attrs := []g.Node{
		g.Attr("data-component", "input"),
		h.Type("date"),
		h.Name(p.Name),
	}
	if p.ID != "" {
		attrs = append(attrs, h.ID(p.ID))
	}
	if p.Value != "" {
		attrs = append(attrs, h.Value(p.Value))
	}
	if p.Min != "" {
		attrs = append(attrs, g.Attr("min", p.Min))
	}
	if p.Max != "" {
		attrs = append(attrs, g.Attr("max", p.Max))
	}
	return h.Input(attrs...)
}

type TimeInputProps struct {
	Name  string
	ID    string
	Value string // e.g. "14:30"
	Min   string
	Max   string
	Step  int // seconds
}

// TimeInput renders a styled native time picker.
func TimeInput(p TimeInputProps) g.Node {
	attrs := []g.Node{
		g.Attr("data-component", "input"),
		h.Type("time"),
		h.Name(p.Name),
	}
	if p.ID != "" {
		attrs = append(attrs, h.ID(p.ID))
	}
	if p.Value != "" {
		attrs = append(attrs, h.Value(p.Value))
	}
	if p.Min != "" {
		attrs = append(attrs, g.Attr("min", p.Min))
	}
	if p.Max != "" {
		attrs = append(attrs, g.Attr("max", p.Max))
	}
	if p.Step > 0 {
		attrs = append(attrs, g.Attr("step", fmt.Sprintf("%d", p.Step)))
	}
	return h.Input(attrs...)
}
