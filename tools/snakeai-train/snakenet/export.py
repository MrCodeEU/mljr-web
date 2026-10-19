"""Export trained net: model.json (pure-Go runtime), model.onnx, golden.json (Go parity tests)."""
import argparse
import json
from pathlib import Path

import numpy as np
import torch

from .env import FEATURE_VERSION, N_FEATURES, WIN, Game, features
from .policy import NumpyMLP
from .teacher import cycle_teacher_action as teacher_action
from .train import layers_of, make_model


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("--model", default="out/model.pt")
    ap.add_argument("--out", default="out")
    a = ap.parse_args()
    out = Path(a.out)
    model = make_model()
    model.load_state_dict(torch.load(a.model, weights_only=True))
    model.eval()
    layers = layers_of(model)

    json.dump({"feature_version": FEATURE_VERSION, "n_features": N_FEATURES, "window_radius": WIN,
               "actions": ["left", "straight", "right"], "activation": "relu",
               "layers": [{"w": w, "b": b} for w, b in layers]},
              open(out / "model.json", "w"))

    torch.onnx.export(model, torch.zeros(1, N_FEATURES), out / "model.onnx",
                      input_names=["features"], output_names=["logits"],
                      dynamic_axes={"features": {0: "batch"}})

    # golden states: raw game state -> features -> logits, so Go can test both
    # its feature extraction and its forward pass
    net = NumpyMLP(layers)
    rng = np.random.default_rng(42)
    golden = []
    for i, n in enumerate((10, 14, 20, 30, 40, 50) * 5):
        g = Game(n, rng)
        for _ in range(int(rng.integers(0, 40 * n * n // 4))):
            if g.dead:
                break
            g.step(teacher_action(g))
        if g.dead:
            continue
        x = features(g)
        golden.append({"n": n, "snake": [list(c) for c in g.snake], "heading": g.heading,
                       "food": list(g.food), "features": x.tolist(),
                       "logits": net.logits(x).tolist(), "action": net.act(x)})
    json.dump({"feature_version": FEATURE_VERSION, "cases": golden}, open(out / "golden.json", "w"))
    print(f"wrote model.json, model.onnx, golden.json ({len(golden)} cases) to {out}/")


if __name__ == "__main__":
    main()
