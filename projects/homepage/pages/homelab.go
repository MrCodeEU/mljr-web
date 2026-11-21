package pages

import (
	"fmt"
	"sort"
	"strings"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"mljr-web/projects/homepage/homelab"
	"mljr-web/ui/swiss"
)

// homelabSection ("Heimnetz"): a control-room board ("Leitstand") that the
// /api/homelab SSE endpoint re-renders every 60s, plus the architecture as
// a numbered flowchart.
func homelabSection(p swiss.SectionProps, snap homelab.Snapshot, lang string) g.Node {
	p.Note = tr(lang,
		"Leitstand: alle 60 s live aktualisiert über Uptime Kuma + VictoriaMetrics. Antwortzeiten sind komplette HTTPS-Abrufe über die öffentliche Adresse, inklusive DNS, TLS und Weg über den VPS.",
		"Control room: refreshed live every 60 s from Uptime Kuma + VictoriaMetrics. Response times are full HTTPS requests to the public address, including DNS, TLS and the hop through the VPS.")
	p.Attrs = []g.Node{g.Attr("data-on-interval__duration.60s", "@get('/api/homelab')")}
	return swiss.Section(p,
		HomelabPanel(snap),
		archFlow(lang),
	)
}

// HomelabPanel renders the patchable Leitstand. Exported: the /api/homelab
// fragment handler re-renders it with a fresh snapshot. Labels are
// bilingual (German / English) because the SSE handler has no language.
func HomelabPanel(snap homelab.Snapshot) g.Node {
	if !snap.KumaOK && snap.ActiveBans < 0 {
		return h.Div(h.ID("homelab-panel"), h.Class("lab-panel"),
			swiss.Cell(swiss.CellProps{},
				swiss.Kicker("Leitstand / Control room", g.Text("◷")),
				h.P(h.Class("cell-text"), g.Text("Telemetrie startet — Live-Daten erscheinen nach der ersten Abfrage. / Telemetry warming up.")),
			),
		)
	}
	num := func(v int) string {
		if v < 0 {
			return "—"
		}
		return fmt.Sprint(v)
	}
	total := len(snap.Services)
	avg := 0.0
	if n := len(snap.PingHistory); n > 0 {
		avg = snap.PingHistory[n-1]
	}
	upVariant := "block"
	if snap.UpCount < total {
		upVariant = "second"
	}
	onlineMesh := 0
	for _, m := range snap.Mesh {
		if m.Online {
			onlineMesh++
		}
	}

	// service matrix: one square per monitor
	svcs := append([]homelab.Service(nil), snap.Services...)
	sort.Slice(svcs, func(i, j int) bool { return svcs[i].Name < svcs[j].Name })
	tiles := make([]g.Node, len(svcs))
	bars := make([]barItem, 0, len(svcs))
	for i, s := range svcs {
		name, host := splitHost(s.Name)
		state := "up"
		if !s.Up {
			state = "down"
		}
		ping := "—"
		if s.Ping > 0 {
			ping = fmt.Sprintf("%.0f ms", s.Ping)
			bars = append(bars, barItem{Label: name, Value: ping, N: s.Ping})
		}
		tiles[i] = h.Li(h.Class("lt-tile lt-"+state), h.Title(s.Name),
			h.Span(h.Class("lt-state"), g.Text(map[string]string{"up": "● in Betrieb", "down": "▲ gestört"}[state])),
			h.B(g.Text(name)),
			h.Span(h.Class("lt-meta"), g.Text(host+" · "+ping)),
		)
	}
	sort.Slice(bars, func(i, j int) bool { return bars[i].N > bars[j].N })
	if len(bars) > 8 {
		bars = bars[:8]
	}

	threats := make([]barItem, 0, len(snap.TopThreats))
	for _, t := range snap.TopThreats {
		threats = append(threats, barItem{Label: threatLabel(t.Name), Value: fmt.Sprint(t.Value), N: float64(t.Value), Variant: "third"})
	}
	var attackDays []dayCount
	for _, d := range snap.AttackDays {
		attackDays = append(attackDays, dayCount{d.Date, d.Count})
	}
	// slowest list beside the mesh table; either takes the full row alone
	sSpan, mSpan := 4, 5
	if len(snap.Mesh) == 0 {
		sSpan = 9
	}
	// threat list beside the attack heatmap; either takes the full row alone
	tSpan, hSpan := 3, 6
	if len(attackDays) == 0 {
		tSpan = 9
	}
	if len(threats) == 0 {
		hSpan = 9
	}
	meshRows := make([]g.Node, len(snap.Mesh))
	for i, m := range snap.Mesh {
		st := "● online"
		if !m.Online {
			st = "○ offline"
		} else if m.Relay {
			st = "◐ relay"
		}
		meshRows[i] = h.Tr(h.Td(h.B(g.Text(m.Name))), h.Td(g.Text(m.OS)), h.Td(g.Text(st)), h.Td(g.Text(fmt.Sprint(len(m.Services)))))
	}

	return h.Div(h.ID("homelab-panel"), h.Class("lab-panel"),
		swiss.Cell(swiss.CellProps{Variant: "flush"},
			figRow([]figure{
				{"Dienste / Services", fmt.Sprintf("%d/%d", snap.UpCount, total), "in Betrieb / up", upVariant, "lucide:server"},
				{"Antwortzeit / Latency", fmt.Sprintf("%.0f", avg), "ms · Ø HTTPS-Abruf / avg HTTPS request", "", "lucide:clock"},
				{"Sperren / Bans", num(snap.ActiveBans), "CrowdSec aktiv / active", "third", "lucide:shield"},
				{"Angriffe / Attacks", num(snap.Attacks24h), "letzte 24 h / last 24 h", "", "lucide:alert-triangle"},
				{"Hosts", num(snap.HostsOnline), "online", "", "lucide:network"},
				{"Alarme / Alerts", num(snap.SecurityEvents), "seit Start / since start", "", "lucide:bell"},
			}),
		),
		swiss.Cell(swiss.CellProps{},
			swiss.Kicker("A. Dienstmatrix / Service matrix", h.Span(h.Class("lt-live"), g.Text("● live · "+snap.FetchedAt.Format("15:04:05")))),
			h.Ul(h.Class("lt-grid"), g.Group(tiles)),
		),
		swiss.Cell(swiss.CellProps{Span: sSpan},
			swiss.Kicker("B. Langsamste Antworten / Slowest responses", g.Text("HTTPS · extern / external")),
			barList(bars),
		),
		g.If(len(meshRows) > 0, swiss.Cell(swiss.CellProps{Span: mSpan},
			swiss.Kicker("C. Tailscale-Mesh", g.Text(fmt.Sprintf("%d/%d online", onlineMesh, len(snap.Mesh)))),
			h.Table(h.Class("lt-table"),
				h.THead(h.Tr(h.Th(g.Text("Host")), h.Th(g.Text("OS")), h.Th(g.Text("Status")), h.Th(g.Text("Dienste / Services")))),
				h.TBody(meshRows...),
			),
		)),
		g.If(len(threats) > 0, swiss.Cell(swiss.CellProps{Span: tSpan},
			swiss.Kicker("D. Abwehr / Defence", g.Text(fmt.Sprintf("%s lokal · %s Community", num(snap.BansLocal), num(snap.BansCommunity)))),
			barList(threats),
		)),
		g.If(len(attackDays) > 0, swiss.Cell(swiss.CellProps{Span: hSpan, Attrs: []g.Node{h.Class("lt-heat")}},
			swiss.Kicker("E. Angriffs-Heatmap / Attack heatmap", g.Text(fmt.Sprintf("CrowdSec · %s %s", num(sumDays(attackDays)), "geblockt / blocked"))),
			calendarSVG(attackDays, "var(--third)", "Angriffe / attacks", "de"),
		)),
	)
}

func sumDays(days []dayCount) int {
	n := 0
	for _, d := range days {
		n += d.Count
	}
	return n
}

// splitHost turns "glance [mljr]" into ("glance", "mljr").
func splitHost(s string) (string, string) {
	if i := strings.LastIndex(s, " ["); i > 0 && strings.HasSuffix(s, "]") {
		return s[:i], s[i+2 : len(s)-1]
	}
	return s, ""
}

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

// archFlow: the request path and the provisioning loop as a numbered
// flowchart. No map — the two locations are named in the boxes.
func archFlow(lang string) g.Node {
	type box struct {
		num, title, role, place string
		items                   []string
		variant                 string
	}
	steps := []box{
		{"1", "Internet", tr(lang, "Besucher", "Visitors"), "", []string{"HTTPS :443"}, "ink"},
		{"2", "mljr · VPS", tr(lang, "Öffentlicher Eingang", "Public entry"), "Frankfurt", []string{"Caddy", "CrowdSec", "Authelia", tr(lang, "öffentliche Apps", "public apps")}, "block"},
		{"3", "Tailscale", tr(lang, "Verschlüsseltes Mesh", "Encrypted mesh"), "WireGuard", []string{tr(lang, "keine offenen Ports daheim", "no open ports at home")}, ""},
		{"4", "nuc", tr(lang, "Heimserver", "Home server"), "Thaling", []string{"VictoriaMetrics", "Grafana", "Loki", tr(lang, "interne Dienste", "internal services")}, ""},
		{"5", "nas", "Unraid NAS", "Thaling", []string{"Storage", "Backups", "Media"}, ""},
	}
	nodes := make([]g.Node, 0, len(steps)*2)
	for i, s := range steps {
		if i > 0 {
			nodes = append(nodes, h.Span(h.Class("af-arrow"), g.Attr("aria-hidden", "true"), g.Text("→")))
		}
		pills := make([]g.Node, len(s.items))
		for j, it := range s.items {
			pills[j] = swiss.Pill(g.Text(it))
		}
		cls := "af-box"
		if s.variant != "" {
			cls += " mono-" + s.variant
		}
		nodes = append(nodes, h.Div(h.Class(cls),
			h.Span(h.Class("af-num"), g.Text(s.num)),
			h.B(h.Class("af-title"), g.Text(s.title)),
			swiss.Micro(g.Text(s.role+g2(s.place))),
			swiss.Pills(pills...),
		))
	}
	return swiss.Cell(swiss.CellProps{},
		swiss.Kicker(tr(lang, "Aufbau", "Architecture"), g.Text("100% IaC")),
		h.Div(h.Class("af"), g.Group(nodes)),
		h.Div(h.Class("af-loop mono-second"),
			h.Span(h.Class("af-num"), g.Text("6")),
			h.P(g.Text(tr(lang,
				"OpenVox-Agenten (Puppet) auf jedem Host wenden dasselbe Repo an — Härtung, Docker-Compose-Dienste, Caddy-Snippets, Monitoring. Kein zentraler Server, kein SSH-Fan-out, keine Schneeflocken.",
				"OpenVox (Puppet) agents on every host apply the same repo — hardening, Docker Compose services, Caddy snippets, monitoring. No central server, no SSH fan-out, no snowflakes."))),
		),
	)
}

func g2(place string) string {
	if place == "" {
		return ""
	}
	return " · " + place
}

const leitstandCSS = `
.lab-panel { display: grid; grid-template-columns: repeat(9, minmax(0, 1fr)); grid-column: 1 / -1; }
.lab-panel > * { grid-column: 1 / -1; }
.lab-panel > [data-span] { grid-column: auto / span var(--span); }
.lt-grid { list-style: none; margin: 0; padding: 0; display: grid; grid-template-columns: repeat(auto-fill, minmax(124px, 1fr)); gap: 3px; }
.lt-tile { min-height: 84px; gap: .4rem; padding: .55rem; display: flex; flex-direction: column; justify-content: space-between; min-width: 0; }
.lt-tile b { font-size: .95rem; letter-spacing: -0.01em; line-height: 1.1; overflow-wrap: anywhere; }
.lt-up { background: var(--block); color: var(--block-fg); }
.lt-down { background: var(--second); color: var(--second-fg); }
.lt-state, .lt-meta { font-family: var(--font-mono); font-size: .64rem; font-weight: 700; text-transform: uppercase; letter-spacing: .05em; }
.lt-meta { opacity: .8; }
.lt-down .lt-meta { opacity: 1; }
.lt-live { color: var(--main-text); }
.lt-sub { display: block; margin: 1rem 0 .4rem; color: var(--muted); }
.lt-table { width: 100%; border-collapse: collapse; font-size: .8rem; }
.lt-table th { text-align: left; font-family: var(--font-mono); font-size: .64rem; text-transform: uppercase; letter-spacing: .06em; border-bottom: 1px solid var(--ink); padding: .25rem .3rem .25rem 0; }
.lt-table td { border-bottom: 1px solid var(--rule-subtle); padding: .3rem .3rem .3rem 0; }
.af { display: flex; align-items: stretch; gap: .4rem; overflow-x: auto; padding-bottom: .25rem; }
.af-box { flex: 1 1 0; min-width: 150px; border: 1px solid var(--ink); padding: .75rem; display: flex; flex-direction: column; gap: .45rem; }
.af-num { font-size: 2.4rem; font-weight: 900; letter-spacing: -0.05em; line-height: .85; }
.af-title { font-size: 1.05rem; letter-spacing: -0.02em; }
.af-arrow { align-self: center; font-size: 1.6rem; font-weight: 900; }
.mono-block [data-component="swiss-pill"], .mono-ink [data-component="swiss-pill"], .mono-second [data-component="swiss-pill"], .mono-third [data-component="swiss-pill"] { border-color: currentColor; }
.af-loop { display: flex; gap: 1rem; align-items: center; padding: .75rem; margin-top: .4rem; }
.af-loop p { margin: 0; font-size: .9rem; font-weight: 600; }
@media (max-width: 1024px) {
  .lab-panel { grid-template-columns: minmax(0, 1fr); }
  .lab-panel > [data-span] { grid-column: 1 / -1; }
  .lt-grid { grid-template-columns: repeat(2, minmax(0, 1fr)); }
  .af { flex-direction: column; }
  .af-arrow { transform: rotate(90deg); }
}
`
