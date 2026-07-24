package main

import (
	"fmt"
	"net/http"
	"time"
)

// handleSSE streams Server-Sent Events: data events every interval and
// comment-line heartbeats every heartbeat period.
func handleSSE(w http.ResponseWriter, r *http.Request) {
	cfg, err := parseStreamConfig(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id := connSeq.Add(1)
	start := time.Now()
	logOpen(id, "/sse", r, cfg)

	// delay holds back the response headers to probe response-start timeouts.
	if !sleepCtx(r.Context(), cfg.Delay) {
		logClose(id, "/sse", start, 0, "closed during delay (before headers)")
		return
	}

	rc := http.NewResponseController(w)
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("X-Accel-Buffering", "no") // nginx-family proxies: do not buffer
	w.WriteHeader(http.StatusOK)
	if err := rc.Flush(); err != nil {
		logClose(id, "/sse", start, 0, "flush failed: "+err.Error())
		return
	}

	sent, err := runTicker(r.Context(), cfg, sink{
		Data: func(seq uint64, elapsed time.Duration) error {
			if _, err := fmt.Fprintf(w, "id: %d\ndata: %s\n\n", seq, dataLine(id, r.Proto, seq, elapsed)); err != nil {
				return err
			}
			return rc.Flush()
		},
		Heartbeat: func(elapsed time.Duration) error {
			// SSE comment line: bytes on the wire without being an event.
			if _, err := fmt.Fprintf(w, ": %s\n\n", heartbeatLine(id, elapsed)); err != nil {
				return err
			}
			return rc.Flush()
		},
	})
	if err == nil {
		fmt.Fprintf(w, "event: done\ndata: %s\n\n", doneLine(id, sent, time.Since(start)))
		_ = rc.Flush()
	}
	logClose(id, "/sse", start, sent, closeReason(err))
}
