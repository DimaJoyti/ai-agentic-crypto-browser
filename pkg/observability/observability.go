package observability

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/ai-agentic-browser/internal/config"
)

// SimpleObservabilityProvider manages basic observability
type SimpleObservabilityProvider struct {
	Logger *Logger
	config *SimpleObservabilityConfig
}

// SimpleObservabilityConfig contains basic configuration
type SimpleObservabilityConfig struct {
	ServiceName    string
	ServiceVersion string
	Environment    string
	LogLevel       string
	LogFormat      string
}

// NewSimpleObservabilityProvider creates a new simple observability provider
func NewSimpleObservabilityProvider(cfg *SimpleObservabilityConfig) (*SimpleObservabilityProvider, error) {
	if cfg == nil {
		cfg = &SimpleObservabilityConfig{
			ServiceName:    "unknown-service",
			ServiceVersion: "unknown",
			Environment:    "development",
			LogLevel:       "info",
			LogFormat:      "json",
		}
	}

	op := &SimpleObservabilityProvider{
		config: cfg,
	}

	// Initialize logger
	loggerConfig := config.ObservabilityConfig{
		ServiceName: cfg.ServiceName,
		LogLevel:    cfg.LogLevel,
		LogFormat:   cfg.LogFormat,
	}
	logger := NewLogger(loggerConfig)
	op.Logger = logger

	return op, nil
}

// Start starts all observability components
func (op *SimpleObservabilityProvider) Start(ctx context.Context) error {
	op.Logger.Info(ctx, "Simple observability provider started", map[string]interface{}{
		"service":     op.config.ServiceName,
		"version":     op.config.ServiceVersion,
		"environment": op.config.Environment,
	})
	return nil
}

// Stop stops all observability components
func (op *SimpleObservabilityProvider) Stop(ctx context.Context) error {
	op.Logger.Info(ctx, "Simple observability provider stopped")
	return nil
}

// GetHTTPMiddleware returns HTTP middleware for observability
func (op *SimpleObservabilityProvider) GetHTTPMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Add request ID if not present
			requestID := r.Header.Get("X-Request-ID")
			if requestID == "" {
				requestID = generateRequestID()
				r.Header.Set("X-Request-ID", requestID)
			}

			// Create context with request ID
			ctx := context.WithValue(r.Context(), "request_id", requestID)
			r = r.WithContext(ctx)

			// Wrap response writer to capture status code
			wrapped := &simpleResponseWriter{ResponseWriter: w, statusCode: 200}

			// Process request
			next.ServeHTTP(wrapped, r)

			// Log request
			duration := time.Since(start)
			op.Logger.Info(ctx, "HTTP request", map[string]interface{}{
				"method":      r.Method,
				"path":        r.URL.Path,
				"status_code": wrapped.statusCode,
				"duration_ms": duration.Milliseconds(),
				"request_id":  requestID,
			})
		})
	}
}

// simpleResponseWriter wraps http.ResponseWriter to capture status code
type simpleResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *simpleResponseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// generateRequestID generates a unique request ID
func generateRequestID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// GetDefaultSimpleConfig returns default observability configuration
func GetDefaultSimpleConfig() *SimpleObservabilityConfig {
	return &SimpleObservabilityConfig{
		ServiceName:    getEnv("SERVICE_NAME", "unknown-service"),
		ServiceVersion: getEnv("SERVICE_VERSION", "unknown"),
		Environment:    getEnv("ENVIRONMENT", "development"),
		LogLevel:       getEnv("LOG_LEVEL", "info"),
		LogFormat:      getEnv("LOG_FORMAT", "json"),
	}
}

// Helper functions for environment variables
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Note: Logger and TracingProvider implementations are in their respective files:
// - logging.go for Logger implementation
// - tracing.go for TracingProvider implementation
