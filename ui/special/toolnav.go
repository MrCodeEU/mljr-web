package special

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"mljr-web/ui/icon"
	"mljr-web/ui/layout"
	"mljr-web/ui/primitive"
	"mljr-web/ui/token"
)

type toolLink struct {
	ID    string
	Label string
	URL   string
}

var liveTools = []toolLink{
	{ID: "regex", Label: "Regex Lab", URL: "https://regex.mljr.eu"},
	{ID: "cron", Label: "Cron Explorer", URL: "https://cron.mljr.eu"},
	{ID: "codec", Label: "Codec", URL: "https://codec.mljr.eu"},
}

// ToolNavbar is the shared navbar for standalone live-tool apps (regex, cron,
// codec): brand link back to the main portfolio, links to sibling tools, and
// the same theme/mode/GitHub actions used on the homepage navbar. current is
// the tool's own ID so it's excluded from its own nav links.
func ToolNavbar(current string) g.Node {
	var navLinks []g.Node
	for _, t := range liveTools {
		if t.ID == current {
			continue
		}
		navLinks = append(navLinks, h.A(h.Href(t.URL), g.Text(t.Label)))
	}

	return layout.Navbar(layout.NavbarProps{},
		h.A(h.Href("https://mljr.eu"),
			h.Img(
				h.Src("/static/img/logo/Logo-h.png"),
				h.Alt("mljr.eu"),
				h.Width("172"),
				h.Height("32"),
				h.Style("height:32px;width:auto"),
			),
		),
		g.Group(navLinks),
		g.Group{
			ThemeToggle(),
			ModeToggle(),
			h.A(
				h.Href("https://github.com/MrCodeEU"),
				g.Attr("target", "_blank"),
				g.Attr("rel", "noopener noreferrer"),
				g.Attr("aria-label", "GitHub"),
				primitive.Button(
					primitive.ButtonProps{
						Variant: token.Outline,
						Size:    token.SizeIcon,
						Attrs:   []g.Node{g.Attr("aria-hidden", "true"), g.Attr("tabindex", "-1")},
					},
					icon.Icon("lucide:github"),
				),
			),
		},
	)
}
