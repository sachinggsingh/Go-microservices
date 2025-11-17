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

RUN go build -o main ./cmd/server

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/product/main .

EXPOSE 8081

CMD ["./main"]

