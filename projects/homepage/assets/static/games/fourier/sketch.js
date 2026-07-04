var CANVAS_W = 1000;
var CANVAS_H = 760;

var DRAW_X = 20, DRAW_Y = 20, DRAW_W = 400, DRAW_H = 400;
var EPI_X = 460, EPI_Y = 20, EPI_W = 520, EPI_H = 400;
var EPI_CX = EPI_X + EPI_W / 2, EPI_CY = EPI_Y + EPI_H / 2;
var WAVE_X = 20, WAVE_Y = 460, WAVE_W = 960, WAVE_H = 260;
var WAVE_BASE_Y = WAVE_Y + WAVE_H / 2;

var RESAMPLE_N = 150;
var DESIRED_RADIUS = 150;

var mode = 'draw'; // 'draw' | 'play'
var rawPoints = [];
var sortedFourier = [];
var t = 0;
var dt;
var tracePoints = [];
var waveHistory = [];
var wavesInput;

function pointInDrawRegion(x, y) {
    return x >= DRAW_X && x <= DRAW_X + DRAW_W && y >= DRAW_Y && y <= DRAW_Y + DRAW_H;
}

function resample(points, n) {
    if (points.length < 2) return [];
    var lens = [0];
    var total = 0;
    for (var i = 1; i < points.length; i++) {
        total += dist(points[i - 1].x, points[i - 1].y, points[i].x, points[i].y);
        lens.push(total);
    }
    if (total === 0) return [];
    var out = [];
    var seg = 0;
    for (var k = 0; k < n; k++) {
        var target = (total * k) / n;
        while (seg < lens.length - 2 && lens[seg + 1] < target) seg++;
        var segLen = lens[seg + 1] - lens[seg];
        var f = segLen > 0 ? (target - lens[seg]) / segLen : 0;
        var a = points[seg], b = points[seg + 1];
        out.push({ x: lerp(a.x, b.x, f), y: lerp(a.y, b.y, f) });
    }
    return out;
}

// Discrete Fourier Transform of a complex-valued signal (re = x, im = y).
// Direct O(n^2) sum - at n=150 that's ~22k multiply-adds, done once when
// Play is pressed, not per frame, so there's no need for an FFT here.
function dft(points) {
    var n = points.length;
    var out = [];
    for (var k = 0; k < n; k++) {
        var re = 0, im = 0;
        for (var i = 0; i < n; i++) {
            var phi = (2 * Math.PI * k * i) / n;
            re += points[i].re * Math.cos(phi) + points[i].im * Math.sin(phi);
            im += points[i].im * Math.cos(phi) - points[i].re * Math.sin(phi);
        }
        re /= n;
        im /= n;
        out.push({ freq: k, amp: Math.sqrt(re * re + im * im), phase: Math.atan2(im, re) });
    }
    return out;
}

function computeFourierFromPath() {
    var resampled = resample(rawPoints, RESAMPLE_N);
    if (resampled.length === 0) return;

    var cx = 0, cy = 0;
    for (var i = 0; i < resampled.length; i++) { cx += resampled[i].x; cy += resampled[i].y; }
    cx /= resampled.length;
    cy /= resampled.length;

    var maxDist = 1;
    for (var j = 0; j < resampled.length; j++) {
        var d = dist(resampled[j].x, resampled[j].y, cx, cy);
        if (d > maxDist) maxDist = d;
    }
    var scale = DESIRED_RADIUS / maxDist;

    var complexPoints = resampled.map(function (p) {
        return { re: (p.x - cx) * scale, im: (p.y - cy) * scale };
    });

    var fourier = dft(complexPoints);
    fourier.sort(function (a, b) { return b.amp - a.amp; });
    sortedFourier = fourier;

    dt = (2 * Math.PI) / RESAMPLE_N;
    t = 0;
    tracePoints = [];
    waveHistory = [];
    mode = 'play';
}

function resetDrawing() {
    mode = 'draw';
    rawPoints = [];
    tracePoints = [];
    waveHistory = [];
    t = 0;
}

function setup() {
    var canvas = createCanvas(CANVAS_W, CANVAS_H);
    canvas.parent('fourier-game');
    wavesInput = document.getElementById('fourier-waves');

    document.getElementById('fourier-play').addEventListener('click', computeFourierFromPath);
    document.getElementById('fourier-clear').addEventListener('click', resetDrawing);
    wavesInput.addEventListener('input', function () {
        document.getElementById('fourier-waves-value').textContent = wavesInput.value;
    });
}

function mousePressed() {
    if (mode === 'draw' && pointInDrawRegion(mouseX, mouseY)) {
        rawPoints = [{ x: mouseX - DRAW_X, y: mouseY - DRAW_Y }];
    }
}

function mouseDragged() {
    if (mode === 'draw' && pointInDrawRegion(mouseX, mouseY) && rawPoints.length > 0) {
        var p = { x: mouseX - DRAW_X, y: mouseY - DRAW_Y };
        var last = rawPoints[rawPoints.length - 1];
        if (dist(p.x, p.y, last.x, last.y) > 2) rawPoints.push(p);
    }
}

function drawPanelFrame(x, y, w, h, label) {
    noFill();
    stroke(90);
    strokeWeight(1);
    rect(x, y, w, h);
    noStroke();
    fill(160);
    textSize(12);
    text(label, x + 6, y - 6);
}

function drawDrawRegion() {
    var labels = window.FOURIER_I18N || {};
    drawPanelFrame(DRAW_X, DRAW_Y, DRAW_W, DRAW_H, mode === 'draw' ? (labels.drawHint || 'Draw here, then press Play') : (labels.yourDrawing || 'Your drawing'));
    noFill();
    stroke(120, 200, 255);
    strokeWeight(2);
    beginShape();
    for (var i = 0; i < rawPoints.length; i++) vertex(DRAW_X + rawPoints[i].x, DRAW_Y + rawPoints[i].y);
    endShape();
}

// Chains rotating vectors tip-to-tail (largest amplitude first, for the
// classic nested-circles look) and returns where the final tip lands -
// that point is both what gets traced on the right and what feeds the
// wave graph below.
function drawEpicycles(count) {
    push();
    translate(EPI_CX, EPI_CY);
    var x = 0, y = 0;
    noFill();
    for (var i = 0; i < count && i < sortedFourier.length; i++) {
        var term = sortedFourier[i];
        var prevx = x, prevy = y;
        x += term.amp * Math.cos(term.freq * t + term.phase);
        y += term.amp * Math.sin(term.freq * t + term.phase);
        stroke(90);
        strokeWeight(1);
        if (term.amp > 0.5) ellipse(prevx, prevy, term.amp * 2, term.amp * 2);
        stroke(200);
        line(prevx, prevy, x, y);
    }
    pop();
    return { x: x, y: y };
}

function draw() {
    background(24);
    var labels = window.FOURIER_I18N || {};
    drawDrawRegion();
    drawPanelFrame(EPI_X, EPI_Y, EPI_W, EPI_H, labels.reconstruction || 'Fourier reconstruction');
    drawPanelFrame(WAVE_X, WAVE_Y, WAVE_W, WAVE_H, labels.waveLabel || 'Layered waves (y over time)');

    if (mode !== 'play' || sortedFourier.length === 0) return;

    var count = Number(wavesInput.value);
    var tip = drawEpicycles(count);

    tracePoints.push({ x: EPI_CX + tip.x, y: EPI_CY + tip.y });
    if (tracePoints.length > RESAMPLE_N) tracePoints.shift();

    noFill();
    stroke(255, 120, 150);
    strokeWeight(2);
    beginShape();
    for (var i = 0; i < tracePoints.length; i++) vertex(tracePoints[i].x, tracePoints[i].y);
    endShape();

    fill(255, 120, 150);
    noStroke();
    circle(EPI_CX + tip.x, EPI_CY + tip.y, 6);

    waveHistory.unshift(tip.y);
    if (waveHistory.length > WAVE_W) waveHistory.pop();

    stroke(90, 160, 255);
    strokeWeight(1);
    line(EPI_CX + tip.x, EPI_CY + tip.y, WAVE_X + WAVE_W, WAVE_BASE_Y + tip.y);

    noFill();
    stroke(90, 160, 255);
    strokeWeight(2);
    beginShape();
    for (var w = 0; w < waveHistory.length; w++) {
        vertex(WAVE_X + WAVE_W - w, WAVE_BASE_Y + waveHistory[w]);
    }
    endShape();

    t += dt;
    if (t > 2 * Math.PI) {
        t -= 2 * Math.PI;
        tracePoints = [];
    }
}
