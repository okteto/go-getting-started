FROM golang:buster

WORKDIR /app
COPY go.mod go.mod
RUN go mod download && go mod verify

ADD . .
RUN go mod tidy && go build -o /usr/local/bin/hello-world

EXPOSE 8080
CMD ["/usr/local/bin/hello-world"]