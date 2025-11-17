FROM golang:1.25.3-alpine3.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cart ./cart

COPY ../pb ./pb

WORKDIR /app/cart

RUN go build -o main ./cmd/server

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/cart/main .

EXPOSE 8082

CMD ["./main"]
