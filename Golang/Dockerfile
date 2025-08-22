FROM golang:1.24-alpine AS builder
RUN apk add --no-cache iptables iproute2

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN curl -fsSL https://cdn.jsdelivr.net/gh/alibaba/loongsuite-go-agent@main/install.sh | sudo bash
RUN ./otel go build -o main .

FROM alpine:latest

RUN apk add --no-cache iptables iproute2

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 3500

CMD ["./main"]
