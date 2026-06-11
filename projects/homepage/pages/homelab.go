package pages

import (
	"fmt"

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
				crowdsecCard(snap),
				pingCard(snap),
				cpuCard(snap),
			),
		),
		attacksHeatmapCard(snap),
		h.Div(
			h.Style("display:flex;align-items:center;gap:var(--sp-2);margin-top:var(--sp-3);font-size:var(--t-xs);color:var(--muted);font-weight:700"),
			h.Span(h.Style("width:8px;height:8px;border-radius:50%;background:#22c55e;animation:pulse-dot 2s ease infinite;flex-shrink:0")),
			g.Text("live — uptime kuma + prometheus over tailscale · updated "+snap.FetchedAt.Format("15:04:05")),
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
	return primitive.Card(primitive.CardProps{Tone: token.ToneViolet},
		h.Div(h.Style("display:flex;align-items:center;gap:var(--sp-2);margin-bottom:var(--sp-3)"),
			icon.Icon("lucide:shield-check", icon.Props{Size: "1.3rem"}),
			h.H3(h.Style("font-size:var(--t-base);font-weight:900;margin:0"), g.Text("CrowdSec perimeter")),
		),
		h.Div(
			h.Style("display:grid;grid-template-columns:repeat(3,1fr);gap:var(--sp-2)"),
			stat(snap.ActiveBans, "Active bans"),
			stat(snap.Attacks24h, "Blocked · 24h"),
			stat(snap.HostsOnline, "Hosts online"),
		),
	)
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
	return h.Div(h.Style("margin-top:var(--sp-4)"),
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
