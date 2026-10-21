// Package snakesim is the fast snake engine used to generate training data and, later, to run the
// live game: array-backed board, any size, relative actions (left, straight, right).
package snakesim

import "math/rand/v2"

// Relative actions.
const (
	Left     = 0
	Straight = 1
	Right    = 2
)

// Clockwise: up, right, down, left. Heading h turns right to (h+1)%4.
var (
	dx   = [4]int{0, 1, 0, -1}
	dy   = [4]int{-1, 0, 1, 0}
	turn = [3]int{3, 0, 1}
)

// Game is one snake game on an n×n board. Cells are indexed y*N+x.
type Game struct {
	N         int
	Heading   int
	Food      int // cell index, -1 when the board is full
	Score     int
	SinceFood int
	Dead      bool
	LegStart  int // head cell when the current apple appeared; the zigzag line runs from here to Food

	body []int32 // ring buffer, head first
	head int     // index of the head in body
	size int
	occ  []bool
	rng  *rand.Rand
}

// NewGame places a length-1 snake in the centre with a random heading.
func NewGame(n int, rng *rand.Rand) *Game {
	g := &Game{N: n, body: make([]int32, n*n), occ: make([]bool, n*n), rng: rng}
	c := (n/2)*n + n/2
	g.body[0] = int32(c)
	g.size = 1
	g.occ[c] = true
	g.Heading = rng.IntN(4)
	g.placeFood()
	return g
}

// Len is the snake length.
func (g *Game) Len() int { return g.size }

// Cell returns the i-th body cell, 0 = head.
func (g *Game) Cell(i int) int { return int(g.body[(g.head+i)%len(g.body)]) }

// Head and Tail cells.
func (g *Game) Head() int { return g.Cell(0) }
func (g *Game) Tail() int { return g.Cell(g.size - 1) }

// Occupied reports whether a body cell sits on c (the tail included).
func (g *Game) Occupied(c int) bool { return g.occ[c] }

func (g *Game) placeFood() {
	free := len(g.occ) - g.size
	if free <= 0 {
		g.Food = -1
		return
	}
	k := g.rng.IntN(free)
	for c, o := range g.occ {
		if o {
			continue
		}
		if k == 0 {
			g.Food = c
			g.LegStart = g.Head()
			return
		}
		k--
	}
}

// Next is the cell reached by a relative action, or -1 outside the board.
func (g *Game) Next(action int) int {
	h := g.Head()
	d := (g.Heading + turn[action]) % 4
	x, y := h%g.N+dx[d], h/g.N+dy[d]
	if x < 0 || x >= g.N || y < 0 || y >= g.N {
		return -1
	}
	return y*g.N + x
}

// Step applies a relative action. Sets Dead on collision, starvation or a full board.
func (g *Game) Step(action int) {
	next := g.Next(action)
	if next < 0 {
		g.Dead = true
		return
	}
	eating := next == g.Food
	if g.occ[next] && (eating || next != g.Tail() || g.size == 1) {
		g.Dead = true
		return
	}
	g.Heading = (g.Heading + turn[action]) % 4
	g.SinceFood++
	if !eating {
		t := g.Tail()
		g.occ[t] = false
		g.size--
	}
	g.head = (g.head - 1 + len(g.body)) % len(g.body)
	g.body[g.head] = int32(next)
	g.size++
	g.occ[next] = true
	if eating {
		g.Score++
		g.SinceFood = 0
		g.placeFood()
	}
	if g.Food < 0 || g.SinceFood > 2*g.N*g.N {
		g.Dead = true
	}
}

// Full reports whether the snake filled the board.
func (g *Game) Full() bool { return g.Food < 0 }

// Snapshot is a serialisable game state, used for warm starts.
type Snapshot struct {
	N         int   `json:"n"`
	Heading   int   `json:"heading"`
	Body      []int `json:"body"` // head first
	Food      int   `json:"food"`
	LegStart  int   `json:"leg_start"`
	Score     int   `json:"score"`
	SinceFood int   `json:"since_food"`
}

// Snapshot captures the state.
func (g *Game) Snapshot() Snapshot {
	b := make([]int, g.size)
	for i := range b {
		b[i] = g.Cell(i)
	}
	return Snapshot{g.N, g.Heading, b, g.Food, g.LegStart, g.Score, g.SinceFood}
}

// Restore builds a game from a snapshot with a fresh random source.
func Restore(s Snapshot, rng *rand.Rand) *Game {
	g := &Game{N: s.N, body: make([]int32, s.N*s.N), occ: make([]bool, s.N*s.N), rng: rng,
		Heading: s.Heading, Food: s.Food, LegStart: s.LegStart, Score: s.Score, SinceFood: s.SinceFood}
	for i, c := range s.Body {
		g.body[i] = int32(c)
		g.occ[c] = true
	}
	g.size = len(s.Body)
	return g
}
