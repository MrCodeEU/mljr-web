# snakeai-train

Offline training pipeline for the live AI snake (`/games/ai-snake-live`). It produces the
weights that `internal/snakeai` runs in plain Go (no ML runtime on the server).

## What it does

1. **Teacher** (`snakenet/teacher.py`): a Hamiltonian cycle plus safe shortcuts. Cells strictly between
   the head and the tail along the cycle are always free, so jumping ahead toward the apple is safe as
   long as it keeps a margin to the tail. Shortcuts are only used while the snake covers less than half
   the board (the cutoff that fills a 20x20 board fastest; 0.5 beat 0.15-0.85 and "always").
2. **Imitation + DAgger** (`snakenet/train.py`): an MLP (133 features -> 128 -> 128 -> 3 relative
   actions: left, straight, right) copies the teacher. After round 0 the network drives and the teacher
   labels every state it visits, with warm-started late-game states and boards up to 50x50.
3. **Export** (`snakenet/export.py`): writes `model.json` and `golden.json` (plus an ONNX file).

`snakenet/ppo.py` is an experimental PPO fine-tune. It did not beat the teacher (see below).

## Commands

```bash
cd tools/snakeai-train
uv sync                                                      # CPU torch
uv run python -m snakenet.train --teacher cycle --rounds 6 --games 200
uv run python -m snakenet.evaluate --size 20 --games 16 --max-steps 25000
uv run python -m snakenet.export
cp out/model.json  ../../internal/snakeai/model.json
cp out/golden.json ../../internal/snakeai/testdata/golden.json
go test ./internal/snakeai/...                               # Python/Go parity check
```

`checkpoints/model_v3.pt` is the network behind the committed `model.json`. Continue from it with
`train --init checkpoints/model_v3.pt`.

## The Python/Go contract

`snakenet/env.py` (`features`) and `internal/snakeai/features.go` must produce identical feature vectors.
`FEATURE_VERSION` is checked when the Go side loads the model. If you change a feature, bump the version in
both places, retrain, re-export, and let `golden.json` (30 real states on 10x10 to 50x50 boards) prove the
Go port still matches. Board sizes must be even (the Hamiltonian cycle needs it).

## Results (cycle teacher, imitation net)

| Board | Budget | Net | Teacher |
|---|---|---|---|
| 10x10 | 3000 steps | 99/99, 0 deaths | 99/99, 0 deaths |
| 20x20 | 25k steps | 399/399, 0 deaths (16/16 filled) | 399/399, 0 deaths |
| 50x50 | 20k steps | 240.6, 0 deaths | 240.5, 0 deaths |

The net matches its teacher, including the teacher's habits: it follows the cycle direction even when a
shortcut or a different route would be as good. Measure speed as score at a fixed step budget, never just
"final score": early experiments looked like the teacher was weak only because of a step cap.

## Known open problem: faster than the teacher

A PPO run with a speed objective (eat reward, step cost, death penalty, decaying teacher anchor) stayed at
the teacher's level. The policy's entropy collapsed to about 0.002, so it never explored. Next attempt:
temperature on the logits plus a larger entropy bonus, start from an imitation of the faster BFS teacher
(`--teacher bfs`) with the cycle features as inputs only (never as labels), and judge by deaths = 0 and
steps-to-fill below the cycle teacher.
