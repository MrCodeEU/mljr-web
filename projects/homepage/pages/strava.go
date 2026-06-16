package pages

import (
	"fmt"
	"strings"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"mljr-web/internal/i18n"
	hpdata "mljr-web/projects/homepage/data"
	"mljr-web/ui/icon"
	"mljr-web/ui/layout"
	"mljr-web/ui/primitive"
	"mljr-web/ui/token"
)

func stravaSection(num string, d hpdata.SiteData, lang string) g.Node {
	if !d.HasStrava() {
		return nil
	}

	s := d.Strava
	recent := s.RecentActivities
	if len(recent) > 8 {
		recent = recent[:8]
	}

	return h.Section(
		h.ID("activity"),
		h.Style("padding:var(--sp-12) 0;border-top:var(--bw-2) solid var(--ink);border-bottom:var(--bw-2) solid var(--ink)"),
		layout.Container(layout.ContainerProps{},
			sectionHeader(num, i18n.T(lang, "sections.activity.title"), publicActivityBadge(s), token.ToneAccent),
			h.Div(
				h.Class("activity-grid"),
				h.Style("display:flex;flex-direction:column;gap:var(--sp-5)"),
				primitive.Card(primitive.CardProps{Tone: token.ToneCyan},
					h.Div(h.Style("display:flex;align-items:center;justify-content:space-between;gap:var(--sp-3);margin-bottom:var(--sp-5)"),
						h.Div(
							h.Div(h.Style("font-size:var(--t-xs);font-weight:900;text-transform:uppercase;letter-spacing:.1em;opacity:.7"), g.Text("Strava snapshot")),
							h.H3(h.Style("font-size:clamp(1.8rem,3vw,3rem);line-height:1;margin:var(--sp-2) 0 0;font-weight:950"), g.Text("Training as data")),
						),
						h.A(h.Href("https://www.strava.com/athletes/123496455"), g.Attr("target", "_blank"), g.Attr("rel", "noopener noreferrer"),
							primitive.Button(primitive.ButtonProps{Variant: token.Outline, Size: token.SizeSM},
								icon.Icon("lucide:external-link"),
								g.Text("Strava"),
							),
						),
					),
					h.Div(
						h.Class("activity-metrics"),
						h.Style("display:grid;grid-template-columns:repeat(6,minmax(0,1fr));gap:var(--sp-3)"),
						h.Div(h.Style("grid-column:span 3"),
							activityMetric("Distance", fmt.Sprintf("%.1f km", hpdata.DistanceKM(s.YearToDateStats.Distance)), "run · hike · ski", "lucide:route")),
						h.Div(h.Style("grid-column:span 3"),
							activityMetric("Moving time", hpdata.DurationHM(s.YearToDateStats.MovingTime), "logged effort", "lucide:timer")),
						h.Div(h.Style("grid-column:span 2"),
							activityMetric("Sessions", fmt.Sprintf("%d", s.YearToDateStats.Count), "year to date", "lucide:flame")),
						h.Div(h.Style("grid-column:span 2"),
							activityMetric("Elevation", fmt.Sprintf("%.0f m", s.YearToDateStats.ElevationGain), "climbed", "lucide:mountain")),
						g.If(s.AvgHeartrate() > 0,
							h.Div(h.Style("grid-column:span 2"),
								activityMetric("Avg HR", fmt.Sprintf("%.0f bpm", s.AvgHeartrate()), "session-weighted", "lucide:heart-pulse"))),
						g.If(s.YTDCalories > 0,
							h.Div(h.Style("grid-column:span 2"),
								activityMetric("Calories", fmt.Sprintf("%.1fk", s.YTDCalories/1000), "burned (est.)", "lucide:zap"))),
					),
					g.If(len(s.Disciplines) > 0,
						h.Div(h.Style("display:flex;flex-wrap:wrap;justify-content:center;gap:var(--sp-2);margin-top:var(--sp-5)"),
							g.Group(func() []g.Node {
								nodes := make([]g.Node, 0, len(s.Disciplines))
								for _, disc := range s.Disciplines {
									nodes = append(nodes, primitive.Tag(primitive.TagProps{Tone: token.ToneLime, Icon: disciplineIcon(disc.Type)},
										g.Text(fmt.Sprintf("%s · %d", disc.Label, disc.Count)),
									))
								}
								return nodes
							}()),
						),
					),
				),
				primitive.Card(primitive.CardProps{Tone: token.TonePink},
					h.Div(h.Style("display:flex;align-items:center;justify-content:space-between;gap:var(--sp-3);margin-bottom:var(--sp-4)"),
						h.H3(h.Style("font-size:var(--t-xl);font-weight:900;margin:0"), g.Text("Recent public activities")),
						primitive.Tag(primitive.TagProps{Tone: token.ToneYellow}, g.Text("aggregated")),
					),
					h.Div(
						h.Class("activity-list-grid"),
						h.Style("display:grid;grid-template-columns:repeat(2,minmax(0,1fr));gap:var(--sp-3)"),
						g.Group(func() []g.Node {
							nodes := make([]g.Node, 0, len(recent))
							for _, a := range recent {
								nodes = append(nodes, activityRow(a))
							}
							return nodes
						}()),
					),
				),
			),
		),
	)
}

func publicActivityBadge(s hpdata.StravaData) string {
	if s.GeneratedAt != "" {
		return "public aggregates"
	}
	if s.Year != "" {
		return s.Year + " aggregates"
	}
	return "public aggregates"
}

func activityMetric(label, value, sub, ic string) g.Node {
	return h.Div(
		h.Style("border:var(--bw-2) solid var(--ink);background:var(--bg);box-shadow:var(--shadow);padding:var(--sp-4);min-height:120px;display:flex;flex-direction:column;justify-content:space-between"),
		h.Div(h.Style("display:flex;align-items:center;justify-content:space-between;gap:var(--sp-2);font-size:var(--t-xs);font-weight:900;text-transform:uppercase;letter-spacing:.08em;color:var(--muted)"),
			g.Text(label),
			icon.Icon(ic, icon.Props{Size: "1.1rem"}),
		),
		h.Div(
			h.Div(h.Style("font-size:clamp(1.6rem,3vw,2.4rem);font-weight:950;line-height:1"), g.Text(value)),
			h.Div(h.Style("font-size:var(--t-xs);font-weight:800;color:var(--muted);margin-top:var(--sp-1)"), g.Text(sub)),
		),
	)
}

func activityRow(a hpdata.StravaActivity) g.Node {
	title := a.Name
	if strings.TrimSpace(title) == "" {
		title = friendlyName(a.Type)
	}
	// Per-type metadata: distance/pace only where it means something —
	// workouts and weight sessions show duration, HR and calories instead.
	meta := []string{a.DisplayDate()}
	if a.MovingTime > 0 {
		meta = append(meta, hpdata.DurationClock(a.MovingTime))
	}
	if showsDistance(a.Type) && a.Distance > 0 {
		meta = append(meta, fmt.Sprintf("%.1f km", hpdata.DistanceKM(a.Distance)))
	}
	if isRun(a.Type) && a.AveragePace > 0 {
		meta = append(meta, hpdata.PaceLabel(a.AveragePace))
	}
	if a.TotalElevationGain > 0 {
		meta = append(meta, fmt.Sprintf("↑ %.0f m", a.TotalElevationGain))
	}
	if a.AverageHeartrate > 0 {
		meta = append(meta, fmt.Sprintf("♥ %.0f bpm", a.AverageHeartrate))
	}
	if a.Calories > 0 {
		meta = append(meta, fmt.Sprintf("%.0f kcal", a.Calories))
	}

	return h.Div(
		h.Style("display:grid;grid-template-columns:auto 1fr;gap:var(--sp-3);align-items:center;border:var(--bw-1) solid var(--ink);background:var(--bg);padding:var(--sp-3)"),
		h.Div(h.Style("width:42px;height:42px;border:var(--bw-2) solid var(--ink);display:grid;place-items:center;background:var(--accent);color:var(--accent-ink);box-shadow:var(--shadow-sm)"),
			icon.Icon(disciplineIcon(a.Type), icon.Props{Size: "1.2rem"}),
		),
		h.Div(
			h.Div(h.Style("font-weight:900;line-height:1.15"), g.Text(title)),
			h.Div(h.Style("display:flex;flex-wrap:wrap;gap:var(--sp-2);margin-top:var(--sp-1);font-size:var(--t-xs);font-weight:800;color:var(--muted)"),
				g.Text(strings.Join(meta, " · ")),
			),
		),
	)
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

func disciplineIcon(kind string) string {
	switch strings.ToLower(kind) {
	case "run", "running", "trailrun", "virtualrun":
		return "lucide:footprints"
	case "ride", "cycling", "spinning", "virtualride", "mountainbikeride", "gravelride", "ebikeride":
		return "lucide:bike"
	case "hike", "hiking", "walk", "snowshoe":
		return "lucide:mountain"
	case "weighttraining", "strength":
		return "lucide:dumbbell"
	case "workout", "training", "crossfit":
		return "lucide:heart-pulse"
	case "alpineski", "backcountryski", "nordicski", "snowboard", "ski", "skiing":
		return "lucide:mountain-snow"
	case "swim", "swimming":
		return "lucide:waves"
	default:
		return "lucide:activity"
	}
}
