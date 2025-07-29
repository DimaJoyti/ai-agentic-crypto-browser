# Multi-stage build for AI Agentic Crypto Browser
FROM golang:1.21-alpine AS builder

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

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ai-agent ./cmd/ai-agent/
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o demo-client ./cmd/demo-client/
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o market-adaptation-demo ./cmd/market-adaptation-demo/

# Production stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Copy certificates and timezone data
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy binaries from builder stage
COPY --from=builder /app/ai-agent /app/
COPY --from=builder /app/demo-client /app/
COPY --from=builder /app/market-adaptation-demo /app/

# Copy configuration files
COPY --from=builder /app/configs/ /app/configs/
COPY --from=builder /app/docs/ /app/docs/

# Create necessary directories
RUN mkdir -p /app/logs /app/data && \
    chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose ports
EXPOSE 8080 8081 8082

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Default command
CMD ["./ai-agent"]
