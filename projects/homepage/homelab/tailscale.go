package homelab

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// cacheRefresh controls how often the public inventory/service-registry
// files are re-fetched. They change rarely, so there's no need to hit
// GitHub on every 60s poll cycle.
const cacheRefresh = 10 * time.Minute

// ServiceEntry is one entry from the public homelab-automation service
// registry (ansible/inventory/group_vars/all/all.yml's `services:` list).
type ServiceEntry struct {
	Name        string
	Domain      string
	Host        string
	Description string
}

// fetchTailscale builds the mesh view: live Tailscale device state, filtered
// down to the infra hosts listed in the public homelab-automation inventory
// (so personal devices — phones, laptops — never appear), annotated with
// each host's services from the public service registry.
func (p *Poller) fetchTailscale(ctx context.Context, s *Snapshot) {
	if time.Since(p.invFetchedAt) > cacheRefresh || p.invCache == nil {
		if inv, err := p.fetchInventory(ctx); err != nil {
			log.Printf("homelab: inventory fetch failed: %v", err)
		} else {
			p.invCache = inv
			p.invFetchedAt = time.Now()
		}
	}
	if time.Since(p.svcFetchedAt) > cacheRefresh || p.svcCache == nil {
		if svc, err := p.fetchServices(ctx); err != nil {
			log.Printf("homelab: services fetch failed: %v", err)
		} else {
			p.svcCache = svc
			p.svcFetchedAt = time.Now()
		}
	}
	if len(p.invCache) == 0 {
		return // can't filter devices down to infra hosts without this — skip the cycle
	}

	devices, err := p.fetchTailscaleDevices(ctx)
	if err != nil {
		p.alertTailscaleFailure(ctx, err)
		return
	}
	p.lastTSAlertAt = time.Time{} // recovered — un-throttle the next failure

	ipToHost := make(map[string]string, len(p.invCache))
	for host, ip := range p.invCache {
		ipToHost[ip] = host
	}
	servicesByHost := map[string][]MeshService{}
	for _, svc := range p.svcCache {
		servicesByHost[svc.Host] = append(servicesByHost[svc.Host], MeshService{
			Name:        svc.Name,
			Domain:      svc.Domain,
			Description: svc.Description,
		})
	}

	mesh := make([]MeshHost, 0, len(p.invCache))
	for _, d := range devices {
		var ip string
		for _, a := range d.Addresses {
			if strings.HasPrefix(a, "100.") {
				ip = a
				break
			}
		}
		host, known := ipToHost[ip]
		if !known {
			continue // not in the public inventory — personal device, drop it
		}
		online := false
		if lastSeen, err := time.Parse(time.RFC3339, d.LastSeen); err == nil {
			online = time.Since(lastSeen) < 5*time.Minute
		}
		mesh = append(mesh, MeshHost{
			Name:        host,
			TailscaleIP: ip,
			OS:          d.OS,
			Online:      online,
			Relay:       d.ClientConnectivity.DERP != "",
			Services:    servicesByHost[host],
		})
	}
	sort.Slice(mesh, func(a, b int) bool { return mesh[a].Name < mesh[b].Name })
	s.Mesh = mesh
	s.MeshOK = len(mesh) > 0
}

// alertTailscaleFailure logs the failure and, throttled to once an hour,
// pings ntfy — the Tailscale API key has a hard 90-day max lifetime, so this
// will eventually fire from expiry alone, not just outages.
func (p *Poller) alertTailscaleFailure(ctx context.Context, err error) {
	log.Printf("homelab: tailscale devices fetch failed: %v", err)
	if !p.lastTSAlertAt.IsZero() && time.Since(p.lastTSAlertAt) < time.Hour {
		return
	}
	p.lastTSAlertAt = time.Now()
	p.notify(ctx, "homepage: Tailscale API failing",
		fmt.Sprintf("Mesh panel can't reach the Tailscale API (check for an expired key — they cap out at 90 days): %v", err))
}

func (p *Poller) notify(ctx context.Context, title, message string) {
	if p.ntfyURL == "" || p.ntfyTopic == "" {
		return
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.ntfyURL+"/"+p.ntfyTopic, strings.NewReader(message))
	if err != nil {
		log.Printf("homelab: ntfy request build failed: %v", err)
		return
	}
	req.Header.Set("Title", title)
	req.Header.Set("Priority", "high")
	resp, err := p.client.Do(req)
	if err != nil {
		log.Printf("homelab: ntfy send failed: %v", err)
		return
	}
	defer func() {
		if err := resp.Body.Close(); err != nil && !errors.Is(err, context.Canceled) {
			log.Printf("homelab: ntfy response body close failed: %v", err)
		}
	}()
	if resp.StatusCode >= 300 {
		log.Printf("homelab: ntfy responded with status %d", resp.StatusCode)
	}
}

// ── Tailscale API ───────────────────────────────────────────────────────────

type tsDevice struct {
	Hostname           string   `json:"hostname"`
	Addresses          []string `json:"addresses"`
	OS                 string   `json:"os"`
	LastSeen           string   `json:"lastSeen"`
	ClientConnectivity struct {
		DERP string `json:"derp"`
	} `json:"clientConnectivity"`
}

type tsDevicesResponse struct {
	Devices []tsDevice `json:"devices"`
}

func (p *Poller) fetchTailscaleDevices(ctx context.Context) ([]tsDevice, error) {
	u := "https://api.tailscale.com/api/v2/tailnet/" + url.PathEscape(p.tsTailnet) + "/devices?fields=all"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+p.tsAPIKey)
	req.Header.Set("Accept", "application/json")
	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil && !errors.Is(err, context.Canceled) {
			log.Printf("homelab: response body close failed: %v", err)
		}
	}()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("tailscale api status %d", resp.StatusCode)
	}
	var out tsDevicesResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return out.Devices, nil
}

// ── Public inventory (homelab-automation, raw YAML) ────────────────────────

type invNode struct {
	Children map[string]invNode `yaml:"children"`
	Hosts    map[string]struct {
		TailscaleIP string `yaml:"tailscale_ip"`
	} `yaml:"hosts"`
}

type inventoryFile struct {
	All invNode `yaml:"all"`
}

func collectInventoryHosts(n invNode, out map[string]string) {
	for name, h := range n.Hosts {
		if h.TailscaleIP != "" {
			out[name] = h.TailscaleIP
		}
	}
	for _, child := range n.Children {
		collectInventoryHosts(child, out)
	}
}

func (p *Poller) fetchInventory(ctx context.Context) (map[string]string, error) {
	body, err := p.getBytes(ctx, p.inventoryURL)
	if err != nil {
		return nil, err
	}
	var inv inventoryFile
	if err := yaml.Unmarshal(body, &inv); err != nil {
		return nil, err
	}
	out := map[string]string{}
	collectInventoryHosts(inv.All, out)
	return out, nil
}

// ── Public service registry (homelab-automation, raw YAML) ─────────────────

type rawService struct {
	Name        string `yaml:"name"`
	Enabled     *bool  `yaml:"enabled"`
	Domain      any    `yaml:"domain"` // string or []string in the source file
	Host        string `yaml:"host"`
	Description string `yaml:"description"`
}

type servicesFile struct {
	Services []rawService `yaml:"services"`
}

func firstDomain(v any) string {
	switch t := v.(type) {
	case string:
		return t
	case []any:
		if len(t) > 0 {
			if s, ok := t[0].(string); ok {
				return s
			}
		}
	}
	return ""
}

func (p *Poller) fetchServices(ctx context.Context) ([]ServiceEntry, error) {
	body, err := p.getBytes(ctx, p.servicesURL)
	if err != nil {
		return nil, err
	}
	var sf servicesFile
	if err := yaml.Unmarshal(body, &sf); err != nil {
		return nil, err
	}
	out := make([]ServiceEntry, 0, len(sf.Services))
	for _, s := range sf.Services {
		if s.Enabled != nil && !*s.Enabled {
			continue
		}
		out = append(out, ServiceEntry{
			Name:        s.Name,
			Domain:      firstDomain(s.Domain),
			Host:        s.Host,
			Description: s.Description,
		})
	}
	return out, nil
}

func (p *Poller) getBytes(ctx context.Context, rawURL string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil && !errors.Is(err, context.Canceled) {
			log.Printf("homelab: response body close failed: %v", err)
		}
	}()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status %d", resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}
