package pages

import (
	"fmt"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"mljr-web/ui"
	"mljr-web/ui/special"
	"mljr-web/ui/swiss"
)

func Home() g.Node {
	initial := EvalCodec(CodecInput{Mode: "hash", Op: "encode", Input: "Hello, World!"})

	return special.ToolShell(special.ToolShellProps{
		ID:          "codec",
		Name:        "Codec",
		Title:       "Codec — Hash & Encode/Decode",
		Description: "Compute MD5/SHA1/SHA256/SHA512 hashes and encode/decode Base64, URL, HTML — all in the browser via Go + Datastar.",
		CSS:         codecCSS,
		Sections: []swiss.NavItem{
			{ID: "tool", Num: "00", Label: "Werkzeug", Sub: "Tool"},
			{ID: "bench", Num: "01", Label: "Werkbank", Sub: "Workbench"},
		},
	},
		swiss.Section(swiss.SectionProps{ID: "tool", Num: "00", Title: "Werkzeug", Sub: "Tool", Field: "block",
			Note: "All operations run server-side in Go. No data is stored or logged."},
			swiss.Cell(swiss.CellProps{},
				swiss.ToolHead("Codec", "Hash text with MD5, SHA-1, SHA-256 and SHA-512, or encode and decode Base64, URL and HTML entities. Results update as you type.",
					"Hash", "Base64", "URL", "HTML", "Go"),
			),
		),
		h.Div(
			ui.Signals(`{codec:"hash",op:"encode",input:"Hello, World!"}`),
			ui.On("input__debounce.200ms", "@post('/api/eval')"),
			swiss.Section(swiss.SectionProps{ID: "bench", Num: "01", Title: "Werkbank", Sub: "Workbench",
				Note: "Pick a mode, type on the left, read the result on the right."},
				swiss.Cell(swiss.CellProps{},
					swiss.Kicker("A. Mode", nil),
					swiss.Choice("Mode",
						modeTab("hash", "Hash", "MD5 · SHA1 · SHA256 · SHA512"),
						modeTab("base64", "Base64", "encode / decode"),
						modeTab("url", "URL", "percent-encoding"),
						modeTab("html", "HTML", "entity encoding"),
					),
					h.Div(g.Attr("data-show", `$codec!="hash"`),
						swiss.Choice("Operation", opBtn("encode", "Encode"), opBtn("decode", "Decode")),
					),
				),
				swiss.Cell(swiss.CellProps{Span: 4},
					swiss.Field("cc-input", "B. Input",
						h.Textarea(h.ID("cc-input"), g.Attr("rows", "10"), g.Attr("spellcheck", "false"), ui.Bind("input")),
					),
				),
				swiss.Cell(swiss.CellProps{Span: 5, Variant: "flush"},
					OutputFragment(initial),
				),
			),
		),
	)
}

func modeTab(value, label, sub string) g.Node {
	return swiss.Btn(false, []g.Node{
		g.Attr("aria-pressed", fmt.Sprint(value == "hash")),
		g.Attr("data-attr:aria-pressed", fmt.Sprintf(`$codec=="%s" ? "true" : "false"`, value)),
		ui.On("click", fmt.Sprintf(`$codec="%s"; @post('/api/eval')`, value)),
	}, g.Text(label), h.Small(g.Text(sub)))
}

func opBtn(value, label string) g.Node {
	return swiss.Btn(false, []g.Node{
		g.Attr("aria-pressed", fmt.Sprint(value == "encode")),
		g.Attr("data-attr:aria-pressed", fmt.Sprintf(`$op=="%s" ? "true" : "false"`, value)),
		ui.On("click", fmt.Sprintf(`$op="%s"; @post('/api/eval')`, value)),
	}, g.Text(label))
}

// OutputFragment renders the result panel (id="cc-output"), for the first
// render and for every SSE patch.
func OutputFragment(r CodecResult) g.Node {
	if r.Err != "" {
		return h.Div(h.ID("cc-output"), swiss.ResultBar("Error", nil), swiss.ErrorBox(r.Err))
	}
	if r.Mode == "hash" && r.Hash != nil {
		hashes := [][2]string{
			{"MD5", r.Hash.MD5},
			{"SHA-1", r.Hash.SHA1},
			{"SHA-256", r.Hash.SHA256},
			{"SHA-512", r.Hash.SHA512},
		}
		blocks := make([]g.Node, len(hashes))
		for i, hh := range hashes {
			blocks[i] = swiss.ResultBlock(fmt.Sprintf("%s · %d chars", hh[0], len(hh[1])), h.Pre(g.Text(hh[1])))
		}
		return h.Div(h.ID("cc-output"), swiss.ResultBar("Hash digests", g.Text("4 algorithms")), g.Group(blocks))
	}
	op := "Encode"
	if r.Op == "decode" {
		op = "Decode"
	}
	return h.Div(h.ID("cc-output"),
		swiss.ResultBar(r.Mode+" · "+op, g.Text(fmt.Sprintf("%d chars", len(r.Output)))),
		swiss.ResultBlock("Output", h.Pre(g.Text(r.Output))),
	)
}

// The hash/encode mode lives in $codec: $mode is the page's light/dark
// signal (ThemeToggleRoot) and must not be reused.
const codecCSS = ``
