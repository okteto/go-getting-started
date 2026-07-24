package main

import (
	"context"
	"errors"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Test tool: accept any origin.
	CheckOrigin: func(*http.Request) bool { return true },
}

// handleWS upgrades to WebSocket and sends text-message ticks every interval,
// ping frames every heartbeat period, and echoes every client message.
func handleWS(w http.ResponseWriter, r *http.Request) {
	cfg, err := parseStreamConfig(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id := connSeq.Add(1)
	start := time.Now()
	logOpen(id, "/ws", r, cfg)

	// delay holds back the 101 response to probe upgrade/response-start timeouts.
	if !sleepCtx(r.Context(), cfg.Delay) {
		logClose(id, "/ws", start, 0, "closed during delay (before upgrade)")
		return
	}
	c, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		// Upgrade already wrote the HTTP error response.
		logClose(id, "/ws", start, 0, "upgrade failed: "+err.Error())
		return
	}
	defer c.Close()
	// No read/write deadlines anywhere: the server must never close a quiet peer.

	// The connection is hijacked, so r.Context() no longer signals peer
	// disconnects. The read pump is the death detector — and it must keep
	// running anyway, since control frames (ping/pong/close) are only
	// processed while a read is pending; the default ping handler auto-pongs.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	incoming := make(chan string, 16)
	var readErr atomic.Value
	go func() {
		defer cancel()
		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				readErr.Store(err)
				return
			}
			select {
			case incoming <- string(msg):
			case <-ctx.Done():
				return
			}
		}
	}()

	sent, err := runTicker(ctx, cfg, sink{
		Data: func(seq uint64, elapsed time.Duration) error {
			return c.WriteMessage(websocket.TextMessage, []byte(dataLine(id, r.Proto, seq, elapsed)))
		},
		Heartbeat: func(elapsed time.Duration) error {
			// A real ping frame; the zero deadline means the write never times out.
			return c.WriteControl(websocket.PingMessage, []byte(heartbeatLine(id, elapsed)), time.Time{})
		},
		Incoming: incoming,
		Echo: func(msg string) error {
			return c.WriteMessage(websocket.TextMessage, []byte("echo "+msg))
		},
	})

	reason := closeReason(err)
	if errors.Is(err, context.Canceled) {
		if re, ok := readErr.Load().(error); ok && re != nil {
			reason = "peer/proxy closed: " + re.Error()
		}
	}
	if err == nil {
		_ = c.WriteMessage(websocket.TextMessage, []byte(doneLine(id, sent, time.Since(start))))
		_ = c.WriteControl(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, "done"),
			time.Now().Add(5*time.Second))
		// Brief grace so the peer's close reply can arrive before teardown.
		select {
		case <-ctx.Done():
		case <-time.After(time.Second):
		}
	}
	logClose(id, "/ws", start, sent, reason)
}
