package registry

import (
	"sort"
	"sync"

	g "maragu.dev/gomponents"
)

// Pattern is a full-page UI composition demonstrating multiple components working together.
type Pattern struct {
	Slug        string
	Name        string
	Category    string // "marketing" | "app" | "auth" | "content"
	Description string
	Tags        []string
	// Render returns a complete HTML page node (should include PageShell).
	Render func(theme, mode string) g.Node
}

var (
	pmu      sync.RWMutex
	patterns = map[string]*Pattern{}
)

func RegisterPattern(p *Pattern) {
	pmu.Lock()
	defer pmu.Unlock()
	patterns[p.Slug] = p
}

func GetPattern(slug string) (*Pattern, bool) {
	pmu.RLock()
	defer pmu.RUnlock()
	p, ok := patterns[slug]
	return p, ok
}

func AllPatterns() []*Pattern {
	pmu.RLock()
	defer pmu.RUnlock()
	out := make([]*Pattern, 0, len(patterns))
	for _, p := range patterns {
		out = append(out, p)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Category == out[j].Category {
			return out[i].Name < out[j].Name
		}
		return out[i].Category < out[j].Category
	})
	return out
}

func PatternCategories() []string {
	seen := map[string]bool{}
	for _, p := range AllPatterns() {
		seen[p.Category] = true
	}
	out := make([]string, 0, len(seen))
	for k := range seen {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
