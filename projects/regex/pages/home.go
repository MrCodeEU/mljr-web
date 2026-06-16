package pages

import (
	"fmt"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"mljr-web/ui"
	"mljr-web/ui/icon"
	"mljr-web/ui/layout"
	"mljr-web/ui/primitive"
	"mljr-web/ui/token"
)

const defaultPattern = `(quick|lazy)`
const defaultInput = "The quick brown fox\njumps over the lazy dog\nPack my box with five dozen liquor jugs"

func Home() g.Node {
	initial := EvalRegex(EvalInput{
		Pattern: defaultPattern,
		Input:   defaultInput,
	})

	return layout.PageShell(
		layout.PageProps{
			Title:       "Regex Lab — live PCRE tester",
			Description: "Live regular expression tester powered by Go + regexp2. Supports backreferences, lookahead, lookbehind. No page reloads.",
			Theme:       token.ThemeSwissBrut,
			Mode:        token.ModeLight,
			HeadExtra: []g.Node{
				g.El("style", g.Raw(regexCSS)),
			},
		},
		h.Main(
			h.Style("min-height:100vh;padding:var(--sp-8) 0 var(--sp-12)"),

			layout.Container(layout.ContainerProps{},

				// ── Header ────────────────────────────────────────────────────
				h.Div(
					h.Style("display:flex;align-items:baseline;gap:var(--sp-4);margin-bottom:var(--sp-6);padding-bottom:var(--sp-4);border-bottom:var(--bw-2) solid var(--ink)"),
					h.H1(
						h.Style("font-size:clamp(2.2rem,5vw,3.6rem);font-weight:900;line-height:1;margin:0;letter-spacing:-.03em"),
						g.Text("REGEX LAB"),
					),
					h.Div(
						h.Style("display:flex;gap:var(--sp-2);margin-left:auto;flex-shrink:0"),
						primitive.Tag(primitive.TagProps{Tone: token.ToneCyan}, g.Text("PCRE")),
						primitive.Tag(primitive.TagProps{Tone: token.ToneLime},
							icon.Icon("simple-icons:go", icon.Props{Size: ".9rem"}),
							g.Text("Go"),
						),
					),
				),

				// ── Signal scope ──────────────────────────────────────────────
				h.Div(
					ui.Signals(fmt.Sprintf(`{pattern:%q,flagI:false,flagM:false,flagS:false,input:%q,replace:"",example:"",showRef:false}`,
						defaultPattern, defaultInput)),
					ui.On("input__debounce.200ms", "@post('/api/eval')"),

					// Examples row
					h.Div(
						h.Style("margin-bottom:var(--sp-5);display:flex;align-items:center;gap:var(--sp-3);flex-wrap:wrap"),
						h.Label(
							h.Style("font-size:var(--t-xs);font-weight:900;text-transform:uppercase;letter-spacing:.12em;opacity:.65;white-space:nowrap"),
							g.Text("EXAMPLES"),
						),
						h.Select(
							h.Class("rx-select"),
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
						// Reference toggle (hidden on large screens via CSS)
						h.Button(
							h.Type("button"),
							h.Class("rx-ref-toggle"),
							ui.On("click", "$showRef=!$showRef"),
							g.Text("Quick Reference"),
							icon.Icon("lucide:chevron-down", icon.Props{Size: ".9rem"}),
						),
					),

					// ── Main layout: [reference | tool] ──────────────────────
					h.Div(
						h.Class("rx-main-layout"),

						// Reference panel
						h.Div(
							h.Class("rx-ref"),
							g.Attr("data-class", `{"rx-ref--closed":!$showRef}`),
							referencePanel(),
						),

						// Tool column
						h.Div(
							// Pattern row
							h.Div(
								h.Style("margin-bottom:var(--sp-4)"),
								h.Label(
									h.Style("display:block;font-size:var(--t-xs);font-weight:900;text-transform:uppercase;letter-spacing:.12em;margin-bottom:var(--sp-2);opacity:.65"),
									g.Text("PATTERN"),
								),
								h.Div(
									h.Style("display:flex;align-items:stretch;gap:var(--sp-3)"),
									h.Div(
										h.Class("rx-pattern-wrap"),
										h.Span(h.Class("rx-slash"), g.Text("/")),
										h.Input(
											h.Type("text"),
											h.Placeholder("pattern…"),
											h.Class("rx-pattern-input"),
											ui.Bind("pattern"),
											g.Attr("autocomplete", "off"),
											g.Attr("autocorrect", "off"),
											g.Attr("autocapitalize", "off"),
											g.Attr("spellcheck", "false"),
										),
										h.Span(h.Class("rx-slash"), g.Text("/")),
									),
									h.Div(
										h.Style("display:flex;gap:var(--sp-1);align-items:center"),
										flagBtn("i", "flagI", "Case insensitive"),
										flagBtn("m", "flagM", "Multiline (^ $ match line boundaries)"),
										flagBtn("s", "flagS", "Dot matches newline"),
									),
								),
							),

							// Inputs / output grid
							h.Div(
								h.Class("rx-grid"),
								// Left: test string + replace
								h.Div(
									h.Div(
										h.Style("margin-bottom:var(--sp-4)"),
										h.Label(
											h.Style("display:block;font-size:var(--t-xs);font-weight:900;text-transform:uppercase;letter-spacing:.12em;margin-bottom:var(--sp-2);opacity:.65"),
											g.Text("TEST STRING"),
										),
										h.Textarea(
											h.Class("rx-textarea"),
											g.Attr("rows", "9"),
											g.Attr("spellcheck", "false"),
											ui.Bind("input"),
										),
									),
									h.Div(
										h.Label(
											h.Style("display:block;font-size:var(--t-xs);font-weight:900;text-transform:uppercase;letter-spacing:.12em;margin-bottom:var(--sp-2);opacity:.65"),
											g.Text("REPLACE"),
										),
										h.Div(
											h.Class("rx-pattern-wrap"),
											h.Input(
												h.Type("text"),
												h.Placeholder("replacement… ($1 groups · ${name} named)"),
												h.Class("rx-pattern-input"),
												ui.Bind("replace"),
												g.Attr("autocomplete", "off"),
												g.Attr("spellcheck", "false"),
											),
										),
									),
								),
								// Right: output panel
								primitive.Card(primitive.CardProps{
									Tone:  token.ToneNone,
									Attrs: []g.Node{h.Style("padding:0;overflow:hidden;min-height:340px")},
								},
									OutputFragment(initial),
								),
							),
						),
					),
				),

				// Engine note
				h.P(
					h.Style("margin-top:var(--sp-8);font-size:var(--t-xs);opacity:.45;font-family:var(--font-mono,monospace)"),
					g.Text("PCRE-compatible engine (regexp2/v2) — backreferences, lookahead/lookbehind, atomic groups, Unicode. Flags: i=case-insensitive  m=multiline  s=dot-matches-newline."),
				),
			),
		),
	)
}

// referencePanel renders the static PCRE cheat sheet.
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
		sec("Anchors", []row{
			{`^`, "start of line / string"},
			{`$`, "end of line / string"},
			{`\b`, "word boundary"},
			{`\B`, "non-word boundary"},
		}),
		sec("Characters", []row{
			{`.`, "any char (s flag: incl. \\n)"},
			{`\d / \D`, "digit / non-digit"},
			{`\w / \W`, "word char / non-word"},
			{`\s / \S`, "whitespace / non-ws"},
			{`[abc]`, "character class"},
			{`[^abc]`, "negated class"},
			{`[a-z]`, "range"},
		}),
		sec("Quantifiers", []row{
			{`*  +  ?`, "0+, 1+, 0-1 (greedy)"},
			{`{n}`, "exactly n"},
			{`{n,m}`, "n to m (greedy)"},
			{`*? +? ??`, "lazy (non-greedy)"},
		}),
		sec("Groups", []row{
			{`(abc)`, "capture group"},
			{`(?:abc)`, "non-capturing"},
			{`(?P<n>abc)`, "named capture"},
		}),
		sec("Lookaround · PCRE", []row{
			{`(?=abc)`, "positive lookahead"},
			{`(?!abc)`, "negative lookahead"},
			{`(?<=abc)`, "positive lookbehind"},
			{`(?<!abc)`, "negative lookbehind"},
		}),
		sec("Backrefs · PCRE", []row{
			{`\1 \2`, "backref in pattern"},
			{`$1 $2`, "backref in replace"},
			{`${name}`, "named group in replace"},
		}),
	)
}

// OutputFragment renders the live result panel (id="output").
// Called both for SSR initial render and SSE patch.
func OutputFragment(r EvalResult) g.Node {
	kicker := "font-size:var(--t-xs);font-weight:900;text-transform:uppercase;letter-spacing:.12em;opacity:.5;margin-bottom:var(--sp-2)"

	if r.Err != "" {
		return h.Div(
			h.ID("output"),
			h.Style("padding:var(--sp-4)"),
			h.Div(
				h.Style("display:flex;align-items:center;gap:var(--sp-2);font-family:var(--font-mono,monospace);font-size:var(--t-sm);color:var(--blush-ink,#9f1239);background:var(--blush-bg,#fde8e8);border:var(--bw-2) solid var(--ink);padding:var(--sp-3)"),
				icon.Icon("lucide:alert-circle", icon.Props{Size: "1rem"}),
				h.Span(g.Text(r.Err)),
			),
		)
	}

	if r.PatternEmpty {
		return h.Div(
			h.ID("output"),
			h.Style("padding:var(--sp-4);display:flex;flex-direction:column;gap:var(--sp-4)"),
			h.Div(h.Style("font-size:var(--t-xs);font-weight:900;text-transform:uppercase;letter-spacing:.12em;opacity:.4"),
				g.Text("TYPE A PATTERN TO SEE MATCHES"),
			),
			h.Pre(
				h.Style("white-space:pre-wrap;font-family:var(--font-mono,monospace);font-size:var(--t-sm);line-height:1.75;margin:0;opacity:.5"),
				g.Raw(r.Highlighted),
			),
		)
	}

	matchLabel := fmt.Sprintf("%d match", r.MatchCount)
	if r.MatchCount != 1 {
		matchLabel = fmt.Sprintf("%d matches", r.MatchCount)
	}

	matchItems := make([]g.Node, 0, len(r.Matches))
	for _, m := range r.Matches {
		val := m.Value
		if len(val) > 80 {
			val = val[:77] + "…"
		}
		rowChildren := []g.Node{
			h.Div(h.Style("font-size:var(--t-xs);font-weight:900;opacity:.35;padding-top:3px"), g.Text(fmt.Sprintf("%d", m.Index))),
			h.Code(h.Style("font-family:var(--font-mono,monospace);font-size:var(--t-sm);word-break:break-all"), g.Text(val)),
			h.Div(h.Style("font-size:var(--t-xs);opacity:.45;white-space:nowrap;padding-top:3px"), g.Text(fmt.Sprintf("%d–%d", m.Start, m.End))),
		}
		for gi, grp := range m.Groups {
			gval := grp
			if len(gval) > 60 {
				gval = gval[:57] + "…"
			}
			rowChildren = append(rowChildren,
				h.Div(h.Style("grid-column:2/-1;font-size:var(--t-xs);opacity:.55;padding-left:var(--sp-3)"),
					g.Text(fmt.Sprintf("↳ group %d: %q", gi+1, gval)),
				),
			)
		}
		matchItems = append(matchItems,
			h.Div(
				h.Style("display:grid;grid-template-columns:20px 1fr auto;gap:var(--sp-2);align-items:start;padding:var(--sp-2) var(--sp-3);border-bottom:var(--bw-1) solid var(--line)"),
				g.Group(rowChildren),
			),
		)
	}

	children := []g.Node{
		h.Div(
			h.Style("display:flex;align-items:center;gap:var(--sp-3);padding:var(--sp-3) var(--sp-4);background:var(--ink);color:var(--bg)"),
			h.Span(h.Style("font-size:var(--t-sm);font-weight:900;font-family:var(--font-mono,monospace)"),
				g.Text("■ "+matchLabel),
			),
		),
		h.Div(
			h.Style("padding:var(--sp-4)"),
			h.Div(h.Style(kicker), g.Text("HIGHLIGHTED")),
			h.Pre(
				h.Style("white-space:pre-wrap;font-family:var(--font-mono,monospace);font-size:var(--t-sm);line-height:1.75;margin:0"),
				g.Raw(r.Highlighted),
			),
		),
		h.Div(
			h.Style("border-top:var(--bw-2) solid var(--ink)"),
			h.Div(h.Style(kicker+";padding:var(--sp-2) var(--sp-4)"), g.Text("MATCHES")),
			h.Div(g.Group(matchItems)),
		),
	}

	if r.ReplaceApplied {
		children = append(children,
			h.Div(
				h.Style("border-top:var(--bw-2) solid var(--ink)"),
				h.Div(h.Style(kicker+";padding:var(--sp-2) var(--sp-4)"), g.Text("REPLACE OUTPUT")),
				h.Div(
					h.Style("padding:var(--sp-3) var(--sp-4)"),
					h.Pre(
						h.Style("white-space:pre-wrap;font-family:var(--font-mono,monospace);font-size:var(--t-sm);line-height:1.75;margin:0"),
						g.Text(r.Replaced),
					),
				),
			),
		)
	}

	return h.Div(h.ID("output"), g.Group(children))
}

func flagBtn(label, signal, title string) g.Node {
	return h.Button(
		h.Type("button"),
		h.Title(title),
		h.Class("rx-flag"),
		ui.On("click", fmt.Sprintf("$%s=!$%s; @post('/api/eval')", signal, signal)),
		g.Attr("data-class", fmt.Sprintf(`{"rx-flag--active":$%s}`, signal)),
		g.Text(label),
	)
}

const regexCSS = `
/* ── Layout ─────────────────────────────────────────────────────── */
.rx-main-layout {
  display: grid;
  grid-template-columns: 220px 1fr;
  gap: var(--sp-6);
  align-items: start;
}
.rx-ref { display: block; }
.rx-ref--closed { display: none; }
.rx-ref-toggle { display: none; }

@media (max-width: 1050px) {
  .rx-main-layout { grid-template-columns: 1fr; }
  .rx-ref { display: none; }
  .rx-ref.rx-ref--closed { display: none; }
  .rx-ref:not(.rx-ref--closed) { display: block; }
  .rx-ref-toggle { display: inline-flex; align-items: center; gap: var(--sp-1); }
}

.rx-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--sp-5);
}
@media (max-width: 720px) {
  .rx-grid { grid-template-columns: 1fr; }
}

/* ── Pattern input ───────────────────────────────────────────────── */
.rx-pattern-wrap {
  display: flex;
  align-items: center;
  flex: 1;
  border: var(--bw-2) solid var(--ink);
  background: var(--bg);
  box-shadow: var(--shadow-sm);
}
.rx-slash {
  padding: var(--sp-2) var(--sp-3);
  font-family: var(--font-mono, monospace);
  font-size: var(--t-lg);
  font-weight: 900;
  opacity: .4;
  user-select: none;
  flex-shrink: 0;
}
.rx-pattern-input {
  flex: 1;
  min-width: 0;
  border: none;
  outline: none;
  background: transparent;
  font-family: var(--font-mono, monospace);
  font-size: var(--t-base);
  color: var(--ink);
  padding: var(--sp-2) 0;
}
.rx-pattern-wrap:focus-within {
  outline: 3px solid var(--accent);
  outline-offset: 1px;
}

/* ── Textarea ────────────────────────────────────────────────────── */
.rx-textarea {
  width: 100%;
  box-sizing: border-box;
  border: var(--bw-2) solid var(--ink);
  background: var(--bg);
  color: var(--ink);
  font-family: var(--font-mono, monospace);
  font-size: var(--t-sm);
  line-height: 1.7;
  padding: var(--sp-3);
  resize: vertical;
  outline: none;
  box-shadow: var(--shadow-sm);
}
.rx-textarea:focus {
  outline: 3px solid var(--accent);
  outline-offset: 1px;
}

/* ── Flag buttons ────────────────────────────────────────────────── */
.rx-flag {
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
.rx-flag:hover { background: var(--surface-2, var(--surface)); }
.rx-flag--active {
  background: var(--ink);
  color: var(--bg);
}

/* ── Examples select ─────────────────────────────────────────────── */
.rx-select {
  flex: 1;
  max-width: 360px;
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
.rx-select:focus { outline: 3px solid var(--accent); outline-offset: 1px; }

/* ── Reference toggle button ─────────────────────────────────────── */
.rx-ref-toggle {
  padding: var(--sp-1) var(--sp-3);
  border: var(--bw-2) solid var(--ink);
  background: var(--bg);
  color: var(--ink);
  font-family: var(--font-mono, monospace);
  font-size: var(--t-sm);
  font-weight: 900;
  cursor: pointer;
  gap: var(--sp-1);
  white-space: nowrap;
}

/* ── Match highlight ─────────────────────────────────────────────── */
mark.rx-match {
  background: var(--accent);
  color: var(--accent-ink, var(--ink));
  border-radius: 2px;
  padding: 0 1px;
}
`
