package snakesim

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

// decode parses one SSE message into its snake state.
func decode(t *testing.T, msg []byte) state {
	t.Helper()
	const prefix = "event: datastar-patch-signals\ndata: signals "
	if !bytes.HasPrefix(msg, []byte(prefix)) {
		t.Fatalf("unexpected message %q", msg)
	}
	var wrap map[string]state
	if err := json.Unmarshal([]byte(strings.TrimSpace(strings.TrimPrefix(string(msg), prefix))), &wrap); err != nil {
		t.Fatal(err)
	}
	return wrap["snake"]
}

// client mirrors the renderer.js logic: rebuild on a new full state, otherwise apply the moves newer
// than the last tick it saw.
type client struct {
	body           []int
	lastFull, tick uint64
}

// apply feeds one message in and returns the body the deltas alone would give for a full state
// (nil when it was not a full state or the client had nothing to compare against).
func (c *client) apply(t *testing.T, s state) []int {
	t.Helper()
	var predicted []int
	if s.Full != 0 && s.Full != c.lastFull {
		if c.body != nil && s.Tick == c.tick+1 && s.Status == "playing" {
			predicted = append([]int(nil), c.body...)
			applyMoves(&predicted, s, c.tick)
		}
		c.body, c.lastFull, c.tick = append([]int(nil), s.Body...), s.Full, s.Tick
		return predicted
	}
	first := s.Tick - uint64(len(s.Moves)/2) + 1
	if first > c.tick+1 {
		t.Fatalf("gap: tick %d covers from %d but the client is at %d", s.Tick, first, c.tick)
	}
	applyMoves(&c.body, s, c.tick)
	c.tick = s.Tick
	return nil
}

func applyMoves(body *[]int, s state, lastTick uint64) {
	k := len(s.Moves) / 2
	first := s.Tick - uint64(k) + 1
	for i := 0; i < k; i++ {
		if first+uint64(i) <= lastTick {
			continue
		}
		*body = append([]int{s.Moves[2*i]}, *body...)
		if s.Moves[2*i+1] == 0 {
			*body = (*body)[:len(*body)-1]
		}
	}
}

// Every periodic full state must equal the body obtained by applying the deltas before it.
func TestFullStatesMatchDeltaReconstruction(t *testing.T) {
	r := NewRoom(Config{N: 20, Speed: 400, Seed: 1, FullEvery: 25, Grace: time.Hour})
	sub := r.Subscribe()
	defer sub.Close()
	var c client
	compared := 0
	deadline := time.After(20 * time.Second)
	for compared < 6 {
		select {
		case <-deadline:
			t.Fatalf("timed out after %d comparisons", compared)
		case msg := <-sub.C:
			s := decode(t, msg)
			if predicted := c.apply(t, s); predicted != nil {
				if !equal(predicted, s.Body) {
					t.Fatalf("tick %d: deltas give %v, full state says %v", s.Tick, predicted, s.Body)
				}
				compared++
			}
		}
	}
}

// A client that misses messages (slow reader, batched updates) catches up from the repeated moves.
func TestClientCatchesUpAfterMissedMessages(t *testing.T) {
	r := NewRoom(Config{N: 20, Speed: 400, Seed: 5, FullEvery: 1000, Grace: time.Hour})
	sub := r.Subscribe()
	defer sub.Close()
	var msgs []state
	for len(msgs) < 120 {
		select {
		case m := <-sub.C:
			msgs = append(msgs, decode(t, m))
		case <-time.After(5 * time.Second):
			t.Fatal("timed out collecting messages")
		}
	}
	var full, skip, ref client
	for _, s := range msgs {
		ref.apply(t, s) // sees every message
	}
	full.apply(t, msgs[0]) // the first message is the full state
	skip.apply(t, msgs[0])
	for i, s := range msgs[1:] {
		full.apply(t, s)
		if i%3 == 2 { // skip sees only every third message
			skip.apply(t, s)
		}
	}
	last := msgs[len(msgs)-1]
	if (len(msgs)-1)%3 != 0 {
		skip.apply(t, last)
	}
	if !equal(full.body, ref.body) || !equal(skip.body, ref.body) {
		t.Fatalf("clients disagree: ref %v full %v skip %v", ref.body, full.body, skip.body)
	}
}

func equal(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestRoomLifecycle(t *testing.T) {
	r := NewRoom(Config{N: 20, Speed: 400, Grace: 50 * time.Millisecond, MaxViewers: 2})
	if r.Running() {
		t.Fatal("room should be idle before the first viewer")
	}
	a, b := r.Subscribe(), r.Subscribe()
	if a == nil || b == nil || r.Subscribe() != nil {
		t.Fatal("expected two viewers and a refusal for the third")
	}
	if !r.Running() || r.Viewers() != 2 {
		t.Fatal("room should be running with 2 viewers")
	}
	a.Close()
	b.Close()
	b.Close() // idempotent
	deadline := time.Now().Add(3 * time.Second)
	for r.Running() && time.Now().Before(deadline) {
		time.Sleep(20 * time.Millisecond)
	}
	if r.Running() {
		t.Fatal("room should stop after the grace period")
	}
}

func TestNewViewerGetsFullStateQuickly(t *testing.T) {
	r := NewRoom(Config{N: 20, Speed: 200, Seed: 2, Grace: time.Hour})
	first := r.Subscribe()
	defer first.Close()
	time.Sleep(300 * time.Millisecond) // game is mid-play
	late := r.Subscribe()
	defer late.Close()
	deadline := time.After(2 * time.Second)
	for {
		select {
		case <-deadline:
			t.Fatal("late viewer never received a full state")
		case msg := <-late.C:
			if s := decode(t, msg); s.Body != nil && s.N == 20 && s.Status == "playing" {
				return
			}
		}
	}
}
