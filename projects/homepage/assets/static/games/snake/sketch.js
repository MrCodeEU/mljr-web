var w = 20;
var size = 34;
var snake = [];
var direction = 'right';
var lastx1 = 0;
var lasty1 = 0;
var lastx2 = 0;
var lasty2 = 0;
var miliseconds;
var gamestate = 'go';
var score = 0;
var apple = {
    x: 0,
    y: 0
}

function resetGame() {
    direction = 'right';
    gamestate = 'go';
    score = 0;
    snake = [];
    snake[0] = { x: 7 * w, y: 5 * w };
    snake[1] = { x: 6 * w, y: 5 * w };
    snake[2] = { x: 5 * w, y: 5 * w };
    apple = {
        x: round(size / 2 * w),
        y: round(size / 2 * w)
    };
}

function setup() {
    createCanvas(w * size, w * size).parent('snake-game');
    miliseconds = millis();
    resetGame();
}

function draw() {
    background(51);

    for (let index = 0; index < snake.length; index++) {
        const element = snake[index];
        fill(index === 0 ? color(120, 200, 120) : 255);
        rect(element.x, element.y, w, w);
    }
    fill(color(255, 0, 0));
    rect(apple.x, apple.y, w, w);

    if (gamestate === 'go' && miliseconds <= millis() - 80) {
        miliseconds = millis();
        for (let index = 0; index < snake.length; index++) {
            const element = snake[index];
            if (index === 0) {
                lastx1 = element.x;
                lasty1 = element.y;
                switch (direction) {
                    case 'right':
                        element.x = element.x + w;
                        break;
                    case 'left':
                        element.x = element.x - w;
                        break;
                    case 'up':
                        element.y = element.y - w;
                        break;
                    case 'down':
                        element.y = element.y + w;
                        break;
                }
                for (let i = 1; i < snake.length; i++) {
                    if (snake[i].x === snake[0].x && snake[i].y === snake[0].y) {
                        gamestate = 'over';
                    }
                }
                if (element.x >= width || element.x < 0 || element.y >= height || element.y < 0) {
                    gamestate = 'over';
                }
                if (apple.x === element.x && apple.y === element.y) {
                    score++;
                    snake.push({ x: lastx2, y: lasty2 });
                    do {
                        var collides = false;
                        apple = {
                            x: round(random(0, size - 1)) * w,
                            y: round(random(0, size - 1)) * w
                        }
                        for (let i = 0; i < snake.length; i++) {
                            if (apple.x === snake[i].x && apple.y === snake[i].y) collides = true;
                        }
                    } while (collides);
                }
            } else {
                lastx2 = element.x;
                lasty2 = element.y;
                element.x = lastx1;
                element.y = lasty1;
                lastx1 = lastx2;
                lasty1 = lasty2;
            }
        }
    }

    textSize(32);
    fill(255);
    text(score, 10, 10, 100, 50);

    if (gamestate === 'over') {
        fill(0, 0, 0, 180);
        rect(0, 0, width, height);
        textAlign(CENTER, CENTER);
        fill(255);
        textSize(40);
        text('Game Over', width / 2, height / 2 - 30);
        textSize(20);
        text('Score: ' + score, width / 2, height / 2 + 10);
        text('Press Space to restart', width / 2, height / 2 + 40);
        textAlign(LEFT, BASELINE);
    }
}

function keyPressed() {
    if (gamestate === 'over') {
        if (keyCode === 32) resetGame();
    } else {
        switch (keyCode) {
            case UP_ARROW:
                if (direction !== 'down') direction = 'up';
                break;
            case DOWN_ARROW:
                if (direction !== 'up') direction = 'down';
                break;
            case LEFT_ARROW:
                if (direction !== 'right') direction = 'left';
                break;
            case RIGHT_ARROW:
                if (direction !== 'left') direction = 'right';
                break;
        }
    }
    if (keyCode === UP_ARROW || keyCode === DOWN_ARROW || keyCode === LEFT_ARROW || keyCode === RIGHT_ARROW || keyCode === 32) {
        return false;
    }
}
