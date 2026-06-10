# Homepage Data Repo Dummy

This folder sketches the separate data-only repository shape for the homepage.
The homepage should consume `generated/site-data.json`; scheduled jobs can update
the generated files without coupling API credentials to the public web server.

## Suggested Files

- `profile.json`: identity, avatar, public links, location, short bio.
- `timeline.json`: work, education, HTL, thesis, and curated milestones.
- `projects.json`: curated projects and manual overrides for GitHub data.
- `generated/github.json`: scraped GitHub repositories and metadata.
- `generated/strava.json`: public Strava aggregates, no maps or GPS traces.
- `generated/site-data.json`: merged, versioned payload consumed by the homepage.
- `assets/`: checked-in avatars, logos, thesis images, and project screenshots.

## Scheduled Jobs

- Strava: daily or weekly job using repository secrets for OAuth refresh.
- GitHub: daily job with ETag caching and manual project overrides.
- Merge: deterministic generator that validates schema versions and writes
  `generated/site-data.json`.
