package snakesim

import (
	"math/rand/v2"
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
