package web

import (
	"io/fs"
	"net/http"
	"os"
	"path/filepath"

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

// MountLogos serves a project's logo variants (assets/static/img/logo) under
// /logo/, so the brand assets can be linked from other sites (e.g. mljr.eu/logo/Logo-h.png).
func MountLogos(e *echo.Echo, embedded fs.FS, devPath string, dev bool) {
	if dev {
		e.Static("/logo", devPath+"/img/logo")
		return
	}
	sub, err := fs.Sub(embedded, "assets/static/img/logo")
	if err != nil {
		panic(err)
	}
	e.GET("/logo/*", echo.WrapHandler(http.StripPrefix("/logo/", http.FileServer(http.FS(sub)))))
}

// MountDataAssets serves the mljr-data tree's assets/ directory (project
// screenshots, thesis PDFs, etc.) under /assets/. dataFile is the resolved
// HOMEPAGE_DATA_FILE path (<repo>/generated/site-data.json); assets/ is its
// sibling. This is plain filesystem serving in both dev and prod — content
// is synced onto disk at runtime, not embedded into the binary.
func MountDataAssets(e *echo.Echo, dataFile string) {
	assetsDir := filepath.Join(filepath.Dir(filepath.Dir(dataFile)), "assets")
	e.Static("/assets", assetsDir)
}

// IsDev reports whether the process is in dev mode.
func IsDev() bool { return os.Getenv("MLJR_ENV") != "prod" }
