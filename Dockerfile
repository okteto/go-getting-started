FROM golang:buster

WORKDIR /app
COPY . .
RUN go build -o /usr/local/bin/hello-world

EXPOSE 8080
CMD ["/usr/local/bin/hello-world"]