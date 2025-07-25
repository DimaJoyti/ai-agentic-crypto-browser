package observability

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ai-agentic-browser/internal/config"
	"go.opentelemetry.io/otel/trace"
)

// LogLevel represents the severity level of a log entry
type LogLevel string

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
)

// LogEntry represents a structured log entry
type LogEntry struct {
	Timestamp string                 `json:"timestamp"`
	Level     LogLevel               `json:"level"`
	Message   string                 `json:"message"`
	Service   string                 `json:"service"`
	TraceID   string                 `json:"trace_id,omitempty"`
	SpanID    string                 `json:"span_id,omitempty"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
	Error     string                 `json:"error,omitempty"`
}

// Logger provides structured logging with OpenTelemetry integration
type Logger struct {
	serviceName string
	logLevel    LogLevel
	format      string
}

// NewLogger creates a new structured logger
func NewLogger(cfg config.ObservabilityConfig) *Logger {
	return &Logger{
		serviceName: cfg.ServiceName,
		logLevel:    LogLevel(cfg.LogLevel),
		format:      cfg.LogFormat,
	}
}

// Debug logs a debug message
func (l *Logger) Debug(ctx context.Context, message string, fields ...map[string]interface{}) {
	if l.shouldLog(LogLevelDebug) {
		l.log(ctx, LogLevelDebug, message, nil, fields...)
	}
}

// Info logs an info message
func (l *Logger) Info(ctx context.Context, message string, fields ...map[string]interface{}) {
	if l.shouldLog(LogLevelInfo) {
		l.log(ctx, LogLevelInfo, message, nil, fields...)
	}
}

// Warn logs a warning message
func (l *Logger) Warn(ctx context.Context, message string, fields ...map[string]interface{}) {
	if l.shouldLog(LogLevelWarn) {
		l.log(ctx, LogLevelWarn, message, nil, fields...)
	}
}

// Error logs an error message
func (l *Logger) Error(ctx context.Context, message string, err error, fields ...map[string]interface{}) {
	if l.shouldLog(LogLevelError) {
		l.log(ctx, LogLevelError, message, err, fields...)
	}
}

// log is the internal logging method
func (l *Logger) log(ctx context.Context, level LogLevel, message string, err error, fields ...map[string]interface{}) {
	entry := LogEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Level:     level,
		Message:   message,
		Service:   l.serviceName,
	}

	// Extract trace information from context
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		entry.TraceID = span.SpanContext().TraceID().String()
		entry.SpanID = span.SpanContext().SpanID().String()
	}

	// Add error if present
	if err != nil {
		entry.Error = err.Error()
	}

	// Merge all field maps
	if len(fields) > 0 {
		entry.Fields = make(map[string]interface{})
		for _, fieldMap := range fields {
			for k, v := range fieldMap {
				entry.Fields[k] = v
			}
		}
	}

	// Output the log entry
	l.output(entry)
}

// output writes the log entry to stdout
func (l *Logger) output(entry LogEntry) {
	if l.format == "json" {
		if data, err := json.Marshal(entry); err == nil {
			fmt.Fprintln(os.Stdout, string(data))
		} else {
			log.Printf("Failed to marshal log entry: %v", err)
		}
	} else {
		// Simple text format
		fmt.Printf("[%s] %s %s: %s\n",
			entry.Timestamp,
			entry.Level,
			entry.Service,
			entry.Message)
	}
}

// shouldLog determines if a message should be logged based on the configured level
func (l *Logger) shouldLog(level LogLevel) bool {
	levels := map[LogLevel]int{
		LogLevelDebug: 0,
		LogLevelInfo:  1,
		LogLevelWarn:  2,
		LogLevelError: 3,
	}

	configuredLevel, exists := levels[l.logLevel]
	if !exists {
		configuredLevel = levels[LogLevelInfo] // Default to info
	}

	messageLevel, exists := levels[level]
	if !exists {
		return false
	}

	return messageLevel >= configuredLevel
}

// WithFields creates a new logger with additional fields
func (l *Logger) WithFields(fields map[string]interface{}) *FieldLogger {
	return &FieldLogger{
		logger: l,
		fields: fields,
	}
}

// FieldLogger is a logger with pre-set fields
type FieldLogger struct {
	logger *Logger
	fields map[string]interface{}
}

// Debug logs a debug message with pre-set fields
func (fl *FieldLogger) Debug(ctx context.Context, message string) {
	fl.logger.Debug(ctx, message, fl.fields)
}

// Info logs an info message with pre-set fields
func (fl *FieldLogger) Info(ctx context.Context, message string) {
	fl.logger.Info(ctx, message, fl.fields)
}

// Warn logs a warning message with pre-set fields
func (fl *FieldLogger) Warn(ctx context.Context, message string) {
	fl.logger.Warn(ctx, message, fl.fields)
}

// Error logs an error message with pre-set fields
func (fl *FieldLogger) Error(ctx context.Context, message string, err error) {
	fl.logger.Error(ctx, message, err, fl.fields)
}

// PerformanceLogger logs performance metrics
type PerformanceLogger struct {
	logger *Logger
}

// NewPerformanceLogger creates a new performance logger
func NewPerformanceLogger(logger *Logger) *PerformanceLogger {
	return &PerformanceLogger{logger: logger}
}

// LogDuration logs the duration of an operation
func (pl *PerformanceLogger) LogDuration(ctx context.Context, operation string, duration time.Duration, fields ...map[string]interface{}) {
	allFields := map[string]interface{}{
		"operation":   operation,
		"duration_ms": duration.Milliseconds(),
		"duration_ns": duration.Nanoseconds(),
		"component":   "performance",
	}

	for _, fieldMap := range fields {
		for k, v := range fieldMap {
			allFields[k] = v
		}
	}

	pl.logger.Info(ctx, fmt.Sprintf("Operation completed: %s", operation), allFields)
}

// LogSlowOperation logs operations that exceed a threshold
func (pl *PerformanceLogger) LogSlowOperation(ctx context.Context, operation string, duration, threshold time.Duration, fields ...map[string]interface{}) {
	if duration <= threshold {
		return
	}

	allFields := map[string]interface{}{
		"operation":    operation,
		"duration_ms":  duration.Milliseconds(),
		"threshold_ms": threshold.Milliseconds(),
		"slow_factor":  float64(duration) / float64(threshold),
		"component":    "performance",
	}

	for _, fieldMap := range fields {
		for k, v := range fieldMap {
			allFields[k] = v
		}
	}

	pl.logger.Warn(ctx, fmt.Sprintf("Slow operation detected: %s", operation), allFields)
}

// SecurityLogger logs security-related events
type SecurityLogger struct {
	logger *Logger
}

// NewSecurityLogger creates a new security logger
func NewSecurityLogger(logger *Logger) *SecurityLogger {
	return &SecurityLogger{logger: logger}
}

// LogAuthEvent logs authentication events
func (sl *SecurityLogger) LogAuthEvent(ctx context.Context, event, userID, ipAddress string, success bool, fields ...map[string]interface{}) {
	allFields := map[string]interface{}{
		"event":      event,
		"user_id":    userID,
		"ip_address": ipAddress,
		"success":    success,
		"component":  "security",
	}

	for _, fieldMap := range fields {
		for k, v := range fieldMap {
			allFields[k] = v
		}
	}

	if success {
		sl.logger.Info(ctx, fmt.Sprintf("Authentication event: %s", event), allFields)
	} else {
		sl.logger.Warn(ctx, fmt.Sprintf("Failed authentication event: %s", event), allFields)
	}
}

// LogSecurityViolation logs security violations
func (sl *SecurityLogger) LogSecurityViolation(ctx context.Context, violation, userID, ipAddress string, severity string, fields ...map[string]interface{}) {
	allFields := map[string]interface{}{
		"violation":  violation,
		"user_id":    userID,
		"ip_address": ipAddress,
		"severity":   severity,
		"component":  "security",
	}

	for _, fieldMap := range fields {
		for k, v := range fieldMap {
			allFields[k] = v
		}
	}

	sl.logger.Error(ctx, fmt.Sprintf("Security violation: %s", violation), fmt.Errorf("security violation: %s", violation), allFields)
}

// AuditLogger logs audit events
type AuditLogger struct {
	logger *Logger
}

// NewAuditLogger creates a new audit logger
func NewAuditLogger(logger *Logger) *AuditLogger {
	return &AuditLogger{logger: logger}
}

// LogUserAction logs user actions for audit purposes
func (al *AuditLogger) LogUserAction(ctx context.Context, action, userID, resource string, fields ...map[string]interface{}) {
	allFields := map[string]interface{}{
		"action":    action,
		"user_id":   userID,
		"resource":  resource,
		"component": "audit",
	}

	for _, fieldMap := range fields {
		for k, v := range fieldMap {
			allFields[k] = v
		}
	}

	al.logger.Info(ctx, fmt.Sprintf("User action: %s", action), allFields)
}

// LogSystemEvent logs system events for audit purposes
func (al *AuditLogger) LogSystemEvent(ctx context.Context, event, component string, fields ...map[string]interface{}) {
	allFields := map[string]interface{}{
		"event":     event,
		"component": component,
		"type":      "system",
	}

	for _, fieldMap := range fields {
		for k, v := range fieldMap {
			allFields[k] = v
		}
	}

	al.logger.Info(ctx, fmt.Sprintf("System event: %s", event), allFields)
}
