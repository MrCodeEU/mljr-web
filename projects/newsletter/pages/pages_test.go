package pages

import (
	"bytes"
	"io"
	"maps"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"mljr-web/projects/newsletter/internal/testutil"

	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
)

// httpApp wraps a tests.TestApp with a single, already-built router/mux, so a
// test can issue several requests against the same app instance without
// hitting PocketBase's known issue of accumulating duplicate OnServe-bound
// routes (e.g. /_/extensions.js) when apis.NewRouter is called more than once
// per app (see https://github.com/pocketbase/pocketbase/discussions/7267,
// which is also why tests.ApiScenario is unsuitable for multi-request flows
// against a shared app).
type httpApp struct {
	app *tests.TestApp
	mux http.Handler
}

// newTestApp boots a fresh PocketBase test app with this package's routes
// and record hooks bound, and builds its HTTP mux once for reuse across
// however many requests a single test needs to make.
func newTestApp(t *testing.T) *httpApp {
	app := testutil.NewApp(t)

	baseRouter, err := apis.NewRouter(app)
	if err != nil {
		t.Fatalf("apis.NewRouter: %v", err)
	}

	serveEvent := &core.ServeEvent{App: app, Router: baseRouter}
	err = app.OnServe().Trigger(serveEvent, func(e *core.ServeEvent) error {
		if err := RegisterRoutes(e); err != nil {
			return err
		}
		return e.Next()
	})
	if err != nil {
		t.Fatalf("trigger OnServe: %v", err)
	}
	RegisterHooks(app)

	mux, err := baseRouter.BuildMux()
	if err != nil {
		t.Fatalf("build mux: %v", err)
	}

	return &httpApp{app: app, mux: mux}
}

type response struct {
	*http.Response
	Body string
}

// do issues a single HTTP request against the app's pre-built mux.
func (h *httpApp) do(t *testing.T, method, url string, body io.Reader, headers map[string]string) response {
	t.Helper()
	req := httptest.NewRequest(method, url, body)
	req.Header.Set("content-type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, strings.TrimSpace(v))
	}
	rec := httptest.NewRecorder()
	h.mux.ServeHTTP(rec, req)
	res := rec.Result()
	return response{Response: res, Body: rec.Body.String()}
}

func formBody(values map[string]string) (*strings.Reader, string) {
	q := url.Values{}
	for k, v := range values {
		q.Set(k, v)
	}
	return strings.NewReader(q.Encode()), "application/x-www-form-urlencoded"
}

func multipartBody(t *testing.T, fields map[string]string) (*bytes.Buffer, string) {
	t.Helper()
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	for k, v := range fields {
		if err := w.WriteField(k, v); err != nil {
			t.Fatalf("write multipart field %s: %v", k, err)
		}
	}
	if err := w.Close(); err != nil {
		t.Fatalf("close multipart writer: %v", err)
	}
	return &buf, w.FormDataContentType()
}

func cookieHeader(t *testing.T, user *core.Record) map[string]string {
	t.Helper()
	return map[string]string{"Cookie": testutil.AuthCookie(t, user)}
}

func mergeHeaders(hdrs ...map[string]string) map[string]string {
	out := map[string]string{}
	for _, m := range hdrs {
		maps.Copy(out, m)
	}
	return out
}

func newRecordHelper(t *testing.T, app *tests.TestApp, col *core.Collection, fields map[string]any) *core.Record {
	t.Helper()
	rec := core.NewRecord(col)
	for k, v := range fields {
		rec.Set(k, v)
	}
	if err := app.Save(rec); err != nil {
		t.Fatalf("save %s record: %v", col.Name, err)
	}
	return rec
}

func expectStatus(t *testing.T, res response, want int) {
	t.Helper()
	if res.StatusCode != want {
		t.Errorf("expected status %d, got %d\nbody: %s", want, res.StatusCode, res.Body)
	}
}

func expectContains(t *testing.T, res response, substrs ...string) {
	t.Helper()
	for _, s := range substrs {
		if !strings.Contains(res.Body, s) {
			t.Errorf("expected body to contain %q\nbody: %s", s, res.Body)
		}
	}
}
