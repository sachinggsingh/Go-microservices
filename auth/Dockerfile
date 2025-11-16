FROM golang:1.25.3-alpine3.22 AS builder

WORKDIR /app

COPY pb ./pb

COPY auth ./auth

WORKDIR /app/auth


RUN go mod download

RUN go build -o main ./cmd/server


FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/auth/main .

EXPOSE 8080

CMD [ "./main" ]