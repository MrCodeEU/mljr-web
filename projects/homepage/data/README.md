# Homepage data

`seed-cache.json` is the single embedded data source for the homepage
(GitHub projects, LinkedIn profile, Strava aggregates). It is currently a
**dummy/seed snapshot** maintained by hand. The plan is to grow this into a
separate data repo with an automated refresh pipeline.

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

| Source   | State                                                        |
|----------|--------------------------------------------------------------|
| Strava   | `scrapers/strava.go` — working OAuth refresh-token client; fetches activities, aggregates stats/disciplines |
| GitHub   | manual — repo list curated by hand                           |
| LinkedIn | manual — no usable API; export is copied in by hand          |

## Automation plan

1. **Split into a data repo** (`mljr-data`): JSON output + fetcher binaries.
   The homepage embeds the latest committed `seed-cache.json` at build time —
   no runtime dependency on third-party APIs, site stays a single binary.
2. **Strava (automate first, scraper already exists)**: scheduled GitHub
   Action (weekly) runs the scraper with `STRAVA_CLIENT_ID/SECRET/
   REFRESH_TOKEN` repo secrets, writes `strava_data`, commits on change.
   Only public-safe aggregates are stored — no GPS traces, no start points.
3. **GitHub stats**: same Action queries the GraphQL API
   (`contributionsCollection` for the real heatmap, repo list for
   stars/languages). Replaces `placeholderContributions()` in
   `pages/github.go` and the `*`-marked sample counters.
4. **LinkedIn**: stays manual (ToS). Keep a small YAML/JSON file edited by
   hand; the Action merges it into the cache.
5. **Trigger**: data repo commit fires `repository_dispatch` → homepage image
   rebuild → deploy via the homelab-automation Ansible flow.

Until the pipeline lands: edit `seed-cache.json` directly and rebuild.
