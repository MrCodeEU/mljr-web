package data

import (
	"bytes"
	"strings"
	"testing"
)

func TestMeshSummarizesServicesWithoutSatelliteNodes(t *testing.T) {
	var out bytes.Buffer
	node := Mesh(MeshProps{}, []MeshNode{
		{Name: "nas", OS: "linux", Online: true, Services: []string{"nextcloud", "grafana", "syncthing"}},
	})
	if err := node.Render(&out); err != nil {
		t.Fatal(err)
	}
	html := out.String()
	if !strings.Contains(html, "3 svc") {
		t.Fatalf("service count missing from mesh: %s", html)
	}
	if !strings.Contains(html, "nas — nextcloud, grafana, syncthing") {
		t.Fatalf("service tooltip missing from mesh: %s", html)
	}
	if strings.Contains(html, "breathe-svc") {
		t.Fatalf("mesh still renders service satellites: %s", html)
	}
}
