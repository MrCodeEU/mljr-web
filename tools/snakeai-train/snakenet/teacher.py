"""Teacher policy: shortest path to food, accepted only if a virtual snake that
follows it can still reach its own tail. Otherwise stall by heading for the
farthest-from-tail safe cell. Returns a relative action."""
from collections import deque

from .env import DIRS, LEFT, RIGHT, STRAIGHT, cyc_dist, cycle_for, flood

ACTIONS = (LEFT, STRAIGHT, RIGHT)


def _bfs(n, start, goal, blocked):
    prev = {start: None}
    q = deque([start])
    while q:
        cur = q.popleft()
        if cur == goal:
            path = []
            while cur is not None:
                path.append(cur)
                cur = prev[cur]
            return path[::-1]
        for dx, dy in DIRS:
            c = (cur[0] + dx, cur[1] + dy)
            if c in prev or not (0 <= c[0] < n and 0 <= c[1] < n) or c in blocked:
                continue
            prev[c] = cur
            q.append(c)
    return None


def _bfs_dist(n, start, blocked):
    dist = {start: 0}
    q = deque([start])
    while q:
        cur = q.popleft()
        for dx, dy in DIRS:
            c = (cur[0] + dx, cur[1] + dy)
            if c in dist or not (0 <= c[0] < n and 0 <= c[1] < n) or c in blocked:
                continue
            dist[c] = dist[cur] + 1
            q.append(c)
    return dist


def _action_for(game, cell):
    for a in ACTIONS:
        if game.cell(a) == cell:
            return a
    return None


def teacher_action(game):
    n = game.n
    body = list(game.snake)
    head, tail = body[0], body[-1]

    # 1) shortest path to food, virtual-snake safety check
    d = DIRS[(game.heading + 2) % 4]  # cell behind the head: no reverse action exists
    behind = (head[0] + d[0], head[1] + d[1])
    path = _bfs(n, head, game.food, set(body[:-1]) | {behind})
    if path and len(path) > 1:
        virt = list(reversed(path)) + body[1:]  # head-first after walking the path
        virt = virt[: len(body) + 1]  # ate the food -> grew by 1
        count, tail_ok = flood(game, virt)
        if tail_ok or len(virt) <= 2:
            return _action_for(game, path[1])

    # 2) stall: safe move that keeps the tail reachable, farthest from the tail
    best, best_key = None, None
    for a in ACTIONS:
        c = game.cell(a)
        if game.blocked(c):
            continue
        virt = [c] + body[:-1]
        count, tail_ok = flood(game, virt)
        dist = _bfs_dist(n, c, set(virt[:-1])).get(virt[-1], -1)
        key = (1 if tail_ok else 0, dist if tail_ok else count, count)
        if best_key is None or key > best_key:
            best, best_key = a, key
    return best if best is not None else STRAIGHT


def cycle_teacher_action(game):
    """Hamiltonian cycle + safe shortcuts (cycle order of the body is invariant, so any
    cell strictly between head and tail along the cycle is free). Jumps ahead toward the
    food as far as the free stretch allows; falls back to pure cycle when the snake is long."""
    n = game.n
    cells, index = cycle_for(n)
    total = len(cells)
    body = list(game.snake)
    length = len(body)
    hi = index[body[0]]
    to_tail = cyc_dist(hi, index[body[-1]], total)
    to_food = cyc_dist(hi, index[game.food], total)
    limit = min(to_food, to_tail - 2) if length < total * 0.5 else 1
    best, best_d = None, 0
    for a in ACTIONS:
        c = game.cell(a)
        if c not in index or game.blocked(c):
            continue
        d = cyc_dist(hi, index[c], total)
        if d == 1 and best is None:
            best, best_d = a, 1  # plain cycle step is always the safe default
        if 1 < d <= limit and d > best_d:
            best, best_d = a, d
    if best is None:  # cycle step blocked (cannot happen while invariant holds)
        return teacher_action(game)
    return best
