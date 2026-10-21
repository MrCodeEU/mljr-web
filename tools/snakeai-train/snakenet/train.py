"""DAgger training loop.

  uv run python -m snakenet.train --rounds 6 --games 200
Round 0: teacher drives (with noise). Later rounds: net drives, teacher labels
every state it visits. Dataset aggregates across rounds.
"""
import argparse
import json
import time
from multiprocessing import Pool
from pathlib import Path

import numpy as np
import torch
import torch.nn as nn

from .env import N_FEATURES
from .rollout import _job

SIZES = (10, 14, 20, 20, 30, 40, 50)
HIDDEN = 128


def make_model():
    return nn.Sequential(nn.Linear(N_FEATURES, HIDDEN), nn.ReLU(),
                         nn.Linear(HIDDEN, HIDDEN), nn.ReLU(),
                         nn.Linear(HIDDEN, 3))


def layers_of(model):
    return [(m.weight.detach().numpy().tolist(), m.bias.detach().numpy().tolist())
            for m in model if isinstance(m, nn.Linear)]


def collect(pool, model, games, seed0, beta, eps, max_steps, teacher, warm=True):
    layers = layers_of(model) if model is not None else None
    jobs = [(SIZES[i % len(SIZES)], seed0 + i, layers, beta, eps, max_steps, True, teacher, warm and i % 4 != 0) for i in range(games)]
    res = pool.map(_job, jobs)
    scores = [r[0] for r in res]
    return np.concatenate([r[2] for r in res if len(r[2])]), np.concatenate([r[3] for r in res if len(r[3])]), scores


def fit(model, X, Y, epochs, lr=2e-3):
    X, Y = torch.from_numpy(X), torch.from_numpy(Y)
    opt = torch.optim.Adam(model.parameters(), lr=lr)
    sched = torch.optim.lr_scheduler.CosineAnnealingLR(opt, epochs, eta_min=lr / 20)
    loss_fn = nn.CrossEntropyLoss()
    for ep in range(epochs):
        perm = torch.randperm(len(X))
        tot = 0.0
        for i in range(0, len(X), 512):
            idx = perm[i:i + 512]
            opt.zero_grad()
            loss = loss_fn(model(X[idx]), Y[idx])
            loss.backward()
            opt.step()
            tot += loss.item() * len(idx)
        sched.step()
    acc = (model(X).argmax(1) == Y).float().mean().item()
    return tot / len(X), acc


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("--rounds", type=int, default=6)
    ap.add_argument("--games", type=int, default=200)
    ap.add_argument("--epochs", type=int, default=20)
    ap.add_argument("--max-steps", type=int, default=2500)
    ap.add_argument("--teacher", default="cycle", choices=["bfs", "cycle"])
    ap.add_argument("--init", default=None, help="warm-start weights (state_dict)")
    ap.add_argument("--out", default="out")
    a = ap.parse_args()
    out = Path(a.out)
    out.mkdir(exist_ok=True)
    torch.manual_seed(0)
    model = make_model()
    if a.init:
        model.load_state_dict(torch.load(a.init, weights_only=True))
    X = Y = None
    with Pool() as pool:
        for r in range(a.rounds):
            t0 = time.time()
            seed0 = 10_000 * (r + 1)
            if r == 0 and not a.init:
                x, y, scores = collect(pool, None, a.games, seed0, 1.0, 0.05, a.max_steps, a.teacher)
            else:
                x, y, scores = collect(pool, model, a.games, seed0, 0.0, 0.0, a.max_steps, a.teacher)
            X = x if X is None else np.concatenate([X, x])
            Y = y if Y is None else np.concatenate([Y, y])
            loss, acc = fit(model, X, Y, a.epochs)
            who = "teacher" if (r == 0 and not a.init) else "net"
            print(f"round {r}: {who} mean score={np.mean(scores):.1f} | data={len(X)} loss={loss:.3f} "
                  f"acc={acc:.3f} | {time.time() - t0:.0f}s", flush=True)
            torch.save(model.state_dict(), out / "model.pt")
    json.dump({"version": 1}, open(out / "train_meta.json", "w"))


if __name__ == "__main__":
    main()
