package snakesim

import (
	"encoding/json"
	"fmt"
)

type layer struct {
	W [][]float32 `json:"w"` // [out][in]
	B []float32   `json:"b"`
}

// Model is a small ReLU MLP: NFeatures inputs -> 3 relative-action logits.
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
		return nil, fmt.Errorf("snakesim: model expects feature v%d/%d, code has v%d/%d",
			m.FeatureVersion, m.NFeatures, FeatureVersion, NFeatures)
	}
	return &m, nil
}

// Act picks the highest-logit relative action, using buf (len ≥ the widest layer) as scratch.
func (m *Model) Act(x []float32, bufA, bufB []float32) int {
	cur := x
	for i, l := range m.Layers {
		out := bufA[:len(l.B)]
		for o := range out {
			s := l.B[o]
			w := l.W[o]
			for j, v := range cur {
				s += w[j] * v
			}
			if i < len(m.Layers)-1 && s < 0 {
				s = 0
			}
			out[o] = s
		}
		cur = out
		bufA, bufB = bufB, bufA
	}
	best := 0
	for i := 1; i < len(cur); i++ {
		if cur[i] > cur[best] {
			best = i
		}
	}
	return best
}

// Logits runs the forward pass and returns the 3 action logits (allocating; use Act in loops).
func (m *Model) Logits(x []float32) []float32 {
	cur := x
	for i, l := range m.Layers {
		out := make([]float32, len(l.B))
		for o := range out {
			s := l.B[o]
			for j, v := range cur {
				s += l.W[o][j] * v
			}
			if i < len(m.Layers)-1 && s < 0 {
				s = 0
			}
			out[o] = s
		}
		cur = out
	}
	return cur
}

// Width is the widest layer output, for sizing scratch buffers.
func (m *Model) Width() int {
	w := 0
	for _, l := range m.Layers {
		if len(l.B) > w {
			w = len(l.B)
		}
	}
	return w
}

// Trace runs the forward pass and returns every layer's output: post-ReLU activations for the hidden layers,
// raw logits last. It allocates, so use it only where the values are shown (live visualisation).
func (m *Model) Trace(x []float32) [][]float32 {
	acts := make([][]float32, 0, len(m.Layers))
	cur := x
	for i, l := range m.Layers {
		out := make([]float32, len(l.B))
		for o := range out {
			s := l.B[o]
			for j, v := range cur {
				s += l.W[o][j] * v
			}
			if i < len(m.Layers)-1 && s < 0 {
				s = 0
			}
			out[o] = s
		}
		acts = append(acts, out)
		cur = out
	}
	return acts
}
