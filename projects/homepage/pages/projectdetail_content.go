package pages

// Deep-dive content for flagship project detail pages. Diagrams are Mermaid
// graph definitions; snippets are real excerpts trimmed from each project's
// own repository (godrive, homelab-automation, nightscout-tray) or from this
// monorepo (mljr-web).

var godriveDetail = projectDetailContent{
	LongDesc: `goDrive treats the filesystem as the source of truth — SQLite only stores metadata: sessions, the search index, trash records, and in-flight upload state. Everything else is rebuildable by reindexing the disk.

Uploads go through the TUS resumable protocol, so a dropped connection on a phone resuming a 4 GB video upload just picks back up rather than restarting. Once a file lands, thumbnails are generated asynchronously across several warmup sizes and cached under a key derived from the file's inode and device number rather than its path — renames and moves on the same filesystem don't invalidate the cache, only a tracked SQLite path lookup needs to update.

A webhook system lets external automation react to file events (upload.complete, file.moved, file.deleted, file.restored), each delivery signed with HMAC-SHA256 so a subscriber can verify it actually came from goDrive. The same Go backend serves a Svelte web UI, a WebDAV mount for native OS clients, and a Flutter mobile app with background upload support on both Android and iOS.`,
	Diagram: `graph TD
  Browser["Svelte / SVAR web UI"] -->|REST + SSE| API[Go backend]
  Mobile["Flutter app (Android/iOS)"] -->|TUS resumable upload| API
  WebDAV["Finder / Files / rclone"] -->|WebDAV| API
  API --> FS[("Filesystem\n(source of truth)")]
  API --> DB[("SQLite\nsessions · index · trash · upload state")]
  API --> Cache[("Inode-keyed\nthumbnail cache")]
  API -->|HMAC-SHA256 signed| Hooks["Webhook subscribers"]`,
	Snippets: []projectSnippet{
		{
			Caption:  "Thumbnails are cached by inode+device, not by path — a file renamed or moved on the same filesystem still hits its existing cached thumbnail instead of regenerating it.",
			Filename: "internal/server/thumbnail.go",
			Language: "go",
			Code: `func thumbnailCachePathInode(cacheRoot string, userID int64, logical string, info os.FileInfo, thumbSize int) string {
	inode, device := inodeKey(info)
	if inode == 0 {
		return thumbnailCachePath(cacheRoot, userID, logical, info.Size(), info.ModTime().UnixNano(), thumbSize)
	}
	sum := sha256.Sum256(fmt.Appendf(nil, "%d\x00%d\x00%d\x00%d\x00%d\x00%d\x00%d",
		thumbnailCacheVersion, userID, inode, device, info.Size(), info.ModTime().UnixNano(), thumbSize))
	return filepath.Join(cacheRoot, "thumbs", hex.EncodeToString(sum[:2]), hex.EncodeToString(sum[:])+".jpg")
}`,
		},
	},
}

var homelabAutomationDetail = projectDetailContent{
	LongDesc: `The whole homelab — Caddy ingress, CrowdSec, Authelia, Grafana/Loki/VictoriaMetrics, and every self-hosted service this portfolio runs alongside — is described as Ansible roles rather than clicked together by hand. Two hosts cover the stack: a public VPS (mljr) terminating TLS and fronting everything through Caddy, and a home NUC running heavier workloads and the monitoring stack, joined over a Tailscale mesh so neither host needs an open inbound port beyond Caddy's 443.

Services are entries in one YAML catalog: name, domain, port, which host runs it, and whether Caddy should sit it behind basic auth or Authelia SSO. A generic role turns that catalog into both a docker-compose deployment and an auto-generated Caddy reverse-proxy snippet — adding a new internet-facing service is a few lines of YAML, not a manual nginx/Caddy edit. Disabling a service in the catalog lets the same automation clean up its stale containers and proxy config instead of leaving orphaned config behind.

Deploys run from GitHub Actions over the Tailscale GitHub Action, with all secrets sealed in an Ansible Vault file that's still committed to git (only the ciphertext). A `+"`make deploy-check`"+` dry run verifies vault decryption and host connectivity before anything touches production, and the same playbooks run identically from a laptop for local iteration.`,
	Diagram: `graph TD
  GHA["GitHub Actions"] -->|Tailscale OAuth| TS[("Tailscale mesh")]
  Dev["make deploy-* (local)"] -->|Ansible Vault secrets| TS
  TS --> VPS["mljr (VPS)\nCaddy ingress · CrowdSec · Authelia"]
  TS --> NUC["nuc (home)\nGrafana · Loki · VictoriaMetrics"]
  VPS -->|reverse_proxy per catalog entry| Services1["Public services\n*.mljr.eu"]
  NUC --> Services2["Heavier workloads\nmonitoring stack"]
  Catalog[("services catalog\nall.yml")] -.generates.-> VPS
  Catalog -.generates.-> NUC`,
	Snippets: []projectSnippet{
		{
			Caption:  "One Jinja2 template turns a service's catalog entry (name, domain, port, auth type) into its Caddy reverse-proxy block — auth import and upstream picked automatically.",
			Filename: "ansible/roles/caddy/templates/service_snippet.caddy.j2",
			Language: "yaml",
			Code: `{% if service.enabled | default(true) and service.domain is defined %}
{% set is_local = service.host == inventory_hostname %}
{% set target_host = 'localhost' if is_local else hostvars[service.host]['ansible_host'] %}
{% set target_port = service.port %}
{% set auth_type = service.caddy_auth | default('') %}

{{ service.domain }} {
    log {
        output file {{ log_path }}/caddy/{{ service.name }}.log {
            roll_size {{ caddy_log_roll_size | default('100mb') }}
        }
    }
{% if auth_type == 'basicauth' %}
    import basicauth_block
{% elif auth_type == 'authelia' %}
    import authelia_auth
{% endif %}
    reverse_proxy {{ target_host }}:{{ target_port }}
}
{% endif %}`,
		},
	},
}

var mljrWebDetail = projectDetailContent{
	LongDesc: `This is the monorepo behind the site you're reading. Every component — gauges, the skills radar, the regex lab, this very project detail page — is a Go function returning HTML via gomponents; there's no JSX, no client-side virtual DOM, and no npm build step in the request path. Datastar (a single ~14 kB script) wires up reactivity by reading data-* attributes already present in the server-rendered HTML, and SSE patches fragments in place for things like the homelab panel without a page reload.

The data layer is a separate Go module (mljr-data) that runs as a scheduled job: it hits the GitHub GraphQL API for contribution stats and language breakdowns, scrapes Strava and LinkedIn where APIs allow it, and writes one site-data.json that the homepage binary embeds as a build-time fallback and can also hot-reload from disk in production. Four themes × two modes are pure CSS custom-property swaps — no separate component variants to maintain.`,
	Diagram: `graph TD
  Generator["mljr-data generator\n(scheduled job)"] -->|GitHub GraphQL, Strava, LinkedIn| JSON[("site-data.json")]
  JSON -->|embedded fallback + hot reload| Binary["homepage Go binary"]
  Binary -->|gomponents render| HTML["Server-rendered HTML"]
  HTML -->|data-* attributes| Datastar["Datastar (~14kB)"]
  Binary -->|SSE PatchElements| Datastar
  Datastar -->|in-place DOM patch| HTML`,
	Snippets: []projectSnippet{
		{
			Caption:  "The skills radar is plain trigonometry rendered to an SVG polygon at request time — no charting library.",
			Filename: "ui/data/radarchart.go",
			Language: "go",
			Code: `// Radar chart: axis i sits at angle 2πi/n, starting at
// 12 o'clock. Values scale along the spoke; the series
// becomes a single <polygon>.
angle := func(i int) float64 {
	return float64(i)*2*math.Pi/float64(n) - math.Pi/2
}
pt := func(radius float64, i int) (float64, float64) {
	a := angle(i)
	return cx + radius*math.Cos(a), cy + radius*math.Sin(a)
}

pts := make([]string, len(p.Axes))
for i := range p.Axes {
	scaled := (s.Values[i] / p.Max) * r
	x, y := pt(scaled, i)
	pts[i] = fmt.Sprintf("%.1f,%.1f", x, y)
}`,
		},
		{
			Caption:  "Live updates without a framework: one attribute on the section, one SSE handler, Datastar morphs the panel in place.",
			Filename: "projects/homepage/handlers.go",
			Language: "go",
			Code: `// One attribute on the section…
h.Section(
	h.ID("homelab"),
	g.Attr("data-on-interval__duration.60s", "@get('/api/homelab')"),
)

// …one SSE handler on the server.
e.GET("/api/homelab", func(c echo.Context) error {
	sse := datastar.NewSSE(c.Response().Writer, c.Request())
	return sse.PatchElements(
		web.RenderToString(pages.HomelabPanel(snapshot())),
	)
})`,
		},
	},
}

var nightscoutTrayDetail = projectDetailContent{
	LongDesc: `nightscout-tray is a Wails3 app (Go backend, native webview frontend) that turns a continuous glucose monitor's Nightscout feed into a system tray icon. Every poll, it rasterizes a fresh 32×32 PNG on the fly with the gg 2D graphics library: a rounded rect colored by urgency status, the current glucose value as anchored text, and a direction arrow for the trend — recomputed each tick rather than swapped between a fixed set of pre-rendered icons.

A separate prediction package implements an OpenAPS-style oref algorithm against recent glucose history to estimate where the value is heading, feeding configurable threshold notifications. The same binary builds natively for Windows, macOS, and Linux through Wails3's cross-platform webview bridge, with mobile builds (Android/iOS) cross-compiled through the project's own CI pipeline.`,
	Diagram: `graph TD
  Nightscout[("Nightscout API")] -->|poll| Service["internal/app/service.go"]
  Service --> Predictor["oref-style predictor\n(internal/prediction)"]
  Service --> IconGen["IconGenerator\n(gg rasterizer)"]
  Predictor -->|trend + forecast| Notifications["Threshold notifications"]
  IconGen -->|fresh PNG per tick| Tray["System tray icon"]
  Service --> Tray`,
	Snippets: []projectSnippet{
		{
			Caption:  "The tray icon isn't a static asset — it's rasterized fresh on every poll: status color, the glucose value as text, and a trend arrow drawn onto a 32×32 canvas.",
			Filename: "internal/tray/icon.go",
			Language: "go",
			Code: `func (g *IconGenerator) GenerateIcon(text string, direction string, status *models.GlucoseStatus) []byte {
	// Wails v3 systray uses PNG on all platforms
	var width, height float64 = 32, 32
	radius := width / 4

	dc := gg.NewContext(int(width), int(height))
	dc.SetRGBA(0, 0, 0, 0)
	dc.Clear()

	bgHex := getStatusColor(status)
	r, ge, b := parseHexColor(bgHex)
	dc.SetRGB255(int(r), int(ge), int(b))
	dc.DrawRoundedRectangle(0, 0, width, height, radius)
	dc.Fill()

	brightness := (int(r)*299 + int(ge)*587 + int(b)*114) / 1000
	if brightness > 128 {
		dc.SetColor(color.Black)
	} else {
		dc.SetColor(color.White)
	}

	fontSize := height * 0.5
	if err := loadFont(dc, fontSize); err == nil {
		dc.DrawStringAnchored(text, width/2, height/2-2, 0.5, 0.5)
	}
	if direction != "" {
		drawArrow(dc, width/2, height-5, height*0.3, direction)
	}

	var buf bytes.Buffer
	png.Encode(&buf, dc.Image())
	return buf.Bytes()
}`,
		},
	},
}
