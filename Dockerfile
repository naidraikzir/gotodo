# Builder
FROM golang:1.13-alpine3.10 AS builder
RUN apk update && \
    apk upgrade && \
    apk add --no-cache ca-certificates gcc musl-dev
WORKDIR /build
COPY . .
RUN go mod download && CGO_ENABLED=1 go build

# App
FROM alpine:latest
RUN apk update && apk upgrade && apk add sqlite
WORKDIR /app
COPY --from=builder /build/todo .
EXPOSE 4321
CMD ["./todo"]
