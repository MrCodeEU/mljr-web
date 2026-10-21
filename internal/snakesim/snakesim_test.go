package snakesim

import (
	"encoding/json"
	"math/rand/v2"
	"os"
	"testing"
)

func newRNG(seed uint64) *rand.Rand {
	return rand.New(rand.NewPCG(seed, seed^0x9e3779b97f4a7c15)) //nolint:gosec // deterministic test RNG
}

func TestStepBasics(t *testing.T) {
	g := NewGame(10, newRNG(1))
	g.Heading = 1 // right
	h := g.Head()
	g.Step(Straight)
	if g.Dead || g.Head() != h+1 || g.Len() != 1 {
		t.Fatalf("straight: head=%d want %d len=%d dead=%v", g.Head(), h+1, g.Len(), g.Dead)
	}
	// wall: walk straight into the right edge
	for i := 0; i < 20 && !g.Dead; i++ {
		g.Step(Straight)
	}
	if !g.Dead {
		t.Fatal("expected a wall collision")
	}
}

func TestEatGrowsAndRespawns(t *testing.T) {
	g := NewGame(10, newRNG(2))
	h := g.Head()
	g.Heading = 1
	g.Food = h + 1
	g.Step(Straight)
	if g.Score != 1 || g.Len() != 2 || g.Food == g.Head() || g.Occupied(g.Food) {
		t.Fatalf("score=%d len=%d food=%d", g.Score, g.Len(), g.Food)
	}
}

func TestOccupancyMatchesBody(t *testing.T) {
	g := NewGame(12, newRNG(3))
	te := NewTeacher(12)
	for i := 0; i < 3000 && !g.Dead; i++ {
		g.Step(te.Action(g))
		count := 0
		for _, o := range g.occ {
			if o {
				count++
			}
		}
		if count != g.Len() {
			t.Fatalf("step %d: occupancy %d != length %d", i, count, g.Len())
		}
		for k := 0; k < g.Len(); k++ {
			if !g.occ[g.Cell(k)] {
				t.Fatalf("body cell %d not marked occupied", k)
			}
		}
	}
}

// The opening (below 25% fill) must be death-free: that is the part viewers see first.
func TestOpeningNeverDies(t *testing.T) {
	for _, n := range []int{10, 20, 30, 50} {
		for seed := uint64(1); seed <= 4; seed++ {
			g := NewGame(n, newRNG(seed*100+uint64(n)))
			te := NewTeacher(n)
			for steps := 0; g.Len()*4 < n*n && steps < 60*n*n; steps++ {
				g.Step(te.Action(g))
				if g.Dead {
					t.Fatalf("n=%d seed=%d died at length %d", n, seed, g.Len())
				}
			}
		}
	}
}

func TestFeaturesShapeAndRange(t *testing.T) {
	for _, n := range []int{10, 30, 50} {
		g := NewGame(n, newRNG(uint64(n)))
		te := NewTeacher(n)
		x := make([]float32, NFeatures)
		for i := 0; i < 400 && !g.Dead; i++ {
			te.Features(g, x)
			for k, v := range x {
				if v != v || v < -1.0001 || v > 2.0001 {
					t.Fatalf("n=%d step %d feature %d out of range: %v", n, i, k, v)
				}
			}
			g.Step(te.Action(g))
		}
	}
}

// Features are egocentric: a state and its mirror image give mirrored per-action blocks.
func TestFeaturesDeterministic(t *testing.T) {
	g := NewGame(20, newRNG(9))
	te := NewTeacher(20)
	for i := 0; i < 200 && !g.Dead; i++ {
		g.Step(te.Action(g))
	}
	a, b := make([]float32, NFeatures), make([]float32, NFeatures)
	te.Features(g, a)
	te.Features(g, b)
	for i := range a {
		if a[i] != b[i] {
			t.Fatalf("feature %d differs between calls", i)
		}
	}
}

// The embedded model's forward pass must match the independent numpy implementation that wrote the
// golden vectors (tools/snakeai-train/snakenet/golden.py). Regenerate both together with the model.
func TestModelMatchesGolden(t *testing.T) {
	data, err := os.ReadFile("testdata/golden.json")
	if err != nil {
		t.Fatal(err)
	}
	var golden struct {
		FeatureVersion int `json:"feature_version"`
		Cases          []struct {
			Features []float32 `json:"features"`
			Logits   []float32 `json:"logits"`
			Action   int       `json:"action"`
		} `json:"cases"`
	}
	if err := json.Unmarshal(data, &golden); err != nil {
		t.Fatal(err)
	}
	if golden.FeatureVersion != FeatureVersion || len(golden.Cases) == 0 {
		t.Fatalf("golden is for feature v%d with %d cases", golden.FeatureVersion, len(golden.Cases))
	}
	m, err := DefaultModel()
	if err != nil {
		t.Fatal(err)
	}
	bufA, bufB := make([]float32, m.Width()), make([]float32, m.Width())
	for i, c := range golden.Cases {
		got := m.Logits(c.Features)
		for k := range got {
			if d := got[k] - c.Logits[k]; d > 1e-3 || d < -1e-3 {
				t.Fatalf("case %d logit %d: go %v python %v", i, k, got[k], c.Logits[k])
			}
		}
		if a := m.Act(c.Features, bufA, bufB); a != c.Action {
			t.Fatalf("case %d: Act %d, golden %d", i, a, c.Action)
		}
	}
}

// The guard must never pick a blocked move or one that loses the tail when a safe one exists, and must
// otherwise follow the net's ranking.
func TestGuardedActVetoesTrappingMoves(t *testing.T) {
	m, err := DefaultModel()
	if err != nil {
		t.Fatal(err)
	}
	feat := make([]float32, NFeatures)
	for trial := 0; trial < 3; trial++ {
		for i := range feat {
			feat[i] = 0
		}
		feat[0*perAction+2] = 1 // left: open, tail reachable
		feat[1*perAction] = 1   // straight: blocked
		feat[2*perAction] = 0   // right: open but trap
		got := GuardedAct(m, feat)
		if got != Left {
			t.Fatalf("expected Left (only safe move), got %d", got)
		}
	}
	feat[0*perAction] = 1 // left blocked too: only the trap remains
	feat[0*perAction+2] = 0
	if got := GuardedAct(m, feat); got != Right {
		t.Fatalf("expected Right (only unblocked move), got %d", got)
	}
}
