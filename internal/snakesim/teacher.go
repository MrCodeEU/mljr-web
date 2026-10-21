package snakesim

import "math"

// Teacher labels states with the "impressive" policy:
//
//  1. Take the shortest path to the apple, tie-broken toward the straight line from where the apple
//     appeared (so open-field runs look like diagonal staircases), if a virtual snake that walked the
//     path could still reach its own tail.
//  2. Optional (off by default), past RiskFrac fill, otherwise take the shortest path anyway when at least StallRoom*length free
//     cells stay reachable after eating (accepting a later death beats a visible loop).
//  3. Otherwise chase the tail; with risk enabled, after 2N steps without food take the shortest path unsafely.
//
// A Teacher owns scratch buffers and is not safe for concurrent use; make one per goroutine.
type Teacher struct {
	// StallRoom is the free-space margin of rule 2 (default 1.2).
	StallRoom float64
	// RiskFrac is the fill fraction past which rules 2 and 4 may take risks. The default (2) turns them off:
	// the safe rules alone never died in 600k steps on 50x50, the risky ones died at ~26% fill.
	RiskFrac float64
	// RiskWait is how many steps without food must pass before rule 2 applies (default 0).
	RiskWait int
	// LastRule is the rule that produced the last action: 1 safe path, 2 room margin, 3 tail chase, 4 unsafe path.
	LastRule int

	dist  []int32
	tdist []int32
	queue []int32
	seen  []int32
	blk   []int32
	gen   int32
	path  []int32
	virt  []int32
}

// NewTeacher allocates scratch space for boards up to n×n.
func NewTeacher(n int) *Teacher {
	return &Teacher{
		StallRoom: 1.2,
		RiskFrac:  2, // risk rules off by default: they killed every game at ~26% fill
		dist:      make([]int32, n*n),
		tdist:     make([]int32, n*n),
		queue:     make([]int32, n*n),
		seen:      make([]int32, n*n),
		blk:       make([]int32, n*n),
		path:      make([]int32, 0, n*n),
		virt:      make([]int32, 0, n*n+1),
	}
}

const unreachable = math.MaxInt32

// foodField fills dist with BFS distances from the food over cells that are not body (the tail
// counts as free: it moves away before the head gets there).
func (t *Teacher) foodField(g *Game) {
	n := g.N
	for i := 0; i < n*n; i++ {
		t.dist[i] = -1
	}
	tail := g.Tail()
	t.dist[g.Food] = 0
	q := t.queue[:0]
	q = append(q, int32(g.Food))
	for h := 0; h < len(q); h++ {
		c := int(q[h])
		cx, cy := c%n, c/n
		for d := 0; d < 4; d++ {
			x, y := cx+dx[d], cy+dy[d]
			if x < 0 || x >= n || y < 0 || y >= n {
				continue
			}
			m := y*n + x
			if t.dist[m] >= 0 || (g.occ[m] && m != tail) {
				continue
			}
			t.dist[m] = t.dist[c] + 1
			q = append(q, int32(m))
		}
	}
	t.queue = q
}

// lineOffset is the distance of cell c from the straight line LegStart -> Food.
func lineOffset(g *Game, c int) float64 {
	n := g.N
	ax, ay := g.LegStart%n, g.LegStart/n
	fx, fy := g.Food%n, g.Food/n
	vx, vy := float64(fx-ax), float64(fy-ay)
	norm := math.Hypot(vx, vy)
	if norm == 0 {
		return 0
	}
	return math.Abs(vx*float64(c/n-ay)-vy*float64(c%n-ax)) / norm
}

// walk follows decreasing food distance from first (a cell with dist >= 0) to the food, preferring
// the neighbour closest to the zigzag line, then going straight on. The path excludes the head and
// includes first and the food.
func (t *Teacher) walk(g *Game, first int, prevDir int) []int32 {
	n := g.N
	p := append(t.path[:0], int32(first))
	cur, dir := first, prevDir
	for cur != g.Food {
		best, bestD, bestOff := -1, 0, math.Inf(1)
		for d := 0; d < 4; d++ {
			x, y := cur%n+dx[d], cur/n+dy[d]
			if x < 0 || x >= n || y < 0 || y >= n {
				continue
			}
			m := y*n + x
			if t.dist[m] != t.dist[cur]-1 {
				continue
			}
			off := lineOffset(g, m)
			if off < bestOff-1e-9 || (math.Abs(off-bestOff) <= 1e-9 && d == dir) {
				best, bestD, bestOff = m, d, off
			}
		}
		cur, dir = best, bestD
		p = append(p, int32(cur))
	}
	t.path = p
	return p
}

// flood reports the free cells reachable from the head of virt (head first; the last cell is the
// tail, counted as free) and whether the tail is among them.
func (t *Teacher) flood(g *Game, virt []int32) (count int, tailOK bool) {
	n := g.N
	t.gen++
	gen := t.gen
	if len(virt) > 2 {
		for _, c := range virt[1 : len(virt)-1] {
			t.blk[c] = gen
		}
	}
	head, tail := int(virt[0]), int(virt[len(virt)-1])
	q := t.queue[:0]
	q = append(q, int32(head))
	t.seen[head] = gen
	tailOK = head == tail
	for h := 0; h < len(q); h++ {
		c := int(q[h])
		cx, cy := c%n, c/n
		for d := 0; d < 4; d++ {
			x, y := cx+dx[d], cy+dy[d]
			if x < 0 || x >= n || y < 0 || y >= n {
				continue
			}
			m := y*n + x
			if t.seen[m] == gen || t.blk[m] == gen {
				continue
			}
			t.seen[m] = gen
			if m == tail {
				tailOK = true
			}
			q = append(q, int32(m))
		}
	}
	t.queue = q
	return len(q), tailOK
}

// floodDist is flood plus the BFS distance from the head of virt to its tail (-1 when unreachable).
func (t *Teacher) floodDist(g *Game, virt []int32) (count int, tailOK bool, tailDist int) {
	n := g.N
	t.gen++
	gen := t.gen
	if len(virt) > 2 {
		for _, c := range virt[1 : len(virt)-1] {
			t.blk[c] = gen
		}
	}
	head, tail := int(virt[0]), int(virt[len(virt)-1])
	q := append(t.queue[:0], int32(head))
	t.seen[head] = gen
	t.tdist[head] = 0
	tailDist = -1
	if head == tail {
		tailOK, tailDist = true, 0
	}
	for h := 0; h < len(q); h++ {
		c := int(q[h])
		cx, cy := c%n, c/n
		for d := 0; d < 4; d++ {
			x, y := cx+dx[d], cy+dy[d]
			if x < 0 || x >= n || y < 0 || y >= n {
				continue
			}
			m := y*n + x
			if t.seen[m] == gen || t.blk[m] == gen {
				continue
			}
			t.seen[m] = gen
			t.tdist[m] = t.tdist[c] + 1
			if m == tail {
				tailOK, tailDist = true, int(t.tdist[m])
			}
			q = append(q, int32(m))
		}
	}
	t.queue = q
	return len(q), tailOK, tailDist
}

// virtualAfter returns the body (head first) after walking path and eating, one cell longer.
func (t *Teacher) virtualAfter(g *Game, path []int32) []int32 {
	v := t.virt[:0]
	for i := len(path) - 1; i >= 0; i-- {
		v = append(v, path[i])
	}
	for i := 0; i < g.size; i++ {
		v = append(v, int32(g.Cell(i)))
	}
	v = v[:g.size+1]
	t.virt = v
	return v
}

type cand struct {
	action int
	cell   int
	length int32
	off    float64
}

// Action returns the teacher's relative action for the current state.
func (t *Teacher) Action(g *Game) int {
	t.LastRule = 0
	t.foodField(g)
	var cs [3]cand
	nc := 0
	tail := g.Tail()
	for a := 0; a < 3; a++ {
		c := g.Next(a)
		if c < 0 || (g.occ[c] && (c != tail || g.size == 1)) {
			continue
		}
		l := int32(unreachable)
		if c == g.Food {
			l = 1
		} else if t.dist[c] >= 0 {
			l = 1 + t.dist[c]
		}
		cs[nc] = cand{a, c, l, lineOffset(g, c)}
		nc++
	}
	if nc == 0 {
		return Straight
	}
	// order by (length, line offset, straight first)
	for i := 1; i < nc; i++ {
		for j := i; j > 0 && less(cs[j], cs[j-1]); j-- {
			cs[j], cs[j-1] = cs[j-1], cs[j]
		}
	}
	// rule 1: first candidate whose full path is safe
	for i := 0; i < nc; i++ {
		c := cs[i]
		if c.length == unreachable {
			break
		}
		ok, _ := t.check(g, c)
		if ok {
			t.LastRule = 1
			return c.action
		}
	}
	// rules 2 and 4 take risks, so they only apply once the opening is over (25% fill)
	late := float64(g.size) >= t.RiskFrac*float64(g.N*g.N)
	waited := g.SinceFood >= t.RiskWait
	// rule 2: shortest path with enough room left
	for i := 0; late && waited && i < nc; i++ {
		c := cs[i]
		if c.length == unreachable {
			break
		}
		if _, count := t.check(g, c); float64(count) >= t.StallRoom*float64(g.size+1) {
			t.LastRule = 2
			return c.action
		}
	}
	// rule 3: chase the tail while the stall is young, then give up and go for the apple
	if late && g.SinceFood > 2*g.N {
		for i := 0; i < nc; i++ {
			if cs[i].length != unreachable {
				t.LastRule = 4
				return cs[i].action
			}
		}
	}
	t.LastRule = 3
	return t.chaseTail(g, cs[:nc])
}

func less(a, b cand) bool {
	if a.length != b.length {
		return a.length < b.length
	}
	if math.Abs(a.off-b.off) > 1e-9 {
		return a.off < b.off
	}
	return a.action == Straight && b.action != Straight
}

// check walks the path through c and reports whether the tail stays reachable and how many cells.
func (t *Teacher) check(g *Game, c cand) (tailOK bool, count int) {
	if int(c.cell) == g.Food {
		v := t.virtualAfter(g, []int32{int32(c.cell)})
		count, tailOK = t.flood(g, v)
		return tailOK || len(v) <= 2, count
	}
	d := (g.Heading + turn[c.action]) % 4
	p := t.walk(g, c.cell, d)
	v := t.virtualAfter(g, p)
	count, tailOK = t.flood(g, v)
	return tailOK || len(v) <= 2, count
}

// chaseTail picks the legal move that keeps the tail reachable and is farthest from it, else the
// move with the most reachable space.
func (t *Teacher) chaseTail(g *Game, cs []cand) int {
	best, bestKey := cs[0].action, [3]int{-1, -1, -1}
	for _, c := range cs {
		virt := t.virt[:0]
		virt = append(virt, int32(c.cell))
		for i := 0; i < g.size-1; i++ {
			virt = append(virt, int32(g.Cell(i)))
		}
		t.virt = virt
		count, tailOK, dist := t.floodDist(g, virt)
		if dist < 0 {
			dist = 0
		}
		ok := 0
		if tailOK {
			ok = 1
		}
		key := [3]int{ok, dist, count}
		if !tailOK {
			key = [3]int{0, count, count}
		}
		if key[0] > bestKey[0] || (key[0] == bestKey[0] && (key[1] > bestKey[1] || (key[1] == bestKey[1] && key[2] > bestKey[2]))) {
			best, bestKey = c.action, key
		}
	}
	return best
}
