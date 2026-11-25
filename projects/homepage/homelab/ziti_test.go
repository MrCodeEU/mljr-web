package homelab

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type series struct {
	labels map[string]string
	value  float64
}

// fakeVM answers instant queries from a fixed table, like VictoriaMetrics.
// A query missing from the table returns an empty result.
func fakeVM(t *testing.T, answers map[string][]series) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		result := []map[string]any{}
		for _, s := range answers[r.URL.Query().Get("query")] {
			result = append(result, map[string]any{
				"metric": s.labels,
				"value":  []any{float64(time.Now().Unix()), fmt.Sprint(s.value)},
			})
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"status": "success",
			"data":   map[string]any{"resultType": "vector", "result": result},
		})
	}))
	t.Cleanup(srv.Close)
	return srv
}

func router(name string, v float64) series {
	return series{map[string]string{"router": name}, v}
}

func link(src, dst string, v float64) series {
	return series{map[string]string{"source": src, "target": dst}, v}
}

func liveMesh(nasOnline float64) map[string][]series {
	return map[string][]series{
		qRouterOnline: {router("mljr", 1), router("nuc", 1), router("nas", nasOnline), router("ugreen", 1)},
		qLinkUp: {
			link("nas", "mljr", 1), link("nas", "nuc", 1), link("nuc", "mljr", 1),
			link("ugreen", "mljr", 1), link("ugreen", "nuc", 1),
		},
		qLinkLatency: {link("nas", "nuc", 0.005), link("nas", "mljr", 0.025)},
		qTerminators: {router("mljr", 12), router("nuc", 20)},
		qCircuits:    {{map[string]string{}, 7}},
	}
}

func meshHost(s Snapshot, name string) MeshHost {
	for _, h := range s.Mesh {
		if h.Name == name {
			return h
		}
	}
	return MeshHost{}
}

func TestZitiMeshAllOnline(t *testing.T) {
	p := New(Options{PromURL: fakeVM(t, liveMesh(1)).URL})
	var s Snapshot
	p.fetchZitiMesh(context.Background(), &s)

	if !s.MeshOK || len(s.Mesh) != 4 || len(s.MeshLinks) != 5 {
		t.Fatalf("MeshOK=%v hosts=%d links=%d", s.MeshOK, len(s.Mesh), len(s.MeshLinks))
	}
	if s.Mesh[0].Name != "mljr" || !s.Mesh[0].Controller {
		t.Errorf("first host = %+v, want the controller mljr", s.Mesh[0])
	}
	if h := meshHost(s, "nuc"); h.LinksUp != 3 || h.LinksTotal != 3 || h.Services != 20 {
		t.Errorf("nuc = %+v", h)
	}
	if h := meshHost(s, "ugreen"); h.Services != 0 || h.Controller {
		t.Errorf("ugreen = %+v", h)
	}
	if l := s.MeshLinks[1]; l.Source != "nas" || l.Target != "nuc" || l.Latency != 5*time.Millisecond {
		t.Errorf("nas>nuc link = %+v", l)
	}
	if s.MeshCircuits != 7 {
		t.Errorf("circuits = %d, want 7", s.MeshCircuits)
	}
}

func TestZitiMeshRouterOffline(t *testing.T) {
	p := New(Options{PromURL: fakeVM(t, liveMesh(0)).URL})
	var s Snapshot
	p.fetchZitiMesh(context.Background(), &s)

	if h := meshHost(s, "nas"); h.Online || h.LinksUp != 0 || h.LinksTotal != 2 {
		t.Errorf("nas = %+v, want offline with 0/2 links", h)
	}
	if h := meshHost(s, "nuc"); h.LinksUp != 2 || h.LinksTotal != 3 {
		t.Errorf("nuc = %+v, want 2/3 links", h)
	}
	for _, l := range s.MeshLinks {
		if (l.Source == "nas" || l.Target == "nas") && (l.Up || l.Latency != 0) {
			t.Errorf("link %+v to offline nas counts as up", l)
		}
	}
}

func TestZitiMeshEmptyKeepsPrevious(t *testing.T) {
	answers := liveMesh(1)
	p := New(Options{PromURL: fakeVM(t, answers).URL})
	p.poll(context.Background())
	if !p.Snapshot().MeshOK {
		t.Fatal("first poll found no mesh")
	}

	delete(answers, qRouterOnline) // exporter stopped
	p.poll(context.Background())
	s := p.Snapshot()
	if !s.MeshOK || len(s.Mesh) != 4 || len(s.MeshLinks) != 5 || s.MeshCircuits != 7 {
		t.Errorf("after empty poll: MeshOK=%v hosts=%d links=%d circuits=%d", s.MeshOK, len(s.Mesh), len(s.MeshLinks), s.MeshCircuits)
	}
}

func TestZitiMeshFirstPollWithoutDataHidesMesh(t *testing.T) {
	p := New(Options{PromURL: fakeVM(t, nil).URL})
	p.poll(context.Background())
	if s := p.Snapshot(); s.MeshOK || s.MeshCircuits != -1 {
		t.Errorf("MeshOK=%v circuits=%d, want hidden", s.MeshOK, s.MeshCircuits)
	}
}

// The queries aggregate away every label but router names, and the mesh
// keeps only routers that report status. Identity names must not leak even
// if a link or terminator series carries one.
func TestZitiMeshNoIdentityNames(t *testing.T) {
	answers := liveMesh(1)
	answers[qLinkUp] = append(answers[qLinkUp], link("phone-mrcode", "mljr", 1))
	answers[qTerminators] = append(answers[qTerminators], router("laptop-mrcode", 3))
	p := New(Options{PromURL: fakeVM(t, answers).URL})
	var s Snapshot
	p.fetchZitiMesh(context.Background(), &s)

	got := fmt.Sprintf("%+v", s)
	if strings.Contains(got, "mrcode") {
		t.Errorf("snapshot leaks an identity name: %s", got)
	}
	for _, q := range []string{qRouterOnline, qLinkUp, qLinkLatency, qTerminators, qCircuits} {
		if strings.Contains(q, "identity") || strings.Contains(q, "client") {
			t.Errorf("query %q groups by an identity label", q)
		}
	}
}
