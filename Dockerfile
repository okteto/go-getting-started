FROM golang:buster

WORKDIR /app
ADD . .
RUN go build -o /usr/local/bin/hello-world

EXPOSE 8443
CMD ["/usr/local/bin/hello-world"]