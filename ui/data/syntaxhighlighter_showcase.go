//go:build showcase

package data

import (
	"mljr-web/ui/registry"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "syntax-highlighter", Name: "Syntax Highlighter", Category: "data",
		Summary: "Server-side syntax highlighting via chroma. Zero client JS — inline styles. Supports 300+ languages, 40+ themes.",
		Code: `data.SyntaxHighlighter(data.SyntaxHighlighterProps{
    Language:        "go",
    Theme:           "dracula",
    ShowLineNumbers: true,
    Filename:        "main.go",
}, code)`,
		Render: func(p map[string]string) g.Node {
			goCode := `package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
	})
	http.ListenAndServe(":8080", nil)
}`
			tsCode := `interface User {
  id: number;
  name: string;
  role: "admin" | "user";
}

async function fetchUser(id: number): Promise<User> {
  const res = await fetch("/api/users/" + id);
  return res.json();
}`
			return h.Div(
				h.Style("display:flex;flex-direction:column;gap:var(--sp-4)"),
				SyntaxHighlighter(SyntaxHighlighterProps{
					Language:        "go",
					Theme:           "dracula",
					ShowLineNumbers: true,
					Filename:        "main.go",
				}, goCode),
				SyntaxHighlighter(SyntaxHighlighterProps{
					Language: "typescript",
					Theme:    "github",
					Filename: "user.ts",
				}, tsCode),
			)
		},
	})
}
