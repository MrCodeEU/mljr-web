var MAX_SEGMENTS = 5;

var n = 3;
var l = [];
var m = [];
var g = 1;
var theta = [];
var omega = [];

var trail;
var originX, originY;
var renderScale = 1;
var reachRadius;
var speedInput;

// Builds the mass matrix A and force vector B for the generalized n-link
// pendulum (Lagrangian mechanics: A(theta) * theta'' = B(theta, omega)).
// M[i] is the combined mass of segment i through the end of the chain -
// it's the "weight this rod has to support", and it's what makes each
// additional link couple back into every other link's equation.
function computeAccel(th, om) {
    var M = new Array(n);
    M[n - 1] = m[n - 1];
    for (var i = n - 2; i >= 0; i--) M[i] = M[i + 1] + m[i];

    var A = [];
    var B = [];
    for (var i = 0; i < n; i++) {
        A.push(new Array(n).fill(0));
        var b = -g * l[i] * M[i] * Math.sin(th[i]);
        for (var j = 0; j < n; j++) {
            var Mij = M[Math.max(i, j)];
            A[i][j] = l[i] * l[j] * Mij * Math.cos(th[i] - th[j]);
            if (j !== i) {
                b -= l[i] * l[j] * Mij * Math.sin(th[i] - th[j]) * om[j] * om[j];
            }
        }
        B.push(b);
    }
    return solveLinear(A, B);
}

// Simple Gaussian elimination with partial pivoting - the matrices here are
// at most 5x5, so this is plenty fast and numerically fine.
function solveLinear(A, B) {
    var nn = B.length;
    var M = A.map(function (row) { return row.slice(); });
    var x = B.slice();

    for (var col = 0; col < nn; col++) {
        var pivot = col;
        for (var r = col + 1; r < nn; r++) {
            if (Math.abs(M[r][col]) > Math.abs(M[pivot][col])) pivot = r;
        }
        if (pivot !== col) {
            var tmpRow = M[col]; M[col] = M[pivot]; M[pivot] = tmpRow;
            var tmpX = x[col]; x[col] = x[pivot]; x[pivot] = tmpX;
        }
        for (var r2 = col + 1; r2 < nn; r2++) {
            var factor = M[r2][col] / M[col][col];
            for (var c2 = col; c2 < nn; c2++) M[r2][c2] -= factor * M[col][c2];
            x[r2] -= factor * x[col];
        }
    }
    var result = new Array(nn);
    for (var i2 = nn - 1; i2 >= 0; i2--) {
        var sum = x[i2];
        for (var j2 = i2 + 1; j2 < nn; j2++) sum -= M[i2][j2] * result[j2];
        result[i2] = sum / M[i2][i2];
    }
    return result;
}

function derivative(th, om) {
    return { dtheta: om.slice(), domega: computeAccel(th, om) };
}

// RK4 keeps a chaotic multi-link pendulum numerically stable far longer
// than symplectic Euler would - important once you go past 2 segments,
// since energy errors compound a lot faster with more coupled links.
function rk4Step(dt) {
    var k1 = derivative(theta, omega);
    var th2 = addScaled(theta, k1.dtheta, dt / 2);
    var om2 = addScaled(omega, k1.domega, dt / 2);
    var k2 = derivative(th2, om2);
    var th3 = addScaled(theta, k2.dtheta, dt / 2);
    var om3 = addScaled(omega, k2.domega, dt / 2);
    var k3 = derivative(th3, om3);
    var th4 = addScaled(theta, k3.dtheta, dt);
    var om4 = addScaled(omega, k3.domega, dt);
    var k4 = derivative(th4, om4);

    for (var i = 0; i < n; i++) {
        theta[i] += (dt / 6) * (k1.dtheta[i] + 2 * k2.dtheta[i] + 2 * k3.dtheta[i] + k4.dtheta[i]);
        omega[i] += (dt / 6) * (k1.domega[i] + 2 * k2.domega[i] + 2 * k3.domega[i] + k4.domega[i]);
    }
}

function addScaled(base, delta, factor) {
    var out = new Array(base.length);
    for (var i = 0; i < base.length; i++) out[i] = base[i] + delta[i] * factor;
    return out;
}

function readParams() {
    n = Number(document.getElementById('pendulum-segments').value);
    g = Number(document.getElementById('pendulum-gravity').value);
    l = [];
    m = [];
    for (var i = 0; i < n; i++) {
        l.push(Number(document.getElementById('pendulum-len-' + i).value));
        m.push(Number(document.getElementById('pendulum-mass-' + i).value));
    }
}

function syncSegmentRows() {
    for (var i = 0; i < MAX_SEGMENTS; i++) {
        var row = document.getElementById('pendulum-row-' + i);
        if (row) row.style.display = i < n ? '' : 'none';
    }
}

function resetPendulum() {
    readParams();
    syncSegmentRows();
    theta = [];
    omega = [];
    // Wide, uneven starting angles produce visibly chaotic motion right
    // away instead of a slow, near-periodic swing.
    for (var i = 0; i < n; i++) {
        theta.push(Math.PI / 2 + i * 0.35);
        omega.push(0);
    }
    trail = [];

    // The chain can (and with enough energy, will) swing all the way around
    // its pivot, so the reachable area is a full disk of radius = total arm
    // length. Scaling every segment's *drawn* length (not its physical
    // length used in the physics) so that disk always fits the canvas is
    // what keeps longer/more-segmented chains from swinging off-screen
    // instead of clipping at the edge.
    var totalLen = 0;
    for (var i = 0; i < n; i++) totalLen += l[i];
    renderScale = reachRadius / totalLen;
}

function bobPositions() {
    var pts = [];
    var x = originX, y = originY;
    for (var i = 0; i < n; i++) {
        x += l[i] * renderScale * Math.sin(theta[i]);
        y += l[i] * renderScale * Math.cos(theta[i]);
        pts.push({ x: x, y: y });
    }
    return pts;
}

function setup() {
    var canvas = createCanvas(1000, 700);
    canvas.parent('double-pendulum-game');
    originX = width / 2;
    originY = 260;
    reachRadius = Math.min(originX, width - originX, originY, height - originY) - 20;
    speedInput = document.getElementById('pendulum-speed');

    document.getElementById('pendulum-reset').addEventListener('click', resetPendulum);
    document.getElementById('pendulum-segments').addEventListener('change', resetPendulum);

    resetPendulum();
}

function draw() {
    background(30);

    var speedVal = Number(speedInput.value);
    var dt = 0.02;
    var steps = Math.max(1, Math.round(speedVal / 4));
    for (var s = 0; s < steps; s++) rk4Step(dt);

    var pts = bobPositions();
    var last = pts[pts.length - 1];
    trail.push({ x: last.x, y: last.y });
    if (trail.length > 800) trail.shift();

    noFill();
    for (var i = 1; i < trail.length; i++) {
        stroke(255, 80, 120, (i / trail.length) * 200);
        line(trail[i - 1].x, trail[i - 1].y, trail[i].x, trail[i].y);
    }

    stroke(255);
    strokeWeight(2);
    var px = originX, py = originY;
    for (var j = 0; j < pts.length; j++) {
        line(px, py, pts[j].x, pts[j].y);
        px = pts[j].x;
        py = pts[j].y;
    }

    noStroke();
    fill(255);
    circle(originX, originY, 6);
    for (var k = 0; k < pts.length; k++) {
        fill(k === pts.length - 1 ? color(255, 120, 150) : color(120, 180, 255));
        circle(pts[k].x, pts[k].y, Math.max(8, Math.min(28, m[k])));
    }
}
