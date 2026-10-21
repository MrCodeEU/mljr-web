package snakesim

import _ "embed"

//go:embed model.json
var embeddedModel []byte

// DefaultModel loads the embedded model.
func DefaultModel() (*Model, error) { return LoadModel(embeddedModel) }

// Agent plays a game with a Model. It owns scratch buffers, so use one per game loop.
type Agent struct {
	m          *Model
	t          *Teacher
	feat       []float32
	bufA, bufB []float32
}

// NewAgent prepares an agent for boards up to n×n.
func NewAgent(m *Model, n int) *Agent {
	w := m.Width()
	return &Agent{m: m, t: NewTeacher(n), feat: make([]float32, NFeatures), bufA: make([]float32, w), bufB: make([]float32, w)}
}

// Act returns the guarded action (see GuardedAct) for the current state.
func (a *Agent) Act(g *Game) int {
	a.t.Features(g, a.feat)
	return GuardedAct(a.m, a.feat)
}

// GuardedAct is the model's best-ranked action that is not blocked and keeps the tail reachable after the
// move (features per action: [0] blocked, [2] tail reachable). If none qualifies it takes the best-ranked
// unblocked action, then the plain argmax. The net still chooses every move it can; the guard only vetoes a
// move that would trap the snake.
func GuardedAct(m *Model, feat []float32) int {
	logits := m.Logits(feat)
	var order [3]int
	for i := range order {
		order[i] = i
	}
	for i := 1; i < 3; i++ {
		for j := i; j > 0 && logits[order[j]] > logits[order[j-1]]; j-- {
			order[j], order[j-1] = order[j-1], order[j]
		}
	}
	open := -1
	for _, a := range order {
		o := a * perAction
		if feat[o] != 0 {
			continue
		}
		if feat[o+2] != 0 {
			return a
		}
		if open < 0 {
			open = a
		}
	}
	if open >= 0 {
		return open
	}
	return order[0]
}
