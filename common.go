package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
	"time"
)

// connSeq numbers every test connection so client output and server logs can
// be correlated.
var connSeq atomic.Int64

// streamConfig holds the query parameters shared by /sse, /stream, /ws and
// /spdy. All values use Go duration syntax (500ms, 30s, 5m, 2h, 0).
type streamConfig struct {
	Duration  time.Duration // total lifetime; 0 = unlimited
	Interval  time.Duration // data tick period; 0 = no data ticks
	Heartbeat time.Duration // protocol keepalive period; 0 = off
	Delay     time.Duration // silence before the first response byte
	IdleAfter time.Duration // go fully silent after this elapsed time; 0 = never
}

func (c streamConfig) String() string {
	return fmt.Sprintf("duration=%s interval=%s heartbeat=%s delay=%s idleAfter=%s",
		c.Duration, c.Interval, c.Heartbeat, c.Delay, c.IdleAfter)
}

const minTick = 10 * time.Millisecond

// parseStreamConfig rejects unknown or invalid parameters: a silently ignored
// typo (e.g. hearbeat=10s) would corrupt the measurement this app exists for.
func parseStreamConfig(r *http.Request) (streamConfig, error) {
	cfg := streamConfig{Duration: 60 * time.Second, Interval: time.Second}
	for key, vals := range r.URL.Query() {
		var dst *time.Duration
		switch key {
		case "duration":
			dst = &cfg.Duration
		case "interval":
			dst = &cfg.Interval
		case "heartbeat":
			dst = &cfg.Heartbeat
		case "delay":
			dst = &cfg.Delay
		case "idleAfter":
			dst = &cfg.IdleAfter
		default:
			return cfg, fmt.Errorf("unknown query parameter %q (known: duration, interval, heartbeat, delay, idleAfter)", key)
		}
		v := vals[len(vals)-1]
		d, err := time.ParseDuration(v)
		if err != nil {
			return cfg, fmt.Errorf("parameter %q: %q is not a Go duration (e.g. 500ms, 30s, 5m, 0)", key, v)
		}
		if d < 0 {
			return cfg, fmt.Errorf("parameter %q must be >= 0", key)
		}
		*dst = d
	}
	if cfg.Interval > 0 && cfg.Interval < minTick {
		return cfg, fmt.Errorf("interval must be >= %s (or 0 to disable)", minTick)
	}
	if cfg.Heartbeat > 0 && cfg.Heartbeat < minTick {
		return cfg, fmt.Errorf("heartbeat must be >= %s (or 0 to disable)", minTick)
	}
	return cfg, nil
}

// sink adapts one protocol connection/stream to runTicker. Every callback runs
// on the runTicker goroutine only — that is the single-writer guarantee the
// WebSocket handler relies on.
type sink struct {
	Data      func(seq uint64, elapsed time.Duration) error
	Heartbeat func(elapsed time.Duration) error
	Incoming  <-chan string          // optional peer payloads to echo (ws/spdy)
	Echo      func(msg string) error // required iff Incoming is set
}

// runTicker drives one connection/stream through the configured scenario.
// Returned error: nil = Duration reached (caller sends its "done" message);
// ctx.Err() = peer/proxy closed first; anything else = a write failed.
func runTicker(ctx context.Context, cfg streamConfig, s sink) (uint64, error) {
	start := time.Now()
	var seq uint64

	var dataC, hbC <-chan time.Time
	if cfg.Interval > 0 && s.Data != nil {
		t := time.NewTicker(cfg.Interval)
		defer t.Stop()
		dataC = t.C
	}
	if cfg.Heartbeat > 0 && s.Heartbeat != nil {
		t := time.NewTicker(cfg.Heartbeat)
		defer t.Stop()
		hbC = t.C
	}
	var durC, idleC <-chan time.Time
	if cfg.Duration > 0 {
		t := time.NewTimer(cfg.Duration)
		defer t.Stop()
		durC = t.C
	}
	if cfg.IdleAfter > 0 {
		t := time.NewTimer(cfg.IdleAfter)
		defer t.Stop()
		idleC = t.C
	}

	for {
		select {
		case <-ctx.Done():
			return seq, ctx.Err()
		case <-durC:
			return seq, nil
		case <-idleC:
			dataC, hbC = nil, nil // go silent, but keep the connection open
		case <-dataC:
			seq++
			if err := s.Data(seq, time.Since(start)); err != nil {
				return seq, err
			}
		case msg, ok := <-s.Incoming:
			if !ok {
				s.Incoming = nil
				continue
			}
			if err := s.Echo(msg); err != nil {
				return seq, err
			}
		case <-hbC:
			if err := s.Heartbeat(time.Since(start)); err != nil {
				return seq, err
			}
		}
	}
}

// sleepCtx waits d (no-op if 0) and reports false if ctx fired first.
func sleepCtx(ctx context.Context, d time.Duration) bool {
	if d <= 0 {
		return true
	}
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return false
	case <-t.C:
		return true
	}
}

func nowStamp() string {
	return time.Now().UTC().Format(time.RFC3339Nano)
}

func dataLine(connID int64, proto string, seq uint64, elapsed time.Duration) string {
	return fmt.Sprintf("seq=%d elapsed=%s time=%s conn=%d proto=%s",
		seq, elapsed.Round(time.Millisecond), nowStamp(), connID, proto)
}

func heartbeatLine(connID int64, elapsed time.Duration) string {
	return fmt.Sprintf("hb conn=%d elapsed=%s time=%s", connID, elapsed.Round(time.Millisecond), nowStamp())
}

func doneLine(connID int64, sent uint64, elapsed time.Duration) string {
	return fmt.Sprintf("done conn=%d sent=%d elapsed=%s time=%s",
		connID, sent, elapsed.Round(time.Millisecond), nowStamp())
}

func logOpen(id int64, endpoint string, r *http.Request, cfg streamConfig) {
	log.Printf("conn=%d open endpoint=%s remote=%s proto=%s cfg={%s}", id, endpoint, r.RemoteAddr, r.Proto, cfg)
}

func logClose(id int64, endpoint string, start time.Time, sent uint64, reason string) {
	log.Printf("conn=%d close endpoint=%s elapsed=%s sent=%d reason=%q",
		id, endpoint, time.Since(start).Round(time.Millisecond), sent, reason)
}

func closeReason(err error) string {
	switch {
	case err == nil:
		return "duration reached"
	case errors.Is(err, context.Canceled), errors.Is(err, context.DeadlineExceeded):
		return "closed by peer/proxy"
	default:
		return "write failed: " + err.Error()
	}
}
