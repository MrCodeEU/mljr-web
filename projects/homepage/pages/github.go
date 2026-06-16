package pages

import (
	"fmt"
	"time"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"mljr-web/internal/i18n"
	hpdata "mljr-web/projects/homepage/data"
	uidata "mljr-web/ui/data"
	"mljr-web/ui/icon"
	"mljr-web/ui/layout"
	"mljr-web/ui/primitive"
	"mljr-web/ui/token"
)

// githubSection renders the open-source activity panel: a contribution
// heatmap, language-share gauges, and headline counters, all sourced from
// the GitHub stats pipeline (mljr-data). Falls back to sample data if
// GitHubStats hasn't been synced yet (e.g. first boot before the data sync
// timer runs).
func githubSection(d hpdata.SiteData, lang string) g.Node {
	repoCount := len(d.GitHub)
	stars := 0
	for _, p := range d.GitHub {
		stars += p.Stars
	}

	stats := d.GitHubStats
	badge := "sample data"
	if stats != nil {
		badge = "live data"
	}

	contributions := placeholderContributions()
	commitsYear := 1200.0
	commitsSuffix := "+"
	streak := 23.0
	streakSuffix := "d"
	gauges := []g.Node{
		uidata.Gauge(uidata.GaugeProps{Value: 58, Label: "Go", Unit: "%", Size: 130}),
		uidata.Gauge(uidata.GaugeProps{Value: 22, Label: "TS / Svelte", Unit: "%", Size: 130, Color: "var(--info)"}),
		uidata.Gauge(uidata.GaugeProps{Value: 20, Label: "Other", Unit: "%", Size: 130, Color: "var(--warning)"}),
	}
	if stats != nil {
		contributions = realContributions(stats.Contributions)
		commitsYear = float64(stats.CommitsYear)
		commitsSuffix = ""
		streak = float64(stats.LongestStreak)
		streakSuffix = "d"
		gauges = realLanguageGauges(stats.LanguageShare)
	}

	commitsLabel := "Commits / year"
	streakLabel := "Longest streak"
	if stats == nil {
		commitsLabel += "*"
		streakLabel += "*"
	}

	return h.Section(
		h.ID("opensource"),
		h.Style("padding:var(--sp-12) 0;border-top:var(--bw-2) solid var(--ink)"),
		layout.Container(layout.ContainerProps{},
			sectionHeader("03", i18n.T(lang, "sections.opensource.title"), badge, token.ToneMint),
			h.Div(
				h.Class("oss-grid"),
				h.Style("display:grid;grid-template-columns:1.4fr 1fr;gap:var(--sp-5);align-items:stretch"),
				// Left: heatmap card
				primitive.Card(primitive.CardProps{Tone: token.ToneMint},
					h.Div(h.Style("display:flex;align-items:center;justify-content:space-between;gap:var(--sp-3);margin-bottom:var(--sp-4)"),
						h.Div(
							h.Div(h.Style("font-size:var(--t-xs);font-weight:900;text-transform:uppercase;letter-spacing:.1em;opacity:.7"), g.Text("Contributions")),
							h.H3(h.Style("font-size:var(--t-xl);font-weight:900;margin:var(--sp-1) 0 0"), g.Text("A year of commits")),
						),
						h.A(h.Href("https://github.com/MrCodeEU"), g.Attr("target", "_blank"), g.Attr("rel", "noopener"),
							primitive.Button(primitive.ButtonProps{Variant: token.Outline, Size: token.SizeSM},
								icon.Icon("simple-icons:github"),
								g.Text("Profile"),
							),
						),
					),
					uidata.Heatmap(uidata.HeatmapProps{
						Weeks: 52, CellSize: 11, Gap: 3,
						ShowMonthLabels: true, ShowDayLabels: true,
					}, contributions),
					// Language share gauges below the heatmap
					h.Div(
						h.Style("display:grid;grid-template-columns:repeat(3,minmax(0,1fr));gap:var(--sp-2);margin-top:var(--sp-4);justify-items:center"),
						g.Group(gauges),
					),
				),
				// Right: counters — same gap as the outer grid so card edges align
				// with the heatmap card across both columns.
				h.Div(
					h.Style("display:grid;grid-template-columns:repeat(2,minmax(0,1fr));grid-template-rows:1fr 1fr;gap:var(--sp-5)"),
					ossStat("nt-oss-repos", float64(repoCount), "", "Public repos", token.ToneCyan, "lucide:folder-git-2"),
					ossStat("nt-oss-stars", float64(stars), "", "GitHub stars", token.ToneYellow, "lucide:star"),
					ossStat("nt-oss-commits", commitsYear, commitsSuffix, commitsLabel, token.ToneViolet, "lucide:git-commit-horizontal"),
					ossStat("nt-oss-streak", streak, streakSuffix, streakLabel, token.TonePink, "lucide:flame"),
				),
			),
			g.If(stats == nil, h.P(h.Style("margin:var(--sp-3) 0 0;font-size:var(--t-xs);color:var(--muted)"),
				g.Text("* heatmap, commit and streak numbers are sample data — live GitHub stats pipeline coming soon."),
			)),
		),
	)
}

// realContributions converts pipeline-provided contribution days into
// heatmap cells.
func realContributions(days []hpdata.ContributionDay) []uidata.HeatmapCell {
	cells := make([]uidata.HeatmapCell, 0, len(days))
	for _, d := range days {
		if d.Count == 0 {
			continue
		}
		date, err := time.Parse("2006-01-02", d.Date)
		if err != nil {
			continue
		}
		cells = append(cells, uidata.HeatmapCell{
			Date:  date,
			Value: d.Count,
			Label: fmt.Sprintf("%d contributions on %s", d.Count, date.Format("Jan 2")),
		})
	}
	return cells
}

// realLanguageGauges renders up to three gauges from pipeline-provided
// language share data, using the same colors as the placeholder gauges.
func realLanguageGauges(shares []hpdata.LanguageShare) []g.Node {
	colors := []string{"", "var(--info)", "var(--warning)"}
	gauges := make([]g.Node, 0, 3)
	for i, s := range shares {
		if i >= len(colors) {
			break
		}
		gauges = append(gauges, uidata.Gauge(uidata.GaugeProps{Value: s.Pct, Label: s.Name, Unit: "%", Size: 130, Color: colors[i]}))
	}
	return gauges
}

func ossStat(id string, value float64, suffix, label string, tone token.Tone, ic string) g.Node {
	return primitive.Card(primitive.CardProps{Tone: tone},
		h.Div(h.Style("display:flex;flex-direction:column;justify-content:center;align-items:flex-start;min-height:150px;height:100%;gap:var(--sp-3);min-width:0"),
			icon.Icon(ic, icon.Props{Size: "1.6rem"}),
			h.Div(
				h.Div(h.Style("font-size:clamp(2rem,3vw,2.8rem);font-weight:900;line-height:1;font-variant-numeric:tabular-nums;max-width:100%;overflow-wrap:anywhere"),
					primitive.NumberTicker(primitive.NumberTickerProps{
						Value: value, Suffix: suffix, TriggerOnView: true, ID: id, Duration: 2600,
					}),
				),
				h.Div(h.Style("font-size:var(--t-xs);font-weight:800;text-transform:uppercase;letter-spacing:.1em;opacity:.7;margin-top:var(--sp-2);overflow-wrap:anywhere"),
					g.Text(label),
				),
			),
		),
	)
}

// placeholderContributions generates a deterministic, plausible-looking year
// of contribution data: weekday-weighted, with hot streaks and quiet weeks.
// Uses a tiny xorshift PRNG so the pattern is stable across renders.
func placeholderContributions() []uidata.HeatmapCell {
	seed := uint64(0x9E3779B97F4A7C15)
	next := func() float64 {
		seed ^= seed << 13
		seed ^= seed >> 7
		seed ^= seed << 17
		return float64(seed%1000) / 1000
	}
	now := time.Now()
	cells := make([]uidata.HeatmapCell, 0, 365)
	for i := 364; i >= 0; i-- {
		date := now.AddDate(0, 0, -i)
		wd := date.Weekday()
		base := 3.0
		if wd == time.Saturday || wd == time.Sunday {
			base = 1.2
		}
		// Quiet weeks and hot streaks based on the week index
		week := i / 7
		switch week % 9 {
		case 2:
			base *= 0.2 // exam / vacation week
		case 5, 6:
			base *= 2.1 // release crunch
		}
		v := int(base * next() * 4)
		if next() < 0.18 {
			v = 0
		}
		if v > 0 {
			cells = append(cells, uidata.HeatmapCell{
				Date:  date,
				Value: v,
				Label: fmt.Sprintf("%d contributions on %s (sample)", v, date.Format("Jan 2")),
			})
		}
	}
	return cells
}
