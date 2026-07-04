var SIZE = 50;
var BLOCK = 14;
var STUCK_LIMIT = 5000;
var gameOver = false;
var score = 0;
var deaths = 0;
var speed;
var strategySelect;
var strategy = 'cycle_pure';

var occupied;
var snake;
var item;
var ticksSinceProgress;
var cycle;
var cycleIndex;

function neighbors4(x, y) {
    var out = [];
    if (x - 1 >= 0) out.push({ x: x - 1, y: y });
    if (x + 1 < SIZE) out.push({ x: x + 1, y: y });
    if (y - 1 >= 0) out.push({ x: x, y: y - 1 });
    if (y + 1 < SIZE) out.push({ x: x, y: y + 1 });
    return out;
}

function cloneOccupied(o) {
    var out = new Array(SIZE);
    for (var i = 0; i < SIZE; i++) out[i] = o[i].slice();
    return out;
}

// Breadth-first search: every move costs the same, so BFS finds the shortest
// path just as well as A* while being far cheaper (no sorting, O(cells)).
// goalAlwaysFree lets the search enter a cell the snake currently occupies -
// used for "path to tail", since the tail will have moved away by the time
// the snake gets there.
function shortestPath(occ, start, goal, goalAlwaysFree) {
    var visited = new Array(SIZE);
    var prev = new Array(SIZE);
    for (var i = 0; i < SIZE; i++) {
        visited[i] = new Array(SIZE).fill(false);
        prev[i] = new Array(SIZE).fill(null);
    }

    var queue = [start];
    var qi = 0;
    visited[start.y][start.x] = true;

    while (qi < queue.length) {
        var cur = queue[qi++];
        if (cur.x === goal.x && cur.y === goal.y) {
            var path = [cur];
            var p = prev[cur.y][cur.x];
            while (p) {
                path.unshift(p);
                p = prev[p.y][p.x];
            }
            return path;
        }
        var neigh = neighbors4(cur.x, cur.y);
        for (var n = 0; n < neigh.length; n++) {
            var next = neigh[n];
            if (visited[next.y][next.x]) continue;
            var isGoal = next.x === goal.x && next.y === goal.y;
            var free = !occ[next.y][next.x] || (goalAlwaysFree && isGoal);
            if (!free) continue;
            visited[next.y][next.x] = true;
            prev[next.y][next.x] = cur;
            queue.push(next);
        }
    }
    return null;
}

// Counts reachable free cells from start - used as a fallback "how much
// room does this direction leave me" heuristic.
function floodFillCount(occ, start) {
    var visited = new Array(SIZE);
    for (var i = 0; i < SIZE; i++) visited[i] = new Array(SIZE).fill(false);
    var stack = [start];
    visited[start.y][start.x] = true;
    var count = 0;
    while (stack.length > 0) {
        var cur = stack.pop();
        count++;
        var neigh = neighbors4(cur.x, cur.y);
        for (var n = 0; n < neigh.length; n++) {
            var next = neigh[n];
            if (visited[next.y][next.x] || occ[next.y][next.x]) continue;
            visited[next.y][next.x] = true;
            stack.push(next);
        }
    }
    return count;
}

// Distance from a cell to the nearest border - used by the heuristic tiers
// to avoid coiling flush against a wall, which removes rerouting options.
function wallDistance(cell) {
    return Math.min(cell.x, SIZE - 1 - cell.x, cell.y, SIZE - 1 - cell.y);
}

// Simulates taking a single step onto `candidate` and reports whether the
// snake could still reach its own tail afterwards. This is a one-step (or
// effectively few-step, since it's re-evaluated fresh every tick) lookahead
// - which is exactly its limitation: it can't see far enough ahead to catch
// every delayed dead end. That limitation is the whole reason the "cycle
// (pure)" strategy exists below.
function isMoveSafe(candidate) {
    var isEating = candidate.x === item.x && candidate.y === item.y;
    var occCopy = cloneOccupied(occupied);
    var bodyCopy = snake.slice();
    bodyCopy.unshift(candidate);
    occCopy[candidate.y][candidate.x] = true;
    if (!isEating) {
        var tail = bodyCopy.pop();
        occCopy[tail.y][tail.x] = false;
    }
    if (bodyCopy.length <= 1) return true;
    var newHead = bodyCopy[0];
    var newTail = bodyCopy[bodyCopy.length - 1];
    return shortestPath(occCopy, newHead, newTail, true) !== null;
}

// A fixed Hamiltonian cycle over the whole grid: a single loop visiting
// every cell exactly once and returning to the start. Following it is a
// structural, mathematical guarantee against ever trapping itself - as
// long as the snake only ever moves to the next cell in this fixed order,
// its body is always one contiguous, non-crossing arc of the same loop.
// Built once and looked up by array index, so following it costs O(1) per
// tick.
//
// Layout: row 0 straight across, then column SIZE-1 down, then the
// remaining columns SIZE-2..1 boustrophedon over rows 1..SIZE-1, then
// column 0 back up to row 1 - closing the loop next to the start at (0,0).
function buildCycle() {
    var out = [];
    for (var x = 0; x < SIZE; x++) out.push({ x: x, y: 0 });
    for (var y = 1; y < SIZE; y++) out.push({ x: SIZE - 1, y: y });

    for (var col = SIZE - 2; col >= 1; col--) {
        var goingUp = (SIZE - 1 - col) % 2 === 1;
        if (goingUp) {
            for (var y2 = SIZE - 1; y2 >= 1; y2--) out.push({ x: col, y: y2 });
        } else {
            for (var y3 = 1; y3 < SIZE; y3++) out.push({ x: col, y: y3 });
        }
    }
    for (var y4 = SIZE - 1; y4 >= 1; y4--) out.push({ x: 0, y: y4 });

    var index = new Array(SIZE);
    for (var i = 0; i < SIZE; i++) index[i] = new Array(SIZE);
    for (var k = 0; k < out.length; k++) index[out[k].y][out[k].x] = k;

    return { cells: out, index: index };
}

// Tier 1 - greedy: always take the shortest path to the food, full stop.
// No safety awareness at all. Dies as soon as the shortest path happens to
// require the snake to cross where its own body will be by the time it
// gets there.
function nextMoveGreedy() {
    var head = snake[0];
    var pathToFood = shortestPath(occupied, head, item, false);
    if (pathToFood && pathToFood.length > 1) return pathToFood[1];
    var free = neighbors4(head.x, head.y).filter(function (n) { return !occupied[n.y][n.x]; });
    return free.length > 0 ? free[0] : null;
}

// Tier 2 - greedy + tail safety: take the shortest path to food, but only if
// a one-step lookahead confirms the snake could still reach its own tail
// afterwards. Otherwise follow the tail to stay alive in open space, and
// failing that, pick whichever free neighbour leaves the most open,
// wall-clear floor space. Survives much longer than pure greedy, but the
// lookahead is fundamentally short-sighted and still occasionally coils
// into a dead end a few steps later.
function nextMoveGreedySafe() {
    var head = snake[0];
    var tail = snake[snake.length - 1];

    var pathToFood = shortestPath(occupied, head, item, false);
    if (pathToFood && pathToFood.length > 1 && isMoveSafe(pathToFood[1])) {
        return pathToFood[1];
    }

    var pathToTail = shortestPath(occupied, head, tail, true);
    if (pathToTail && pathToTail.length > 1) return pathToTail[1];

    var freeNeighbors = neighbors4(head.x, head.y).filter(function (n) { return !occupied[n.y][n.x]; });
    if (freeNeighbors.length === 0) return null;

    var safeNeighbors = freeNeighbors.filter(isMoveSafe);
    var pool = safeNeighbors.length > 0 ? safeNeighbors : freeNeighbors;

    var best = null;
    var bestScore = -1;
    for (var i = 0; i < pool.length; i++) {
        var s = floodFillCount(occupied, pool[i]) * 4 + wallDistance(pool[i]);
        if (s > bestScore) {
            bestScore = s;
            best = pool[i];
        }
    }
    return best;
}

// Tier 3 - cycle + shortcuts: follow the fixed Hamiltonian cycle by default,
// but jump straight to the food whenever a one-step-safe shortcut is
// available. Much faster than pure cycle-following since it beelines for
// food instead of waiting for the cycle to sweep past it - but the
// shortcuts are arbitrary BFS routes, not moves that stay consistent with
// the cycle's own direction, so the one-step safety check can't fully
// vouch for them either. Rare, but it can still end up boxed in.
function nextMoveCycleShortcuts() {
    var head = snake[0];

    var pathToFood = shortestPath(occupied, head, item, false);
    if (pathToFood && pathToFood.length > 1 && isMoveSafe(pathToFood[1])) {
        return pathToFood[1];
    }

    var headIdx = cycleIndex[head.y][head.x];
    var cycleNext = cycle[(headIdx + 1) % cycle.length];
    if (!occupied[cycleNext.y][cycleNext.x]) {
        return cycleNext;
    }

    var freeNeighbors = neighbors4(head.x, head.y).filter(function (n) { return !occupied[n.y][n.x]; });
    if (freeNeighbors.length === 0) return null;

    var safeNeighbors = freeNeighbors.filter(isMoveSafe);
    var pool = safeNeighbors.length > 0 ? safeNeighbors : freeNeighbors;

    var best = null;
    var bestScore = -1;
    for (var i = 0; i < pool.length; i++) {
        var s = floodFillCount(occupied, pool[i]) * 4 + wallDistance(pool[i]);
        if (s > bestScore) {
            bestScore = s;
            best = pool[i];
        }
    }
    return best;
}

// Tier 4 - cycle (pure), the default: always follow the fixed cycle to the
// next cell after the head, full stop. No shortcuts, no heuristics, nothing
// to get subtly wrong - the body is always a contiguous, non-crossing arc
// of the same loop, so the next cycle cell is guaranteed free right up
// until the snake fills the entire board. It never beelines for food, so
// it's slower to grow, but it's the only tier that's actually provable
// rather than merely well-tested.
function nextMoveCyclePure() {
    var head = snake[0];
    var headIdx = cycleIndex[head.y][head.x];
    return cycle[(headIdx + 1) % cycle.length];
}

function nextMove() {
    switch (strategy) {
        case 'greedy': return nextMoveGreedy();
        case 'greedy_safe': return nextMoveGreedySafe();
        case 'cycle_shortcuts': return nextMoveCycleShortcuts();
        default: return nextMoveCyclePure();
    }
}

function placeItem() {
    do {
        item = { x: Math.floor(Math.random() * SIZE), y: Math.floor(Math.random() * SIZE) };
    } while (occupied[item.y][item.x]);
}

function resetGame() {
    gameOver = false;
    score = 0;
    occupied = new Array(SIZE);
    for (var i = 0; i < SIZE; i++) occupied[i] = new Array(SIZE).fill(false);

    var cx = Math.floor(SIZE / 2);
    var cy = Math.floor(SIZE / 2);
    snake = [{ x: cx, y: cy }];
    occupied[cy][cx] = true;
    ticksSinceProgress = 0;
    placeItem();
}

function onDeath() {
    deaths++;
    gameOver = true;
}

function setup() {
    createCanvas(SIZE * BLOCK, SIZE * BLOCK).parent('ai-snake-game');
    noSmooth();
    speed = document.getElementById('ai-snake-speed');
    strategySelect = document.getElementById('ai-snake-strategy');
    strategy = strategySelect.value;

    strategySelect.addEventListener('change', function () {
        strategy = strategySelect.value;
        deaths = 0;
        resetGame();
    });

    var built = buildCycle();
    cycle = built.cells;
    cycleIndex = built.index;
    resetGame();
}

function tick() {
    // The pure-cycle tier never needs this - it's structurally guaranteed to
    // keep making progress. The other tiers can get stuck circling forever
    // without ever eating: food being BFS-reachable in principle doesn't
    // mean the safety check will ever actually let the snake take that
    // path, so a "is food reachable" check can stay permanently satisfied
    // while the snake just orbits its own tail loop, never once getting
    // closer. Tracking real progress (ticks since the score last went up)
    // instead of reachability is what actually catches that.
    if (strategy !== 'cycle_pure') {
        ticksSinceProgress++;
        if (ticksSinceProgress > STUCK_LIMIT) {
            onDeath();
            return;
        }
    }

    var next = nextMove();
    if (!next) {
        onDeath();
        return;
    }

    var isEating = next.x === item.x && next.y === item.y;
    var tail = snake[snake.length - 1];
    var isTailVacate = !isEating && next.x === tail.x && next.y === tail.y;

    // Hard invariant, enforced at the single point every move gets
    // committed: never move onto an occupied cell unless it's the tail
    // slot that's vacating this same tick. Whatever chose `next` upstream
    // is expected to already guarantee this - this is the backstop that
    // turns any bug in that selection into a clean game over instead of
    // the head silently passing through the body.
    if (occupied[next.y][next.x] && !isTailVacate) {
        onDeath();
        return;
    }

    snake.unshift(next);
    occupied[next.y][next.x] = true;

    if (isEating) {
        score++;
        ticksSinceProgress = 0;
        placeItem();
    } else {
        var removed = snake.pop();
        occupied[removed.y][removed.x] = false;
    }
}

function draw() {
    background(51);

    if (gameOver) {
        textAlign(CENTER, CENTER);
        fill(255);
        textSize(28);
        text('Game Over - Score ' + score, width / 2, height / 2);
        textAlign(LEFT, BASELINE);
        setTimeout(resetGame, 1200);
        return;
    }

    var speedVal = Number(speed.value);
    frameRate(Math.min(speedVal, 60));
    var stepsPerFrame = Math.max(1, Math.floor(speedVal / 60));
    for (var s = 0; s < stepsPerFrame && !gameOver; s++) tick();

    noStroke();
    fill(255, 0, 0);
    rect(BLOCK * item.x, BLOCK * item.y, BLOCK, BLOCK);

    // A line traced through consecutive segment centers, in snake[] order
    // (not grid adjacency), makes the actual connectivity of the body
    // visible - two cells that are merely next to each other but not
    // actually consecutive in the body never get a line between them.
    // The head-to-tail color gradient shows which end is which along it.
    noFill();
    strokeWeight(BLOCK * 0.3);
    strokeCap(ROUND);
    for (var i = 0; i < snake.length - 1; i++) {
        var t = snake.length > 1 ? i / (snake.length - 1) : 0;
        stroke(lerpColor(color(120, 200, 120), color(90, 110, 220), t));
        line(
            BLOCK * snake[i].x + BLOCK / 2, BLOCK * snake[i].y + BLOCK / 2,
            BLOCK * snake[i + 1].x + BLOCK / 2, BLOCK * snake[i + 1].y + BLOCK / 2
        );
    }

    strokeWeight(1);
    stroke(30);
    for (var j = 0; j < snake.length; j++) {
        var tj = snake.length > 1 ? j / (snake.length - 1) : 0;
        fill(j === 0 ? color(120, 200, 120) : lerpColor(color(120, 200, 120), color(90, 110, 220), tj));
        rect(BLOCK * snake[j].x + BLOCK * 0.15, BLOCK * snake[j].y + BLOCK * 0.15, BLOCK * 0.7, BLOCK * 0.7);
    }

    fill(255);
    textSize(20);
    text('Score ' + score + '   Deaths ' + deaths, 8, 8, 300, 30);
}
