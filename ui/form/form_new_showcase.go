//go:build showcase

package form

import (
	"mljr-web/ui"
	"mljr-web/ui/icon"
	"mljr-web/ui/registry"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "number-input", Name: "Number Input", Category: "form",
		Summary: "Stepper with − and + buttons. Min/max/step enforced via Datastar expressions.",
		Code: `// declare signal in ancestor
ui.Signals("{count:1}")

form.NumberInput(form.NumberInputProps{
    Signal: "count",
    Min:    0,
    Max:    99,
    Step:   1,
})`,
		Controls: []registry.Control{
			{Name: "min", Type: registry.ControlEnum, Options: []string{"0", "-10"}, Default: "0"},
			{Name: "max", Type: registry.ControlEnum, Options: []string{"10", "100", "999"}, Default: "10"},
			{Name: "step", Type: registry.ControlEnum, Options: []string{"1", "5", "10"}, Default: "1"},
		},
		Render: func(p map[string]string) g.Node {
			return h.Div(
				ui.Signals(`{"numVal":1}`),
				Field(FieldProps{Label: "Quantity"},
					NumberInput(NumberInputProps{
						Signal: "numVal",
						Min:    0,
						Max:    10,
						Step:   1,
						Value:  1,
					}),
				),
				h.P(h.Style("font-size:var(--t-sm);margin-top:var(--sp-2)"),
					g.Text("Value: "), h.Span(g.Attr("data-text", "$numVal")),
				),
			)
		},
	})

	registry.Register(&registry.Component{
		Slug: "slider", Name: "Slider", Category: "form",
		Summary: "Styled range input with Datastar binding and optional value label.",
		Code: `form.Slider(form.SliderProps{
    Signal:    "vol",
    Label:     "Volume",
    Min:       0,
    Max:       100,
    Step:      1,
    ShowValue: true,
})`,
		Controls: []registry.Control{
			{Name: "label", Type: registry.ControlBool, Default: "true"},
		},
		Render: func(p map[string]string) g.Node {
			return h.Div(
				ui.Signals(`{"sliderVal":60}`),
				Slider(SliderProps{
					Signal:    "sliderVal",
					Label:     "Volume",
					ShowValue: p["label"] == "true",
					Min:       0,
					Max:       100,
					Step:      1,
				}),
			)
		},
	})

	registry.Register(&registry.Component{
		Slug: "input-group", Name: "Input Group", Category: "form",
		Summary: "Input with prefix/suffix addons — icons, text, or buttons.",
		Code: `form.InputGroup(form.InputGroupProps{
    Prefix: icon.Icon("lucide:search"),
},
    form.Input(form.InputProps{Type: "search", Placeholder: "Search…", Signal: "q"}),
)`,
		Render: func(p map[string]string) g.Node {
			return h.Div(
				ui.Signals(`{"searchQ":"","urlQ":""}`),
				h.Div(h.Style("display:flex;flex-direction:column;gap:var(--sp-4)"),
					Field(FieldProps{Label: "Search"},
						InputGroup(InputGroupProps{Prefix: icon.Icon("lucide:search")},
							Input(InputProps{Type: "search", Placeholder: "Search…", Signal: "searchQ"}),
						),
					),
					Field(FieldProps{Label: "URL"},
						InputGroup(InputGroupProps{
							Prefix: h.Span(g.Text("https://")),
							Suffix: h.Span(g.Text(".com")),
						},
							Input(InputProps{Placeholder: "example", Signal: "urlQ"}),
						),
					),
				),
			)
		},
	})
}
