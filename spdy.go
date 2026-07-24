package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"k8s.io/apimachinery/pkg/util/httpstream"
	"k8s.io/apimachinery/pkg/util/httpstream/spdy"
)

// handleSPDY performs a Kubernetes-style SPDY/3.1 upgrade (the protocol
// kubectl exec/port-forward use). Each stream the client opens gets ticks
// every interval and echoes every line the client writes; heartbeat > 0 sends
// connection-level SPDY PING frames at that period.
func handleSPDY(w http.ResponseWriter, r *http.Request) {
	cfg, err := parseStreamConfig(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id := connSeq.Add(1)
	start := time.Now()
	logOpen(id, "/spdy", r, cfg)

	// delay holds back the 101 response to probe upgrade/response-start timeouts.
	if !sleepCtx(r.Context(), cfg.Delay) {
		logClose(id, "/spdy", start, 0, "closed during delay (before upgrade)")
		return
	}

	streams := make(chan httpstream.Stream, 8)
	upgrader := spdy.NewResponseUpgraderWithPings(cfg.Heartbeat)
	conn := upgrader.UpgradeResponse(w, r, func(s httpstream.Stream, replySent <-chan struct{}) error {
		// Must return quickly: the SYN_REPLY that unblocks the client's
		// CreateStream is only sent after this handler returns.
		select {
		case streams <- s:
			return nil
		default:
			return fmt.Errorf("too many concurrent streams")
		}
	})
	if conn == nil {
		// UpgradeResponse already wrote the HTTP error response.
		logClose(id, "/spdy", start, 0, "upgrade failed (SPDY needs an HTTP/1.1 Upgrade request)")
		return
	}
	defer conn.Close()
	// Never conn.SetIdleTimeout(...): the server must not reap idle connections.

	// The connection is hijacked, so r.Context() no longer signals peer
	// disconnects; CloseChan does.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		<-conn.CloseChan()
		cancel()
	}()

	var total atomic.Uint64
	var wg sync.WaitGroup
	go func() {
		for {
			select {
			case s := <-streams:
				wg.Add(1)
				go func() {
					defer wg.Done()
					serveSPDYStream(ctx, cfg, id, s, &total)
				}()
			case <-ctx.Done():
				return
			}
		}
	}()

	var durC <-chan time.Time
	if cfg.Duration > 0 {
		t := time.NewTimer(cfg.Duration)
		defer t.Stop()
		durC = t.C
	}
	select {
	case <-ctx.Done():
		logClose(id, "/spdy", start, total.Load(), "connection closed by peer/proxy")
	case <-durC:
		// Per-stream tickers share cfg.Duration and started marginally later,
		// so they finish (and write their done lines) first.
		wg.Wait()
		time.Sleep(200 * time.Millisecond) // let in-flight done lines reach the client before Close
		logClose(id, "/spdy", start, total.Load(), "duration reached")
	}
}

func serveSPDYStream(ctx context.Context, cfg streamConfig, connID int64, st httpstream.Stream, total *atomic.Uint64) {
	start := time.Now()
	log.Printf("conn=%d stream=%d open endpoint=/spdy", connID, st.Identifier())

	incoming := make(chan string, 16)
	go func() {
		// Read pump. EOF here is only a half-close by the client — connection
		// death is signalled via CloseChan/ctx, not via this pump.
		sc := bufio.NewScanner(st)
		for sc.Scan() {
			select {
			case incoming <- sc.Text():
			case <-ctx.Done():
				return
			}
		}
	}()

	sent, err := runTicker(ctx, cfg, sink{
		Data: func(seq uint64, elapsed time.Duration) error {
			_, werr := fmt.Fprintln(st, dataLine(connID, "SPDY/3.1", seq, elapsed))
			return werr
		},
		// Heartbeats are connection-level SPDY pings, handled by the upgrader.
		Incoming: incoming,
		Echo: func(msg string) error {
			_, werr := fmt.Fprintf(st, "echo %s\n", msg)
			return werr
		},
	})
	total.Add(sent)
	if err == nil {
		fmt.Fprintln(st, doneLine(connID, sent, time.Since(start)))
	}
	_ = st.Close()
	log.Printf("conn=%d stream=%d close elapsed=%s sent=%d reason=%q",
		connID, st.Identifier(), time.Since(start).Round(time.Millisecond), sent, closeReason(err))
}
