package pages

import (
	"fmt"
	"strings"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	hpdata "mljr-web/projects/homepage/data"
	"mljr-web/ui/icon"
	"mljr-web/ui/layout"
	"mljr-web/ui/primitive"
	"mljr-web/ui/token"
)

func stravaSection(d hpdata.SiteData) g.Node {
	if !d.HasStrava() {
		return nil
	}

	s := d.Strava
	recent := s.RecentActivities
	if len(recent) > 5 {
		recent = recent[:5]
	}

	return h.Section(
		h.ID("activity"),
		h.Style("padding:var(--sp-12) 0;border-top:var(--bw-2) solid var(--ink);border-bottom:var(--bw-2) solid var(--ink)"),
		layout.Container(layout.ContainerProps{},
			sectionHeader("06", "Activity", publicActivityBadge(s), token.ToneAccent),
			h.Div(
				h.Class("activity-grid"),
				h.Style("display:grid;grid-template-columns:1.1fr .9fr;gap:var(--sp-5);align-items:stretch"),
				primitive.Card(primitive.CardProps{Tone: token.ToneCyan},
					h.Div(h.Style("display:flex;align-items:center;justify-content:space-between;gap:var(--sp-3);margin-bottom:var(--sp-5)"),
						h.Div(
							h.Div(h.Style("font-size:var(--t-xs);font-weight:900;text-transform:uppercase;letter-spacing:.1em;opacity:.7"), g.Text("Strava snapshot")),
							h.H3(h.Style("font-size:clamp(1.8rem,3vw,3rem);line-height:1;margin:var(--sp-2) 0 0;font-weight:950"), g.Text("Training as data")),
						),
						h.A(h.Href("https://www.strava.com/athletes/mrcode"), g.Attr("target", "_blank"), g.Attr("rel", "noopener noreferrer"),
							primitive.Button(primitive.ButtonProps{Variant: token.Outline, Size: token.SizeSM},
								icon.Icon("lucide:external-link"),
								g.Text("Strava"),
							),
						),
					),
					h.Div(
						h.Style("display:grid;grid-template-columns:repeat(2,minmax(0,1fr));gap:var(--sp-3)"),
						activityMetric("Runs", fmt.Sprintf("%d", s.YearToDateStats.Count), "year to date", "lucide:activity"),
						activityMetric("Distance", fmt.Sprintf("%.1f km", hpdata.DistanceKM(s.YearToDateStats.Distance)), "year to date", "lucide:arrow-right"),
						activityMetric("Moving time", fmt.Sprintf("%.1f h", hpdata.DurationHours(s.YearToDateStats.MovingTime)), "logged effort", "lucide:calendar"),
						activityMetric("Elevation", fmt.Sprintf("%.0f m", s.YearToDateStats.ElevationGain), "climbed", "lucide:bar-chart-2"),
					),
					g.If(len(s.Disciplines) > 0,
						h.Div(h.Style("display:flex;flex-wrap:wrap;gap:var(--sp-2);margin-top:var(--sp-5)"),
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
					h.Div(h.Style("display:flex;flex-direction:column;gap:var(--sp-3)"),
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
	meta := []string{a.DisplayDate(), fmt.Sprintf("%.1f km", hpdata.DistanceKM(a.Distance))}
	if a.AveragePace > 0 {
		meta = append(meta, hpdata.PaceLabel(a.AveragePace))
	} else if a.MovingTime > 0 {
		meta = append(meta, fmt.Sprintf("%.1f h", hpdata.DurationHours(a.MovingTime)))
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

func disciplineIcon(kind string) string {
	switch strings.ToLower(kind) {
	case "run", "running", "trailrun", "virtualrun":
		return "lucide:activity"
	case "ride", "cycling", "virtualride", "mountainbikeride", "gravelride":
		return "lucide:zap"
	case "workout", "training", "weighttraining":
		return "lucide:heart"
	default:
		return "lucide:activity"
	}
}
