"""Play games, record (features, teacher action) pairs. Used by DAgger + eval."""
import numpy as np

from .env import Game, features
from .teacher import cycle_teacher_action, teacher_action

TEACHERS = {"bfs": teacher_action, "cycle": cycle_teacher_action}


def play(n, seed, net=None, beta=0.0, eps=0.0, max_steps=3000, record=True, teacher="cycle", warm=False):
    """beta = prob of executing the teacher action, else the net's (or teacher if no net).
    eps = prob of a random action (exploration noise). Labels are always the teacher's."""
    rng = np.random.default_rng(seed)
    game = Game(n, rng)
    label = TEACHERS[teacher]
    if warm:  # fast-forward with the teacher (no features) to reach late-game states
        for _ in range(int(rng.integers(0, 50 * n * n))):
            if game.dead:
                break
            game.step(label(game))
        if game.dead or game.food is None:
            game = Game(n, rng)
    xs, ys = [], []
    steps = 0
    while not game.dead and steps < max_steps:
        x = features(game) if (record or net is not None) else None
        t = label(game)
        if record:
            xs.append(x)
            ys.append(t)
        if rng.random() < eps:
            a = int(rng.integers(3))
        elif net is None or rng.random() < beta:
            a = t
        else:
            a = net.act(x)
        game.step(a)
        steps += 1
    return game.score, steps, xs, ys


def _job(args):
    from .policy import NumpyMLP
    n, seed, layers, beta, eps, max_steps, record = args[:7]
    teacher = args[7] if len(args) > 7 else "cycle"
    warm = args[8] if len(args) > 8 else False
    net = NumpyMLP(layers) if layers is not None else None
    score, steps, xs, ys = play(n, seed, net, beta, eps, max_steps, record, teacher, warm)
    return score, steps, np.array(xs, np.float32), np.array(ys, np.int64)
