var W = 900, H = 600;
var PERCEPTION = 50;
var MAX_SPEED = 3.2;
var MAX_FORCE = 0.15;

var boids = [];
var sepInput, aliInput, cohInput, countInput;

function Boid() {
    this.pos = createVector(Math.random() * W, Math.random() * H);
    this.vel = p5.Vector.random2D().mult(Math.random() * 2 + 1);
    this.acc = createVector(0, 0);
}

Boid.prototype.edges = function () {
    if (this.pos.x > W) this.pos.x = 0;
    if (this.pos.x < 0) this.pos.x = W;
    if (this.pos.y > H) this.pos.y = 0;
    if (this.pos.y < 0) this.pos.y = H;
};

// The three classic Reynolds rules, each computed over boids within
// PERCEPTION radius: steer away from close neighbours (separation), match
// their average heading (alignment), and drift toward their average
// position (cohesion). Flocking is what falls out of blending these three
// simple, purely local rules - no boid knows about the flock as a whole.
Boid.prototype.flock = function (others, weights) {
    var steerSep = createVector(0, 0);
    var steerAli = createVector(0, 0);
    var steerCoh = createVector(0, 0);
    var totalSep = 0, totalAli = 0, totalCoh = 0;

    for (var i = 0; i < others.length; i++) {
        var other = others[i];
        if (other === this) continue;
        var d = p5.Vector.dist(this.pos, other.pos);
        if (d === 0 || d > PERCEPTION) continue;

        var diff = p5.Vector.sub(this.pos, other.pos);
        diff.div(d * d);
        steerSep.add(diff);
        totalSep++;

        steerAli.add(other.vel);
        totalAli++;

        steerCoh.add(other.pos);
        totalCoh++;
    }

    if (totalSep > 0) {
        steerSep.div(totalSep);
        steerSep.setMag(MAX_SPEED);
        steerSep.sub(this.vel);
        steerSep.limit(MAX_FORCE);
    }
    if (totalAli > 0) {
        steerAli.div(totalAli);
        steerAli.setMag(MAX_SPEED);
        steerAli.sub(this.vel);
        steerAli.limit(MAX_FORCE);
    }
    if (totalCoh > 0) {
        steerCoh.div(totalCoh);
        steerCoh.sub(this.pos);
        steerCoh.setMag(MAX_SPEED);
        steerCoh.sub(this.vel);
        steerCoh.limit(MAX_FORCE);
    }

    steerSep.mult(weights.sep);
    steerAli.mult(weights.ali);
    steerCoh.mult(weights.coh);

    this.acc.add(steerSep);
    this.acc.add(steerAli);
    this.acc.add(steerCoh);
};

Boid.prototype.update = function () {
    this.pos.add(this.vel);
    this.vel.add(this.acc);
    this.vel.limit(MAX_SPEED);
    this.acc.mult(0);
};

Boid.prototype.render = function () {
    var angle = this.vel.heading() + Math.PI / 2;
    push();
    translate(this.pos.x, this.pos.y);
    rotate(angle);
    noStroke();
    fill(120, 190, 255);
    triangle(0, -8, -5, 6, 5, 6);
    pop();
};

function rebuildFlock() {
    var count = Number(countInput.value);
    boids = [];
    for (var i = 0; i < count; i++) boids.push(new Boid());
}

function setup() {
    var canvas = createCanvas(W, H);
    canvas.parent('boids-game');
    sepInput = document.getElementById('boids-separation');
    aliInput = document.getElementById('boids-alignment');
    cohInput = document.getElementById('boids-cohesion');
    countInput = document.getElementById('boids-count');

    countInput.addEventListener('change', rebuildFlock);
    rebuildFlock();
}

function draw() {
    background(20);
    var weights = {
        sep: Number(sepInput.value),
        ali: Number(aliInput.value),
        coh: Number(cohInput.value),
    };
    for (var i = 0; i < boids.length; i++) {
        boids[i].flock(boids, weights);
    }
    for (var j = 0; j < boids.length; j++) {
        boids[j].update();
        boids[j].edges();
        boids[j].render();
    }
}
