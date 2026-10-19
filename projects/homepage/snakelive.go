package main

import (
	"time"

	"github.com/labstack/echo/v4"

	"mljr-web/internal/snakeai"
	"mljr-web/projects/homepage/pages"
)

// snakeRoom is the one shared game. It stays idle (no game, no model in memory)
// until a viewer connects and stops shortly after the last one leaves.
var snakeRoom = snakeai.NewRoom(snakeai.Config{
	N:          30,
	Tick:       80 * time.Millisecond,
	Grace:      20 * time.Second,
	MaxViewers: 50,
})

// registerSnakeLive mounts the SSE feed. It lives under /api/, which the shared
// gzip middleware skips, so frames are never buffered.
func registerSnakeLive(e *echo.Echo) {
	e.GET(pages.AISnakeStreamPath, echo.WrapHandler(snakeai.StreamHandler(snakeRoom)))
}
