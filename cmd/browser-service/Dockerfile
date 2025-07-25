# Build stage
FROM golang:1.22-alpine AS builder

# Install git and ca-certificates
RUN apk add --no-cache git ca-certificates

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o browser-service ./cmd/browser-service

# Final stage - Use a Chrome-compatible base image
FROM zenika/alpine-chrome:with-chromedriver

# Switch to root to install packages
USER root

# Install ca-certificates and other dependencies
RUN apk --no-cache add ca-certificates dumb-init

# Create app directory
WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/browser-service .

# Create non-root user for the app
RUN adduser -D -s /bin/sh appuser && \
    chown appuser:appuser browser-service

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8083

# Use dumb-init to handle signals properly
ENTRYPOINT ["dumb-init", "--"]

# Run the application
CMD ["./browser-service"]
