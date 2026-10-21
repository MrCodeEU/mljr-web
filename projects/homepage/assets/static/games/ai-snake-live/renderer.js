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

    window.snakeLiveRender = function (s) {
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
