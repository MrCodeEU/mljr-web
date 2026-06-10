//go:build showcase

package special

import (
	"mljr-web/ui/registry"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "language-toggle", Name: "Language Toggle", Category: "special",
		Summary: "Cycles through locales, persists in cookie, optional page reload. Drop-in for multilingual sites.",
		Code: `special.LanguageToggle(special.LanguageToggleProps{
    Current: "en",
    Cookie:  "lang",
    ReloadOnChange: true,
    Languages: []special.Language{
        {Code: "en", Label: "EN", Title: "English"},
        {Code: "de", Label: "DE", Title: "German"},
        {Code: "fr", Label: "FR", Title: "French"},
    },
})`,
		Render: func(p map[string]string) g.Node {
			return h.Div(
				h.Style("display:flex;flex-direction:column;gap:var(--sp-5)"),
				h.Div(
					h.P(h.Style("font-size:var(--t-xs);opacity:.5;font-weight:700;text-transform:uppercase;letter-spacing:.06em;margin:0 0 var(--sp-2)"), g.Text("2 languages")),
					LanguageToggle(LanguageToggleProps{
						Current: "en",
						Languages: []Language{
							{Code: "en", Label: "EN", Title: "English"},
							{Code: "de", Label: "DE", Title: "German"},
						},
					}),
				),
				h.Div(
					h.P(h.Style("font-size:var(--t-xs);opacity:.5;font-weight:700;text-transform:uppercase;letter-spacing:.06em;margin:0 0 var(--sp-2)"), g.Text("4 languages")),
					LanguageToggle(LanguageToggleProps{
						Current: "de",
						Languages: []Language{
							{Code: "en", Label: "EN"},
							{Code: "de", Label: "DE"},
							{Code: "fr", Label: "FR"},
							{Code: "ja", Label: "JA"},
						},
					}),
				),
				h.P(h.Style("color:var(--muted);font-size:var(--t-xs)"), g.Text("Cookie is set on click. Add ReloadOnChange: true to reload the page after switching.")),
			)
		},
	})
}
