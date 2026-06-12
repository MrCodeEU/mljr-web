# Homepage Data Repo Dummy

This folder sketches the separate data-only repository shape for the homepage.
`generated/site-data.json` mirrors the current homepage seed-cache payload and
is now the default runtime data file for local/dev homepage runs. The embedded
`projects/homepage/data/seed-cache.json` remains as a fallback so a broken or
missing data checkout does not prevent the site binary from starting.

## Suggested Files

- `profile.json`: identity, avatar, public links, location, short bio.
- `timeline.json`: work, education, HTL, thesis, and curated milestones.
- `projects.json`: curated projects and manual overrides for GitHub data.
- `generated/github.json`: scraped GitHub repositories and metadata.
- `generated/strava.json`: public Strava aggregates, no maps or GPS traces.
- `generated/site-data.json`: merged, versioned payload consumed by the homepage.
- `assets/`: checked-in avatars, logos, thesis images, and project screenshots.
- `schemas/site-data.schema.json`: JSON Schema contract for the generated
  homepage artifact.

## Contract

`generated/site-data.json` must validate against
`schemas/site-data.schema.json`.

Required top-level fields:

- `schema_version`: currently `site-data.v1`.
- `generated_at`: RFC3339 timestamp for the generated artifact.
- `github_projects`: merged GitHub/project cards.
- `linkedin_data`: public profile, experience, education, and skills.
- `strava_data`: public aggregate activity data only.

Generators should validate the file before committing it. The homepage also
parses the same file with the Go `SiteData` types and keeps the previous valid
data if a reload fails.

## Homepage Runtime

The homepage reads `HOMEPAGE_DATA_FILE` at startup and periodically checks its
mtime. By default this points at `data-repo-dummy/generated/site-data.json`.
`HOMEPAGE_DATA_RELOAD_SECONDS` controls the check interval and defaults to 300.

That means the future deploy flow can update or replace the data checkout and
the already-running web process will pick up the new JSON on the next reload
check. No homepage rebuild is required for JSON-only content changes.

## Scheduled Jobs

- Strava: daily or weekly job using repository secrets for OAuth refresh.
- GitHub: daily job with ETag caching and manual project overrides.
- Merge: deterministic generator that validates schema versions and writes
  `generated/site-data.json`.
