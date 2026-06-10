import "./altcha.js";

const shaWorkerURL = new URL("./altcha-workers/sha.js", import.meta.url);

for (const algorithm of ["SHA-1", "SHA-256", "SHA-512"]) {
  if (!globalThis.$altcha.algorithms.has(algorithm)) {
    globalThis.$altcha.algorithms.set(
      algorithm,
      () => new Worker(shaWorkerURL),
    );
  }
}
