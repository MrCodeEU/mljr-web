"""Experimental teachers: BFS shortest path with virtual-snake safety, plus ways out of the late-game stall loop.

  v1  noise     : while stalled (no food for > n steps) pick randomly among tail-safe moves half the time
  v2  relax     : while stalled, accept the food path if enough free space remains after eating, even if the tail is not reachable
  v3  both
"""
import random

from .env import DIRS, flood
from .teacher import ACTIONS, STRAIGHT, _action_for, _bfs, _bfs_dist, teacher_action

_rng = random.Random(1)


def _path(game):
    n, body = game.n, list(game.snake)
    head = body[0]
    d = DIRS[(game.heading + 2) % 4]
    behind = (head[0] + d[0], head[1] + d[1])
    return _bfs(n, head, game.food, set(body[:-1]) | {behind}), body


def _make(noise, relax):
    def act(game):
        n = game.n
        stalled = game.since_food > n
        path, body = _path(game)
        if path and len(path) > 1:
            virt = (list(reversed(path)) + body[1:])[: len(body) + 1]
            count, tail_ok = flood(game, virt)
            if tail_ok or len(virt) <= 2 or (relax and stalled and count >= len(virt)):
                return _action_for(game, path[1])
        cands = []
        for a in ACTIONS:
            c = game.cell(a)
            if game.blocked(c):
                continue
            virt = [c] + body[:-1]
            count, tail_ok = flood(game, virt)
            dist = _bfs_dist(n, c, set(virt[:-1])).get(virt[-1], -1)
            cands.append(((1 if tail_ok else 0, dist if tail_ok else count, count), a, tail_ok))
        if not cands:
            return STRAIGHT
        if noise and stalled:
            safe = [c for c in cands if c[2]]
            if safe and _rng.random() < 0.5:
                return _rng.choice(safe)[1]
        return max(cands)[1]
    return act


v1 = _make(True, False)
v2 = _make(False, True)
v3 = _make(True, True)
