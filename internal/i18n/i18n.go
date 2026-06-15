// Package i18n provides simple JSON-catalog based translations for the
// homepage UI. Catalogs are flat maps of dot-path keys to translated strings,
// embedded at build time from locales/en.json and locales/de.json.
package i18n

import (
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"slices"
)

//go:embed locales/*.json
var localesFS embed.FS

// DefaultLang is used as a fallback when the requested locale or key is
// missing.
const DefaultLang = "en"

// Supported lists the locales the site can render in.
var Supported = []string{"en", "de"}

var catalogs map[string]map[string]string

func init() {
	catalogs = make(map[string]map[string]string, len(Supported))
	for _, lang := range Supported {
		b, err := localesFS.ReadFile("locales/" + lang + ".json")
		if err != nil {
			log.Fatalf("i18n: read locale %q: %v", lang, err)
		}
		var nested map[string]any
		if err := json.Unmarshal(b, &nested); err != nil {
			log.Fatalf("i18n: parse locale %q: %v", lang, err)
		}
		flat := make(map[string]string)
		flatten("", nested, flat)
		catalogs[lang] = flat
	}
}

// flatten turns nested JSON objects into dot-path keys, e.g.
// {"nav":{"projects":"Projects"}} -> {"nav.projects": "Projects"}.
func flatten(prefix string, in map[string]any, out map[string]string) {
	for k, v := range in {
		key := k
		if prefix != "" {
			key = prefix + "." + k
		}
		switch val := v.(type) {
		case string:
			out[key] = val
		case map[string]any:
			flatten(key, val, out)
		default:
			log.Fatalf("i18n: key %q has unsupported value type %T", key, v)
		}
	}
}

// IsSupported reports whether lang is one of the supported locale codes.
func IsSupported(lang string) bool {
	return slices.Contains(Supported, lang)
}

// T returns the translated string for key in lang, falling back to
// DefaultLang and finally to the key itself if no translation is found.
func T(lang, key string, args ...any) string {
	if cat, ok := catalogs[lang]; ok {
		if s, ok := cat[key]; ok {
			return format(s, args)
		}
	}
	if lang != DefaultLang {
		if cat, ok := catalogs[DefaultLang]; ok {
			if s, ok := cat[key]; ok {
				return format(s, args)
			}
		}
	}
	log.Printf("i18n: missing key %q for lang %q", key, lang)
	return key
}

func format(s string, args []any) string {
	if len(args) == 0 {
		return s
	}
	return fmt.Sprintf(s, args...)
}

// Catalog returns the flattened key set for lang (for testing).
func Catalog(lang string) map[string]string {
	return catalogs[lang]
}
