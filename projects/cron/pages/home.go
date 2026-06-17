package pages

import (
	"fmt"
	"time"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"mljr-web/ui"
	"mljr-web/ui/icon"
	"mljr-web/ui/layout"
	"mljr-web/ui/primitive"
	"mljr-web/ui/token"
)

const defaultExpr = "*/5 * * * *"

func Home() g.Node {
	initial := EvalCron(CronInput{Expression: defaultExpr, Count: 10})

	return layout.PageShell(
		layout.PageProps{
			Title:       "Cron Explorer — cron expression tester",
			Description: "Parse and preview cron expressions live. See next execution times, field breakdown, and human-readable descriptions.",
			Theme:       token.ThemeSwissBrut,
			Mode:        token.ModeLight,
			HeadExtra: []g.Node{
				g.El("style", g.Raw(cronCSS)),
			},
		},
		h.Main(
			h.Style("min-height:100vh;padding:var(--sp-8) 0 var(--sp-12)"),

			layout.Container(layout.ContainerProps{Attrs: []g.Node{h.Style("padding-left:clamp(.75rem,2vw,1.5rem);padding-right:clamp(.75rem,2vw,1.5rem)")}},

				// Header
				h.Div(
					h.Style("display:flex;align-items:baseline;gap:var(--sp-4);margin-bottom:var(--sp-6);padding-bottom:var(--sp-4);border-bottom:var(--bw-2) solid var(--ink)"),
					h.H1(
						h.Style("font-size:clamp(2.2rem,5vw,3.6rem);font-weight:900;line-height:1;margin:0;letter-spacing:-.03em"),
						g.Text("CRON EXPLORER"),
					),
					h.Div(
						h.Style("display:flex;gap:var(--sp-2);margin-left:auto;flex-shrink:0"),
						primitive.Tag(primitive.TagProps{Tone: token.ToneCyan}, g.Text("CRON")),
						primitive.Tag(primitive.TagProps{Tone: token.ToneLime},
							icon.Icon("simple-icons:go", icon.Props{Size: ".9rem"}),
							g.Text("Go"),
						),
					),
				),

				// Signal scope
				h.Div(
					ui.Signals(fmt.Sprintf(`{expr:%q,withSec:false,example:""}`, defaultExpr)),
					ui.On("input__debounce.200ms", "@post('/api/eval')"),

					// Examples row
					h.Div(
						h.Style("margin-bottom:var(--sp-5);display:flex;align-items:center;gap:var(--sp-3);flex-wrap:wrap"),
						h.Label(
							h.Style("font-size:var(--t-xs);font-weight:900;text-transform:uppercase;letter-spacing:.12em;opacity:.65;white-space:nowrap"),
							g.Text("EXAMPLES"),
						),
						h.Select(
							h.Class("cr-select"),
							ui.Bind("example"),
							ui.On("change__debounce.0ms", "@post('/api/example')"),
							h.Option(h.Value(""), g.Text("— choose an example —")),
							g.Group(func() []g.Node {
								nodes := make([]g.Node, 0, len(ExampleGroups))
								for _, grp := range ExampleGroups {
									opts := make([]g.Node, len(grp.Items))
									for i, ex := range grp.Items {
										opts[i] = h.Option(h.Value(ex.Key), g.Text(ex.Name))
									}
									nodes = append(nodes, g.El("optgroup",
										g.Attr("label", grp.Name),
										g.Group(opts),
									))
								}
								return nodes
							}()),
						),
					),

					// Input row
					h.Div(
						h.Style("margin-bottom:var(--sp-5)"),
						h.Label(
							h.Style("display:block;font-size:var(--t-xs);font-weight:900;text-transform:uppercase;letter-spacing:.12em;margin-bottom:var(--sp-2);opacity:.65"),
							g.Text("EXPRESSION"),
						),
						h.Div(
							h.Style("display:flex;align-items:stretch;gap:var(--sp-3);flex-wrap:wrap"),
							h.Div(
								h.Class("cr-input-wrap"),
								h.Input(
									h.Type("text"),
									h.Placeholder("* * * * *"),
									h.Class("cr-input"),
									ui.Bind("expr"),
									g.Attr("autocomplete", "off"),
									g.Attr("autocorrect", "off"),
									g.Attr("autocapitalize", "off"),
									g.Attr("spellcheck", "false"),
								),
							),
							h.Button(
								h.Type("button"),
								h.Title("Enable 6-field cron (with seconds as first field)"),
								h.Class("cr-flag"),
								ui.On("click", "$withSec=!$withSec; @post('/api/eval')"),
								g.Attr("data-class", `{"cr-flag--active":$withSec}`),
								g.Text("6-field"),
							),
						),
						// Field labels guide
						h.Div(
							h.Class("cr-field-guide"),
							g.Attr("data-class", `{"cr-field-guide--sec":$withSec}`),
							fieldLabel("sec", "Second", true),
							fieldLabel("min", "Minute", false),
							fieldLabel("hr", "Hour", false),
							fieldLabel("dom", "Day (month)", false),
							fieldLabel("mon", "Month", false),
							fieldLabel("dow", "Day (week)", false),
						),
					),

					// Output
					h.Div(
						h.Class("cr-grid"),
						// Left: field breakdown
						primitive.Card(primitive.CardProps{
							Tone:  token.ToneNone,
							Attrs: []g.Node{h.Style("padding:0;overflow:hidden")},
						},
							OutputFragment(initial),
						),
						// Right: cheat sheet
						referencePanel(),
					),
				),

				// Note
				h.P(
					h.Style("margin-top:var(--sp-8);font-size:var(--t-xs);opacity:.45;font-family:var(--font-mono,monospace)"),
					g.Text("Standard 5-field (min hr dom mon dow) or 6-field with leading second. Supports */step, ranges a-b, lists a,b,c, @aliases."),
				),
			),
		),
	)
}

func fieldLabel(abbr, full string, secOnly bool) g.Node {
	cls := "cr-fl-item"
	if secOnly {
		cls += " cr-fl-sec"
	}
	return h.Div(
		h.Class(cls),
		h.Span(h.Class("cr-fl-abbr"), g.Text(abbr)),
		h.Span(h.Class("cr-fl-full"), g.Text(full)),
	)
}

func referencePanel() g.Node {
	type row [2]string
	sec := func(title string, rows []row) g.Node {
		items := make([]g.Node, len(rows))
		for i, r := range rows {
			items[i] = h.Div(
				h.Style("display:grid;grid-template-columns:auto 1fr;gap:2px var(--sp-2);align-items:baseline;padding:3px 0;border-bottom:1px solid var(--line)"),
				h.Code(h.Style("font-family:var(--font-mono,monospace);font-size:11px;font-weight:900;white-space:nowrap;color:var(--ink)"), g.Text(r[0])),
				h.Span(h.Style("font-size:11px;opacity:.7;line-height:1.4"), g.Text(r[1])),
			)
		}
		return h.Div(
			h.Style("margin-bottom:var(--sp-4)"),
			h.Div(
				h.Style("font-size:10px;font-weight:900;text-transform:uppercase;letter-spacing:.12em;opacity:.45;margin-bottom:var(--sp-1)"),
				g.Text(title),
			),
			g.Group(items),
		)
	}

	return primitive.Card(primitive.CardProps{
		Tone:  token.ToneNone,
		Attrs: []g.Node{h.Style("padding:var(--sp-4);height:100%;box-sizing:border-box")},
	},
		h.Div(
			h.Style("font-size:var(--t-xs);font-weight:900;text-transform:uppercase;letter-spacing:.12em;margin-bottom:var(--sp-4);padding-bottom:var(--sp-2);border-bottom:var(--bw-2) solid var(--ink)"),
			g.Text("QUICK REFERENCE"),
		),
		sec("Fields (5-field)", []row{
			{"min", "0–59"},
			{"hour", "0–23"},
			{"dom", "1–31 (day of month)"},
			{"month", "1–12 or JAN–DEC"},
			{"dow", "0–7 (0&7=Sun) or SUN–SAT"},
		}),
		sec("Special Values", []row{
			{"*", "every value"},
			{"*/n", "every n-th value"},
			{"a-b", "range from a to b"},
			{"a,b,c", "list of values"},
			{"?", "any (dom/dow only)"},
		}),
		sec("@ Shortcuts", []row{
			{"@yearly", "0 0 1 1 *"},
			{"@monthly", "0 0 1 * *"},
			{"@weekly", "0 0 * * 0"},
			{"@daily", "0 0 * * *"},
			{"@hourly", "0 * * * *"},
			{"@every 5m", "interval (e.g. 5m, 1h30m)"},
		}),
		sec("6-Field (with seconds)", []row{
			{"sec", "0–59 (prepend to 5-field)"},
			{"*/10 * * * * *", "every 10 seconds"},
		}),
	)
}

func OutputFragment(r CronResult) g.Node {
	if r.Err != "" {
		return h.Div(
			h.ID("cr-output"),
			h.Style("padding:var(--sp-4)"),
			h.Div(
				h.Style("display:flex;align-items:center;gap:var(--sp-2);font-family:var(--font-mono,monospace);font-size:var(--t-sm);color:var(--blush-ink,#9f1239);background:var(--blush-bg,#fde8e8);border:var(--bw-2) solid var(--ink);padding:var(--sp-3)"),
				icon.Icon("lucide:alert-circle", icon.Props{Size: "1rem"}),
				h.Span(g.Text(r.Err)),
			),
		)
	}

	kicker := "font-size:var(--t-xs);font-weight:900;text-transform:uppercase;letter-spacing:.12em;opacity:.5"

	var sections []g.Node

	// Header bar
	sections = append(sections,
		h.Div(
			h.Style("display:flex;align-items:center;gap:var(--sp-3);padding:var(--sp-3) var(--sp-4);background:var(--ink);color:var(--bg)"),
			h.Span(h.Style("font-size:var(--t-sm);font-weight:900;font-family:var(--font-mono,monospace)"),
				g.Text("■ "+r.Expression),
			),
		),
	)

	// Human-readable
	if r.Human != "" {
		sections = append(sections,
			h.Div(
				h.Style("padding:var(--sp-3) var(--sp-4);border-bottom:var(--bw-1) solid var(--line)"),
				h.Span(h.Style("font-size:var(--t-base);font-weight:700"), g.Text(r.Human)),
			),
		)
	}

	// Field breakdown
	if len(r.Fields) > 0 {
		fieldNodes := make([]g.Node, len(r.Fields))
		fieldTones := []token.Tone{token.ToneCyan, token.ToneLime, token.ToneViolet, token.ToneYellow, token.TonePink, token.ToneMint}
		for i, f := range r.Fields {
			tone := fieldTones[i%len(fieldTones)]
			fieldNodes[i] = primitive.Tag(
				primitive.TagProps{Tone: tone},
				h.Span(h.Style("opacity:.55;margin-right:3px"), g.Text(f.Name+":")),
				h.Code(h.Style("font-family:var(--font-mono,monospace);font-weight:900"), g.Text(f.Value)),
				g.If(f.Desc != "" && f.Desc != "at "+f.Value,
					h.Span(h.Style("opacity:.6;margin-left:4px;font-size:.9em"), g.Text("("+f.Desc+")")),
				),
			)
		}
		sections = append(sections,
			h.Div(
				h.Style("padding:var(--sp-3) var(--sp-4);border-bottom:var(--bw-1) solid var(--line)"),
				h.Div(h.Style(kicker+";margin-bottom:var(--sp-2)"), g.Text("FIELDS")),
				h.Div(h.Style("display:flex;flex-wrap:wrap;gap:var(--sp-2)"), g.Group(fieldNodes)),
			),
		)
	}

	// Next executions
	if len(r.Next) > 0 {
		rows := make([]g.Node, len(r.Next))
		now := time.Now()
		for i, t := range r.Next {
			diff := t.Sub(now)
			var diffStr string
			if diff < time.Minute {
				diffStr = fmt.Sprintf("%ds", int(diff.Seconds()))
			} else if diff < time.Hour {
				diffStr = fmt.Sprintf("%dm %ds", int(diff.Minutes()), int(diff.Seconds())%60)
			} else if diff < 24*time.Hour {
				diffStr = fmt.Sprintf("%dh %dm", int(diff.Hours()), int(diff.Minutes())%60)
			} else {
				diffStr = fmt.Sprintf("%dd %dh", int(diff.Hours()/24), int(diff.Hours())%24)
			}
			rows[i] = h.Div(
				h.Style("display:grid;grid-template-columns:20px 1fr auto;gap:var(--sp-2);align-items:center;padding:var(--sp-2) var(--sp-3);border-bottom:var(--bw-1) solid var(--line)"),
				h.Div(h.Style("font-size:var(--t-xs);font-weight:900;opacity:.35"), g.Text(fmt.Sprintf("%d", i+1))),
				h.Code(h.Style("font-family:var(--font-mono,monospace);font-size:var(--t-sm)"),
					g.Text(t.Format("2006-01-02 15:04:05")),
				),
				h.Div(h.Style("font-size:var(--t-xs);opacity:.45;white-space:nowrap;font-family:var(--font-mono,monospace)"),
					g.Text("in "+diffStr),
				),
			)
		}
		sections = append(sections,
			h.Div(
				h.Style("border-top:var(--bw-2) solid var(--ink)"),
				h.Div(h.Style(kicker+";padding:var(--sp-2) var(--sp-4)"), g.Text("NEXT 10 EXECUTIONS")),
				h.Div(g.Group(rows)),
			),
		)
	}

	return h.Div(h.ID("cr-output"), g.Group(sections))
}

const cronCSS = `
.cr-grid {
  display: grid;
  grid-template-columns: 1fr 220px;
  gap: var(--sp-6);
  align-items: start;
}
@media (max-width: 900px) {
  .cr-grid { grid-template-columns: 1fr; }
}

.cr-input-wrap {
  display: flex;
  align-items: center;
  flex: 1;
  border: var(--bw-2) solid var(--ink);
  background: var(--bg);
  box-shadow: var(--shadow-sm);
  min-width: 200px;
}
.cr-input {
  flex: 1;
  min-width: 0;
  border: none;
  outline: none;
  background: transparent;
  font-family: var(--font-mono, monospace);
  font-size: var(--t-lg);
  font-weight: 900;
  color: var(--ink);
  letter-spacing: .08em;
  padding: var(--sp-2) var(--sp-3);
}
.cr-input-wrap:focus-within {
  outline: 3px solid var(--accent);
  outline-offset: 1px;
}

.cr-flag {
  padding: var(--sp-1) var(--sp-3);
  border: var(--bw-2) solid var(--ink);
  background: var(--bg);
  color: var(--ink);
  font-family: var(--font-mono, monospace);
  font-weight: 900;
  font-size: var(--t-sm);
  cursor: pointer;
  user-select: none;
  transition: background .1s, color .1s;
}
.cr-flag:hover { background: var(--surface-2, var(--surface)); }
.cr-flag--active { background: var(--ink); color: var(--bg); }

.cr-select {
  flex: 1;
  max-width: 400px;
  border: var(--bw-2) solid var(--ink);
  background: var(--bg);
  color: var(--ink);
  font-family: var(--font-mono, monospace);
  font-size: var(--t-sm);
  padding: var(--sp-1) var(--sp-3);
  cursor: pointer;
  outline: none;
  box-shadow: var(--shadow-sm);
  appearance: none;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='8' viewBox='0 0 12 8'%3E%3Cpath d='M1 1l5 5 5-5' stroke='%23000' stroke-width='2' fill='none'/%3E%3C/svg%3E");
  background-repeat: no-repeat;
  background-position: right var(--sp-3) center;
  padding-right: var(--sp-7);
}
.cr-select:focus { outline: 3px solid var(--accent); outline-offset: 1px; }

/* Field guide */
.cr-field-guide {
  display: flex;
  gap: var(--sp-3);
  margin-top: var(--sp-2);
  flex-wrap: wrap;
}
.cr-fl-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
  min-width: 48px;
}
.cr-fl-abbr {
  font-family: var(--font-mono, monospace);
  font-weight: 900;
  font-size: var(--t-sm);
}
.cr-fl-full {
  font-size: 10px;
  opacity: .5;
  white-space: nowrap;
}
/* Second field hidden unless 6-field mode */
.cr-fl-sec { display: none; }
.cr-field-guide--sec .cr-fl-sec { display: flex; }
`
