package observability

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// ObservabilityMiddleware provides comprehensive observability for HTTP requests
type ObservabilityMiddleware struct {
	tracer         trace.Tracer
	metrics        *MetricsProvider
	logger         *Logger
	performanceLog *PerformanceLogger
	securityLog    *SecurityLogger
	auditLog       *AuditLogger
	serviceName    string
	slowThreshold  time.Duration
}

// MiddlewareConfig contains configuration for observability middleware
type MiddlewareConfig struct {
	ServiceName    string
	ServiceVersion string
	SlowThreshold  time.Duration
	EnableTracing  bool
	EnableMetrics  bool
	EnableLogging  bool
	EnableSecurity bool
	EnableAudit    bool
}

// NewObservabilityMiddleware creates a new observability middleware
func NewObservabilityMiddleware(
	metrics *MetricsProvider,
	logger *Logger,
	config MiddlewareConfig,
) *ObservabilityMiddleware {
	tracer := otel.Tracer(config.ServiceName)

	slowThreshold := config.SlowThreshold
	if slowThreshold == 0 {
		slowThreshold = 1 * time.Second
	}

	return &ObservabilityMiddleware{
		tracer:         tracer,
		metrics:        metrics,
		logger:         logger,
		performanceLog: NewPerformanceLogger(logger),
		securityLog:    NewSecurityLogger(logger),
		auditLog:       NewAuditLogger(logger),
		serviceName:    config.ServiceName,
		slowThreshold:  slowThreshold,
	}
}

// GinMiddleware returns a Gin middleware for observability
func (om *ObservabilityMiddleware) GinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Generate request ID
		requestID := uuid.New().String()
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)

		// Extract trace context from headers
		ctx := otel.GetTextMapPropagator().Extract(c.Request.Context(), propagation.HeaderCarrier(c.Request.Header))

		// Start span
		spanName := fmt.Sprintf("%s %s", c.Request.Method, c.FullPath())
		ctx, span := om.tracer.Start(ctx, spanName)
		defer span.End()

		// Set span attributes
		span.SetAttributes(
			attribute.String("http.method", c.Request.Method),
			attribute.String("http.url", c.Request.URL.String()),
			attribute.String("http.route", c.FullPath()),
			attribute.String("http.user_agent", c.Request.UserAgent()),
			attribute.String("http.remote_addr", c.ClientIP()),
			attribute.String("request.id", requestID),
			attribute.String("service.name", om.serviceName),
		)

		// Add trace context to Gin context
		c.Request = c.Request.WithContext(ctx)

		// Log request start
		om.logger.Info(ctx, "HTTP request started", map[string]interface{}{
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"user_agent": c.Request.UserAgent(),
			"remote_ip":  c.ClientIP(),
			"request_id": requestID,
		})

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(start)
		statusCode := c.Writer.Status()

		// Set final span attributes
		span.SetAttributes(
			attribute.Int("http.status_code", statusCode),
			attribute.Int64("http.response_size", int64(c.Writer.Size())),
			attribute.Float64("http.duration_ms", float64(duration.Nanoseconds())/1e6),
		)

		// Set span status based on HTTP status code
		if statusCode >= 400 {
			span.SetAttributes(attribute.Bool("error", true))
			if statusCode >= 500 {
				span.RecordError(fmt.Errorf("HTTP %d", statusCode))
			}
		}

		// Record metrics
		if om.metrics != nil {
			om.metrics.RecordHTTPRequest(
				ctx,
				c.Request.Method,
				c.FullPath(),
				strconv.Itoa(statusCode),
				duration,
			)
		}

		// Log request completion
		logFields := map[string]interface{}{
			"method":      c.Request.Method,
			"path":        c.Request.URL.Path,
			"status_code": statusCode,
			"duration_ms": duration.Milliseconds(),
			"request_id":  requestID,
			"user_agent":  c.Request.UserAgent(),
			"remote_ip":   c.ClientIP(),
		}

		if statusCode >= 400 {
			om.logger.Warn(ctx, "HTTP request completed with error", logFields)
		} else {
			om.logger.Info(ctx, "HTTP request completed", logFields)
		}

		// Log slow requests
		if duration > om.slowThreshold {
			om.performanceLog.LogSlowOperation(
				ctx,
				fmt.Sprintf("%s %s", c.Request.Method, c.FullPath()),
				duration,
				om.slowThreshold,
				logFields,
			)
		}

		// Security logging for authentication endpoints
		if om.isAuthEndpoint(c.FullPath()) {
			userID := om.getUserID(c)
			success := statusCode < 400

			om.securityLog.LogAuthEvent(
				ctx,
				fmt.Sprintf("%s %s", c.Request.Method, c.FullPath()),
				userID,
				c.ClientIP(),
				success,
				logFields,
			)
		}

		// Audit logging for sensitive operations
		if om.isSensitiveEndpoint(c.FullPath()) && statusCode < 400 {
			userID := om.getUserID(c)
			resource := om.extractResource(c.FullPath())

			om.auditLog.LogUserAction(
				ctx,
				fmt.Sprintf("%s %s", c.Request.Method, c.FullPath()),
				userID,
				resource,
				logFields,
			)
		}
	}
}

// HTTPMiddleware returns a standard HTTP middleware for observability
func (om *ObservabilityMiddleware) HTTPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Generate request ID
		requestID := uuid.New().String()
		w.Header().Set("X-Request-ID", requestID)

		// Extract trace context from headers
		ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))

		// Start span
		spanName := fmt.Sprintf("%s %s", r.Method, r.URL.Path)
		ctx, span := om.tracer.Start(ctx, spanName)
		defer span.End()

		// Set span attributes
		span.SetAttributes(
			attribute.String("http.method", r.Method),
			attribute.String("http.url", r.URL.String()),
			attribute.String("http.user_agent", r.UserAgent()),
			attribute.String("http.remote_addr", r.RemoteAddr),
			attribute.String("request.id", requestID),
			attribute.String("service.name", om.serviceName),
		)

		// Create response writer wrapper to capture status code and size
		rw := &responseWriter{ResponseWriter: w, statusCode: 200}

		// Add trace context to request
		r = r.WithContext(ctx)

		// Log request start
		om.logger.Info(ctx, "HTTP request started", map[string]interface{}{
			"method":      r.Method,
			"path":        r.URL.Path,
			"user_agent":  r.UserAgent(),
			"remote_addr": r.RemoteAddr,
			"request_id":  requestID,
		})

		// Process request
		next.ServeHTTP(rw, r)

		// Calculate duration
		duration := time.Since(start)
		statusCode := rw.statusCode

		// Set final span attributes
		span.SetAttributes(
			attribute.Int("http.status_code", statusCode),
			attribute.Int64("http.response_size", int64(rw.size)),
			attribute.Float64("http.duration_ms", float64(duration.Nanoseconds())/1e6),
		)

		// Set span status based on HTTP status code
		if statusCode >= 400 {
			span.SetAttributes(attribute.Bool("error", true))
			if statusCode >= 500 {
				span.RecordError(fmt.Errorf("HTTP %d", statusCode))
			}
		}

		// Record metrics
		if om.metrics != nil {
			om.metrics.RecordHTTPRequest(
				ctx,
				r.Method,
				r.URL.Path,
				strconv.Itoa(statusCode),
				duration,
			)
		}

		// Log request completion
		logFields := map[string]interface{}{
			"method":      r.Method,
			"path":        r.URL.Path,
			"status_code": statusCode,
			"duration_ms": duration.Milliseconds(),
			"request_id":  requestID,
			"user_agent":  r.UserAgent(),
			"remote_addr": r.RemoteAddr,
		}

		if statusCode >= 400 {
			om.logger.Warn(ctx, "HTTP request completed with error", logFields)
		} else {
			om.logger.Info(ctx, "HTTP request completed", logFields)
		}

		// Log slow requests
		if duration > om.slowThreshold {
			om.performanceLog.LogSlowOperation(
				ctx,
				fmt.Sprintf("%s %s", r.Method, r.URL.Path),
				duration,
				om.slowThreshold,
				logFields,
			)
		}
	})
}

// responseWriter wraps http.ResponseWriter to capture status code and response size
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *responseWriter) Write(data []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(data)
	rw.size += size
	return size, err
}

// Helper methods

func (om *ObservabilityMiddleware) isAuthEndpoint(path string) bool {
	authPaths := []string{
		"/api/auth/login",
		"/api/auth/register",
		"/api/auth/logout",
		"/api/auth/refresh",
		"/api/auth/mfa",
	}

	for _, authPath := range authPaths {
		if path == authPath {
			return true
		}
	}
	return false
}

func (om *ObservabilityMiddleware) isSensitiveEndpoint(path string) bool {
	sensitivePaths := []string{
		"/api/workflows",
		"/api/users",
		"/api/teams",
		"/api/settings",
		"/api/admin",
	}

	for _, sensitivePath := range sensitivePaths {
		if len(path) >= len(sensitivePath) && path[:len(sensitivePath)] == sensitivePath {
			return true
		}
	}
	return false
}

func (om *ObservabilityMiddleware) getUserID(c *gin.Context) string {
	if userID, exists := c.Get("user_id"); exists {
		if uid, ok := userID.(string); ok {
			return uid
		}
		if uid, ok := userID.(uuid.UUID); ok {
			return uid.String()
		}
	}
	return "anonymous"
}

func (om *ObservabilityMiddleware) extractResource(path string) string {
	// Extract resource from path (e.g., /api/workflows/123 -> workflows)
	parts := []string{}
	for _, part := range []string{"api", "workflows", "users", "teams", "settings"} {
		if len(path) > len(part) && path[1:len(part)+1] == part {
			parts = append(parts, part)
		}
	}

	if len(parts) > 1 {
		return parts[1] // Return the resource type
	}
	return "unknown"
}

// TraceMiddleware provides basic tracing without full observability
func TraceMiddleware(serviceName string) gin.HandlerFunc {
	tracer := otel.Tracer(serviceName)

	return func(c *gin.Context) {
		// Extract trace context from headers
		ctx := otel.GetTextMapPropagator().Extract(c.Request.Context(), propagation.HeaderCarrier(c.Request.Header))

		// Start span
		spanName := fmt.Sprintf("%s %s", c.Request.Method, c.FullPath())
		ctx, span := tracer.Start(ctx, spanName)
		defer span.End()

		// Set basic span attributes
		span.SetAttributes(
			attribute.String("http.method", c.Request.Method),
			attribute.String("http.route", c.FullPath()),
			attribute.String("service.name", serviceName),
		)

		// Add trace context to Gin context
		c.Request = c.Request.WithContext(ctx)

		// Process request
		c.Next()

		// Set final span attributes
		span.SetAttributes(
			attribute.Int("http.status_code", c.Writer.Status()),
		)

		// Set span status based on HTTP status code
		if c.Writer.Status() >= 400 {
			span.SetAttributes(attribute.Bool("error", true))
			if c.Writer.Status() >= 500 {
				span.RecordError(fmt.Errorf("HTTP %d", c.Writer.Status()))
			}
		}
	}
}

// MetricsMiddleware provides basic metrics collection
func MetricsMiddleware(metrics *MetricsProvider) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process request
		c.Next()

		// Record metrics
		if metrics != nil {
			duration := time.Since(start)
			metrics.RecordHTTPRequest(
				c.Request.Context(),
				c.Request.Method,
				c.FullPath(),
				strconv.Itoa(c.Writer.Status()),
				duration,
			)
		}
	}
}
