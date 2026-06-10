package special

import (
	"strings"

	"mljr-web/ui/primitive"
	"mljr-web/ui/token"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type Language struct {
	Code  string // e.g. "en", "de", "fr"
	Label string // display label, e.g. "EN"
	Title string // full name for aria-label, e.g. "English"
}

type LanguageToggleProps struct {
	// Languages lists the available locales in cycle order.
	Languages []Language
	// Current is the active language code.
	Current string
	// Cookie is the cookie name to persist the choice (default "lang").
	Cookie string
	// ReloadOnChange triggers a full page reload after setting the cookie.
	ReloadOnChange bool
	// BasePath is the path prefix to prepend when building locale URLs.
	// If empty, uses JS cookie + reload approach instead of URL switching.
	BasePath string
}

// LanguageToggle renders a button that cycles through available languages.
// Persists the choice via cookie and optionally reloads the page.
func LanguageToggle(p LanguageToggleProps) g.Node {
	if p.Cookie == "" {
		p.Cookie = "lang"
	}
	if len(p.Languages) == 0 {
		return g.Group{}
	}

	// Build the cycle: current → next
	codes := make([]string, len(p.Languages))
	for i, l := range p.Languages {
		codes[i] = `'` + l.Code + `'`
	}
	codesJS := "[" + strings.Join(codes, ",") + "]"

	currentLabel := p.Languages[0].Label
	for _, l := range p.Languages {
		if l.Code == p.Current {
			currentLabel = l.Label
		}
	}

	clickExpr := `var langs=` + codesJS + `;` +
		`var cur=document.cookie.split(';').find(c=>c.trim().startsWith('` + p.Cookie + `='));` +
		`var idx=cur?langs.indexOf(cur.split('=')[1].trim()):langs.indexOf('` + p.Current + `');` +
		`var next=langs[(idx+1)%langs.length];` +
		`document.cookie='` + p.Cookie + `='+next+';path=/;max-age=31536000';`

	if p.ReloadOnChange {
		clickExpr += `location.reload();`
	}

	return primitive.Button(
		primitive.ButtonProps{
			Variant: token.Outline,
			Attrs: []g.Node{
				g.Attr("aria-label", "Switch language"),
				g.Attr("title", "Switch language"),
				g.Attr("onclick", clickExpr),
			},
		},
		h.Span(h.Style("font-size:var(--t-xs);font-weight:900;letter-spacing:.04em;font-family:var(--font-display)"),
			g.Text(currentLabel),
		),
	)
}
