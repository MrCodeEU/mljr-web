// Package icon renders pre-generated Iconify SVGs. The icons map is produced
// by tools/icongen from tools/icongen/icons.txt.
package icon

import (
	"fmt"
	"sort"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type Props struct {
	Size  string // CSS size; default "1.25em"
	Label string // if set: role="img" + aria-label; else aria-hidden="true"
	Attrs []g.Node
}

// Icon renders one icon by "set:name". Renders an HTML comment if missing.
func Icon(name string, p ...Props) g.Node {
	var pp Props
	if len(p) > 0 {
		pp = p[0]
	}
	if pp.Size == "" {
		pp.Size = "1.25em"
	}
	svg, ok := icons[name]
	if !ok {
		return g.Raw(fmt.Sprintf(`<!-- icon %q missing -->`, name))
	}
	attrs := []g.Node{
		g.Attr("data-component", "icon"),
		h.Style("display:inline-flex;width:" + pp.Size + ";height:" + pp.Size + ";font-size:" + pp.Size + ";flex-shrink:0"),
	}
	if pp.Label != "" {
		attrs = append(attrs, h.Role("img"), g.Attr("aria-label", pp.Label))
	} else {
		attrs = append(attrs, g.Attr("aria-hidden", "true"))
	}
	attrs = append(attrs, pp.Attrs...)
	attrs = append(attrs, g.Raw(svg))
	return h.Span(attrs...)
}

// Has reports whether the named icon is registered.
func Has(name string) bool { _, ok := icons[name]; return ok }

// All returns all registered icon names, sorted alphabetically.
func All() []string {
	names := make([]string, 0, len(icons))
	for k := range icons {
		names = append(names, k)
	}
	sort.Strings(names)
	return names
}
