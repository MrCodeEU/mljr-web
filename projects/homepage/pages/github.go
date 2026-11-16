package pages

import (
	"fmt"
	"time"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	hpdata "mljr-web/projects/homepage/data"
	"mljr-web/ui/icon"
	"mljr-web/ui/swiss"
)

// githubSection ("Quelltext"): the commit calendar as a full-width poster,
// language share as one proportional bar, and four headline figures. Falls
// back to sample data until the GitHub stats pipeline has synced.
func githubSection(p swiss.SectionProps, d hpdata.SiteData, lang string) g.Node {
	stars := 0
	for _, pr := range d.GitHub {
		stars += pr.Stars
	}
	stats := d.GitHubStats
	sample := stats == nil

	var days []dayCount
	commits, streak := 1200, 23
	langs := []segment{{"Go", 58}, {"TypeScript / Svelte", 22}, {"Other", 20}}
	if sample {
		days = sampleContributions()
	} else {
		for _, c := range stats.Contributions {
			if t, err := time.Parse("2006-01-02", c.Date); err == nil {
				days = append(days, dayCount{t, c.Count})
			}
		}
		commits, streak = stats.CommitsYear, stats.LongestStreak
		langs = langs[:0]
		// top five, the long tail folded into "Other"
		other := 0.0
		for i, s := range stats.LanguageShare {
			if i < 5 {
				langs = append(langs, segment{s.Name, s.Pct})
			} else {
				other += s.Pct
			}
		}
		if other > 0 {
			langs = append(langs, segment{tr(lang, "Übrige", "Other"), other})
		}
	}
	total := 0
	for _, c := range days {
		total += c.Count
	}
	mark := ""
	if sample {
		mark = "*"
	}
	p.Note = tr(lang, "Live-Daten aus der GitHub-Pipeline.", "Live data from the GitHub pipeline.")
	if sample {
		p.Note = tr(lang, "* Beispieldaten — die GitHub-Pipeline hat noch nicht synchronisiert.", "* Sample data — the GitHub pipeline hasn't synced yet.")
	}
	return swiss.Section(p,
		swiss.Cell(swiss.CellProps{},
			swiss.Kicker("A. "+tr(lang, "Kalender · 12 Monate", "Calendar · 12 months"), g.Text(fmt.Sprintf("%d %s%s", total, tr(lang, "Beiträge", "contributions"), mark))),
			calendarSVG(days, "var(--accent)", tr(lang, "Beiträge", "contributions"), lang),
		),
		swiss.Cell(swiss.CellProps{Variant: "flush"},
			figRow([]figure{
				{tr(lang, "Repos", "Repos"), fmt.Sprint(len(d.GitHub)), tr(lang, "öffentlich", "public"), "", "lucide:folder-git-2"},
				{tr(lang, "Sterne", "Stars"), fmt.Sprint(stars), "GitHub", "", "lucide:star"},
				{tr(lang, "Commits / Jahr", "Commits / year"), fmt.Sprint(commits) + mark, tr(lang, "letzte 12 Monate", "last 12 months"), "block", "lucide:git-commit-horizontal"},
				{tr(lang, "Längste Serie", "Longest streak"), fmt.Sprintf("%d%s", streak, mark), tr(lang, "Tage in Folge", "days in a row"), "", "lucide:flame"},
			}),
		),
		swiss.Cell(swiss.CellProps{},
			swiss.Kicker("B. "+tr(lang, "Sprachanteil", "Language share"), g.Text(tr(lang, "nach Codezeilen", "by lines of code")+mark)),
			propBar(langs, "%"),
			h.Div(h.Class("cat-links cell-foot"),
				swiss.BtnLink("https://github.com/MrCodeEU", true, []g.Node{g.Attr("target", "_blank"), g.Attr("rel", "noopener noreferrer")},
					icon.Icon("simple-icons:github", icon.Props{Size: "1em"}), g.Text(tr(lang, "Profil öffnen ↗", "Open profile ↗"))),
			),
		),
	)
}

// sampleContributions: deterministic plausible year of commits.
func sampleContributions() []dayCount {
	seed := uint64(0x9e3779b97f4a7c15)
	next := func() float64 {
		seed ^= seed << 13
		seed ^= seed >> 7
		seed ^= seed << 17
		return float64(seed%1000) / 1000
	}
	now := time.Now()
	out := make([]dayCount, 0, 365)
	for i := 364; i >= 0; i-- {
		date := now.AddDate(0, 0, -i)
		base := 3.0
		if wd := date.Weekday(); wd == time.Saturday || wd == time.Sunday {
			base = 1.2
		}
		switch (i / 7) % 9 {
		case 2:
			base *= 0.2
		case 5, 6:
			base *= 2.1
		}
		v := int(base * next() * 4)
		if next() < 0.18 {
			v = 0
		}
		out = append(out, dayCount{date, v})
	}
	return out
}
