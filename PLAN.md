# mljr-web v2 — Plan Summary

Single Go module monorepo for `data-*`-driven web stack. No Node at runtime, no CDN, no telemetry.

## Stack

- **Go 1.26** + **Echo** (HTTP), **gomponents v1** (HTML), **Datastar 1.0.x** (reactivity + SSE)
- **Tailwind v4** standalone binary (build-time CSS only)
- **Iconify** icons via build-time Go codegen → `g.Raw(svg)` (gomponents v1 dropped `svg` pkg)
- **altcha** captcha (self-hosted, `altcha-lib-go`)
- **SQLite** via `modernc.org/sqlite` (pure Go, scratch-friendly)

## Layout

```
mljr-web/
├── go.mod                       # single module
├── Makefile                     # PROJECT-parameterized
├── lefthook.yml                 # pre-commit/pre-push
├── .golangci.yml
├── PLAN.md                      # this file
├── agents.md                    # impl guide for agents
├── bin/                         # tailwind binary + builds (gitignored)
├── ui/                          # headless component library (mljr-ui)
│   ├── css/core.css             # design system tokens + component rules
│   ├── token/                   # Go mirror of CSS enums
│   ├── icon/                    # Icon() + generated icons_gen.go
│   ├── primitive/               # Button, Tag, Card, Heading, Display, Icon
│   ├── form/                    # field, input, select…
│   ├── layout/                  # Container, Stack, Grid, Navbar, Footer, PageShell
│   ├── overlay/                 # Modal, Toast, Portal
│   ├── data/                    # tables, lists, stats, server-SVG charts (line/bar/radar/gauge/heatmap)
│   ├── feedback/                # alert, spinner, skeleton
│   ├── special/                 # ThemeToggle, ModeToggle, Captcha, Tour, Confetti, OpenMap
│   └── registry/                # showcase metadata (build tag: showcase)
├── internal/
│   ├── web/                     # Echo bootstrap, render, sse, security, assets
│   ├── config/                  # typed env config (incl. HOMELAB_* panel sources)
│   └── version/                 # ldflags-injected build info
├── projects/
│   ├── homepage/                # portfolio binary (homelab/ = live-panel poller)
│   ├── showcase/                # component catalogue binary
│   └── regex/                   # live RE2 regex tester (Datastar SSE, port 8092)
└── tools/icongen/               # Iconify → Go SVG generator
```

## Theming (v2)

Two **themes** × two **modes** = 4 skins:

- `swissbrut-light`, `swissbrut-dark` — combined Swiss × neobrutalism (Swiss-red dominant, hard black lines, hard shadows)
- `ink-light`, `ink-dark` — refined editorial ink palette

CSS contract: `<html data-theme="…" data-mode="…">`. Components emit **structure only** (`data-component`, `data-variant`, `data-size`, `data-tone`, `data-state`). **Never `class=` in `ui/**.go`** — enforced by `make guard-classes`.

Per-project default theme set by `PageShell`. Runtime toggles rewrite `data-theme` / `data-mode` via Datastar signals.

## Datastar contract

- Signals: `data-signals`, `data-bind-X`, `data-text`, `data-show`, `data-attr`, `data-class`
- Events: `data-on-click`, `data-on-input__debounce.200ms`, `data-on-keydown__window`, `data-on-interval__5s`
- Backend: `@get/@post/@put/@patch/@delete` → server SSE `sse.PatchElements(html)` / `sse.PatchSignals(json)` / `sse.RemoveElements(selector)` / `sse.ExecuteScript(js)` / `sse.Redirect(url)`
- Read client signals: `datastar.ReadSignals(r, &state)`
- Vendored `datastar.js` at `/static/datastar.js` (no CDN)

## Icons

`tools/icongen` reads `tools/icongen/icons.txt` (one `set:name` per line), fetches Iconify `@iconify-json/<set>/icons.json` from jsdelivr, caches under `tools/icongen/.cache/`, emits `ui/icon/icons_gen.go` with `map[string]icon{svg}`. `Icon("lucide:github")` renders inline SVG via `g.Raw` with `fill="currentColor"` + `1em` sizing.

## Per-project Docker

Multi-stage: golang-alpine fetches tailwind, builds CSS into `projects/<P>/assets/static/app.css`, then `CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" ./projects/<P>` → distroless `static-debian12:nonroot`. Build context = repo root (so `ui/` is visible), but only project's `assets/static` is embedded.

## Initial scope (first run)

Components:
- **primitive**: Button, Tag, Card, Heading, Display, Icon
- **layout**: Container, Stack, Grid, Navbar, Footer, PageShell
- **overlay**: Modal, Toast, Portal
- **special**: ThemeToggle (cycles themes), ModeToggle (light/dark)

Projects:
- **homepage**: hero + skills stub with `PageShell` showing both toggles
- **showcase**: catalogue listing all components with live preview + theme/mode toggle
- **regex**: live RE2 regex tester — pattern/input/flags via Datastar signals, highlights via SSE `PatchElements`, replace output

Both projects buildable to static binary; showcase compiled with `-tags showcase`.

## Quality gates

- `golangci-lint` (govet, staticcheck, errcheck, gosec, bodyclose, errorlint, …)
- `govulncheck`
- `make guard-classes` — grep enforces no `Class(` / `class=` in `ui/**.go`
- `gomponents` render tests + Echo `httptest` handler tests
- Playwright e2e (theme toggle, mode toggle, modal, toast)
- lefthook pre-commit (fmt, guard-classes, lint) + pre-push (test, vuln)

## Quickstart

```bash
make setup        # tailwind binary + datastar.js + altcha.js + go mod tidy
make icons        # regenerate ui/icon/icons_gen.go
make dev PROJECT=homepage           # → :8090
make dev-showcase                   # → :8091
make dev-regex                      # → :8092
make check        # full local gate
make docker PROJECT=homepage TAG=v1
```

## Homepage Contact Mail

The contact form verifies the honeypot and ALTCHA first, then sends email via
SMTP. In development, missing SMTP config falls back to logging accepted
messages. In production (`MLJR_ENV=prod`), missing mail config fails startup.

Required production env:

```bash
ALTCHA_HMAC_KEY="$(openssl rand -hex 32)"
SMTP_HOST="smtp.example.com"
SMTP_PORT="587"
SMTP_USER="smtp-user"       # optional if the relay allows unauthenticated local send
SMTP_PASS="smtp-password"   # optional if the relay allows unauthenticated local send
SMTP_FROM="Portfolio <portfolio@example.com>"
CONTACT_TO="Michael Reinegger <reinemic2.0@gmail.com>"
```

The sender is always `SMTP_FROM`; the visitor email is set as `Reply-To` to
avoid SPF/DMARC failures.

## Homepage Homelab Live Panel

A background poller (60s) aggregates Uptime Kuma's public status-page JSON and
PromQL stats (CrowdSec, hosts) into an in-memory snapshot; the section
re-fetches `/api/homelab` via `data-on-interval` and patches by id. Visitors
never hit upstream sources.

```bash
HOMELAB_KUMA_URL="https://uptime.mljr.eu"   # default; public, no auth
HOMELAB_KUMA_SLUG="all"                     # default
HOMELAB_PROM_URL="http://nuc.tail33930.ts.net:19090"  # Tailscale-only; empty disables PromQL stats
```

The metrics backend is VictoriaMetrics (PromQL-compatible; deployed via the
homelab-automation repo). Dev without reachable sources renders
`homelab.Sample()`.
