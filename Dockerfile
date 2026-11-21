FROM golang:1.26.5-trixie AS dev

# The cache mounts keep the Go module cache and the Go build cache in
# BuildKit, which reuses them between builds
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    GOBIN=/usr/local/bin go install github.com/air-verse/air@v1.67.4 && \
    GOBIN=/usr/local/bin go install github.com/go-delve/delve/cmd/dlv@v1.27.1

WORKDIR /usr/src/app
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -o /usr/local/bin/hello-world

FROM debian:trixie-slim AS prod

COPY --from=dev /usr/local/bin/hello-world /usr/local/bin/hello-world
CMD ["/usr/local/bin/hello-world"]
