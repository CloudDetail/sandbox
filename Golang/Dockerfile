FROM golang:1.24-alpine AS builder
RUN apk add --no-cache iptables iproute2
RUN wget -O otel https://github.com/alibaba/loongsuite-go-agent/releases/download/v1.0.0/otel-linux-amd64

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN ./otel go build -o main .

FROM alpine:latest

RUN apk add --no-cache iptables iproute2

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 3500

CMD ["./main"]
