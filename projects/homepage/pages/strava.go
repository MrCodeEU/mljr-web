package pages

import (
	"fmt"
	"strings"
	"time"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	hpdata "mljr-web/projects/homepage/data"
	"mljr-web/ui/icon"
	"mljr-web/ui/swiss"
)

func stravaSection(p swiss.SectionProps, d hpdata.SiteData, lang string) g.Node {
	if !d.HasStrava() {
		return nil
	}
	s := d.Strava
	ytd := s.YearToDateStats
	year := s.Year
	if year == "" {
		year = fmt.Sprint(time.Now().Year())
	}
	p.Note = tr(lang, "Öffentliche Aggregate aus Strava und intervals.icu. Keine Gesundheitsdaten.", "Public aggregates from Strava and intervals.icu. No health data.")

	figs := []figure{
		{tr(lang, "Distanz ", "Distance ") + year, fmt.Sprintf("%.0f km", hpdata.DistanceKM(ytd.Distance)), tr(lang, "Laufen · Rad · Berg", "run · ride · hike"), "block", "lucide:route"},
		{tr(lang, "Zeit in Bewegung", "Moving time"), fmt.Sprintf("%d h", ytd.MovingTime/3600), tr(lang, "seit Jänner", "since January"), "", "lucide:timer"},
		{tr(lang, "Einheiten", "Sessions"), fmt.Sprint(ytd.Count), tr(lang, "dieses Jahr", "this year"), "", "lucide:activity"},
		{tr(lang, "Höhenmeter", "Elevation"), fmt.Sprintf("%.0f m", ytd.ElevationGain), tr(lang, "aufwärts", "climbed"), "third", "lucide:mountain"},
	}
	if hr := s.AvgHeartrate(); hr > 0 {
		figs = append(figs, figure{tr(lang, "Ø Puls", "Avg HR"), fmt.Sprintf("%.0f bpm", hr), tr(lang, "gewichtet nach Einheiten", "session-weighted"), "", "lucide:heart-pulse"})
	}

	// monthly columns: distance per month, January ticks labelled with the year
	months := s.MonthlyActivities
	if len(months) > 36 {
		months = months[len(months)-36:]
	}
	cols := make([]colPoint, len(months))
	for i, m := range months {
		t, _ := time.Parse("2006-01", m.Month)
		cols[i] = colPoint{
			Label: t.Format("2006"),
			Tick:  t.Month() == time.January,
			Value: hpdata.DistanceKM(m.Distance),
			Title: fmt.Sprintf("%s · %.0f km · %d×", t.Format("01/2006"), hpdata.DistanceKM(m.Distance), m.Count),
		}
	}

	// personal records as a timetable
	prRows := make([]g.Node, 0, len(s.PersonalRecords))
	for _, r := range s.PersonalRecords {
		prRows = append(prRows, h.Tr(
			h.Td(h.B(g.Text(recordLabel(r.Type, lang)))),
			h.Td(h.Class("bz-time"), g.Text(hpdata.DurationClock(r.Time))),
			h.Td(g.Text(paceFor(r.Time, r.Distance))),
			h.Td(g.Text(shortDate(r.Date))),
		))
	}
	best := s.BestActivities
	type rec struct {
		label, value string
		a            hpdata.StravaActivity
	}
	recs := []rec{
		{tr(lang, "Längste Distanz", "Longest distance"), fmt.Sprintf("%.1f km", hpdata.DistanceKM(best.LongestDistance.Distance)), best.LongestDistance},
		{tr(lang, "Längste Dauer", "Longest duration"), hpdata.DurationHM(best.LongestTime.MovingTime), best.LongestTime},
		{tr(lang, "Schnellste Pace", "Fastest pace"), hpdata.PaceLabel(best.FastestPace.AveragePace), best.FastestPace},
		{tr(lang, "Meiste Höhenmeter", "Most elevation"), fmt.Sprintf("%.0f m", best.MostElevation.TotalElevationGain), best.MostElevation},
	}
	recRows := make([]g.Node, 0, len(recs))
	for _, r := range recs {
		if r.a.MovingTime == 0 {
			continue
		}
		recRows = append(recRows, h.Tr(h.Td(g.Text(r.label)), h.Td(h.Class("bz-time"), g.Text(r.value)), h.Td(h.Class("bz-sport"), h.Span(h.Class("bz-ico"), icon.Icon(sportIcon(r.a.Type), icon.Props{Size: "1rem"})), g.Text(friendlyName(r.a.Type))), h.Td(g.Text(r.a.DisplayDate()))))
	}

	// disciplines as a proportional bar by time
	var discs []segment
	for _, ds := range s.Disciplines {
		discs = append(discs, segment{fmt.Sprintf("%s · %d×", ds.Label, ds.Count), float64(ds.TotalTime)})
	}

	recent := s.RecentActivities
	if len(recent) > 5 {
		recent = recent[:5]
	}
	actRows := make([]g.Node, len(recent))
	for i, a := range recent {
		title := strings.TrimSpace(a.Name)
		if title == "" {
			title = friendlyName(a.Type)
		}
		dist, pace := "—", "—"
		if showsDistance(a.Type) && a.Distance > 0 {
			dist = fmt.Sprintf("%.1f km", hpdata.DistanceKM(a.Distance))
		}
		if isRun(a.Type) && a.AveragePace > 0 {
			pace = hpdata.PaceLabel(a.AveragePace)
		}
		actRows[i] = h.Tr(
			h.Td(g.Text(a.DisplayDate())),
			h.Td(h.Class("bz-sport"), h.Span(h.Class("bz-ico"), icon.Icon(sportIcon(a.Type), icon.Props{Size: "1rem"})), h.Span(h.B(g.Text(title)), h.Span(h.Class("bz-kind"), g.Text(friendlyName(a.Type))))),
			h.Td(h.Class("bz-time"), g.Text(hpdata.DurationClock(a.MovingTime))),
			h.Td(g.Text(dist)),
			h.Td(h.Class("bz-hide"), g.Text(pace)),
			h.Td(h.Class("bz-hide"), g.Text(orDash(a.TotalElevationGain, "↑ %.0f m"))),
			h.Td(h.Class("bz-hide"), g.Text(orDash(a.AverageHeartrate, "♥ %.0f"))),
			h.Td(h.Class("bz-hide"), g.Text(orDash(a.Calories, "%.0f kcal"))),
		)
	}

	return swiss.Section(p,
		swiss.Cell(swiss.CellProps{Variant: "flush"}, figRow(figs)),
		g.If(len(cols) > 0, swiss.Cell(swiss.CellProps{},
			swiss.Kicker("A. "+tr(lang, "Kilometer pro Monat", "Kilometres per month"), g.Text(fmt.Sprintf("%d %s", len(cols), tr(lang, "Monate", "months")))),
			colChart(cols, "var(--accent)", 180),
		)),
		g.If(len(prRows) > 0, swiss.Cell(swiss.CellProps{Span: 4, Variant: "ink"},
			swiss.Kicker("B. "+tr(lang, "Bestzeiten", "Best times"), g.Text(tr(lang, "Laufen", "Running"))),
			h.Table(h.Class("bz"),
				h.THead(h.Tr(h.Th(g.Text(tr(lang, "Strecke", "Distance"))), h.Th(g.Text(tr(lang, "Zeit", "Time"))), h.Th(g.Text("Pace")), h.Th(g.Text(tr(lang, "Datum", "Date"))))),
				h.TBody(prRows...),
			),
		)),
		g.If(len(recRows) > 0, swiss.Cell(swiss.CellProps{Span: 5},
			swiss.Kicker("C. "+tr(lang, "Rekorde", "Records"), g.Text(tr(lang, "alle Disziplinen", "all disciplines"))),
			h.Table(h.Class("bz"), h.TBody(recRows...)),
			g.If(len(discs) > 0, h.Div(h.Class("bz-disc"),
				swiss.Micro(g.Text(tr(lang, "Disziplinen nach Zeit", "Disciplines by time"))),
				propBar(discs, "%"),
			)),
		)),
		g.If(len(actRows) > 0, swiss.Cell(swiss.CellProps{},
			swiss.Kicker("D. "+tr(lang, "Letzte Einheiten", "Latest sessions"), h.A(h.Href("https://www.strava.com/athletes/123496455"), g.Attr("target", "_blank"), g.Attr("rel", "noopener noreferrer"), g.Text("Strava ↗"))),
			h.Table(h.Class("bz bz-acts"),
				h.THead(h.Tr(h.Th(g.Text(tr(lang, "Datum", "Date"))), h.Th(g.Text(tr(lang, "Einheit", "Session"))), h.Th(g.Text(tr(lang, "Dauer", "Duration"))), h.Th(g.Text(tr(lang, "Distanz", "Distance"))), h.Th(h.Class("bz-hide"), g.Text("Pace")), h.Th(h.Class("bz-hide"), g.Text(tr(lang, "Höhe", "Elev."))), h.Th(h.Class("bz-hide"), g.Text(tr(lang, "Puls", "HR"))), h.Th(h.Class("bz-hide"), g.Text("kcal")))),
				h.TBody(actRows...),
			),
		)),
		fitnessCell(d.WellnessData, lang),
	)
}

func recordLabel(kind, lang string) string {
	switch kind {
	case "half_marathon":
		return tr(lang, "Halbmarathon", "Half marathon")
	case "marathon":
		return "Marathon"
	}
	return strings.ToUpper(strings.ReplaceAll(kind, "_", " "))
}

func paceFor(sec int, meters float64) string {
	if meters <= 0 {
		return "—"
	}
	return hpdata.PaceLabel(float64(sec) / 60 / (meters / 1000))
}

func shortDate(iso string) string {
	t, err := time.Parse("2006-01-02", iso)
	if err != nil {
		return iso
	}
	return t.Format("02.01.2006")
}

// fitnessCell: training load (CTL/ATL/form) for the last 90 days.
func fitnessCell(days []hpdata.WellnessDay, lang string) g.Node {
	if len(days) < 2 {
		return nil
	}
	if len(days) > 90 {
		days = days[len(days)-90:]
	}
	ctl := make([]float64, len(days))
	atl := make([]float64, len(days))
	form := make([]float64, len(days))
	for i, day := range days {
		ctl[i], atl[i], form[i] = day.CTL, day.ATL, day.Form
	}
	latest := days[len(days)-1]
	return swiss.Cell(swiss.CellProps{},
		swiss.Kicker("E. "+tr(lang, "Trainingsbelastung · 90 Tage", "Training load · 90 days"), g.Text(fmt.Sprintf("Form %+.1f", latest.Form))),
		lineChart([]lineSeries{
			{Label: tr(lang, "Fitness (CTL)", "Fitness (CTL)"), Color: "var(--accent)", Points: ctl},
			{Label: tr(lang, "Ermüdung (ATL)", "Fatigue (ATL)"), Color: "var(--second)", Points: atl},
			{Label: "Form", Color: "var(--ink)", Points: form},
		}, 200, "training load"),
	)
}

// orDash formats v, or "—" when there is no value.
func orDash(v float64, format string) string {
	if v <= 0 {
		return "—"
	}
	return fmt.Sprintf(format, v)
}

// sportIcon maps a Strava activity type to a lucide icon.
func sportIcon(kind string) string {
	switch k := strings.ToLower(kind); {
	case isRun(k):
		return "lucide:footprints"
	case strings.Contains(k, "ride"), strings.Contains(k, "bike"):
		return "lucide:bike"
	case strings.Contains(k, "hike"), strings.Contains(k, "walk"):
		return "lucide:mountain"
	case strings.Contains(k, "swim"):
		return "lucide:waves"
	case strings.Contains(k, "weight"), strings.Contains(k, "workout"), strings.Contains(k, "training"), strings.Contains(k, "crossfit"):
		return "lucide:dumbbell"
	}
	return "lucide:activity"
}

func isRun(kind string) bool {
	switch strings.ToLower(kind) {
	case "run", "running", "trailrun", "virtualrun":
		return true
	}
	return false
}

// showsDistance reports whether distance is meaningful for the activity type
// (stationary workouts and weight sessions have none worth showing).
func showsDistance(kind string) bool {
	switch strings.ToLower(kind) {
	case "workout", "training", "crossfit", "weighttraining", "yoga":
		return false
	}
	return true
}

const bewegungCSS = `
.bz-sport { display: flex; gap: .6rem; align-items: center; }
.bz-ico { flex: none; width: 1.8rem; height: 1.8rem; display: grid; place-items: center; border-radius: 50%; background: var(--second); color: var(--second-fg); }
.bz { width: 100%; border-collapse: collapse; font-size: .9rem; }
.bz th { text-align: left; font-family: var(--font-mono); font-size: .64rem; font-weight: 700; text-transform: uppercase; letter-spacing: .06em; padding: 0 .6rem .4rem 0; border-bottom: 1px solid currentColor; }
.bz td { padding: .45rem .6rem .45rem 0; border-bottom: 1px solid var(--rule-subtle); vertical-align: baseline; }
[data-variant="ink"] .bz td { border-bottom-color: color-mix(in srgb, currentColor 25%, transparent); }
.bz tr:last-child td { border-bottom: 0; }
.bz-time { font-family: var(--font-mono); font-weight: 700; font-variant-numeric: tabular-nums; }
.bz td:last-child, .bz th:last-child { padding-right: 0; }
.bz-kind { margin-left: .5rem; font-family: var(--font-mono); font-size: .64rem; text-transform: uppercase; letter-spacing: .05em; color: var(--muted); }
.bz-acts td:first-child { font-family: var(--font-mono); font-size: .78rem; white-space: nowrap; }
.bz-disc { margin-top: 1.25rem; }
.bz-disc [data-component="swiss-micro"] { display: block; margin-bottom: .4rem; }
@media (max-width: 640px) { .bz-hide { display: none; } .bz-kind { display: block; margin: 0; } }
`
