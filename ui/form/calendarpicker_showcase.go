//go:build showcase

package form

import (
	"mljr-web/ui/registry"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "calendar-picker", Name: "Calendar Picker", Category: "form",
		Summary: "Custom styled date picker. Month grid with keyboard-friendly day selection. Submits ISO 8601 value via hidden input. No external library.",
		Code: `form.CalendarPicker(form.CalendarPickerProps{
    Name:        "start_date",
    Value:       "2026-06-10",
    Min:         "2026-01-01",
    Placeholder: "Pick a date…",
})`,
		Render: func(p map[string]string) g.Node {
			return h.Div(
				h.Style("display:flex;flex-direction:column;gap:var(--sp-5)"),
				h.Div(
					h.Label(h.Style("display:block;font-weight:700;margin-bottom:var(--sp-2)"), g.Text("Start date")),
					CalendarPicker(CalendarPickerProps{
						Name:        "start_date",
						Value:       "2026-06-15",
						Placeholder: "Select start date…",
						Signal:      "_cal1",
					}),
				),
				h.Div(
					h.Label(h.Style("display:block;font-weight:700;margin-bottom:var(--sp-2)"), g.Text("Restricted range (June 2026 only)")),
					CalendarPicker(CalendarPickerProps{
						Name:        "event_date",
						Min:         "2026-06-01",
						Max:         "2026-06-30",
						Placeholder: "June dates only…",
						Signal:      "_cal2",
					}),
				),
			)
		},
	})
}
