package pages

import (
	"fmt"
	"time"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"mljr-web/ui"
	"mljr-web/ui/special"
	"mljr-web/ui/swiss"
)

const defaultExpr = "*/5 * * * *"

func Home() g.Node {
	initial := EvalCron(CronInput{Expression: defaultExpr, Count: 10})

	return special.ToolShell(special.ToolShellProps{
		ID:          "cron",
		Name:        "Cron Explorer",
		Title:       "Cron Explorer — cron expression tester",
		Description: "Parse and preview cron expressions live. See next execution times, field breakdown, and human-readable descriptions.",
		CSS:         cronCSS,
		Sections: []swiss.NavItem{
			{ID: "tool", Num: "00", Label: "Werkzeug", Sub: "Tool"},
			{ID: "bench", Num: "01", Label: "Werkbank", Sub: "Workbench"},
			{ID: "reference", Num: "02", Label: "Referenz", Sub: "Reference"},
		},
	},
		swiss.Section(swiss.SectionProps{ID: "tool", Num: "00", Title: "Werkzeug", Sub: "Tool", Field: "block",
			Note: "Standard 5-field or 6-field with leading second. Supports */step, ranges, lists and @aliases."},
			swiss.Cell(swiss.CellProps{},
				swiss.ToolHead("Cron Explorer", "Type a cron expression and see what it means, field by field, and when it runs next.",
					"Cron", "@aliases", "Go"),
			),
		),
		h.Div(
			ui.Signals(fmt.Sprintf(`{expr:%q,withSec:false,example:""}`, defaultExpr)),
			ui.On("input__debounce.200ms", "@post('/api/eval')"),
			swiss.Section(swiss.SectionProps{ID: "bench", Num: "01", Title: "Werkbank", Sub: "Workbench",
				Note: "Run times are in the server's time zone, named above the list."},
				swiss.Cell(swiss.CellProps{Span: 3},
					swiss.Field("cr-example", "A. Example",
						h.Select(h.ID("cr-example"),
							ui.Bind("example"),
							ui.On("change__debounce.0ms", "@post('/api/example')"),
							h.Option(h.Value(""), g.Text("— choose an example —")),
							g.Group(exampleOptions()),
						),
					),
				),
				swiss.Cell(swiss.CellProps{Span: 6},
					swiss.Field("cr-expr", "B. Expression",
						h.Input(h.Type("text"), h.ID("cr-expr"), h.Class("cr-input"), h.Placeholder("* * * * *"), ui.Bind("expr"),
							g.Attr("autocomplete", "off"), g.Attr("autocorrect", "off"), g.Attr("autocapitalize", "off"), g.Attr("spellcheck", "false")),
						swiss.Choice("Format",
							swiss.Btn(false, []g.Node{
								h.Title("Enable 6-field cron (with seconds as first field)"),
								g.Attr("aria-pressed", "false"),
								g.Attr("data-attr:aria-pressed", `$withSec ? "true" : "false"`),
								ui.On("click", "$withSec=!$withSec; @post('/api/eval')"),
							}, g.Text("6-field")),
						),
					),
					h.Div(h.Class("cr-guide"), g.Attr("aria-hidden", "true"),
						g.Attr("data-class", `{"cr-guide--sec":$withSec}`),
						fieldLabel("sec", "Second", true),
						fieldLabel("min", "Minute", false),
						fieldLabel("hr", "Hour", false),
						fieldLabel("dom", "Day (month)", false),
						fieldLabel("mon", "Month", false),
						fieldLabel("dow", "Day (week)", false),
					),
				),
				swiss.Cell(swiss.CellProps{Variant: "flush"},
					OutputFragment(initial),
				),
			),
		),
		swiss.Section(swiss.SectionProps{ID: "reference", Num: "02", Title: "Referenz", Sub: "Reference",
			Note: "Field ranges, special values and shortcuts."},
			referenceCells()...,
		),
	)
}

func exampleOptions() []g.Node {
	nodes := make([]g.Node, 0, len(ExampleGroups))
	for _, grp := range ExampleGroups {
		opts := make([]g.Node, len(grp.Items))
		for i, ex := range grp.Items {
			opts[i] = h.Option(h.Value(ex.Key), g.Text(ex.Name))
		}
		nodes = append(nodes, g.El("optgroup", g.Attr("label", grp.Name), g.Group(opts)))
	}
	return nodes
}

func fieldLabel(abbr, full string, secOnly bool) g.Node {
	cls := "cr-fl"
	if secOnly {
		cls += " cr-fl-sec"
	}
	return h.Div(h.Class(cls), h.B(g.Text(abbr)), h.Span(g.Text(full)))
}

// referenceCells renders the static cron cheat sheet, one cell per topic.
func referenceCells() []g.Node {
	cell := func(span int, title string, rows [][2]string) g.Node {
		return swiss.Cell(swiss.CellProps{Span: span}, swiss.RefList(title, rows))
	}
	return []g.Node{
		cell(3, "Fields (5-field)", [][2]string{
			{"min", "0–59"},
			{"hour", "0–23"},
			{"dom", "1–31 (day of month)"},
			{"month", "1–12 or JAN–DEC"},
			{"dow", "0–7 (0&7=Sun) or SUN–SAT"},
		}),
		cell(3, "Special values", [][2]string{
			{"*", "every value"},
			{"*/n", "every n-th value"},
			{"a-b", "range from a to b"},
			{"a,b,c", "list of values"},
			{"?", "any (dom/dow only)"},
		}),
		cell(3, "@ Shortcuts", [][2]string{
			{"@yearly", "0 0 1 1 *"},
			{"@monthly", "0 0 1 * *"},
			{"@weekly", "0 0 * * 0"},
			{"@daily", "0 0 * * *"},
			{"@hourly", "0 * * * *"},
			{"@every 5m", "interval (e.g. 5m, 1h30m)"},
		}),
		cell(9, "6-field (with seconds)", [][2]string{
			{"sec", "0–59 (prepend to 5-field)"},
			{"*/10 * * * * *", "every 10 seconds"},
		}),
	}
}

// OutputFragment renders the result panel (id="cr-output"), for the first
// render and for every SSE patch.
func OutputFragment(r CronResult) g.Node {
	if r.Err != "" {
		return h.Div(h.ID("cr-output"), swiss.ResultBar("Error", nil), swiss.ErrorBox(r.Err))
	}

	var fields []g.Node
	for _, f := range r.Fields {
		fields = append(fields, h.Div(h.Class("cr-field"),
			swiss.Micro(g.Text(f.Name)),
			h.Code(g.Text(f.Value)),
			g.If(f.Desc != "" && f.Desc != "at "+f.Value, h.Span(g.Text(f.Desc))),
		))
	}

	var rows []g.Node
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
		rows = append(rows, swiss.Row(fmt.Sprint(i+1), h.Code(g.Text(t.Format("2006-01-02 15:04:05"))), g.Text("in "+diffStr)))
	}

	return h.Div(h.ID("cr-output"), h.Class("cr-out"),
		swiss.ResultBar(r.Expression, g.Text(fmt.Sprintf("%d fields", len(r.Fields)))),
		g.If(r.Human != "", h.P(h.Class("cr-human"), g.Text(r.Human))),
		h.Div(h.Class("cr-out-grid"),
			g.If(len(fields) > 0, swiss.ResultBlock("Fields", h.Div(h.Class("cr-fields"), g.Group(fields)))),
			g.If(len(rows) > 0, swiss.ResultBlock(fmt.Sprintf("Next %d executions · %s", len(rows), now.Format("MST")), g.Group(rows))),
		),
	)
}

const cronCSS = `
.cr-input { font-size: 1.35rem !important; font-weight: 700; letter-spacing: .08em; }
.cr-guide { display: flex; flex-wrap: wrap; gap: .25rem 1rem; font-family: var(--font-mono); font-size: .68rem; }
.cr-fl { display: flex; flex-direction: column; min-width: 3.2rem; }
.cr-fl span { color: var(--muted); }
.cr-fl-sec { display: none; }
.cr-guide--sec .cr-fl-sec { display: flex; }
.cr-human { margin: 0; padding: var(--swiss-pad); font-size: clamp(1.2rem, 1rem + .8vw, 1.8rem); font-weight: 800; letter-spacing: -.02em; border-bottom: 1px solid var(--ink); }
.cr-out-grid { display: grid; grid-template-columns: minmax(0, 1fr) minmax(0, 1.4fr); }
.cr-out-grid > * + * { box-shadow: -1px 0 0 var(--ink); }
.cr-out-grid > [data-component="swiss-result-block"] { border-bottom: 0; }
.cr-fields { display: grid; grid-template-columns: repeat(auto-fill, minmax(8rem, 1fr)); border-top: 1px solid var(--ink); border-left: 1px solid var(--ink); }
.cr-field { display: flex; flex-direction: column; gap: .2rem; padding: .55rem .65rem; border-right: 1px solid var(--ink); border-bottom: 1px solid var(--ink); }
.cr-field code { font-family: var(--font-mono); font-size: 1.05rem; font-weight: 700; }
.cr-field span { font-size: .75rem; color: var(--muted); }
@media (max-width: 1024px) {
  .cr-out-grid { grid-template-columns: minmax(0, 1fr); }
  .cr-out-grid > * + * { box-shadow: 0 -1px 0 var(--ink); }
}
`
