package main

import (
	"fmt"
	"log"
	"net/http"
	"sort"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.LUTC)

	mux := http.NewServeMux()
	mux.HandleFunc("/", helloServer)
	mux.HandleFunc("/info", handleInfo)
	mux.HandleFunc("/sse", handleSSE)
	mux.HandleFunc("/stream", handleStream)
	mux.HandleFunc("/ws", handleWS)
	mux.HandleFunc("/spdy", handleSPDY)

	// h2c only intercepts HTTP/2 prior-knowledge and "Upgrade: h2c" requests;
	// WebSocket and SPDY upgrades stay on HTTP/1.1 where hijacking works.
	srv := &http.Server{
		Addr:    ":8080",
		Handler: h2c.NewHandler(mux, &http2.Server{}),
		// This app exists to measure proxy timeouts, so it must never be the
		// one closing connections: every server timeout stays disabled.
		ReadTimeout:       0,
		ReadHeaderTimeout: 0,
		WriteTimeout:      0,
		IdleTimeout:       0,
	}

	log.Println("Starting hello-world server on :8080 (endpoints: /, /info, /sse, /stream, /ws, /spdy)")
	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}

func helloServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello world!")
}

// handleInfo reports how the request arrived (protocol, TLS, headers) so the
// protocol of the proxy->app hop can be confirmed.
func handleInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "time: %s\n", nowStamp())
	fmt.Fprintf(w, "proto: %s\n", r.Proto)
	fmt.Fprintf(w, "method: %s\n", r.Method)
	fmt.Fprintf(w, "host: %s\n", r.Host)
	fmt.Fprintf(w, "url: %s\n", r.URL.String())
	fmt.Fprintf(w, "remote: %s\n", r.RemoteAddr)
	fmt.Fprintf(w, "tls: %v\n", r.TLS != nil)
	keys := make([]string, 0, len(r.Header))
	for k := range r.Header {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	fmt.Fprintln(w, "headers:")
	for _, k := range keys {
		for _, v := range r.Header[k] {
			fmt.Fprintf(w, "  %s: %s\n", k, v)
		}
	}
}
