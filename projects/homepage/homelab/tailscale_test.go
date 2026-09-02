package homelab

import (
	"reflect"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestInventoryFileParsesTailscaleHosts(t *testing.T) {
	// Mirrors the real homelab-automation shape since the OpenVox cutover:
	// a flat hostname -> tailscale_ip map, plus unrelated sibling keys
	// (services_catalog etc.) that this type must simply ignore.
	src := `
tailscale_hosts:
  mljr: '100.100.20.1'
  nuc: '100.100.10.1'
  nas: '100.100.10.2'
  homeassistant: '100.100.10.200'

services_catalog:
  - name: authelia
`
	var inv inventoryFile
	if err := yaml.Unmarshal([]byte(src), &inv); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	want := map[string]string{
		"mljr":          "100.100.20.1",
		"nuc":           "100.100.10.1",
		"nas":           "100.100.10.2",
		"homeassistant": "100.100.10.200",
	}
	if !reflect.DeepEqual(inv.TailscaleHosts, want) {
		t.Fatalf("inventoryFile.TailscaleHosts = %#v, want %#v", inv.TailscaleHosts, want)
	}
}

func TestFirstDomain(t *testing.T) {
	tests := []struct {
		name string
		in   any
		want string
	}{
		{"single string", "auth.mljr.eu", "auth.mljr.eu"},
		{"list of strings", []any{"git.mljr.eu", "forge.mljr.eu"}, "git.mljr.eu"},
		{"empty list", []any{}, ""},
		{"nil", nil, ""},
		{"wrong type in list", []any{42}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := firstDomain(tt.in); got != tt.want {
				t.Errorf("firstDomain(%#v) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestFetchServicesSkipsDisabled(t *testing.T) {
	src := `
services_catalog:
  - name: authelia
    enabled: true
    domain: "auth.mljr.eu"
    host: mljr
    description: "Authelia SSO"
  - name: staging-site
    enabled: false
    domain: "mljr.eu"
    host: mljr
  - name: nightscout
    domain:
      - "nightscout.mljr.eu"
      - "ns.mljr.eu"
    host: nuc
    description: "Diabetes Management"
`
	var sf servicesFile
	if err := yaml.Unmarshal([]byte(src), &sf); err != nil {
		t.Fatalf("unmarshal: %v", err)
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

	want := []ServiceEntry{
		{Name: "authelia", Domain: "auth.mljr.eu", Host: "mljr", Description: "Authelia SSO"},
		{Name: "nightscout", Domain: "nightscout.mljr.eu", Host: "nuc", Description: "Diabetes Management"},
	}
	if !reflect.DeepEqual(out, want) {
		t.Fatalf("got %#v, want %#v", out, want)
	}
}
