// Command snakesim benchmarks the teacher policy and renders replays.
//
//	go run ./tools/snakesim bench -n 50 -games 40
//	go run ./tools/snakesim gif -n 50 -apples 40 -seed 1 -out /tmp/teacher.gif
package main

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"math"
	"math/rand/v2"
	"os"
	"sort"
	"strconv"
	"strings"
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
	case "gen":
		gen(os.Args[2:])
	case "warmlib":
		warmlib(os.Args[2:])
	default:
		fmt.Fprintln(os.Stderr, "unknown command", os.Args[1])
		os.Exit(2)
	}
}

type result struct {
	score, steps  int
	died          bool
	openingDeath  bool
	spa150        float64 // steps per apple over the first 150 apples, 0 if not reached
	stallApples   int     // apples that took more than 4N steps
	reached25     bool
	stepsTo25pct  int
	diagonalShare float64 // share of moves in the first 150 apples that turned
	bandSteps     [5]int
	bandApples    [5]int
}

type cfg struct {
	model    *snakesim.Model
	riskFrac float64
	riskWait int
	room     float64
}

func play(n int, seed uint64, maxSteps int, c cfg) result {
	rng := rand.New(rand.NewPCG(seed, seed^0x9e3779b97f4a7c15)) //nolint:gosec // deterministic benchmark RNG
	g := snakesim.NewGame(n, rng)
	t := snakesim.NewTeacher(n)
	t.RiskFrac, t.RiskWait, t.StallRoom = c.riskFrac, c.riskWait, c.room
	var r result
	steps150, turns, moves := 0, 0, 0
	lastAte := 0
	prevHeading := g.Heading
	var feat, bufA, bufB []float32
	if c.model != nil {
		feat = make([]float32, snakesim.NFeatures)
		bufA, bufB = make([]float32, c.model.Width()), make([]float32, c.model.Width())
	}
	for !g.Dead && r.steps < maxSteps {
		var a int
		if c.model != nil {
			t.Features(g, feat)
			a = c.model.Act(feat, bufA, bufB)
		} else {
			a = t.Action(g)
		}
		before := g.Score
		b := band(g.Len(), n)
		g.Step(a)
		r.steps++
		r.bandSteps[b]++
		r.bandApples[b] += g.Score - before
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
	r.died = g.Dead
	r.openingDeath = g.Dead && !g.Full() && g.Len() < n*n/4
	if steps150 > 0 {
		r.spa150 = float64(steps150) / 150
	}
	if moves > 0 {
		r.diagonalShare = float64(turns) / float64(moves)
	}
	return r
}

// band maps snake length to a fill band: <10%, <25%, <50%, <75%, rest.
func band(length, n int) int {
	f := float64(length) / float64(n*n)
	switch {
	case f < 0.10:
		return 0
	case f < 0.25:
		return 1
	case f < 0.50:
		return 2
	case f < 0.75:
		return 3
	}
	return 4
}

func bench(args []string) {
	fs := flag.NewFlagSet("bench", flag.ExitOnError)
	n := fs.Int("n", 50, "board size")
	games := fs.Int("games", 40, "games")
	maxSteps := fs.Int("steps", 0, "step cap (default 40*n*n)")
	seed := fs.Uint64("seed", 777000, "first seed")
	riskFrac := fs.Float64("risk", 2, "fill fraction past which risky rules apply")
	riskWait := fs.Int("wait", 0, "steps without food before the room rule applies")
	room := fs.Float64("room", 1.2, "free-space margin of the room rule")
	modelPath := fs.String("model", "", "play with this model.json instead of the teacher")
	_ = fs.Parse(args)
	c := cfg{nil, *riskFrac, *riskWait, *room}
	if *modelPath != "" {
		c.model = loadModel(*modelPath)
	}
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
			res[i] = play(*n, *seed+uint64(i), *maxSteps, c)
			<-sem
		}()
	}
	wg.Wait()
	scores := make([]int, len(res))
	deaths, reach25, stalls, died := 0, 0, 0, 0
	stepSum := 0
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
		stepSum += r.steps
		if r.died {
			died++
		}
		if r.spa150 > 0 {
			spaSum += r.spa150
			spaN++
		}
		turnSum += r.diagonalShare
	}
	var bs, ba [5]int
	for _, r := range res {
		for i := range bs {
			bs[i] += r.bandSteps[i]
			ba[i] += r.bandApples[i]
		}
	}
	sort.Ints(scores)
	spa := 0.0
	if spaN > 0 {
		spa = spaSum / float64(spaN)
	}
	fmt.Printf("risk=%.2f wait=%d room=%.1f  ", *riskFrac, *riskWait, *room)
	fmt.Printf("n=%d games=%d  median score=%d (min %d, max %d)  opening deaths (<25%% fill)=%d/%d  reached 25%%=%d/%d  died=%d/%d  mean steps=%d\n",
		*n, *games, scores[len(scores)/2], scores[0], scores[len(scores)-1], deaths, *games, reach25, *games, died, *games, stepSum / *games)
	fmt.Print("steps/apple by fill [<10% <25% <50% <75% rest]:")
	for i := range bs {
		if ba[i] == 0 {
			fmt.Print("      -")
			continue
		}
		fmt.Printf(" %6.1f", float64(bs[i])/float64(ba[i]))
	}
	fmt.Println()
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

func loadModel(path string) *snakesim.Model {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	m, err := snakesim.LoadModel(data)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return m
}

// gen plays games and writes (teacher label, features) records for DAgger: 1 byte label then
// NFeatures little-endian float32. With -model the net drives (probability 1-beta per step) and the
// teacher still labels every state.
func gen(args []string) {
	fs := flag.NewFlagSet("gen", flag.ExitOnError)
	games := fs.Int("games", 400, "games to play")
	sizes := fs.String("sizes", "20,30,40,50", "board sizes, cycled")
	maxSteps := fs.Int("steps", 1500, "recorded steps per game")
	seed := fs.Uint64("seed", 1, "first seed")
	beta := fs.Float64("beta", 1, "probability of executing the teacher action")
	warm := fs.Float64("warm", 0.75, "share of games that start from a library snapshot")
	libPath := fs.String("lib", "", "warm-start library from the warmlib command")
	modelPath := fs.String("model", "", "model.json driving the non-teacher steps")
	out := fs.String("out", "data.bin", "output file")
	procs := fs.Int("procs", 16, "parallel workers")
	_ = fs.Parse(args)
	var model *snakesim.Model
	if *modelPath != "" {
		model = loadModel(*modelPath)
	}
	var ns []int
	for _, f := range strings.Split(*sizes, ",") {
		v, err := strconv.Atoi(strings.TrimSpace(f))
		if err != nil {
			fmt.Fprintln(os.Stderr, "bad -sizes:", err)
			os.Exit(2)
		}
		ns = append(ns, v)
	}
	lib := map[int][]snakesim.Snapshot{}
	if *libPath != "" {
		data, err := os.ReadFile(*libPath)
		if err == nil {
			err = json.Unmarshal(data, &lib)
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
	f, err := os.Create(*out)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	w := bufio.NewWriterSize(f, 1<<20)
	var mu sync.Mutex
	records, deaths := 0, 0
	var wg sync.WaitGroup
	sem := make(chan struct{}, *procs)
	for i := 0; i < *games; i++ {
		wg.Add(1)
		sem <- struct{}{}
		go func() {
			defer wg.Done()
			defer func() { <-sem }()
			n := ns[i%len(ns)]
			sd := *seed*1_000_003 + uint64(i)
			rng := rand.New(rand.NewPCG(sd, sd^0x9e3779b97f4a7c15)) //nolint:gosec // deterministic data generation
			g := snakesim.NewGame(n, rng)
			t := snakesim.NewTeacher(n)
			if snaps := lib[n]; len(snaps) > 0 && rng.Float64() < *warm {
				g = snakesim.Restore(snaps[rng.IntN(len(snaps))], rng)
			}
			feat := make([]float32, snakesim.NFeatures)
			var bufA, bufB []float32
			if model != nil {
				bufA, bufB = make([]float32, model.Width()), make([]float32, model.Width())
			}
			buf := make([]byte, 0, *maxSteps*(1+4*snakesim.NFeatures))
			count := 0
			for step := 0; step < *maxSteps && !g.Dead && !g.Full(); step++ {
				label := t.Action(g)
				t.Features(g, feat)
				buf = append(buf, byte(label))
				for _, v := range feat {
					buf = binary.LittleEndian.AppendUint32(buf, math.Float32bits(v))
				}
				count++
				a := label
				if model != nil && rng.Float64() >= *beta {
					a = model.Act(feat, bufA, bufB)
				}
				g.Step(a)
			}
			mu.Lock()
			_, _ = w.Write(buf)
			records += count
			if g.Dead && !g.Full() {
				deaths++
			}
			mu.Unlock()
		}()
	}
	wg.Wait()
	if err := w.Flush(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	_ = f.Close()
	fmt.Printf("wrote %s: %d records (%d bytes each), %d/%d games ended in a death\n", *out, records, 1+4*snakesim.NFeatures, deaths, *games)
}

// warmlib plays teacher games per board size and snapshots them every 2% of fill up to -max, so that
// training rounds can start from late-game states without replaying the whole game.
func warmlib(args []string) {
	fs := flag.NewFlagSet("warmlib", flag.ExitOnError)
	sizes := fs.String("sizes", "20,30,40,50", "board sizes")
	per := fs.Int("games", 24, "teacher games per size")
	maxFill := fs.Float64("max", 0.6, "stop a game at this fill fraction")
	out := fs.String("out", "warmlib.json", "output file")
	_ = fs.Parse(args)
	lib := map[int][]snakesim.Snapshot{}
	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, 16)
	for _, f := range strings.Split(*sizes, ",") {
		n, err := strconv.Atoi(strings.TrimSpace(f))
		if err != nil {
			fmt.Fprintln(os.Stderr, "bad -sizes:", err)
			os.Exit(2)
		}
		for i := 0; i < *per; i++ {
			wg.Add(1)
			sem <- struct{}{}
			go func() {
				defer wg.Done()
				defer func() { <-sem }()
				sd := uint64(9_000_000 + n*1000 + i)
				rng := rand.New(rand.NewPCG(sd, sd^0x9e3779b97f4a7c15)) //nolint:gosec // deterministic library
				g := snakesim.NewGame(n, rng)
				t := snakesim.NewTeacher(n)
				var snaps []snakesim.Snapshot
				next := 0.02
				for k := 0; k < 400*n*n && !g.Dead && float64(g.Len()) < *maxFill*float64(n*n); k++ {
					g.Step(t.Action(g))
					if float64(g.Len()) >= next*float64(n*n) && !g.Dead {
						snaps = append(snaps, g.Snapshot())
						next += 0.02
					}
				}
				mu.Lock()
				lib[n] = append(lib[n], snaps...)
				mu.Unlock()
			}()
		}
	}
	wg.Wait()
	data, _ := json.Marshal(lib)
	if err := os.WriteFile(*out, data, 0o644); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	for n, s := range lib {
		fmt.Printf("n=%d: %d snapshots\n", n, len(s))
	}
}
