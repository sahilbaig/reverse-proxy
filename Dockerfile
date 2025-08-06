# Development image
FROM golang:1.24.5-alpine

WORKDIR /app

# Install air (using the correct package)
RUN go install github.com/air-verse/air@latest

# Copy files
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Default air command (no config file needed)
CMD ["air"]