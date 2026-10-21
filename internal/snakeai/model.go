package snakeai

import (
	_ "embed"
	"encoding/json"
	"fmt"
)

//go:embed model.json
var embeddedModel []byte

type layer struct {
	W [][]float32 `json:"w"` // [out][in]
	B []float32   `json:"b"`
}

// Model is a small ReLU MLP: features -> 3 relative-action logits.
type Model struct {
	FeatureVersion int     `json:"feature_version"`
	NFeatures      int     `json:"n_features"`
	Layers         []layer `json:"layers"`
}

// LoadModel parses model.json bytes and checks the feature contract.
func LoadModel(data []byte) (*Model, error) {
	var m Model
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	if m.FeatureVersion != FeatureVersion || m.NFeatures != NFeatures {
		return nil, fmt.Errorf("snakeai: model expects feature v%d/%d, code has v%d/%d",
			m.FeatureVersion, m.NFeatures, FeatureVersion, NFeatures)
	}
	return &m, nil
}

// DefaultModel loads the embedded model.
func DefaultModel() (*Model, error) { return LoadModel(embeddedModel) }

// Logits runs the forward pass.
func (m *Model) Logits(x []float32) []float32 {
	for i, l := range m.Layers {
		y := make([]float32, len(l.B))
		for o := range y {
			s := l.B[o]
			w := l.W[o]
			for j, v := range x {
				s += w[j] * v
			}
			if i < len(m.Layers)-1 && s < 0 {
				s = 0
			}
			y[o] = s
		}
		x = y
	}
	return x
}

// Act picks the highest-logit relative action.
func (m *Model) Act(x []float32) int {
	z := m.Logits(x)
	best := 0
	for i := 1; i < len(z); i++ {
		if z[i] > z[best] {
			best = i
		}
	}
	return best
}
