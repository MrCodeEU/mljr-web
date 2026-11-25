package pages

import (
	"fmt"
	stdhtml "html"
	"math"
	"sort"
	"strings"
	"time"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	hpdata "mljr-web/projects/homepage/data"
	"mljr-web/projects/homepage/homelab"
	"mljr-web/ui/swiss"
)

// homelabSection ("Heimnetz"): a control-room board ("Leitstand") that the
// /api/homelab SSE endpoint re-renders every 60s, the architecture as a
// numbered flowchart, and the zero-trust migration folded away.
func homelabSection(p swiss.SectionProps, snap homelab.Snapshot, mig *hpdata.MigrationContent, lang string) g.Node {
	p.Note = tr(lang,
		"Leitstand: alle 60 s live aktualisiert über Uptime Kuma + VictoriaMetrics, das Mesh aus den Metriken des Ziti-Controllers. Antwortzeiten sind komplette HTTPS-Abrufe über die öffentliche Adresse, inklusive DNS, TLS und Weg über den VPS.",
		"Control room: refreshed live every 60 s from Uptime Kuma + VictoriaMetrics, the mesh from the Ziti controller's metrics. Response times are full HTTPS requests to the public address, including DNS, TLS and the hop through the VPS.")
	p.Attrs = []g.Node{g.Attr("data-on-interval__duration.60s", "@get('/api/homelab')")}
	return swiss.Section(p,
		HomelabPanel(snap),
		archFlow(lang),
		migrationCell(mig, lang),
	)
}

// HomelabPanel renders the patchable Leitstand. Exported: the /api/homelab
// fragment handler re-renders it with a fresh snapshot. Labels are
// bilingual (German / English) because the SSE handler has no language.
func HomelabPanel(snap homelab.Snapshot) g.Node {
	if !snap.KumaOK && snap.ActiveBans < 0 && !snap.MeshOK {
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
	routersUp, linksUp, dark := 0, 0, 0
	for _, m := range snap.Mesh {
		if m.Online {
			routersUp++
		}
		dark += m.Services
	}
	for _, l := range snap.MeshLinks {
		if l.Up {
			linksUp++
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
	// slowest list beside the threat list; either takes the full row alone
	sSpan, tSpan := 4, 5
	if len(threats) == 0 {
		sSpan = 9
	}

	figs := []figure{
		{"Dienste / Services", fmt.Sprintf("%d/%d", snap.UpCount, total), "in Betrieb / up", upVariant, "lucide:server"},
	}
	if snap.MeshOK {
		meshVariant := "block"
		if routersUp < len(snap.Mesh) || linksUp < len(snap.MeshLinks) {
			meshVariant = "second"
		}
		figs = append(figs,
			figure{"Ziti-Router", fmt.Sprintf("%d/%d", routersUp, len(snap.Mesh)), fmt.Sprintf("online · %d/%d Links", linksUp, len(snap.MeshLinks)), meshVariant, "lucide:network"},
			figure{"Dunkle Dienste / Dark services", fmt.Sprint(dark), fmt.Sprintf("nur übers Overlay · %s offene Verbindungen / overlay only · open circuits", num(snap.MeshCircuits)), "", "lucide:lock"},
		)
	} else {
		figs = append(figs,
			figure{"Hosts", num(snap.HostsOnline), "online", "", "lucide:network"},
			figure{"Alarme / Alerts", num(snap.SecurityEvents), "seit Start / since start", "", "lucide:bell"},
		)
	}
	figs = append(figs,
		figure{"Antwortzeit / Latency", fmt.Sprintf("%.0f", avg), "ms · Ø HTTPS-Abruf / avg HTTPS request", "", "lucide:clock"},
		figure{"Sperren / Bans", num(snap.ActiveBans), "CrowdSec aktiv / active", "third", "lucide:shield"},
		figure{"Angriffe / Attacks", num(snap.Attacks24h), "letzte 24 h / last 24 h", "", "lucide:alert-triangle"},
	)

	// Kicker letters follow the cells that are actually shown.
	letter := 'A'
	kick := func(label string, aside g.Node) g.Node {
		k := swiss.Kicker(string(letter)+". "+label, aside)
		letter++
		return k
	}
	var meshCell, matrixCell, slowCell, threatCell, heatCell g.Node
	if snap.MeshOK {
		meshCell = swiss.Cell(swiss.CellProps{},
			kick("OpenZiti-Mesh", g.Text(fmt.Sprintf("%d/%d Router online", routersUp, len(snap.Mesh)))),
			zitiMesh(snap),
		)
	}
	matrixCell = swiss.Cell(swiss.CellProps{},
		kick("Dienstmatrix / Service matrix", h.Span(h.Class("lt-live"), g.Text("● live · "+snap.FetchedAt.Format("15:04:05")))),
		h.Ul(h.Class("lt-grid"), g.Group(tiles)),
	)
	slowCell = swiss.Cell(swiss.CellProps{Span: sSpan},
		kick("Langsamste Antworten / Slowest responses", g.Text("HTTPS · extern / external")),
		barList(bars),
	)
	if len(threats) > 0 {
		threatCell = swiss.Cell(swiss.CellProps{Span: tSpan},
			kick("Abwehr / Defence", g.Text(fmt.Sprintf("%s lokal · %s Community", num(snap.BansLocal), num(snap.BansCommunity)))),
			barList(threats),
		)
	}
	if len(attackDays) > 0 {
		heatCell = swiss.Cell(swiss.CellProps{Attrs: []g.Node{h.Class("lt-heat")}},
			kick("Angriffs-Heatmap / Attack heatmap", g.Text(fmt.Sprintf("CrowdSec · %s %s", num(sumDays(attackDays)), "geblockt / blocked"))),
			calendarSVG(attackDays, "var(--third)", "Angriffe / attacks", "de"),
		)
	}

	return h.Div(h.ID("homelab-panel"), h.Class("lab-panel"),
		swiss.Cell(swiss.CellProps{Variant: "flush"}, figRow(figs)),
		meshCell, matrixCell, slowCell, threatCell, heatCell,
	)
}

// zitiMesh shows the overlay as a topology drawing (controller on the
// public side, home routers below the line, every link with its latency)
// next to a table with the same facts for screen readers and small screens.
func zitiMesh(snap homelab.Snapshot) g.Node {
	rows := make([]g.Node, len(snap.Mesh))
	for i, m := range snap.Mesh {
		role := "Heim · nur ausgehend / home, outbound only"
		if m.Controller {
			role = "Controller · öffentlich / public"
		}
		st := "● online"
		if !m.Online {
			st = "○ offline"
		}
		rows[i] = h.Tr(
			h.Td(h.B(g.Text(m.Name))),
			h.Td(h.Class("zm-role"), g.Text(role)),
			h.Td(h.Class("zm-nowrap"), g.Text(st)),
			h.Td(g.Text(fmt.Sprintf("%d/%d", m.LinksUp, m.LinksTotal))),
			h.Td(g.Text(fmt.Sprint(m.Services))),
		)
	}
	links := make([]string, len(snap.MeshLinks))
	for i, l := range snap.MeshLinks {
		links[i] = l.Source + " ↔ " + l.Target + " " + linkLatency(l)
	}
	return h.Div(h.Class("zm"),
		h.Figure(h.Class("zm-fig"), meshSVG(snap.Mesh, snap.MeshLinks),
			h.FigCaption(swiss.Micro(g.Text("Links wählen sich von daheim nach außen: daheim lauscht kein Port im Internet. / Links dial out from home: nothing at home listens on the internet."))),
		),
		h.Div(
			h.Table(h.Class("lt-table"),
				h.THead(h.Tr(h.Th(g.Text("Router")), h.Th(h.Class("zm-role"), g.Text("Rolle / Role")), h.Th(g.Text("Status")), h.Th(g.Text("Links")), h.Th(g.Text("Dienste / Services")))),
				h.TBody(rows...),
			),
			g.If(len(links) > 0, h.P(h.Class("zm-links"), g.Text(strings.Join(links, " · ")))),
		),
	)
}

func linkLatency(l homelab.MeshLink) string {
	switch {
	case !l.Up:
		return "✕"
	case l.Latency <= 0:
		return "—"
	case l.Latency < time.Millisecond:
		return "<1 ms"
	}
	return fmt.Sprintf("%d ms", l.Latency.Round(time.Millisecond).Milliseconds())
}

// meshSVG draws controllers as squares above the internet line and home
// routers as circles below it. Arrows point the way a link was dialled.
func meshSVG(hosts []homelab.MeshHost, links []homelab.MeshLink) g.Node {
	const w, hgt, boundary = 440.0, 300.0, 112.0
	const sqW, sqH, r = 132.0, 56.0, 30.0
	type pos struct {
		x, y float64
		ctrl bool
	}
	var ctrls, homes []homelab.MeshHost
	for _, m := range hosts {
		if m.Controller {
			ctrls = append(ctrls, m)
		} else {
			homes = append(homes, m)
		}
	}
	at := map[string]pos{}
	for i, m := range ctrls {
		at[m.Name] = pos{w * float64(i+1) / float64(len(ctrls)+1), 62, true}
	}
	homeIdx := map[string]int{}
	for i, m := range homes {
		at[m.Name] = pos{w * float64(i+1) / float64(len(homes)+1), 212, false}
		homeIdx[m.Name] = i
	}
	// edge returns where the segment from a towards b leaves a's shape.
	edge := func(a, b pos) (float64, float64) {
		dx, dy := b.x-a.x, b.y-a.y
		d := dx*dx + dy*dy
		if d == 0 {
			return a.x, a.y
		}
		if !a.ctrl {
			l := math.Sqrt(d)
			return a.x + dx/l*r, a.y + dy/l*r
		}
		// clip to the square's border
		tx, ty := 1e9, 1e9
		if dx != 0 {
			tx = (sqW / 2) / math.Abs(dx)
		}
		if dy != 0 {
			ty = (sqH / 2) / math.Abs(dy)
		}
		t := min(tx, ty)
		return a.x + dx*t, a.y + dy*t
	}

	var b strings.Builder
	fmt.Fprintf(&b, `<svg viewBox="0 0 %.0f %.0f" class="zm-svg" role="img" aria-label="%s">`, w, hgt, stdhtml.EscapeString(meshLabel(hosts, links)))
	b.WriteString(`<defs><marker id="zm-arrow" viewBox="0 0 10 10" refX="9" refY="5" markerWidth="7" markerHeight="7" orient="auto-start-reverse"><path d="M0 0L10 5L0 10z" class="zm-arrowhead"/></marker></defs>`)
	// the internet line
	fmt.Fprintf(&b, `<line x1="0" y1="%.0f" x2="%.0f" y2="%.0f" class="zm-boundary"/>`, boundary, w, boundary)
	fmt.Fprintf(&b, `<text x="0" y="%.0f" class="zm-zone">INTERNET</text>`, boundary-8)
	fmt.Fprintf(&b, `<text x="0" y="%.0f" class="zm-zone">HEIMNETZ / HOME</text>`, boundary+18)

	type label struct {
		x, y float64
		text string
	}
	var labels []label
	for _, l := range links {
		s, okS := at[l.Source]
		t, okT := at[l.Target]
		if !okS || !okT {
			continue
		}
		cls := "zm-link"
		marker := ` marker-end="url(#zm-arrow)"`
		if !l.Up {
			cls += " zm-down"
			marker = ""
		}
		x1, y1 := edge(s, t)
		x2, y2 := edge(t, s)
		i, j := homeIdx[l.Source], homeIdx[l.Target]
		if !s.ctrl && !t.ctrl && (i-j > 1 || j-i > 1) {
			// home routers that are not neighbours: arc below the row
			cy := s.y + 80
			fmt.Fprintf(&b, `<path d="M%.1f %.1f Q%.1f %.1f %.1f %.1f" class="%s"%s/>`, s.x, s.y+r, (s.x+t.x)/2, cy, t.x, t.y+r, cls, marker)
			labels = append(labels, label{(s.x + t.x) / 2, (s.y+r+cy)/2 + 10, linkLatency(l)})
			continue
		}
		fmt.Fprintf(&b, `<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" class="%s"%s/>`, x1, y1, x2, y2, cls, marker)
		ly := (y1+y2)/2 + 4
		if !s.ctrl && !t.ctrl {
			ly = (y1+y2)/2 - 9 // above a link between neighbours
		}
		labels = append(labels, label{(x1 + x2) / 2, ly, linkLatency(l)})
	}
	for _, m := range ctrls {
		p := at[m.Name]
		cls := "zm-ctrl"
		if !m.Online {
			cls += " zm-off"
		}
		fmt.Fprintf(&b, `<g class="%s"><rect x="%.1f" y="%.1f" width="%.0f" height="%.0f"/>`, cls, p.x-sqW/2, p.y-sqH/2, sqW, sqH)
		fmt.Fprintf(&b, `<text x="%.1f" y="%.1f" class="zm-name">%s</text>`, p.x, p.y-2, stdhtml.EscapeString(m.Name))
		fmt.Fprintf(&b, `<text x="%.1f" y="%.1f" class="zm-meta">CONTROLLER</text></g>`, p.x, p.y+16)
		fmt.Fprintf(&b, `<text x="%.1f" y="%.1f" class="zm-count">%d DIENSTE</text>`, p.x, p.y-sqH/2-10, m.Services)
	}
	for _, m := range homes {
		p := at[m.Name]
		cls := "zm-home"
		if !m.Online {
			cls += " zm-off"
		}
		fmt.Fprintf(&b, `<g class="%s"><circle cx="%.1f" cy="%.1f" r="%.0f"/>`, cls, p.x, p.y, r)
		fmt.Fprintf(&b, `<text x="%.1f" y="%.1f" class="zm-name">%s</text></g>`, p.x, p.y+5, stdhtml.EscapeString(m.Name))
		fmt.Fprintf(&b, `<text x="%.1f" y="%.1f" class="zm-count">%d DIENSTE</text>`, p.x, p.y+r+18, m.Services)
	}
	for _, l := range labels {
		fmt.Fprintf(&b, `<text x="%.1f" y="%.1f" class="zm-lat">%s</text>`, l.x, l.y, stdhtml.EscapeString(l.text))
	}
	b.WriteString(`</svg>`)
	return g.Raw(b.String())
}

func meshLabel(hosts []homelab.MeshHost, links []homelab.MeshLink) string {
	up, lup := 0, 0
	for _, m := range hosts {
		if m.Online {
			up++
		}
	}
	for _, l := range links {
		if l.Up {
			lup++
		}
	}
	return fmt.Sprintf("OpenZiti-Mesh: %d/%d Router online, %d/%d Links aktiv / routers online, links up", up, len(hosts), lup, len(links))
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
// flowchart, then the two planes: browsers through Caddy, machines over
// OpenZiti. No map — the two locations are named in the boxes.
func archFlow(lang string) g.Node {
	type box struct {
		num, title, role, place string
		items                   []string
		variant                 string
	}
	steps := []box{
		{"1", "Internet", tr(lang, "Besucher", "Visitors"), "", []string{"HTTPS :443"}, "ink"},
		{"2", "mljr · VPS", tr(lang, "Öffentlicher Eingang", "Public entry"), "Frankfurt", []string{"Caddy", "CrowdSec", "Authelia", tr(lang, "Ziti-Controller", "Ziti controller")}, "block"},
		{"3", "OpenZiti", tr(lang, "Zero-Trust-Overlay", "Zero-trust overlay"), "mTLS", []string{tr(lang, "nur ausgehend von daheim", "outbound only from home"), tr(lang, "Zugriff pro Dienst & Identität", "access per service & identity"), tr(lang, "Default-Deny, Policies als Code", "default deny, policies as code")}, ""},
		{"4", "nuc", tr(lang, "Heimserver", "Home server"), "Thaling", []string{"VictoriaMetrics", "Grafana", "Loki", tr(lang, "interne Dienste", "internal services")}, ""},
		{"5", "nas · ugreen", tr(lang, "Unraid-NAS + Backup-NAS", "Unraid NAS + backup NAS"), "Thaling", []string{"Storage", "Backups", "Media", "Syncthing"}, ""},
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
	plane := func(title, text string) g.Node {
		return h.Div(h.Class("af-plane"),
			swiss.Micro(g.Text(title)),
			h.P(g.Text(text)),
		)
	}
	return swiss.Cell(swiss.CellProps{},
		swiss.Kicker(tr(lang, "Aufbau", "Architecture"), g.Text("100% IaC")),
		h.Div(h.Class("af"), g.Group(nodes)),
		h.Div(h.Class("af-planes"),
			plane(tr(lang, "Web-Ebene", "Web plane"), tr(lang,
				"Browser → Caddy mit Authelia (OIDC). Kein Client nötig, auch nicht für Handy und Familie.",
				"Browser → Caddy with Authelia (OIDC). No client needed, not even for phones and family.")),
			plane(tr(lang, "Maschinen-Ebene", "Machine plane"), tr(lang,
				"Exporter, Backups, SSH und Caddys Upstreams → OpenZiti. Jede Verbindung mit Identität, alles andere verboten.",
				"Exporters, backups, SSH and Caddy's upstreams → OpenZiti. Every connection carries an identity, everything else is denied.")),
		),
		h.Div(h.Class("af-loop mono-second"),
			h.Span(h.Class("af-num"), g.Text("6")),
			h.P(g.Text(tr(lang,
				"OpenVox-Agenten (Puppet) auf jedem Host wenden dasselbe Repo an — Härtung, Docker-Compose-Dienste, Caddy, Ziti-Policies (OpenTofu), Monitoring. CI hat keinen Zugang zu den Hosts, kein SSH-Fan-out, keine Schneeflocken.",
				"OpenVox (Puppet) agents on every host apply the same repo — hardening, Docker Compose services, Caddy, Ziti policies (OpenTofu), monitoring. CI has no access to the hosts, no SSH fan-out, no snowflakes."))),
		),
	)
}

// migrationCell folds the zero-trust migration into one line that opens to
// the full list of gated phases and a link to the field note.
func migrationCell(mig *hpdata.MigrationContent, lang string) g.Node {
	if mig == nil || len(mig.Phases) == 0 {
		return nil
	}
	return swiss.Cell(swiss.CellProps{Variant: "flush"},
		h.Details(h.Class("zt"),
			h.Summary(
				h.Span(h.Class("zt-title"), g.Text(mig.Title)),
				h.Span(h.Class("zt-status"), g.Text(migrationStatus(mig, lang))),
				h.Span(h.Class("zt-open"), g.Attr("aria-hidden", "true"), g.Text("+")),
			),
			h.Div(h.Class("zt-body"),
				g.If(mig.Intro != "", h.P(h.Class("zt-intro"), g.Text(mig.Intro))),
				migrationList(mig, lang),
				h.P(h.Class("zt-more"), swiss.BtnLink(tr(lang, "/de/notes/zero-trust", "/notes/zero-trust"), false, nil,
					g.Text(tr(lang, "Die ganze Geschichte: Von Tailscale zu OpenZiti →", "The whole story: from Tailscale to OpenZiti →")))),
			),
		),
	)
}

// migrationStatus names the running phase, or how many are done.
func migrationStatus(mig *hpdata.MigrationContent, lang string) string {
	done := 0
	for _, ph := range mig.Phases {
		switch ph.Status {
		case "active":
			return fmt.Sprintf(tr(lang, "Phase %s von %s läuft", "Phase %s of %s in progress"), ph.Num, mig.Phases[len(mig.Phases)-1].Num)
		case "done":
			done++
		}
	}
	return fmt.Sprintf(tr(lang, "%d von %d Phasen erledigt", "%d of %d phases done"), done, len(mig.Phases))
}

// migrationList renders the phases with their status and gate date.
func migrationList(mig *hpdata.MigrationContent, lang string) g.Node {
	rows := make([]g.Node, len(mig.Phases))
	for i, ph := range mig.Phases {
		mark := map[string]string{"done": "✓", "active": "●", "planned": "○"}[ph.Status]
		state := map[string]string{
			"done":    tr(lang, "erledigt", "done"),
			"active":  tr(lang, "läuft", "in progress"),
			"planned": tr(lang, "geplant", "planned"),
		}[ph.Status]
		if ph.Date != "" {
			if t, err := time.Parse("2006-01-02", ph.Date); err == nil {
				state += " · " + t.Format("02.01.2006")
			}
		}
		rows[i] = h.Li(h.Class("zt-phase zt-"+ph.Status),
			h.Span(h.Class("zt-num"), g.Text(ph.Num)),
			h.Div(
				h.B(g.Text(ph.Title)),
				g.If(ph.Summary != "", h.P(g.Text(ph.Summary))),
			),
			h.Span(h.Class("zt-state"), g.Text(mark+" "+state)),
		)
	}
	return h.Ol(h.Class("zt-list"), g.Group(rows))
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
.zm { display: grid; grid-template-columns: minmax(0, 5fr) minmax(0, 4fr); gap: 1.5rem 2rem; align-items: start; }
.zm-fig { margin: 0; display: grid; gap: .5rem; }
.zm-svg { display: block; width: 100%; max-width: 560px; height: auto; font-family: var(--font-mono); }
.zm-boundary { stroke: var(--ink); stroke-width: 1; stroke-dasharray: 4 4; }
.zm-zone { fill: var(--muted); font-size: 10px; font-weight: 700; letter-spacing: .08em; }
.zm-link { stroke: var(--ink); stroke-width: 2; fill: none; }
.zm-link.zm-down { stroke: var(--muted); stroke-dasharray: 5 5; }
.zm-arrowhead { fill: var(--ink); }
.zm-ctrl rect { fill: var(--block); }
.zm-ctrl text { fill: var(--block-fg); }
.zm-home circle { fill: var(--second); }
.zm-home text { fill: var(--second-fg); }
.zm-off rect, .zm-off circle { fill: var(--bg); stroke: var(--ink); stroke-width: 2; stroke-dasharray: 5 4; }
.zm-off text { fill: var(--ink); }
.zm-name { text-anchor: middle; font-family: var(--font-display, inherit); font-size: 14px; font-weight: 900; }
.zm-meta { text-anchor: middle; font-size: 9px; font-weight: 700; letter-spacing: .06em; }
.zm-count { text-anchor: middle; fill: var(--ink); font-size: 10px; font-weight: 700; letter-spacing: .05em; }
.zm-lat { text-anchor: middle; fill: var(--ink); font-size: 11px; font-weight: 700; paint-order: stroke; stroke: var(--bg); stroke-width: 4px; stroke-linejoin: round; }
.zm-nowrap { white-space: nowrap; }
.zm-links { margin: .75rem 0 0; font-family: var(--font-mono); font-size: .72rem; color: var(--muted); }
.af { display: flex; align-items: stretch; gap: .4rem; overflow-x: auto; padding-bottom: .25rem; }
.af-box { flex: 1 1 0; min-width: 150px; border: 1px solid var(--ink); padding: .75rem; display: flex; flex-direction: column; gap: .45rem; }
.af-num { font-size: 2.4rem; font-weight: 900; letter-spacing: -0.05em; line-height: .85; }
.af-title { font-size: 1.05rem; letter-spacing: -0.02em; }
.af-arrow { align-self: center; font-size: 1.6rem; font-weight: 900; }
.mono-block [data-component="swiss-pill"], .mono-ink [data-component="swiss-pill"], .mono-second [data-component="swiss-pill"], .mono-third [data-component="swiss-pill"] { border-color: currentColor; }
.af-planes { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: .4rem; margin-top: .4rem; }
.af-plane { border: 1px solid var(--ink); padding: .75rem; display: grid; gap: .35rem; }
.af-plane p { margin: 0; font-size: .9rem; }
.af-loop { display: flex; gap: 1rem; align-items: center; padding: .75rem; margin-top: .4rem; }
.af-loop p { margin: 0; font-size: .9rem; font-weight: 600; }
.zt > summary { display: flex; flex-wrap: wrap; align-items: baseline; gap: .25rem 1.25rem; padding: .9rem var(--swiss-pad); cursor: pointer; list-style: none; }
.zt > summary::-webkit-details-marker { display: none; }
.zt > summary:hover { background: var(--act); color: var(--act-fg); }
.zt > summary:hover * { color: inherit; }
.zt > summary:focus-visible { outline: 3px solid var(--act); outline-offset: -3px; }
.zt-title { font-weight: 900; font-size: 1.05rem; letter-spacing: -0.02em; text-transform: uppercase; }
.zt-status { font-family: var(--font-mono); font-size: .72rem; font-weight: 700; text-transform: uppercase; letter-spacing: .06em; color: var(--main-text); }
.zt-open { margin-left: auto; font-family: var(--font-mono); font-size: 1.2rem; font-weight: 700; line-height: 1; }
.zt[open] .zt-open { transform: rotate(45deg); }
.zt-body { padding: 0 var(--swiss-pad) var(--swiss-pad); }
.zt-intro { margin: 0 0 .75rem; max-width: 70ch; font-size: .9rem; }
.zt-list { list-style: none; margin: 0; padding: 0; }
.zt-phase { display: grid; grid-template-columns: 3rem minmax(0, 1fr) 11rem; gap: 1rem; align-items: baseline; padding: .55rem 0; border-top: 1px solid var(--rule-subtle); }
.zt-phase p { margin: .15rem 0 0; font-size: .85rem; color: var(--muted); }
.zt-num { font-size: 1.6rem; font-weight: 900; letter-spacing: -0.04em; line-height: 1; }
.zt-state { font-family: var(--font-mono); font-size: .68rem; font-weight: 700; text-transform: uppercase; letter-spacing: .05em; }
.zt-active { background: var(--block); color: var(--block-fg); padding-left: .5rem; padding-right: .5rem; }
.zt-active p { color: inherit; }
.zt-planned .zt-num, .zt-planned b { color: var(--muted); }
.zt-more { margin: 1rem 0 0; }
@media (max-width: 1024px) {
  .lab-panel { grid-template-columns: minmax(0, 1fr); }
  .lab-panel > [data-span] { grid-column: 1 / -1; }
  .lt-grid { grid-template-columns: repeat(2, minmax(0, 1fr)); }
  .zm { grid-template-columns: minmax(0, 1fr); }
  .af { flex-direction: column; }
  .af-arrow { transform: rotate(90deg); }
  .af-planes { grid-template-columns: minmax(0, 1fr); }
  .zt-phase { grid-template-columns: 2.4rem minmax(0, 1fr); }
  .zt-state { grid-column: 2; }
}
@media (max-width: 600px) {
  .zm-role { display: none; }
}
`
