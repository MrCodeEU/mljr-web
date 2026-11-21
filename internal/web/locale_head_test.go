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

func TestURLLanguageWinsWithoutHTMLAccept(t *testing.T) {
	e := echo.New()
	var got string
	h := LangPrefix()(LocaleMiddleware()(func(c echo.Context) error {
		got = Lang(c)
		return nil
	}))
	for _, accept := range []string{"text/html", "*/*", ""} {
		req := httptest.NewRequest("GET", "/de/", nil)
		if accept != "" {
			req.Header.Set("Accept", accept)
		}
		if err := h(e.NewContext(req, httptest.NewRecorder())); err != nil {
			t.Fatal(err)
		}
		if got != "de" {
			t.Errorf("Accept %q: lang = %q, want de", accept, got)
		}
	}
}
