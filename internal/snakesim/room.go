package snakesim

import (
	"encoding/json"
	"math/rand/v2"
	"sync"
	"time"
)

// Config tunes a Room. Zero values get defaults.
type Config struct {
	N            int           // board size (default 50)
	Speed        float64       // moves per second, timed against the real clock (default 30)
	Grace        time.Duration // keep playing this long after the last viewer leaves (default 20s)
	RestartDelay time.Duration // pause on the final frame before a new game (default 3s)
	MaxViewers   int           // 0 = unlimited
	FullEvery    int           // resend the full state every this many moves (default 60)
	Seed         uint64        // game randomness; 0 = time based
}

func (c *Config) defaults() {
	if c.N == 0 {
		c.N = 50
	}
	if c.Speed == 0 {
		c.Speed = 30
	}
	if c.Grace == 0 {
		c.Grace = 20 * time.Second
	}
	if c.RestartDelay == 0 {
		c.RestartDelay = 3 * time.Second
	}
	if c.FullEvery == 0 {
		c.FullEvery = 60
	}
}

// Room is the single shared game: one server-side game state, many viewers.
// The game and the model exist only while at least one viewer is connected
// (plus a grace period), so an idle site costs nothing.
type Room struct {
	cfg Config

	mu       sync.Mutex
	subs     map[*Subscription]struct{}
	running  bool
	idleAt   time.Time
	needFull bool
	best     int
	games    int
	tick     uint64
	fullSeq  uint64
	moves    []int // last moves, head/grew pairs
	last     []int // test hook: body of the last broadcast state (head first, cell indices)
}

// NewRoom creates an idle room.
func NewRoom(cfg Config) *Room {
	cfg.defaults()
	return &Room{cfg: cfg, subs: map[*Subscription]struct{}{}}
}

// Subscription receives pre-encoded SSE messages. Slow readers miss messages, never block the game; the
// periodic full state resynchronises them.
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

// Subscribe joins the room, starting the game if it is idle, and asks for a full state on the next move.
// Returns nil when the room is full.
func (r *Room) Subscribe() *Subscription {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.cfg.MaxViewers > 0 && len(r.subs) >= r.cfg.MaxViewers {
		return nil
	}
	ch := make(chan []byte, 8)
	s := &Subscription{C: ch, c: ch, room: r}
	r.subs[s] = struct{}{}
	r.needFull = true
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

// state is the shared signal object. Full states carry the body; every state carries the last few moves,
// so a client that missed a message catches up without waiting for the next full state.
type state struct {
	N       int    `json:"n,omitempty"`
	Full    uint64 `json:"full,omitempty"` // increments with every full state; clients rebuild when it changes
	Body    []int  `json:"body,omitempty"` // cell indices y*n+x, head first
	Moves   []int  `json:"m,omitempty"`    // last moves as head cell, grew (0/1) pairs, oldest first; the newest is this tick
	Food    int    `json:"food"`
	Score   int    `json:"score"`
	Best    int    `json:"best"`
	Games   int    `json:"games"`
	Viewers int    `json:"viewers"`
	Status  string `json:"status"` // playing | dead | complete
	Tick    uint64 `json:"tick"`
}

// encode wraps a state as a Datastar patch-signals SSE message.
func encode(s state) []byte {
	payload, err := json.Marshal(map[string]state{"snake": s})
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
		default: // viewer is behind: drop its oldest message
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

// stopIfIdle ends the loop when nobody has watched for the grace period.
func (r *Room) stopIfIdle() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(r.subs) == 0 && time.Since(r.idleAt) >= r.cfg.Grace {
		r.running = false
		return true
	}
	return false
}

// movesKept is how many recent moves each delta repeats.
const movesKept = 4

// emit builds and broadcasts the message for this move. withBody forces a full state.
func (r *Room) emit(g *Game, ate bool, status string, withBody bool) {
	r.mu.Lock()
	r.tick++
	if r.needFull {
		withBody = true
		r.needFull = false
	}
	grew := 0
	if ate {
		grew = 1
	}
	r.moves = append(r.moves, g.Head(), grew)
	if len(r.moves) > 2*movesKept {
		r.moves = r.moves[len(r.moves)-2*movesKept:]
	}
	s := state{Food: g.Food, Score: g.Score, Best: r.best, Games: r.games,
		Viewers: len(r.subs), Status: status, Tick: r.tick}
	if withBody {
		r.fullSeq++
		s.N, s.Full = g.N, r.fullSeq
		s.Body = make([]int, g.Len())
		for i := range s.Body {
			s.Body[i] = g.Cell(i)
		}
		r.last = s.Body
	}
	s.Moves = append([]int(nil), r.moves...)
	r.mu.Unlock()
	r.broadcast(encode(s))
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
	r.idleAt = time.Now()
	seed := r.cfg.Seed
	r.mu.Unlock()
	if seed == 0 {
		seed = uint64(time.Now().UnixNano())
	}
	rng := rand.New(rand.NewPCG(seed, 0x5eed)) //nolint:gosec // game randomness, not security-sensitive
	agent := NewAgent(model, r.cfg.N)
	tick := time.Duration(float64(time.Second) / r.cfg.Speed)
	next := time.Now()
	wait := func() bool { // sleeps until the next move is due on the real clock; false when the room went idle
		next = next.Add(tick)
		if d := time.Until(next); d > 0 {
			time.Sleep(d)
		} else if d < -5*tick { // far behind (suspended process): skip ahead instead of bursting
			next = time.Now()
		}
		return !r.stopIfIdle()
	}
	for {
		g := NewGame(r.cfg.N, rng)
		r.emit(g, false, "playing", true)
		for moves := 1; !g.Dead; moves++ {
			if !wait() {
				return
			}
			before := g.Score
			g.Step(agent.Act(g))
			if !g.Dead {
				r.emit(g, g.Score > before, "playing", moves%r.cfg.FullEvery == 0)
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
		r.emit(g, false, status, true)
		for end := time.Now().Add(r.cfg.RestartDelay); time.Now().Before(end); {
			if !wait() {
				return
			}
		}
	}
}
