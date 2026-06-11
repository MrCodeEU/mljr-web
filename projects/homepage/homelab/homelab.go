// Package homelab polls live infrastructure data for the homepage panel:
// Uptime Kuma's public status-page API and PromQL stats (CrowdSec, hosts).
// A background goroutine refreshes an in-memory snapshot; visitors never
// trigger upstream calls.
package homelab

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"sync"
	"time"
)

type Service struct {
	Name  string
	Group string
	Up    bool
	Ping  float64 // last response time in ms, 0 = unknown
}

type Snapshot struct {
	FetchedAt time.Time

	// Uptime Kuma
	KumaOK   bool
	Services []Service
	UpCount  int

	// Avg response time across all monitors, oldest→newest (for sparkline).
	PingHistory []float64

	// PromQL stats; -1 = unavailable
	ActiveBans     int
	Attacks24h     int
	HostsOnline    int
	SecurityEvents int // total alert scenarios triggered since CrowdSec start
}

// Sample returns a plausible snapshot for dev environments without
// Tailscale access to the homelab.
func Sample() Snapshot {
	return Snapshot{
		FetchedAt: time.Now(),
		KumaOK:    true,
		Services: []Service{
			{Name: "Caddy ingress", Group: "Core", Up: true, Ping: 12},
			{Name: "Authelia", Group: "Core", Up: true, Ping: 23},
			{Name: "Uptime Kuma", Group: "Core", Up: true, Ping: 8},
			{Name: "Grafana", Group: "Monitoring", Up: true, Ping: 41},
			{Name: "Umami", Group: "Monitoring", Up: true, Ping: 35},
			{Name: "Mailcow", Group: "Apps", Up: true, Ping: 64},
			{Name: "ntfy", Group: "Apps", Up: true, Ping: 19},
			{Name: "Home Assistant", Group: "Apps", Up: false, Ping: 0},
		},
		UpCount:        7,
		PingHistory:    []float64{31, 29, 35, 28, 33, 41, 30, 27, 29, 38, 32, 28, 26, 31, 36, 29, 27, 30, 33, 28},
		ActiveBans:     142,
		Attacks24h:     37,
		HostsOnline:    3,
		SecurityEvents: 2841,
	}
}

type Poller struct {
	kumaURL  string
	kumaSlug string
	promURL  string
	client   *http.Client

	mu   sync.RWMutex
	snap Snapshot
}

// New creates a poller. Empty kumaURL or promURL disables that source.
func New(kumaURL, kumaSlug, promURL string) *Poller {
	if kumaSlug == "" {
		kumaSlug = "all"
	}
	return &Poller{
		kumaURL:  kumaURL,
		kumaSlug: kumaSlug,
		promURL:  promURL,
		client:   &http.Client{Timeout: 10 * time.Second},
	}
}

// Snapshot returns the latest snapshot (zero value before the first poll).
func (p *Poller) Snapshot() Snapshot {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.snap
}

// Start polls immediately, then every interval until ctx is done.
func (p *Poller) Start(ctx context.Context, interval time.Duration) {
	go func() {
		p.poll(ctx)
		t := time.NewTicker(interval)
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				p.poll(ctx)
			}
		}
	}()
}

func (p *Poller) poll(ctx context.Context) {
	next := Snapshot{
		FetchedAt:      time.Now(),
		ActiveBans:     -1,
		Attacks24h:     -1,
		HostsOnline:    -1,
		SecurityEvents: -1,
	}

	if p.kumaURL != "" {
		if err := p.fetchKuma(ctx, &next); err != nil {
			log.Printf("homelab: kuma fetch failed: %v", err)
		}
	}
	if p.promURL != "" {
		p.fetchProm(ctx, &next)
	}

	p.mu.Lock()
	prev := p.snap
	// Keep last good Kuma data through transient failures.
	if !next.KumaOK && prev.KumaOK {
		next.KumaOK = true
		next.Services = prev.Services
		next.UpCount = prev.UpCount
		next.PingHistory = prev.PingHistory
		next.FetchedAt = prev.FetchedAt
	}
	p.snap = next
	p.mu.Unlock()
}

// ── Uptime Kuma public status page ────────────────────────────────────────────

type kumaStatusPage struct {
	PublicGroupList []struct {
		Name        string `json:"name"`
		MonitorList []struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
		} `json:"monitorList"`
	} `json:"publicGroupList"`
}

type kumaBeat struct {
	Status int      `json:"status"` // 1 up, 0 down, 2 pending, 3 maintenance
	Ping   *float64 `json:"ping"`
}

type kumaHeartbeats struct {
	HeartbeatList map[string][]kumaBeat `json:"heartbeatList"`
}

func (p *Poller) fetchKuma(ctx context.Context, s *Snapshot) error {
	var page kumaStatusPage
	if err := p.getJSON(ctx, p.kumaURL+"/api/status-page/"+p.kumaSlug, &page); err != nil {
		return fmt.Errorf("status page: %w", err)
	}
	var beats kumaHeartbeats
	if err := p.getJSON(ctx, p.kumaURL+"/api/status-page/heartbeat/"+p.kumaSlug, &beats); err != nil {
		return fmt.Errorf("heartbeats: %w", err)
	}

	maxHistory := 0
	for _, group := range page.PublicGroupList {
		for _, m := range group.MonitorList {
			hb := beats.HeartbeatList[strconv.FormatInt(m.ID, 10)]
			svc := Service{Name: m.Name, Group: group.Name}
			if len(hb) > 0 {
				last := hb[len(hb)-1]
				svc.Up = last.Status == 1
				if last.Ping != nil {
					svc.Ping = *last.Ping
				}
			}
			if svc.Up {
				s.UpCount++
			}
			if len(hb) > maxHistory {
				maxHistory = len(hb)
			}
			s.Services = append(s.Services, svc)
		}
	}

	// Average ping per beat index (aligned from the newest backwards).
	const sparkPoints = 30
	if maxHistory > sparkPoints {
		maxHistory = sparkPoints
	}
	for i := maxHistory; i >= 1; i-- {
		sum, n := 0.0, 0
		for _, group := range page.PublicGroupList {
			for _, m := range group.MonitorList {
				hb := beats.HeartbeatList[strconv.FormatInt(m.ID, 10)]
				if len(hb) >= i {
					b := hb[len(hb)-i]
					if b.Ping != nil && *b.Ping > 0 {
						sum += *b.Ping
						n++
					}
				}
			}
		}
		if n > 0 {
			s.PingHistory = append(s.PingHistory, sum/float64(n))
		}
	}

	// Stable order: by group, then name.
	sort.SliceStable(s.Services, func(a, b int) bool {
		if s.Services[a].Group != s.Services[b].Group {
			return s.Services[a].Group < s.Services[b].Group
		}
		return s.Services[a].Name < s.Services[b].Name
	})

	s.KumaOK = len(s.Services) > 0
	return nil
}

// ── PromQL ────────────────────────────────────────────────────────────────────

func (p *Poller) fetchProm(ctx context.Context, s *Snapshot) {
	queries := []struct {
		dst   *int
		query string
	}{
		{&s.ActiveBans, `sum(cs_active_decisions)`},
		{&s.Attacks24h, `sum(increase(cs_bucket_overflowed_total[24h]))`},
		{&s.SecurityEvents, `sum(cs_alerts)`},
		{&s.HostsOnline, `count(count by (host) (up == 1))`},
	}
	for _, q := range queries {
		v, err := p.promQuery(ctx, q.query)
		if err != nil {
			log.Printf("homelab: promql %q failed: %v", q.query, err)
			continue
		}
		*q.dst = int(v)
	}
}

func (p *Poller) promQuery(ctx context.Context, query string) (float64, error) {
	var res struct {
		Status string `json:"status"`
		Data   struct {
			Result []struct {
				Value [2]any `json:"value"`
			} `json:"result"`
		} `json:"data"`
	}
	u := p.promURL + "/api/v1/query?query=" + url.QueryEscape(query)
	if err := p.getJSON(ctx, u, &res); err != nil {
		return 0, err
	}
	if res.Status != "success" || len(res.Data.Result) == 0 {
		return 0, fmt.Errorf("no data")
	}
	str, ok := res.Data.Result[0].Value[1].(string)
	if !ok {
		return 0, fmt.Errorf("unexpected value type")
	}
	return strconv.ParseFloat(str, 64)
}

func (p *Poller) getJSON(ctx context.Context, rawURL string, dst any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	resp, err := p.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status %d", resp.StatusCode)
	}
	return json.NewDecoder(resp.Body).Decode(dst)
}
