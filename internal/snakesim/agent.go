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

// Act returns the model's relative action for the current state.
func (a *Agent) Act(g *Game) int {
	a.t.Features(g, a.feat)
	return a.m.Act(a.feat, a.bufA, a.bufB)
}
