FROM golang:1.14.1-alpine3.11 as builder
RUN apk --update --no-cache add bash
WORKDIR /app
ADD . .
RUN go build -o app

FROM alpine:3.11 as prod
WORKDIR /app
COPY --from=builder /app/app /app/app
EXPOSE 8080
CMD ["/app/app"]