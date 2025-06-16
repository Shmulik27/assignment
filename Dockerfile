# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/app ./cmd/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/bin/app .
COPY --from=builder /app/config/app_configuration.json ./config/

# Create a non-root user
RUN adduser -D -g '' appuser
USER appuser

# Expose metrics and health check ports
EXPOSE 8080 8081

# Run the application
CMD ["./app"] 