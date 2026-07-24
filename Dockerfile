FROM golang:1.24-bookworm

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /usr/local/bin/hello-world .

EXPOSE 8080
CMD ["/usr/local/bin/hello-world"]
