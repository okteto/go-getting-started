// Command client drives the test endpoints over protocols curl cannot speak
// (WebSocket, SPDY/3.1) and prints the arrival time of every message, so
// proxy idle/duration limits can be measured from the client side.
//
// It deliberately sets no timeouts of its own: if the connection dies before
// the server's `duration` elapses, something between client and app killed it.
package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"golang.org/x/net/http2"
	"k8s.io/apimachinery/pkg/util/httpstream/spdy"
)

var (
	start    = time.Now()
	received int
)

func p(format string, args ...any) {
	fmt.Printf("%s elapsed=%s | %s\n",
		time.Now().UTC().Format(time.RFC3339Nano),
		time.Since(start).Round(time.Millisecond),
		fmt.Sprintf(format, args...))
}

func fatal(format string, args ...any) {
	p("FATAL "+format, args...)
	os.Exit(1)
}

// finish reports the total lifetime and exits: 0 = clean end (server reached
// its duration), 1 = the connection was killed early.
func finish(clean bool, reason string) {
	p("closed received=%d clean=%v reason=%s", received, clean, reason)
	if clean {
		os.Exit(0)
	}
	os.Exit(1)
}

func main() {
	urlFlag := flag.String("url", "", "endpoint URL (http://, https://, ws:// or wss://)")
	mode := flag.String("mode", "auto", "auto|http|sse|ws|spdy (auto: ws(s):// -> ws, path /spdy -> spdy, else http)")
	useHTTP2 := flag.Bool("http2", false, "force HTTP/2 (prior-knowledge h2c on http://, ALPN on https://)")
	sendInterval := flag.Duration("send-interval", 0, "ws/spdy: send a client message at this period (0 = never)")
	insecure := flag.Bool("insecure", false, "skip TLS certificate verification")
	flag.Parse()
	if *urlFlag == "" {
		flag.Usage()
		os.Exit(2)
	}
	u, err := url.Parse(*urlFlag)
	if err != nil {
		fatal("bad url: %v", err)
	}

	m := *mode
	if m == "auto" {
		switch {
		case u.Scheme == "ws" || u.Scheme == "wss":
			m = "ws"
		case strings.Contains(u.Path, "/spdy"):
			m = "spdy"
		default:
			m = "http"
		}
	}
	if m == "sse" {
		m = "http"
	}

	tlsCfg := &tls.Config{InsecureSkipVerify: *insecure}
	switch m {
	case "http":
		runHTTP(u, *useHTTP2, tlsCfg)
	case "ws":
		runWS(u, *sendInterval, tlsCfg)
	case "spdy":
		runSPDY(u, *sendInterval, tlsCfg)
	default:
		fatal("unknown mode %q", m)
	}
}

func runHTTP(u *url.URL, forceH2 bool, tlsCfg *tls.Config) {
	if u.Scheme != "http" && u.Scheme != "https" {
		fatal("http mode needs an http(s):// url")
	}
	var rt http.RoundTripper
	if forceH2 {
		if u.Scheme == "http" {
			// Prior-knowledge cleartext HTTP/2 (h2c).
			rt = &http2.Transport{
				AllowHTTP: true,
				DialTLSContext: func(ctx context.Context, network, addr string, _ *tls.Config) (net.Conn, error) {
					var d net.Dialer
					return d.DialContext(ctx, network, addr)
				},
			}
		} else {
			rt = &http2.Transport{TLSClientConfig: tlsCfg}
		}
	} else {
		t := http.DefaultTransport.(*http.Transport).Clone()
		t.TLSClientConfig = tlsCfg
		rt = t
	}
	client := &http.Client{Transport: rt} // Timeout 0: only the proxy may kill this request
	resp, err := client.Get(u.String())
	if err != nil {
		fatal("request failed: %v", err)
	}
	defer resp.Body.Close()
	p("connected status=%s proto=%s", resp.Status, resp.Proto)

	br := bufio.NewReader(resp.Body)
	for {
		line, err := br.ReadString('\n')
		if line != "" {
			received++
			p("recv %s", strings.TrimRight(line, "\n"))
		}
		if err != nil {
			clean := errors.Is(err, io.EOF) && resp.StatusCode == http.StatusOK
			finish(clean, err.Error())
		}
	}
}

func runWS(u *url.URL, sendInterval time.Duration, tlsCfg *tls.Config) {
	wsURL := *u
	switch wsURL.Scheme {
	case "http":
		wsURL.Scheme = "ws"
	case "https":
		wsURL.Scheme = "wss"
	}
	d := websocket.Dialer{
		HandshakeTimeout: 0, // none: the proxy's upgrade timeout is what's being measured
		TLSClientConfig:  tlsCfg,
		Proxy:            http.ProxyFromEnvironment,
	}
	c, resp, err := d.Dial(wsURL.String(), nil)
	if err != nil {
		if resp != nil {
			fatal("dial failed: %v (status=%s)", err, resp.Status)
		}
		fatal("dial failed: %v", err)
	}
	defer c.Close()
	p("connected status=%s", resp.Status)

	c.SetPingHandler(func(appData string) error {
		p("recv ping %q", appData)
		err := c.WriteControl(websocket.PongMessage, []byte(appData), time.Now().Add(10*time.Second))
		if err == websocket.ErrCloseSent {
			return nil
		} else if e, ok := err.(net.Error); ok && e.Timeout() {
			return nil
		}
		return err
	})
	c.SetPongHandler(func(appData string) error {
		p("recv pong %q", appData)
		return nil
	})

	if sendInterval > 0 {
		go func() {
			t := time.NewTicker(sendInterval)
			defer t.Stop()
			for i := 1; ; i++ {
				<-t.C
				msg := fmt.Sprintf("client msg=%d time=%s", i, time.Now().UTC().Format(time.RFC3339Nano))
				if err := c.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
					return
				}
			}
		}()
	}

	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			var ce *websocket.CloseError
			if errors.As(err, &ce) {
				finish(ce.Code == websocket.CloseNormalClosure, fmt.Sprintf("close code=%d text=%q", ce.Code, ce.Text))
			}
			finish(false, err.Error())
		}
		received++
		p("recv %s", string(msg))
	}
}

func runSPDY(u *url.URL, sendInterval time.Duration, tlsCfg *tls.Config) {
	if u.Scheme != "http" && u.Scheme != "https" {
		fatal("spdy mode needs an http(s):// url")
	}
	rt, err := spdy.NewRoundTripper(tlsCfg)
	if err != nil {
		fatal("round tripper: %v", err)
	}
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		fatal("bad request: %v", err)
	}
	resp, err := rt.RoundTrip(req) // adds Connection: Upgrade + Upgrade: SPDY/3.1 itself
	if err != nil {
		fatal("upgrade request failed: %v", err)
	}
	conn, err := rt.NewConnection(resp)
	if err != nil {
		fatal("upgrade rejected: %v (status=%s)", err, resp.Status)
	}
	defer conn.Close()
	p("connected status=%s (SPDY/3.1)", resp.Status)

	st, err := conn.CreateStream(http.Header{})
	if err != nil {
		fatal("create stream: %v", err)
	}
	p("stream opened")

	if sendInterval > 0 {
		go func() {
			t := time.NewTicker(sendInterval)
			defer t.Stop()
			for i := 1; ; i++ {
				<-t.C
				if _, err := fmt.Fprintf(st, "client msg=%d time=%s\n", i, time.Now().UTC().Format(time.RFC3339Nano)); err != nil {
					return
				}
			}
		}()
	}

	lines := make(chan string)
	scanErr := make(chan error, 1)
	go func() {
		sc := bufio.NewScanner(st)
		for sc.Scan() {
			lines <- sc.Text()
		}
		scanErr <- sc.Err() // nil on EOF
		close(lines)
	}()

	sawDone := false
	for {
		select {
		case line, ok := <-lines:
			if !ok {
				if err := <-scanErr; err != nil {
					finish(false, "stream read error: "+err.Error())
				}
				finish(sawDone, "stream closed by server (EOF)")
			}
			received++
			if strings.HasPrefix(line, "done ") {
				sawDone = true
			}
			p("recv %s", line)
		case <-conn.CloseChan():
			finish(sawDone, "connection closed")
		}
	}
}
