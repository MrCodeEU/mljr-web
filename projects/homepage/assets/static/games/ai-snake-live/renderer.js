// Canvas renderer for the shared snake game. Datastar patches the `snake`
// signal on every server tick; data-effect calls window.snakeAIRender.
(function () {
    var CELL = 20;
    function css(name, fallback) {
        var v = getComputedStyle(document.documentElement).getPropertyValue(name).trim();
        return v || fallback;
    }
    window.snakeAIRender = function (s) {
        if (!s || !s.n) return;
        var canvas = document.getElementById('snake-ai-canvas');
        if (!canvas) return;
        var size = s.n * CELL;
        if (canvas.width !== size) { canvas.width = size; canvas.height = size; }
        var ctx = canvas.getContext('2d');
        ctx.fillStyle = css('--surface', '#f2efe6');
        ctx.fillRect(0, 0, size, size);
        ctx.fillStyle = css('--accent', '#e5484d');
        ctx.fillRect(s.food[0] * CELL, s.food[1] * CELL, CELL, CELL);
        var ink = css('--ink', '#111');
        var cells = s.cells || [];
        for (var i = 0; i < cells.length; i += 2) {
            ctx.fillStyle = i === 0 ? css('--accent', '#e5484d') : ink;
            ctx.fillRect(cells[i] * CELL + 1, cells[i + 1] * CELL + 1, CELL - 2, CELL - 2);
        }
    };
})();
