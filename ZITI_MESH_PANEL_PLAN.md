# Plan: replace the Tailscale mesh panel with an OpenZiti one

Status: plan, written 2026-11-24. Do it in its own session, in this repo,
then do the small cleanup in homelab-automation.

## Why

The homelab left Tailscale for OpenZiti (homelab-automation
`docs/ZERO_TRUST_PLAN.md`, Phase 5). Every host-to-host and admin path now
runs over Ziti. tailscaled is about to be uninstalled, so the homepage's
"C. Tailscale-Mesh" table goes empty. The Tailscale API key also expires
every 90 days, and an expired key pages ntfy. The texts that explain the
setup ("live via Tailscale", the architecture box "3 · Tailscale ·
WireGuard") are wrong now too.

## What the new panel shows

The data is already in VictoriaMetrics. `ziti-mesh-exporter` runs on mljr:
it reads the Ziti controller's management API with a read-only identity,
publishes metrics and VictoriaMetrics scrapes them (label `host="mljr"`).
The homepage already runs PromQL against VictoriaMetrics
(`HOMELAB_PROM_URL=http://ziti-proxy:19090`, over Ziti through the
`homepage-mljr` identity). So the panel needs **no new credentials, no API
key and no new network path**.

Metrics to use (live values from 2026-11-24 in brackets):

| Metric | Labels | Meaning |
|---|---|---|
| `ziti_mesh_router_online` | `router` (mljr, nuc, nas, ugreen) | 1 if the edge router is connected to the controller |
| `ziti_mesh_link_up` | `source`, `target` | 1 per router-to-router link (5 links: nas>mljr, nas>nuc, nuc>mljr, ugreen>mljr, ugreen>nuc) |
| `ziti_mesh_link_latency_seconds` | `source`, `target`, `measured_by` | link latency as measured from each end (nas>nuc ~5 ms, nas>mljr ~25 ms); the exporter already drops values above 60 s |
| `ziti_mesh_terminators` | `router`, `service` | services hosted behind each router |
| `ziti_mesh_circuits` | `client`, `service`, `ingress`, `egress` | open circuits (live connections) |

Do **not** use `ziti_mesh_identity_online`, `ziti_mesh_node` or
`ziti_mesh_edge` (the Grafana node-graph metrics). They carry identity
names such as `phone-mrcode` and `laptop-mrcode`: personal devices must
never appear on the public page, the same rule as the old
inventory filter. Of `ziti_mesh_circuits` and `ziti_mesh_terminators`,
publish only counts. Service and client names include internal ones
(`facts-ssh-*`, `backup-sftp-ugreen`, `hawser-*`) that need not be public.

PromQL (all instant queries, through the existing `promQueryVector`
helper or a sibling of it):

```promql
max by (router) (ziti_mesh_router_online)                     # status per router
max by (source, target) (ziti_mesh_link_up)                    # links
max by (source, target) (ziti_mesh_link_latency_seconds)       # worse end of each link
count by (router) (ziti_mesh_terminators > 0)                  # services per router
sum(ziti_mesh_circuits)                                        # open connections, one number
```

Panel layout (same cell and span as now, `lt-table`):

- Kicker: `C. OpenZiti-Mesh` with `4/4 Router online`.
- Table rows, one per router: **Host** (mljr, nuc, nas, ugreen),
  **Rolle / Role** (mljr: `Controller · öffentlich / public`, others:
  `Heimnetz / home, nur ausgehend / outbound only`), **Status**
  (`● online` / `○ offline`), **Links** (number up, e.g. `2/2`),
  **Dienste / Services** (terminator count).
- One line under the table: the links with their latency, e.g.
  `nas ↔ nuc 5 ms · nas ↔ mljr 25 ms · …`, sorted by name. Links are
  dialled outbound from home to mljr, so no router at home listens on the
  internet.
- Optional third figure: `N offene Verbindungen / open circuits` from
  `sum(ziti_mesh_circuits)`.

Rename the Go types to match: `MeshHost{Name, Role, Online, LinksUp,
LinksTotal, Services int}`, plus `MeshLink{Source, Target string, Latency
time.Duration}` and `Snapshot.MeshLinks`. The `TailscaleIP`, `OS` and
`Relay` fields go: Ziti has no per-host overlay IP, the exporter has no
OS, and nothing goes through a relay (links are direct router-to-router
TLS).

Failure behaviour: keep the current rule. A failed or empty query keeps
the last good mesh in the snapshot (the comment at `homelab.go:299`); the
first poll without data leaves `MeshOK=false` and the cell hidden. No ntfy
alert: VictoriaMetrics failures already show in the other figures, and the
exporter has its own monitoring in homelab-automation.

## Steps in mljr-web

1. **Data** (`projects/homepage/homelab/`):
   - New `ziti.go` with `fetchZitiMesh(ctx, *Snapshot)` that runs the
     queries above and fills `Mesh` and `MeshLinks`. Router role comes from
     a small constant map (`mljr` = controller + public router), not from
     the inventory.
   - Delete `tailscale.go` and `tailscale_test.go`, with
     `fetchTailscaleDevices`, `alertTailscaleFailure` and the inventory
     cache (`invCache`, `invFetchedAt`, `fetchInventory`). Check whether
     `fetchServices` (`HOMELAB_SERVICES_URL`, the catalog) is still used
     by anything else; if it only fed the mesh badges, delete it too.
   - `homelab.go`: drop `TailscaleAPIKey`, `TailscaleTailnet`,
     `InventoryURL` (and `ServicesURL` if unused) from `Options` and the
     poller; call `fetchZitiMesh` where `fetchTailscale` was called
     (line ~286). The mesh needs only `promURL`, so it is on whenever
     PromQL is.
   - `Sample()`: new sample mesh with the four routers and five links
     (values from the table above), so dev mode still renders.
   - Tests: `ziti_test.go` with an `httptest` VictoriaMetrics answering
     the five queries. Cover: all online; nas offline (its links count as
     down); empty answers keep the previous mesh; no identity names reach
     the snapshot.
2. **Config** (`internal/config/config.go`, `projects/homepage/main.go`):
   remove the Tailscale fields and env vars (`TAILSCALE_API_KEY`,
   `TAILSCALE_TAILNET`, `HOMELAB_INVENTORY_URL`). Fix the `PromURL`
   comment: VictoriaMetrics is reached over Ziti, not Tailscale.
3. **Page** (`projects/homepage/pages/homelab.go`): the new table and link
   line from the layout above; kicker `C. OpenZiti-Mesh`.

## Rewrite the explanation texts

All of these describe the old setup and get rewritten, in German and
English where the page is bilingual:

- `internal/i18n/locales/de.json` / `en.json`, `homelab.sub`:
  `live via Tailscale` -> `live über OpenZiti` / `live via OpenZiti`.
- `pages/homelab.go`, `archFlow()`, box 3. Today
  `{"3", "Tailscale", "Verschlüsseltes Mesh" / "Encrypted mesh", "WireGuard", ["keine offenen Ports daheim" / "no open ports at home"]}`.
  New:
  - title `OpenZiti`, role `Zero-Trust-Overlay` / `Zero-trust overlay`,
    place `mTLS`;
  - pills: `nur ausgehend von daheim` / `outbound only from home`,
    `Zugriff pro Dienst & Identität` / `access per service & identity`,
    `Default-Deny, Policies als Code` / `default deny, policies as code`.
- `archFlow()`, box 2 (mljr): add `Ziti-Controller` to its pills, so it
  is clear where the overlay is run.
- `archFlow()`, boxes 4 and 5: unchanged hosts. Add a box for `ugreen`
  (`Backup-NAS`, Thaling, pills `Backups`, `Syncthing`) only if the flow
  stays readable on a phone; otherwise mention it in box 5's pills.
- `archFlow()`, loop box 6: the agents still apply the same repo, but they
  are now triggered over Ziti, and CI has no access to the hosts. Suggested
  text:
  "OpenVox-Agenten (Puppet) auf jedem Host wenden dasselbe Repo an —
  Härtung, Docker-Compose-Dienste, Caddy, Ziti-Policies (OpenTofu),
  Monitoring. CI hat keinen Zugang zu den Hosts, kein SSH-Fan-out, keine
  Schneeflocken." /
  "OpenVox (Puppet) agents on every host apply the same repo — hardening,
  Docker Compose services, Caddy, Ziti policies (OpenTofu), monitoring. CI
  has no access to the hosts, no SSH fan-out, no snowflakes."
- `projects/homepage/data/types.go:467`, the "Infra / Homelab" skill row:
  `Tailscale` -> `OpenZiti`. Line 907, icon map: add an `openziti` entry
  (check `tools/icongen/icons.txt` for `simple-icons:openziti`; if it does
  not exist, use `lucide:network`) and regenerate the icons with
  `tools/icongen`. Keep the `tailscale` icon only if something else still
  uses it.
- `agents.md` lines ~704-708 (homelab live panel): sources are Kuma and
  PromQL over Ziti; the mesh comes from `ziti_mesh_*` metrics; the
  architecture diagram is `archFlow()` (the text still says `archCard()`,
  a "dashed Tailscale-mesh box" and an "Ansible IaC bar", all outdated).
- `PLAN.md:141`: the example `HOMELAB_PROM_URL` with `nuc.tail33930.ts.net`
  -> `http://ziti-proxy:19090` (over Ziti).
- Read the whole homelab section on the live page once more after the
  change, in both languages, for anything else that still says Tailscale
  or WireGuard.

## Steps in homelab-automation (after the new image is live)

Only once the deployed homepage no longer reads the Tailscale variables:

1. `services/homepage/docker-compose.yml`: remove `TAILSCALE_API_KEY`,
   `TAILSCALE_TAILNET`, `HOMELAB_INVENTORY_URL` (and
   `HOMELAB_SERVICES_URL` if mljr-web dropped it) plus their comment;
   `make sync-openvox-services`.
2. `openvox/modules/roles/manifests/services.pp`: the two
   `tailscale_api_key` lookups (homepage on mljr and nuc, lines ~156 and
   ~281). `openvox/modules/roles/templates/services/env.epp`: the
   `TAILSCALE_API_KEY=` line.
3. `openvox/data/secrets/mljr.int.mljr.eu.eyaml` and
   `nuc.int.mljr.eu.eyaml`: remove `vault_homepage_tailscale_api_key`.
4. `openvox/data/common.yaml`: remove the `tailscale_hosts` block and its
   comment (nothing else reads it since #726).
5. Deploy mljr, then nuc, one at a time.
6. Owner: revoke the API key in the Tailscale admin console.

## Done when

- The live homepage shows `C. OpenZiti-Mesh` with four routers, their
  links and latency, and no identity or personal device names.
- Stopping a router (or the exporter) does not blank the cell at once;
  the previous mesh stays until data comes back.
- No `Tailscale`, `WireGuard` or `tail33930` left in mljr-web outside
  history (`git grep -i -e tailscale -e wireguard -e tail33930`), except
  an icon that something else still uses.
- homelab-automation has no `TAILSCALE_*` variable, no
  `vault_homepage_tailscale_api_key` and no `tailscale_hosts`.
