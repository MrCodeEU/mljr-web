package main

import (
	"bufio"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"mljr-web/internal/web"
	"mljr-web/projects/homepage/pages"
)

// The stream must reach the client uncompressed and frame by frame, even when the
// browser advertises gzip: the shared gzip middleware has to skip /api/ paths.
func TestSnakeLiveStreamIsUncompressedSSE(t *testing.T) {
	e := web.NewEcho()
	registerSnakeLive(e)
	srv := httptest.NewServer(e)
	defer srv.Close()

	req, _ := http.NewRequest(http.MethodGet, srv.URL+pages.AISnakeStreamPath, nil)
	req.Header.Set("Accept-Encoding", "gzip")
	client := &http.Client{Timeout: 5 * time.Second, Transport: &http.Transport{DisableCompression: true}}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = resp.Body.Close() }()
	if ct := resp.Header.Get("Content-Type"); ct != "text/event-stream" {
		t.Fatalf("content-type %q", ct)
	}
	if ce := resp.Header.Get("Content-Encoding"); ce != "" {
		t.Fatalf("stream was compressed (%q); SSE must not be buffered", ce)
	}
	sc := bufio.NewScanner(resp.Body)
	var event, data string
	for sc.Scan() && (event == "" || data == "") {
		line := sc.Text()
		if strings.HasPrefix(line, "event: ") {
			event = line
		}
		if strings.HasPrefix(line, "data: signals ") {
			data = line
		}
	}
	if event != "event: datastar-patch-signals" || !strings.Contains(data, `"snake":{"n":30`) {
		t.Fatalf("unexpected first frame: %q / %q", event, data)
	}
}
