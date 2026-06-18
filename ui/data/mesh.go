package data

import (
	"fmt"
	stdhtml "html"
	"math"
	"strings"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

// MeshNode is one infra host in a Tailscale-style mesh diagram.
type MeshNode struct {
	Name     string
	OS       string   // OS family, e.g. "linux", "windows", "macOS"
	Online   bool
	Relay    bool     // true if currently reachable only via a DERP relay (not direct)
	Services []string // service names hosted here
}

type MeshProps struct {
	// Size is the square viewBox size in px (default 320).
	Size int
}

// Mesh renders a hub-and-spoke network diagram as an SVG: a central tailnet
// hub with one spoke per host, colored by online/relay state. Server-side
// only — zero JS, same approach as Heatmap.
func Mesh(p MeshProps, nodes []MeshNode) g.Node {
	if p.Size == 0 {
		p.Size = 320
	}
	if len(nodes) == 0 {
		return nil
	}

	size := float64(p.Size)
	cx, cy := size/2, size/2
	hubR := size * 0.07
	nodeR := size * 0.085
	orbit := size*0.5 - nodeR - 6
	n := len(nodes)

	var sb strings.Builder
	fmt.Fprintf(&sb, `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %.0f %.0f" style="width:100%%;height:auto;max-width:%.0fpx" data-component="mesh">`,
		size, size, size)

	// Edges first so node circles draw on top of them.
	for i, node := range nodes {
		x, y := meshPoint(cx, cy, orbit, i, n)
		stroke := "var(--accent)"
		dash := ""
		switch {
		case !node.Online:
			stroke = "var(--muted)"
			dash = ` stroke-dasharray="4 4"`
		case node.Relay:
			stroke = "color-mix(in srgb, var(--accent) 50%, orange)"
		}
		fmt.Fprintf(&sb, `<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="%s" stroke-width="2"%s/>`,
			cx, cy, x, y, stroke, dash)
	}

	// Hub.
	fmt.Fprintf(&sb, `<circle cx="%.1f" cy="%.1f" r="%.1f" fill="var(--ink)"/>`, cx, cy, hubR)
	fmt.Fprintf(&sb, `<text x="%.1f" y="%.1f" fill="var(--bg)" font-size="%.1f" font-family="var(--font-mono)" text-anchor="middle" dominant-baseline="middle" font-weight="900">tailnet</text>`,
		cx, cy+1, hubR*0.42)

	// Host nodes.
	for i, node := range nodes {
		x, y := meshPoint(cx, cy, orbit, i, n)

		fill := "var(--surface-2)"
		if node.Online {
			fill = "color-mix(in srgb, var(--accent) 25%, var(--surface-2))"
		}

		title := node.Name
		if len(node.Services) > 0 {
			title = fmt.Sprintf("%s — %s", node.Name, strings.Join(node.Services, ", "))
		}

		sb.WriteString(`<g>`)
		fmt.Fprintf(&sb, `<title>%s</title>`, stdhtml.EscapeString(title))
		fmt.Fprintf(&sb, `<circle cx="%.1f" cy="%.1f" r="%.1f" fill="%s" stroke="var(--ink)" stroke-width="2"/>`,
			x, y, nodeR, fill)

		dotColor := "#22c55e"
		if !node.Online {
			dotColor = "var(--muted)"
		}
		fmt.Fprintf(&sb, `<circle cx="%.1f" cy="%.1f" r="3.5" fill="%s"/>`, x+nodeR*0.6, y-nodeR*0.6, dotColor)

		fmt.Fprintf(&sb, `<text x="%.1f" y="%.1f" fill="var(--ink)" font-size="%.1f" font-family="var(--font-mono)" text-anchor="middle" dominant-baseline="middle" font-weight="800">%s</text>`,
			x, y+1, nodeR*0.32, stdhtml.EscapeString(osAbbrev(node.OS)))

		labelY := y + nodeR + 12
		fmt.Fprintf(&sb, `<text x="%.1f" y="%.1f" fill="var(--ink)" font-size="9" font-family="var(--font-mono)" text-anchor="middle" font-weight="700">%s</text>`,
			x, labelY, stdhtml.EscapeString(node.Name))
		if len(node.Services) > 0 {
			fmt.Fprintf(&sb, `<text x="%.1f" y="%.1f" fill="var(--muted)" font-size="8" font-family="var(--font-mono)" text-anchor="middle">%s</text>`,
				x, labelY+11, stdhtml.EscapeString(serviceSummary(node.Services)))
		}
		sb.WriteString(`</g>`)
	}

	sb.WriteString(`</svg>`)

	return h.Div(
		g.Attr("data-component", "mesh-wrap"),
		h.Style("overflow-x:auto"),
		g.Raw(sb.String()),
	)
}

func meshPoint(cx, cy, orbit float64, i, n int) (float64, float64) {
	angle := -math.Pi/2 + float64(i)*(2*math.Pi/float64(n))
	return cx + orbit*math.Cos(angle), cy + orbit*math.Sin(angle)
}

func osAbbrev(os string) string {
	switch strings.ToLower(os) {
	case "linux":
		return "LNX"
	case "windows":
		return "WIN"
	case "macos", "darwin":
		return "MAC"
	case "ios":
		return "iOS"
	case "android":
		return "AND"
	default:
		if len(os) >= 3 {
			return strings.ToUpper(os[:3])
		}
		return strings.ToUpper(os)
	}
}

func serviceSummary(services []string) string {
	if len(services) <= 2 {
		return strings.Join(services, ", ")
	}
	return fmt.Sprintf("%s +%d more", strings.Join(services[:2], ", "), len(services)-2)
}
