package main

import (
	"fmt"
	"net/http"
	"time"
)

// handleStream is a plain chunked text stream. It works over HTTP/1.1 and
// HTTP/2 (h2c, or h2 via a TLS-terminating proxy) and announces the
// negotiated protocol in its first line, which makes it the HTTP/2 test
// endpoint.
func handleStream(w http.ResponseWriter, r *http.Request) {
	cfg, err := parseStreamConfig(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id := connSeq.Add(1)
	start := time.Now()
	logOpen(id, "/stream", r, cfg)

	if !sleepCtx(r.Context(), cfg.Delay) {
		logClose(id, "/stream", start, 0, "closed during delay (before headers)")
		return
	}

	rc := http.NewResponseController(w)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("X-Accel-Buffering", "no")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "connected conn=%d proto=%s remote=%s time=%s\n", id, r.Proto, r.RemoteAddr, nowStamp())
	if err := rc.Flush(); err != nil {
		logClose(id, "/stream", start, 0, "flush failed: "+err.Error())
		return
	}

	sent, err := runTicker(r.Context(), cfg, sink{
		Data: func(seq uint64, elapsed time.Duration) error {
			if _, err := fmt.Fprintln(w, dataLine(id, r.Proto, seq, elapsed)); err != nil {
				return err
			}
			return rc.Flush()
		},
		Heartbeat: func(elapsed time.Duration) error {
			if _, err := fmt.Fprintf(w, "# %s\n", heartbeatLine(id, elapsed)); err != nil {
				return err
			}
			return rc.Flush()
		},
	})
	if err == nil {
		fmt.Fprintln(w, doneLine(id, sent, time.Since(start)))
		_ = rc.Flush()
	}
	logClose(id, "/stream", start, sent, closeReason(err))
}
