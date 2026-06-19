package web

import (
	"github.com/pocketbase/pocketbase/core"
	g "maragu.dev/gomponents"
)

// RenderPB writes a gomponents node as a UTF-8 HTML response for a
// PocketBase request event (sibling to Render, which targets Echo).
func RenderPB(e *core.RequestEvent, status int, n g.Node) error {
	e.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	e.Response.WriteHeader(status)
	return n.Render(e.Response)
}
