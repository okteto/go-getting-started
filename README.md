# Getting Started on Okteto with Go

This example shows how to use the [Okteto CLI](https://github.com/okteto/okteto) to develop a Go Sample App directly in Kubernetes. The Go Sample App is deployed using Kubernetes manifests.

This is the application used for the [Getting Started on Okteto with Go](https://www.okteto.com/docs/samples/golang/) tutorial.

## Testing long-lived connections

This branch adds streaming endpoints (`/sse`, `/stream`, `/ws`, `/spdy`, `/info`) to measure how proxies treat long-lived connections: idle timeouts, max request duration, and heartbeats over SSE, HTTP/2, WebSocket and SPDY/3.1. See [TESTING.md](TESTING.md).
