package primitive

import (
	"mljr-web/ui/icon"
	"mljr-web/ui/token"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type ShareButtonProps struct {
	URL     string // URL to share (defaults to current page via JS)
	Title   string // share title
	Text    string // share body text
	Variant token.Variant
	Size    token.Size
	Label   string // button label (default "Share")
}

// ShareButton uses the Web Share API when available, falls back to copy-link.
// No server needed — all client-side.
func ShareButton(p ShareButtonProps) g.Node {
	if p.Label == "" {
		p.Label = "Share"
	}
	if p.Variant == "" {
		p.Variant = token.Outline
	}

	url := p.URL
	if url == "" {
		url = "window.location.href"
	} else {
		url = `'` + url + `'`
	}

	clickExpr := `(function(){` +
		`var u=` + url + `;` +
		`var t=` + jsStr(p.Title) + `;` +
		`var x=` + jsStr(p.Text) + `;` +
		`if(navigator.share){navigator.share({url:u,title:t,text:x}).catch(function(){});}` +
		`else{navigator.clipboard.writeText(u).then(function(){alert('Link copied!');}).catch(function(){prompt('Copy this link:',u);});}` +
		`})()`

	return Button(ButtonProps{Variant: p.Variant, Size: p.Size},
		g.Attr("data-on:click", clickExpr),
		icon.Icon("lucide:share"),
		g.If(p.Label != "", h.Span(g.Text(" "+p.Label))),
	)
}

// jsStr quotes a Go string as a JS string literal (single quotes, escaped).
func jsStr(s string) string {
	if s == "" {
		return "''"
	}
	out := make([]byte, 0, len(s)+2)
	out = append(out, '\'')
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '\'' {
			out = append(out, '\\', '\'')
		} else if c == '\\' {
			out = append(out, '\\', '\\')
		} else {
			out = append(out, c)
		}
	}
	return string(append(out, '\''))
}
