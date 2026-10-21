"""DAgger driver for the v4 policy: Go generates teacher-labelled data (and plays the net), torch fits.

  uv run python -m snakenet.dagger_go --sim /path/to/snakesim --out out50 --hours 4

Round 0 imitates the teacher; later rounds let the net drive (teacher still labels every state), pool all data
and keep fitting. Each round the Go bench plays the net alone and the checks are tracked:
  (a) share of games that reach 25% fill without dying
  (b) steps per apple over the first 150 apples vs the teacher
  (c) median score at the step cap vs the teacher
The best model by these checks is kept as <out>/model_best.json.
"""
import argparse
import json
import re
import subprocess
import time
from pathlib import Path

import numpy as np
import torch
import torch.nn as nn

N_FEATURES = 105
FEATURE_VERSION = 5
REC = np.dtype([("y", "u1"), ("rule", "u1"), ("x", "<f4", (N_FEATURES,))])


def make_model(hidden):
    return nn.Sequential(nn.Linear(N_FEATURES, hidden), nn.ReLU(), nn.Linear(hidden, hidden), nn.ReLU(),
                         nn.Linear(hidden, 3))


def export(model, path):
    layers = [{"w": m.weight.detach().numpy().tolist(), "b": m.bias.detach().numpy().tolist()}
              for m in model if isinstance(m, nn.Linear)]
    json.dump({"feature_version": FEATURE_VERSION, "n_features": N_FEATURES, "actions": ["left", "straight", "right"],
               "activation": "relu", "layers": layers}, open(path, "w"))


def run(cmd):
    return subprocess.run(cmd, check=True, capture_output=True, text=True).stdout


def bench(sim, n, games, steps, model=None):
    cmd = [sim, "bench", "-n", str(n), "-games", str(games), "-steps", str(steps)]
    if model:
        cmd += ["-model", str(model)]
    out = run(cmd)
    m = re.search(r"median score=(\d+).*opening deaths \(<25% fill\)=(\d+)/(\d+)\s+reached 25%=(\d+)/(\d+)\s+died=(\d+)/(\d+)", out)
    spa = re.search(r"first 150 apples=([\d.]+)", out)
    return {"median": int(m[1]), "open_deaths": int(m[2]), "reached25": int(m[4]), "games": int(m[3]),
            "died": int(m[6]), "spa150": float(spa[1]) if spa else float("nan")}


def checks(net, teacher):
    games = net["games"]
    a = net["reached25"] / games
    b = net["spa150"] / teacher["spa150"] if teacher["spa150"] else float("nan")
    c = net["median"] / max(teacher["median"], 1)
    return a, b, c, (a >= 0.95 and b <= 1.10 and c >= 0.60)


def fit(model, X, Y, W, epochs, lr, batch=1024):
    X, Y, W = torch.from_numpy(X), torch.from_numpy(Y), torch.from_numpy(W)
    opt = torch.optim.Adam(model.parameters(), lr=lr)
    sched = torch.optim.lr_scheduler.CosineAnnealingLR(opt, epochs, eta_min=lr / 20)
    loss_fn = nn.CrossEntropyLoss(reduction="none")
    loss_v = 0.0
    for _ in range(epochs):
        perm = torch.randperm(len(X))
        tot = 0.0
        for i in range(0, len(X), batch):
            idx = perm[i:i + batch]
            opt.zero_grad()
            loss = (loss_fn(model(X[idx]), Y[idx]) * W[idx]).sum() / W[idx].sum()
            loss.backward()
            opt.step()
            tot += loss.item() * len(idx)
        sched.step()
        loss_v = tot / len(X)
    with torch.no_grad():
        acc = sum((model(X[i:i + 100000]).argmax(1) == Y[i:i + 100000]).sum().item() for i in range(0, len(X), 100000)) / len(X)
    return loss_v, acc


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("--sim", required=True, help="path to the snakesim binary")
    ap.add_argument("--out", default="out50")
    ap.add_argument("--lib", default=None, help="warm-start library (built if missing)")
    ap.add_argument("--hidden", type=int, default=128)
    ap.add_argument("--hours", type=float, default=4.0)
    ap.add_argument("--patience", type=int, default=3)
    ap.add_argument("--games", type=int, default=600, help="games per round")
    ap.add_argument("--epochs0", type=int, default=14)
    ap.add_argument("--epochs", type=int, default=6)
    ap.add_argument("--lr", type=float, default=2e-3)
    ap.add_argument("--eval-games", type=int, default=24)
    ap.add_argument("--eval-steps", type=int, default=40000)
    ap.add_argument("--max-records", type=int, default=7_000_000)
    ap.add_argument("--chase-weight", type=float, default=3.0, help="loss weight of tail-chase states (teacher rule 3)")
    ap.add_argument("--disagree-weight", type=float, default=5.0, help="loss weight where the net disagreed with the teacher")
    ap.add_argument("--predeath-weight", type=float, default=5.0, help="loss weight for the 60 moves before a death")
    ap.add_argument("--init", default=None, help="warm-start weights (state_dict .pt)")
    a = ap.parse_args()
    out = Path(a.out)
    out.mkdir(exist_ok=True)
    lib = Path(a.lib) if a.lib else out / "warmlib.json"
    t_end = time.time() + a.hours * 3600
    if not lib.exists():
        print("building warm-start library ...", flush=True)
        run([a.sim, "warmlib", "-games", "16", "-out", str(lib)])

    teacher = bench(a.sim, 50, a.eval_games, a.eval_steps)
    print(f"teacher baseline 50x50: {teacher}", flush=True)

    model = make_model(a.hidden)
    if a.init:
        model.load_state_dict(torch.load(a.init, weights_only=True))
    data_x, data_y, data_w = [], [], []
    best_score, best_round, bad = -1.0, -1, 0
    rnd = 0
    while time.time() < t_end:
        t0 = time.time()
        cur = out / "model_current.json"
        betas = [1.0, 0.5, 0.2, 0.0, 0.0, 0.0]
        beta = betas[min(rnd, len(betas) - 1)]
        f = out / f"data_{rnd}.bin"
        cmd = [a.sim, "gen", "-games", str(a.games), "-seed", str(rnd + 1), "-lib", str(lib), "-out", str(f)]
        if rnd > 0:
            cmd += ["-model", str(cur), "-beta", str(beta)]
        gen_out = run(cmd).strip()
        rec = np.fromfile(f, dtype=REC)
        f.unlink()
        data_x.append(rec["x"].copy())
        data_y.append(rec["y"].astype(np.int64))
        rule, flags = rec["rule"] & 15, rec["rule"]
        w = np.where(rule == 3, a.chase_weight, 1.0)
        w = w * np.where(flags & 16, a.disagree_weight, 1.0) * np.where(flags & 32, a.predeath_weight, 1.0)
        data_w.append(w.astype(np.float32))
        while sum(len(d) for d in data_y) > a.max_records and len(data_y) > 1:
            data_x.pop(0)
            data_y.pop(0)
            data_w.pop(0)
        X, Y, Wt = np.concatenate(data_x), np.concatenate(data_y), np.concatenate(data_w)
        t1 = time.time()
        loss, acc = fit(model, X, Y, Wt, a.epochs0 if rnd == 0 else a.epochs, a.lr if rnd == 0 else a.lr / 2)
        export(model, cur)
        t2 = time.time()
        net = bench(a.sim, 50, a.eval_games, a.eval_steps, cur)
        ca, cb, cc, ok = checks(net, teacher)
        score = (ca >= 0.95) * 1.0 + min(cc, 1.0) + max(0.0, 2.0 - cb)  # higher is better
        tag = ""
        if score > best_score:
            best_score, best_round, bad = score, rnd, 0
            export(model, out / "model_best.json")
            torch.save(model.state_dict(), out / "model_best.pt")
            tag = " *best*"
        else:
            bad += 1
        print(f"round {rnd}: beta={beta} data={len(Y)} loss={loss:.4f} acc={acc:.4f} | {net} | "
              f"(a) reach25={ca:.2f} (b) spa/teacher={cb:.2f} (c) score/teacher={cc:.2f} pass={ok}{tag} | "
              f"gen {t1 - t0:.0f}s fit {t2 - t1:.0f}s eval {time.time() - t2:.0f}s", flush=True)
        torch.save(model.state_dict(), out / "model_current.pt")
        if bad >= a.patience:
            print(f"no improvement for {a.patience} rounds, stopping", flush=True)
            break
        rnd += 1
    print(f"best round {best_round}, saved {out / 'model_best.json'}", flush=True)


if __name__ == "__main__":
    main()
