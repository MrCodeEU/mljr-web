"""Benchmark teachers alone: score at a fixed step budget, deaths, and steps per apple by length band.

  uv run python -m snakenet.teacher_bench --n 20 --games 24 --teachers cycle,bfs
"""
import argparse
import importlib
from multiprocessing import Pool

import numpy as np

from .env import Game

BANDS = (0.25, 0.5, 0.75, 1.01)


def get_teacher(name):
    if ":" in name:
        mod, fn = name.split(":")
        return getattr(importlib.import_module(mod), fn)
    from .rollout import TEACHERS
    return TEACHERS[name]


def run(args):
    name, n, seed, cap = args
    act = get_teacher(name)
    g = Game(n, np.random.default_rng(seed))
    band_steps = [0] * len(BANDS)
    band_apples = [0] * len(BANDS)
    steps = 0
    full = n * n - 1
    while not g.dead and steps < cap and g.score < full:
        frac = len(g.snake) / (n * n)
        b = next(i for i, edge in enumerate(BANDS) if frac < edge)
        before = g.score
        g.step(act(g))
        steps += 1
        band_steps[b] += 1
        band_apples[b] += g.score - before
    died = g.dead and g.score < full and steps < cap
    return g.score, steps, died, g.score >= full, band_steps, band_apples


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("--n", type=int, default=20)
    ap.add_argument("--games", type=int, default=24)
    ap.add_argument("--cap", type=int, default=None)
    ap.add_argument("--teachers", default="cycle,bfs")
    ap.add_argument("--seed", type=int, default=777_000)
    ap.add_argument("--procs", type=int, default=16)
    a = ap.parse_args()
    cap = a.cap or 40 * a.n * a.n
    with Pool(a.procs) as pool:
        for name in a.teachers.split(","):
            res = pool.map(run, [(name, a.n, a.seed + i, cap) for i in range(a.games)])
            bs = np.sum([r[4] for r in res], axis=0)
            ba = np.sum([r[5] for r in res], axis=0)
            per = " ".join(f"{bs[i] / max(ba[i], 1):.1f}" for i in range(len(BANDS)))
            fs = [r[1] for r in res if r[3]]
            fill = f"{np.mean(fs):.0f}" if fs else "-"
            print(f"n={a.n} {name:22} fill_steps={fill:>6} score={np.mean([r[0] for r in res]):6.1f} deaths={sum(r[2] for r in res)}/{a.games} "
                  f"filled={sum(r[3] for r in res)}/{a.games} steps/apple by length band "
                  f"[<25% <50% <75% rest]: {per}", flush=True)


if __name__ == "__main__":
    main()
