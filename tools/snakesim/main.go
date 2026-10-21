// Command snakesim benchmarks the teacher policy and renders replays.
//
//	go run ./tools/snakesim bench -n 50 -games 40
//	go run ./tools/snakesim gif -n 50 -apples 40 -seed 1 -out /tmp/teacher.gif
package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"math/rand/v2"
	"os"
	"sort"
	"sync"

	"mljr-web/internal/snakesim"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: snakesim bench|gif [flags]")
		os.Exit(2)
	}
	switch os.Args[1] {
	case "bench":
		bench(os.Args[2:])
	case "gif":
		render(os.Args[2:])
	default:
		fmt.Fprintln(os.Stderr, "unknown command", os.Args[1])
		os.Exit(2)
	}
}

type result struct {
	score, steps  int
	openingDeath  bool
	spa150        float64 // steps per apple over the first 150 apples, 0 if not reached
	stallApples   int     // apples that took more than 4N steps
	reached25     bool
	stepsTo25pct  int
	diagonalShare float64 // share of moves in the first 150 apples that turned
}

func play(n int, seed uint64, maxSteps int) result {
	rng := rand.New(rand.NewPCG(seed, seed^0x9e3779b97f4a7c15)) //nolint:gosec // deterministic benchmark RNG
	g := snakesim.NewGame(n, rng)
	t := snakesim.NewTeacher(n)
	var r result
	steps150, turns, moves := 0, 0, 0
	lastAte := 0
	prevHeading := g.Heading
	for !g.Dead && r.steps < maxSteps {
		a := t.Action(g)
		before := g.Score
		g.Step(a)
		r.steps++
		if g.Score < 150 || (g.Score == 150 && before < 150) {
			moves++
			if g.Heading != prevHeading {
				turns++
			}
		}
		prevHeading = g.Heading
		if g.Score > before {
			if r.steps-lastAte > 4*n {
				r.stallApples++
			}
			lastAte = r.steps
			if g.Score == 150 {
				steps150 = r.steps
			}
		}
		if !r.reached25 && g.Len() >= n*n/4 {
			r.reached25, r.stepsTo25pct = true, r.steps
		}
	}
	r.score = g.Score
	r.openingDeath = g.Dead && !g.Full() && g.Len() < n*n/4
	if steps150 > 0 {
		r.spa150 = float64(steps150) / 150
	}
	if moves > 0 {
		r.diagonalShare = float64(turns) / float64(moves)
	}
	return r
}

func bench(args []string) {
	fs := flag.NewFlagSet("bench", flag.ExitOnError)
	n := fs.Int("n", 50, "board size")
	games := fs.Int("games", 40, "games")
	maxSteps := fs.Int("steps", 0, "step cap (default 40*n*n)")
	seed := fs.Uint64("seed", 777000, "first seed")
	_ = fs.Parse(args)
	if *maxSteps == 0 {
		*maxSteps = 40 * *n * *n
	}
	res := make([]result, *games)
	var wg sync.WaitGroup
	sem := make(chan struct{}, 16)
	for i := range res {
		wg.Add(1)
		sem <- struct{}{}
		go func() {
			defer wg.Done()
			res[i] = play(*n, *seed+uint64(i), *maxSteps)
			<-sem
		}()
	}
	wg.Wait()
	scores := make([]int, len(res))
	deaths, reach25, stalls := 0, 0, 0
	var spaSum, turnSum float64
	spaN := 0
	for i, r := range res {
		scores[i] = r.score
		if r.openingDeath {
			deaths++
		}
		if r.reached25 {
			reach25++
		}
		stalls += r.stallApples
		if r.spa150 > 0 {
			spaSum += r.spa150
			spaN++
		}
		turnSum += r.diagonalShare
	}
	sort.Ints(scores)
	spa := 0.0
	if spaN > 0 {
		spa = spaSum / float64(spaN)
	}
	fmt.Printf("n=%d games=%d  median score=%d (min %d, max %d)  opening deaths (<25%% fill)=%d/%d  reached 25%%=%d/%d\n",
		*n, *games, scores[len(scores)/2], scores[0], scores[len(scores)-1], deaths, *games, reach25, *games)
	fmt.Printf("steps/apple over first 150 apples=%.1f (%d games reached it)  turn share of moves=%.2f  slow apples (>4N steps) per game=%.1f\n",
		spa, spaN, turnSum/float64(*games), float64(stalls)/float64(*games))
}

func render(args []string) {
	fs := flag.NewFlagSet("gif", flag.ExitOnError)
	n := fs.Int("n", 50, "board size")
	apples := fs.Int("apples", 40, "stop after this many apples")
	seed := fs.Uint64("seed", 1, "seed")
	cell := fs.Int("cell", 10, "pixels per cell")
	delay := fs.Int("delay", 3, "frame delay in 1/100 s")
	out := fs.String("out", "teacher.gif", "output file")
	_ = fs.Parse(args)
	rng := rand.New(rand.NewPCG(*seed, *seed^0x9e3779b97f4a7c15)) //nolint:gosec // deterministic replay RNG
	g := snakesim.NewGame(*n, rng)
	t := snakesim.NewTeacher(*n)
	pal := color.Palette{color.RGBA{250, 247, 240, 255}, color.RGBA{20, 20, 20, 255}, color.RGBA{225, 60, 50, 255}, color.RGBA{60, 140, 90, 255}}
	var anim gif.GIF
	frame := func() {
		px := *n * *cell
		img := image.NewPaletted(image.Rect(0, 0, px, px), pal)
		fill := func(c int, col uint8) {
			x0, y0 := (c%*n)**cell, (c / *n)**cell
			for y := y0; y < y0+*cell-1; y++ {
				for x := x0; x < x0+*cell-1; x++ {
					img.SetColorIndex(x, y, col)
				}
			}
		}
		for i := 1; i < g.Len(); i++ {
			fill(g.Cell(i), 1)
		}
		fill(g.Head(), 3)
		if g.Food >= 0 {
			fill(g.Food, 2)
		}
		anim.Image = append(anim.Image, img)
		anim.Delay = append(anim.Delay, *delay)
	}
	frame()
	for !g.Dead && g.Score < *apples && len(anim.Image) < 20000 {
		g.Step(t.Action(g))
		frame()
	}
	f, err := os.Create(*out)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer f.Close()
	if err := gif.EncodeAll(f, &anim); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("wrote %s: %d frames, %d apples, dead=%v\n", *out, len(anim.Image), g.Score, g.Dead)
}
