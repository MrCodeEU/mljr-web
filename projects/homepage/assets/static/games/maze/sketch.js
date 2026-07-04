var COLS = 30, ROWS = 20;
var CELL = 26;
var W = COLS * CELL, H = ROWS * CELL;

var grid = [];
var stack = [];
var current;
var genDone = false;

var solveVisited, solveQueue, solveParent, solveDone, solvePath, solveIndex;

var speedInput;

function Cell(col, row) {
    this.col = col;
    this.row = row;
    this.walls = { top: true, right: true, bottom: true, left: true };
    this.visited = false;
}

function idx(col, row) {
    if (col < 0 || row < 0 || col >= COLS || row >= ROWS) return -1;
    return col + row * COLS;
}

function neighborsOf(cell, requireUnvisited) {
    var out = [];
    var dirs = [
        { dc: 0, dr: -1, self: 'top', other: 'bottom' },
        { dc: 1, dr: 0, self: 'right', other: 'left' },
        { dc: 0, dr: 1, self: 'bottom', other: 'top' },
        { dc: -1, dr: 0, self: 'left', other: 'right' },
    ];
    for (var i = 0; i < dirs.length; i++) {
        var d = dirs[i];
        var ni = idx(cell.col + d.dc, cell.row + d.dr);
        if (ni === -1) continue;
        var n = grid[ni];
        if (requireUnvisited && n.visited) continue;
        out.push({ cell: n, dir: d });
    }
    return out;
}

// Randomized depth-first search ("recursive backtracker"): carve a passage
// to a random unvisited neighbour, push the old cell so we can backtrack to
// it once the new branch dead-ends. Run one step per call so the caller can
// animate the carving instead of generating the whole maze instantly.
function generateStep() {
    if (genDone) return;
    var options = neighborsOf(current, true);
    if (options.length > 0) {
        var choice = options[Math.floor(Math.random() * options.length)];
        current.walls[choice.dir.self] = false;
        choice.cell.walls[choice.dir.other] = false;
        choice.cell.visited = true;
        stack.push(current);
        current = choice.cell;
    } else if (stack.length > 0) {
        current = stack.pop();
    } else {
        genDone = true;
        startSolve();
    }
}

function startSolve() {
    solveVisited = new Array(grid.length).fill(false);
    solveParent = new Array(grid.length).fill(-1);
    solveQueue = [0];
    solveVisited[0] = true;
    solveDone = false;
    solvePath = [];
    solveIndex = 0;
}

// Breadth-first search over the carved passages guarantees the shortest
// route from the top-left cell to the bottom-right one.
function solveStep() {
    if (solveDone) return;
    if (solveQueue.length === 0) {
        solveDone = true;
        return;
    }
    var ci = solveQueue.shift();
    var goal = grid.length - 1;
    if (ci === goal) {
        var path = [ci];
        var p = solveParent[ci];
        while (p !== -1) {
            path.push(p);
            p = solveParent[p];
        }
        solvePath = path;
        solveDone = true;
        return;
    }
    var cell = grid[ci];
    var passable = neighborsOf(cell, false).filter(function (n) {
        return !cell.walls[n.dir.self];
    });
    for (var i = 0; i < passable.length; i++) {
        var n = passable[i];
        var ni = idx(n.cell.col, n.cell.row);
        if (!solveVisited[ni]) {
            solveVisited[ni] = true;
            solveParent[ni] = ci;
            solveQueue.push(ni);
        }
    }
}

function newMaze() {
    grid = [];
    for (var r = 0; r < ROWS; r++) {
        for (var c = 0; c < COLS; c++) {
            grid.push(new Cell(c, r));
        }
    }
    grid[0].visited = true;
    current = grid[0];
    stack = [];
    genDone = false;
    solveDone = true;
    solveVisited = [];
    solvePath = [];
}

function setup() {
    var canvas = createCanvas(W, H);
    canvas.parent('maze-game');
    speedInput = document.getElementById('maze-speed');
    document.getElementById('maze-regenerate').addEventListener('click', newMaze);
    newMaze();
}

function drawMaze() {
    stroke(200);
    strokeWeight(2);
    for (var i = 0; i < grid.length; i++) {
        var cell = grid[i];
        var x = cell.col * CELL, y = cell.row * CELL;
        if (cell.walls.top) line(x, y, x + CELL, y);
        if (cell.walls.right) line(x + CELL, y, x + CELL, y + CELL);
        if (cell.walls.bottom) line(x, y + CELL, x + CELL, y + CELL);
        if (cell.walls.left) line(x, y, x, y + CELL);
        if (cell.visited && !genDone) {
            noStroke();
            fill(90, 60, 60);
            rect(x + 1, y + 1, CELL - 2, CELL - 2);
            stroke(200);
        }
    }
}

function draw() {
    background(20);

    var speed = Number(speedInput.value);
    if (!genDone) {
        for (var s = 0; s < speed; s++) generateStep();
    } else if (!solveDone) {
        for (var s2 = 0; s2 < speed; s2++) solveStep();
    }

    drawMaze();

    if (!genDone) {
        noStroke();
        fill(255, 210, 90);
        circle(current.col * CELL + CELL / 2, current.row * CELL + CELL / 2, CELL * 0.5);
    } else if (solveVisited) {
        noStroke();
        fill(90, 140, 230, 90);
        for (var v = 0; v < solveVisited.length; v++) {
            if (solveVisited[v]) {
                var c = grid[v];
                rect(c.col * CELL + 4, c.row * CELL + 4, CELL - 8, CELL - 8);
            }
        }
        if (solveDone && solvePath.length > 0) {
            noFill();
            stroke(255, 120, 150);
            strokeWeight(4);
            beginShape();
            for (var p = solvePath.length - 1; p >= 0; p--) {
                var pc = grid[solvePath[p]];
                vertex(pc.col * CELL + CELL / 2, pc.row * CELL + CELL / 2);
            }
            endShape();
        }
    }

    noStroke();
    fill(120, 220, 140);
    circle(CELL / 2, CELL / 2, CELL * 0.4);
    fill(230, 90, 100);
    circle(W - CELL / 2, H - CELL / 2, CELL * 0.4);
}
