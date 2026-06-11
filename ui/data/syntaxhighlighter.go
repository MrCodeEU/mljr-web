package data

import (
	"bytes"
	"html/template"
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"

	g "maragu.dev/gomponents"
	ghtml "maragu.dev/gomponents/html"
)

type SyntaxHighlighterProps struct {
	// Language (e.g. "go", "typescript", "bash"). Empty = auto-detect.
	Language string
	// Theme: "github", "dracula", "monokai", "vs", "solarized-dark", "nord" (default "github").
	Theme string
	// ShowLineNumbers renders line numbers.
	ShowLineNumbers bool
	// Filename shown above the code block.
	Filename string
}

// SyntaxHighlighter renders source code with server-side syntax highlighting.
// Uses the chroma library — zero client JavaScript. Themed HTML output.
func SyntaxHighlighter(p SyntaxHighlighterProps, code string) g.Node {
	if p.Theme == "" {
		p.Theme = "github"
	}

	// Get lexer
	var lexer chroma.Lexer
	if p.Language != "" {
		lexer = lexers.Get(p.Language)
	}
	if lexer == nil {
		lexer = lexers.Analyse(code)
	}
	if lexer == nil {
		lexer = lexers.Fallback
	}
	lexer = chroma.Coalesce(lexer)

	// Get style
	style := styles.Get(p.Theme)
	if style == nil {
		style = styles.Fallback
	}

	// Format to HTML
	formatter := html.New(
		html.WithClasses(false), // inline styles — no extra CSS needed
		html.WithLineNumbers(p.ShowLineNumbers),
		html.TabWidth(4),
	)

	tokens, err := lexer.Tokenise(nil, code)
	var buf bytes.Buffer
	if err == nil {
		_ = formatter.Format(&buf, style, tokens)
	} else {
		buf.WriteString("<pre><code>" + template.HTMLEscapeString(code) + "</code></pre>")
	}

	lang := lexer.Config().Name

	var headerNode g.Node
	if p.Filename != "" || lang != "" {
		headerNode = ghtml.Div(
			g.Attr("data-slot", "header"),
			ghtml.Div(g.Attr("data-slot", "filename"), g.Text(p.Filename)),
			ghtml.Span(g.Attr("data-slot", "lang"), g.Text(lang)),
		)
	}

	return ghtml.Div(
		g.Attr("data-component", "syntax-highlighter"),
		g.Attr("data-lang", strings.ToLower(lang)),
		headerNode,
		g.Raw(buf.String()),
	)
}
