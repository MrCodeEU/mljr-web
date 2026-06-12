package scrapers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	hpdata "mljr-web/projects/homepage/data"
)

const (
	defaultStravaAPIBase  = "https://www.strava.com/api/v3"
	defaultStravaTokenURL = "https://www.strava.com/oauth/token" //nolint:gosec // not a credential, just an OAuth endpoint URL
)

type StravaConfig struct {
	ClientID     string
	ClientSecret string
	RefreshToken string
	After        time.Time
	MaxPages     int
	HTTPClient   *http.Client
	APIBase      string
	TokenURL     string
	Now          func() time.Time
}

type StravaScraper struct {
	cfg         StravaConfig
	client      *http.Client
	accessToken string
	tokenExpiry time.Time
}

func NewStravaScraper(cfg StravaConfig) *StravaScraper {
	client := cfg.HTTPClient
	if client == nil {
		client = &http.Client{Timeout: 30 * time.Second}
	}
	if cfg.APIBase == "" {
		cfg.APIBase = defaultStravaAPIBase
	}
	if cfg.TokenURL == "" {
		cfg.TokenURL = defaultStravaTokenURL
	}
	if cfg.MaxPages <= 0 {
		cfg.MaxPages = 5
	}
	if cfg.Now == nil {
		cfg.Now = time.Now
	}
	return &StravaScraper{cfg: cfg, client: client}
}

func (s *StravaScraper) Scrape(ctx context.Context) (hpdata.StravaData, error) {
	if strings.TrimSpace(s.cfg.ClientID) == "" || strings.TrimSpace(s.cfg.ClientSecret) == "" || strings.TrimSpace(s.cfg.RefreshToken) == "" {
		return hpdata.StravaData{}, fmt.Errorf("missing Strava credentials")
	}
	if err := s.ensureAccessToken(ctx); err != nil {
		return hpdata.StravaData{}, err
	}

	athleteID, err := s.fetchAthleteID(ctx)
	if err != nil {
		return hpdata.StravaData{}, err
	}
	stats, err := s.fetchStats(ctx, athleteID)
	if err != nil {
		return hpdata.StravaData{}, err
	}
	rawActivities, err := s.fetchActivities(ctx)
	if err != nil {
		return hpdata.StravaData{}, err
	}

	data := buildStravaData(stats, rawActivities, s.cfg.Now())
	data.Normalize()
	return data, nil
}

type stravaTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
}

type stravaActivityDTO struct {
	ID                 int64   `json:"id"`
	Name               string  `json:"name"`
	Distance           float64 `json:"distance"`
	MovingTime         int     `json:"moving_time"`
	ElapsedTime        int     `json:"elapsed_time"`
	TotalElevationGain float64 `json:"total_elevation_gain"`
	Type               string  `json:"type"`
	StartDate          string  `json:"start_date"`
	StartDateLocal     string  `json:"start_date_local"`
	AverageSpeed       float64 `json:"average_speed"`
	MaxSpeed           float64 `json:"max_speed"`
	AverageHeartrate   float64 `json:"average_heartrate"`
	MaxHeartrate       float64 `json:"max_heartrate"`
	Calories           float64 `json:"calories"`
	Kilojoules         float64 `json:"kilojoules"`
}

type stravaStatsDTO struct {
	AllRunTotals stravaStatsTotals `json:"all_run_totals"`
	YTDRunTotals stravaStatsTotals `json:"ytd_run_totals"`
}

type stravaStatsTotals struct {
	Count         int     `json:"count"`
	Distance      float64 `json:"distance"`
	MovingTime    int     `json:"moving_time"`
	ElapsedTime   int     `json:"elapsed_time"`
	ElevationGain float64 `json:"elevation_gain"`
}

func (s *StravaScraper) ensureAccessToken(ctx context.Context) error {
	if s.accessToken != "" && s.cfg.Now().Before(s.tokenExpiry.Add(-time.Minute)) {
		return nil
	}

	form := url.Values{}
	form.Set("client_id", s.cfg.ClientID)
	form.Set("client_secret", s.cfg.ClientSecret)
	form.Set("grant_type", "refresh_token")
	form.Set("refresh_token", s.cfg.RefreshToken)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.cfg.TokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	var token stravaTokenResponse
	if err := s.doJSON(req, &token); err != nil {
		return fmt.Errorf("refresh Strava token: %w", err)
	}
	if token.AccessToken == "" {
		return fmt.Errorf("refresh Strava token: missing access token")
	}
	s.accessToken = token.AccessToken
	s.tokenExpiry = time.Unix(token.ExpiresAt, 0)
	if token.RefreshToken != "" {
		s.cfg.RefreshToken = token.RefreshToken
	}
	return nil
}

func (s *StravaScraper) fetchAthleteID(ctx context.Context) (int64, error) {
	var athlete struct {
		ID int64 `json:"id"`
	}
	req, err := s.authRequest(ctx, http.MethodGet, s.cfg.APIBase+"/athlete")
	if err != nil {
		return 0, err
	}
	if err := s.doJSON(req, &athlete); err != nil {
		return 0, fmt.Errorf("fetch Strava athlete: %w", err)
	}
	if athlete.ID == 0 {
		return 0, fmt.Errorf("fetch Strava athlete: missing athlete id")
	}
	return athlete.ID, nil
}

func (s *StravaScraper) fetchStats(ctx context.Context, athleteID int64) (stravaStatsDTO, error) {
	var stats stravaStatsDTO
	req, err := s.authRequest(ctx, http.MethodGet, fmt.Sprintf("%s/athletes/%d/stats", s.cfg.APIBase, athleteID))
	if err != nil {
		return stats, err
	}
	if err := s.doJSON(req, &stats); err != nil {
		return stats, fmt.Errorf("fetch Strava stats: %w", err)
	}
	return stats, nil
}

func (s *StravaScraper) fetchActivities(ctx context.Context) ([]stravaActivityDTO, error) {
	var out []stravaActivityDTO
	after := ""
	if !s.cfg.After.IsZero() {
		after = fmt.Sprintf("&after=%d", s.cfg.After.Unix())
	}
	for page := 1; page <= s.cfg.MaxPages; page++ {
		endpoint := fmt.Sprintf("%s/athlete/activities?per_page=200&page=%d%s", s.cfg.APIBase, page, after)
		req, err := s.authRequest(ctx, http.MethodGet, endpoint)
		if err != nil {
			return nil, err
		}
		var batch []stravaActivityDTO
		if err := s.doJSON(req, &batch); err != nil {
			return nil, fmt.Errorf("fetch Strava activities page %d: %w", page, err)
		}
		if len(batch) == 0 {
			break
		}
		out = append(out, batch...)
		if len(batch) < 200 {
			break
		}
	}
	return out, nil
}

func (s *StravaScraper) authRequest(ctx context.Context, method, endpoint string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+s.accessToken)
	return req, nil
}

func (s *StravaScraper) doJSON(req *http.Request, target any) error {
	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return fmt.Errorf("status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}
	return json.NewDecoder(resp.Body).Decode(target)
}

func buildStravaData(stats stravaStatsDTO, raw []stravaActivityDTO, now time.Time) hpdata.StravaData {
	activities := make([]hpdata.StravaActivity, 0, len(raw))
	for _, a := range raw {
		activities = append(activities, convertStravaActivity(a))
	}
	sort.Slice(activities, func(i, j int) bool {
		return activityTime(activities[i]).After(activityTime(activities[j]))
	})

	runs := filterActivities(activities, "running")
	recent := activities
	if len(recent) > 10 {
		recent = recent[:10]
	}
	return hpdata.StravaData{
		GeneratedAt:       now.UTC().Format(time.RFC3339),
		Year:              now.Format("2006"),
		TotalStats:        convertStats(stats.AllRunTotals),
		YearToDateStats:   convertStats(stats.YTDRunTotals),
		RecentActivities:  recent,
		BestActivities:    bestActivities(runs),
		PersonalRecords:   personalRecords(runs),
		Disciplines:       disciplines(activities),
		MonthlyActivities: monthBuckets(activities),
	}
}

func convertStats(t stravaStatsTotals) hpdata.StravaStats {
	return hpdata.StravaStats{
		Count:         t.Count,
		Distance:      t.Distance,
		MovingTime:    t.MovingTime,
		ElapsedTime:   t.ElapsedTime,
		ElevationGain: t.ElevationGain,
	}
}

func convertStravaActivity(a stravaActivityDTO) hpdata.StravaActivity {
	calories := a.Calories
	if calories == 0 && a.Kilojoules > 0 {
		calories = a.Kilojoules * 0.239
	}
	out := hpdata.StravaActivity{
		ID:                 a.ID,
		Name:               a.Name,
		Distance:           a.Distance,
		MovingTime:         a.MovingTime,
		ElapsedTime:        a.ElapsedTime,
		TotalElevationGain: a.TotalElevationGain,
		Type:               a.Type,
		StartDate:          firstNonEmpty(a.StartDateLocal, a.StartDate),
		StartDateLocal:     a.StartDateLocal,
		AverageSpeed:       a.AverageSpeed,
		MaxSpeed:           a.MaxSpeed,
		AverageHeartrate:   a.AverageHeartrate,
		MaxHeartrate:       a.MaxHeartrate,
		Calories:           calories,
	}
	out.Normalize()
	return out
}

func disciplineType(kind string) string {
	switch strings.ToLower(kind) {
	case "run", "trailrun", "virtualrun":
		return "running"
	case "ride", "virtualride", "mountainbikeride", "gravelride", "ebikeride", "emountainbikeride":
		return "cycling"
	default:
		return "training"
	}
}

func filterActivities(activities []hpdata.StravaActivity, discipline string) []hpdata.StravaActivity {
	out := make([]hpdata.StravaActivity, 0)
	for _, a := range activities {
		if disciplineType(a.Type) == discipline {
			out = append(out, a)
		}
	}
	return out
}

func bestActivities(activities []hpdata.StravaActivity) hpdata.StravaBestRecords {
	if len(activities) == 0 {
		return hpdata.StravaBestRecords{}
	}
	best := hpdata.StravaBestRecords{
		LongestDistance: activities[0],
		LongestTime:     activities[0],
		FastestPace:     activities[0],
		MostElevation:   activities[0],
	}
	for _, a := range activities {
		if a.Distance > best.LongestDistance.Distance {
			best.LongestDistance = a
		}
		if a.MovingTime > best.LongestTime.MovingTime {
			best.LongestTime = a
		}
		if a.AveragePace > 0 && (best.FastestPace.AveragePace == 0 || a.AveragePace < best.FastestPace.AveragePace) {
			best.FastestPace = a
		}
		if a.TotalElevationGain > best.MostElevation.TotalElevationGain {
			best.MostElevation = a
		}
	}
	return best
}

func personalRecords(activities []hpdata.StravaActivity) []hpdata.StravaRecord {
	targets := []struct {
		name string
		m    float64
	}{
		{"5k", 5000},
		{"10k", 10000},
		{"half_marathon", 21097.5},
		{"marathon", 42195},
	}
	records := make([]hpdata.StravaRecord, 0, len(targets))
	for _, target := range targets {
		var best hpdata.StravaActivity
		for _, a := range activities {
			tolerance := target.m * 0.02
			if a.Distance < target.m-tolerance || a.Distance > target.m+tolerance {
				continue
			}
			if best.MovingTime == 0 || a.MovingTime < best.MovingTime {
				best = a
			}
		}
		if best.MovingTime > 0 {
			records = append(records, hpdata.StravaRecord{
				Type:     target.name,
				Time:     best.MovingTime,
				Distance: best.Distance,
				Date:     best.DisplayDate(),
				Activity: best,
			})
		}
	}
	return records
}

func disciplines(activities []hpdata.StravaActivity) []hpdata.StravaDiscipline {
	type bucket struct {
		label    string
		items    []hpdata.StravaActivity
		time     int
		distance float64
		hrTotal  float64
		hrCount  int
	}
	buckets := map[string]*bucket{
		"running":  {label: "Running"},
		"cycling":  {label: "Cycling"},
		"training": {label: "Training"},
	}
	for _, a := range activities {
		key := disciplineType(a.Type)
		b := buckets[key]
		b.items = append(b.items, a)
		b.time += a.MovingTime
		b.distance += a.Distance
		if a.AverageHeartrate > 0 {
			b.hrTotal += a.AverageHeartrate
			b.hrCount++
		}
	}
	order := []string{"running", "cycling", "training"}
	out := make([]hpdata.StravaDiscipline, 0, len(order))
	for _, key := range order {
		b := buckets[key]
		if len(b.items) == 0 {
			continue
		}
		items := b.items
		if len(items) > 5 {
			items = items[:5]
		}
		avgHR := 0.0
		if b.hrCount > 0 {
			avgHR = b.hrTotal / float64(b.hrCount)
		}
		out = append(out, hpdata.StravaDiscipline{
			Type:          key,
			Label:         b.label,
			Count:         len(b.items),
			TotalTime:     b.time,
			TotalDistance: b.distance,
			AvgHeartrate:  avgHR,
			Activities:    items,
		})
	}
	return out
}

func monthBuckets(activities []hpdata.StravaActivity) []hpdata.StravaMonthBucket {
	byMonth := map[string]*hpdata.StravaMonthBucket{}
	for _, a := range activities {
		t := activityTime(a)
		if t.IsZero() {
			continue
		}
		key := t.Format("2006-01")
		b := byMonth[key]
		if b == nil {
			b = &hpdata.StravaMonthBucket{Month: key}
			byMonth[key] = b
		}
		b.Count++
		b.Distance += a.Distance
		b.Time += a.MovingTime
	}
	keys := make([]string, 0, len(byMonth))
	for key := range byMonth {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	out := make([]hpdata.StravaMonthBucket, 0, len(keys))
	for _, key := range keys {
		out = append(out, *byMonth[key])
	}
	return out
}

func activityTime(a hpdata.StravaActivity) time.Time {
	raw := firstNonEmpty(a.StartDate, a.StartDateLocal)
	t, _ := time.Parse(time.RFC3339, raw)
	return t
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}
