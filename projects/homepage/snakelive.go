package main

import (
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	"mljr-web/internal/snakesim"
	"mljr-web/projects/homepage/pages"
)

// snakeRoom is the one shared game. It stays idle (no game, no model in memory)
// until a viewer connects and stops shortly after the last one leaves.
var snakeRoom = snakesim.NewRoom(snakesim.Config{
	N:          50,
	Speed:      snakeSpeed(),
	Grace:      20 * time.Second,
	MaxViewers: 50,
})

// registerSnakeLive mounts the SSE feed. It lives under /api/, which the shared
// gzip middleware skips, so frames are never buffered.
func registerSnakeLive(e *echo.Echo) {
	e.GET(pages.AISnakeStreamPath, echo.WrapHandler(snakesim.StreamHandler(snakeRoom)))
}

// snakeSpeed is the number of moves per second, timed against the real clock: 30 by default,
// overridable with SNAKE_AI_MPS (5-100).
func snakeSpeed() float64 {
	if v, err := strconv.ParseFloat(os.Getenv("SNAKE_AI_MPS"), 64); err == nil && v >= 5 && v <= 100 {
		return v
	}
	return 30
}
