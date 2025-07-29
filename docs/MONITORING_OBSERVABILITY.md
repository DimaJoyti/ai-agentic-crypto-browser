# AI Agentic Crypto Browser - Monitoring & Observability Guide

## 📊 Overview

This guide provides comprehensive monitoring and observability strategies for the AI Agentic Crypto Browser, ensuring optimal performance, reliability, and operational insights across all system components.

## 🎯 Observability Strategy

### Three Pillars of Observability

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│     METRICS     │    │      LOGS       │    │     TRACES      │
│                 │    │                 │    │                 │
│ • Performance   │    │ • Events        │    │ • Request Flow  │
│ • Business KPIs │    │ • Errors        │    │ • Dependencies  │
│ • System Health │    │ • Audit Trail   │    │ • Latency       │
│ • Alerts        │    │ • Debug Info    │    │ • Bottlenecks   │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 │
                    ┌─────────────────┐
                    │   DASHBOARDS    │
                    │                 │
                    │ • Real-time     │
                    │ • Historical    │
                    │ • Alerting      │
                    │ • Analytics     │
                    └─────────────────┘
```

## 📈 Metrics Collection

### Core System Metrics

```go
// Prometheus metrics for the AI system
var (
    // HTTP metrics
    httpRequestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "endpoint", "status_code"},
    )
    
    httpRequestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "HTTP request duration in seconds",
            Buckets: prometheus.ExponentialBuckets(0.001, 2, 15),
        },
        []string{"method", "endpoint"},
    )
    
    // AI Engine metrics
    aiRequestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "ai_requests_total",
            Help: "Total number of AI requests",
        },
        []string{"engine_type", "request_type", "status"},
    )
    
    aiProcessingDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "ai_processing_duration_seconds",
            Help:    "AI processing duration in seconds",
            Buckets: prometheus.ExponentialBuckets(0.01, 2, 12),
        },
        []string{"engine_type", "request_type"},
    )
    
    // Market Pattern metrics
    patternsDetected = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "patterns_detected_total",
            Help: "Total number of patterns detected",
        },
        []string{"pattern_type", "asset", "confidence_level"},
    )
    
    strategyAdaptations = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "strategy_adaptations_total",
            Help: "Total number of strategy adaptations",
        },
        []string{"strategy_type", "adaptation_reason", "success"},
    )
    
    // Business metrics
    userBehaviorEvents = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "user_behavior_events_total",
            Help: "Total number of user behavior events",
        },
        []string{"event_type", "user_segment"},
    )
    
    marketDataLatency = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "market_data_latency_seconds",
            Help:    "Market data processing latency",
            Buckets: prometheus.ExponentialBuckets(0.001, 2, 10),
        },
        []string{"data_source", "asset"},
    )
)

// Metrics middleware
func MetricsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        // Wrap response writer to capture status code
        wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
        
        next.ServeHTTP(wrapped, r)
        
        // Record metrics
        duration := time.Since(start).Seconds()
        endpoint := getEndpointLabel(r.URL.Path)
        
        httpRequestsTotal.WithLabelValues(
            r.Method,
            endpoint,
            fmt.Sprintf("%d", wrapped.statusCode),
        ).Inc()
        
        httpRequestDuration.WithLabelValues(
            r.Method,
            endpoint,
        ).Observe(duration)
    })
}
```

### Custom Business Metrics

```go
type BusinessMetrics struct {
    // Trading performance metrics
    TotalReturn          prometheus.Gauge
    SharpeRatio         prometheus.Gauge
    MaxDrawdown         prometheus.Gauge
    WinRate             prometheus.Gauge
    
    // Pattern detection metrics
    PatternAccuracy     prometheus.Gauge
    FalsePositiveRate   prometheus.Gauge
    PatternCoverage     prometheus.Gauge
    
    // User engagement metrics
    ActiveUsers         prometheus.Gauge
    SessionDuration     prometheus.Histogram
    FeatureUsage        prometheus.CounterVec
}

func (bm *BusinessMetrics) UpdateTradingMetrics(performance *PerformanceMetrics) {
    bm.TotalReturn.Set(performance.TotalReturn)
    bm.SharpeRatio.Set(performance.SharpeRatio)
    bm.MaxDrawdown.Set(performance.MaxDrawdown)
    bm.WinRate.Set(performance.WinRate)
}

func (bm *BusinessMetrics) UpdatePatternMetrics(accuracy, falsePositive, coverage float64) {
    bm.PatternAccuracy.Set(accuracy)
    bm.FalsePositiveRate.Set(falsePositive)
    bm.PatternCoverage.Set(coverage)
}
```

## 📝 Structured Logging

### Comprehensive Logging Strategy

```go
type StructuredLogger struct {
    logger     *slog.Logger
    traceID    string
    userID     string
    component  string
    version    string
}

func NewStructuredLogger(component string) *StructuredLogger {
    handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level: slog.LevelInfo,
        AddSource: true,
    })
    
    logger := slog.New(handler)
    
    return &StructuredLogger{
        logger:    logger,
        component: component,
        version:   GetVersion(),
    }
}

func (sl *StructuredLogger) WithContext(ctx context.Context) *StructuredLogger {
    traceID := GetTraceIDFromContext(ctx)
    userID := GetUserIDFromContext(ctx)
    
    return &StructuredLogger{
        logger:    sl.logger,
        traceID:   traceID,
        userID:    userID,
        component: sl.component,
        version:   sl.version,
    }
}

func (sl *StructuredLogger) Info(msg string, fields map[string]interface{}) {
    attrs := sl.buildAttributes(fields)
    sl.logger.Info(msg, attrs...)
}

func (sl *StructuredLogger) Error(msg string, err error, fields map[string]interface{}) {
    attrs := sl.buildAttributes(fields)
    attrs = append(attrs, slog.String("error", err.Error()))
    sl.logger.Error(msg, attrs...)
}

func (sl *StructuredLogger) buildAttributes(fields map[string]interface{}) []slog.Attr {
    attrs := []slog.Attr{
        slog.String("component", sl.component),
        slog.String("version", sl.version),
        slog.Time("timestamp", time.Now()),
    }
    
    if sl.traceID != "" {
        attrs = append(attrs, slog.String("trace_id", sl.traceID))
    }
    
    if sl.userID != "" {
        attrs = append(attrs, slog.String("user_id", sl.userID))
    }
    
    for key, value := range fields {
        attrs = append(attrs, slog.Any(key, value))
    }
    
    return attrs
}

// Logging middleware
func LoggingMiddleware(logger *StructuredLogger) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            traceID := uuid.New().String()
            
            // Add trace ID to context
            ctx := context.WithValue(r.Context(), "trace_id", traceID)
            r = r.WithContext(ctx)
            
            // Add trace ID to response headers
            w.Header().Set("X-Trace-ID", traceID)
            
            wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
            
            next.ServeHTTP(wrapped, r)
            
            duration := time.Since(start)
            
            logger.WithContext(ctx).Info("HTTP request", map[string]interface{}{
                "method":      r.Method,
                "path":        r.URL.Path,
                "status_code": wrapped.statusCode,
                "duration_ms": duration.Milliseconds(),
                "user_agent":  r.UserAgent(),
                "remote_addr": r.RemoteAddr,
            })
        })
    }
}
```

### Log Aggregation and Analysis

```go
// Log aggregation for pattern analysis
type LogAggregator struct {
    elasticsearch *elasticsearch.Client
    buffer        []LogEntry
    bufferSize    int
    flushInterval time.Duration
    mutex         sync.Mutex
}

type LogEntry struct {
    Timestamp time.Time              `json:"@timestamp"`
    Level     string                 `json:"level"`
    Message   string                 `json:"message"`
    Component string                 `json:"component"`
    TraceID   string                 `json:"trace_id"`
    UserID    string                 `json:"user_id"`
    Fields    map[string]interface{} `json:"fields"`
}

func (la *LogAggregator) AddLog(entry LogEntry) {
    la.mutex.Lock()
    defer la.mutex.Unlock()
    
    la.buffer = append(la.buffer, entry)
    
    if len(la.buffer) >= la.bufferSize {
        go la.flush()
    }
}

func (la *LogAggregator) flush() {
    la.mutex.Lock()
    entries := make([]LogEntry, len(la.buffer))
    copy(entries, la.buffer)
    la.buffer = la.buffer[:0]
    la.mutex.Unlock()
    
    // Bulk index to Elasticsearch
    var buf bytes.Buffer
    for _, entry := range entries {
        meta := map[string]interface{}{
            "index": map[string]interface{}{
                "_index": fmt.Sprintf("ai-browser-logs-%s", time.Now().Format("2006.01.02")),
            },
        }
        
        metaJSON, _ := json.Marshal(meta)
        entryJSON, _ := json.Marshal(entry)
        
        buf.Write(metaJSON)
        buf.WriteByte('\n')
        buf.Write(entryJSON)
        buf.WriteByte('\n')
    }
    
    la.elasticsearch.Bulk(bytes.NewReader(buf.Bytes()))
}
```

## 🔍 Distributed Tracing

### OpenTelemetry Integration

```go
import (
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/attribute"
    "go.opentelemetry.io/otel/exporters/jaeger"
    "go.opentelemetry.io/otel/sdk/resource"
    "go.opentelemetry.io/otel/sdk/trace"
    semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func InitTracing(serviceName, jaegerEndpoint string) error {
    // Create Jaeger exporter
    exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(jaegerEndpoint)))
    if err != nil {
        return fmt.Errorf("failed to create Jaeger exporter: %w", err)
    }
    
    // Create trace provider
    tp := trace.NewTracerProvider(
        trace.WithBatcher(exp),
        trace.WithResource(resource.NewWithAttributes(
            semconv.SchemaURL,
            semconv.ServiceNameKey.String(serviceName),
            semconv.ServiceVersionKey.String(GetVersion()),
        )),
    )
    
    otel.SetTracerProvider(tp)
    return nil
}

// Tracing for AI operations
func (m *MarketAdaptationEngine) DetectPatternsWithTracing(ctx context.Context, data map[string]interface{}) ([]*DetectedPattern, error) {
    tracer := otel.Tracer("market-adaptation")
    ctx, span := tracer.Start(ctx, "detect_patterns")
    defer span.End()
    
    // Add attributes
    span.SetAttributes(
        attribute.String("operation", "pattern_detection"),
        attribute.Int("data_points", len(data)),
    )
    
    // Detect patterns
    patterns, err := m.detectPatterns(ctx, data)
    if err != nil {
        span.RecordError(err)
        span.SetStatus(codes.Error, err.Error())
        return nil, err
    }
    
    // Add result attributes
    span.SetAttributes(
        attribute.Int("patterns_found", len(patterns)),
        attribute.Float64("avg_confidence", calculateAvgConfidence(patterns)),
    )
    
    return patterns, nil
}

// Trace strategy adaptation
func (m *MarketAdaptationEngine) AdaptStrategyWithTracing(ctx context.Context, strategy *AdaptiveStrategy, patterns []*DetectedPattern) error {
    tracer := otel.Tracer("market-adaptation")
    ctx, span := tracer.Start(ctx, "adapt_strategy")
    defer span.End()
    
    span.SetAttributes(
        attribute.String("strategy_id", strategy.ID),
        attribute.String("strategy_type", strategy.Type),
        attribute.Int("pattern_count", len(patterns)),
    )
    
    // Create child span for parameter calculation
    _, paramSpan := tracer.Start(ctx, "calculate_parameters")
    newParams, err := m.calculateOptimalParameters(ctx, strategy, patterns)
    paramSpan.End()
    
    if err != nil {
        span.RecordError(err)
        return err
    }
    
    // Create child span for parameter update
    _, updateSpan := tracer.Start(ctx, "update_parameters")
    err = m.updateStrategyParameters(ctx, strategy, newParams)
    updateSpan.End()
    
    if err != nil {
        span.RecordError(err)
        return err
    }
    
    span.SetAttributes(
        attribute.Bool("adaptation_successful", true),
        attribute.Int("parameters_changed", len(newParams)),
    )
    
    return nil
}
```

## 📊 Dashboard Configuration

### Grafana Dashboard JSON

```json
{
  "dashboard": {
    "id": null,
    "title": "AI Agentic Crypto Browser - Overview",
    "tags": ["ai", "crypto", "monitoring"],
    "timezone": "browser",
    "panels": [
      {
        "id": 1,
        "title": "Request Rate",
        "type": "stat",
        "targets": [
          {
            "expr": "rate(http_requests_total[5m])",
            "legendFormat": "{{method}} {{endpoint}}"
          }
        ],
        "fieldConfig": {
          "defaults": {
            "unit": "reqps",
            "thresholds": {
              "steps": [
                {"color": "green", "value": null},
                {"color": "yellow", "value": 100},
                {"color": "red", "value": 500}
              ]
            }
          }
        }
      },
      {
        "id": 2,
        "title": "AI Processing Latency",
        "type": "graph",
        "targets": [
          {
            "expr": "histogram_quantile(0.95, rate(ai_processing_duration_seconds_bucket[5m]))",
            "legendFormat": "95th percentile"
          },
          {
            "expr": "histogram_quantile(0.50, rate(ai_processing_duration_seconds_bucket[5m]))",
            "legendFormat": "50th percentile"
          }
        ]
      },
      {
        "id": 3,
        "title": "Pattern Detection Success Rate",
        "type": "stat",
        "targets": [
          {
            "expr": "rate(patterns_detected_total[5m])",
            "legendFormat": "{{pattern_type}}"
          }
        ]
      },
      {
        "id": 4,
        "title": "Strategy Adaptation Rate",
        "type": "graph",
        "targets": [
          {
            "expr": "rate(strategy_adaptations_total{success=\"true\"}[5m])",
            "legendFormat": "Successful"
          },
          {
            "expr": "rate(strategy_adaptations_total{success=\"false\"}[5m])",
            "legendFormat": "Failed"
          }
        ]
      }
    ],
    "time": {
      "from": "now-1h",
      "to": "now"
    },
    "refresh": "5s"
  }
}
```

## 🚨 Alerting Rules

### Prometheus Alerting Rules

```yaml
# alerts.yml
groups:
- name: ai-browser-alerts
  rules:
  # High error rate
  - alert: HighErrorRate
    expr: rate(http_requests_total{status_code=~"5.."}[5m]) > 0.1
    for: 5m
    labels:
      severity: critical
      service: ai-browser
    annotations:
      summary: "High error rate detected"
      description: "Error rate is {{ $value }} errors per second"
      
  # AI processing latency
  - alert: HighAILatency
    expr: histogram_quantile(0.95, rate(ai_processing_duration_seconds_bucket[5m])) > 5
    for: 10m
    labels:
      severity: warning
      service: ai-browser
    annotations:
      summary: "High AI processing latency"
      description: "95th percentile latency is {{ $value }}s"
      
  # Pattern detection failure
  - alert: PatternDetectionFailure
    expr: rate(patterns_detected_total[10m]) == 0
    for: 15m
    labels:
      severity: critical
      service: ai-browser
    annotations:
      summary: "Pattern detection stopped"
      description: "No patterns detected in the last 15 minutes"
      
  # Strategy adaptation failure
  - alert: StrategyAdaptationFailure
    expr: rate(strategy_adaptations_total{success="false"}[5m]) > 0.1
    for: 5m
    labels:
      severity: warning
      service: ai-browser
    annotations:
      summary: "High strategy adaptation failure rate"
      description: "Strategy adaptation failure rate is {{ $value }}"
      
  # Memory usage
  - alert: HighMemoryUsage
    expr: process_resident_memory_bytes / 1024 / 1024 > 2000
    for: 10m
    labels:
      severity: warning
      service: ai-browser
    annotations:
      summary: "High memory usage"
      description: "Memory usage is {{ $value }}MB"
      
  # Database connection issues
  - alert: DatabaseConnectionFailure
    expr: up{job="postgres"} == 0
    for: 1m
    labels:
      severity: critical
      service: ai-browser
    annotations:
      summary: "Database connection failed"
      description: "Cannot connect to PostgreSQL database"
```

## 📱 Alert Management

### Alert Manager Configuration

```go
type AlertManager struct {
    webhookURL    string
    emailSMTP     *SMTPConfig
    slackWebhook  string
    pagerDutyKey  string
    logger        *StructuredLogger
}

type Alert struct {
    ID          string                 `json:"id"`
    Title       string                 `json:"title"`
    Description string                 `json:"description"`
    Severity    string                 `json:"severity"`
    Service     string                 `json:"service"`
    Timestamp   time.Time              `json:"timestamp"`
    Labels      map[string]string      `json:"labels"`
    Annotations map[string]string      `json:"annotations"`
    Status      string                 `json:"status"`
}

func (am *AlertManager) SendAlert(ctx context.Context, alert *Alert) error {
    // Log alert
    am.logger.WithContext(ctx).Error("Alert triggered", nil, map[string]interface{}{
        "alert_id":    alert.ID,
        "title":       alert.Title,
        "severity":    alert.Severity,
        "service":     alert.Service,
    })
    
    // Send to multiple channels based on severity
    switch alert.Severity {
    case "critical":
        go am.sendToPagerDuty(alert)
        go am.sendToSlack(alert)
        go am.sendEmail(alert)
    case "warning":
        go am.sendToSlack(alert)
        go am.sendEmail(alert)
    case "info":
        go am.sendToSlack(alert)
    }
    
    return nil
}

func (am *AlertManager) sendToSlack(alert *Alert) error {
    payload := map[string]interface{}{
        "text": fmt.Sprintf("🚨 *%s*", alert.Title),
        "attachments": []map[string]interface{}{
            {
                "color": am.getSeverityColor(alert.Severity),
                "fields": []map[string]interface{}{
                    {"title": "Service", "value": alert.Service, "short": true},
                    {"title": "Severity", "value": alert.Severity, "short": true},
                    {"title": "Description", "value": alert.Description, "short": false},
                },
                "ts": alert.Timestamp.Unix(),
            },
        },
    }
    
    jsonPayload, _ := json.Marshal(payload)
    resp, err := http.Post(am.slackWebhook, "application/json", bytes.NewBuffer(jsonPayload))
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    return nil
}
```

This comprehensive monitoring and observability guide ensures complete visibility into the AI Agentic Crypto Browser's performance, health, and business metrics, enabling proactive issue detection and resolution.
