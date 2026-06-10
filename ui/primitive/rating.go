package primitive

import (
	"fmt"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type RatingProps struct {
	Signal   string
	Max      int // default 5
	ReadOnly bool
	Attrs    []g.Node
}

// Rating renders a star rating widget bound to a Datastar signal.
// Stars filled state is driven reactively; click sets the signal value.
func Rating(p RatingProps, attrs ...g.Node) g.Node {
	max := p.Max
	if max <= 0 {
		max = 5
	}
	sig := p.Signal

	stars := make([]g.Node, max)
	for i := 1; i <= max; i++ {
		n := i
		starAttrs := []g.Node{
			g.Attr("data-slot", "star"),
			g.Attr("data-attr", fmt.Sprintf(`{"data-state":$%s>=%d?"active":""}`, sig, n)),
			g.Text("★"),
		}
		if !p.ReadOnly {
			starAttrs = append(starAttrs,
				h.Type("button"),
				g.Attr("data-on:click", fmt.Sprintf("$%s=%d", sig, n)),
				g.Attr("aria-label", fmt.Sprintf("%d star", n)),
			)
			stars[n-1] = h.Button(starAttrs...)
		} else {
			stars[n-1] = h.Span(starAttrs...)
		}
	}

	nodes := []g.Node{
		g.Attr("data-component", "rating"),
		g.If(p.ReadOnly, g.Attr("data-state", "readonly")),
		g.Group(p.Attrs),
		g.Group(attrs),
	}
	nodes = append(nodes, stars...)
	return h.Div(nodes...)
}
