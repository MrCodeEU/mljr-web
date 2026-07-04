var W = 900, H = 500;
var values = [];
var highlight = [];
var gen = null;
var comparisons = 0, swaps = 0;

var algoInput, sizeInput, speedInput, statusEl;

function shuffleBars() {
    var n = Number(sizeInput.value);
    values = [];
    for (var i = 0; i < n; i++) values.push(Math.random() * 0.9 + 0.05);
    highlight = [];
    comparisons = 0;
    swaps = 0;
    gen = null;
    updateStatus();
}

function swap(arr, i, j) {
    var tmp = arr[i];
    arr[i] = arr[j];
    arr[j] = tmp;
    swaps++;
}

// Every algorithm below is written as a generator that yields after each
// comparison/swap - that step is what gets rendered as one animation frame,
// so the visualizer can drive any of them through the exact same loop
// without each algorithm needing its own custom stepping logic.
function* bubbleSort(arr) {
    for (var i = 0; i < arr.length; i++) {
        var swappedAny = false;
        for (var j = 0; j < arr.length - i - 1; j++) {
            highlight = [j, j + 1];
            comparisons++;
            yield;
            if (arr[j] > arr[j + 1]) {
                swap(arr, j, j + 1);
                swappedAny = true;
                yield;
            }
        }
        if (!swappedAny) break;
    }
}

function* selectionSort(arr) {
    for (var i = 0; i < arr.length; i++) {
        var min = i;
        for (var j = i + 1; j < arr.length; j++) {
            highlight = [min, j];
            comparisons++;
            yield;
            if (arr[j] < arr[min]) min = j;
        }
        if (min !== i) { swap(arr, i, min); yield; }
    }
}

function* insertionSort(arr) {
    for (var i = 1; i < arr.length; i++) {
        var j = i;
        while (j > 0) {
            highlight = [j - 1, j];
            comparisons++;
            yield;
            if (arr[j - 1] > arr[j]) {
                swap(arr, j - 1, j);
                j--;
                yield;
            } else {
                break;
            }
        }
    }
}

function* merge(arr, lo, mid, hi) {
    var left = arr.slice(lo, mid + 1);
    var right = arr.slice(mid + 1, hi + 1);
    var i = 0, j = 0, k = lo;
    while (i < left.length && j < right.length) {
        highlight = [k, k];
        comparisons++;
        yield;
        if (left[i] <= right[j]) { arr[k] = left[i]; i++; }
        else { arr[k] = right[j]; j++; }
        swaps++;
        k++;
        yield;
    }
    while (i < left.length) { arr[k] = left[i]; i++; k++; swaps++; yield; }
    while (j < right.length) { arr[k] = right[j]; j++; k++; swaps++; yield; }
}

function* mergeSort(arr, lo, hi) {
    if (lo === undefined) { lo = 0; hi = arr.length - 1; }
    if (lo >= hi) return;
    var mid = Math.floor((lo + hi) / 2);
    yield* mergeSort(arr, lo, mid);
    yield* mergeSort(arr, mid + 1, hi);
    yield* merge(arr, lo, mid, hi);
}

function* partition(arr, lo, hi) {
    var pivot = arr[hi];
    var i = lo - 1;
    for (var j = lo; j < hi; j++) {
        highlight = [j, hi];
        comparisons++;
        yield;
        if (arr[j] < pivot) {
            i++;
            swap(arr, i, j);
            yield;
        }
    }
    swap(arr, i + 1, hi);
    yield;
    return i + 1;
}

function* quickSort(arr, lo, hi) {
    if (lo === undefined) { lo = 0; hi = arr.length - 1; }
    if (lo >= hi) return;
    var p = yield* partition(arr, lo, hi);
    yield* quickSort(arr, lo, p - 1);
    yield* quickSort(arr, p + 1, hi);
}

function makeGenerator() {
    var algo = algoInput.value;
    if (algo === 'bubble') return bubbleSort(values);
    if (algo === 'selection') return selectionSort(values);
    if (algo === 'insertion') return insertionSort(values);
    if (algo === 'merge') return mergeSort(values);
    return quickSort(values);
}

function updateStatus() {
    statusEl.textContent = 'Comparisons: ' + comparisons + '  Swaps: ' + swaps;
}

function setup() {
    var canvas = createCanvas(W, H);
    canvas.parent('sorting-game');
    colorMode(HSB, 360, 100, 100);
    algoInput = document.getElementById('sorting-algo');
    sizeInput = document.getElementById('sorting-size');
    speedInput = document.getElementById('sorting-speed');
    statusEl = document.getElementById('sorting-status');

    document.getElementById('sorting-shuffle').addEventListener('click', shuffleBars);
    document.getElementById('sorting-start').addEventListener('click', function () {
        if (!gen) gen = makeGenerator();
    });

    shuffleBars();
}

function draw() {
    background(20);

    if (gen) {
        var stepsPerFrame = Number(speedInput.value);
        for (var s = 0; s < stepsPerFrame; s++) {
            var res = gen.next();
            if (res.done) { gen = null; break; }
        }
        updateStatus();
    }

    var n = values.length;
    var barW = W / n;
    for (var i = 0; i < n; i++) {
        var h = values[i] * (H - 20);
        var isHighlighted = highlight.indexOf(i) !== -1;
        // Hue tracks the bar's own value, not its position - so a sorted
        // array reads as a clean rainbow gradient, and a shuffled one shows
        // the same colors scattered out of order.
        var hue = values[i] * 300;
        if (isHighlighted) fill(hue, 60, 100);
        else fill(hue, 80, 80);
        noStroke();
        rect(i * barW, H - h, Math.max(1, barW - 1), h);
    }
}
