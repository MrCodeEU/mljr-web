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

func Home() g.Node {
	initial := EvalCodec(CodecInput{Mode: "hash", Op: "encode", Input: "Hello, World!"})

	return layout.PageShell(
		layout.PageProps{
			Title:       "Codec — Hash & Encode/Decode",
			Description: "Compute MD5/SHA1/SHA256/SHA512 hashes and encode/decode Base64, URL, HTML — all in the browser via Go + Datastar.",
			Theme:       token.ThemeSwissBrut,
			Mode:        token.ModeLight,
			HeadExtra: []g.Node{
				g.El("style", g.Raw(codecCSS)),
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
						g.Text("CODEC"),
					),
					h.Div(
						h.Style("display:flex;gap:var(--sp-2);margin-left:auto;flex-shrink:0"),
						primitive.Tag(primitive.TagProps{Tone: token.ToneCyan}, g.Text("HASH")),
						primitive.Tag(primitive.TagProps{Tone: token.ToneViolet}, g.Text("ENCODE")),
						primitive.Tag(primitive.TagProps{Tone: token.ToneLime},
							icon.Icon("simple-icons:go", icon.Props{Size: ".9rem"}),
							g.Text("Go"),
						),
					),
				),

				// Signal scope
				h.Div(
					ui.Signals(`{mode:"hash",op:"encode",input:"Hello, World!"}`),
					ui.On("input__debounce.200ms", "@post('/api/eval')"),

					// Mode tabs
					h.Div(
						h.Style("margin-bottom:var(--sp-5)"),
						h.Label(
							h.Style("display:block;font-size:var(--t-xs);font-weight:900;text-transform:uppercase;letter-spacing:.12em;margin-bottom:var(--sp-2);opacity:.65"),
							g.Text("MODE"),
						),
						h.Div(
							h.Class("cc-tabs"),
							modeTab("hash", "Hash", "MD5·SHA1·SHA256·SHA512"),
							modeTab("base64", "Base64", "encode / decode"),
							modeTab("url", "URL", "percent-encoding"),
							modeTab("html", "HTML", "entity encoding"),
						),
					),

					// Encode/Decode toggle (hidden in hash mode)
					h.Div(
						h.Style("margin-bottom:var(--sp-5)"),
						g.Attr("data-show", `$mode!="hash"`),
						h.Div(
							h.Style("display:flex;gap:var(--sp-2);align-items:center"),
							h.Label(
								h.Style("font-size:var(--t-xs);font-weight:900;text-transform:uppercase;letter-spacing:.12em;opacity:.65"),
								g.Text("OPERATION"),
							),
							h.Div(
								h.Style("display:flex;gap:var(--sp-1)"),
								opBtn("encode", "Encode"),
								opBtn("decode", "Decode"),
							),
						),
					),

					// Input
					h.Div(
						h.Style("margin-bottom:var(--sp-5)"),
						h.Label(
							h.Style("display:block;font-size:var(--t-xs);font-weight:900;text-transform:uppercase;letter-spacing:.12em;margin-bottom:var(--sp-2);opacity:.65"),
							g.Text("INPUT"),
						),
						h.Textarea(
							h.Class("cc-textarea"),
							g.Attr("rows", "6"),
							g.Attr("spellcheck", "false"),
							ui.Bind("input"),
						),
					),

					// Output
					primitive.Card(primitive.CardProps{
						Tone:  token.ToneNone,
						Attrs: []g.Node{h.Style("padding:0;overflow:hidden;min-height:200px")},
					},
						OutputFragment(initial),
					),
				),

				// Note
				h.P(
					h.Style("margin-top:var(--sp-8);font-size:var(--t-xs);opacity:.45;font-family:var(--font-mono,monospace)"),
					g.Text("All operations run server-side in Go. No data is stored or logged."),
				),
			),
		),
	)
}

func modeTab(value, label, sub string) g.Node {
	return h.Button(
		h.Type("button"),
		h.Class("cc-tab"),
		g.Attr("data-class", fmt.Sprintf(`{"cc-tab--active":$mode=="%s"}`, value)),
		ui.On("click", fmt.Sprintf(`$mode="%s"; @post('/api/eval')`, value)),
		h.Div(h.Style("font-weight:900;font-size:var(--t-sm)"), g.Text(label)),
		h.Div(h.Style("font-size:10px;opacity:.55;font-family:var(--font-mono,monospace)"), g.Text(sub)),
	)
}

func opBtn(value, label string) g.Node {
	return h.Button(
		h.Type("button"),
		h.Class("cc-flag"),
		g.Attr("data-class", fmt.Sprintf(`{"cc-flag--active":$op=="%s"}`, value)),
		ui.On("click", fmt.Sprintf(`$op="%s"; @post('/api/eval')`, value)),
		g.Text(label),
	)
}

func OutputFragment(r CodecResult) g.Node {
	kicker := "font-size:var(--t-xs);font-weight:900;text-transform:uppercase;letter-spacing:.12em;opacity:.5"

	if r.Err != "" {
		return h.Div(
			h.ID("cc-output"),
			h.Style("padding:var(--sp-4)"),
			h.Div(
				h.Style("display:flex;align-items:center;gap:var(--sp-2);font-family:var(--font-mono,monospace);font-size:var(--t-sm);color:var(--blush-ink,#9f1239);background:var(--blush-bg,#fde8e8);border:var(--bw-2) solid var(--ink);padding:var(--sp-3)"),
				icon.Icon("lucide:alert-circle", icon.Props{Size: "1rem"}),
				h.Span(g.Text(r.Err)),
			),
		)
	}

	var sections []g.Node

	if r.Mode == "hash" && r.Hash != nil {
		// Header
		sections = append(sections,
			h.Div(
				h.Style("display:flex;align-items:center;gap:var(--sp-3);padding:var(--sp-3) var(--sp-4);background:var(--ink);color:var(--bg)"),
				h.Span(h.Style("font-size:var(--t-sm);font-weight:900;font-family:var(--font-mono,monospace)"), g.Text("■ HASH DIGESTS")),
			),
		)
		hashes := [][2]string{
			{"MD5", r.Hash.MD5},
			{"SHA-1", r.Hash.SHA1},
			{"SHA-256", r.Hash.SHA256},
			{"SHA-512", r.Hash.SHA512},
		}
		tones := []token.Tone{token.ToneCyan, token.ToneViolet, token.ToneLime, token.ToneYellow}
		for i, hh := range hashes {
			sections = append(sections, hashRow(hh[0], hh[1], tones[i]))
		}
	} else {
		opLabel := "ENCODE"
		if r.Op == "decode" {
			opLabel = "DECODE"
		}
		sections = append(sections,
			h.Div(
				h.Style("display:flex;align-items:center;gap:var(--sp-3);padding:var(--sp-3) var(--sp-4);background:var(--ink);color:var(--bg)"),
				h.Span(h.Style("font-size:var(--t-sm);font-weight:900;font-family:var(--font-mono,monospace)"),
					g.Text(fmt.Sprintf("■ %s · %s", r.Mode, opLabel)),
				),
			),
			h.Div(
				h.Style("padding:var(--sp-4)"),
				h.Div(h.Style(kicker+";margin-bottom:var(--sp-2)"), g.Text("OUTPUT")),
				h.Pre(
					h.Style("white-space:pre-wrap;word-break:break-all;font-family:var(--font-mono,monospace);font-size:var(--t-sm);line-height:1.75;margin:0"),
					g.Text(r.Output),
				),
			),
		)
	}

	return h.Div(h.ID("cc-output"), g.Group(sections))
}

func hashRow(algo, digest string, tone token.Tone) g.Node {
	return h.Div(
		h.Style("padding:var(--sp-3) var(--sp-4);border-bottom:var(--bw-1) solid var(--line)"),
		h.Div(h.Style("display:flex;align-items:center;gap:var(--sp-3);margin-bottom:var(--sp-1)"),
			primitive.Tag(primitive.TagProps{Tone: tone}, g.Text(algo)),
			h.Div(h.Style("font-size:var(--t-xs);opacity:.45;font-family:var(--font-mono,monospace)"),
				g.Text(fmt.Sprintf("%d chars", len(digest))),
			),
		),
		h.Code(
			h.Style("font-family:var(--font-mono,monospace);font-size:var(--t-sm);word-break:break-all;line-height:1.6"),
			g.Text(digest),
		),
	)
}

const codecCSS = `
.cc-tabs {
  display: flex;
  gap: var(--sp-2);
  flex-wrap: wrap;
}

.cc-tab {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  padding: var(--sp-2) var(--sp-4);
  border: var(--bw-2) solid var(--ink);
  background: var(--bg);
  color: var(--ink);
  cursor: pointer;
  transition: background .1s, color .1s;
  min-width: 110px;
}
.cc-tab:hover { background: var(--surface-2, var(--surface)); }
.cc-tab--active {
  background: var(--ink);
  color: var(--bg);
}
.cc-tab--active > div:last-child { opacity: .55; }

.cc-flag {
  padding: var(--sp-1) var(--sp-4);
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
.cc-flag:hover { background: var(--surface-2, var(--surface)); }
.cc-flag--active { background: var(--ink); color: var(--bg); }

.cc-textarea {
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
.cc-textarea:focus {
  outline: 3px solid var(--accent);
  outline-offset: 1px;
}
`
