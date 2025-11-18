FROM golang:1.25.3-alpine3.22 AS builder

WORKDIR /app

COPY pb ./pb

COPY auth ./auth

WORKDIR /app/auth


RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o main ./cmd/server

FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /root/

COPY --from=builder /app/auth/main .

EXPOSE 8080

CMD [ "./main" ]