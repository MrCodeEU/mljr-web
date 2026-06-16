// Package data holds site content types and the loader that parses seed-cache.json.
// Re-run the scraper in the source repo and replace seed-cache.json to refresh.
package data

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

//go:embed seed-cache.json
var seedJSON []byte

type SiteData struct {
	GitHub        []Project      `json:"github_projects"`
	LinkedIn      LinkedInData   `json:"linkedin_data"`
	Strava        StravaData     `json:"strava_data"`
	GitHubStats   *GitHubStats   `json:"github_stats,omitempty"`
	SchemaVersion string         `json:"schema_version,omitempty"`
	GeneratedAt   string         `json:"generated_at,omitempty"`
	Content       map[string]SiteContent `json:"content"`
	Thesis        map[string][]Thesis    `json:"thesis"`
	Timeline      []TimelineItem `json:"timeline"`
}

// ContentFor returns the hand-authored copy for lang, falling back to
// English if the locale isn't present.
func (d SiteData) ContentFor(lang string) SiteContent {
	if c, ok := d.Content[lang]; ok {
		return c
	}
	return d.Content["en"]
}

// ThesisFor returns the thesis entries for lang, falling back to English if
// the locale isn't present.
func (d SiteData) ThesisFor(lang string) []Thesis {
	if t, ok := d.Thesis[lang]; ok {
		return t
	}
	return d.Thesis["en"]
}

// SiteContent holds hand-authored copy for sections that change
// independently of the GitHub/LinkedIn/Strava data feeds (mljr-data/content.json).
type SiteContent struct {
	Hero    HeroContent    `json:"hero"`
	Contact ContactContent `json:"contact"`
}

type HeroContent struct {
	StatusTag    string       `json:"status_tag"`
	TaglineLines []string     `json:"tagline_lines"`
	Description  string       `json:"description"`
	Bento        BentoContent `json:"bento"`
}

type BentoContent struct {
	Focus     BentoCell `json:"focus"`
	Status    BentoCell `json:"status"`
	Education BentoCell `json:"education"`
	Homelab   BentoCell `json:"homelab"`
}

type BentoCell struct {
	Label string `json:"label"`
	Value string `json:"value"`
	Sub   string `json:"sub"`
}

type ContactContent struct {
	Intro     string   `json:"intro"`
	Currently []string `json:"currently"`
}

// Thesis is a hand-authored education thesis entry with an optional PDF link.
type Thesis struct {
	Title       string `json:"title"`
	Type        string `json:"type"`
	Description string `json:"description"`
	PDF         string `json:"pdf"`
}

// TimelineItem is a hand-authored experience/education entry from timeline.json,
// embedded in site-data.json by the generator. Always in English.
type TimelineItem struct {
	ID           string   `json:"id"`
	Kind         string   `json:"kind"` // "work", "education", "thesis"
	Title        string   `json:"title"`
	TitleDE      string   `json:"title_de,omitempty"`
	Organization string   `json:"organization"`
	Start        string   `json:"start"` // "YYYY-MM" or "YYYY"
	End          *string  `json:"end"`   // null = present
	Location     string   `json:"location,omitempty"`
	Summary      string   `json:"summary"`
	SummaryDE    string   `json:"summary_de,omitempty"`
	Details      []string `json:"details,omitempty"`
	DetailsDE    []string `json:"details_de,omitempty"`
	Tags         []string `json:"tags,omitempty"`
	Logo         string   `json:"logo,omitempty"`
	Repo         string   `json:"repo,omitempty"`
}

// TitleFor returns the localized title, falling back to English.
func (t TimelineItem) TitleFor(lang string) string {
	if lang == "de" && t.TitleDE != "" {
		return t.TitleDE
	}
	return t.Title
}

// DetailsFor returns localized details, falling back to English.
func (t TimelineItem) DetailsFor(lang string) []string {
	if lang == "de" && len(t.DetailsDE) > 0 {
		return t.DetailsDE
	}
	return t.Details
}

// FormatPeriod returns a human-readable date range, e.g. "Nov 2025 – Present".
func (t TimelineItem) FormatPeriod() string {
	return formatYYYYMM(t.Start) + " – " + formatEndDate(t.End)
}

// FormatDuration returns an approximate duration string, e.g. "6 months", "3 years".
func (t TimelineItem) FormatDuration() string {
	start := parseYYYYMM(t.Start)
	var end time.Time
	if t.End == nil {
		end = time.Now()
	} else {
		end = parseYYYYMM(*t.End)
		// End month is inclusive — advance to end of month
		end = end.AddDate(0, 1, 0)
	}
	months := (end.Year()-start.Year())*12 + int(end.Month()) - int(start.Month())
	if months < 1 {
		months = 1
	}
	if months < 12 {
		if months == 1 {
			return "1 month"
		}
		return fmt.Sprintf("%d months", months)
	}
	years := months / 12
	rem := months % 12
	if rem == 0 {
		if years == 1 {
			return "1 year"
		}
		return fmt.Sprintf("%d years", years)
	}
	if years == 1 {
		return fmt.Sprintf("1 yr %d mo", rem)
	}
	return fmt.Sprintf("%d yrs %d mo", years, rem)
}

func parseYYYYMM(s string) time.Time {
	if len(s) == 7 {
		t, _ := time.Parse("2006-01", s)
		return t
	}
	t, _ := time.Parse("2006", s)
	return t
}

func formatYYYYMM(s string) string {
	t := parseYYYYMM(s)
	if t.IsZero() {
		return s
	}
	if len(s) == 7 {
		return t.Format("Jan 2006")
	}
	return t.Format("2006")
}

func formatEndDate(end *string) string {
	if end == nil {
		return "Present"
	}
	return formatYYYYMM(*end)
}

// WorkItems returns timeline items with kind == "work", in order.
func (d SiteData) WorkItems() []TimelineItem {
	var out []TimelineItem
	for _, item := range d.Timeline {
		if item.Kind == "work" {
			out = append(out, item)
		}
	}
	return out
}

// EduItems returns timeline items with kind == "education", in order.
func (d SiteData) EduItems() []TimelineItem {
	var out []TimelineItem
	for _, item := range d.Timeline {
		if item.Kind == "education" {
			out = append(out, item)
		}
	}
	return out
}

// GitHubStats holds contribution/heatmap data produced by the mljr-data
// generator. Nil only before the first data sync (e.g. fresh deploy), in
// which case pages fall back to sample data.
type GitHubStats struct {
	CommitsYear   int               `json:"commits_year"`
	LongestStreak int               `json:"longest_streak"`
	Contributions []ContributionDay `json:"contributions"`
	LanguageShare []LanguageShare   `json:"language_share"`
}

type ContributionDay struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

type LanguageShare struct {
	Name string  `json:"name"`
	Pct  float64 `json:"pct"`
}

type LinkedInData struct {
	Profile    Profile  `json:"profile"`
	Name       string   `json:"name"`
	Headline   string   `json:"headline"`
	Location   string   `json:"location"`
	About      string   `json:"about"`
	Experience []Job    `json:"experience"`
	Education  []School `json:"education"`
	Skills     []string `json:"skills"`
}

type Profile struct {
	Name     string `json:"name"`
	PhotoURL string `json:"photo_url"`
}

type Job struct {
	Title    string `json:"title"`
	Company  string `json:"company"`
	Type     string `json:"type"`
	Period   string `json:"period"`
	Duration string `json:"duration"`
	Desc     string `json:"description"`
}

type School struct {
	School string `json:"school"`
	Degree string `json:"degree"`
	Period string `json:"period"`
}

type Project struct {
	Name     string        `json:"name"`
	Desc     string        `json:"description"`
	URL      string        `json:"url"`
	Stars    int           `json:"stars"`
	Language string        `json:"language"`
	Topics   []string      `json:"topics"`
	Images   []string      `json:"images"`
	Featured bool          `json:"featured"`
	Order    int           `json:"order,omitempty"`
	Links    []ProjectLink `json:"links"`
}

type ProjectLink struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// StravaData is generated by the data pipeline and embedded into the homepage.
// Distances are stored in meters and durations in seconds to keep the data
// source precise; display helpers convert them at render time.
type StravaData struct {
	GeneratedAt       string              `json:"generated_at,omitempty"`
	YTDCalories       float64             `json:"ytd_calories,omitempty"`
	TotalStats        StravaStats         `json:"total_stats"`
	YearToDateStats   StravaStats         `json:"year_to_date_stats"`
	RecentActivities  []StravaActivity    `json:"recent_activities"`
	BestActivities    StravaBestRecords   `json:"best_activities"`
	PersonalRecords   []StravaRecord      `json:"personal_records"`
	Disciplines       []StravaDiscipline  `json:"disciplines"`
	MonthlyActivities []StravaMonthBucket `json:"monthly_activities,omitempty"`
	Year              string              `json:"year,omitempty"`
}

type StravaStats struct {
	Count         int     `json:"count"`
	Distance      float64 `json:"distance"`
	MovingTime    int     `json:"moving_time"`
	ElapsedTime   int     `json:"elapsed_time"`
	ElevationGain float64 `json:"elevation_gain"`
}

type StravaActivity struct {
	ID                 int64   `json:"id,omitempty"`
	Name               string  `json:"name"`
	Distance           float64 `json:"distance"`
	MovingTime         int     `json:"moving_time"`
	ElapsedTime        int     `json:"elapsed_time,omitempty"`
	TotalElevationGain float64 `json:"total_elevation_gain,omitempty"`
	Type               string  `json:"type"`
	StartDate          string  `json:"start_date,omitempty"`
	StartDateLocal     string  `json:"start_date_local,omitempty"`
	AveragePace        float64 `json:"average_pace,omitempty"`
	AverageSpeed       float64 `json:"average_speed,omitempty"`
	MaxSpeed           float64 `json:"max_speed,omitempty"`
	AverageHeartrate   float64 `json:"average_heartrate,omitempty"`
	MaxHeartrate       float64 `json:"max_heartrate,omitempty"`
	Calories           float64 `json:"calories,omitempty"`
}

type StravaBestRecords struct {
	LongestDistance StravaActivity `json:"longest_distance"`
	LongestTime     StravaActivity `json:"longest_time"`
	FastestPace     StravaActivity `json:"fastest_pace"`
	MostElevation   StravaActivity `json:"most_elevation"`
}

type StravaRecord struct {
	Type     string         `json:"type"`
	Time     int            `json:"time"`
	Distance float64        `json:"distance"`
	Date     string         `json:"date"`
	Activity StravaActivity `json:"activity"`
}

type StravaDiscipline struct {
	Type          string           `json:"type"`
	Label         string           `json:"label"`
	Count         int              `json:"count"`
	TotalTime     int              `json:"total_time"`
	TotalDistance float64          `json:"total_distance"`
	AvgHeartrate  float64          `json:"avg_heartrate"`
	Activities    []StravaActivity `json:"activities"`
}

type StravaMonthBucket struct {
	Month    string  `json:"month"`
	Count    int     `json:"count"`
	Distance float64 `json:"distance"`
	Time     int     `json:"time"`
}

// SkillGroup is a hand-curated grouping of the flat LinkedIn skills list plus
// extra skills not exported by LinkedIn.
type SkillGroup struct {
	Label  string
	Short  string // short label for chart axes
	Tone   string
	Icon   string // Iconify icon for the group
	Level  int    // self-assessed depth 0–100, drives the radar chart
	Skills []string
}

var skillGroups = []SkillGroup{
	{"Languages", "Languages", "violet", "lucide:code-2", 85, []string{"Go", "Rust", "TypeScript", "JavaScript", "Python", "Kotlin", "Java", "C#"}},
	{"Web", "Web", "sky", "lucide:globe", 75, []string{"HTML", "CSS", "Svelte", "SvelteKit", "AngularJS", "Angular", "ASP.NET", ".NET"}},
	{"Infra / Homelab", "Infra", "lime", "lucide:server", 82, []string{"Docker", "Ansible", "Linux", "Unraid", "Tailscale", "Caddy", "CI/CD"}},
	{"Security", "Security", "yellow", "lucide:shield", 78, []string{"Cybersecurity", "Network Security", "IAM", "Prolog", "PAM"}},
	{"Embedded / BLE", "Embedded", "pink", "lucide:cpu", 60, []string{"ESP32", "BLE", "Dexcom G7", "Nightscout", "Kotlin Multiplatform"}},
	{"Ops / Data", "Ops/Data", "mint", "lucide:database", 70, []string{"SQLite", "VictoriaMetrics", "Grafana", "slog", "WebGL"}},
}

// Load parses and returns the embedded fallback site data. Panics on bad JSON
// because the embedded seed is committed with the binary.
func Load() SiteData {
	d, err := parse(seedJSON)
	if err != nil {
		log.Fatalf("data: parse seed-cache.json: %v", err)
	}
	return d
}

func LoadFile(path string) (SiteData, error) {
	// #nosec G304 -- HOMEPAGE_DATA_FILE is an operator-controlled data source.
	b, err := os.ReadFile(path)
	if err != nil {
		return SiteData{}, err
	}
	return parse(b)
}

func parse(b []byte) (SiteData, error) {
	var d SiteData
	if err := json.Unmarshal(b, &d); err != nil {
		return SiteData{}, err
	}
	d.Strava.Normalize()
	return d, nil
}

type Store struct {
	path        string
	reloadEvery time.Duration

	mu        sync.RWMutex
	current   SiteData
	lastCheck time.Time
	lastMod   time.Time
}

func NewStore(path, reloadSeconds string) *Store {
	reloadEvery := 5 * time.Minute
	if seconds, err := strconv.Atoi(strings.TrimSpace(reloadSeconds)); err == nil && seconds > 0 {
		reloadEvery = time.Duration(seconds) * time.Second
	}

	s := &Store{
		path:        strings.TrimSpace(path),
		reloadEvery: reloadEvery,
		current:     Load(),
	}
	if s.path != "" {
		if err := s.reloadIfChanged(true); err != nil {
			log.Printf("data: using embedded fallback; %s: %v", s.path, err)
		}
	}
	return s
}

func (s *Store) Current() SiteData {
	if s == nil {
		return Load()
	}
	s.mu.RLock()
	shouldCheck := s.path != "" && time.Since(s.lastCheck) >= s.reloadEvery
	s.mu.RUnlock()
	if shouldCheck {
		if err := s.reloadIfChanged(false); err != nil {
			log.Printf("data: keeping previous data; %s: %v", s.path, err)
		}
	}

	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.current
}

func (s *Store) reloadIfChanged(force bool) error {
	info, err := os.Stat(s.path)
	s.mu.Lock()
	defer s.mu.Unlock()
	s.lastCheck = time.Now()
	if err != nil {
		return err
	}
	if !force && !info.ModTime().After(s.lastMod) {
		return nil
	}
	d, err := LoadFile(s.path)
	if err != nil {
		return err
	}
	s.current = d
	s.lastMod = info.ModTime()
	log.Printf("data: loaded %s", s.path)
	return nil
}

// FeaturedProjects returns projects marked featured, excluding meta-entries,
// ordered by the curated "order" field (ascending, lower = shown first).
func (d SiteData) FeaturedProjects() []Project {
	var out []Project
	for _, p := range d.GitHub {
		if p.Featured && !strings.EqualFold(p.Name, "homepage") {
			out = append(out, p)
		}
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Order != out[j].Order {
			return out[i].Order < out[j].Order
		}
		return out[i].Name < out[j].Name
	})
	return out
}

// AllProjects returns all non-featured, non-meta projects, ordered by the
// curated "order" field (ascending, lower = shown first), then by name.
func (d SiteData) AllProjects() []Project {
	var out []Project
	for _, p := range d.GitHub {
		if !p.Featured && !strings.EqualFold(p.Name, "homepage") {
			out = append(out, p)
		}
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Order != out[j].Order {
			return out[i].Order < out[j].Order
		}
		return out[i].Name < out[j].Name
	})
	return out
}

// HasStrava returns true when the generated data contains enough public
// activity aggregate data to render a useful section.
func (d SiteData) HasStrava() bool {
	return d.Strava.TotalStats.Count > 0 || len(d.Strava.RecentActivities) > 0
}

// SkillGroups returns the curated skill groups.
func SkillGroups() []SkillGroup { return skillGroups }

// DisplayName returns the short first name.
func (d LinkedInData) DisplayName() string {
	parts := strings.Fields(d.Name)
	if len(parts) > 0 {
		return parts[0]
	}
	return d.Name
}

// RelevantExperience returns the top N non-trivial experience items.
func (d LinkedInData) RelevantExperience(n int) []Job {
	var out []Job
	for _, j := range d.Experience {
		out = append(out, j)
		if len(out) >= n {
			break
		}
	}
	return out
}

// LocalImages returns displayable images for a project: site-relative paths
// and absolute http(s) URLs (e.g. placeholder images), excluding bare repo
// paths that haven't been resolved to a URL.
func (p Project) LocalImages() []string {
	var out []string
	for _, img := range p.Images {
		if strings.HasPrefix(img, "/") || strings.HasPrefix(img, "http://") || strings.HasPrefix(img, "https://") {
			out = append(out, img)
		}
	}
	return out
}

func (s *StravaData) Normalize() {
	for i := range s.RecentActivities {
		s.RecentActivities[i].Normalize()
	}
	for i := range s.Disciplines {
		for j := range s.Disciplines[i].Activities {
			s.Disciplines[i].Activities[j].Normalize()
		}
	}
	s.BestActivities.LongestDistance.Normalize()
	s.BestActivities.LongestTime.Normalize()
	s.BestActivities.FastestPace.Normalize()
	s.BestActivities.MostElevation.Normalize()
	for i := range s.PersonalRecords {
		s.PersonalRecords[i].Activity.Normalize()
		if s.PersonalRecords[i].Date == "" {
			s.PersonalRecords[i].Date = s.PersonalRecords[i].Activity.DisplayDate()
		}
	}
	if s.Year == "" {
		s.Year = time.Now().Format("2006")
	}
}

func (a *StravaActivity) Normalize() {
	if a.StartDate == "" {
		a.StartDate = a.StartDateLocal
	}
	if a.ElapsedTime == 0 {
		a.ElapsedTime = a.MovingTime
	}
	if a.AveragePace == 0 && a.AverageSpeed > 0 {
		a.AveragePace = 1000 / (a.AverageSpeed * 60)
	}
}

func (a StravaActivity) DisplayDate() string {
	raw := a.StartDate
	if raw == "" {
		raw = a.StartDateLocal
	}
	if raw == "" {
		return ""
	}
	for _, layout := range []string{time.RFC3339, "2006-01-02T15:04:05Z07:00", "2006-01-02"} {
		t, err := time.Parse(layout, raw)
		if err == nil {
			return t.Format("Jan 2006")
		}
	}
	if len(raw) >= 7 {
		return raw[:7]
	}
	return raw
}

func DistanceKM(meters float64) float64 {
	return math.Round(meters/100) / 10
}

func DurationHours(seconds int) float64 {
	return math.Round(float64(seconds)/36) / 100
}

// DurationClock renders seconds as "37:26" or "2:51:00" — exact, not decimal.
func DurationClock(seconds int) string {
	hours := seconds / 3600
	mins := (seconds % 3600) / 60
	secs := seconds % 60
	if hours > 0 {
		return fmt.Sprintf("%d:%02d:%02d", hours, mins, secs)
	}
	return fmt.Sprintf("%d:%02d", mins, secs)
}

// DurationHM renders seconds as "27 h 39 m" for aggregate tiles.
func DurationHM(seconds int) string {
	hours := seconds / 3600
	mins := (seconds % 3600) / 60
	if hours == 0 {
		return fmt.Sprintf("%d m", mins)
	}
	return fmt.Sprintf("%d h %02d m", hours, mins)
}

// AvgHeartrate returns the session-weighted average heartrate across
// disciplines, or 0 if no discipline reports one.
func (s StravaData) AvgHeartrate() float64 {
	var sum float64
	var n int
	for _, d := range s.Disciplines {
		if d.AvgHeartrate > 0 && d.Count > 0 {
			sum += d.AvgHeartrate * float64(d.Count)
			n += d.Count
		}
	}
	if n == 0 {
		return 0
	}
	return sum / float64(n)
}

func PaceLabel(minutesPerKM float64) string {
	if minutesPerKM <= 0 {
		return "n/a"
	}
	minutes := int(minutesPerKM)
	seconds := int(math.Round((minutesPerKM - float64(minutes)) * 60))
	if seconds == 60 {
		minutes++
		seconds = 0
	}
	return fmt.Sprintf("%d:%02d /km", minutes, seconds)
}

// LogoForCompany maps LinkedIn company names to local logo paths.
func LogoForCompany(company string) string {
	m := map[string]string{
		"Dynatrace":                        "/static/logos/dynatrace.jpg",
		"Johannes Kepler Universität Linz": "/static/logos/jku.jpg",
		"ventopay gmbh":                    "/static/logos/ventopay.jpg",
		"Bosch":                            "/static/logos/bosch.jpg",
		"Bosch Rexroth":                    "/static/logos/bosch-rexroth.jpg",
		"ENGEL":                            "/static/logos/engel.jpg",
		"HerzReha Bad Ischl":               "/static/logos/herzreha-bad-ischl.jpg",
	}
	return m[company]
}

// LogoForSchool maps school names to local logo paths.
func LogoForSchool(school string) string {
	m := map[string]string{
		"Johannes Kepler Universität Linz": "/static/logos/jku.jpg",
		"HTL Steyr":                        "/static/logos/htl-steyr.png",
	}
	return m[school]
}

// TechIcon maps a lowercase tech/topic keyword to an Iconify icon name.
// Returns "" if no icon is known for the keyword.
func TechIcon(keyword string) string {
	m := map[string]string{
		// ── languages ────────────────────────────────────────────────
		"go":         "simple-icons:go",
		"golang":     "simple-icons:go",
		"rust":       "simple-icons:rust",
		"typescript": "simple-icons:typescript",
		"javascript": "simple-icons:javascript",
		"python":     "simple-icons:python",
		"kotlin":     "simple-icons:kotlin",
		"java":       "lucide:code-2",
		"c#":         "simple-icons:csharp",
		"csharp":     "simple-icons:csharp",
		"nodejs":     "simple-icons:nodedotjs",
		"node.js":    "simple-icons:nodedotjs",
		// ── web frameworks ────────────────────────────────────────────
		"svelte":         "simple-icons:svelte",
		"sveltekit":      "simple-icons:svelte",
		"tailwindcss":    "simple-icons:tailwindcss",
		"html":           "simple-icons:html5",
		"css":            "simple-icons:css3",
		"angular":        "simple-icons:angular",
		"angularjs":      "simple-icons:angularjs",
		"asp.net":        "simple-icons:dotnet",
		".net":           "simple-icons:dotnet",
		".net-framework": "simple-icons:dotnet",
		"dotnet":         "simple-icons:dotnet",
		// ── infra / homelab ──────────────────────────────────────────
		"docker":     "simple-icons:docker",
		"linux":      "simple-icons:linux",
		"ansible":    "simple-icons:ansible",
		"tailscale":  "simple-icons:tailscale",
		"devops":     "lucide:workflow",
		"cicd":       "lucide:git-branch",
		"ci/cd":      "lucide:git-branch",
		"automation": "lucide:bot",
		"homelab":    "lucide:server",
		"unraid":     "lucide:hard-drive",
		"caddy":      "lucide:shield",
		"template":   "lucide:layers",
		// ── web general ──────────────────────────────────────────────
		"web":           "lucide:globe",
		"webgl":         "lucide:box",
		"3d":            "lucide:box",
		"animation":     "lucide:zap",
		"multiplatform": "lucide:layers",
		"plugin":        "lucide:puzzle",
		// ── data / ops ───────────────────────────────────────────────
		"database":          "lucide:database",
		"sqlite":            "lucide:database",
		"monitoring":        "lucide:activity",
		"victoriametrics":   "simple-icons:victoriametrics",
		"grafana":           "simple-icons:grafana",
		"slog":              "lucide:scroll",
		"networking":        "lucide:network",
		"power-consumption": "lucide:zap",
		// ── security ─────────────────────────────────────────────────
		"security":           "lucide:shield",
		"cybersecurity":      "lucide:shield",
		"network security":   "lucide:network",
		"netzwerksicherheit": "lucide:network",
		"pam":                "lucide:key",
		"iam":                "lucide:user-check",
		"prolog":             "lucide:brain",
		// ── embedded / BLE ───────────────────────────────────────────
		"esp32":                "lucide:cpu",
		"ble":                  "lucide:bluetooth",
		"dexcom g7":            "lucide:activity",
		"nightscout":           "lucide:activity",
		"kotlin multiplatform": "simple-icons:kotlin",
		"diabetes":             "lucide:activity",
		"healthcare":           "lucide:activity",
		// ── research / ai ────────────────────────────────────────────
		"ai":               "lucide:brain",
		"machine-learning": "lucide:brain",
		"lstm":             "lucide:brain",
		"deep-learning":    "lucide:brain",
		"research":         "lucide:microscope",
		"thesis":           "lucide:microscope",
		// ── media ────────────────────────────────────────────────────
		"music": "lucide:music",
		"audio": "lucide:music",
		// ── project tools ────────────────────────────────────────────
		"flexbar":    "lucide:cpu",
		"streamdeck": "lucide:cpu",
		"desktop":    "lucide:cpu",
		"mobile":     "lucide:cpu",
		"wails":      "simple-icons:go",
		"github":     "simple-icons:github",
		// ── linkedin exported names (German) ─────────────────────────
		"it-infrastruktur":            "lucide:server",
		"it-management":               "lucide:user-check",
		"softwareentwicklung":         "lucide:code-2",
		"kontinuierliche integration": "lucide:git-branch",
		"netzwerkadministration":      "lucide:network",
	}
	return m[strings.ToLower(keyword)]
}
