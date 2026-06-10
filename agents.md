# agents.md — implementation guide for mljr-web v2

This file tells future Claude/agent runs **how to extend this repo without breaking its invariants.** Read `PLAN.md` first for the high-level architecture.

## Hard invariants (CI-enforced)

1. **`ui/**.go` MUST NOT contain `Class(` or `class=`.** Components emit `data-component`, `data-variant`, `data-size`, `data-tone`, `data-state` — all visual styling lives in `ui/css/core.css`. Enforced by `make guard-classes`.
2. **No runtime CDN, no telemetry.** Every JS/CSS/font asset is vendored under `projects/<P>/assets/static/` and served from `/static/`.
3. **gomponents v1 has no `svg` package.** Use `g.Raw(svgString)` only. Generator at `tools/icongen` is the single source of SVG strings.
4. **One Go module.** Do not add nested modules unless `tools/` grows heavy deps — then split via `go.work`.
5. **No `class=` even via `g.Attr("class", ...)` in `ui/`.** The grep guard catches both.
6. **Datastar SDK is `github.com/starfederation/datastar-go/datastar`** (v1.0.x). API: `NewSSE`, `PatchElements`, `PatchSignals`, `RemoveElements`, `ExecuteScript`, `Redirect`, `ReadSignals`, `MarshalAndPatchSignals`.

## Module path

`mljr-web`. Imports: `mljr-web/ui/primitive`, `mljr-web/ui/layout`, `mljr-web/internal/web`, etc.

## Conventions

### Component file shape

```go
package primitive

import (
    g "maragu.dev/gomponents"
    h "maragu.dev/gomponents/html"
    "mljr-web/ui/token"
)

type XProps struct {
    Variant token.Variant
    Size    token.Size
    Tone    token.Tone
    Attrs   []g.Node          // pass-through: id, aria-*, data-on-*, Datastar bindings
}

func X(p XProps, children ...g.Node) g.Node {
    if p.Variant == "" { p.Variant = token.Primary }
    if p.Size == "" { p.Size = token.SizeMD }
    return h.Tag(
        g.Attr("data-component", "x"),
        g.Attr("data-variant", string(p.Variant)),
        g.Attr("data-size", string(p.Size)),
        g.If(p.Tone != "", g.Attr("data-tone", string(p.Tone))),
        g.Group(p.Attrs),
        g.Group(children),
    )
}
```

- Import aliases fixed: `g` = `maragu.dev/gomponents`, `h` = `maragu.dev/gomponents/html`. No dot-imports in `ui/`.
- One component per file. Sibling `*_showcase.go` with `//go:build showcase` for registry metadata.
- All enums via `mljr-web/ui/token` for compile-time safety.
- Default values applied at top of function. Zero values must be safe.
- Pass-through attrs via `Attrs []g.Node` field + `g.Group(p.Attrs)`. Always last before children.

### CSS rules

- Add new rules to the appropriate `ui/css/_*.css` partial under `@layer components` keyed off `[data-component="X"]` + `[data-variant=…]`.
- The partials are: `_base.css`, `_layout.css`, `_primitive.css`, `_overlay.css`, `_form.css`, `_data.css`, `_feedback.css`.
- `core.css` imports all partials. `projects/*/assets/css/input.css` imports them directly (nested imports may not propagate in Tailwind v4).
- Tokens (palette, semantic, scale) live in `_base.css`. Themes are role-assignment blocks there.
- Theme blocks: `[data-theme="swissbrut"][data-mode="light"] { … }` etc. Four blocks total.
- Use `var(--accent)`, `var(--surface)`, `var(--ink)`, `var(--line)`, `var(--shadow)` — never raw palette vars in components.

**CRITICAL: After editing any `ui/css/_*.css` partial, rebuild the compiled CSS:**
```bash
bin/tailwindcss -i projects/showcase/assets/css/input.css -o projects/showcase/assets/static/app.css
bin/tailwindcss -i projects/homepage/assets/css/input.css  -o projects/homepage/assets/static/app.css
```
In dev mode (`MLJR_ENV != "prod"`), the server reads CSS from disk — no server restart needed after rebuild. But the browser needs a hard reload (Ctrl+Shift+R) to bypass cache.

**CSS custom properties for variants:** Components that vary by variant should NOT rely solely on `--tone`/`--on-tone` for backgrounds (those may equal `var(--surface)` = page bg = invisible). Add explicit fallback:
```css
[data-component="toast"] {
    background: var(--surface-2);  /* visible even without variant */
    color: var(--ink);
}
[data-component="toast"][data-variant="success"] { background: var(--success); color: #fff; }
```

**Multiple CSS animations on one element:** When two animations both affect `transform`/`opacity`, the **last listed wins**. Design exit + entrance animations to avoid conflicts, or use separate `animation-name` on pseudo-elements.

**Icon colors in alerts/toasts:** Icon SVGs use `fill="currentColor"` / `stroke="currentColor"`. They inherit parent `color`. Explicitly set `color` per variant on the icon slot:
```css
[data-component="alert"][data-variant="info"] [data-slot="icon"] { color: var(--info); }
```

### Datastar helpers

Use `ui/datastar.go` wrappers:
- `ui.On(event, expr)` → `data-on-EVENT`
- `ui.Bind(signal)` → `data-bind-SIGNAL`
- `ui.Text(expr)` → `data-text`
- `ui.Show(expr)` → `data-show`
- `ui.Signals(json)` → `data-signals`

Add new wrappers there if Datastar attribute spelling shifts.

### SSE fragments

Fragments for `PatchElements` are **the same component called directly** — render to a `bytes.Buffer` via `internal/web.RenderToString(node)`. Always include a stable `id` on the root element so Datastar morphs by id.

### Registry (showcase metadata)

```go
//go:build showcase

package primitive

import "mljr-web/ui/registry"

func init() {
    registry.Register(&registry.Component{
        Slug: "button", Name: "Button", Category: "primitive",
        Summary: "…",
        Controls: []registry.Control{ … },
        Render: func(p map[string]string) g.Node { … },
    })
}
```

Prod projects build **without** `-tags showcase` → registry is tree-shaken.

## Theme model (v2)

- Themes: `swissbrut`, `ink`. Modes: `light`, `dark`. Total 4 skins.
- Root attrs: `<html data-theme="swissbrut" data-mode="light">`.
- CSS keys off both: `[data-theme="swissbrut"][data-mode="dark"] { … }`.
- `ThemeToggle` cycles `$theme` signal across `["swissbrut","ink"]` and patches `data-theme` on `<html>`.
- `ModeToggle` flips `$mode` across `["light","dark"]` and patches `data-mode` on `<html>`.
- Persist to `localStorage` via small inline `<script>` in `<head>` to avoid FOUC.
- `prefers-color-scheme` is only the initial fallback if no stored value.

## Adding an icon

1. Append `set:name` to `tools/icongen/icons.txt` (e.g. `lucide:github`).
2. `make icons` → regenerates `ui/icon/icons_gen.go` (deterministic; commit it).
3. Use `icon.Icon("lucide:github")` anywhere.

Supported sets (cached): `lucide`, `simple-icons`, `mdi`. Add more by adding entries; generator fetches new sets on demand.

## Adding a project

1. `mkdir -p projects/<name>/{pages,assets/css,assets/static/fonts}`
2. Copy `projects/homepage/main.go` as template. Adjust `Page.Theme` / port.
3. Create `projects/<name>/assets/css/input.css` importing `ui/css/core.css`.
4. Copy `projects/homepage/Dockerfile`, substitute name.
5. `make dev PROJECT=<name>` works automatically.

## Adding a component

1. `ui/<category>/<name>.go` — component + `Props` struct.
2. `ui/<category>/<name>_showcase.go` — registry entry (build tag `showcase`).
3. Add CSS rules to `ui/css/core.css` under `@layer components`.
4. Add render test under `ui/<category>/<name>_test.go` — assert structural `data-*` attrs present, no `class=`.

## Datastar quick reference

| Need | Attribute |
|---|---|
| Local signal | `data-signals="{open:false}"` |
| Bind input | `data-bind:name` (**colon**, not hyphen) |
| Text reactive | `data-text="$name"` |
| Show/hide | `data-show="$open"` |
| Reactive attr | `data-attr="{'data-theme':$theme}"` |
| Click | `data-on:click="$open=!$open"` (**colon** — Datastar parses `plugin:key__mod`) |
| Debounced | `data-on:input__debounce.300ms="@get('/api/x')"` |
| Window event | `data-on:keydown__window="…"` |
| Polling | `data-on:interval__5s="@get('/api/health')"` |
| Loading indicator | `data-indicator:fetching` (**colon**) |
| Backend call | `@get('/x')`, `@post('/x')`, etc. |
| Skip from server | prefix signal `_` (e.g. `_local`) |

**Critical Datastar v1.x syntax rule:** plugins with a key use **colon** as separator (`data-on:click`, `data-bind:name`, `data-indicator:loading`). Plugins without a key use no separator (`data-signals`, `data-show`, `data-text`, `data-effect`). Using hyphens (`data-on-click`) silently fails — no error, no handler registered.

## Datastar 1.0.2 — verified available plugins

Confirmed by inspecting `projects/showcase/assets/static/datastar.js` (Datastar v1.0.2).
Extract plugin names: `grep -o 'name:"[a-z-]*"' datastar.js | sort -u`

### Client-side attribute plugins

| Attribute | Requirement | Notes |
|---|---|---|
| `data-signals` | value optional, key optional | Declare/merge signal state |
| `data-computed:name` | value required | Derived reactive signal. `data-computed:total="$a+$b"` creates `$total` |
| `data-text` | value required, no key | Sets `el.textContent` reactively |
| `data-show` | value required, no key | Toggles `display:none`. Uses CSS, no JS DOM manipulation |
| `data-bind:signal` | key XOR value ("exclusive") | Two-way binding. For range/number add `data-on:input` too |
| `data-attr` | value required | Sets HTML attributes from JS object expression |
| `data-style` | value required | Sets inline CSS properties (merges, doesn't replace) |
| `data-class` | value required | Adds/removes individual CSS classes |
| `data-ref:name` | key XOR value | Captures DOM element as signal `$name` |
| `data-on:EVENT` | key required, value required | DOM event handler |
| `data-on-interval` | value required, no key | `setInterval`. Duration: `data-on-interval__duration.1s` |
| `data-on-intersect` | value required, no key | IntersectionObserver. Mods: `__full`, `__half`, `__threshold.50`, `__once`, `__exit` |
| `data-on-signal-patch` | value required | Fires when `PatchSignals` SSE arrives. `argNames:["patch"]` |
| `data-effect` | value required, no key | Reactive side-effect. Re-runs when any accessed signal changes |
| `data-init` | value required, no key | One-shot on mount. Does NOT re-run |
| `data-indicator:name` | key XOR value | Sets `$name=true` while SSE request is in-flight |

### SSE action plugins (server → client)

Triggered by `datastar-patch-elements` / `datastar-patch-signals` SSE events:
`outer`, `inner`, `append`, `prepend`, `before`, `after`, `remove`, `replace`

### **`data-for` does NOT exist in Datastar 1.0.2**

There is no loop/iteration attribute. Do NOT use `data-for` — it will be silently ignored.

**Alternative for dynamic lists (e.g. toast queue, comment feed):**

Option A — **DOM manipulation via `<script>` helper** (recommended for showcase/demos):
```go
// In Render(): embed a <script> that defines a JS function
h.Script(g.Raw(`
    window._addItem = function(data) {
        var container = document.querySelector('[data-component="list"]');
        var el = document.createElement('div');
        el.setAttribute('data-component', 'item');
        el.innerHTML = data.title;
        container.appendChild(el);
        if (container.children.length > 6) container.firstChild.remove();
    };
`)),
// Button calls the function
g.Attr("data-on:click", "_addItem({title:'Hello'})"),
```

Option B — **SSE `PatchElements`** (recommended for production):
Server patches individual elements into a mounted container. Each toast/item has a stable `id`.
```go
sse.PatchElements(web.RenderToString(
    overlay.Toast(overlay.ToastProps{ID: "t-"+id, ...}, ...),
))
```

Option C — **Pre-rendered fixed slots** (for ≤ N items, e.g. wizard steps):
Render N slots server-side, each with `data-show="$slot_N.active"`.

### `data-attr` gotchas

`data-attr` sets HTML *attributes* (not CSS properties) via JavaScript object expression:
```
data-attr='{"data-variant": $theme}'    ✓ sets attribute data-variant
data-attr='{"style": "color:red"}'      ✓ replaces entire style attribute
```

**Setting CSS custom property via style:**
```
data-attr='{"style":"--val:"+$signal}   ✓ works — sets style="--val:0.6"
```
CSS `calc(var(--val, 0) * 100%)` then picks it up. BUT `calc(number * percentage)` requires
browser support for CSS math — use `left:X%` directly when possible:
```go
// Preferred — set left% directly from signal
expr := fmt.Sprintf("($%s-%d)/(%d-%d)*100", sig, min, max, min)
g.Attr("data-attr", fmt.Sprintf(`{"style":"left:clamp(1rem,"+(%s).toFixed(1)+"%%,calc(100%% - 1rem))"}`, expr))
```

**Don't use `data-attr` for static attributes** — use `g.Attr()` directly. Only use `data-attr` when the value is truly reactive (depends on a signal).

Wrong:
```go
g.Attr("data-attr", fmt.Sprintf(`{"data-align":$open?"none":"%s"}`, align))
// Bug: when open=true, sets data-align="none" (inverted logic)
```
Right:
```go
g.If(align == "right", g.Attr("data-align", "right"))  // static, no data-attr needed
```

### `data-bind` on range/number inputs

`data-bind:signal` on `<input type="range">` may only fire on `change` (mouse-up), not `input` (while dragging). For live value display, add `data-on:input` explicitly:
```go
g.Attr("data-bind:"+sig),
g.Attr("data-on:input", "$"+sig+"=Number(evt.target.value)"),
```

### Initial value of bound inputs

`data-bind:signal` sets input value FROM signal reactively, but the initial HTML render shows the element's default (`0` for number). Always set the HTML `value` attribute for the initial render:
```go
h.Input(
    h.Type("number"),
    g.Attr("data-bind:"+sig),
    h.Value(fmt.Sprintf("%d", p.Value)),  // ← required for correct initial display
)
```

## Server-side patterns

```go
import "github.com/starfederation/datastar-go/datastar"

func handler(c echo.Context) error {
    // Read client signals (GET: query params, POST: JSON body)
    var s State
    if err := datastar.ReadSignals(c.Request(), &s); err != nil {
        return echo.NewHTTPError(400, err.Error())
    }
    sse := datastar.NewSSE(c.Response().Writer, c.Request())

    // Patch DOM fragment (morphs by #id)
    sse.PatchElements(web.RenderToString(MyFragment(s)))

    // Update client signals  — NOTE: method on sse, NOT package-level function
    sse.MarshalAndPatchSignals(map[string]any{"loading": false})
    return nil
}

// Persistent SSE stream (live clock, real-time updates)
func streamHandler(c echo.Context) error {
    sse := datastar.NewSSE(c.Response().Writer, c.Request())
    ticker := time.NewTicker(time.Second)
    defer ticker.Stop()
    for {
        select {
        case <-c.Request().Context().Done():
            return nil
        case t := <-ticker.C:
            sse.MarshalAndPatchSignals(map[string]any{
                "serverTime": t.Format("15:04:05"),
            })
        }
    }
}
```

**Critical:** `sse.MarshalAndPatchSignals(map)` is a METHOD on the SSE object (not `datastar.MarshalAndPatchSignals(sse, map)`). Same for `sse.PatchElements(html)`.
Always check `c.Request().Context().Done()` in streaming loops — client disconnect returns nil.
```

Echo gzip middleware **must skip `/sse/*` and `/api/*` SSE routes** or it buffers and breaks streaming.

## Security

- CSP set in `internal/web/security.go`. **`'unsafe-eval'` is required** — Datastar v1.x evaluates all `data-*` expressions via `new Function()` at runtime. Removing it breaks all interactivity. `'unsafe-inline'` covers the pre-paint FOUC-prevention inline `<script>`. This is an accepted trade-off; Datastar has no precompile mode.
- Self-hosted fonts in `assets/static/fonts/`. No Google Fonts.
- altcha + honeypot on every public form. Verify server-side with `altcha.VerifySolution`.
- `gosec` runs in CI. `govulncheck` in pre-push.
- Rate-limit `/api/*` via Echo middleware.

## Testing

- **Render tests** — `gomponents.Render(buf, node)` + string assertions on `data-component`, `data-variant`. Grep `class=` and fail if found.
- **Handler tests** — `httptest` + `echo.New()`. For SSE, assert `text/event-stream` body contains `datastar-patch-elements` + expected fragment `id="…"`.
- **Registry test** — every registered component renders non-empty HTML at default props.
- **E2E** — Playwright (Node spec or `playwright-go`). Theme/mode toggle persistence, modal open/close, toast dismiss.

## Make targets cheat sheet

```
make setup                          # fetch tailwind + datastar.js + altcha.js + go mod tidy
make icons                          # regen ui/icon/icons_gen.go
make dev PROJECT=homepage           # tailwind --watch + air
make dev-showcase                   # PROJECT=showcase with -tags showcase
make build PROJECT=homepage         # static binary → bin/homepage
make check                          # fmt + vet + lint + guard-classes + test + vuln
make guard-classes                  # enforce no class= in ui/
make docker PROJECT=homepage TAG=v1 # per-project image
make upgrade-deps                   # bump tailwind + datastar.js + altcha.js to TAG vars
```

Version pins live at the top of the `Makefile`. To upgrade: edit `TAILWIND_VERSION` / `DATASTAR_VERSION` / `ALTCHA_VERSION`, run `make upgrade-deps`, commit the refreshed files.

## Animation library — Motion v10

**Decision: Motion v10 only.** Single library, ~24KB UMD, available as `window.Motion` in all showcase previews.

Vendored at: `projects/showcase/assets/static/motion.min.js`
Loaded in: `projects/showcase/pages/preview.go` via `<script src="/static/motion.min.js">`

**Key API:**
```js
// Animate element
const ctrl = Motion.animate(el, { x: 100, opacity: [0, 1], rotate: 180, scale: [0.8, 1] },
    { duration: 0.3, easing: 'ease-out', delay: 0.1 })

// Loop / sway (ping-pong)
Motion.animate(el, { x: [0, 50] }, { repeat: Infinity, direction: 'alternate', duration: 2 })

// Stagger
Motion.animate('.list-item', { opacity: [0, 1], y: [-8, 0] },
    { duration: 0.35, delay: Motion.stagger(0.06) })

// Timeline
Motion.timeline([
    ['.a', { opacity: [0,1] }, { duration: 0.25 }],
    ['.b', { opacity: [0,1] }, { duration: 0.25, at: '-0.1' }],
])

// Scroll reveal
Motion.inView('.section', ({target}) => Motion.animate(target, { opacity: [0,1] }, { duration: 0.5 }))

// Stop an animation
ctrl.stop()   // or ctrl.cancel()
```

**Important:** Motion uses **seconds** (not ms). `duration: 0.3` = 300ms.

**SVG scatter animation pattern** (used in logo showcase):
- Wrap each `<path>` in `<g style="transform-box:fill-box;transform-origin:center">`
- **CRITICAL:** `transform: translateX(Xpx)` on SVG child elements uses X as SVG USER UNITS, not CSS pixels. The `px` unit inside SVG coordinate space = 1 SVG unit (not 1 layout pixel).
- **Always multiply CSS pixel displacements by svgScale** before using in transform strings:
  ```js
  var svgScale = svg.viewBox.baseVal.width / svgR.width; // e.g. 2666/280 ≈ 9.52
  var dx = (target_x_css - center_x_css) * svgScale;   // → SVG user units
  ```
- Sway amplitudes, offsets — same scaling needed. Rotation (deg) and scale() are dimensionless — no scaling.
- Do NOT use Motion's `x`/`y` shorthand on SVG elements — those map to CSS custom properties `--x`/`--y` which SVG elements don't inherit. Use full `transform` string: `'translateX('+dx+'px) translateY('+dy+'px) rotate('+r+'deg) scale('+s+')'`
- `setTimeout(sway, dur*1000 + 100)` for post-scatter sway (reliable onComplete replacement)

**Do NOT vendor anime.js or GSAP** — breaks the lightweight single-library principle.

## What NOT to do

- Don't add `class=` to `ui/**.go`. Use `data-*` + CSS.
- Don't `g.Raw` user input. Use `g.Text` / `g.Textf`.
- Don't import `mljr-web/ui/registry` from prod (`projects/homepage`). Showcase only.
- Don't add a CDN `<link>` or `<script src="https://…">`. Vendor it.
- Don't dot-import `gomponents` in `ui/`. Use the `g` / `h` aliases.
- Don't add inline `style=` for theming. Use CSS vars + `data-tone`.
- Don't add a new theme without also adding both `-light` and `-dark` mode blocks.
- Don't bypass `Page()` / `PageShell()` — they set `data-theme` / `data-mode` on `<html>` and mount `Portal("portal")` for overlays.
- **Don't use `data-for`** — it does not exist in Datastar 1.0.2.
- Don't assume `data-bind` fires in real-time for range/number inputs — add `data-on:input` too.
- Don't rely on CSS `--tone`/`--on-tone` as the only visible indicator for a component variant. Always add a fallback background so the component is visible even when the variant attribute hasn't been set yet.
- Don't use `data-attr` for static attribute values — use `g.Attr()` directly. `data-attr` is only for reactive attributes that depend on a Datastar signal.
- Don't forget to set HTML `value` attribute on bound inputs — `data-bind` doesn't set the initial display on first render.
- Don't use a nested CSS import chain for Tailwind v4 — `projects/*/assets/css/input.css` must import all `ui/css/_*.css` partials directly (not through `core.css`) because Tailwind v4 may not follow nested `@import` chains for source scanning.

## Render a node to string (Go)

When you need an HTML fragment as a string (e.g. embedding SVG in JS):
```go
import "strings"
import g "maragu.dev/gomponents"

func renderStr(n g.Node) string {
    var b strings.Builder
    _ = n.Render(&b)
    return b.String()
}

// Escape for single-quoted JS string:
func jsStr(s string) string {
    s = strings.ReplaceAll(s, `\`, `\\`)
    s = strings.ReplaceAll(s, `'`, `\'`)
    s = strings.ReplaceAll(s, "\n", "")
    return s
}
```

Or use `internal/web.RenderToString(n)` which does the same via `bytes.Buffer`.

---

## Component checklist — COMPONENTS.md

**MANDATORY RULE: After every component add/remove/rename, update `COMPONENTS.md` in the repo root.**

`COMPONENTS.md` is the authoritative source of:
- All implemented components (slug, name, package, status)
- All planned but not-yet-built components
- Showcase infrastructure status

The full inventory lives there. Don't duplicate it here. When asked "what components exist?" or "what's missing?", read `COMPONENTS.md` first.

Current count (as of last update): **85 components** across primitive, layout, form, data, overlay, feedback, special, datastar, animation.

---

## gomponents — `g.Group` does NOT propagate attributes to parent elements

**Critical gotcha:** `primitive.Button(props, h.ID("foo"), g.Text("bar"))` does NOT set `id="foo"` on the button.

`primitive.Button` passes children via `g.Group(children)` to `h.Button(...)`. `g.Group` does not implement `placer` (the gomponents interface for attribute placement). So `h.ID("foo")` inside the group renders as **text content** inside the button, not as an attribute.

**Rule: Never pass attribute nodes as children of wrapper components.** Use one of:
1. `Attrs []g.Node` field on the Props struct — passed as `g.Group(p.Attrs)` (same issue!). Instead use direct `h.Button(h.ID("foo"), ...)` if you need IDs.
2. Wrap in a `h.Div(h.ID("foo"), ...)` outer element and listen on the div.
3. Use `h.Button(h.ID("foo"), g.Attr("data-component","button"), ...)` directly — bypass the wrapper.

This matters for: magnetic buttons (need `id` for JS), copy buttons, any component where JS needs to `getElementById` something inside a primitive wrapper.

---

## Showcase catalogue — iframe-per-component pattern

The catalogue overview (`projects/showcase/pages/catalogue.go`) renders each component in its own `<iframe src="/components/{slug}/preview?theme=...&mode=...">` rather than calling `c.Render()` inline.

**Why iframes:**
- Motion v10 and Datastar load fully in each preview (not available in the catalogue page itself)
- CSS isolation: component body styles can't bleed into the catalogue page
- No `overflow:hidden` breaking catalogue scroll (critical for the logo scatter animation)
- Theme/mode sync via `data-attr` → Datastar signal → iframe src rebuild:

```go
g.Attr("data-attr", fmt.Sprintf(`{"src":"/components/%s/preview?theme="+$theme+"&mode="+$mode}`, slug))
```

When `$theme` or `$mode` changes, Datastar rebuilds the src attribute → iframe reloads with correct theme.

**IMPORTANT:** The `Render()` function on `registry.Component` is still called inline in the detail page's variant matrix. Keep `Render()` functions safe to call without Motion/Datastar being loaded — they should render valid static HTML.

---

## Datastar expression conflicts with HTML attribute names

Datastar intercepts ALL `data-*` attributes it knows about. **Any attribute starting with `data-` + a registered plugin name will be evaluated as a JS expression.**

Known Datastar plugins that are easy to accidentally use for custom HTML data:
- `data-effect` — re-runs on signal change; `data-effect="tr-blur"` tries to evaluate `tr - blur` as JS
- `data-show` — evaluated as boolean expression
- `data-text` — sets textContent
- `data-bind` — two-way binding

**Rule:** For non-Datastar custom data attributes that happen to share a plugin name, use a namespace prefix:
```go
// WRONG — Datastar intercepts this
g.Attr("data-effect", "tr-blur")   // Error: tr is not defined

// CORRECT — unique name Datastar won't intercept
g.Attr("data-anim", "tr-blur")     // safe
g.Attr("data-x-effect", "tr-blur") // safe
```

---

## Showcase detail page — responsive layout

Detail page uses CSS classes for responsive stacking:
- `.detail-outer-grid` → collapses `220px 1fr` to `1fr` at ≤768px (hides nav sidebar)
- `.detail-controls-grid` → collapses `260px 1fr` to `1fr` at ≤768px

CSS lives in `ui/css/_layout.css`:
```css
@media (max-width: 768px) {
  .detail-outer-grid { grid-template-columns: 1fr !important; }
  .detail-outer-grid > nav[aria-label="Component navigation"] { display: none; }
  .detail-controls-grid { grid-template-columns: 1fr !important; }
}
```

Code blocks (`.usage-snippet`, controls panel `<code>`) use `overflow-x:auto` on the `<pre>` + `min-width:0` on the grid main area for mobile horizontal scrolling.

---

## Shared components across projects

`ui/special/logo_scatter.go` — `LogoScatter(LogoScatterProps)` — used by both:
1. `ui/datastar/animation_logo_showcase.go` (showcase, build tag `showcase`) — loop mode, 280px
2. `projects/homepage/pages/animatedlogo.go` (production) — scroll mode, 420px, with background div

`LogoScatterProps`:
```go
type LogoScatterProps struct {
    ID             string  // SVG id + gradient prefix (required, must be unique per page)
    SVGStyle       string  // full inline style; default = centered absolute
    Size           string  // CSS size (used only when SVGStyle is empty)
    InitialOpacity float64 // assembled opacity (default 0.65)
    Mode           string  // "loop" | "scroll"
    TriggerID      string  // scroll mode: IntersectionObserver element id
    WithBackground bool    // wrap in full-page abs div (id = ID+"-bg")
    WrapInLoad     bool    // wrap script in window.addEventListener('load', ...)
}
```

SVG gradient IDs are namespaced as `{ID}-lgN` to prevent collisions when multiple instances appear on one page.
