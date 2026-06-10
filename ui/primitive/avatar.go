package primitive

import (
	"mljr-web/ui/token"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type AvatarShape string

const (
	AvatarCircle AvatarShape = ""
	AvatarSquare AvatarShape = "square"
)

type AvatarStatus string

const (
	AvatarOnline  AvatarStatus = "online"
	AvatarAway    AvatarStatus = "away"
	AvatarOffline AvatarStatus = "offline"
)

type AvatarProps struct {
	Src      string // image URL — if empty, shows initials
	Initials string // fallback text when no Src
	Alt      string
	Size     token.Size   // sm | md (default) | lg | xl
	Shape    AvatarShape  // "" = circle | "square"
	Status   AvatarStatus // "" | online | away | offline
	Tone     token.Tone   // background tone for initials avatar
	Attrs    []g.Node
}

// Avatar renders a user avatar with image or initials fallback, optional status dot.
func Avatar(p AvatarProps) g.Node {
	var statusDot g.Node
	if p.Status != "" {
		statusDot = h.Span(
			g.Attr("data-slot", "status"),
			g.Attr("data-state", string(p.Status)),
		)
	}

	var inner g.Node
	if p.Src != "" {
		alt := p.Alt
		if alt == "" {
			alt = p.Initials
		}
		inner = h.Img(h.Src(p.Src), h.Alt(alt))
	} else {
		inner = h.Span(g.Attr("data-slot", "initials"), g.Text(p.Initials))
	}

	return h.Div(
		g.Attr("data-component", "avatar"),
		g.If(p.Size != "", g.Attr("data-size", string(p.Size))),
		g.If(p.Shape != "", g.Attr("data-shape", string(p.Shape))),
		g.If(p.Tone != "", g.Attr("data-tone", string(p.Tone))),
		g.Group(p.Attrs),
		inner,
		statusDot,
	)
}
