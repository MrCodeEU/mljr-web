package i18n

import (
	"sort"
	"testing"
)

// TestCatalogsHaveSameKeys ensures every locale defines exactly the same set
// of translation keys, so a missing translation fails the build instead of
// silently falling back at runtime.
func TestCatalogsHaveSameKeys(t *testing.T) {
	en := Catalog("en")
	de := Catalog("de")

	if len(en) == 0 || len(de) == 0 {
		t.Fatalf("catalogs not loaded: en=%d keys, de=%d keys", len(en), len(de))
	}

	missingInDE := diffKeys(en, de)
	missingInEN := diffKeys(de, en)

	if len(missingInDE) > 0 {
		sort.Strings(missingInDE)
		t.Errorf("keys present in en.json but missing in de.json: %v", missingInDE)
	}
	if len(missingInEN) > 0 {
		sort.Strings(missingInEN)
		t.Errorf("keys present in de.json but missing in en.json: %v", missingInEN)
	}
}

func diffKeys(a, b map[string]string) []string {
	var missing []string
	for k := range a {
		if _, ok := b[k]; !ok {
			missing = append(missing, k)
		}
	}
	return missing
}
