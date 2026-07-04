(function () {
    var SIZE = 20;
    var x = 0;
    var y = 0;
    var Minenfeld = [];
    var Moved = [];
    var row = 0;
    var Input = NaN;
    var lost = false;

    function newField() {
        Minenfeld = [];
        Moved = [];
        row = 0;
        Input = NaN;
        lost = false;
        for (var i = 0; i < SIZE; i++) {
            Minenfeld[i] = [];
            for (var j = 0; j < SIZE; j++) {
                Minenfeld[i][j] = false;
            }
        }
        for (var i = 0; i < SIZE; i++) {
            Minenfeld[i][Math.round(Math.random() * (SIZE - 1))] = true;
        }
    }
    newField();

    var minenfelCanvas = document.getElementById("minenfeld-canvas");
    var context = minenfelCanvas.getContext("2d");

    var inputText = document.getElementById("minenfeld-input-label");
    var inputBox = document.getElementById("minenfeld-input-box");
    var inputBtn = document.getElementById("minenfeld-input-btn");
    var inputWrap = document.getElementById("minenfeld-input-wrap");
    var message = document.getElementById("minenfeld-message");
    var fehler = document.getElementById("minenfeld-error");
    var restartBtn = document.getElementById("minenfeld-restart");

    function showError(text) {
        fehler.style.display = "inline";
        fehler.textContent = text;
        message.style.display = "none";
    }

    function showMessage(html) {
        fehler.style.display = "none";
        message.style.display = "inline";
        message.innerHTML = html;
    }

    function startMove(step) {
        if (Minenfeld[row][step]) {
            lost = true;
            Moved[row] = step;
            DrawWay();
            Lost();
        } else if (row === SIZE - 1) {
            Moved[row] = step;
            DrawWay();
            Win();
        } else {
            Moved[row] = step;
            DrawWay();
            row++;
        }
    }

    function submitStart() {
        var value = Number(inputBox.value);
        if (value >= SIZE) {
            showError("Die Zahl ist zu groß!");
            return;
        }
        if (value <= -1) {
            showError("Die Zahl ist zu klein!");
            return;
        }
        if (isNaN(value)) {
            showError("Geben Sie bitte eine Zahl ein");
            return;
        }
        showMessage("Sie bewegen sich durch das Minenfeld mit den Pfeiltasten.<br />Nach links[←]/nach unten[↓]/nach rechts[→]");
        inputWrap.style.display = "none";
        Input = value;
        startMove(value);
    }

    inputBtn.onclick = submitStart;
    inputBox.addEventListener("keydown", function (e) {
        if (e.key === "Enter") {
            e.preventDefault();
            submitStart();
        }
    });

    document.addEventListener("keydown", function (e) {
        if (e.key === "ArrowLeft" || e.key === "ArrowRight" || e.key === "ArrowDown" || e.key === "ArrowUp") {
            e.preventDefault();
        }
        if (isNaN(Input) || row > SIZE - 1 || lost) return;
        var next = Input;
        if (e.key === "ArrowLeft" && Input >= 1) {
            next = Input - 1;
        } else if (e.key === "ArrowDown") {
            next = Input;
        } else if (e.key === "ArrowRight" && Input <= SIZE - 2) {
            next = Input + 1;
        } else {
            return;
        }
        Input = next;
        startMove(next);
    });

    function DrawMinefield(rowtodraw) {
        context.strokeStyle = "black";
        context.lineWidth = 2;
        x = minenfelCanvas.width / SIZE;
        y = minenfelCanvas.height / SIZE;
        for (var i = 0; i <= rowtodraw; i++) {
            for (var j = 0; j < SIZE; j++) {
                context.fillStyle = Minenfeld[i][j] ? "red" : "lawngreen";
                context.fillRect(x * j, y * i, x, y);
                context.strokeRect(x * j, y * i, x, y);
            }
        }
    }

    function DrawWay() {
        var lineX = minenfelCanvas.width / SIZE;
        var lineY = minenfelCanvas.height / SIZE;
        context.clearRect(0, 0, minenfelCanvas.width, minenfelCanvas.height);
        DrawMinefield(row);

        context.strokeStyle = "blue";
        context.lineWidth = 4;
        context.lineCap = "round";
        context.beginPath();
        context.moveTo(lineX * Moved[0] + lineX / 2, 0);
        for (var i = 0; i <= row; i++) {
            context.lineTo(lineX * Moved[i] + lineX / 2, lineY * i + lineY / 2);
        }
        context.stroke();

        context.beginPath();
        context.lineWidth = 4;
        context.lineCap = "round";
        if (lost) {
            context.strokeStyle = "red";
            context.arc(lineX * Input + lineX / 2, lineY * row + lineY / 2, 17, 1.5 * Math.PI, 3.5 * Math.PI);
        } else if (row === SIZE - 1) {
            context.strokeStyle = "green";
            context.arc(lineX * Input + lineX / 2, lineY * row + lineY / 2, 5, 1.5 * Math.PI, 3.5 * Math.PI);
        }
        context.stroke();
    }

    function Lost() {
        showMessage("Boom!! Sie haben eine Bombe erwischt!");
        restartBtn.style.display = "inline-block";
    }

    function Win() {
        showMessage("Gratulation!! Sie haben es durch das Minenfeld geschafft!");
        restartBtn.style.display = "inline-block";
        lost = true;
    }

    restartBtn.onclick = function () {
        newField();
        inputBox.value = "";
        inputWrap.style.display = "";
        restartBtn.style.display = "none";
        message.style.display = "none";
        fehler.style.display = "none";
        context.clearRect(0, 0, minenfelCanvas.width, minenfelCanvas.height);
        DrawMinefield(SIZE - 1);
    };

    DrawMinefield(SIZE - 1);
})();
