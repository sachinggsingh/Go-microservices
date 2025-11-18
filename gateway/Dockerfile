FROM golang:1.25.3-alpine3.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY gateway ./gateway

COPY ../pb ./pb

WORKDIR /app/gateway

# Build statically linked, stripped binary for Linux AMD64
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o main ./cmd/server

FROM alpine:latest

# Install CA certificates for HTTPS support
RUN apk add --no-cache ca-certificates

WORKDIR /root/

COPY  --from=builder /app/gateway/main .

EXPOSE 8085

CMD ["./main"]