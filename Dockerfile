# Multi-stage Dockerfile for AI Agentic Browser
# This Dockerfile builds all Go services in a single image with service selection at runtime

# Build stage
FROM golang:1.22-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build all services
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/api-gateway ./cmd/api-gateway
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/ai-agent ./cmd/ai-agent
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/browser-service ./cmd/browser-service
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/web3-service ./cmd/web3-service
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/auth-service ./cmd/auth-service

# Runtime stage
FROM alpine:3.18

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata curl netcat-openbsd

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Copy binaries from builder stage
COPY --from=builder /app/bin/ ./bin/

# Copy configuration files
COPY --from=builder /app/configs/ ./configs/
COPY --from=builder /app/migrations/ ./migrations/

# Copy startup script
COPY scripts/docker-entrypoint.sh ./
RUN chmod +x docker-entrypoint.sh

# Create directories for logs and data
RUN mkdir -p /app/logs /app/data && \
    chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:${PORT:-8080}/health || exit 1

# Expose port (will be overridden by service-specific values)
EXPOSE 8080

# Default command (will be overridden by service-specific commands)
ENTRYPOINT ["./docker-entrypoint.sh"]
CMD ["api-gateway"]
