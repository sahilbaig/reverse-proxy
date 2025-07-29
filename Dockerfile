# ====== Builder stage ======
FROM golang:1.24.5-alpine AS builder

WORKDIR /app

# Copy Go mod files first (better caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build binary
RUN go build -o server main.go

# ====== Final image ======
FROM alpine:latest

WORKDIR /app

# Copy the built binary
COPY --from=builder /app/server .

# Expose port
EXPOSE 7001

# Set default env var (can be overridden during runtime)
ENV PROXY_TARGET=http://localhost:8080

# Run the binary
ENTRYPOINT ["./server"]
