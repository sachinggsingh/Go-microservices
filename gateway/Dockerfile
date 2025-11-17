FROM golang:1.25.3-alpine3.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download


COPY gateway ./gateway
COPY ../pb ./pb

WORKDIR /app/gateway

RUN go build -o main ./cmd/server

FROM alpine:latest

WORKDIR /root/

COPY  --from=builder /app/gateway/main .

EXPOSE 8085

CMD ["./main"]