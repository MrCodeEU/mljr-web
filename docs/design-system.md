# mljr.eu design system — Swiss (International Typographic Style)

How the `swiss` theme looks, which building blocks exist, and how to build a
new page with them. Applies to the homepage, its sub pages (games, legal,
project details, 404) and the live tools (regex, cron, codec). The newsletter
app stays on `swissbrut`.

Source of truth: tokens in `ui/css/_base.css` (`[data-theme="swiss"]`),
component CSS in `ui/css/_swiss.css`, components in `ui/swiss/` and
`ui/special/toolshell.go`. The homepage's own section CSS lives in
`projects/homepage/pages/*.go` (`…CSS` consts).

## 1. Principles

- **Grid first.** Full screen width, 12 columns: a 3-column *anchor* (index
  numeral, German title, English sub, note) beside a 9-column *body*. Every
  page is a stack of these rows, separated by a 3px ink rule. Cells inside a
  body are separated by 1px ink rules (drawn as box-shadows).
- **Flat and hard-edged.** No radius, no shadows, no gradients. Depth comes
  from colour fields, not elevation.
- **Type does the work.** Inter for text and display (900 weight, negative
  tracking for big titles), JetBrains Mono for labels, numbers and code.
  Labels are UPPERCASE mono, small, letter-spaced.
- **Bilingual titles.** Section titles are German and UPPERCASE, with the
  English name as a small mono sub label (`Werkkatalog` / `Catalogue`). Body
  copy follows the page language.
- **Lettered kickers.** Cells open with a kicker: `A. Label` left, an aside
  right (count, unit, "● Live"). Letters run A, B, C… within a section.
- **Yellow means "you can act".** `--act` (#ffd400) is reserved for hover and
  focus of interactive things (buttons, tabs, links, pager squares, focus
  rings). Never use it as a decorative or data colour.
- **Accessibility is not optional.** Every text pair ≥ 4.5:1, graphics ≥ 3:1,
  in light *and* dark. Animations stop under `prefers-reduced-motion`.

## 2. Colour

| Token | Light "Klein-Blau" | Dark "Klein-Nacht" | Use |
|---|---|---|---|
| `--bg` / `--ink` | #f3f1ea / #0e0f14 | #0b0d16 / #eef0f6 | paper and text |
| `--block` (+`--block-fg`) | #002fa7 | #002fa7 | main field colour (hero block, anchors) |
| `--second` (+`-fg`) | #d96200 orange | #5a78ff | 2nd field / graphics |
| `--third` (+`-fg`) | #a40054 magenta | #e2357f | 3rd field / graphics |
| `--main-text` | #002fa7 | #ff7a1a | coloured text that must stay readable |
| `--muted` | | | secondary text (still ≥ 4.5:1) |
| `--rule-subtle` | | | hairlines inside lists |
| `--act` / `--act-fg` | #ffd400 / #0b0d16 | same | interactive only |
| `--board` | #0e0f14 | #05060a | surfaces that stay dark in both modes |

Rules:

- Text on a coloured field always uses the matching `-fg` token.
- A surface that carries yellow text or a yellow aside (result bars, tab
  leads, the departure board) uses `--board` with fixed `#f3f1ea` text:
  `--ink` turns light in dark mode and yellow on it fails contrast.
- Light and dark are different colourways on purpose (blue loses punch when
  lightened). Do not derive one from the other.

## 3. Layout building blocks (`ui/swiss`)

| Component | What it is |
|---|---|
| `swiss.Page(sheetID, bg, rows…)` | the sheet; `bg` is an optional background layer (logo grid) |
| `swiss.Section(SectionProps, cells…)` | one row: anchor + body. `Field` paints the anchor (`block`, `second`, `third`, `ink`); `Wide` turns the anchor into a header strip; `NoScatter` keeps logo bars out; `AnchorExtra` adds content to the anchor (Bauhaus mark, logo slot) |
| `swiss.Cell(CellProps, …)` | ruled box in a body. `Span` 1–9 (no span = full row); `Variant` `block`/`second`/`third`/`ink`/`elevated`/`flush` |
| `swiss.Kicker(label, aside)` | the mono line that opens a cell |
| `swiss.Figure`, `swiss.Micro`, `swiss.Pill(s)` | big number, small mono label, spec tag |
| `swiss.Btn` / `swiss.BtnLink` | square mono button; `primary` = filled; `aria-pressed="true"` = ink fill |
| `swiss.Colophon` + `swiss.Register` | sticky top bar with the section ruler (labels below 1800px, numbers above) and the full-screen register on mobile; `SkipLabel` adds the skip link to `#main`; `Primary` holds links that must always be reachable (CV, contact): in the bar on desktop, at the top of the register on mobile |
| `swiss.Gallery`, `swiss.FlapBoard`, `swiss.Fahrplan` | figure gallery, departure board (last column `Value`, up to 5 chars), timeline |
| `swiss.SROnly` | text for screen readers only |

Tool and form blocks (`ui/swiss/tool.go`):

| Component | What it is |
|---|---|
| `swiss.ToolHead(name, lede, pills…)` | opening cell of a tool: huge uppercase name, lede, pills |
| `swiss.Field(id, label, control…)` | labelled control; inputs, textareas and selects inside get the Swiss form look (1px ink border, mono, yellow focus ring) |
| `swiss.Affix(text)` | text glued to an input (regex slashes) |
| `swiss.Choice(label, buttons…)` | segmented group of `swiss.Btn` with `aria-pressed` |
| `swiss.ResultBar(label, aside)` | dark strip heading a result; aside in yellow |
| `swiss.ResultBlock(title, …)` | titled part of a result |
| `swiss.Row(num, main, aside, detail…)` | numbered result line |
| `swiss.RefList(title, rows)` | cheat-sheet `<dl>` |
| `swiss.ErrorBox(msg)` | `role="alert"` error |
| `swiss.Footer(note, links, bottom…)` | slim poster footer for standalone pages |

Kit rule: no `Class(` in `ui/**.go` (`make guard-classes`). Components emit
`data-component` / `data-slot`; CSS selects on those. Pages under `projects/`
may use classes for page-specific layout.

## 4. Decoration

- **Bauhaus marks** (`projects/homepage/pages/bauhaus.go`): a 4×4 grid
  composition of circles, half/quarter circles, triangles, rings and stripes
  in each section's anchor, one composition per section. A colour role that
  equals the section's field colour is swapped. The sticky anchor head has an
  opaque background so the title always reads over the mark. Hidden on wide
  sections and below 1024px.
- **Icon shapes**: figure rows put their icon in a square, circle or triangle
  (`.fig-shape`), cycling through the three field colours.
- **Footer poster**: wordmark block, circle, triangle.
- **Logo grid** (`special.LogoGrid`): the logo assembles in the hero anchor
  and scatters onto the column lines on scroll. Static under reduced motion.
  Sections with dense UI set `NoScatter`.

## 5. Interaction patterns

- **Tabs + pager** (`projects/homepage/pages/tabs.go`, `tabsJS`): a
  `[data-tabs]` root with `[data-tab="n"]` buttons (`aria-pressed`) and
  `[data-tab-page="n"]` pages; a pager under each page names the previous and
  next item. Arrow keys move between tabs. A tab lead cell ("Pick a plate",
  board colours) says what the buttons do.
- **Accordion**: native `<details name="…">` so only one row is open; the
  first row starts open to show the pattern.
- **Lazy heavy parts**: iframes get `data-src` and load when visible (lab
  games start when the lab nears the viewport); Leaflet loads on
  IntersectionObserver; the captcha script loads near the form and starts its
  proof-of-work on focus (`special.CaptchaProps{Lazy: true}`).
- **Live tools**: one signal scope around the workbench section,
  `input__debounce.200ms` posts to `/api/eval`, the server patches the result
  fragment by id.

## 6. Accessibility checklist

- One `<main id="main">` per page; skip link via `SkipLabel`.
- Every section is labelled (`aria-labelledby` on the title).
- Buttons that toggle use `aria-pressed` with explicit strings: in Datastar
  write `data-attr:aria-pressed="$x ? 'true' : 'false'"` — a bare boolean
  renders `aria-pressed=""` or removes it, which is invalid and unstyled.
- Never name a page signal `$mode` or `$theme`: `ThemeToggleRoot` owns them
  and mirrors `$mode` to `<html data-mode>` (a tool once broke dark mode by
  using `$mode` for its own state).
- Focus: yellow `--act` outline or fill; SVG controls show focus with a fill
  (outline does not render on `<g>`).
- Decorative shapes are `aria-hidden`; charts are `role="img"` with a label.
- Scrollable code blocks are focusable (`tabindex="0"`) and labelled.
- Check with axe in light, dark and at 390px before shipping
  (0 violations is the bar).

## 7. Recipes

### A homepage section

```go
// home.go: add to `all`, give it a field colour in sectionField if wanted
{"papers", "Schriften", "Papers", func(p swiss.SectionProps) g.Node { return papersSection(p, d, lang) }},

func papersSection(p swiss.SectionProps, d hpdata.SiteData, lang string) g.Node {
	p.Note = tr(lang, "…", "…")
	return swiss.Section(p,
		swiss.Cell(swiss.CellProps{Span: 3},
			swiss.Kicker("A. Label", g.Text("aside")),
			swiss.Figure(g.Text("42")),
		),
		swiss.Cell(swiss.CellProps{Span: 6}, /* … */),
	)
}
```

Numbering, ruler and Bauhaus mark follow automatically from the order in
`all`.

### A sub page of the homepage (game, legal, detail)

```go
h.Main(h.ID("main"), h.Class("sub-page"),
	subHead("/#playground", tr(lang, "Labor", "Lab"), mark,
		h.H1(g.Text("Boids")),
		h.P(h.Class("sub-lede"), g.Text(lede)),
	),
	h.Div(h.Class("sub-body"), /* content */),
),
subRow("", tr(lang, "So funktioniert's", "How it works"), "", /* body */),
```

`subHead`, `subRow` and their CSS live in
`projects/homepage/pages/subpage.go`.

### A standalone tool

```go
special.ToolShell(special.ToolShellProps{
	ID: "codec", Name: "Codec", Title: "…", Description: "…", CSS: pageCSS,
	Sections: []swiss.NavItem{{ID: "tool", Num: "00", Label: "Werkzeug", Sub: "Tool"}, …},
},
	swiss.Section(swiss.SectionProps{ID: "tool", Num: "00", Title: "Werkzeug", Sub: "Tool", Field: "block", Note: "…"},
		swiss.Cell(swiss.CellProps{}, swiss.ToolHead("Codec", "lede", "Hash", "Go")),
	),
	h.Div(ui.Signals(`{…}`), ui.On("input__debounce.200ms", "@post('/api/eval')"),
		swiss.Section(swiss.SectionProps{ID: "bench", Num: "01", Title: "Werkbank", Sub: "Workbench"},
			swiss.Cell(swiss.CellProps{Span: 4}, swiss.Field("in", "A. Input", h.Textarea(h.ID("in"), ui.Bind("input")))),
			swiss.Cell(swiss.CellProps{Span: 5, Variant: "flush"}, OutputFragment(initial)),
		),
	),
)
```

Section rhythm for tools: `00 Werkzeug` (blue field, head) → `01 Werkbank`
(inputs left, result right) → `02 Referenz` (cheat sheet, optional).

## 8. Before shipping a UI change

1. `make guard-classes`, `go vet ./...`, `go test ./...`.
2. Rebuild Tailwind for every project that imports `ui/css`.
3. Screenshot desktop (1440) and mobile (390), light and dark; no
   horizontal overflow.
4. axe: 0 violations; tab through the page once.
5. If `ui/css` changed, bump the `app.css?v=` string in
   `ui/layout/pageshell.go`.
