package homelab

import (
	"reflect"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestCollectInventoryHosts(t *testing.T) {
	// Mirrors the real homelab-automation inventory shape: hosts can sit at
	// any depth under nested "children" groups, and not every group has a
	// "hosts" key (e.g. proxy_only's parent has none, only its children do).
	src := `
all:
  children:
    managed:
      children:
        rocky:
          hosts:
            mljr:
              tailscale_ip: 100.100.20.1
            nuc:
              tailscale_ip: 100.100.10.1
        unraid:
          hosts:
            nas:
              tailscale_ip: 100.100.10.2
    proxy_only:
      hosts:
        homeassistant:
          tailscale_ip: 100.100.10.200
        monitoring:
          ansible_host: 192.168.50.175
`
	var inv inventoryFile
	if err := yaml.Unmarshal([]byte(src), &inv); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	out := map[string]string{}
	collectInventoryHosts(inv.All, out)

	want := map[string]string{
		"mljr":          "100.100.20.1",
		"nuc":           "100.100.10.1",
		"nas":           "100.100.10.2",
		"homeassistant": "100.100.10.200",
		// "monitoring" has no tailscale_ip and must be excluded, not zero-valued.
	}
	if !reflect.DeepEqual(out, want) {
		t.Fatalf("collectInventoryHosts() = %#v, want %#v", out, want)
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
services:
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
