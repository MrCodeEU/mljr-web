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

// Home renders the full page with a pre-computed initial result.
func Home() g.Node {
	initial := EvalRegex(EvalInput{
		Pattern: defaultPattern,
		Input:   defaultInput,
	})

	return layout.PageShell(
		layout.PageProps{
			Title:       "Regex Lab — live RE2 tester",
			Description: "Live regular expression tester powered by Go's RE2 engine. No page reloads — matches highlight as you type.",
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
					h.Style("display:flex;align-items:baseline;gap:var(--sp-4);margin-bottom:var(--sp-8);padding-bottom:var(--sp-4);border-bottom:var(--bw-2) solid var(--ink)"),
					h.H1(
						h.Style("font-size:clamp(2.2rem,5vw,3.6rem);font-weight:900;line-height:1;margin:0;letter-spacing:-.03em"),
						g.Text("REGEX LAB"),
					),
					h.Div(
						h.Style("display:flex;gap:var(--sp-2);margin-left:auto"),
						primitive.Tag(primitive.TagProps{Tone: token.ToneCyan}, g.Text("RE2")),
						primitive.Tag(primitive.TagProps{Tone: token.ToneLime},
							icon.Icon("simple-icons:go", icon.Props{Size: ".9rem"}),
							g.Text("Go"),
						),
					),
				),

				// ── Tool ─────────────────────────────────────────────────────
				h.Div(
					// Signal scope + trigger wrapper
					ui.Signals(fmt.Sprintf(`{pattern:%q,flagI:false,flagM:false,flagS:false,input:%q,replace:""}`,
						defaultPattern, defaultInput)),
					ui.On("input__debounce.200ms", "@post('/api/eval')"),

					// Pattern row
					h.Div(
						h.Style("margin-bottom:var(--sp-4)"),
						h.Label(
							h.Style("display:block;font-size:var(--t-xs);font-weight:900;text-transform:uppercase;letter-spacing:.12em;margin-bottom:var(--sp-2);opacity:.65"),
							g.Text("PATTERN"),
						),
						h.Div(
							h.Style("display:flex;align-items:stretch;gap:var(--sp-3)"),
							// /pattern/ input
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
							// Flag toggles
							h.Div(
								h.Style("display:flex;gap:var(--sp-1);align-items:center"),
								flagBtn("i", "flagI", "Case insensitive"),
								flagBtn("m", "flagM", "Multiline (^ $ match line boundaries)"),
								flagBtn("s", "flagS", "Dot matches newline"),
							),
						),
					),

					// Two-column: inputs left, output right
					h.Div(
						h.Class("rx-grid"),
						// ── Left: test string + replace ───────────────────────
						h.Div(
							// Test string
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
							// Replace
							h.Div(
								h.Label(
									h.Style("display:block;font-size:var(--t-xs);font-weight:900;text-transform:uppercase;letter-spacing:.12em;margin-bottom:var(--sp-2);opacity:.65"),
									g.Text("REPLACE"),
								),
								h.Div(
									h.Class("rx-pattern-wrap"),
									h.Input(
										h.Type("text"),
										h.Placeholder("replacement… ($1 for groups)"),
										h.Class("rx-pattern-input"),
										ui.Bind("replace"),
										g.Attr("autocomplete", "off"),
										g.Attr("spellcheck", "false"),
									),
								),
							),
						),

						// ── Right: output panel ───────────────────────────────
						primitive.Card(primitive.CardProps{
							Tone:  token.ToneNone,
							Attrs: []g.Node{h.Style("padding:0;overflow:hidden;min-height:340px")},
						},
							OutputFragment(initial),
						),
					),
				),

				// RE2 note
				h.P(
					h.Style("margin-top:var(--sp-8);font-size:var(--t-xs);opacity:.45;font-family:var(--font-mono,monospace)"),
					g.Text("RE2 engine — no backreferences, lookahead, or lookbehind. Unicode supported. Flags: i=case-insensitive  m=multiline  s=dot-matches-newline."),
				),
			),
		),
	)
}

// OutputFragment renders the live result panel (id="output").
// Called both for SSR initial render and SSE patch.
func OutputFragment(r EvalResult) g.Node {
	kicker := "font-size:var(--t-xs);font-weight:900;text-transform:uppercase;letter-spacing:.12em;opacity:.5;margin-bottom:var(--sp-2)"

	// Error state
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

	// Empty pattern
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
		// Stats bar
		h.Div(
			h.Style("display:flex;align-items:center;gap:var(--sp-3);padding:var(--sp-3) var(--sp-4);background:var(--ink);color:var(--bg)"),
			h.Span(h.Style("font-size:var(--t-sm);font-weight:900;font-family:var(--font-mono,monospace)"),
				g.Text("■ "+matchLabel),
			),
		),
		// Highlighted test string
		h.Div(
			h.Style("padding:var(--sp-4)"),
			h.Div(h.Style(kicker), g.Text("HIGHLIGHTED")),
			h.Pre(
				h.Style("white-space:pre-wrap;font-family:var(--font-mono,monospace);font-size:var(--t-sm);line-height:1.75;margin:0"),
				g.Raw(r.Highlighted),
			),
		),
		// Match list
		h.Div(
			h.Style("border-top:var(--bw-2) solid var(--ink)"),
			h.Div(h.Style(kicker+";padding:var(--sp-2) var(--sp-4)"), g.Text("MATCHES")),
			h.Div(g.Group(matchItems)),
		),
	}

	// Replace output
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
		// Toggle signal, then re-evaluate
		ui.On("click", fmt.Sprintf("$%s=!$%s; @post('/api/eval')", signal, signal)),
		// Reflect active state via data-active attribute
		g.Attr("data-attr", fmt.Sprintf(`{"data-active":$%s?"true":"false"}`, signal)),
		g.Text(label),
	)
}

const regexCSS = `
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
.rx-textarea:focus, .rx-pattern-wrap:focus-within {
  outline: 3px solid var(--accent);
  outline-offset: 1px;
}
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
.rx-flag:hover { background: var(--surface-2); }
.rx-flag[data-active="true"] {
  background: var(--ink);
  color: var(--bg);
}
mark.rx-match {
  background: var(--accent);
  color: var(--accent-ink, var(--ink));
  border-radius: 2px;
  padding: 0 1px;
}
.rx-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--sp-5);
}
@media (max-width: 720px) {
  .rx-grid { grid-template-columns: 1fr; }
}
`
