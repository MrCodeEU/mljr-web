// Package registry stores showcase metadata for the catalogue app. Components
// register themselves via init() in sibling *_showcase.go files compiled with
// -tags showcase. Prod binaries omit the tag, so this registry is tree-shaken.
package registry

import (
	"sort"
	"strings"
	"sync"

	g "maragu.dev/gomponents"
)

type ControlType string

const (
	ControlEnum ControlType = "enum"
	ControlBool ControlType = "bool"
	ControlText ControlType = "text"
)

type Control struct {
	Name    string
	Type    ControlType
	Options []string
	Default string
}

type Example struct {
	Title string
	Node  func() g.Node
	Code  string
}

type Component struct {
	Slug          string
	Name          string
	Category      string
	Summary       string
	Code          string // Go usage example shown on detail page
	PreviewHeight string // override iframe height (default "480px"), e.g. "720px"
	Controls      []Control
	Render        func(p map[string]string) g.Node
	Examples      []Example
}

var (
	mu  sync.RWMutex
	reg = map[string]*Component{}
)

func Register(c *Component) {
	mu.Lock()
	defer mu.Unlock()
	reg[c.Slug] = c
}

func Get(slug string) (*Component, bool) {
	mu.RLock()
	defer mu.RUnlock()
	c, ok := reg[slug]
	return c, ok
}

func All() []*Component {
	mu.RLock()
	defer mu.RUnlock()
	out := make([]*Component, 0, len(reg))
	for _, c := range reg {
		out = append(out, c)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Category == out[j].Category {
			return out[i].Name < out[j].Name
		}
		return out[i].Category < out[j].Category
	})
	return out
}

func Categories() []string {
	seen := map[string]bool{}
	for _, c := range All() {
		seen[c.Category] = true
	}
	out := make([]string, 0, len(seen))
	for k := range seen {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// DefaultProps returns each control's default value keyed by name.
func DefaultProps(c *Component) map[string]string {
	p := map[string]string{}
	for _, ctl := range c.Controls {
		p[ctl.Name] = ctl.Default
	}
	return p
}

// Combinations returns the cartesian product of all enum controls.
// Bool and text controls are fixed at their defaults in every combination.
func Combinations(c *Component) []map[string]string {
	defaults := DefaultProps(c)
	var enumCtls []Control
	for _, ctl := range c.Controls {
		if ctl.Type == ControlEnum {
			enumCtls = append(enumCtls, ctl)
		}
	}
	if len(enumCtls) == 0 {
		return []map[string]string{defaults}
	}
	result := []map[string]string{{}}
	for _, ctl := range enumCtls {
		var next []map[string]string
		for _, existing := range result {
			for _, opt := range ctl.Options {
				combo := make(map[string]string, len(defaults))
				for k, v := range defaults {
					combo[k] = v
				}
				for k, v := range existing {
					combo[k] = v
				}
				combo[ctl.Name] = opt
				next = append(next, combo)
			}
		}
		result = next
	}
	return result
}

// ComboLabel returns a human-readable label for a combination (enum values joined by ·).
func ComboLabel(c *Component, combo map[string]string) string {
	var parts []string
	for _, ctl := range c.Controls {
		if ctl.Type == ControlEnum {
			if v := combo[ctl.Name]; v != "" {
				parts = append(parts, v)
			}
		}
	}
	if len(parts) == 0 {
		return "default"
	}
	return strings.Join(parts, " · ")
}
