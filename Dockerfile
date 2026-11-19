FROM golang:1.26.5-trixie AS dev

RUN GOBIN=/usr/local/bin go install github.com/air-verse/air@v1.67.4 && \
    GOBIN=/usr/local/bin go install github.com/go-delve/delve/cmd/dlv@v1.27.1

WORKDIR /usr/src/app
COPY . .
RUN go build -o /usr/local/bin/hello-world

FROM debian:trixie-slim AS prod

COPY --from=dev /usr/local/bin/hello-world /usr/local/bin/hello-world
CMD ["/usr/local/bin/hello-world"]
