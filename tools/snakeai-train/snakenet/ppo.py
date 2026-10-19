"""PPO fine-tuning, warm-started from the imitation net.

  uv run python -m snakenet.ppo --init out/model.pt --iters 60

Guards against stalls/collapse: critic warm-up with a frozen actor, decaying
behaviour-cloning anchor to the teacher (labelled on-policy), entropy bonus,
potential-based shaping, starvation timeout, periodic fixed-seed eval with
best-checkpoint keeping and rollback on regression.
"""
import argparse
import time
from multiprocessing import Pool
from pathlib import Path

import numpy as np
import torch
import torch.nn as nn

from .env import N_FEATURES, Game, features
from .policy import NumpyMLP
from .rollout import _job
from .teacher import cycle_teacher_action as teacher_action
from .train import HIDDEN, layers_of, make_model

GAMMA, LAM = 0.99, 0.95
R_EAT, R_DEATH, SHAPE, R_STEP = 1.0, -5.0, 0.3, -0.002  # speed objective: time costs, death forfeits future food
BUDGET = {10: 800, 14: 2000, 20: 5000, 30: 12000}  # eval step budgets (below the cycle fill time)


def _episode(args):
    n, seed, layers, max_steps, warm = args
    net = NumpyMLP(layers)
    rng = np.random.default_rng(seed)
    game = Game(n, rng)
    if warm:  # start from a late-game state reached by the safe teacher
        for _ in range(int(rng.integers(0, 40 * n * n))):
            if game.dead:
                break
            game.step(teacher_action(game))
        if game.dead or game.food is None:
            game = Game(n, rng)
    X, A, LP, R, T = [], [], [], [], []
    steps, dead = 0, False
    while not game.dead and steps < max_steps:
        x = features(game)
        z = net.logits(x)
        z = z - z.max()
        p = np.exp(z)
        p /= p.sum()
        a = int(rng.choice(3, p=p))
        head, food = game.snake[0], game.food
        d0 = abs(food[0] - head[0]) + abs(food[1] - head[1])
        score0 = game.score
        T.append(teacher_action(game))
        game.step(a)
        r = 0.0
        if game.dead and not (game.food is None):
            r = R_DEATH
        elif game.score > score0:
            r = R_EAT  # food teleported: skip shaping across the jump
        else:
            nh = game.snake[0]
            d1 = abs(food[0] - nh[0]) + abs(food[1] - nh[1])
            # potential phi = -dist/(2n): F = gamma*phi(s') - phi(s)
            r = SHAPE * (-GAMMA * d1 + d0) / (2 * n)
        r += R_STEP
        X.append(x)
        A.append(a)
        LP.append(float(np.log(p[a] + 1e-9)))
        R.append(r)
        steps += 1
    dead = game.dead
    final_x = features(game) if not dead else np.zeros(N_FEATURES, np.float32)
    return (np.array(X, np.float32), np.array(A), np.array(LP, np.float32), np.array(R, np.float32),
            np.array(T), final_x, dead, game.score)


def gae(values, rewards, last_value, dead):
    adv = np.zeros_like(rewards)
    nxt, last = (0.0 if dead else last_value), 0.0
    for t in reversed(range(len(rewards))):
        delta = rewards[t] + GAMMA * nxt - values[t]
        last = delta + GAMMA * LAM * last
        adv[t] = last
        nxt = values[t]
    return adv


def evaluate(pool, model, sizes=(10, 20), games=24, teacher=False):
    layers = None if teacher else layers_of(model)
    out = {}
    for n in sizes:
        cap = BUDGET[n]
        res = pool.map(_job, [(n, 555_000 + i, layers, 0.0, 0.0, cap, False, "cycle") for i in range(games)])
        deaths = sum(1 for r in res if r[1] < cap and r[0] < n * n - 1)
        out[n] = (round(float(np.mean([r[0] for r in res])), 1), deaths)
    return out


def fitness(ev, games=24):
    """Mean board-fill fraction at the step budget, minus a penalty per death."""
    return sum(sc / (n * n) - 0.5 * d / games for n, (sc, d) in ev.items()) / len(ev)


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("--init", default="out/model.pt")
    ap.add_argument("--out", default="out")
    ap.add_argument("--iters", type=int, default=60)
    ap.add_argument("--episodes", type=int, default=64)
    ap.add_argument("--max-steps", type=int, default=1500)
    ap.add_argument("--lr", type=float, default=3e-4)
    ap.add_argument("--warmup", type=int, default=4, help="iterations training only the critic")
    ap.add_argument("--bc", type=float, default=0.5, help="initial behaviour-cloning weight, decays to 0.05")
    ap.add_argument("--bc-end", type=float, default=0.3)
    ap.add_argument("--sizes", default="10,14,20,20,30")
    ap.add_argument("--ent", type=float, default=0.01)
    ap.add_argument("--eval-every", type=int, default=4)
    a = ap.parse_args()
    out = Path(a.out)
    sizes = tuple(int(x) for x in a.sizes.split(","))

    actor = make_model()
    actor.load_state_dict(torch.load(a.init, weights_only=True))
    critic = nn.Sequential(nn.Linear(N_FEATURES, HIDDEN), nn.ReLU(), nn.Linear(HIDDEN, HIDDEN), nn.ReLU(),
                           nn.Linear(HIDDEN, 1))
    opt_a = torch.optim.Adam(actor.parameters(), lr=a.lr)
    opt_c = torch.optim.Adam(critic.parameters(), lr=1e-3)

    best, best_state, bad = -1.0, None, 0
    with Pool() as pool:
        print("teacher baseline", evaluate(pool, actor, teacher=True), flush=True)
        ev = evaluate(pool, actor)
        best = fitness(ev)
        best_state = {k: v.clone() for k, v in actor.state_dict().items()}
        print(f"init eval {ev} fitness={best:.4f}", flush=True)
        for it in range(a.iters):
            t0 = time.time()
            layers = layers_of(actor)
            jobs = [(sizes[(it * a.episodes + i) % len(sizes)], 1_000_000 + it * 1000 + i, layers, a.max_steps, i % 2 == 0)
                    for i in range(a.episodes)]
            eps = pool.map(_episode, jobs)

            X = np.concatenate([e[0] for e in eps])
            Xt = torch.from_numpy(X)
            with torch.no_grad():
                vals_all = critic(Xt).squeeze(1).numpy()
                lastv = [critic(torch.from_numpy(e[5])).item() if not e[6] else 0.0 for e in eps]
            advs, rets, pos = [], [], 0
            for e, lv in zip(eps, lastv):
                m = len(e[3])
                v = vals_all[pos:pos + m]
                ad = gae(v, e[3], lv, e[6])
                advs.append(ad)
                rets.append(ad + v)
                pos += m
            ADV = np.concatenate(advs)
            RET = torch.from_numpy(np.concatenate(rets))
            ADV = torch.from_numpy((ADV - ADV.mean()) / (ADV.std() + 1e-8))
            ACT = torch.from_numpy(np.concatenate([e[1] for e in eps]))
            OLP = torch.from_numpy(np.concatenate([e[2] for e in eps]))
            TEA = torch.from_numpy(np.concatenate([e[4] for e in eps]))

            frac = it / max(a.iters - 1, 1)
            bc_w = a.bc + (a.bc_end - a.bc) * frac
            train_actor = it >= a.warmup
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
                    logits = actor(Xt[idx])
                    logp_all = torch.log_softmax(logits, 1)
                    lp = logp_all.gather(1, ACT[idx].unsqueeze(1)).squeeze(1)
                    ratio = torch.exp(lp - OLP[idx])
                    pg = -torch.min(ratio * ADV[idx], torch.clamp(ratio, 0.8, 1.2) * ADV[idx]).mean()
                    ent = -(logp_all.exp() * logp_all).sum(1).mean()
                    bc = nn.functional.cross_entropy(logits, TEA[idx])
                    loss = pg - a.ent * ent + bc_w * bc
                    opt_a.zero_grad()
                    loss.backward()
                    nn.utils.clip_grad_norm_(actor.parameters(), 0.5)
                    opt_a.step()

            ep_scores = [e[7] for e in eps]
            msg = (f"iter {it}: steps={len(X)} mean_score={np.mean(ep_scores):.1f} "
                   f"ret={float(RET.mean()):.3f} v_loss={v_loss.item():.4f} ent={ent.item() if train_actor else 0:.3f} "
                   f"bc_w={bc_w:.2f} {'train' if train_actor else 'critic-warmup'} {time.time() - t0:.0f}s")
            if train_actor and (it + 1) % a.eval_every == 0:
                ev = evaluate(pool, actor)
                fit = fitness(ev)
                msg += f" | eval {ev} fitness={fit:.4f}"
                if fit > best:
                    best, bad = fit, 0
                    best_state = {k: v.clone() for k, v in actor.state_dict().items()}
                    torch.save(best_state, out / "model_ppo.pt")
                    msg += " *best*"
                else:
                    bad += 1
                    if bad >= 2 and fit < 0.9 * best:
                        actor.load_state_dict(best_state)
                        opt_a = torch.optim.Adam(actor.parameters(), lr=a.lr * 0.5)
                        a.lr *= 0.5
                        bad = 0
                        msg += " ROLLBACK, lr halved"
            print(msg, flush=True)
    torch.save(best_state, out / "model_ppo.pt")
    print(f"best fitness {best:.4f} saved to {out / 'model_ppo.pt'}")


if __name__ == "__main__":
    main()
