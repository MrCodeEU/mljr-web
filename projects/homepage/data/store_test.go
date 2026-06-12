package data

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestStoreReloadsChangedDataFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "site-data.json")
	first := `{"github_projects":[],"linkedin_data":{"name":"First","profile":{"name":"First"}},"strava_data":{}}`
	second := `{"github_projects":[],"linkedin_data":{"name":"Second","profile":{"name":"Second"}},"strava_data":{}}`

	if err := os.WriteFile(path, []byte(first), 0o600); err != nil {
		t.Fatal(err)
	}
	store := NewStore(path, "1")
	if got := store.Current().LinkedIn.Name; got != "First" {
		t.Fatalf("initial Current() = %q, want First", got)
	}

	time.Sleep(1100 * time.Millisecond)
	if err := os.WriteFile(path, []byte(second), 0o600); err != nil {
		t.Fatal(err)
	}
	if got := store.Current().LinkedIn.Name; got != "Second" {
		t.Fatalf("reloaded Current() = %q, want Second", got)
	}
}

func TestStoreKeepsPreviousDataWhenReloadFails(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "site-data.json")
	if err := os.WriteFile(path, []byte(`{"github_projects":[],"linkedin_data":{"name":"Good"},"strava_data":{}}`), 0o600); err != nil {
		t.Fatal(err)
	}
	store := NewStore(path, "1")

	time.Sleep(1100 * time.Millisecond)
	if err := os.WriteFile(path, []byte(`{bad json`), 0o600); err != nil {
		t.Fatal(err)
	}
	if got := store.Current().LinkedIn.Name; got != "Good" {
		t.Fatalf("Current() after bad reload = %q, want previous Good", got)
	}
}

func TestGeneratedSiteDataHasContractMetadata(t *testing.T) {
	d, err := LoadFile("../../../data-repo-dummy/generated/site-data.json")
	if err != nil {
		t.Fatal(err)
	}
	if d.SchemaVersion != "site-data.v1" {
		t.Fatalf("SchemaVersion = %q, want site-data.v1", d.SchemaVersion)
	}
	if _, err := time.Parse(time.RFC3339, d.GeneratedAt); err != nil {
		t.Fatalf("GeneratedAt is not RFC3339: %q", d.GeneratedAt)
	}
	if len(d.GitHub) == 0 {
		t.Fatal("generated data has no GitHub projects")
	}
	if len(d.LinkedIn.Experience) == 0 {
		t.Fatal("generated data has no experience entries")
	}
}
