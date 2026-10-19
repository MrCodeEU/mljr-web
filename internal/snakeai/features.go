package snakeai

// Feature layout version; must match FEATURE_VERSION in train/snakenet/env.py.
const (
	FeatureVersion = 3
	win            = 3 // window radius
	grid           = 5 // coarse ego density map is grid×grid bins
	// NFeatures = 3 actions×6 + 3 global + 4 topology + 2 density maps + window + 3 actions×2 cycle + 3 cycle global.
	NFeatures = 3*6 + 3 + 4 + 2*grid*grid + (2*win+1)*(2*win+1) + 3*2 + 3
)

// board is a flat occupancy grid used by the BFS helpers.
type board struct {
	n int
}

func (b board) in(x, y int) bool { return x >= 0 && x < b.n && y >= 0 && y < b.n }

// flood counts free cells reachable from snake[0] (tail counts as free) and whether the tail is reached.
func flood(n int, snake []Point) (int, bool) {
	b := board{n}
	start, tail := snake[0], snake[len(snake)-1]
	blocked := make([]bool, n*n)
	for i := 1; i < len(snake)-1; i++ {
		blocked[snake[i].Y*n+snake[i].X] = true
	}
	seen := make([]bool, n*n)
	seen[start.Y*n+start.X] = true
	q := []Point{start}
	count, tailReached := 1, start == tail
	for len(q) > 0 {
		p := q[0]
		q = q[1:]
		for _, d := range dirs {
			c := Point{p.X + d.X, p.Y + d.Y}
			if !b.in(c.X, c.Y) || seen[c.Y*n+c.X] || blocked[c.Y*n+c.X] {
				continue
			}
			seen[c.Y*n+c.X] = true
			count++
			if c == tail {
				tailReached = true
			}
			q = append(q, c)
		}
	}
	return count, tailReached
}

// bfsDist returns BFS distances from start (-1 = unreachable); blocked cells are impassable.
func bfsDist(n int, start Point, blocked []bool) []int {
	b := board{n}
	dist := make([]int, n*n)
	for i := range dist {
		dist[i] = -1
	}
	dist[start.Y*n+start.X] = 0
	q := []Point{start}
	for len(q) > 0 {
		p := q[0]
		q = q[1:]
		for _, d := range dirs {
			c := Point{p.X + d.X, p.Y + d.Y}
			if !b.in(c.X, c.Y) || dist[c.Y*n+c.X] >= 0 || blocked[c.Y*n+c.X] {
				continue
			}
			dist[c.Y*n+c.X] = dist[p.Y*n+p.X] + 1
			q = append(q, c)
		}
	}
	return dist
}

// freeRegions returns (component count, largest, total free) of cells not in blocked.
func freeRegions(n int, blocked []bool) (int, int, int) {
	b := board{n}
	seen := append([]bool(nil), blocked...)
	count, largest, total := 0, 0, 0
	for sx := 0; sx < n; sx++ {
		for sy := 0; sy < n; sy++ {
			if seen[sy*n+sx] {
				continue
			}
			seen[sy*n+sx] = true
			q := []Point{{sx, sy}}
			size := 0
			for len(q) > 0 {
				p := q[0]
				q = q[1:]
				size++
				for _, d := range dirs {
					c := Point{p.X + d.X, p.Y + d.Y}
					if !b.in(c.X, c.Y) || seen[c.Y*n+c.X] {
						continue
					}
					seen[c.Y*n+c.X] = true
					q = append(q, c)
				}
			}
			count++
			total += size
			if size > largest {
				largest = size
			}
		}
	}
	return count, largest, total
}

func mask(n int, cells []Point) []bool {
	m := make([]bool, n*n)
	for _, c := range cells {
		m[c.Y*n+c.X] = true
	}
	return m
}

func clip(v, lo, hi float64) float64 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func minF(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

// Features builds the policy input for the current state. Requires a food cell.
func (g *Game) Features() []float32 {
	n := g.N
	body := g.Snake
	head := body[0]
	length := len(body)
	f := dirs[g.Heading]
	r := dirs[(g.Heading+1)%4]
	food := g.Food
	dist0 := abs(food.X-head.X) + abs(food.Y-head.Y)
	capD := 2.0 * float64(n)
	out := make([]float32, NFeatures)

	for i, a := range [3]int{Left, Straight, Right} {
		c := g.Cell(a)
		eating := c == food
		blk := g.blocked(c, !eating)
		o := i * 6
		if blk {
			out[o] = 1
		}
		out[o+4], out[o+5] = 1, 1
		if blk {
			continue
		}
		keep := length
		if !eating {
			keep = length - 1
		}
		virt := append([]Point{c}, body[:keep]...)
		count, tailOK := flood(n, virt)
		den := length
		if den < 1 {
			den = 1
		}
		out[o+1] = float32(minF(float64(count)/float64(den), 2.0) / 2.0)
		if tailOK {
			out[o+2] = 1
		}
		out[o+3] = float32(dist0 - (abs(food.X-c.X) + abs(food.Y-c.Y)))
		var inner []Point
		if len(virt) > 2 {
			inner = virt[1 : len(virt)-1]
		}
		dist := bfsDist(n, c, mask(n, inner))
		if d := dist[virt[len(virt)-1].Y*n+virt[len(virt)-1].X]; d >= 0 {
			out[o+4] = float32(minF(float64(d)/capD, 1.0))
		}
		if d := dist[food.Y*n+food.X]; d >= 0 {
			out[o+5] = float32(minF(float64(d)/capD, 1.0))
		}
	}

	dx, dy := food.X-head.X, food.Y-head.Y
	half := float64(n) / 2
	out[18] = float32(clip(float64(dx*r.X+dy*r.Y)/half, -1, 1))
	out[19] = float32(clip(float64(dx*f.X+dy*f.Y)/half, -1, 1))
	out[20] = float32(float64(length) / float64(n*n))

	// free-space topology: tail counts as free, head as blocked
	bodyNoTail := body[:length-1]
	nReg, largest, total := freeRegions(n, mask(n, bodyNoTail))
	out[21] = float32(float64(minInt(nReg, 4)) / 4.0)
	if total > 0 {
		out[22] = float32(float64(largest) / float64(total))
	}
	out[23] = float32(float64(total) / float64(n*n))

	var inner []Point
	if length > 2 {
		inner = body[1 : length-1]
	}
	dist := bfsDist(n, head, mask(n, inner))
	out[24] = 1
	if d := dist[food.Y*n+food.X]; d >= 0 {
		out[24] = float32(minF(float64(d)/capD, 1.0))
	}

	// coarse ego density maps: body fraction + in-board fraction per bin
	bodyMask := mask(n, bodyNoTail)
	bodyCnt := make([]float32, grid*grid)
	areaCnt := make([]float32, grid*grid)
	bin := func(x, y int) int {
		ox, oy := x-head.X, y-head.Y
		u := ox*r.X + oy*r.Y
		v := ox*f.X + oy*f.Y
		bu := minInt((u+n)*grid/(2*n), grid-1)
		bv := minInt((n-v)*grid/(2*n), grid-1)
		return bv*grid + bu
	}
	for x := 0; x < n; x++ {
		for y := 0; y < n; y++ {
			areaCnt[bin(x, y)]++
			if bodyMask[y*n+x] {
				bodyCnt[bin(x, y)]++
			}
		}
	}
	cellArea := float32(float64(2*n) / float64(grid) * (float64(2*n) / float64(grid)))
	for i := 0; i < grid*grid; i++ {
		out[25+i] = bodyCnt[i] / cellArea
		out[25+grid*grid+i] = areaCnt[i] / cellArea
	}

	k := 25 + 2*grid*grid
	for v := win; v >= -win; v-- {
		for u := -win; u <= win; u++ {
			x := head.X + u*r.X + v*f.X
			y := head.Y + u*r.Y + v*f.Y
			if x < 0 || x >= n || y < 0 || y >= n || bodyMask[y*n+x] {
				out[k] = 1
			}
			k++
		}
	}

	// cycle features
	cyc := cycleFor(n)
	cyclen := len(cyc.cells)
	hi := cyc.index[head.Y*n+head.X]
	for i, a := range [3]int{Left, Straight, Right} {
		c := g.Cell(a)
		if !g.inside(c) {
			out[k+i*2] = 1
			out[k+i*2+1] = 1
			continue
		}
		out[k+i*2+1] = float32(float64(cycDist(hi, cyc.index[c.Y*n+c.X], cyclen)) / float64(cyclen))
	}
	tail := body[length-1]
	out[k+6] = float32(float64(cycDist(hi, cyc.index[tail.Y*n+tail.X], cyclen)) / float64(cyclen))
	out[k+7] = float32(float64(cycDist(hi, cyc.index[food.Y*n+food.X], cyclen)) / float64(cyclen))
	out[k+8] = float32(float64(length) / float64(cyclen))
	return out
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
