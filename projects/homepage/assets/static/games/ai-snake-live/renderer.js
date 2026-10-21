// Canvas renderer for the live AI snake page. Datastar patches the `snake` signal on every server move;
// data-effect calls window.snakeLiveRender. A full state (new `full` counter, with `body`) rebuilds the
// board; every state also carries the last few moves (`m`: head cell, grew pairs, newest = `tick`), which
// are applied incrementally, so a missed message is caught up from the next one.
(function () {
    var cell = 16, n = 0, body = [], food = -1, lastFull = 0, lastTick = 0;
    var canvas = null, ctx = null, bg = '#f2efe6', accent = '#e5484d', ink = '#111';

    function css(name, fallback) {
        var v = getComputedStyle(document.documentElement).getPropertyValue(name).trim();
        return v || fallback;
    }
    function square(i, color) {
        ctx.fillStyle = color;
        ctx.fillRect((i % n) * cell + 1, Math.floor(i / n) * cell + 1, cell - 2, cell - 2);
    }
    function clear(i) {
        ctx.fillStyle = bg;
        ctx.fillRect((i % n) * cell, Math.floor(i / n) * cell, cell, cell);
    }
    function drawAll() {
        ctx.fillStyle = bg;
        ctx.fillRect(0, 0, canvas.width, canvas.height);
        for (var i = body.length - 1; i >= 1; i--) square(body[i], ink);
        if (body.length) square(body[0], accent);
        if (food >= 0) square(food, accent);
    }
    function rebuild(s) {
        n = s.n;
        cell = Math.max(4, Math.floor(800 / n));
        bg = css('--surface', bg); accent = css('--accent', accent); ink = css('--ink', ink);
        canvas.width = canvas.height = n * cell;
        body = s.body.slice();
        food = s.food;
        lastFull = s.full;
        lastTick = s.tick;
        drawAll();
    }
    function advance(s) {
        var m = s.m || [], k = m.length / 2, first = s.tick - k + 1;
        if (first > lastTick + 1) return; // gap larger than the repeated moves: wait for the next full state
        for (var i = 0; i < k; i++) {
            if (first + i <= lastTick) continue;
            body.unshift(m[2 * i]);
            if (!m[2 * i + 1]) clear(body.pop());
            if (body.length > 1) square(body[1], ink);
            square(body[0], accent);
        }
        lastTick = s.tick;
        if (s.food !== food) { food = s.food; if (food >= 0) square(food, accent); }
    }

    // ---- network picture: one forward pass, heavily abstracted (94,000 weights are not drawn) ----
    var netCanvas = null, netCtx = null, lastNet = 0;
    var LABELS = ['left', 'straight', 'right'];

    function bytes(b64) {
        var raw = atob(b64), out = new Uint8Array(raw.length);
        for (var i = 0; i < raw.length; i++) out[i] = raw.charCodeAt(i);
        return out;
    }
    function heat(x, y, size, gap, v, signed) {
        // v: 0..1 magnitude; signed: draw below-0.5 values in blue (inputs are mapped to 0..1 around 0.5)
        var c = netCtx, a;
        if (signed) {
            var d = (v - 0.5) * 2;
            a = Math.min(1, Math.abs(d));
            c.fillStyle = d >= 0 ? accent : '#78b4ff';
        } else {
            a = v;
            c.fillStyle = accent;
        }
        c.globalAlpha = 0.08 + 0.92 * a;
        c.fillRect(x, y, size - gap, size - gap);
        c.globalAlpha = 1;
    }
    function block(x, y, cols, rows, cell, vals, offset, signed) {
        for (var r = 0; r < rows; r++)
            for (var q = 0; q < cols; q++)
                heat(x + q * cell, y + r * cell, cell, 2, vals[offset + r * cols + q] / 255, signed);
    }
    function topCells(vals, k) {
        var idx = [];
        for (var i = 0; i < vals.length; i++) idx.push(i);
        idx.sort(function (a, b) { return vals[b] - vals[a]; });
        return idx.slice(0, k);
    }
    function label(text, x, y, align) {
        netCtx.fillStyle = ink; netCtx.globalAlpha = 0.65;
        netCtx.font = '11px sans-serif'; netCtx.textAlign = align || 'left';
        netCtx.fillText(text, x, y); netCtx.globalAlpha = 1;
    }
    function drawNet(net) {
        if (!netCanvas) {
            netCanvas = document.getElementById('ai-snake-live-net');
            if (!netCanvas) return;
            netCtx = netCanvas.getContext('2d');
            netCanvas.width = 1600; netCanvas.height = 520;
        }
        var c = netCtx;
        c.setTransform(2, 0, 0, 2, 0, 0);
        c.clearRect(0, 0, 800, 260);
        var inp = bytes(net.in), h1 = bytes(net.h1), h2 = bytes(net.h2);
        // inputs: 3x7 per-move terms, 10 globals, 7x7 window around the head, 5x5 body map
        label('inputs · 105', 20, 18);
        block(20, 30, 7, 3, 16, inp, 0, true);
        for (var a = 0; a < 3; a++) label(LABELS[a][0].toUpperCase(), 14, 42 + a * 16, 'right');
        block(20, 92, 10, 1, 14, inp, 21, true);
        block(20, 118, 7, 7, 14, inp, 31, true);
        block(132, 118, 5, 5, 14, inp, 80, true);
        label('per move', 20, 86); label('globals', 20, 110 + 0); label('window', 20, 230); label('body map', 132, 230);
        // hidden layers as 16x16 heatmaps
        var H = 11, y0 = 40, x1 = 285, x2 = 500;
        label('layer 1 · 256', x1, 18); label('layer 2 · 256', x2, 18);
        block(x1, y0, 16, 16, H, h1, 0, false);
        block(x2, y0, 16, 16, H, h2, 0, false);
        // funnels between the blocks, faint
        c.fillStyle = ink; c.globalAlpha = 0.06;
        [[250, 30, 125, x1, y0, 176], [x1 + 176, y0, 176, x2, y0, 176], [x2 + 176, y0, 176, 715, 90, 90]].forEach(function (f) {
            c.beginPath(); c.moveTo(f[0], f[1]); c.lineTo(f[3], f[4]); c.lineTo(f[3], f[4] + f[5]); c.lineTo(f[0], f[1] + f[2]); c.closePath(); c.fill();
        });
        c.globalAlpha = 1;
        // outline the strongest units of each hidden layer
        c.strokeStyle = ink; c.lineWidth = 1;
        [[x1, h1], [x2, h2]].forEach(function (L) {
            topCells(L[1], 6).forEach(function (i) {
                c.strokeRect(L[0] + (i % 16) * H - 0.5, y0 + Math.floor(i / 16) * H - 0.5, H - 1, H - 1);
            });
        });
        // outputs: probability bars, net's favourite and the move actually played
        label('output', 715, 18);
        for (var o = 0; o < 3; o++) {
            var by = 90 + o * 34, w = 70 * net.p[o];
            c.fillStyle = ink; c.globalAlpha = 0.12; c.fillRect(715, by, 70, 18);
            c.globalAlpha = 1; c.fillStyle = o === net.chosen ? accent : ink; c.globalAlpha = o === net.chosen ? 1 : 0.45;
            c.fillRect(715, by, w, 18); c.globalAlpha = 1;
            label(LABELS[o] + ' ' + Math.round(net.p[o] * 100) + '%', 715, by + 32, 'left');
        }
        if (net.top !== net.chosen) {
            c.fillStyle = accent; c.font = 'bold 11px sans-serif'; c.textAlign = 'left';
            c.fillText('safety check overrode ' + LABELS[net.top], 715, 214);
        }
    }

    window.snakeLiveRender = function (s) {
        if (s && s.net && s.net.t !== lastNet) { lastNet = s.net.t; drawNet(s.net); }
        if (!s || !s.full || !s.body) return;
        if (!canvas) {
            canvas = document.getElementById('ai-snake-live-canvas');
            if (!canvas) return;
            ctx = canvas.getContext('2d');
        }
        if (s.full !== lastFull) rebuild(s);
        else if (s.tick > lastTick) advance(s);
    };
})();
