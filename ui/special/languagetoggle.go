package special

import (
	"strings"

	"mljr-web/ui/icon"
	"mljr-web/ui/primitive"
	"mljr-web/ui/token"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type Language struct {
	Code  string // e.g. "en", "de"
	Label string // display label, e.g. "EN"
	Title string // full name for aria-label, e.g. "English"
	Flag  string // iconify icon, e.g. "circle-flags:gb"
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
// Shows a circular flag icon beside the locale label. Styled like ThemeToggle.
func LanguageToggle(p LanguageToggleProps) g.Node {
	if p.Cookie == "" {
		p.Cookie = "lang"
	}
	if len(p.Languages) == 0 {
		return g.Group{}
	}

	// Build cycle JS
	codes := make([]string, len(p.Languages))
	for i, l := range p.Languages {
		codes[i] = `'` + l.Code + `'`
	}
	codesJS := "[" + strings.Join(codes, ",") + "]"

	var current Language
	current = p.Languages[0]
	for _, l := range p.Languages {
		if l.Code == p.Current {
			current = l
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

	ariaLabel := "Switch language"
	if current.Title != "" {
		ariaLabel = "Language: " + current.Title + " — click to switch"
	}

	return primitive.Button(
		primitive.ButtonProps{
			Variant: token.Outline,
			Size:    token.SizeIcon,
			Attrs: []g.Node{
				g.Attr("aria-label", ariaLabel),
				g.Attr("title", ariaLabel),
				g.Attr("onclick", clickExpr),
				g.Attr("style", "aspect-ratio:unset;padding:.55rem .75rem;display:inline-flex;align-items:center;gap:.4rem;"),
			},
		},
		g.If(current.Flag != "",
			h.Span(h.Style("display:inline-flex;align-items:center;flex-shrink:0;line-height:0;"),
				icon.Icon(current.Flag, icon.Props{Size: "1em"}),
			),
		),
		h.Span(h.Style("font-size:var(--t-xs);font-weight:900;letter-spacing:.06em;font-family:var(--font-display);line-height:1"),
			g.Text(current.Label),
		),
	)
}
