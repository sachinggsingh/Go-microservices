FROM golang:1.25.3-alpine3.22 AS deps

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

FROM golang:1.25.3-alpine3.22 AS builder

WORKDIR /app

COPY --from=deps /app/go.mod /app/go.mod

COPY --from=deps /app/go.sum /app/go.sum

COPY product ./product

COPY ../pb ./pb

WORKDIR /app/product

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o main ./cmd/server

FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /root/

COPY --from=builder /app/product/main .

EXPOSE 8081

CMD ["./main"]

