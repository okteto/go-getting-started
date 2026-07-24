# Testing long-lived connections through a proxy

These endpoints measure how the ingress/reverse-proxy layer in front of the
app treats long-lived connections: idle timeouts, max request duration, and
whether heartbeats keep a connection alive — over SSE, plain HTTP streaming
(HTTP/1.1 and HTTP/2), WebSocket, and SPDY/3.1.

Two ground rules make the measurements trustworthy:

- **The server never closes connections on its own**: no HTTP server timeouts,
  no WebSocket read deadlines, no SPDY idle timeout. If a connection dies
  before `duration` elapses, something between the client and the app killed it.
- **Everything is timestamped**: every payload line carries `seq`, `elapsed`,
  wall-clock `time` and a `conn` id, and the server logs every connection's
  open (with its config) and close (with elapsed time and reason).

## Endpoints

| Path | Protocol | Drive it with |
|---|---|---|
| `/sse` | Server-Sent Events (heartbeats are `:` comment lines) | `curl -N` |
| `/stream` | Plain chunked stream over HTTP/1.1 or HTTP/2 (h2c); first line reports the negotiated protocol | `curl -N` |
| `/ws` | WebSocket (heartbeats are ping frames; echoes client messages) | `go run ./cmd/client` |
| `/spdy` | SPDY/3.1 kubectl-style upgrade (heartbeats are SPDY PING frames; echoes lines per stream) | `go run ./cmd/client` |
| `/info` | Reports the protocol/TLS/headers the request arrived with | `curl` |

## Query parameters (shared by /sse, /stream, /ws, /spdy)

Values use Go duration syntax (`500ms`, `30s`, `5m`, `2h`, `0`). Unknown or
invalid parameters fail with `400` so a typo can't silently skew a measurement.

| Param | Default | Meaning |
|---|---|---|
| `duration` | `60s` | Total lifetime; ends with a `done` message and a clean close. `0` = unlimited (run until something kills it). |
| `interval` | `1s` | Period between data messages. `0` = no data → idle connection. |
| `heartbeat` | `0` (off) | Period between protocol-level keepalives (SSE comment, WS ping, SPDY ping, `#` line on /stream). |
| `delay` | `0` | Silence before the first response byte — headers / the 101 upgrade are held back. Probes response-start and upgrade timeouts. |
| `idleAfter` | `0` (never) | Stop sending everything after this much elapsed time → active-then-idle. |

## Scenarios

Set `BASE` to the app URL, e.g. `BASE=https://hello-world-<namespace>.<okteto-domain>`
(or `BASE=http://localhost:8080` when running locally with `go run .`).

**Idle from the beginning** — how long until a silent connection is cut?

    curl -sN "$BASE/sse?interval=0&duration=0"
    go run ./cmd/client -url "$BASE/ws?interval=0&duration=0"
    go run ./cmd/client -mode spdy -url "$BASE/spdy?interval=0&duration=0"

**Active for a few seconds, then idle**:

    curl -sN "$BASE/sse?interval=1s&idleAfter=5s&duration=0"

**Heartbeat with a customizable interval** — does a 20s heartbeat survive the idle timeout?

    curl -sN "$BASE/sse?interval=0&heartbeat=20s&duration=0"
    go run ./cmd/client -url "$BASE/ws?interval=0&heartbeat=20s&duration=0"
    go run ./cmd/client -mode spdy -url "$BASE/spdy?interval=0&heartbeat=20s&duration=0"

**Max request duration** — steady traffic; when is an *active* connection cut?

    curl -sN "$BASE/stream?interval=5s&duration=0"

**Response-start timeout** — server silent before headers / the 101:

    time curl -sN "$BASE/stream?delay=120s&duration=1s"
    go run ./cmd/client -url "$BASE/ws?delay=120s&duration=1s"

**HTTP/2**:

    curl -sN --http2-prior-knowledge "http://localhost:8080/stream?duration=5s"  # direct h2c
    curl -sN --http2 "$BASE/stream?duration=5s"        # h2 to the ingress
    go run ./cmd/client -url "$BASE/stream?duration=5s" -http2

**SPDY/3.1 with client traffic** (client lines are echoed back):

    go run ./cmd/client -mode spdy -url "$BASE/spdy?duration=30s&interval=1s" -send-interval 2s

**WebSocket kept alive by client-side traffic only**:

    go run ./cmd/client -url "$BASE/ws?interval=0&duration=0" -send-interval 5s

## Reading the results

- **Client side**: `go run ./cmd/client` prints each message with its arrival
  time and `elapsed`; the final `closed ...` line carries the measured
  lifetime. Exit code 0 = clean (`duration` reached), 1 = killed early.
  `curl` streams stop printing when the proxy cuts the connection — wrap in
  `time` to capture when.
- **Server side**: correlate with the app log:
  `conn=N open endpoint=/sse ... cfg={...}` and
  `conn=N close endpoint=/sse elapsed=... reason=...`.
  A close reason of `"closed by peer/proxy"` before `duration` elapsed means
  the proxy (or the client) cut the connection.
- **Binary-search the thresholds**: if `heartbeat=30s` connections die but
  `heartbeat=10s` ones survive, the idle timeout is between 10s and 30s.
- The ingress usually talks **HTTP/1.1 to the app** even when the browser
  speaks HTTP/2 to the ingress; `/info` and the `/stream` preamble show the
  protocol of the last hop. WebSocket and SPDY upgrades always ride HTTP/1.1.

## Client flags

| Flag | Meaning |
|---|---|
| `-url` | Endpoint URL (`http://`, `https://`, `ws://`, `wss://`) |
| `-mode` | `auto` (default), `http`, `sse`, `ws`, `spdy` |
| `-http2` | Force HTTP/2: prior-knowledge h2c on `http://`, ALPN on `https://` |
| `-send-interval` | ws/spdy: send a client message at this period (keeps client→server traffic flowing) |
| `-insecure` | Skip TLS certificate verification |
