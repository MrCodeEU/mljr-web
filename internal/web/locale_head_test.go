package web

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestLocalizedHeadEscapesPath(t *testing.T) {
	e := echo.New()
	c := e.NewContext(httptest.NewRequest("GET", "/", nil), httptest.NewRecorder())
	c.Set(pagePathKey, `/x"><script>alert(1)</script>`)
	out := string(localizedHead(c))
	if strings.Contains(out, "<script>") || strings.Contains(out, `x">`) {
		t.Fatalf("path not escaped: %s", out)
	}
	if !strings.Contains(out, `href="`+SiteURL+`/x%22%3E%3Cscript%3Ealert%281%29%3C/script%3E"`) {
		t.Fatalf("unexpected canonical: %s", out)
	}
}
