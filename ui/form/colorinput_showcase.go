//go:build showcase

package form

import (
	"mljr-web/ui/registry"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "color-input", Name: "Color Input", Category: "form",
		Summary: "Styled native color picker with optional hex label. Datastar signal tracks the selected value.",
		Code: `form.ColorInput(form.ColorInputProps{
    Name:    "brand",
    Value:   "#f28d1d",
    ShowHex: true,
})`,
		Render: func(p map[string]string) g.Node {
			return h.Div(
				h.Style("display:flex;flex-direction:column;gap:var(--sp-5)"),
				Field(FieldProps{Label: "Brand color"},
					ColorInput(ColorInputProps{Name: "brand", Value: "#f28d1d", ShowHex: true}),
				),
				Field(FieldProps{Label: "Accent color"},
					ColorInput(ColorInputProps{Name: "accent", Value: "#a40054", ShowHex: true}),
				),
				h.Div(
					h.P(h.Style("font-size:var(--t-xs);opacity:.5;font-weight:700;text-transform:uppercase;letter-spacing:.06em;margin:0 0 var(--sp-3)"), g.Text("Palette")),
					h.Div(
						h.Style("display:flex;gap:var(--sp-3);flex-wrap:wrap"),
						ColorInput(ColorInputProps{Name: "c1", Value: "#ff6b6b"}),
						ColorInput(ColorInputProps{Name: "c2", Value: "#ffd93d"}),
						ColorInput(ColorInputProps{Name: "c3", Value: "#6bcb77"}),
						ColorInput(ColorInputProps{Name: "c4", Value: "#4d96ff"}),
					),
				),
			)
		},
	})

	registry.Register(&registry.Component{
		Slug: "range-pair", Name: "Range Pair", Category: "form",
		Summary: "Min/max dual-handle range selector. Two overlapping range inputs with Datastar signals enforcing low ≤ high.",
		Code: `form.RangePair(form.RangePairProps{
    Name:    "price",
    Min:     0,
    Max:     1000,
    LowVal:  200,
    HighVal: 800,
})`,
		Render: func(p map[string]string) g.Node {
			return h.Div(
				h.Style("display:flex;flex-direction:column;gap:var(--sp-5);max-width:400px"),
				Field(FieldProps{Label: "Price range (€)", Hint: "Filter by price"},
					RangePair(RangePairProps{Name: "price", Min: 0, Max: 1000, LowVal: 200, HighVal: 800}),
				),
				Field(FieldProps{Label: "Year range"},
					RangePair(RangePairProps{Name: "year", Min: 2000, Max: 2026, LowVal: 2018, HighVal: 2025, Step: 1}),
				),
			)
		},
	})
}
