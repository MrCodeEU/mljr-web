package snakeai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"math/rand/v2"
	"os"
	"testing"
	"time"
)

type goldenCase struct {
	N        int       `json:"n"`
	Snake    [][2]int  `json:"snake"`
	Heading  int       `json:"heading"`
	Food     [2]int    `json:"food"`
	Features []float32 `json:"features"`
	Logits   []float32 `json:"logits"`
	Action   int       `json:"action"`
}

func loadGolden(t *testing.T) []goldenCase {
	t.Helper()
	raw, err := os.ReadFile("testdata/golden.json")
	if err != nil {
		t.Fatal(err)
	}
	var g struct {
		FeatureVersion int          `json:"feature_version"`
		Cases          []goldenCase `json:"cases"`
	}
	if err := json.Unmarshal(raw, &g); err != nil {
		t.Fatal(err)
	}
	if g.FeatureVersion != FeatureVersion {
		t.Fatalf("golden feature version %d != %d", g.FeatureVersion, FeatureVersion)
	}
	return g.Cases
}

func gameFrom(c goldenCase) *Game {
	g := &Game{N: c.N, Heading: c.Heading, Food: Point{c.Food[0], c.Food[1]}, HasFood: true,
		rng: rand.New(rand.NewPCG(1, 2))} //nolint:gosec // deterministic test RNG
	for _, p := range c.Snake {
		g.Snake = append(g.Snake, Point{p[0], p[1]})
	}
	return g
}

func TestFeaturesMatchPython(t *testing.T) {
	for i, c := range loadGolden(t) {
		got := gameFrom(c).Features()
		if len(got) != len(c.Features) {
			t.Fatalf("case %d: %d features, want %d", i, len(got), len(c.Features))
		}
		for k := range got {
			if math.Abs(float64(got[k]-c.Features[k])) > 1e-5 {
				t.Errorf("case %d (n=%d) feature %d: got %v want %v", i, c.N, k, got[k], c.Features[k])
			}
		}
	}
}

func TestModelMatchesPython(t *testing.T) {
	m, err := DefaultModel()
	if err != nil {
		t.Fatal(err)
	}
	for i, c := range loadGolden(t) {
		z := m.Logits(c.Features)
		for k := range z {
			if math.Abs(float64(z[k]-c.Logits[k])) > 1e-3 {
				t.Errorf("case %d logit %d: got %v want %v", i, k, z[k], c.Logits[k])
			}
		}
		if got := m.Act(c.Features); got != c.Action {
			t.Errorf("case %d action: got %d want %d", i, got, c.Action)
		}
	}
}

// End to end: the net should survive a full 10x10 game and fill it.
func TestSelfPlay10x10(t *testing.T) {
	m, err := DefaultModel()
	if err != nil {
		t.Fatal(err)
	}
	for seed := uint64(0); seed < 4; seed++ {
		g := NewGame(10, rand.New(rand.NewPCG(seed, 99))) //nolint:gosec // deterministic test RNG
		for steps := 0; !g.Dead && steps < 20000; steps++ {
			g.Step(m.Act(g.Features()))
		}
		if g.Score < 99 {
			t.Errorf("seed %d: score %d (dead=%v), want 99", seed, g.Score, g.Dead)
		}
	}
}

func TestRoomLifecycle(t *testing.T) {
	r := NewRoom(Config{N: 10, Tick: time.Millisecond, Grace: 60 * time.Millisecond, RestartDelay: 5 * time.Millisecond})
	if r.Running() {
		t.Fatal("room should start idle")
	}
	a, b := r.Subscribe(), r.Subscribe()
	if !r.Running() || r.Viewers() != 2 {
		t.Fatalf("running=%v viewers=%d", r.Running(), r.Viewers())
	}
	// both viewers see the same shared state: identical bytes for any tick both received
	got := [2]map[string][]byte{{}, {}}
	for i := 0; i < 40; i++ {
		for k, sub := range []*Subscription{a, b} {
			msg := <-sub.C
			if !bytes.HasPrefix(msg, []byte("event: datastar-patch-signals\ndata: signals {")) {
				t.Fatalf("bad SSE message: %q", msg)
			}
			var env map[string]frame
			payload := bytes.TrimSuffix(bytes.TrimPrefix(msg, []byte("event: datastar-patch-signals\ndata: signals ")), []byte("\n\n"))
			if err := json.Unmarshal(payload, &env); err != nil {
				t.Fatal(err)
			}
			got[k][fmt.Sprint(env["snake"].Tick)] = msg
		}
	}
	common := 0
	for tick, ma := range got[0] {
		if mb, ok := got[1][tick]; ok {
			common++
			if !bytes.Equal(ma, mb) {
				t.Fatalf("tick %s differs between viewers", tick)
			}
		}
	}
	if common < 5 {
		t.Fatalf("only %d common ticks", common)
	}
	a.Close()
	b.Close()
	b.Close() // idempotent
	if r.Viewers() != 0 {
		t.Fatalf("viewers=%d", r.Viewers())
	}
	deadline := time.Now().Add(2 * time.Second)
	for r.Running() && time.Now().Before(deadline) {
		time.Sleep(10 * time.Millisecond)
	}
	if r.Running() {
		t.Fatal("room should stop after the grace period")
	}
	r.mu.Lock()
	loaded := r.model != nil
	r.mu.Unlock()
	if loaded {
		t.Fatal("model should be unloaded when idle")
	}
	c := r.Subscribe() // restart on demand
	defer c.Close()
	if !r.Running() {
		t.Fatal("room should restart for a new viewer")
	}
	<-c.C
}

func TestRoomViewerCap(t *testing.T) {
	r := NewRoom(Config{N: 10, Tick: time.Millisecond, Grace: 10 * time.Millisecond, MaxViewers: 1})
	a := r.Subscribe()
	defer a.Close()
	if r.Subscribe() != nil {
		t.Fatal("second viewer should be rejected")
	}
}
