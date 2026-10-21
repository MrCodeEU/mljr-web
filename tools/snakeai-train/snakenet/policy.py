"""NumPy MLP forward pass (same math the Go runtime runs)."""
import numpy as np


class NumpyMLP:
    def __init__(self, layers):
        self.layers = [(np.asarray(w, np.float32), np.asarray(b, np.float32)) for w, b in layers]

    def logits(self, x):
        for i, (w, b) in enumerate(self.layers):
            x = x @ w.T + b
            if i < len(self.layers) - 1:
                x = np.maximum(x, 0)
        return x

    def act(self, x):
        return int(np.argmax(self.logits(x)))
