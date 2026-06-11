package pages

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	hpdata "mljr-web/projects/homepage/data"
	"mljr-web/ui/icon"
	"mljr-web/ui/primitive"
)

// statsSection renders the Number Ticker metrics + Marquee tech strip
// placed between the hero and experience section.
func statsSection(d hpdata.SiteData) g.Node {
	jobs := d.LinkedIn.RelevantExperience(100)
	companies := map[string]bool{}
	for _, j := range jobs {
		companies[j.Company] = true
	}

	// Marquee items: tech logos
	techItems := []struct{ ic, label string }{
		{"simple-icons:go", "Go"},
		{"simple-icons:rust", "Rust"},
		{"simple-icons:typescript", "TypeScript"},
		{"simple-icons:python", "Python"},
		{"simple-icons:docker", "Docker"},
		{"simple-icons:linux", "Linux"},
		{"simple-icons:ansible", "Ansible"},
		{"simple-icons:tailscale", "Tailscale"},
		{"simple-icons:grafana", "Grafana"},
		{"simple-icons:svelte", "Svelte"},
		{"simple-icons:kotlin", "Kotlin"},
		{"lucide:shield", "Security"},
		{"lucide:network", "Networking"},
		{"lucide:server", "Homelab"},
		{"lucide:brain", "IAM / PAM"},
		{"lucide:cpu", "Embedded"},
		{"simple-icons:postgresql", "PostgreSQL"},
		{"simple-icons:sqlite", "SQLite"},
	}
	statTones := []struct{ bg, ink string }{
		{"var(--yellow-bg,#fef08a)", "var(--yellow-ink,#713f12)"},
		{"var(--cyan-bg,#a5f3fc)", "var(--cyan-ink,#164e63)"},
		{"var(--lime-bg,#d9f99d)", "var(--lime-ink,#365314)"},
		{"var(--violet-bg,#ddd6fe)", "var(--violet-ink,#3b0764)"},
		{"var(--pink-bg,#fbcfe8)", "var(--pink-ink,#831843)"},
		{"var(--sky-bg,#bae6fd)", "var(--sky-ink,#0c4a6e)"},
		{"var(--mint-bg,#d1fae5)", "var(--mint-ink,#064e3b)"},
		{"var(--blush-bg,#fde8e8)", "var(--blush-ink,#7f1d1d)"},
	}
	marqueeItems := make([]g.Node, len(techItems))
	for i, t := range techItems {
		tv := statTones[i%len(statTones)]
		marqueeItems[i] = h.Div(
			h.Style("display:flex;align-items:center;gap:var(--sp-2);padding:var(--sp-2) var(--sp-4);border:var(--bw-2) solid var(--ink);border-radius:calc(var(--radius)*2);font-size:var(--t-sm);font-weight:700;white-space:nowrap;background:"+tv.bg+";color:"+tv.ink),
			icon.Icon(t.ic, icon.Props{Size: "1rem"}),
			g.Text(t.label),
		)
	}

	return h.Div(
		h.Style("border-top:var(--bw-2) solid var(--ink);border-bottom:var(--bw-2) solid var(--ink)"),
		// Stats row
		h.Div(
			h.Class("hero-stat-grid"),
			h.Style("display:grid;grid-template-columns:repeat(4,1fr);gap:0;margin-bottom:0"),
			// dividers between stat tiles via border-right
			statTile("Yrs coding", primitive.NumberTickerProps{Value: 8, Suffix: "+", TriggerOnView: true, ID: "nt2-yrs"}),
			statTile("Projects", primitive.NumberTickerProps{Value: float64(len(d.GitHub)), TriggerOnView: true, ID: "nt2-proj"}),
			statTile("Companies", primitive.NumberTickerProps{Value: float64(len(companies)), TriggerOnView: true, ID: "nt2-comp"}),
			statTile("Degrees", primitive.NumberTickerProps{Value: float64(len(d.LinkedIn.Education)), TriggerOnView: true, ID: "nt2-edu"}),
		),
		// Marquee strip
		h.Div(
			h.Style("border-top:var(--bw-1) solid var(--line);padding:var(--sp-3) 0"),
			primitive.Marquee(primitive.MarqueeProps{
				Speed:        "30s",
				PauseOnHover: true,
				Gap:          "var(--sp-3)",
			}, marqueeItems...),
		),
	)
}

func statTile(label string, p primitive.NumberTickerProps) g.Node {
	if p.Duration == 0 {
		p.Duration = 2500
	}
	return h.Div(
		h.Class("hero-stat-tile"),
		h.Style("text-align:center;padding:var(--sp-6) var(--sp-4);display:flex;flex-direction:column;align-items:center;justify-content:center;border-right:var(--bw-1) solid var(--line)"),
		h.Div(h.Style("font-size:clamp(2.5rem,4vw,3.8rem);font-weight:900;line-height:1"),
			primitive.NumberTicker(p),
		),
		h.Div(h.Style("font-size:var(--t-xs);color:var(--muted);font-weight:800;text-transform:uppercase;letter-spacing:.1em;margin-top:var(--sp-2)"),
			g.Text(label),
		),
	)
}
