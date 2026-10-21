package snakesim

import "math"

// Feature layout version of the policy input. Bump it whenever the layout or any formula changes and
// regenerate the model: the Go runtime refuses a model with another version.
const (
	FeatureVersion = 5
	win            = 3 // ego window radius
	grid           = 5 // coarse ego body-density map is grid×grid bins
	perAction      = 7
	nGlobal        = 10
	// NFeatures = 3 actions×7 + 10 global + (2·win+1)² window + grid² density.
	NFeatures = 3*perAction + nGlobal + (2*win+1)*(2*win+1) + grid*grid
)

// regions returns the number of free-space components, the largest and the total free cells (tail free).
func (t *Teacher) regions(g *Game) (count, largest, total int) {
	n := g.N
	t.gen++
	gen := t.gen
	tail := g.Tail()
	for sc := 0; sc < n*n; sc++ {
		if t.seen[sc] == gen || (g.occ[sc] && sc != tail) {
			continue
		}
		t.seen[sc] = gen
		q := append(t.queue[:0], int32(sc))
		for h := 0; h < len(q); h++ {
			c := int(q[h])
			cx, cy := c%n, c/n
			for d := 0; d < 4; d++ {
				x, y := cx+dx[d], cy+dy[d]
				if x < 0 || x >= n || y < 0 || y >= n {
					continue
				}
				m := y*n + x
				if t.seen[m] == gen || (g.occ[m] && m != tail) {
					continue
				}
				t.seen[m] = gen
				q = append(q, int32(m))
			}
		}
		t.queue = q
		count++
		total += len(q)
		if len(q) > largest {
			largest = len(q)
		}
	}
	return count, largest, total
}

func clip(v, lo, hi float64) float32 {
	return float32(math.Max(lo, math.Min(hi, v)))
}

// Features writes the NFeatures-long policy input for the current state into out (egocentric:
// actions are left, straight, right). It requires a food on the board.
func (t *Teacher) Features(g *Game, out []float32) {
	for i := range out[:NFeatures] {
		out[i] = 0
	}
	n := g.N
	nf := float64(n)
	t.foodField(g)
	tail := g.Tail()
	for a := 0; a < 3; a++ {
		o := a * perAction
		c := g.Next(a)
		if c < 0 || (g.occ[c] && (c != tail || g.size == 1)) {
			out[o] = 1
			out[o+3], out[o+4], out[o+5] = 1, 1, 1
			continue
		}
		var v []int32
		if c == g.Food {
			v = t.virtualAfter(g, []int32{int32(c)})
		} else {
			v = append(t.virt[:0], int32(c))
			for i := 0; i < g.size-1; i++ {
				v = append(v, int32(g.Cell(i)))
			}
			t.virt = v
		}
		count, tailOK, tailDist := t.floodDist(g, v)
		out[o+1] = clip(float64(count)/float64(g.size), 0, 2) / 2
		if tailOK {
			out[o+2] = 1
		}
		l := int32(unreachable)
		if c == g.Food {
			l = 1
		} else if t.dist[c] >= 0 {
			l = 1 + t.dist[c]
		}
		out[o+3] = 1
		if l != unreachable {
			out[o+3] = clip(float64(l)/(2*nf), 0, 1)
		}
		out[o+4] = 1 // distance from the new head to the tail after the move, as the tail chase measures it
		if tailDist >= 0 {
			out[o+4] = clip(float64(tailDist)/(2*nf), 0, 1)
		}
		off := lineOffset(g, c)
		out[o+5] = clip(off/3, 0, 1)
		if l != unreachable {
			if ok, _ := t.check(g, cand{a, c, l, off}); ok {
				out[o+6] = 1
			}
		}
	}

	head := g.Head()
	hx, hy := head%n, head/n
	fx, fy := g.Food%n, g.Food/n
	fw, rt := (g.Heading)%4, (g.Heading+1)%4
	ddx, ddy := float64(fx-hx), float64(fy-hy)
	k := 3 * perAction
	out[k] = clip((ddx*float64(dx[rt])+ddy*float64(dy[rt]))/(nf/2), -1, 1)
	out[k+1] = clip((ddx*float64(dx[fw])+ddy*float64(dy[fw]))/(nf/2), -1, 1)
	out[k+2] = float32(g.size) / float32(n*n)
	out[k+3] = clip(float64(g.SinceFood)/(2*nf), 0, 1)
	// side of the head relative to the zigzag line, in the heading frame (positive: line lies to the right)
	ax, ay := float64(g.LegStart%n), float64(g.LegStart/n)
	vx, vy := float64(fx)-ax, float64(fy)-ay
	if norm := math.Hypot(vx, vy); norm > 0 {
		nx, ny := -vy/norm, vx/norm
		s := (float64(hx)-ax)*nx + (float64(hy)-ay)*ny
		out[k+4] = clip(-s*(nx*float64(dx[rt])+ny*float64(dy[rt]))/3, -1, 1)
	}
	nreg, largest, total := t.regions(g)
	out[k+5] = clip(float64(nreg)/4, 0, 1)
	if total > 0 {
		out[k+6] = float32(largest) / float32(total)
	}
	out[k+7] = float32(total) / float32(n*n)
	out[k+8] = clip((math.Abs(ddx)+math.Abs(ddy))/(2*nf), 0, 1)
	out[k+9] = float32(nf / 50)

	// ego window: body or wall within win cells of the head (the tail cell counts as free)
	w := k + nGlobal
	side := 2*win + 1
	for ry := win; ry >= -win; ry-- {
		for rx := -win; rx <= win; rx++ {
			x := hx + rx*dx[rt] + ry*dx[fw]
			y := hy + rx*dy[rt] + ry*dy[fw]
			idx := (win-ry)*side + (rx + win)
			if x < 0 || x >= n || y < 0 || y >= n || (g.occ[y*n+x] && y*n+x != tail) {
				out[w+idx] = 1
			}
		}
	}
	// coarse body density over the whole board, in the heading frame
	d := w + side*side
	bw := nf / grid
	var bins [grid * grid]float32
	for i := 1; i < g.size; i++ {
		c := g.Cell(i)
		cdx, cdy := c%n-hx, c/n-hy
		rx := cdx*dx[rt] + cdy*dy[rt]
		ry := cdx*dx[fw] + cdy*dy[fw]
		bx := int(math.Floor((float64(rx) + nf/2) / bw))
		by := int(math.Floor((float64(-ry) + nf/2) / bw))
		if bx < 0 {
			bx = 0
		} else if bx >= grid {
			bx = grid - 1
		}
		if by < 0 {
			by = 0
		} else if by >= grid {
			by = grid - 1
		}
		bins[by*grid+bx]++
	}
	for i, b := range bins {
		out[d+i] = clip(float64(b)/(bw*bw), 0, 1)
	}
}
