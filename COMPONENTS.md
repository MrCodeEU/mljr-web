# mljr-web — Component Checklist

**RULE: This file must be kept up to date at all times.**
Every time a component is added, removed, or substantially changed, update this file.
Agents: treat any ✅ without a matching `*_showcase.go` as a bug. Treat any ☐ as a valid next target.

Legend: ✅ Done · 🚧 Partial · ☐ Not started

Tech stack: Go + gomponents · Datastar 1.0.2 · Tailwind v4 · Motion v10

---

## `ui/primitive/` — Atoms and display primitives

| Status | Name | Slug | Notes |
|--------|------|------|-------|
| ✅ | Avatar | `avatar` | image + initials fallback, sizes, shapes, status dot |
| ✅ | Avatar Group | `avatar-group` | stacked overlap, overflow count |
| ✅ | Badge | `badge` | status pills, dot variant, tones |
| ✅ | Button | `button` | primary/secondary/outline/ghost/danger/tone, sizes, disabled |
| ✅ | Callout | `callout` | info/success/warning/danger variants |
| ✅ | Card | `card` | themed surface container, interactive variant, tones |
| ✅ | Chip | `chip` | dismissable pill with tone |
| ✅ | Code Block | `code-block` | monospace, copy button |
| ✅ | Color Swatch | `color-swatch` | square/circle shape, label, palette group |
| ✅ | Copy Button | `copy-button` | clipboard + Datastar ✓ feedback |
| ✅ | Display | `display` | hero-sized headline text |
| ✅ | Feature Card | `feature-card` | icon box + title + description, tone |
| ✅ | Heading | `heading` | h1–h5 with size scale |
| ✅ | Icon | `icon` | Iconify SVG, 107 icons across lucide + simple-icons |
| ✅ | Kbd | `kbd` | keyboard shortcut display |
| ✅ | Pricing Card | `pricing-card` | price + feature list + CTA, highlighted tier |
| ✅ | Progress | `progress` | bar with variants |
| ✅ | Progress Ring | `progress-ring` | SVG circular progress |
| ✅ | Rating | `rating` | star rating, read-only mode |
| ✅ | Segmented Control | `segmented` | radio group as connected button bar |
| ✅ | Tag | `tag` | toned label |
| ✅ | Toggle Group | `toggle-group` | exclusive/multi-select button set |
| ✅ | Tooltip | `tooltip` | hover label, placement variants |
| ✅ | Announcement Bar | `announcement-bar` | Datastar signal, dismiss, CTA link, variants |
| ✅ | Button Group | `button-group` | attached mode collapses shared borders |
| ✅ | Collapse | `collapse` | max-height CSS transition, Datastar signal |
| ✅ | Countdown | `countdown` | JS setInterval, dd:hh:mm:ss + compact mode |
| ✅ | CTA Banner | `cta-banner` | full-width accent strip, dual CTA |
| ✅ | FAB | `fab` | fixed corner button, sizes, primary variant |
| ✅ | Icon Button | `icon-button` | aria-label wrapper over Button size=icon |
| ✅ | Speed Dial | `fab` | expandable mini-action list, Datastar signal, shares fab.go |
| ✅ | Share Button | `share-button` | Web Share API + clipboard fallback, all client-side |
| ✅ | Number Ticker | `number-ticker` | requestAnimationFrame counter, ease-out cubic, IntersectionObserver trigger |
| ✅ | Scroll Area | `scroll-area` | thin themed scrollbars, vertical/horizontal/both |
| ✅ | Split Button | `split-button` | primary action + chevron dropdown, Datastar signal |
| ✅ | Marquee | `marquee` | CSS-only infinite horizontal scroll, pause-on-hover, direction |
| ✅ | Word Rotate | `word-rotate` | fade+slide cycling words, setInterval, configurable interval |
| ✅ | Typewriter | `typewriter` | type-and-delete loop, configurable speed/pause, blinking cursor |
| ✅ | Media Card | `media-card` | image top + content body, badge overlay, lazy loading, hover zoom |
| ✅ | Flip Card | `flip-card` | CSS 3D rotateY flip, hover or Datastar-signal click trigger |
| ✅ | Gradient Text | `gradient-text` | CSS background-clip:text, any tag, from/via/to, theme-aware |
| ✅ | Scroll To Top | `scroll-to-top` | fixed button, IntersectionObserver threshold, smooth scroll |
| ✅ | Read Progress | `read-progress` | fixed top bar, fills on scroll, window or custom container target |
| ✅ | Video Player | `video-player` | native `<video>` + custom controls (play, scrub, volume, speed, fullscreen); zero external library |
| ✅ | Audio Player | `audio-player` | native `<audio>` + Canvas waveform animation (40 bars); Web Audio API analyser |
| ✅ | Tilt Card | `tilt-card` | pointer-driven 3D perspective tilt + optional shine gradient; zero JS libs |

---

## `ui/layout/` — Structural containers and navigation

| Status | Name | Slug | Notes |
|--------|------|------|-------|
| ✅ | Accordion | `accordion` | collapsible sections, `<details>` based |
| ✅ | Banner | `banner` | dismissable info/warning/success/danger strip |
| ✅ | Breadcrumb | `breadcrumb` | hierarchical path, separator |
| ✅ | Container | `container` | max-width centering wrapper |
| ✅ | Divider | `divider` | horizontal/vertical rule, optional label |
| ✅ | Footer | `footer` | page footer |
| ✅ | Grid / Col | `grid` | 12-column grid system |
| ✅ | Navbar | `navbar` | top bar with brand + nav + actions slots |
| ✅ | Sidebar | `sidebar` | collapsible left nav, Datastar signal driven |
| ✅ | Stack | `stack` | vertical/horizontal flex stack with gap |
| ✅ | Stepper | `stepper` | multi-step wizard progress |
| ✅ | Tabs | `tabs` | horizontal tab switcher, Datastar signal |
| ✅ | App Shell | `app-shell` | sidebar + main flex layout |
| ✅ | Auth Layout | `auth-layout` | centered card for login/register |
| ✅ | Background | `background` | dots/grid/lines/diagonal/cross/gradient CSS patterns |
| ✅ | Resizable Panels | `resizable-panels` | pointer+touch drag handle, horizontal/vertical |
| ✅ | Bento Grid | `bento-grid` | mosaic CSS grid, BentoItem with col/row span, landing pages/dashboards |
| ✅ | Masonry Grid | `masonry` | CSS `column-count` masonry layout; zero JS; `MasonryItem` prevents column breaks |
| ✅ | Table of Contents | `table-of-contents` | IntersectionObserver scroll-spy; auto-detects h2/h3/h4 or accepts explicit items; sticky option |

---

## `ui/form/` — User input

| Status | Name | Slug | Notes |
|--------|------|------|-------|
| ✅ | Checkbox | `checkbox` | styled with label |
| ✅ | Combobox | `combobox` | filterable select, keyboard nav |
| ✅ | Contact Form | `contact-form` | altcha + honeypot, full validation |
| ✅ | File Drop Zone | `file-drop-zone` | drag-and-drop upload, file list display |
| ✅ | File Input | `file-input` | styled native file input |
| ✅ | Input | `input` | text/email/number with field wrapper |
| ✅ | Input Group | `input-group` | prefix/suffix addons |
| ✅ | Number Input | `number-input` | stepper with ±buttons |
| ✅ | OTP Input | `otp-input` | single-digit boxes, auto-advance, paste |
| ✅ | Password Input | `password-input` | show/hide toggle, Datastar type swap |
| ✅ | Radio Group | `radio` | styled radio buttons |
| ✅ | Rating | `rating` | star input (also in primitive for display) |
| ✅ | Select | `select` | native select, styled wrapper |
| ✅ | Slider | `slider` | range with floating value label |
| ✅ | Switch | `switch` | toggle/checkbox hybrid |
| ✅ | Tags Input | `tags-input` | multi-value Enter/comma add, × remove |
| ✅ | Textarea | `textarea` | auto-resize, field wrapper |
| ✅ | Color Input | `color-input` | native color picker, hex label, Datastar signal |
| ✅ | Date Input | `date-input` | native date picker, styled |
| ✅ | Range Pair | `range-pair` | dual-handle min/max slider, Datastar enforces low≤high |
| ✅ | Search Input | `search-input` | debounced @get, clear button, Datastar |
| ✅ | Time Input | `time-input` | native time picker, step support |
| ✅ | Multi Select | `multi-select` | chip-based multi-value, dropdown, max selection, hidden inputs |
| ✅ | Calendar Picker | `calendar-picker` | custom month-grid date picker; ISO 8601 hidden input; min/max restriction; zero external library |

---

## `ui/data/` — Structured content display

| Status | Name | Slug | Notes |
|--------|------|------|-------|
| ✅ | Bar Chart | `bar-chart` | pure SVG, values, grid, custom colors |
| ✅ | Carousel | `carousel` | touch swipeable image gallery |
| ✅ | Description List | `description-list` | key/value pairs |
| ✅ | Donut Chart | `donut-chart` | stroke-dasharray SVG, legend, center label |
| ✅ | Line Chart | `line-chart` | multi-series, bezier, area fill, grid |
| ✅ | List | `list` | ordered/divided variants |
| ✅ | Pagination | `pagination` | page controls |
| ✅ | Sparkline | `sparkline` | tiny inline trend line |
| ✅ | Stat Card | `stat-card` | number + label + delta trend |
| ✅ | Table | `table` | responsive data table, striped |
| ✅ | Timeline | `timeline` | vertical event timeline, tones |
| ✅ | Data Grid | `data-grid` | sort, filter, pagination — all client-side JS |
| ✅ | Lazy Image | `lazy-image` | IntersectionObserver, skeleton fade-in |
| ✅ | Lightbox | `lightbox` | thumbnail grid, fullscreen viewer, ←/→ keyboard nav |
| ✅ | Pie Chart | `pie-chart` | SVG arc paths, legend, percentage labels |
| ✅ | Tree View | `tree-view` | recursive details elements, no JS |
| ✅ | Virtual List | `virtual-list` | CSS content-visibility:auto viewport culling, zero JS |
| ✅ | Sortable List | `sortable` | HTML5 drag-and-drop reorder, optional grab handle, onChange callback |
| ✅ | Snake Timeline | `snake-timeline` | serpentine layout, alternating LTR/RTL rows, curved connectors, 2–4 cols |
| ✅ | GitHub Heatmap | `heatmap` | contribution grid SVG, server-side, color-mix intensity, month/day labels |
| ✅ | Radar Chart | `radar-chart` | multi-series spider chart SVG, concentric grid, server-side Go math |
| ✅ | Gauge / Meter | `gauge` | 270° SVG arc, tick marks, center value+label, server-side |
| ✅ | Syntax Highlighter | `syntax-highlighter` | server-side chroma v2, 300+ languages, 40+ themes, optional line numbers |
| ✅ | Horizontal Timeline | `horizontal-timeline` | CSS scroll-snap strip, dot+card per item, zero JS |
| ✅ | Infinite Scroll | `infinite-scroll` | sentinel `data-on:intersect__once` triggers Datastar `@get`; server appends HTML |

---

## `ui/overlay/` — Layered UI

| Status | Name | Slug | Notes |
|--------|------|------|-------|
| ✅ | Alert Dialog | `alert-dialog` | confirmation with cancel/confirm, Datastar |
| ✅ | Command Palette | `command` | ⌘K search modal, grouped items, filter |
| ✅ | Drawer | `drawer` | slide-in panel left/right, sizes |
| ✅ | Dropdown | `dropdown` | anchored floating menu |
| ✅ | Modal | `modal` | accessible dialog, backdrop, sizes |
| ✅ | Popover | `popover` | rich hover/click card, placement |
| ✅ | Toast | `toast` | ephemeral notification, auto-dismiss, queue |
| ✅ | Context Menu | `context-menu` | right-click, cursor-anchored, data-ctx trigger |
| ✅ | Notification Stack | `notification-stack` | fixed-position feed, window._pushNotification(), auto-dismiss |
| ✅ | Sheet | `sheet` | full-edge slide-in (bottom/right/left), Datastar signal |
| ✅ | Hover Card | `hover-card` | CSS :hover rich info card, 4 placements, no JS |

---

## `ui/feedback/` — Status communication

| Status | Name | Slug | Notes |
|--------|------|------|-------|
| ✅ | Alert | `alert` | info/success/warning/danger, dismiss |
| ✅ | Empty State | `empty-state` | icon + title + action for zero results |
| ✅ | Skeleton | `skeleton` | loading placeholder, text/block/circle |
| ✅ | Spinner | `spinner` | loading spinner, 5 variants, sizes |
| ✅ | Loading Overlay | `loading-overlay` | full-surface spinner, Datastar signal |
| ✅ | Notification Badge | `notification-badge` | count badge, dot mode, max cap |
| ✅ | Shimmer | `shimmer` | animated gradient loading placeholder, lines/circle variants |

---

## `ui/special/` — Cross-cutting concerns

| Status | Name | Slug | Notes |
|--------|------|------|-------|
| ✅ | Captcha (Altcha) | `captcha` | altcha widget wrapper, honeypot |
| ✅ | Theme / Mode Toggle | `theme-toggle` | Datastar `$theme`/`$mode` signals |
| ✅ | Logo Scatter | — | Motion v10 scatter animation, shared component in `ui/special/logo_scatter.go` |
| ✅ | Language Toggle | `language-toggle` | cycle locales, cookie persist, optional reload |
| ✅ | User Menu | `user-menu` | composite Avatar + Datastar dropdown with identity header + actions |
| ✅ | Login Form | `login-form` | composite: email + password (show/hide) + remember me + submit |
| ✅ | Cookie Banner | `cookie-banner` | GDPR consent banner, localStorage persist, position top/bottom |
| ✅ | OpenMap | `open-map` | Leaflet 1.9 self-hosted (148KB), OSM tiles, drop pins with popups |
| ✅ | Confetti | `confetti` | Canvas particle burst, rAF animated, trigger: click/load/manual |
| ✅ | Onboarding Tour | `onboarding-tour` | clip-path spotlight overlay, floating tooltip, step navigation, zero library |

---

## `ui/datastar/` — Datastar feature showcases

| Status | Name | Slug |
|--------|------|------|
| ✅ | Signals | `ds-signals` |
| ✅ | Reactivity | `ds-reactivity` |
| ✅ | Events | `ds-events` |
| ✅ | Effects & Refs | `ds-effects` |
| ✅ | Fetch (SSE) | `ds-fetch` |
| ✅ | Server Push (SSE) | `ds-push` |
| ✅ | Patterns | `ds-patterns` |
| ✅ | Modifiers | `ds-modifiers` |
| ✅ | Animation (Motion) | `ds-animation` |

---

## Animation showcases (`ui/datastar/animation_*`)

| Status | Name | Slug | Technique |
|--------|------|------|-----------|
| ✅ | Logo Scatter | `ds-anim-logo` | SVG user-unit transform, Motion.animate, sway loop |
| ✅ | Spring Physics | `ds-anim-spring` | `Motion.spring({stiffness,damping,mass})` |
| ✅ | inView Reveals | `ds-anim-inview` | `Motion.inView` with custom scroll root |
| ✅ | Text Reveal | `ds-anim-text` | char split + stagger, 3 distinct effects |
| ✅ | Gesture & Hover | `ds-anim-gesture` | hover lift, magnetic button, click ripple |
| ✅ | Scroll Progress | `ds-anim-scroll` | Motion.scroll + inView reveals inside scrollable container |
| ✅ | Loading Morphs | `ds-anim-morph` | Motion.timeline skeleton→content, stagger entrance |
| ✅ | Page Transitions | `ds-anim-transition` | Motion.timeline directional slide between simulated views |

---

## Showcase infrastructure

| Status | Feature |
|--------|---------|
| ✅ | Auto-registration via `init()` + build tag `showcase` |
| ✅ | Catalogue page: search + category filter + iframes |
| ✅ | Iframe theme/mode sync via Datastar `$theme`/`$mode` signals |
| ✅ | Stack intro section (philosophy + tech cards + request flow) |
| ✅ | Detail page: controls panel, viewport tabs, usage code, props table |
| ✅ | Detail page: responsive (single-col on ≤768px) |
| ✅ | Preview iframe: Motion v10 + Datastar loaded |
| ✅ | Icon showcase: searchable grid, grouped by set, click to copy |
| ✅ | 160 showcase entries registered (recount: `grep -rc "registry.Register" --include="*.go" ui/ projects/`) |
| ✅ | Prev/next keyboard navigation (←/→) via ArrowLeft/ArrowRight in detail.go |
| ✅ | Patterns page: `/patterns` listing with iframe previews + theme sync |
| ✅ | Pattern detail: `/patterns/{slug}` with controls + full iframe |
| ✅ | Pattern: Login Page (auth-login) |
| ✅ | Pattern: Register Page (auth-register) |
| ✅ | Pattern: App Dashboard (app-dashboard) |
| ✅ | Pattern: Pricing Page (marketing-pricing) |
| ✅ | Pattern: Settings Page (app-settings) |
| ☐ | Markdown export report |

---

## Homepage project (`projects/homepage`)

| Status | Feature | Notes |
|--------|---------|-------|
| ✅ | Umami analytics | Self-hosted, proxied via `/umami/*`. `AnalyticsConfig{UmamiScriptSrc,UmamiWebsiteID,UmamiHostURL,UmamiDomains}`. Both fields required or `AnalyticsHead()` returns nil. |
| ✅ | Analytics reverse proxy | `analytics_proxy.go` — proxies to upstream Umami; stub JS on failure |
| ✅ | Strava integration | `pages/strava.go` — YTD stats + recent activities; only renders when `d.HasStrava()` is true |
| ✅ | Legal pages | `/impressum` + `/datenschutz` (Austrian ECG / DSGVO requirements) in `pages/legal.go` |
| ✅ | Shared navbar/footer | `siteNavbar()` + `siteFooter()` extracted to `pages/legal.go` — used by home + legal pages |
| ✅ | Hero section | Typewriter tagline + GradientText headline + Bento Grid (photo, year ticker, project ticker, MSc, Dynatrace) |
| ✅ | Stats strip | 4× NumberTicker + Marquee of colorful tech-stack pills |
| ✅ | Snake timeline | Experience in serpentine layout (3 cols desktop) + vertical fallback on mobile |
| ✅ | Skills marquee | 3 rows, alternating direction, colored pills (8-tone cycle) |
| ✅ | Read progress bar | 8px accent bar at top via `primitive.ReadProgress` |
| ✅ | Logo scatter (dual) | Primary scroll-triggered (540px, left of center) + secondary loop (380px, lower-right) |
| ✅ | Page animations | Motion.inView stagger entrance per section; hero entrance stagger on load; project pagination fade-slide via MutationObserver |
| ✅ | Mobile responsive | `homepageCSS` const: ≤900px hero stack + new-section single-col, ≤768px bento collapse + snake→vertical timeline, ≤480px CTA full-width |
| ✅ | Numbered section headers | Swiss-editorial outlined index (`.section-num`), `sectionHeader(num, heading, sub, tone)` in `pages/skills.go`. Order: 01 Experience · 02 Featured · 03 Projects · 04 Open Source · 05 Homelab · 06 Activity · 07 Skills · 08 Under the hood · 09 Contact |
| ✅ | Featured work section | `pages/featured.go` — asymmetric spotlight grid, TiltCard + shine, first project double-height; main projects grid then shows non-featured only |
| ✅ | Open Source section | `pages/github.go` — contribution Heatmap + 3 language Gauges + 4 NumberTicker counters. Repos/stars real; heatmap/commits/streak are deterministic placeholder data (labeled) until a GitHub stats pipeline lands |
| ✅ | Homelab live panel | `homelab/` poller (60s) → Uptime Kuma public status-page API + PromQL (CrowdSec bans/attacks, hosts online). `data-on-interval` refresh via `/api/homelab` SSE fragment. Env: `HOMELAB_KUMA_URL`, `HOMELAB_PROM_URL`; dev fallback `homelab.Sample()` |
| ✅ | Under the hood section | `pages/codeshowcase.go` — chroma-highlighted `gaugeExcerpt` (keep in sync with `ui/data/gauge.go`) + fact chips |
| ✅ | No section backgrounds | Sections must not set `background:` — it blocks the logo-scatter page background |
