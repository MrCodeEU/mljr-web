package web

import (
	"bytes"

	"github.com/labstack/echo/v4"
	"golang.org/x/net/html"
	g "maragu.dev/gomponents"
)

// Render writes a gomponents node as a UTF-8 HTML response. The node is first
// rendered to a buffer so the CSP can be rebuilt with sha256 hashes for every
// inline <script> in the page (see SetCSPWithScriptHashes).
func Render(c echo.Context, status int, n g.Node) error {
	var buf bytes.Buffer
	if err := n.Render(&buf); err != nil {
		return err
	}
	SetCSPWithScriptHashes(c, inlineScripts(buf.Bytes()))
	c.Response().Header().Set(echo.HeaderContentType, "text/html; charset=utf-8")
	c.Response().WriteHeader(status)
	_, err := c.Response().Writer.Write(buf.Bytes())
	return err
}

// RenderToString materializes a node into an HTML string, e.g. for
// sse.PatchElements(html string).
func RenderToString(n g.Node) string {
	var buf bytes.Buffer
	_ = n.Render(&buf)
	return buf.String()
}

// inlineScripts returns the exact byte content of every <script> element
// without a src attribute in the document. <script> is a raw-text element, so
// the parser preserves the content byte-for-byte — the hashes match what the
// browser computes.
func inlineScripts(doc []byte) [][]byte {
	root, err := html.Parse(bytes.NewReader(doc))
	if err != nil {
		return nil
	}
	var out [][]byte
	var walk func(n *html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "script" && !hasSrcAttr(n) {
			var content bytes.Buffer
			for ch := n.FirstChild; ch != nil; ch = ch.NextSibling {
				if ch.Type == html.TextNode {
					content.WriteString(ch.Data)
				}
			}
			out = append(out, content.Bytes())
		}
		for ch := n.FirstChild; ch != nil; ch = ch.NextSibling {
			walk(ch)
		}
	}
	walk(root)
	return out
}

func hasSrcAttr(n *html.Node) bool {
	for _, a := range n.Attr {
		if a.Key == "src" {
			return true
		}
	}
	return false
}
