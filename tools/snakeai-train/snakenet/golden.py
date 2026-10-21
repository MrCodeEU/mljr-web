"""Write golden vectors for the Go forward pass: features from `snakesim gen` and logits from an independent
numpy implementation of the exported model.

  snakesim gen -games 12 -out /tmp/g.bin
  uv run python -m snakenet.golden --model out50/model_best.json --data /tmp/g.bin \
      --out ../../internal/snakesim/testdata/golden.json
"""
import argparse
import json

import numpy as np

from .dagger_go import FEATURE_VERSION, N_FEATURES, REC


def forward(layers, x):
    for i, l in enumerate(layers):
        x = x @ np.asarray(l["w"], np.float32).T + np.asarray(l["b"], np.float32)
        if i < len(layers) - 1:
            x = np.maximum(x, 0)
    return x


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("--model", required=True)
    ap.add_argument("--data", required=True)
    ap.add_argument("--out", required=True)
    ap.add_argument("--count", type=int, default=40)
    a = ap.parse_args()
    model = json.load(open(a.model))
    assert model["feature_version"] == FEATURE_VERSION and model["n_features"] == N_FEATURES
    rec = np.fromfile(a.data, dtype=REC)
    idx = np.random.default_rng(7).choice(len(rec), a.count, replace=False)
    x = rec["x"][idx]
    logits = forward(model["layers"], x)
    cases = [{"features": x[i].tolist(), "logits": logits[i].tolist(), "action": int(np.argmax(logits[i]))}
             for i in range(len(x))]
    json.dump({"feature_version": FEATURE_VERSION, "cases": cases}, open(a.out, "w"))
    print(f"wrote {len(cases)} golden cases to {a.out}")


if __name__ == "__main__":
    main()
