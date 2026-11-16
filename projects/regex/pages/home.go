package pages

import (
	"fmt"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"mljr-web/ui"
	"mljr-web/ui/special"
	"mljr-web/ui/swiss"
)

const defaultPattern = `(quick|lazy)`
const defaultInput = "The quick brown fox\njumps over the lazy dog\nPack my box with five dozen liquor jugs"

func Home() g.Node {
	initial := EvalRegex(EvalInput{
		Pattern: defaultPattern,
		Input:   defaultInput,
	})

	return special.ToolShell(special.ToolShellProps{
		ID:          "regex",
		Name:        "Regex Lab",
		Title:       "Regex Lab — live PCRE tester",
		Description: "Live regular expression tester powered by Go + regexp2. Supports backreferences, lookahead, lookbehind. No page reloads.",
		CSS:         regexCSS,
		Sections: []swiss.NavItem{
			{ID: "tool", Num: "00", Label: "Werkzeug", Sub: "Tool"},
			{ID: "bench", Num: "01", Label: "Werkbank", Sub: "Workbench"},
			{ID: "reference", Num: "02", Label: "Referenz", Sub: "Reference"},
		},
	},
		swiss.Section(swiss.SectionProps{ID: "tool", Num: "00", Title: "Werkzeug", Sub: "Tool", Field: "block",
			Note: "PCRE-compatible engine (regexp2): backreferences, lookaround, atomic groups, Unicode."},
			swiss.Cell(swiss.CellProps{},
				swiss.ToolHead("Regex Lab", "Test regular expressions live: matches are highlighted as you type, with capture groups and a replace preview.",
					"PCRE", "regexp2", "Go"),
			),
		),
		h.Div(
			ui.Signals(fmt.Sprintf(`{pattern:%q,flagI:false,flagM:false,flagS:false,input:%q,replace:"",example:""}`,
				defaultPattern, defaultInput)),
			ui.On("input__debounce.200ms", "@post('/api/eval')"),
			swiss.Section(swiss.SectionProps{ID: "bench", Num: "01", Title: "Werkbank", Sub: "Workbench",
				Note: "Flags: i = case-insensitive, m = multiline, s = dot matches newline."},
				swiss.Cell(swiss.CellProps{Span: 3},
					swiss.Field("rx-example", "A. Example",
						h.Select(h.ID("rx-example"),
							ui.Bind("example"),
							ui.On("change__debounce.0ms", "@post('/api/example')"),
							h.Option(h.Value(""), g.Text("— choose an example —")),
							g.Group(exampleOptions()),
						),
					),
				),
				swiss.Cell(swiss.CellProps{Span: 6},
					swiss.Field("rx-pattern", "B. Pattern",
						h.Div(h.Class("rx-pattern"),
							swiss.Affix("/"),
							h.Input(h.Type("text"), h.ID("rx-pattern"), h.Placeholder("pattern…"), ui.Bind("pattern"),
								g.Attr("autocomplete", "off"), g.Attr("autocorrect", "off"), g.Attr("autocapitalize", "off"), g.Attr("spellcheck", "false")),
							swiss.Affix("/"),
						),
						swiss.Choice("Flags",
							flagBtn("i", "flagI", "Case insensitive"),
							flagBtn("m", "flagM", "Multiline (^ $ match line boundaries)"),
							flagBtn("s", "flagS", "Dot matches newline"),
						),
					),
				),
				swiss.Cell(swiss.CellProps{Span: 4},
					swiss.Field("rx-input", "C. Test string",
						h.Textarea(h.ID("rx-input"), g.Attr("rows", "9"), g.Attr("spellcheck", "false"), ui.Bind("input")),
					),
					swiss.Field("rx-replace", "D. Replace",
						h.Input(h.Type("text"), h.ID("rx-replace"), h.Placeholder("$1 groups · ${name} named"), ui.Bind("replace"),
							g.Attr("autocomplete", "off"), g.Attr("spellcheck", "false")),
					),
				),
				swiss.Cell(swiss.CellProps{Span: 5, Variant: "flush"},
					OutputFragment(initial),
				),
			),
		),
		swiss.Section(swiss.SectionProps{ID: "reference", Num: "02", Title: "Referenz", Sub: "Reference",
			Note: "The syntax regexp2 understands, at a glance."},
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

// referenceCells renders the static PCRE cheat sheet, one cell per topic.
func referenceCells() []g.Node {
	cell := func(title string, rows [][2]string) g.Node {
		return swiss.Cell(swiss.CellProps{Span: 3}, swiss.RefList(title, rows))
	}
	return []g.Node{
		cell("Anchors", [][2]string{
			{`^`, "start of line / string"},
			{`$`, "end of line / string"},
			{`\b`, "word boundary"},
			{`\B`, "non-word boundary"},
		}),
		cell("Characters", [][2]string{
			{`.`, "any char (s flag: incl. \\n)"},
			{`\d / \D`, "digit / non-digit"},
			{`\w / \W`, "word char / non-word"},
			{`\s / \S`, "whitespace / non-ws"},
			{`[abc]`, "character class"},
			{`[^abc]`, "negated class"},
			{`[a-z]`, "range"},
		}),
		cell("Quantifiers", [][2]string{
			{`*  +  ?`, "0+, 1+, 0-1 (greedy)"},
			{`{n}`, "exactly n"},
			{`{n,m}`, "n to m (greedy)"},
			{`*? +? ??`, "lazy (non-greedy)"},
		}),
		cell("Groups", [][2]string{
			{`(abc)`, "capture group"},
			{`(?:abc)`, "non-capturing"},
			{`(?P<n>abc)`, "named capture"},
		}),
		cell("Lookaround · PCRE", [][2]string{
			{`(?=abc)`, "positive lookahead"},
			{`(?!abc)`, "negative lookahead"},
			{`(?<=abc)`, "positive lookbehind"},
			{`(?<!abc)`, "negative lookbehind"},
		}),
		cell("Backrefs · PCRE", [][2]string{
			{`\1 \2`, "backref in pattern"},
			{`$1 $2`, "backref in replace"},
			{`${name}`, "named group in replace"},
		}),
	}
}

// OutputFragment renders the live result panel (id="output").
// Called both for SSR initial render and SSE patch.
func OutputFragment(r EvalResult) g.Node {
	if r.Err != "" {
		return h.Div(h.ID("output"), swiss.ResultBar("Error", nil), swiss.ErrorBox(r.Err))
	}
	if r.PatternEmpty {
		return h.Div(h.ID("output"),
			swiss.ResultBar("Type a pattern to see matches", nil),
			swiss.ResultBlock("Test string", h.Pre(g.Raw(r.Highlighted))),
		)
	}

	matchLabel := fmt.Sprintf("%d match", r.MatchCount)
	if r.MatchCount != 1 {
		matchLabel = fmt.Sprintf("%d matches", r.MatchCount)
	}
	rows := make([]g.Node, 0, len(r.Matches))
	for _, m := range r.Matches {
		val := m.Value
		if len(val) > 80 {
			val = val[:77] + "…"
		}
		var groups []g.Node
		for gi, grp := range m.Groups {
			gval := grp
			if len(gval) > 60 {
				gval = gval[:57] + "…"
			}
			groups = append(groups, h.Div(g.Text(fmt.Sprintf("↳ group %d: %q", gi+1, gval))))
		}
		rows = append(rows, swiss.Row(fmt.Sprint(m.Index), h.Code(g.Text(val)), g.Text(fmt.Sprintf("%d–%d", m.Start, m.End)), groups...))
	}

	return h.Div(h.ID("output"),
		swiss.ResultBar("Result", g.Text(matchLabel)),
		swiss.ResultBlock("Highlighted", h.Pre(g.Raw(r.Highlighted))),
		g.If(len(rows) > 0, swiss.ResultBlock("Matches", g.Group(rows))),
		g.If(r.ReplaceApplied, swiss.ResultBlock("Replace output", h.Pre(g.Text(r.Replaced)))),
	)
}

func flagBtn(label, signal, title string) g.Node {
	return swiss.Btn(false, []g.Node{
		h.Title(title),
		g.Attr("aria-label", title),
		g.Attr("aria-pressed", "false"),
		g.Attr("data-attr:aria-pressed", "$"+signal+` ? "true" : "false"`),
		ui.On("click", fmt.Sprintf("$%s=!$%s; @post('/api/eval')", signal, signal)),
	}, g.Text(label))
}

const regexCSS = `
.rx-pattern { display: flex; flex: 1 1 16rem; min-width: 0; }
.rx-pattern input { font-size: 1.05rem !important; }
mark.rx-match { background: var(--block); color: var(--block-fg); padding: 0 1px; }
`
