var SIZE = 50;
var BLOCK = 14;
var gameOver = false;
var score = 0;
var speed;

var occupied;
var snake;
var item;

function inBounds(x, y) {
    return x >= 0 && y >= 0 && x < SIZE && y < SIZE;
}

function neighbors4(x, y) {
    var out = [];
    if (inBounds(x - 1, y)) out.push({ x: x - 1, y: y });
    if (inBounds(x + 1, y)) out.push({ x: x + 1, y: y });
    if (inBounds(x, y - 1)) out.push({ x: x, y: y - 1 });
    if (inBounds(x, y + 1)) out.push({ x: x, y: y + 1 });
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
// room does this direction leave me" heuristic when no path to food or tail
// exists.
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

// Simulates the snake sliding along `path` to the food (growing only on the
// final step), returning the resulting occupied grid + body so we can check
// whether the snake would still be able to reach its own tail afterwards.
function simulateEatPath(occ, body, path) {
    var occCopy = cloneOccupied(occ);
    var bodyCopy = body.slice();
    for (var i = 1; i < path.length; i++) {
        var head = path[i];
        bodyCopy.unshift(head);
        occCopy[head.y][head.x] = true;
        if (i < path.length - 1) {
            var tail = bodyCopy.pop();
            occCopy[tail.y][tail.x] = false;
        }
    }
    return { occupied: occCopy, body: bodyCopy };
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
    placeItem();
}

function setup() {
    createCanvas(SIZE * BLOCK, SIZE * BLOCK).parent('ai-snake-game');
    speed = document.getElementById('ai-snake-speed');
    resetGame();
}

// Decide the next cell for the snake head:
//  1. shortest path to food, but only if the snake can still reach its own
//     tail after eating (otherwise it might seal itself into a dead end)
//  2. otherwise follow the tail to stay alive in open space
//  3. otherwise pick whichever free neighbour leaves the most open room
function nextMove() {
    var head = snake[0];
    var tail = snake[snake.length - 1];

    var pathToFood = shortestPath(occupied, head, item, false);
    if (pathToFood && pathToFood.length > 1) {
        var sim = simulateEatPath(occupied, snake, pathToFood);
        var newHead = sim.body[0];
        var newTail = sim.body[sim.body.length - 1];
        var stillSafe = sim.body.length <= 1 || shortestPath(sim.occupied, newHead, newTail, true) !== null;
        if (stillSafe) return pathToFood[1];
    }

    var pathToTail = shortestPath(occupied, head, tail, true);
    if (pathToTail && pathToTail.length > 1) return pathToTail[1];

    var candidates = neighbors4(head.x, head.y).filter(function (n) { return !occupied[n.y][n.x]; });
    var best = null;
    var bestScore = -1;
    for (var i = 0; i < candidates.length; i++) {
        var s = floodFillCount(occupied, candidates[i]);
        if (s > bestScore) {
            bestScore = s;
            best = candidates[i];
        }
    }
    return best;
}

function tick() {
    var next = nextMove();
    if (!next) {
        gameOver = true;
        return;
    }

    snake.unshift(next);
    occupied[next.y][next.x] = true;

    if (next.x === item.x && next.y === item.y) {
        score++;
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
    for (var i = 0; i < snake.length; i++) {
        fill(i === 0 ? color(120, 200, 120) : 255);
        rect(BLOCK * snake[i].x, BLOCK * snake[i].y, BLOCK, BLOCK);
    }

    fill(255);
    textSize(20);
    text(score, 8, 8, 100, 30);
}
