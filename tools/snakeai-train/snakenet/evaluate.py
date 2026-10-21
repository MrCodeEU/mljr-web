"""Net vs teacher at a fixed step budget, reporting score and deaths.

  uv run python -m snakenet.evaluate --size 20 --games 16 --max-steps 8000
A game that ends before the budget is a death (or a full board: score == n*n-1).
"""
import argparse
from multiprocessing import Pool

import numpy as np
import torch

from .rollout import _job
from .train import layers_of, make_model


def run(pool, n, games, layers, max_steps, teacher):
    jobs = [(n, 777_000 + i, layers, 0.0, 0.0, max_steps, False, teacher) for i in range(games)]
    res = pool.map(_job, jobs)
    full = n * n - 1
    deaths = sum(1 for r in res if r[1] < max_steps and r[0] < full)
    fills = sum(1 for r in res if r[0] >= full)
    return [r[0] for r in res], deaths, fills


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("--size", type=int, default=20)
    ap.add_argument("--games", type=int, default=20)
    ap.add_argument("--max-steps", type=int, default=8000)
    ap.add_argument("--model", default="out/model.pt")
    ap.add_argument("--teacher", default="cycle", choices=["bfs", "cycle"])
    a = ap.parse_args()
    model = make_model()
    model.load_state_dict(torch.load(a.model, weights_only=True))
    with Pool() as pool:
        for name, layers in (("net", layers_of(model)), (a.teacher, None)):
            sc, deaths, fills = run(pool, a.size, a.games, layers, a.max_steps, a.teacher)
            print(f"size={a.size} budget={a.max_steps} {name:7s} mean={np.mean(sc):.1f} min={min(sc)} "
                  f"max={max(sc)} deaths={deaths}/{a.games} filled={fills}/{a.games}")


if __name__ == "__main__":
    main()
