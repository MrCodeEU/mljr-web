package web

import (
	"io/fs"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

// MountStatic serves a project's embedded assets/static under /static/ in prod,
// and the on-disk path under /static/ in dev (so tailwind --watch and edits
// land without rebuild).
func MountStatic(e *echo.Echo, embedded fs.FS, devPath string, dev bool) {
	if dev {
		e.Static("/static", devPath)
		return
	}
	sub, err := fs.Sub(embedded, "assets/static")
	if err != nil {
		panic(err)
	}
	e.GET("/static/*", echo.WrapHandler(http.StripPrefix("/static/", http.FileServer(http.FS(sub)))))
}

// IsDev reports whether the process is in dev mode.
func IsDev() bool { return os.Getenv("MLJR_ENV") != "prod" }
