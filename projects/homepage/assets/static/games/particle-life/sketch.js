var W = 900, H = 600;
var TYPES = 5;
var COLORS = [
    [230, 80, 90], [90, 200, 120], [90, 150, 230], [230, 200, 90], [200, 100, 230]
];

var particles = [];
var rules = [];
var maxRadius = 80;
var collisionRadius = 12;
var maxForce = 0.6;
var friction = 0.62;

var countInput, forceInput, radiusInput;

function randomRules() {
    rules = [];
    for (var i = 0; i < TYPES; i++) {
        rules.push([]);
        for (var j = 0; j < TYPES; j++) {
            rules[i].push(Math.random() * 2 - 1);
        }
    }
}

// The standard particle-life force curve: hard repulsion at very close
// range so particles never fully overlap, then the (possibly negative)
// per-type-pair rule takes over as a smooth attract/repel bump peaking at
// the midpoint of the interaction radius, fading to zero at the edge.
function forceAt(r, a) {
    var beta = collisionRadius / maxRadius;
    if (r < beta) return r / beta - 1;
    if (r < 1) return a * (1 - Math.abs(2 * r - 1 - beta) / (1 - beta));
    return 0;
}

function spawnParticles() {
    var count = Number(countInput.value);
    particles = [];
    for (var i = 0; i < count; i++) {
        particles.push({
            x: Math.random() * W,
            y: Math.random() * H,
            vx: 0,
            vy: 0,
            type: Math.floor(Math.random() * TYPES),
        });
    }
}

function setup() {
    var canvas = createCanvas(W, H);
    canvas.parent('particle-life-game');
    countInput = document.getElementById('particle-life-count');
    forceInput = document.getElementById('particle-life-force');
    radiusInput = document.getElementById('particle-life-radius');

    document.getElementById('particle-life-randomize').addEventListener('click', function () {
        randomRules();
    });
    document.getElementById('particle-life-reset').addEventListener('click', spawnParticles);

    randomRules();
    spawnParticles();
}

function step() {
    maxForce = Number(forceInput.value);
    maxRadius = Number(radiusInput.value);

    for (var i = 0; i < particles.length; i++) {
        var p = particles[i];
        var fx = 0, fy = 0;
        for (var j = 0; j < particles.length; j++) {
            if (i === j) continue;
            var q = particles[j];
            var dx = q.x - p.x, dy = q.y - p.y;
            var d = Math.sqrt(dx * dx + dy * dy);
            if (d === 0 || d > maxRadius) continue;
            var r = d / maxRadius;
            var f = forceAt(r, rules[p.type][q.type]) * maxForce;
            fx += (dx / d) * f;
            fy += (dy / d) * f;
        }
        p.vx = (p.vx + fx) * friction;
        p.vy = (p.vy + fy) * friction;
    }

    for (var k = 0; k < particles.length; k++) {
        var pt = particles[k];
        pt.x += pt.vx;
        pt.y += pt.vy;
        if (pt.x < 0) { pt.x = 0; pt.vx *= -1; }
        if (pt.x > W) { pt.x = W; pt.vx *= -1; }
        if (pt.y < 0) { pt.y = 0; pt.vy *= -1; }
        if (pt.y > H) { pt.y = H; pt.vy *= -1; }
    }
}

function draw() {
    background(15);
    step();
    noStroke();
    for (var i = 0; i < particles.length; i++) {
        var p = particles[i];
        var c = COLORS[p.type];
        fill(c[0], c[1], c[2]);
        circle(p.x, p.y, 6);
    }
}
