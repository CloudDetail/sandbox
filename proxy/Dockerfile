
# Step 1: Build the application
FROM golang:1.24.4-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum and download dependencies
COPY go.mod ./go.mod
COPY go.sum ./go.sum
RUN go mod download

# Copy the source code
COPY main.go ./main.go

# Build the application
RUN CGO_ENABLED=0 go build -o /proxy main.go

# Step 2: Run the application
FROM alpine:latest

WORKDIR /app

# Copy the built executable from the builder stage
COPY --from=builder /proxy /app/proxy

# Expose the port the proxy listens on (default 20000)
EXPOSE 20000

# Run the application
ENTRYPOINT ["/app/proxy"]
