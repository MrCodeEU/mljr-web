package pages

import (
	"fmt"
	"strings"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"mljr-web/projects/homepage/homelab"
	uidata "mljr-web/ui/data"
	"mljr-web/ui/icon"
	"mljr-web/ui/layout"
	"mljr-web/ui/primitive"
	"mljr-web/ui/token"
)

// homelabSection renders the live infrastructure panel. The inner panel
// (#homelab-panel) is re-fetched every 60s via Datastar and patched in place.
func homelabSection(snap homelab.Snapshot) g.Node {
	return h.Section(
		h.ID("homelab"),
		h.Style("padding:var(--sp-12) 0;border-top:var(--bw-2) solid var(--ink)"),
		g.Attr("data-on-interval__duration.60s", "@get('/api/homelab')"),
		layout.Container(layout.ContainerProps{},
			sectionHeader("05", "Homelab", "live via tailscale", token.ToneBlush),
			HomelabPanel(snap),
		),
	)
}

// HomelabPanel renders the patchable panel. Exported: the /api/homelab
// fragment handler re-renders it with a fresh snapshot.
func HomelabPanel(snap homelab.Snapshot) g.Node {
	if !snap.KumaOK && snap.ActiveBans < 0 {
		// No data from any source yet (first poll pending or sources offline).
		return h.Div(h.ID("homelab-panel"),
			primitive.Callout(primitive.CalloutProps{Variant: primitive.CalloutInfo},
				g.Text("Telemetry warming up — live homelab data appears here once the first poll lands."),
			),
			archCard(),
		)
	}

	total := len(snap.Services)

	return h.Div(
		h.ID("homelab-panel"),
		h.Div(
			h.Class("homelab-grid"),
			h.Style("display:grid;grid-template-columns:1.3fr 1fr;gap:var(--sp-5);align-items:stretch"),
			servicesCard(snap, total),
			h.Div(
				h.Style("display:flex;flex-direction:column;gap:var(--sp-4)"),
				pingCard(snap),
				cpuCard(snap),
			),
		),
		h.Div(
			h.Class("homelab-grid"),
			h.Style("display:grid;grid-template-columns:1fr 1fr;gap:var(--sp-5);align-items:stretch;margin-top:var(--sp-5)"),
			crowdsecCard(snap),
			threatsCard(snap),
		),
		attacksHeatmapCard(snap),
		archCard(),
		h.Div(
			h.Style("display:flex;align-items:center;gap:var(--sp-2);margin-top:var(--sp-3);font-size:var(--t-xs);color:var(--muted);font-weight:700"),
			h.Span(h.Style("width:8px;height:8px;border-radius:50%;background:#22c55e;animation:pulse-dot 2s ease infinite;flex-shrink:0")),
			g.Text("live — uptime kuma + victoriametrics over tailscale · updated "+snap.FetchedAt.Format("15:04:05")),
		),
	)
}

func servicesCard(snap homelab.Snapshot, total int) g.Node {
	chips := make([]g.Node, 0, total)
	for _, svc := range snap.Services {
		dotColor := "#22c55e"
		if !svc.Up {
			dotColor = "var(--danger,#ef4444)"
		}
		ping := ""
		if svc.Up && svc.Ping > 0 {
			ping = fmt.Sprintf("%.0f ms", svc.Ping)
		}
		chips = append(chips, h.Div(
			h.Style("display:flex;align-items:center;gap:var(--sp-2);border:var(--bw-1) solid var(--ink);background:var(--bg);padding:var(--sp-2) var(--sp-3);min-width:0"),
			h.Span(h.Style("width:9px;height:9px;border-radius:50%;flex-shrink:0;border:1px solid var(--ink);background:"+dotColor)),
			h.Span(h.Style("font-size:var(--t-xs);font-weight:800;overflow:hidden;text-overflow:ellipsis;white-space:nowrap"), g.Text(svc.Name)),
			g.If(ping != "", h.Span(h.Style("font-size:var(--t-xs);font-family:var(--font-mono,monospace);color:var(--muted);margin-left:auto;flex-shrink:0"), g.Text(ping))),
		))
	}

	upBadgeTone := token.ToneLime
	if snap.UpCount < total {
		upBadgeTone = token.ToneYellow
	}

	return primitive.Card(primitive.CardProps{Tone: token.ToneBlush},
		h.Div(h.Style("display:flex;align-items:center;justify-content:space-between;gap:var(--sp-3);margin-bottom:var(--sp-4)"),
			h.Div(
				h.Div(h.Style("font-size:var(--t-xs);font-weight:900;text-transform:uppercase;letter-spacing:.1em;opacity:.7"), g.Text("Services")),
				h.H3(h.Style("font-size:var(--t-xl);font-weight:900;margin:var(--sp-1) 0 0"), g.Text("Self-hosted fleet")),
			),
			primitive.Tag(primitive.TagProps{Tone: upBadgeTone},
				g.Text(fmt.Sprintf("%d / %d up", snap.UpCount, total))),
		),
		h.Div(
			h.Style("display:grid;grid-template-columns:repeat(auto-fill,minmax(170px,1fr));gap:var(--sp-2)"),
			g.Group(chips),
		),
	)
}

func crowdsecCard(snap homelab.Snapshot) g.Node {
	stat := func(value int, label string) g.Node {
		display := "—"
		if value >= 0 {
			display = fmt.Sprintf("%d", value)
		}
		return h.Div(
			h.Style("border:var(--bw-2) solid var(--ink);background:var(--bg);padding:var(--sp-3);text-align:center"),
			h.Div(h.Style("font-size:clamp(1.6rem,2.5vw,2.2rem);font-weight:900;line-height:1;font-variant-numeric:tabular-nums"), g.Text(display)),
			h.Div(h.Style("font-size:var(--t-xs);font-weight:800;color:var(--muted);text-transform:uppercase;letter-spacing:.06em;margin-top:var(--sp-1)"), g.Text(label)),
		)
	}

	// Community vs local origin split as a stacked bar.
	var originSplit g.Node
	if snap.BansCommunity >= 0 && snap.BansLocal >= 0 && snap.BansCommunity+snap.BansLocal > 0 {
		totalBans := snap.BansCommunity + snap.BansLocal
		commPct := float64(snap.BansCommunity) / float64(totalBans) * 100
		originSplit = h.Div(
			h.Style("margin-top:var(--sp-3)"),
			h.Div(h.Style("display:flex;justify-content:space-between;gap:var(--sp-2);font-size:var(--t-xs);font-weight:800;margin-bottom:var(--sp-1)"),
				h.Span(g.Textf("Community blocklist · %d", snap.BansCommunity)),
				h.Span(h.Style("color:var(--muted)"), g.Textf("caught here · %d", snap.BansLocal)),
			),
			h.Div(
				h.Style("display:flex;height:14px;border:var(--bw-2) solid var(--ink);background:var(--bg);overflow:hidden"),
				h.Div(h.Style(fmt.Sprintf("width:%.1f%%;background:var(--violet-bg,#ddd6fe);border-right:var(--bw-1) solid var(--ink)", commPct))),
				h.Div(h.Style("flex:1;background:var(--accent)")),
			),
		)
	}

	return primitive.Card(primitive.CardProps{Tone: token.ToneViolet},
		h.Div(h.Style("display:flex;align-items:center;gap:var(--sp-2);margin-bottom:var(--sp-3)"),
			icon.Icon("lucide:shield-check", icon.Props{Size: "1.3rem"}),
			h.H3(h.Style("font-size:var(--t-base);font-weight:900;margin:0"), g.Text("CrowdSec perimeter")),
		),
		h.Div(
			h.Style("display:grid;grid-template-columns:repeat(2,minmax(0,1fr));gap:var(--sp-2)"),
			stat(snap.ActiveBans, "Active bans"),
			stat(snap.Attacks24h, "Blocked · 24h"),
			stat(snap.SecurityEvents, "Alerts · total"),
			stat(snap.HostsOnline, "Hosts online"),
		),
		originSplit,
	)
}

// threatsCard ranks what the perimeter is actually blocking right now.
func threatsCard(snap homelab.Snapshot) g.Node {
	if len(snap.TopThreats) == 0 {
		return nil
	}
	maxVal := snap.TopThreats[0].Value
	if maxVal < 1 {
		maxVal = 1
	}
	rows := make([]g.Node, 0, len(snap.TopThreats))
	for _, t := range snap.TopThreats {
		pct := float64(t.Value) / float64(maxVal) * 100
		if pct < 2 {
			pct = 2
		}
		rows = append(rows, h.Div(
			h.Div(h.Style("display:flex;justify-content:space-between;gap:var(--sp-2);font-size:var(--t-xs);font-weight:800;margin-bottom:2px"),
				h.Span(g.Text(threatLabel(t.Name))),
				h.Span(h.Style("font-family:var(--font-mono,monospace);color:var(--muted)"), g.Textf("%d", t.Value)),
			),
			h.Div(h.Style("height:12px;border:var(--bw-1) solid var(--ink);background:var(--bg)"),
				h.Div(h.Style(fmt.Sprintf("width:%.1f%%;height:100%%;background:var(--accent)", pct))),
			),
		))
	}
	return primitive.Card(primitive.CardProps{Tone: token.ToneYellow},
		h.Div(h.Style("display:flex;align-items:center;gap:var(--sp-2);margin-bottom:var(--sp-3)"),
			icon.Icon("lucide:radar", icon.Props{Size: "1.3rem"}),
			h.H3(h.Style("font-size:var(--t-base);font-weight:900;margin:0"), g.Text("Top threats · active bans")),
		),
		h.Div(h.Style("display:flex;flex-direction:column;gap:var(--sp-2)"), g.Group(rows)),
	)
}

// threatLabel maps CrowdSec decision reasons to readable names.
func threatLabel(raw string) string {
	m := map[string]string{
		"http:scan":       "HTTP scanning",
		"http:bruteforce": "HTTP brute force",
		"http:exploit":    "HTTP exploits",
		"http:crawl":      "Aggressive crawling",
		"http:dos":        "HTTP DoS attempts",
		"ssh:bruteforce":  "SSH brute force",
	}
	if v, ok := m[raw]; ok {
		return v
	}
	s := strings.TrimPrefix(raw, "crowdsecurity/")
	s = strings.TrimPrefix(s, "LePresidente/")
	return strings.ReplaceAll(s, "-", " ")
}

func cpuCard(snap homelab.Snapshot) g.Node {
	if len(snap.CPUUtil) < 2 {
		return nil
	}
	last := snap.CPUUtil[len(snap.CPUUtil)-1]
	return primitive.Card(primitive.CardProps{Tone: token.ToneLime},
		h.Div(h.Style("display:flex;align-items:baseline;justify-content:space-between;gap:var(--sp-2);margin-bottom:var(--sp-2)"),
			h.H3(h.Style("font-size:var(--t-base);font-weight:900;margin:0"), g.Text("Avg CPU · all hosts · 24h")),
			h.Span(h.Style("font-family:var(--font-mono,monospace);font-weight:800;font-size:var(--t-sm)"),
				g.Text(fmt.Sprintf("%.0f%%", last))),
		),
		uidata.LineChart(uidata.LineChartProps{
			Height: 70,
			Labels: snap.CPULabels,
			Series: []uidata.LineChartSeries{{
				Points: snap.CPUUtil,
				Color:  "var(--accent)",
				Fill:   true,
			}},
		}),
	)
}

func attacksHeatmapCard(snap homelab.Snapshot) g.Node {
	if len(snap.AttackDays) == 0 {
		return nil
	}
	cells := make([]uidata.HeatmapCell, len(snap.AttackDays))
	total := 0
	for i, d := range snap.AttackDays {
		cells[i] = uidata.HeatmapCell{
			Date:  d.Date,
			Value: d.Count,
			Label: fmt.Sprintf("%d attacks blocked on %s", d.Count, d.Date.Format("Jan 2")),
		}
		total += d.Count
	}
	return h.Div(h.Style("margin-top:var(--sp-5)"),
		primitive.Card(primitive.CardProps{Tone: token.ToneNone},
			h.Div(h.Style("display:flex;align-items:center;justify-content:space-between;gap:var(--sp-3);margin-bottom:var(--sp-3);flex-wrap:wrap"),
				h.Div(
					h.Div(h.Style("font-size:var(--t-xs);font-weight:900;text-transform:uppercase;letter-spacing:.1em;opacity:.7"), g.Text("CrowdSec · last 12 months")),
					h.H3(h.Style("font-size:var(--t-xl);font-weight:900;margin:var(--sp-1) 0 0"), g.Text("Attacks blocked per day")),
				),
				primitive.Tag(primitive.TagProps{Tone: token.ToneViolet},
					g.Text(fmt.Sprintf("%d total · recording since Jun 2026", total))),
			),
			uidata.Heatmap(uidata.HeatmapProps{
				Weeks: 52, CellSize: 11, Gap: 3,
				ShowMonthLabels: true, ShowDayLabels: true,
			}, cells),
		),
	)
}

func pingCard(snap homelab.Snapshot) g.Node {
	if len(snap.PingHistory) < 2 {
		return nil
	}
	last := snap.PingHistory[len(snap.PingHistory)-1]
	return primitive.Card(primitive.CardProps{Tone: token.ToneSky},
		h.Div(h.Style("display:flex;align-items:baseline;justify-content:space-between;gap:var(--sp-2);margin-bottom:var(--sp-2)"),
			h.H3(h.Style("font-size:var(--t-base);font-weight:900;margin:0"), g.Text("Avg response time")),
			h.Span(h.Style("font-family:var(--font-mono,monospace);font-weight:800;font-size:var(--t-sm)"),
				g.Text(fmt.Sprintf("%.0f ms", last))),
		),
		uidata.LineChart(uidata.LineChartProps{
			Height: 70,
			Series: []uidata.LineChartSeries{{
				Points: snap.PingHistory,
				Color:  "var(--accent)",
				Fill:   true,
			}},
		}),
	)
}

// archCard explains how the homelab hangs together: three devices on a
// Tailscale mesh, one public Caddy ingress, everything provisioned by
// Ansible. Static content — no live data needed.
func archCard() g.Node {
	deviceBox := func(ic, name, role string, items []string) g.Node {
		tags := make([]g.Node, len(items))
		for i, it := range items {
			tags[i] = h.Span(
				h.Style("border:var(--bw-1) solid var(--ink);background:var(--bg);padding:1px var(--sp-2);font-size:var(--t-xs);font-weight:700;white-space:nowrap"),
				g.Text(it),
			)
		}
		return h.Div(
			h.Style("border:var(--bw-2) solid var(--ink);background:var(--surface);box-shadow:var(--shadow-sm);padding:var(--sp-4);min-width:0"),
			h.Div(h.Style("display:flex;align-items:center;gap:var(--sp-2);margin-bottom:var(--sp-1)"),
				icon.Icon(ic, icon.Props{Size: "1.3rem"}),
				h.Div(h.Style("font-weight:900;font-size:var(--t-base);line-height:1.2"), g.Text(name)),
			),
			h.Div(h.Style("font-size:var(--t-xs);font-weight:800;color:var(--muted);text-transform:uppercase;letter-spacing:.06em;margin-bottom:var(--sp-2)"), g.Text(role)),
			h.Div(h.Style("display:flex;flex-wrap:wrap;gap:var(--sp-1)"), g.Group(tags)),
		)
	}

	return h.Div(h.Style("margin-top:var(--sp-5)"),
		primitive.Card(primitive.CardProps{Tone: token.ToneNone},
			h.Div(h.Style("display:flex;align-items:center;justify-content:space-between;gap:var(--sp-3);margin-bottom:var(--sp-3);flex-wrap:wrap"),
				h.Div(
					h.Div(h.Style("font-size:var(--t-xs);font-weight:900;text-transform:uppercase;letter-spacing:.1em;opacity:.7"), g.Text("Architecture")),
					h.H3(h.Style("font-size:var(--t-xl);font-weight:900;margin:var(--sp-1) 0 0"), g.Text("How it hangs together")),
				),
				primitive.Tag(primitive.TagProps{Tone: token.ToneLime, Icon: "simple-icons:ansible"}, g.Text("100% IaC")),
			),
			h.P(h.Style("font-size:var(--t-sm);color:var(--muted);line-height:1.6;margin:0 0 var(--sp-4);max-width:78ch"),
				g.Text("All public traffic enters through Caddy on the VPS, with CrowdSec banning attackers at the edge and Authelia guarding private apps. Behind that, three machines talk over an encrypted Tailscale mesh — no open ports at home. Every host, container and config file is declared in one Ansible repo: a single make deploy converges the whole fleet."),
			),
			// Internet → ingress
			h.Div(
				h.Style("display:flex;align-items:center;gap:var(--sp-3);flex-wrap:wrap;margin-bottom:var(--sp-2)"),
				h.Div(h.Style("display:flex;align-items:center;gap:var(--sp-2);border:var(--bw-2) solid var(--ink);background:var(--bg);padding:var(--sp-2) var(--sp-3);font-weight:900;font-size:var(--t-sm)"),
					icon.Icon("lucide:globe", icon.Props{Size: "1.1rem"}),
					g.Text("Internet"),
				),
				h.Span(h.Style("font-family:var(--font-mono,monospace);font-size:var(--t-xs);font-weight:700;color:var(--muted)"), g.Text("→ HTTPS :443 · Caddy ingress · CrowdSec at the edge →")),
			),
			// Tailscale mesh containing the three devices
			h.Div(
				h.Style("border:var(--bw-2) dashed var(--ink);padding:var(--sp-4);position:relative;background:color-mix(in srgb,var(--surface) 60%,transparent)"),
				h.Div(h.Style("position:absolute;top:-11px;left:var(--sp-4);background:var(--surface);padding:0 var(--sp-2);font-size:var(--t-xs);font-weight:900;text-transform:uppercase;letter-spacing:.1em;display:flex;align-items:center;gap:var(--sp-1)"),
					icon.Icon("simple-icons:tailscale", icon.Props{Size: ".9rem"}),
					g.Text("Tailscale mesh · WireGuard"),
				),
				h.Div(
					h.Class("homelab-arch-grid"),
					h.Style("display:grid;grid-template-columns:repeat(3,minmax(0,1fr));gap:var(--sp-3)"),
					deviceBox("lucide:cloud", "mljr", "VPS · public entry",
						[]string{"Caddy", "CrowdSec", "Authelia", "public apps"}),
					deviceBox("lucide:server", "nuc", "home server",
						[]string{"VictoriaMetrics", "Grafana", "Loki", "internal services"}),
					deviceBox("lucide:hard-drive", "nas", "Unraid NAS",
						[]string{"storage", "backups", "media"}),
				),
			),
			// Ansible bar
			h.Div(
				h.Style("display:flex;align-items:center;gap:var(--sp-2);border:var(--bw-2) solid var(--ink);background:var(--lime-bg,#d9f99d);padding:var(--sp-2) var(--sp-3);margin-top:var(--sp-2);font-size:var(--t-xs);font-weight:800"),
				icon.Icon("simple-icons:ansible", icon.Props{Size: "1rem"}),
				g.Text("Ansible provisions all three hosts — inventory, hardening, Docker Compose services, deploy hooks. No snowflakes."),
			),
		),
	)
}
