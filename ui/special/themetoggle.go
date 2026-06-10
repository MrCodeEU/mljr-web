// Package special holds composite components that wire app-wide concerns
// (theming, captcha, etc.). These compose primitives + icons + Datastar.
package special

import (
	"mljr-web/ui/icon"
	"mljr-web/ui/primitive"
	"mljr-web/ui/token"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

// ThemeToggleRoot must be rendered once per page (inside <body>). It declares
// the $theme and $mode signals seeded from window.__mljr* (populated by the
// pre-paint script in PageShell), and uses a data-effect to mirror signal
// changes back onto <html data-*> and into localStorage.
func ThemeToggleRoot(defaultTheme token.Theme, defaultMode token.Mode) g.Node {
	if defaultTheme == "" {
		defaultTheme = token.ThemeSwissBrut
	}
	if defaultMode == "" {
		defaultMode = token.ModeLight
	}
	seed := `{theme: window.__mljrTheme || '` + string(defaultTheme) + `', mode: window.__mljrMode || '` + string(defaultMode) + `'}`
	syncEffect := `document.documentElement.setAttribute('data-theme', $theme);` +
		`document.documentElement.setAttribute('data-mode', $mode);` +
		`try{localStorage.setItem('mljr-theme',$theme);localStorage.setItem('mljr-mode',$mode)}catch(e){}`
	return h.Div(
		g.Attr("data-signals", seed),
		g.Attr("data-effect", syncEffect),
	)
}

// ThemeToggle cycles $theme across the available themes.
func ThemeToggle() g.Node {
	return primitive.Button(
		primitive.ButtonProps{
			Variant: token.Outline,
			Size:    token.SizeIcon,
			Attrs: []g.Node{
				g.Attr("aria-label", "Toggle theme"),
				g.Attr("title", "Toggle theme"),
				g.Attr("data-on:click", `$theme = $theme === 'swissbrut' ? 'ink' : 'swissbrut'`),
			},
		},
		icon.Icon("lucide:palette"),
	)
}

// ModeToggle flips $mode between light and dark.
func ModeToggle() g.Node {
	return primitive.Button(
		primitive.ButtonProps{
			Variant: token.Outline,
			Size:    token.SizeIcon,
			Attrs: []g.Node{
				g.Attr("aria-label", "Toggle light/dark mode"),
				g.Attr("title", "Toggle light/dark mode"),
				g.Attr("data-on:click", `$mode = $mode === 'light' ? 'dark' : 'light'`),
			},
		},
		h.Span(g.Attr("data-show", "$mode === 'light'"), icon.Icon("lucide:moon")),
		h.Span(g.Attr("data-show", "$mode === 'dark'"), icon.Icon("lucide:sun")),
	)
}
