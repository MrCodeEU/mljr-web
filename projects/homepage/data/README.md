# Homepage data

`seed-cache.json` is the embedded fallback data source for the homepage
(GitHub projects, LinkedIn profile, Strava aggregates, GitHub stats). It's a
copy of `mljr-data/generated/site-data.json` baked in at build time and only
used as the cold-start default — see the root `agents.md` "Runtime data
refresh" section for how live data is synced without a rebuild.

## Shape

```
github_projects[]   — repo name, desc, stars, language, topics, images, links
linkedin_data       — profile, experience[], education[], skills[]
strava_data         — total/ytd stats, recent_activities[], disciplines[],
                      monthly_activities[], ytd_calories
```

Types live in `types.go`; distances are meters, durations are seconds.
Display conversion happens at render time (`DistanceKM`, `DurationClock`,
`DurationHM`, `PaceLabel`).

## Refresh pipeline — current state

All sources are automated via the `mljr-data` repo's nightly generator
(`generator/cmd/generate`, see its `README.md`):

| Source   | State                                                          |
|----------|-----------------------------------------------------------------|
| Strava   | OAuth refresh-token client; fetches activities, aggregates stats/disciplines |
| GitHub   | GraphQL/REST: per-project stars/language/topics + account-level `github_stats` (contribution heatmap, commit count, language share) |
| LinkedIn | manual — no usable API; `profile.json`/`timeline.json` edited by hand |

## Updating live data

1. Edit `mljr-data/projects.json` (or `profile.json`/`timeline.json`), commit
   and push — the generator workflow regenerates `generated/site-data.json`.
2. The prod systemd timer syncs that file to the homepage container's mounted
   volume within ~15 min; the running process hot-reloads it (no rebuild,
   no redeploy).
3. To update the embedded cold-start fallback, run `make data-update` (or
   copy `mljr-data/generated/site-data.json` to `seed-cache.json`) and rebuild.
