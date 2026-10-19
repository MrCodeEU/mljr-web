package snakeai

import (
	"encoding/json"
	"math/rand/v2"
	"sync"
	"time"
)

// Config tunes a Room. Zero values get defaults.
type Config struct {
	N            int           // board size, must be even (default 30)
	Tick         time.Duration // time between moves (default 80ms)
	Grace        time.Duration // keep playing this long after the last viewer leaves (default 20s)
	RestartDelay time.Duration // pause on the final frame before a new game (default 3s)
	MaxViewers   int           // 0 = unlimited
}

func (c *Config) defaults() {
	if c.N == 0 {
		c.N = 30
	}
	if c.Tick == 0 {
		c.Tick = 80 * time.Millisecond
	}
	if c.Grace == 0 {
		c.Grace = 20 * time.Second
	}
	if c.RestartDelay == 0 {
		c.RestartDelay = 3 * time.Second
	}
}

// Room is the single shared game: one server-side game state, many viewers.
// The game and the model exist only while at least one viewer is connected
// (plus a grace period), so an idle site costs nothing.
type Room struct {
	cfg Config

	mu      sync.Mutex
	subs    map[*Subscription]struct{}
	running bool
	idleAt  time.Time
	model   *Model // non-nil only while running
	best    int
	games   int
	frames  uint64
}

// NewRoom creates an idle room.
func NewRoom(cfg Config) *Room {
	cfg.defaults()
	return &Room{cfg: cfg, subs: map[*Subscription]struct{}{}}
}

// Subscription receives pre-encoded SSE messages. Slow readers miss frames, never block the game.
type Subscription struct {
	C    <-chan []byte
	c    chan []byte
	room *Room
	once sync.Once
}

// Close unsubscribes.
func (s *Subscription) Close() {
	s.once.Do(func() {
		s.room.mu.Lock()
		delete(s.room.subs, s)
		if len(s.room.subs) == 0 {
			s.room.idleAt = time.Now()
		}
		s.room.mu.Unlock()
	})
}

// Subscribe joins the room, starting the game if it is idle. Returns nil when the room is full.
func (r *Room) Subscribe() *Subscription {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.cfg.MaxViewers > 0 && len(r.subs) >= r.cfg.MaxViewers {
		return nil
	}
	ch := make(chan []byte, 2)
	s := &Subscription{C: ch, c: ch, room: r}
	r.subs[s] = struct{}{}
	if !r.running {
		r.running = true
		go r.loop()
	}
	return s
}

// Running reports whether a game (and the model) is currently loaded.
func (r *Room) Running() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.running
}

// Viewers is the current subscriber count.
func (r *Room) Viewers() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.subs)
}

type frame struct {
	N       int    `json:"n"`
	Cells   []int  `json:"cells"` // x,y pairs, head first
	Food    [2]int `json:"food"`
	Score   int    `json:"score"`
	Best    int    `json:"best"`
	Games   int    `json:"games"`
	Viewers int    `json:"viewers"`
	Status  string `json:"status"` // playing | dead | complete
	Tick    uint64 `json:"tick"`
}

// encode wraps a frame as a Datastar patch-signals SSE message.
func encode(f frame) []byte {
	payload, err := json.Marshal(map[string]frame{"snake": f})
	if err != nil {
		return nil // unreachable for this struct; broadcast skips empty messages
	}
	return append(append([]byte("event: datastar-patch-signals\ndata: signals "), payload...), '\n', '\n')
}

func (r *Room) broadcast(msg []byte) {
	if len(msg) == 0 {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	for s := range r.subs {
		select {
		case s.c <- msg:
		default: // viewer is behind: drop this frame
			select {
			case <-s.c:
			default:
			}
			select {
			case s.c <- msg:
			default:
			}
		}
	}
}

func (r *Room) frameOf(g *Game, status string) frame {
	cells := make([]int, 0, 2*len(g.Snake))
	for _, p := range g.Snake {
		cells = append(cells, p.X, p.Y)
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.frames++
	f := frame{N: g.N, Cells: cells, Score: g.Score, Best: r.best, Games: r.games,
		Viewers: len(r.subs), Status: status, Tick: r.frames}
	if g.HasFood {
		f.Food = [2]int{g.Food.X, g.Food.Y}
	}
	return f
}

// stopIfIdle ends the loop when nobody has watched for the grace period.
func (r *Room) stopIfIdle() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(r.subs) == 0 && time.Since(r.idleAt) >= r.cfg.Grace {
		r.running = false
		r.model = nil
		return true
	}
	return false
}

func (r *Room) loop() {
	model, err := DefaultModel()
	if err != nil {
		r.mu.Lock()
		r.running = false
		r.mu.Unlock()
		return
	}
	r.mu.Lock()
	r.model = model
	r.idleAt = time.Now()
	r.mu.Unlock()

	rng := rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), 0x5eed)) //nolint:gosec // game randomness, not security-sensitive
	ticker := time.NewTicker(r.cfg.Tick)
	defer ticker.Stop()
	for {
		g := NewGame(r.cfg.N, rng)
		for !g.Dead {
			<-ticker.C
			if r.stopIfIdle() {
				return
			}
			g.Step(model.Act(g.Features()))
			if !g.Dead {
				r.broadcast(encode(r.frameOf(g, "playing")))
			}
		}
		status := "dead"
		if g.Full() {
			status = "complete"
		}
		r.mu.Lock()
		r.games++
		if g.Score > r.best {
			r.best = g.Score
		}
		r.mu.Unlock()
		r.broadcast(encode(r.frameOf(g, status)))
		for end := time.Now().Add(r.cfg.RestartDelay); time.Now().Before(end); {
			<-ticker.C
			if r.stopIfIdle() {
				return
			}
		}
	}
}
