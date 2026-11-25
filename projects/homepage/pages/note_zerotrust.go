package pages

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	hpdata "mljr-web/projects/homepage/data"
	"mljr-web/ui/layout"
	"mljr-web/ui/special"
	"mljr-web/ui/swiss"
	"mljr-web/ui/token"
)

// ZeroTrustNote is the field note on the homelab's move from a VPN mesh to
// an OpenZiti zero-trust overlay: why, the architecture, the gated phases
// (from content.json), what went wrong on the way and the trade-offs.
// Facts come from homelab-automation's docs/ZERO_TRUST_PLAN.md. Accepted
// risks are described in general terms only, without host details.
func ZeroTrustNote(d hpdata.SiteData, lang string, a AnalyticsConfig) g.Node {
	home := tr(lang, "/de/", "/")
	var phases g.Node
	if mig := d.ContentFor(lang).Migration; mig != nil && len(mig.Phases) > 0 {
		phases = g.Group([]g.Node{
			swiss.Micro(h.Class("zn-status"), g.Text(migrationStatus(mig, lang))),
			migrationList(mig, lang),
		})
	}

	type item struct{ title, text string }
	list := func(items []item) g.Node {
		rows := make([]g.Node, len(items))
		for i, it := range items {
			rows[i] = h.Li(h.Class("zn-item"),
				h.B(g.Text(it.title)),
				h.P(g.Text(it.text)),
			)
		}
		return h.Ol(h.Class("zn-list"), g.Group(rows))
	}

	architecture := []item{
		{tr(lang, "Controller am einzigen öffentlichen Host", "Controller on the only public host"),
			tr(lang, "Der Ziti-Controller und ein öffentlicher Edge-Router laufen auf dem VPS, mit eigenen Ports neben Caddy: Keiner hängt vom anderen ab.",
				"The Ziti controller and a public edge router run on the VPS, on their own ports next to Caddy, so neither depends on the other.")},
		{tr(lang, "Daheim nur ausgehend", "Outbound only at home"),
			tr(lang, "Die Router auf dem NUC, dem Unraid-NAS und dem Backup-NAS verbinden sich nach außen zum VPS. Daheim lauscht kein Port im Internet.",
				"The routers on the NUC, the Unraid NAS and the backup NAS connect out to the VPS. Nothing at home listens on the internet.")},
		{tr(lang, "Root-CA offline", "Offline root CA"),
			tr(lang, "Auf dem VPS liegt nur eine Zwischen-CA. Die Laufzeiten von Root, Zwischen-CA und Router-Zertifikaten werden täglich geprüft.",
				"Only an intermediate CA lives on the VPS. The lifetimes of the root, the intermediate and the router certificates are checked daily.")},
		{tr(lang, "Policies als Code, Default-Deny", "Policies as code, default deny"),
			tr(lang, "OpenTofu erzeugt Ziti-Dienste und Regeln aus demselben Dienstkatalog, aus dem OpenVox deployt. Was keine Regel erlaubt, ist nicht erreichbar.",
				"OpenTofu generates the Ziti services and rules from the same service catalog OpenVox deploys from. Whatever no rule allows is unreachable.")},
		{tr(lang, "Menschen über einen Identity Provider", "People through one identity provider"),
			tr(lang, "Web-Apps laufen weiter über Caddy mit OIDC, ohne Client auf Handy oder Laptop. SSH meldet sich mit opkssh über denselben IdP an, statt mit Schlüsseln pro Gerät.",
				"Web apps keep running through Caddy with OIDC, with no client on phones or laptops. SSH logs in with opkssh through the same IdP instead of keys per device.")},
		{tr(lang, "Wegwerf-Identitäten", "Disposable identities"),
			tr(lang, "Beim Neuaufbau eines Hosts wird er neu eingeschrieben, statt Identitäten aus Backups zurückzuholen. Sie müssen also nie gesichert werden.",
				"A rebuilt host is enrolled again instead of restoring identities from backups, so they never need backing up.")},
	}

	lessons := []item{
		{tr(lang, "Der Durchsatz lag am Puffer, nicht an der Leitung", "Throughput was a buffer problem, not the line"),
			tr(lang, "Vom NAS zum VPS schaffte der Tunneler nur 53 % des Ausgangswerts, obwohl die Leitung mehr hergab. Die Router-Metriken zeigten rund 450 Neuübertragungen pro Sekunde: Der TCP-Sendepuffer des NAS stapelte Sekunden an Daten, das Overlay sah eine riesige Laufzeit und schickte Daten neu, die nie verloren waren. Ein Router auf dem NAS plus net.ipv4.tcp_notsent_lowat=131072 brachte 97–99 %, im Nachtest 101 % in beide Richtungen.",
				"From the NAS to the VPS the tunneler reached only 53 % of the baseline, although the line had room. The router metrics showed about 450 retransmissions per second: the NAS's TCP send buffer queued seconds of data, the overlay saw a huge round trip and resent data that was never lost. A router on the NAS plus net.ipv4.tcp_notsent_lowat=131072 brought 97–99 %, and 101 % both ways in the retest.")},
		{tr(lang, "Lokaler Verkehr blieb nicht lokal", "Local traffic did not stay local"),
			tr(lang, "Zwischen zwei Geräten im selben LAN liefen 2 von 3 Messungen über den VPS, obwohl die gemessenen Link-Kosten dagegen sprachen. Eine gemeinsame Policy für alle Router reicht nicht; LAN-Dienste brauchen Policies nach Router-Attributen.",
				"Between two devices on the same LAN, 2 of 3 runs went through the VPS, although the measured link costs said otherwise. One shared policy for all routers is not enough; LAN services need policies by router attribute.")},
		{tr(lang, "Docker umgeht die Host-Firewall", "Docker bypasses the host firewall"),
			tr(lang, "Von Docker veröffentlichte Ports laufen an firewalld vorbei. Controller und Router nutzen deshalb das Host-Netz, und die Dienste binden nur an Loopback: Erreichbar sind sie nur noch übers Overlay.",
				"Ports published by Docker go around firewalld. So the controller and the routers use the host network, and the services bind to loopback only: they are reachable over the overlay alone.")},
		{tr(lang, "Verweigerungen erzeugen kein Ereignis", "Denials leave no event"),
			tr(lang, "Ein Client, der einen Dienst nicht wählen darf, bekommt gar keine Sitzung dafür. Alarme zählen deshalb fehlgeschlagene Logins und gescheiterte Verbindungen. Jeder Alarm wurde mit einem echten Ausfall getestet: Controller gestoppt, Router gestoppt, Link gekappt.",
				"A client that may not dial a service is never given a session for it. Alerts therefore count failed logins and failed circuits. Every alert was tested with a real outage: controller stopped, router stopped, link cut.")},
		{tr(lang, "Die Reihenfolge zählt", "Order matters"),
			tr(lang, "Ein Router, dessen erste Anfrage für einen Dienst abgelehnt wird, versucht es bis zum Neustart nicht wieder. Gemeinsame Policies kommen deshalb vor den Diensten.",
				"A router whose first request for a service is rejected does not retry until it restarts. Shared policies therefore come before the services.")},
		{tr(lang, "Nur echte Neustarts finden Boot-Races", "Only real reboots find boot races"),
			tr(lang, "Appliances wie das Unraid-NAS, dessen Root-Dateisystem bei jedem Start neu entsteht, werden mit einem echten Reboot abgenommen. Frühere Tests haben Rennen beim Start nur so gefunden.",
				"Appliances such as the Unraid NAS, whose root filesystem is rebuilt at every boot, are signed off with a real reboot. Earlier tests found start-up races only that way.")},
	}

	tradeoffs := []item{
		{tr(lang, "Der öffentliche Host trägt den Controller", "The public host carries the controller"),
			tr(lang, "Er ist der exponierteste Host. Dagegen: Root-CA offline, Management-API nur lokal, ein Alarm bei jeder neuen Identität.",
				"It is the most exposed host. Against that: an offline root CA, the management API on localhost only, and an alert for every new identity.")},
		{tr(lang, "Die Web-Ebene schrumpft nicht", "The web plane does not shrink"),
			tr(lang, "Der Reverse Proxy erreicht weiter alle Web-Backends. Ziti verengt den Verkehr zwischen Maschinen, nicht den Weg der Browser.",
				"The reverse proxy still reaches every web backend. Ziti narrows traffic between machines, not the path browsers take.")},
		{tr(lang, "Der Identity Provider wird kritisch", "The identity provider becomes critical"),
			tr(lang, "Fällt er aus, kommt man nur noch über den Notfallzugang hinein. Er bekommt deshalb dieselbe Sorgfalt beim Wiederaufbau wie der Controller.",
				"When it is down, only break-glass access gets in. So it gets the same rebuild care as the controller.")},
	}

	title := tr(lang, "Von Tailscale zu OpenZiti", "From Tailscale to OpenZiti")
	return layout.PageShell(
		layout.PageProps{
			Title:       title + " - mljr.eu",
			Description: tr(lang, "Wie das Homelab vom VPN-Mesh auf ein Zero-Trust-Overlay mit OpenZiti umgezogen ist: Aufbau, Phasen mit Gate, Messungen und Lehren.", "How the homelab moved from a VPN mesh to a zero-trust overlay with OpenZiti: architecture, gated phases, measurements and lessons."),
			Theme:       token.ThemeSwiss,
			AutoMode:    true,
			LockTheme:   true,
			Mode:        token.ModeLight,
			Lang:        lang,
			HeadExtra:   append([]g.Node{g.El("style", g.Raw(homepageCSS+leitstandCSS+zeroTrustNoteCSS))}, AnalyticsHead(a)...),
		},
		special.ThemeToggleRoot(token.ThemeSwiss, token.ModeLight),
		siteNavbar(lang),
		h.Main(h.ID("main"),
			h.Div(h.Class("sub-page"),
				subHead(home+"#homelab", tr(lang, "Heimnetz", "Homelab"), 3,
					h.H1(g.Text(title)),
					h.P(h.Class("sub-lede"), g.Text(tr(lang,
						"Wie das Homelab vom VPN-Mesh auf ein Zero-Trust-Overlay umgezogen ist: in Phasen mit Gate, mit Messungen und mit dem, was dabei schiefging.",
						"How the homelab moved from a VPN mesh to a zero-trust overlay: in gated phases, with measurements, and with what went wrong on the way."))),
				),
			),
			subRow("zn-row", tr(lang, "Warum", "Why"), "",
				h.P(h.Class("zn-text"), g.Text(tr(lang,
					"Das Ziel: seitliche Bewegung bremsen. Jede Verbindung zwischen Maschinen wird über eine Identität authentifiziert und pro Dienst erlaubt, alles andere ist verboten. Ein übernommener Host oder Container erreicht nur, was seine Rolle braucht. Und es ist ein Lernprojekt: Alles, was ich selbst dafür schreibe, ist Go.",
					"The goal: slow down lateral movement. Every connection between machines is authenticated by an identity and allowed per service, everything else is denied. A compromised host or container reaches only what its role needs. And it is a learning project: everything I write for it myself is Go."))),
			),
			subRow("zn-row", tr(lang, "Aufbau", "Architecture"), "", list(architecture)),
			subRow("zn-row", tr(lang, "Vorgehen", "Method"), "",
				h.P(h.Class("zn-text"), g.Text(tr(lang,
					"Ein Umbau, zerlegt in Phasen. Jede lässt das Netz voll funktionsfähig, hat einen schriftlichen Rollback und endet mit Weiter, Kurswechsel oder Stopp. Vor der Produktion lief ein Wegwerf-PoC gegen 15 schriftliche Kriterien, von der Offline-Root-CA bis zum SSH-Login über den IdP.",
					"One migration, cut into phases. Each leaves the network fully working, has a written rollback and ends in go, change course or stop. Before production, a throwaway PoC ran against 15 written criteria, from the offline root CA to SSH login through the IdP."))),
				phases,
			),
			subRow("zn-row", tr(lang, "Gelernt", "Lessons"), "", list(lessons)),
			subRow("zn-row", tr(lang, "Kompromisse", "Trade-offs"), "", list(tradeoffs)),
			subRow("zn-row", tr(lang, "Weiter", "More"), "",
				h.Div(h.Class("zn-links"),
					swiss.BtnLink(home+"#homelab", true, nil, g.Text(tr(lang, "Heimnetz live", "Homelab live"))),
					swiss.BtnLink(tr(lang, "/de/projects/homelab-automation", "/projects/homelab-automation"), false, nil, g.Text(tr(lang, "Projekt: homelab-automation", "Project: homelab-automation"))),
					swiss.BtnLink("https://github.com/MrCodeEU/homelab-automation", false, []g.Node{g.Attr("target", "_blank"), g.Attr("rel", "noopener")}, g.Text("GitHub ↗")),
				),
			),
		),
		siteFooter(lang),
	)
}

const zeroTrustNoteCSS = `
.zn-text { margin: 0; max-width: 68ch; font-size: 1.05rem; line-height: 1.55; }
.zn-list { list-style: none; margin: 0; padding: 0; counter-reset: zn; max-width: 80ch; }
.zn-item { counter-increment: zn; display: grid; grid-template-columns: 3rem minmax(0, 1fr); column-gap: 1rem; padding: .8rem 0; border-top: 1px solid var(--rule-subtle); }
.zn-item:first-child { border-top: 0; padding-top: 0; }
.zn-item::before { content: counter(zn, decimal-leading-zero); grid-row: span 2; font-size: 1.6rem; font-weight: 900; letter-spacing: -0.04em; line-height: 1; }
.zn-item b { font-size: 1.05rem; letter-spacing: -0.01em; }
.zn-item p { margin: .3rem 0 0; line-height: 1.55; }
.zn-status { display: block; margin: 1.25rem 0 .5rem; color: var(--main-text); }
.zn-links { display: flex; flex-wrap: wrap; gap: .5rem; }
`
