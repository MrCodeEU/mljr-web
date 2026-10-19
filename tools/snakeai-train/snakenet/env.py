"""Snake environment + feature extraction.

The feature layout here is the contract with the Go runtime: any change
must bump FEATURE_VERSION and be mirrored in Go (golden.json guards parity).
"""
from collections import deque

import numpy as np

FEATURE_VERSION = 3
# clockwise order: heading h turns right to (h+1)%4, left to (h+3)%4
DIRS = [(0, -1), (1, 0), (0, 1), (-1, 0)]  # up, right, down, left
# relative actions
LEFT, STRAIGHT, RIGHT = 0, 1, 2
TURN = {LEFT: 3, STRAIGHT: 0, RIGHT: 1}
WIN = 3  # window radius -> (2*WIN+1)^2 cells
GRID = 5  # coarse ego density map is GRID x GRID bins
N_FEATURES = 3 * 6 + 3 + 4 + 2 * GRID * GRID + (2 * WIN + 1) ** 2 + 3 * 2 + 3


class Game:
    def __init__(self, n, rng):
        self.n = n
        self.rng = rng
        self.snake = deque([(n // 2, n // 2)])  # head first
        self.heading = int(rng.integers(4))
        self.score = 0
        self.since_food = 0
        self.dead = False
        self.food = None
        self.place_food()

    def place_food(self):
        occ = set(self.snake)
        free = [(x, y) for x in range(self.n) for y in range(self.n) if (x, y) not in occ]
        if not free:
            self.food = None
            return
        self.food = free[int(self.rng.integers(len(free)))]

    def blocked(self, c, tail_free=True):
        x, y = c
        if not (0 <= x < self.n and 0 <= y < self.n):
            return True
        if tail_free and len(self.snake) > 1:
            return c in list(self.snake)[:-1]
        return c in self.snake

    def cell(self, action):
        d = DIRS[(self.heading + TURN[action]) % 4]
        h = self.snake[0]
        return (h[0] + d[0], h[1] + d[1])

    def step(self, action):
        nxt = self.cell(action)
        eating = nxt == self.food
        if self.blocked(nxt, tail_free=not eating):
            self.dead = True
            return
        self.heading = (self.heading + TURN[action]) % 4
        self.snake.appendleft(nxt)
        self.since_food += 1
        if eating:
            self.score += 1
            self.since_food = 0
            self.place_food()
        else:
            self.snake.pop()
        if self.food is None or self.since_food > 2 * self.n * self.n:
            self.dead = True


def flood(game, snake):
    """Free cells reachable from snake[0] (the head); the tail counts as free.
    Returns (count, tail_reached)."""
    start, tail = snake[0], snake[-1]
    blocked = set(snake[1:-1])
    n = game.n
    seen = {start}
    q = deque([start])
    tail_reached = start == tail
    while q:
        x, y = q.popleft()
        for dx, dy in DIRS:
            c = (x + dx, y + dy)
            if c in seen or not (0 <= c[0] < n and 0 <= c[1] < n) or c in blocked:
                continue
            seen.add(c)
            if c == tail:
                tail_reached = True
            q.append(c)
    return len(seen), tail_reached


_CYCLES = {}


def cycle_for(n):
    """Hamiltonian cycle (even n): same comb layout as ai-snake/sketch.js buildCycle.
    Returns (cells, index) where index[(x, y)] is the position along the cycle."""
    if n in _CYCLES:
        return _CYCLES[n]
    assert n % 2 == 0, "Hamiltonian cycle needs an even board size"
    cells = [(x, 0) for x in range(n)] + [(n - 1, y) for y in range(1, n)]
    for col in range(n - 2, 0, -1):
        if (n - 1 - col) % 2 == 1:
            cells += [(col, y) for y in range(n - 1, 0, -1)]
        else:
            cells += [(col, y) for y in range(1, n)]
    cells += [(0, y) for y in range(n - 1, 0, -1)]
    index = {c: i for i, c in enumerate(cells)}
    _CYCLES[n] = (cells, index)
    return _CYCLES[n]


def cyc_dist(a, b, length):
    """Steps forward along the cycle from position a to b (length if a == b)."""
    d = (b - a) % length
    return d if d else length


def bfs_dist(n, start, blocked):
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


def free_regions(n, blocked):
    """Connected components of free cells -> (count, largest, total_free)."""
    seen = set(blocked)
    sizes = []
    for sx in range(n):
        for sy in range(n):
            if (sx, sy) in seen:
                continue
            seen.add((sx, sy))
            q = deque([(sx, sy)])
            size = 0
            while q:
                x, y = q.popleft()
                size += 1
                for dx, dy in DIRS:
                    c = (x + dx, y + dy)
                    if c in seen or not (0 <= c[0] < n and 0 <= c[1] < n):
                        continue
                    seen.add(c)
                    q.append(c)
            sizes.append(size)
    return len(sizes), (max(sizes) if sizes else 0), sum(sizes)


def features(game):
    n = game.n
    head = game.snake[0]
    body_list = list(game.snake)
    length = len(body_list)
    f = DIRS[game.heading]
    r = DIRS[(game.heading + 1) % 4]
    fx, fy = game.food
    dist0 = abs(fx - head[0]) + abs(fy - head[1])
    cap = 2.0 * n
    out = np.zeros(N_FEATURES, dtype=np.float32)
    for i, a in enumerate((LEFT, STRAIGHT, RIGHT)):
        c = game.cell(a)
        eating = c == game.food
        blk = game.blocked(c, tail_free=not eating)
        o = i * 6
        out[o] = 1.0 if blk else 0.0
        out[o + 4] = 1.0  # unreachable defaults
        out[o + 5] = 1.0
        if not blk:
            virt = [c] + body_list[: length - (0 if eating else 1)]
            count, tail_ok = flood(game, virt)
            out[o + 1] = min(count / max(length, 1), 2.0) / 2.0
            out[o + 2] = 1.0 if tail_ok else 0.0
            out[o + 3] = float(dist0 - (abs(fx - c[0]) + abs(fy - c[1])))
            dist = bfs_dist(n, c, set(virt[1:-1]))
            if virt[-1] in dist:
                out[o + 4] = min(dist[virt[-1]] / cap, 1.0)
            if game.food in dist:
                out[o + 5] = min(dist[game.food] / cap, 1.0)
    dx, dy = fx - head[0], fy - head[1]
    out[18] = float(np.clip((dx * r[0] + dy * r[1]) / (n / 2), -1, 1))  # food right(+)/left(-)
    out[19] = float(np.clip((dx * f[0] + dy * f[1]) / (n / 2), -1, 1))  # food ahead(+)/behind(-)
    out[20] = length / (n * n)
    # global: free-space topology (tail counts as free, head as blocked)
    blocked_now = set(body_list[:-1])
    n_reg, largest, total = free_regions(n, blocked_now)
    out[21] = min(n_reg, 4) / 4.0
    out[22] = largest / total if total else 0.0
    out[23] = total / (n * n)
    # food reachable from head at all (path length, 1 = unreachable)
    dist = bfs_dist(n, head, set(body_list[1:-1]))
    out[24] = min(dist[game.food] / cap, 1.0) if game.food in dist else 1.0
    # coarse ego density maps: body fraction + in-board fraction per bin
    body = set(body_list[:-1])
    body_cnt = np.zeros((GRID, GRID), np.float32)
    area_cnt = np.zeros((GRID, GRID), np.float32)
    def bins(xs, ys):
        ox, oy = xs - head[0], ys - head[1]
        u = ox * r[0] + oy * r[1]
        v = ox * f[0] + oy * f[1]
        bu = np.minimum(((u + n) * GRID) // (2 * n), GRID - 1)
        bv = np.minimum(((n - v) * GRID) // (2 * n), GRID - 1)  # forward = top row
        return bv.astype(int), bu.astype(int)

    gx, gy = np.indices((n, n))
    bv, bu = bins(gx.ravel(), gy.ravel())
    np.add.at(area_cnt, (bv, bu), 1)
    if body:
        bx = np.array([c[0] for c in body])
        by = np.array([c[1] for c in body])
        bv, bu = bins(bx, by)
        np.add.at(body_cnt, (bv, bu), 1)
    cell_area = (2.0 * n / GRID) ** 2
    out[25:25 + GRID * GRID] = (body_cnt / cell_area).ravel()
    out[25 + GRID * GRID:25 + 2 * GRID * GRID] = (area_cnt / cell_area).ravel()
    k = 25 + 2 * GRID * GRID
    for v in range(WIN, -WIN - 1, -1):  # forward rows first
        for u in range(-WIN, WIN + 1):
            x = head[0] + u * r[0] + v * f[0]
            y = head[1] + u * r[1] + v * f[1]
            wall = not (0 <= x < n and 0 <= y < n)
            out[k] = 1.0 if (wall or (x, y) in body) else 0.0
            k += 1
    # cycle features: forward cycle distance of each move cell, then global tail/food distances
    cells, index = cycle_for(n)
    total = len(cells)
    hi = index[head]
    base = 3 * 6 + 3 + 4 + 2 * GRID * GRID + (2 * WIN + 1) ** 2
    for i, a in enumerate((LEFT, STRAIGHT, RIGHT)):
        c = game.cell(a)
        out[base + i * 2] = 1.0 if c not in index else 0.0  # off-board marker
        out[base + i * 2 + 1] = cyc_dist(hi, index[c], total) / total if c in index else 1.0
    out[base + 6] = cyc_dist(hi, index[body_list[-1]], total) / total
    out[base + 7] = cyc_dist(hi, index[game.food], total) / total
    out[base + 8] = length / total
    return out
