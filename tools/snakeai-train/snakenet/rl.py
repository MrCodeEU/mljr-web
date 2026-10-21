"""Masked PPO for speed, warm-started from an imitation net.

  python -m snakenet.rl --init checkpoints/model_v3.pt --mask cycle --iters 80
  python -m snakenet.rl --init checkpoints/model_v3.pt --mask tail  --iters 80

Why it differs from ppo.py (which stalled at the teacher):
  * Safety mask: unsafe actions get probability 0 during training *and* evaluation, so
    exploration cannot kill the snake and PPO only has to learn speed.
      cycle: only moves that keep the body in Hamiltonian-cycle order (cycle step or a
             shortcut that leaves a margin to the tail). Provably safe, but the snake
             always circulates in the cycle's direction.
      tail : any free move after which the head can still reach the tail. Lets the
             snake leave the loop; safe in practice, not provably.
  * Exploration: policy = softmax(logits / T), with T calibrated so the masked policy
    has real entropy (the imitation net is near-deterministic otherwise).
  * Reward: +1 per apple, potential-based shaping, death penalty; no teacher anchor.
  * Eval reports score at a fixed step budget relative to the cycle teacher, deaths, and
    how often the *unmasked* argmax already respected the mask (is the mask redundant?).
"""
import argparse
import time
from multiprocessing import Pool
from pathlib import Path

import numpy as np
import torch
import torch.nn as nn

from .env import N_FEATURES, Game, cyc_dist, cycle_for, features
from .policy import NumpyMLP
from .rollout import _job
from .teacher import ACTIONS, cycle_teacher_action, teacher_action
from .train import HIDDEN, layers_of, make_model

GAMMA, LAM = 0.995, 0.95
R_EAT, R_DEATH, SHAPE = 1.0, -5.0, 0.5
BUDGET = {10: 800, 14: 2000, 20: 5000}
NEG = -1e9


def mask_cycle(game, x):
    cells, index = cycle_for(game.n)
    total = len(cells)
    body = list(game.snake)
    hi = index[body[0]]
    to_tail = cyc_dist(hi, index[body[-1]], total)
    m = []
    for a in ACTIONS:
        c = game.cell(a)
        ok = False
        if c in index and not game.blocked(c):
            d = cyc_dist(hi, index[c], total)
            ok = d == 1 or 1 < d <= to_tail - 2
        m.append(ok)
    return m


def mask_tail(game, x):
    return [bool(x[i * 6] == 0 and x[i * 6 + 2] == 1) for i in range(3)]


MASKS = {"cycle": mask_cycle, "tail": mask_tail}


def masked_probs(logits, mask, temp):
    z = np.where(mask, logits / temp, NEG)
    z = z - z.max()
    p = np.exp(z)
    return p / p.sum()


def _episode(args):
    n, seed, layers, temp, horizon, warm, mask_mode = args
    net = NumpyMLP(layers)
    rng = np.random.default_rng(seed)
    game = Game(n, rng)
    mfn = MASKS[mask_mode]
    if warm:  # fast-forward with the safe teacher to a random fill level
        target = int(rng.uniform(0, 0.9) * n * n)
        for _ in range(60 * n * n):
            if game.dead or len(game.snake) >= target:
                break
            game.step(cycle_teacher_action(game))
        if game.dead or game.food is None:
            game = Game(n, rng)
    start_score = game.score
    X, M, A, LP, R, W = [], [], [], [], [], []
    for _ in range(horizon):
        if game.dead:
            break
        x = features(game)
        m = mfn(game, x)
        head, food = game.snake[0], game.food
        d0 = abs(food[0] - head[0]) + abs(food[1] - head[1])
        score0 = game.score
        if any(m):
            p = masked_probs(net.logits(x), np.array(m), temp)
            a = int(rng.choice(3, p=p))
            lp, w, mm = float(np.log(p[a] + 1e-9)), 1.0, m
        else:  # no safe action (rare, e.g. the very first move): forced teacher move, not trained on
            a = teacher_action(game)
            lp, w, mm = 0.0, 0.0, [True, True, True]
        game.step(a)
        if game.dead and game.food is not None:
            r = R_DEATH
        elif game.score > score0:
            r = R_EAT
        else:
            nh = game.snake[0]
            d1 = abs(food[0] - nh[0]) + abs(food[1] - nh[1])
            r = SHAPE * (-GAMMA * d1 + d0) / (2 * n)
        X.append(x)
        M.append(mm)
        A.append(a)
        LP.append(lp)
        R.append(r)
        W.append(w)
    dead = game.dead
    final_x = features(game) if not dead else np.zeros(N_FEATURES, np.float32)
    return (np.array(X, np.float32), np.array(M, bool), np.array(A), np.array(LP, np.float32),
            np.array(R, np.float32), np.array(W, np.float32), final_x, dead, game.score - start_score, len(X))


def _calib_states(args):
    n, seed, layers, mask_mode, steps = args
    rng = np.random.default_rng(seed)
    game = Game(n, rng)
    mfn = MASKS[mask_mode]
    for _ in range(int(rng.integers(0, 30 * n * n))):
        if game.dead:
            break
        game.step(cycle_teacher_action(game))
    if game.dead or game.food is None:
        return np.zeros((0, N_FEATURES), np.float32), np.zeros((0, 3), bool)
    X, M = [], []
    for _ in range(steps):
        if game.dead:
            break
        x = features(game)
        m = mfn(game, x)
        X.append(x)
        M.append(m)
        game.step(cycle_teacher_action(game))
    return np.array(X, np.float32), np.array(M, bool)


def calibrate_temp(pool, layers, mask_mode, target):
    jobs = [(n, 900 + i, layers, mask_mode, 150) for i, n in enumerate((10, 14, 20) * 8)]
    res = pool.map(_calib_states, jobs)
    X = np.concatenate([r[0] for r in res if len(r[0])])
    M = np.concatenate([r[1] for r in res if len(r[1])])
    multi = M.sum(1) >= 2
    net = NumpyMLP(layers)
    logits = net.logits(X)
    print(f"calibration: {len(X)} states, {multi.mean():.1%} have >=2 safe actions", flush=True)
    for temp in (1, 2, 4, 8, 16, 32, 64):
        ent = []
        for z, m in zip(logits[multi], M[multi]):
            p = masked_probs(z, m, temp)
            ent.append(-(p * np.log(p + 1e-12)).sum())
        mean = float(np.mean(ent)) if ent else 0.0
        print(f"  T={temp:<3} mean entropy over multi-choice states = {mean:.3f}", flush=True)
        if mean >= target:
            return float(temp)
    return 64.0


def _eval_game(args):
    n, seed, layers, cap, mask_mode = args
    net = NumpyMLP(layers)
    rng = np.random.default_rng(seed)
    game = Game(n, rng)
    mfn = MASKS[mask_mode]
    steps = violations = decisions = 0
    while not game.dead and steps < cap:
        x = features(game)
        m = mfn(game, x)
        z = net.logits(x)
        if any(m):
            a = int(np.argmax(np.where(m, z, NEG)))
            if sum(m) >= 2:
                decisions += 1
                violations += 0 if m[int(np.argmax(z))] else 1
        else:
            a = teacher_action(game)
        game.step(a)
        steps += 1
    return game.score, steps, game.dead, violations, decisions


def evaluate(pool, model, mask_mode, games, baseline):
    layers = layers_of(model)
    out, fit = {}, []
    for n, cap in BUDGET.items():
        res = pool.map(_eval_game, [(n, 555_000 + i, layers, cap, mask_mode) for i in range(games)])
        score = float(np.mean([r[0] for r in res]))
        deaths = sum(1 for r in res if r[1] < cap and r[0] < n * n - 1)
        viol = sum(r[3] for r in res)
        dec = sum(r[4] for r in res)
        out[n] = {"score": round(score, 1), "teacher": baseline[n], "deaths": deaths,
                  "mask_violation": round(viol / max(dec, 1), 3)}
        fit.append(score / baseline[n] - 0.5 * deaths / games)
    return out, float(np.mean(fit))


def teacher_baseline(pool, games):
    base = {}
    for n, cap in BUDGET.items():
        res = pool.map(_job, [(n, 555_000 + i, None, 0.0, 0.0, cap, False, "cycle") for i in range(games)])
        base[n] = float(np.mean([r[0] for r in res]))
    return base


def gae(values, rewards, last_value, dead):
    adv = np.zeros_like(rewards)
    nxt, last = (0.0 if dead else last_value), 0.0
    for t in reversed(range(len(rewards))):
        delta = rewards[t] + GAMMA * nxt - values[t]
        last = delta + GAMMA * LAM * last
        adv[t] = last
        nxt = values[t]
    return adv


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("--init", default="checkpoints/model_v3.pt")
    ap.add_argument("--mask", default="cycle", choices=list(MASKS))
    ap.add_argument("--out", default="out")
    ap.add_argument("--tag", default=None)
    ap.add_argument("--iters", type=int, default=80)
    ap.add_argument("--episodes", type=int, default=48)
    ap.add_argument("--horizon", type=int, default=1200)
    ap.add_argument("--lr", type=float, default=2e-4)
    ap.add_argument("--ent", type=float, default=0.02)
    ap.add_argument("--target-entropy", type=float, default=0.25)
    ap.add_argument("--warmup", type=int, default=5, help="iterations training only the critic")
    ap.add_argument("--eval-every", type=int, default=5)
    ap.add_argument("--eval-games", type=int, default=16)
    ap.add_argument("--procs", type=int, default=8)
    a = ap.parse_args()
    tag = a.tag or a.mask
    out = Path(a.out)
    out.mkdir(exist_ok=True)

    actor = make_model()
    actor.load_state_dict(torch.load(a.init, weights_only=True))
    critic = nn.Sequential(nn.Linear(N_FEATURES, HIDDEN), nn.ReLU(), nn.Linear(HIDDEN, HIDDEN), nn.ReLU(),
                           nn.Linear(HIDDEN, 1))
    opt_a = torch.optim.Adam(actor.parameters(), lr=a.lr)
    opt_c = torch.optim.Adam(critic.parameters(), lr=1e-3)
    sizes = (10, 14, 20)

    with Pool(a.procs) as pool:
        temp = calibrate_temp(pool, layers_of(actor), a.mask, a.target_entropy)
        print(f"[{tag}] temperature T={temp}", flush=True)
        base = teacher_baseline(pool, a.eval_games)
        print(f"[{tag}] teacher baseline score at budget {BUDGET}: {base}", flush=True)
        ev, best = evaluate(pool, actor, a.mask, a.eval_games, base)
        print(f"[{tag}] init eval {ev} fitness={best:.4f}", flush=True)
        best_state = {k: v.clone() for k, v in actor.state_dict().items()}
        torch.save({"state": best_state, "T": temp, "mask": a.mask}, out / f"model_rl_{tag}.pt")
        bad = 0
        for it in range(a.iters):
            t0 = time.time()
            layers = layers_of(actor)
            jobs = [(sizes[(it * a.episodes + i) % len(sizes)], 2_000_000 + it * 1000 + i, layers, temp,
                     a.horizon, i % 4 != 0, a.mask) for i in range(a.episodes)]
            eps = pool.map(_episode, jobs)
            X = np.concatenate([e[0] for e in eps])
            M = torch.from_numpy(np.concatenate([e[1] for e in eps]))
            Xt = torch.from_numpy(X)
            with torch.no_grad():
                vals = critic(Xt).squeeze(1).numpy()
                lastv = [critic(torch.from_numpy(e[6])).item() if not e[7] else 0.0 for e in eps]
            advs, rets, pos = [], [], 0
            for e, lv in zip(eps, lastv):
                m = len(e[4])
                v = vals[pos:pos + m]
                ad = gae(v, e[4], lv, e[7])
                advs.append(ad)
                rets.append(ad + v)
                pos += m
            ADV = np.concatenate(advs)
            RET = torch.from_numpy(np.concatenate(rets))
            W = torch.from_numpy(np.concatenate([e[5] for e in eps]))
            ADV = torch.from_numpy((ADV - ADV.mean()) / (ADV.std() + 1e-8))
            ACT = torch.from_numpy(np.concatenate([e[3] * 0 + e[2] for e in eps]).astype(np.int64))
            OLP = torch.from_numpy(np.concatenate([e[3] for e in eps]))
            train_actor = it >= a.warmup
            v_loss = torch.tensor(0.0)
            ent = torch.tensor(0.0)
            for _ in range(4):
                perm = torch.randperm(len(Xt))
                for i in range(0, len(Xt), 2048):
                    idx = perm[i:i + 2048]
                    v_loss = ((critic(Xt[idx]).squeeze(1) - RET[idx]) ** 2).mean()
                    opt_c.zero_grad()
                    v_loss.backward()
                    opt_c.step()
                    if not train_actor:
                        continue
                    z = actor(Xt[idx]) / temp
                    z = z.masked_fill(~M[idx], NEG)
                    logp_all = torch.log_softmax(z, 1)
                    lp = logp_all.gather(1, ACT[idx].unsqueeze(1)).squeeze(1)
                    ratio = torch.exp(lp - OLP[idx])
                    surr = torch.min(ratio * ADV[idx], torch.clamp(ratio, 0.8, 1.2) * ADV[idx])
                    multi = (M[idx].sum(1) >= 2).float() * W[idx]
                    ent_each = -(logp_all.exp() * logp_all.clamp(min=-30)).sum(1)
                    ent = (ent_each * multi).sum() / multi.sum().clamp(min=1)
                    loss = -(surr * W[idx]).sum() / W[idx].sum().clamp(min=1) - a.ent * ent
                    opt_a.zero_grad()
                    loss.backward()
                    nn.utils.clip_grad_norm_(actor.parameters(), 0.5)
                    opt_a.step()
            eats = sum(e[8] for e in eps)
            steps = sum(e[9] for e in eps)
            msg = (f"[{tag}] iter {it}: steps={steps} apples/1k steps={1000 * eats / max(steps, 1):.2f} "
                   f"v_loss={v_loss.item():.3f} ent={float(ent.detach()):.3f} "
                   f"{'train' if train_actor else 'critic-warmup'} {time.time() - t0:.0f}s")
            if train_actor and (it + 1) % a.eval_every == 0:
                ev, fit = evaluate(pool, actor, a.mask, a.eval_games, base)
                msg += f" | eval {ev} fitness={fit:.4f}"
                if fit > best:
                    best, bad = fit, 0
                    best_state = {k: v.clone() for k, v in actor.state_dict().items()}
                    torch.save({"state": best_state, "T": temp, "mask": a.mask}, out / f"model_rl_{tag}.pt")
                    msg += " *best*"
                else:
                    bad += 1
                    if bad >= 2 and fit < 0.95 * best:
                        actor.load_state_dict(best_state)
                        a.lr *= 0.5
                        opt_a = torch.optim.Adam(actor.parameters(), lr=a.lr)
                        bad = 0
                        msg += " ROLLBACK, lr halved"
            print(msg, flush=True)
    print(f"[{tag}] best fitness {best:.4f} (1.0 = cycle teacher) saved to {out / f'model_rl_{tag}.pt'}")


if __name__ == "__main__":
    main()
