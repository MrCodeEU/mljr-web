package main

import (
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	"mljr-web/internal/snakeai"
	"mljr-web/projects/homepage/pages"
)

// snakeRoom is the one shared game. It stays idle (no game, no model in memory)
// until a viewer connects and stops shortly after the last one leaves.
var snakeRoom = snakeai.NewRoom(snakeai.Config{
	N:          30,
	Tick:       snakeTick(),
	Grace:      20 * time.Second,
	MaxViewers: 50,
})

// registerSnakeLive mounts the SSE feed. It lives under /api/, which the shared
// gzip middleware skips, so frames are never buffered.
func registerSnakeLive(e *echo.Echo) {
	e.GET(pages.AISnakeStreamPath, echo.WrapHandler(snakeai.StreamHandler(snakeRoom)))
}

// snakeTick is the time between moves: 40 ms by default (about 25 moves/s; a move costs ~0.3 ms),
// overridable with SNAKE_AI_TICK_MS (10-1000).
func snakeTick() time.Duration {
	if ms, err := strconv.Atoi(os.Getenv("SNAKE_AI_TICK_MS")); err == nil && ms >= 10 && ms <= 1000 {
		return time.Duration(ms) * time.Millisecond
	}
	return 40 * time.Millisecond
}
