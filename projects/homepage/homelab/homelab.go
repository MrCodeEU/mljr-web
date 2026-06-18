// Package homelab polls live infrastructure data for the homepage panel:
// Uptime Kuma's public status-page API and PromQL stats (CrowdSec, hosts).
// A background goroutine refreshes an in-memory snapshot; visitors never
// trigger upstream calls.
package homelab

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
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

	// CrowdSec detail; -1 = unavailable
	BansCommunity int         // active bans from the community blocklist (CAPI/lists)
	BansLocal     int         // active bans detected by local scenarios
	TopThreats    []NameValue // top categories by active decisions, desc

	// PromQL range data
	AttackDays []DayValue // attacks blocked per day, oldest→newest (≤365)
	CPUUtil    []float64  // avg CPU % across hosts, last 24h, oldest→newest
	CPULabels  []string   // hour labels aligned to CPUUtil

	// Tailscale mesh, filtered to infra hosts listed in the public
	// homelab-automation inventory (personal devices never appear here).
	MeshOK bool
	Mesh   []MeshHost
}

// MeshHost is one infra host in the Tailscale mesh view.
type MeshHost struct {
	Name        string // inventory hostname, e.g. "mljr"
	TailscaleIP string
	OS          string // OS family only (linux/windows/macOS/...), no version
	Online      bool
	Relay       bool // true if this host is currently reachable only via a DERP relay
	Services    []MeshService
}

// MeshService is one service hosted on a MeshHost (from the public
// homelab-automation service registry).
type MeshService struct {
	Name        string
	Domain      string
	Description string
}

// NameValue is a labeled counter (e.g. threat category → active bans).
type NameValue struct {
	Name  string
	Value int
}

// DayValue is one day of an aggregated counter (for heatmaps).
type DayValue struct {
	Date  time.Time
	Count int
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
		BansCommunity:  128,
		BansLocal:      14,
		TopThreats: []NameValue{
			{Name: "http:scan", Value: 78},
			{Name: "ssh:bruteforce", Value: 31},
			{Name: "http:bruteforce", Value: 17},
			{Name: "http:exploit", Value: 9},
			{Name: "http:crawl", Value: 7},
		},
		AttackDays: sampleAttackDays(),
		CPUUtil:    []float64{12, 14, 11, 18, 22, 19, 15, 13, 24, 31, 28, 21, 17, 14, 16, 19, 23, 26, 20, 15, 13, 12, 14, 17},
		CPULabels:  sampleHourLabels(24),
		MeshOK:     true,
		Mesh:       sampleMesh(),
	}
}

func sampleMesh() []MeshHost {
	return []MeshHost{
		{Name: "mljr", TailscaleIP: "100.100.20.1", OS: "linux", Online: true, Services: []MeshService{
			{Name: "authelia", Domain: "auth.mljr.eu", Description: "Authelia SSO"},
			{Name: "ntfy", Domain: "ntfy.mljr.eu", Description: "Push Notifications"},
		}},
		{Name: "nuc", TailscaleIP: "100.100.10.1", OS: "linux", Online: true, Services: []MeshService{
			{Name: "uptime-kuma", Domain: "uptime.mljr.eu", Description: "Uptime Monitoring"},
			{Name: "grafana", Domain: "grafana.mljr.eu", Description: "Monitoring"},
		}},
		{Name: "nas", TailscaleIP: "100.100.10.2", OS: "linux", Online: true, Relay: true},
		{Name: "homeassistant", TailscaleIP: "100.100.10.200", OS: "linux", Online: false, Services: []MeshService{
			{Name: "home-assistant", Domain: "home.mljr.eu", Description: "Home Automation"},
		}},
	}
}

func sampleAttackDays() []DayValue {
	seed := uint64(0xC0FFEE)
	next := func() uint64 {
		seed ^= seed << 13
		seed ^= seed >> 7
		seed ^= seed << 17
		return seed
	}
	now := time.Now()
	out := make([]DayValue, 0, 90)
	for i := 89; i >= 0; i-- {
		out = append(out, DayValue{
			Date:  now.AddDate(0, 0, -i),
			Count: int(next() % 60),
		})
	}
	return out
}

func sampleHourLabels(n int) []string {
	now := time.Now()
	out := make([]string, n)
	for i := range out {
		out[i] = now.Add(time.Duration(i-n+1) * time.Hour).Format("15h")
	}
	return out
}

// Options configures a Poller. Each data source is independently optional;
// the panel degrades gracefully when a source's fields are left empty.
type Options struct {
	KumaURL  string
	KumaSlug string
	PromURL  string

	TailscaleAPIKey  string // empty disables the mesh panel entirely
	TailscaleTailnet string

	InventoryURL string // public homelab-automation Ansible inventory (raw YAML)
	ServicesURL  string // public homelab-automation service registry (raw YAML)

	NtfyURL   string // ops alert target, e.g. failed Tailscale API auth
	NtfyTopic string
}

type Poller struct {
	kumaURL  string
	kumaSlug string
	promURL  string

	tsAPIKey     string
	tsTailnet    string
	inventoryURL string
	servicesURL  string
	ntfyURL      string
	ntfyTopic    string

	client *http.Client

	// Inventory/service-registry caches and mesh-alert throttle state are
	// only ever touched from the single poll() goroutine — no locking needed.
	invCache      map[string]string // inventory hostname -> tailscale_ip
	invFetchedAt  time.Time
	svcCache      []ServiceEntry
	svcFetchedAt  time.Time
	lastTSAlertAt time.Time

	mu   sync.RWMutex
	snap Snapshot
}

// New creates a poller. Empty KumaURL or PromURL disables that source;
// empty TailscaleAPIKey disables the mesh panel.
func New(opts Options) *Poller {
	if opts.KumaSlug == "" {
		opts.KumaSlug = "all"
	}
	if opts.TailscaleTailnet == "" {
		opts.TailscaleTailnet = "-"
	}
	return &Poller{
		kumaURL:      opts.KumaURL,
		kumaSlug:     opts.KumaSlug,
		promURL:      opts.PromURL,
		tsAPIKey:     opts.TailscaleAPIKey,
		tsTailnet:    opts.TailscaleTailnet,
		inventoryURL: opts.InventoryURL,
		servicesURL:  opts.ServicesURL,
		ntfyURL:      opts.NtfyURL,
		ntfyTopic:    opts.NtfyTopic,
		client:       &http.Client{Timeout: 10 * time.Second},
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
		BansCommunity:  -1,
		BansLocal:      -1,
	}

	if p.kumaURL != "" {
		if err := p.fetchKuma(ctx, &next); err != nil {
			log.Printf("homelab: kuma fetch failed: %v", err)
		}
	}
	if p.promURL != "" {
		p.fetchProm(ctx, &next)
	}
	if p.tsAPIKey != "" {
		p.fetchTailscale(ctx, &next)
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
	// Same for the mesh: a transient Tailscale API hiccup shouldn't blank
	// out the panel (the ntfy alert already covers persistent failures).
	if !next.MeshOK && prev.MeshOK {
		next.MeshOK = true
		next.Mesh = prev.Mesh
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
		{&s.ActiveBans, `sum(cs_active_decisions) or vector(0)`},
		{&s.Attacks24h, `sum(increase(cs_bucket_overflowed_total[24h])) or vector(0)`},
		{&s.SecurityEvents, `sum(cs_alerts) or vector(0)`},
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

	// Top threat categories by currently active decisions.
	if vec, err := p.promQueryVector(ctx, `topk(6, sum by (reason) (cs_active_decisions))`, "reason"); err != nil {
		log.Printf("homelab: promql top threats failed: %v", err)
	} else {
		s.TopThreats = vec
	}

	// Ban origin split: community blocklist (CAPI/lists) vs local detections.
	if vec, err := p.promQueryVector(ctx, `sum by (origin) (cs_active_decisions)`, "origin"); err != nil {
		if s.ActiveBans == 0 {
			s.BansCommunity, s.BansLocal = 0, 0
		} else {
			log.Printf("homelab: promql ban origins failed: %v", err)
		}
	} else {
		comm, local := 0, 0
		for _, nv := range vec {
			switch strings.ToLower(nv.Name) {
			case "capi", "lists":
				comm += nv.Value
			default:
				local += nv.Value
			}
		}
		s.BansCommunity, s.BansLocal = comm, local
	}

	now := time.Now()

	// Attacks blocked per day, last ~12 months (heatmap fills as data accrues).
	// VictoriaMetrics range queries spanning >40 days can return empty results
	// (per-day index cutoff), so fetch the year in 30-day chunks.
	seenDay := map[string]bool{}
	for chunk := 11; chunk >= 0; chunk-- {
		start := now.AddDate(0, 0, -30*(chunk+1))
		end := now.AddDate(0, 0, -30*chunk)
		pts, err := p.promQueryRange(ctx,
			`sum(increase(cs_bucket_overflowed_total[1d]))`,
			start, end, 24*time.Hour,
		)
		if err != nil {
			continue // empty chunks are expected until a year of data exists
		}
		for _, pt := range pts {
			day := time.Unix(int64(pt[0]), 0)
			key := day.Format("2006-01-02")
			if !seenDay[key] {
				seenDay[key] = true
				s.AttackDays = append(s.AttackDays, DayValue{Date: day, Count: int(pt[1])})
			}
		}
	}

	// Avg CPU across hosts, last 24h, 30-minute resolution.
	if pts, err := p.promQueryRange(ctx,
		`100 - avg(rate(node_cpu_seconds_total{mode="idle"}[10m])) * 100`,
		now.Add(-24*time.Hour), now, 30*time.Minute,
	); err != nil {
		log.Printf("homelab: promql cpu range failed: %v", err)
	} else {
		for _, pt := range pts {
			s.CPUUtil = append(s.CPUUtil, pt[1])
			label := ""
			t := time.Unix(int64(pt[0]), 0)
			if t.Minute() == 0 && t.Hour()%6 == 0 {
				label = t.Format("15h")
			}
			s.CPULabels = append(s.CPULabels, label)
		}
	}
}

// promQueryRange runs a PromQL range query and returns [timestamp, value]
// pairs of the first result series.
func (p *Poller) promQueryRange(ctx context.Context, query string, start, end time.Time, step time.Duration) ([][2]float64, error) {
	var res struct {
		Status string `json:"status"`
		Data   struct {
			Result []struct {
				Values [][2]any `json:"values"`
			} `json:"result"`
		} `json:"data"`
	}
	u := p.promURL + "/api/v1/query_range?query=" + url.QueryEscape(query) +
		fmt.Sprintf("&start=%d&end=%d&step=%d", start.Unix(), end.Unix(), int(step.Seconds()))
	if err := p.getJSON(ctx, u, &res); err != nil {
		return nil, err
	}
	if res.Status != "success" || len(res.Data.Result) == 0 {
		return nil, fmt.Errorf("no data")
	}
	out := make([][2]float64, 0, len(res.Data.Result[0].Values))
	for _, v := range res.Data.Result[0].Values {
		ts, ok1 := v[0].(float64)
		str, ok2 := v[1].(string)
		if !ok1 || !ok2 {
			continue
		}
		f, err := strconv.ParseFloat(str, 64)
		if err != nil {
			continue
		}
		out = append(out, [2]float64{ts, f})
	}
	return out, nil
}

// promQueryVector runs an instant PromQL query and returns one NameValue per
// result series, named by the given label, sorted by value descending.
func (p *Poller) promQueryVector(ctx context.Context, query, label string) ([]NameValue, error) {
	var res struct {
		Status string `json:"status"`
		Data   struct {
			Result []struct {
				Metric map[string]string `json:"metric"`
				Value  [2]any            `json:"value"`
			} `json:"result"`
		} `json:"data"`
	}
	u := p.promURL + "/api/v1/query?query=" + url.QueryEscape(query)
	if err := p.getJSON(ctx, u, &res); err != nil {
		return nil, err
	}
	if res.Status != "success" || len(res.Data.Result) == 0 {
		return nil, fmt.Errorf("no data")
	}
	out := make([]NameValue, 0, len(res.Data.Result))
	for _, r := range res.Data.Result {
		str, ok := r.Value[1].(string)
		if !ok {
			continue
		}
		f, err := strconv.ParseFloat(str, 64)
		if err != nil {
			continue
		}
		out = append(out, NameValue{Name: r.Metric[label], Value: int(f)})
	}
	sort.Slice(out, func(a, b int) bool { return out[a].Value > out[b].Value })
	return out, nil
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
	defer func() {
		if err := resp.Body.Close(); err != nil && !errors.Is(err, context.Canceled) {
			log.Printf("homelab: response body close failed: %v", err)
		}
	}()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status %d", resp.StatusCode)
	}
	return json.NewDecoder(resp.Body).Decode(dst)
}
