// Package snakeai runs a self-playing snake driven by a small neural policy.
// Game rules and feature extraction mirror train/snakenet/env.py exactly;
// testdata/golden.json guards that parity.
package snakeai

import (
	"math/rand/v2"
	"sync"
)

// Point is a board cell.
type Point struct{ X, Y int }

// Clockwise: up, right, down, left. Heading h turns right to (h+1)%4.
var dirs = [4]Point{{0, -1}, {1, 0}, {0, 1}, {-1, 0}}

// Relative actions.
const (
	Left     = 0
	Straight = 1
	Right    = 2
)

var turn = [3]int{3, 0, 1}

// Game is one snake game on an n×n board (n must be even for the cycle features).
type Game struct {
	N         int
	Snake     []Point // head first
	Heading   int
	Food      Point
	HasFood   bool
	Score     int
	SinceFood int
	Dead      bool

	rng *rand.Rand
}

// NewGame places a length-1 snake in the centre with a random heading.
func NewGame(n int, rng *rand.Rand) *Game {
	g := &Game{N: n, Snake: []Point{{n / 2, n / 2}}, Heading: rng.IntN(4), rng: rng}
	g.placeFood()
	return g
}

func (g *Game) placeFood() {
	occ := make([]bool, g.N*g.N)
	for _, c := range g.Snake {
		occ[c.Y*g.N+c.X] = true
	}
	free := make([]int, 0, g.N*g.N-len(g.Snake))
	for i, o := range occ {
		if !o {
			free = append(free, i)
		}
	}
	if len(free) == 0 {
		g.HasFood = false
		return
	}
	i := free[g.rng.IntN(len(free))]
	g.Food = Point{i % g.N, i / g.N}
	g.HasFood = true
}

func (g *Game) inside(c Point) bool { return c.X >= 0 && c.X < g.N && c.Y >= 0 && c.Y < g.N }

// Cell is the cell reached by a relative action.
func (g *Game) Cell(action int) Point {
	d := dirs[(g.Heading+turn[action])%4]
	h := g.Snake[0]
	return Point{h.X + d.X, h.Y + d.Y}
}

// blocked reports wall or body; the tail cell counts as free when tailFree.
func (g *Game) blocked(c Point, tailFree bool) bool {
	if !g.inside(c) {
		return true
	}
	body := g.Snake
	if tailFree && len(body) > 1 {
		body = body[:len(body)-1]
	}
	for _, b := range body {
		if b == c {
			return true
		}
	}
	return false
}

// Step applies a relative action. Sets Dead on collision, starvation or a full board.
func (g *Game) Step(action int) {
	next := g.Cell(action)
	eating := g.HasFood && next == g.Food
	if g.blocked(next, !eating) {
		g.Dead = true
		return
	}
	g.Heading = (g.Heading + turn[action]) % 4
	g.Snake = append([]Point{next}, g.Snake...)
	g.SinceFood++
	if eating {
		g.Score++
		g.SinceFood = 0
		g.placeFood()
	} else {
		g.Snake = g.Snake[:len(g.Snake)-1]
	}
	if !g.HasFood || g.SinceFood > 2*g.N*g.N {
		g.Dead = true
	}
}

// Full reports whether the snake filled the board.
func (g *Game) Full() bool { return !g.HasFood }

var (
	cyclesMu sync.Mutex
	cycles   = map[int]*cycle{}
)

type cycle struct {
	cells []Point
	index []int // y*n+x -> position along the cycle
}

// cycleFor builds the Hamiltonian cycle (same comb layout as the training code).
func cycleFor(n int) *cycle {
	cyclesMu.Lock()
	defer cyclesMu.Unlock()
	if c, ok := cycles[n]; ok {
		return c
	}
	cells := make([]Point, 0, n*n)
	for x := 0; x < n; x++ {
		cells = append(cells, Point{x, 0})
	}
	for y := 1; y < n; y++ {
		cells = append(cells, Point{n - 1, y})
	}
	for col := n - 2; col >= 1; col-- {
		if (n-1-col)%2 == 1 {
			for y := n - 1; y >= 1; y-- {
				cells = append(cells, Point{col, y})
			}
		} else {
			for y := 1; y < n; y++ {
				cells = append(cells, Point{col, y})
			}
		}
	}
	for y := n - 1; y >= 1; y-- {
		cells = append(cells, Point{0, y})
	}
	c := &cycle{cells: cells, index: make([]int, n*n)}
	for i, p := range cells {
		c.index[p.Y*n+p.X] = i
	}
	cycles[n] = c
	return c
}

func cycDist(a, b, length int) int {
	d := ((b-a)%length + length) % length
	if d == 0 {
		return length
	}
	return d
}
