"""Alternative Hamiltonian cycles: boundary walk around a random spanning tree of the (n/2 x n/2) block grid.
Teacher spec for the bench: snakenet.cycles:tree<seed>  (e.g. tree3), snakenet.cycles:comb"""
import random

from . import env
from .teacher import cycle_teacher_action

_ORIG = {}


def tree_cycle(n, seed):
    m = n // 2
    rng = random.Random(seed)
    parent = list(range(m * m))

    def find(a):
        while parent[a] != a:
            parent[a] = parent[parent[a]]
            a = parent[a]
        return a

    edges = [((x, y), (x + 1, y)) for x in range(m - 1) for y in range(m)] + \
            [((x, y), (x, y + 1)) for x in range(m) for y in range(m - 1)]
    rng.shuffle(edges)
    adj = {}

    def link(a, b):
        adj.setdefault(a, set()).add(b)
        adj.setdefault(b, set()).add(a)

    # all 2x2 blocks start as 4-cycles
    for bx in range(m):
        for by in range(m):
            c = [(2 * bx, 2 * by), (2 * bx + 1, 2 * by), (2 * bx + 1, 2 * by + 1), (2 * bx, 2 * by + 1)]
            for i in range(4):
                link(c[i], c[(i + 1) % 4])
    for (ax, ay), (bx, by) in edges:
        ra, rb = find(ay * m + ax), find(by * m + bx)
        if ra == rb:
            continue
        parent[ra] = rb
        if bx == ax + 1:  # horizontal tree edge: A right column <-> B left column
            a1, a2 = (2 * ax + 1, 2 * ay), (2 * ax + 1, 2 * ay + 1)
            b1, b2 = (2 * bx, 2 * by), (2 * bx, 2 * by + 1)
        else:  # vertical: A bottom row <-> B top row
            a1, a2 = (2 * ax, 2 * ay + 1), (2 * ax + 1, 2 * ay + 1)
            b1, b2 = (2 * bx, 2 * by), (2 * bx + 1, 2 * by)
        adj[a1].discard(a2); adj[a2].discard(a1)
        adj[b1].discard(b2); adj[b2].discard(b1)
        link(a1, b1); link(a2, b2)
    cells, prev, cur = [(0, 0)], None, (0, 0)
    while True:
        nxt = next(c for c in adj[cur] if c != prev) if prev is not None else min(adj[cur])
        if nxt == (0, 0):
            break
        cells.append(nxt)
        prev, cur = cur, nxt
    assert len(cells) == n * n, (len(cells), n * n)
    return cells


def install(n, cells):
    env._CYCLES[n] = (cells, {c: i for i, c in enumerate(cells)})


def _teacher(tag, builder):
    cur = {}

    def act(game):
        n = game.n
        if cur.get("n") != (n, tag):
            if n not in _ORIG:
                _ORIG[n] = env._CYCLES.get(n) or env.cycle_for(n)
            install(n, builder(n)) if builder else env._CYCLES.__setitem__(n, _ORIG[n])
            cur["n"] = (n, tag)
        return cycle_teacher_action(game)
    return act


def __getattr__(name):
    if name == "comb":
        return _teacher("comb", None)
    if name.startswith("tree"):
        seed = int(name[4:])
        return _teacher(name, lambda n: tree_cycle(n, seed))
    raise AttributeError(name)
