FROM golang:1.24.4-alpine AS builder

WORKDIR /app

COPY go.mod ./go.mod
COPY go.sum ./go.sum
RUN go mod download

COPY main.go ./main.go

RUN CGO_ENABLED=0 go build -o /proxy main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /proxy /app/proxy

EXPOSE 20000

ENTRYPOINT ["/app/proxy"]
