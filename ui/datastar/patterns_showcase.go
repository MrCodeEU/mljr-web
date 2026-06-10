//go:build showcase

package datastar

import (
	"mljr-web/ui"
	"mljr-web/ui/form"
	"mljr-web/ui/layout"
	"mljr-web/ui/primitive"
	"mljr-web/ui/registry"
	"mljr-web/ui/token"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "ds-patterns", Name: "Patterns", Category: "datastar",
		Summary: "Real-world patterns: client-side search filter, multi-step wizard, optimistic toggle.",
		Code: `// Pattern 1: client-side search filter
// Filter a list by signal value — no server needed
data-show="$q===''||item.toLowerCase().includes($q.toLowerCase())"

// Pattern 2: wizard / multi-step
data-show="$step===1"   // show step 1 panel
data-on:click="$step=Math.min($totalSteps,$step+1)"  // next
data-show="$step>1"     // show back button

// Pattern 3: optimistic toggle
// 1. Flip UI immediately
data-on:click="$liked=!$liked"
// 2. Confirm or revert via SSE
data-on:click="$liked=!$liked;@post('/api/like')"
// Server can revert: MarshalAndPatchSignals(sse, {"liked": false})`,
		Render: func(p map[string]string) g.Node {
			fruits := []string{
				"Apple", "Apricot", "Avocado", "Banana", "Blueberry",
				"Cherry", "Coconut", "Grape", "Guava", "Kiwi",
				"Lemon", "Lime", "Lychee", "Mango", "Melon",
				"Orange", "Papaya", "Peach", "Pear", "Pineapple",
				"Plum", "Pomegranate", "Raspberry", "Strawberry", "Watermelon",
			}

			fruitItems := make([]g.Node, len(fruits))
			for i, f := range fruits {
				fruit := f
				fruitItems[i] = h.Div(
					h.Style("padding:var(--sp-2) var(--sp-3);border-radius:var(--radius);font-size:var(--t-sm)"),
					// show if query is empty OR fruit contains query (case-insensitive)
					g.Attr("data-show", `$_q===''||'`+fruit+`'.toLowerCase().includes($_q.toLowerCase())`),
					h.Style("display:none"),
					g.Text(fruit),
				)
			}

			return h.Div(
				ui.Signals(`{"_q":"","step":1,"liked":false,"likesCount":142,"wizardName":"","wizardPlan":"free","_wizardDone":false}`),
				h.Div(h.Style("display:flex;flex-direction:column;gap:var(--sp-5)"),

					// Client-side search filter
					primitive.Card(primitive.CardProps{},
						h.P(h.Style("font-size:var(--t-xs);text-transform:uppercase;font-weight:700;letter-spacing:.06em;opacity:.5;margin-bottom:var(--sp-3)"), g.Text("Pattern 1 — client-side search filter")),
						form.Input(form.InputProps{Signal: "_q", Placeholder: "Filter fruits…"}),
						h.Div(
							h.Style("margin-top:var(--sp-3);display:grid;grid-template-columns:repeat(auto-fill,minmax(120px,1fr));gap:var(--sp-1);max-height:200px;overflow-y:auto"),
							g.Group(fruitItems),
						),
						h.P(h.Style("font-size:var(--t-xs);opacity:.5;margin-top:var(--sp-2)"),
							g.Text("25 items, zero server requests. Each item has data-show that checks signal against its own text."),
						),
					),

					// Multi-step wizard
					primitive.Card(primitive.CardProps{},
						h.P(h.Style("font-size:var(--t-xs);text-transform:uppercase;font-weight:700;letter-spacing:.06em;opacity:.5;margin-bottom:var(--sp-3)"), g.Text("Pattern 2 — multi-step wizard")),
						h.Div(
							g.Attr("data-show", "!$_wizardDone"),
							h.Style("display:none"),
							// Progress indicator
							layout.Stepper(layout.StepperProps{},
								layout.Step{Label: "Account", State: layout.StepComplete},
								layout.Step{Label: "Plan"},
								layout.Step{Label: "Done"},
							),

							// Step 1
							h.Div(
								g.Attr("data-show", "$step===1"),
								h.Style("display:none;margin-top:var(--sp-4)"),
								form.Field(form.FieldProps{Label: "Your name"},
									form.Input(form.InputProps{Signal: "wizardName", Placeholder: "e.g. Alice"}),
								),
								h.Div(h.Style("margin-top:var(--sp-4)"),
									primitive.Button(primitive.ButtonProps{Variant: token.Primary},
										g.Attr("data-on:click", "$step=2"),
										g.Attr("data-attr", `{"disabled":$wizardName===''}`),
										g.Text("Next →"),
									),
								),
							),
							// Step 2
							h.Div(
								g.Attr("data-show", "$step===2"),
								h.Style("display:none;margin-top:var(--sp-4)"),
								form.Field(form.FieldProps{Label: "Choose plan"},
									form.RadioGroup(form.RadioGroupProps{
										Signal: "wizardPlan", Name: "wplan",
										Options: []form.RadioOption{
											{Value: "free", Label: "Free"},
											{Value: "pro", Label: "Pro"},
											{Value: "enterprise", Label: "Enterprise"},
										},
									}),
								),
								h.Div(h.Style("display:flex;gap:var(--sp-2);margin-top:var(--sp-4)"),
									primitive.Button(primitive.ButtonProps{Variant: token.Ghost},
										g.Attr("data-on:click", "$step=1"), g.Text("← Back"),
									),
									primitive.Button(primitive.ButtonProps{Variant: token.Primary},
										g.Attr("data-on:click", "$step=3;$_wizardDone=true"),
										g.Text("Finish"),
									),
								),
							),
						),
						// Done state
						h.Div(
							g.Attr("data-show", "$_wizardDone"),
							h.Style("display:none;text-align:center;padding:var(--sp-5)"),
							h.Div(h.Style("font-size:2.5rem"), g.Text("🎉")),
							h.P(h.Style("font-weight:700;margin-top:var(--sp-2)"),
								g.Attr("data-text", "`Welcome, ${$wizardName}! Plan: ${$wizardPlan}`"),
							),
							primitive.Button(primitive.ButtonProps{Variant: token.Ghost, Size: token.SizeSM},
								g.Attr("data-on:click", "$step=1;$_wizardDone=false;$wizardName=''"),
								g.Text("Reset"),
							),
						),
					),

					// Optimistic toggle
					primitive.Card(primitive.CardProps{},
						h.P(h.Style("font-size:var(--t-xs);text-transform:uppercase;font-weight:700;letter-spacing:.06em;opacity:.5;margin-bottom:var(--sp-3)"), g.Text("Pattern 3 — optimistic UI toggle")),
						h.Div(h.Style("display:flex;align-items:center;gap:var(--sp-4)"),
							h.Button(
								g.Attr("data-component", "button"),
								g.Attr("data-attr", `{"data-variant":$liked?"primary":"outline"}`),
								g.Attr("data-on:click", "$liked=!$liked;$likesCount+=$liked?1:-1"),
								g.Attr("data-text", "`${$liked?'♥':'♡'} ${$likesCount}`"),
							),
							h.P(h.Style("font-size:var(--t-sm);opacity:.6"),
								g.Text("UI updates instantly (optimistic). In production, add "),
								h.Code(g.Text("@post('/api/like')")),
								g.Text(" — server confirms or reverts via "),
								h.Code(g.Text("PatchSignals")),
								g.Text("."),
							),
						),
					),
				),
			)
		},
	})
}
