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
	OS       string // OS family, e.g. "linux", "windows", "macOS"
	Online   bool
	Relay    bool     // true if currently reachable only via a DERP relay (not direct)
	Services []string // service names hosted here
}

type MeshProps struct {
	// Size is the square viewBox size in px (default 480).
	Size int
}

type meshKind int

const (
	meshKindHub meshKind = iota
	meshKindHost
)

// meshSimNode is one body in the force simulation: the tailnet hub or a host.
type meshSimNode struct {
	x, y, vx, vy float64
	r            float64
	kind         meshKind
	label        string
	title        string
	online       bool
	relay        bool
}

type meshEdge struct {
	from, to int
	ideal    float64
	strength float64
}

// Mesh renders the tailnet as a force-directed graph with one node per host.
// Service names remain available in each host's tooltip and their count is
// shown below the host; omitting individual service satellites keeps the
// topology readable as the inventory grows.
func Mesh(p MeshProps, nodes []MeshNode) g.Node {
	if p.Size == 0 {
		p.Size = 480
	}
	if len(nodes) == 0 {
		return nil
	}

	size := float64(p.Size)
	hubR := size * 0.055
	hostR := size * 0.07

	sim := []meshSimNode{{x: 0, y: 0, r: hubR, kind: meshKindHub, label: "tailnet"}}
	var edges []meshEdge

	n := len(nodes)
	hostOrbit := size * 0.32
	for i, node := range nodes {
		angle := -math.Pi/2 + float64(i)*(2*math.Pi/float64(n))
		hx := hostOrbit * math.Cos(angle)
		hy := hostOrbit * math.Sin(angle)
		hostIdx := len(sim)

		title := node.Name
		if len(node.Services) > 0 {
			title = fmt.Sprintf("%s — %s", node.Name, strings.Join(node.Services, ", "))
		}
		sim = append(sim, meshSimNode{
			x: hx, y: hy, r: hostR, kind: meshKindHost,
			label: node.Name, title: title, online: node.Online, relay: node.Relay,
		})
		edges = append(edges, meshEdge{from: 0, to: hostIdx, ideal: hostOrbit, strength: 0.06})
	}

	simulateMesh(sim, edges)
	cx, cy, scale := meshFit(sim, size)

	var sb strings.Builder
	fmt.Fprintf(&sb, `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %.0f %.0f" style="width:100%%;height:auto;max-width:%.0fpx" data-component="mesh">`,
		size, size, size)

	project := func(s meshSimNode) (float64, float64, float64) {
		return cx + s.x*scale, cy + s.y*scale, s.r * math.Max(scale, 0.55)
	}

	// Edges first so node circles draw on top of them.
	for _, e := range edges {
		x1, y1, _ := project(sim[e.from])
		x2, y2, _ := project(sim[e.to])
		to := sim[e.to]
		stroke := "var(--muted)"
		dash := ` stroke-dasharray="3 3"`
		strokeW := "1"
		switch to.kind {
		case meshKindHost:
			strokeW = "2"
			dash = ""
			stroke = "var(--accent)"
			if !to.online {
				stroke = "var(--muted)"
				dash = ` stroke-dasharray="4 4"`
			} else if to.relay {
				stroke = "color-mix(in srgb, var(--accent) 50%, orange)"
			}
		}
		fmt.Fprintf(&sb, `<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="%s" stroke-width="%s"%s/>`,
			x1, y1, x2, y2, stroke, strokeW, dash)
	}

	for _, s := range sim {
		x, y, r := project(s)
		switch s.kind {
		case meshKindHub:
			sb.WriteString(`<g>`)
		case meshKindHost:
			sb.WriteString(`<g data-mesh="breathe-host" style="transform-box:fill-box;transform-origin:50% 50%">`)
		}
		if s.title != "" {
			fmt.Fprintf(&sb, `<title>%s</title>`, stdhtml.EscapeString(s.title))
		}

		switch s.kind {
		case meshKindHub:
			fmt.Fprintf(&sb, `<circle cx="%.1f" cy="%.1f" r="%.1f" fill="var(--ink)"/>`, x, y, r)
			fmt.Fprintf(&sb, `<text x="%.1f" y="%.1f" fill="var(--bg)" font-size="%.1f" font-family="var(--font-mono)" text-anchor="middle" dominant-baseline="middle" font-weight="900">tailnet</text>`,
				x, y+1, r*0.42)

		case meshKindHost:
			fill := "var(--surface-2)"
			if s.online {
				fill = "color-mix(in srgb, var(--accent) 25%, var(--surface-2))"
			}
			fmt.Fprintf(&sb, `<circle cx="%.1f" cy="%.1f" r="%.1f" fill="%s" stroke="var(--ink)" stroke-width="2"/>`,
				x, y, r, fill)
			dotColor := "#22c55e"
			if !s.online {
				dotColor = "var(--muted)"
			}
			fmt.Fprintf(&sb, `<circle cx="%.1f" cy="%.1f" r="4" fill="%s" data-mesh="pulse"/>`, x+r*0.6, y-r*0.6, dotColor)
			fmt.Fprintf(&sb, `<text x="%.1f" y="%.1f" fill="var(--ink)" font-size="%.1f" font-family="var(--font-mono)" text-anchor="middle" dominant-baseline="middle" font-weight="800">%s</text>`,
				x, y+1, r*0.32, stdhtml.EscapeString(osAbbrev(findHostOS(nodes, s.label))))
			fmt.Fprintf(&sb, `<text x="%.1f" y="%.1f" fill="var(--ink)" font-size="11" font-family="var(--font-mono)" text-anchor="middle" font-weight="700">%s</text>`,
				x, y+r+14, stdhtml.EscapeString(s.label))
			if count := serviceCount(nodes, s.label); count > 0 {
				fmt.Fprintf(&sb, `<text x="%.1f" y="%.1f" fill="var(--muted)" font-size="11" font-family="var(--font-mono)" text-anchor="middle" font-weight="700">%d svc</text>`,
					x, y+r+25, count)
			}
		}
		sb.WriteString(`</g>`)
	}

	sb.WriteString(`</svg>`)
	sb.WriteString(meshFloatScript)

	return h.Div(
		g.Attr("data-component", "mesh-wrap"),
		h.Style("overflow-x:auto"),
		g.Raw(sb.String()),
	)
}

// meshFloatScript gives mesh nodes a gentle idle breathing pulse and pulses
// the online-status dots. Scale-only (no translate) so edges, which are
// static SVG lines anchored to each node's rest position, never visibly
// detach. Purely decorative — the layout itself is fixed, server-computed
// SVG; this is the same Motion One library already loaded for the homepage
// hero, used the same way as ui/special/logo_scatter.go.
const meshFloatScript = `<script>(function(root){
  if(!root||!window.Motion) return;
  root.querySelectorAll('[data-mesh="breathe-host"]').forEach(function(el){
    Motion.animate(el,{scale:[1,1.045,1]},
      {duration:2.6+Math.random()*1.2,easing:'ease-in-out',repeat:Infinity,delay:Math.random()*0.8});
  });
  root.querySelectorAll('[data-mesh="pulse"]').forEach(function(el){
    Motion.animate(el,{opacity:[1,0.35,1]},
      {duration:1.6+Math.random()*0.8,repeat:Infinity,delay:Math.random()*0.6});
  });
})(document.currentScript.previousElementSibling)</script>`

// simulateMesh runs a fixed number of basic spring/repulsion steps in place.
// Deterministic given deterministic input positions (no global rand use).
func simulateMesh(sim []meshSimNode, edges []meshEdge) {
	const steps = 280
	const repelK = 2200.0
	const centerK = 0.01
	const damping = 0.82

	for range steps {
		fx := make([]float64, len(sim))
		fy := make([]float64, len(sim))

		for i := range sim {
			for j := i + 1; j < len(sim); j++ {
				dx := sim[i].x - sim[j].x
				dy := sim[i].y - sim[j].y
				d2 := dx*dx + dy*dy
				if d2 < 1 {
					d2 = 1
				}
				f := repelK / d2
				d := math.Sqrt(d2)
				ux, uy := dx/d, dy/d
				fx[i] += ux * f
				fy[i] += uy * f
				fx[j] -= ux * f
				fy[j] -= uy * f
			}
		}

		for _, e := range edges {
			a, b := sim[e.from], sim[e.to]
			dx := b.x - a.x
			dy := b.y - a.y
			d := math.Sqrt(dx*dx + dy*dy)
			if d < 0.01 {
				d = 0.01
			}
			diff := (d - e.ideal) * e.strength
			ux, uy := dx/d, dy/d
			fx[e.from] += ux * diff
			fy[e.from] += uy * diff
			fx[e.to] -= ux * diff
			fy[e.to] -= uy * diff
		}

		for i := range sim {
			if sim[i].kind == meshKindHub {
				sim[i].x, sim[i].y = 0, 0
				sim[i].vx, sim[i].vy = 0, 0
				continue
			}
			fx[i] += -sim[i].x * centerK
			fy[i] += -sim[i].y * centerK
			sim[i].vx = (sim[i].vx + fx[i]) * damping
			sim[i].vy = (sim[i].vy + fy[i]) * damping
			sim[i].x += sim[i].vx
			sim[i].y += sim[i].vy
		}
	}
}

// meshFit computes a center offset and uniform scale so the simulated layout
// (centered on the hub at origin) fits within a size x size canvas with
// margin for node radii and labels.
func meshFit(sim []meshSimNode, size float64) (cx, cy, scale float64) {
	maxExtent := 1.0
	for _, s := range sim {
		extent := math.Max(math.Abs(s.x), math.Abs(s.y)) + s.r + 20
		if extent > maxExtent {
			maxExtent = extent
		}
	}
	scale = (size / 2) / maxExtent
	return size / 2, size / 2, scale
}

func findHostOS(nodes []MeshNode, name string) string {
	for _, n := range nodes {
		if n.Name == name {
			return n.OS
		}
	}
	return ""
}

func serviceCount(nodes []MeshNode, name string) int {
	for _, n := range nodes {
		if n.Name == name {
			return len(n.Services)
		}
	}
	return 0
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
