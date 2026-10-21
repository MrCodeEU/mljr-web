package snakesim

import "net/http"

// StreamHandler is the SSE endpoint: GET -> Datastar patch-signals stream of the shared game.
// Mount it under a path the gzip middleware skips (SSE must not be buffered).
func StreamHandler(r *Room) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		fl, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "streaming unsupported", http.StatusInternalServerError)
			return
		}
		sub := r.Subscribe()
		if sub == nil {
			http.Error(w, "too many viewers, try again shortly", http.StatusServiceUnavailable)
			return
		}
		defer sub.Close()
		h := w.Header()
		h.Set("Content-Type", "text/event-stream")
		h.Set("Cache-Control", "no-cache")
		h.Set("Connection", "keep-alive")
		h.Set("X-Accel-Buffering", "no") // reverse proxies: do not buffer
		w.WriteHeader(http.StatusOK)
		fl.Flush()
		for {
			select {
			case <-req.Context().Done():
				return
			case msg := <-sub.C:
				if _, err := w.Write(msg); err != nil {
					return
				}
				fl.Flush()
			}
		}
	}
}
