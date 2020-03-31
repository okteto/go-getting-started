FROM golang:alpine as builder

WORKDIR /app
ADD . .
RUN go build -o app

##########################

FROM builder as dev

RUN apk --update --no-cache add bash

COPY bashrc /root/.bashrc

RUN go get github.com/codegangsta/gin && \
    go get github.com/go-delve/delve/cmd/dlv && \
    go get golang.org/x/tools/gopls

##########################

FROM alpine as prod

WORKDIR /app
COPY --from=builder /app/app /app/app
EXPOSE 8080
CMD ["./app"]