package observability

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otelprom "go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

// MetricsProvider manages OpenTelemetry metrics and Prometheus integration
type MetricsProvider struct {
	meterProvider *sdkmetric.MeterProvider
	meter         metric.Meter
	registry      *prometheus.Registry

	// Application metrics
	httpRequestsTotal     metric.Int64Counter
	httpRequestDuration   metric.Float64Histogram
	workflowExecutions    metric.Int64Counter
	workflowDuration      metric.Float64Histogram
	aiRequestsTotal       metric.Int64Counter
	aiRequestDuration     metric.Float64Histogram
	browserSessionsActive metric.Int64UpDownCounter
	web3TransactionsTotal metric.Int64Counter
	errorRate             metric.Float64Gauge
	systemResourceUsage   metric.Float64Gauge
}

// MetricsConfig contains metrics configuration
type MetricsConfig struct {
	ServiceName    string
	ServiceVersion string
	Namespace      string
	Port           int
	Enabled        bool
}

// NewMetricsProvider creates a new metrics provider
func NewMetricsProvider(cfg MetricsConfig) (*MetricsProvider, error) {
	if !cfg.Enabled {
		return &MetricsProvider{}, nil
	}

	// Create Prometheus registry
	registry := prometheus.NewRegistry()

	// Create Prometheus exporter
	exporter, err := otelprom.New(
		otelprom.WithRegisterer(registry),
		otelprom.WithNamespace(cfg.Namespace),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create Prometheus exporter: %w", err)
	}

	// Create resource
	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(cfg.ServiceName),
			semconv.ServiceVersionKey.String(cfg.ServiceVersion),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Create meter provider
	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(exporter),
	)

	// Set global meter provider
	otel.SetMeterProvider(meterProvider)

	// Create meter
	meter := meterProvider.Meter(cfg.ServiceName)

	// Initialize metrics
	mp := &MetricsProvider{
		meterProvider: meterProvider,
		meter:         meter,
		registry:      registry,
	}

	if err := mp.initializeMetrics(); err != nil {
		return nil, fmt.Errorf("failed to initialize metrics: %w", err)
	}

	return mp, nil
}

// initializeMetrics creates all application metrics
func (mp *MetricsProvider) initializeMetrics() error {
	var err error

	// HTTP metrics
	mp.httpRequestsTotal, err = mp.meter.Int64Counter(
		"http_requests_total",
		metric.WithDescription("Total number of HTTP requests"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return fmt.Errorf("failed to create http_requests_total counter: %w", err)
	}

	mp.httpRequestDuration, err = mp.meter.Float64Histogram(
		"http_request_duration_seconds",
		metric.WithDescription("HTTP request duration in seconds"),
		metric.WithUnit("s"),
		metric.WithExplicitBucketBoundaries(0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10),
	)
	if err != nil {
		return fmt.Errorf("failed to create http_request_duration histogram: %w", err)
	}

	// Workflow metrics
	mp.workflowExecutions, err = mp.meter.Int64Counter(
		"workflow_executions_total",
		metric.WithDescription("Total number of workflow executions"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return fmt.Errorf("failed to create workflow_executions_total counter: %w", err)
	}

	mp.workflowDuration, err = mp.meter.Float64Histogram(
		"workflow_execution_duration_seconds",
		metric.WithDescription("Workflow execution duration in seconds"),
		metric.WithUnit("s"),
		metric.WithExplicitBucketBoundaries(1, 5, 10, 30, 60, 120, 300, 600, 1200, 3600),
	)
	if err != nil {
		return fmt.Errorf("failed to create workflow_duration histogram: %w", err)
	}

	// AI metrics
	mp.aiRequestsTotal, err = mp.meter.Int64Counter(
		"ai_requests_total",
		metric.WithDescription("Total number of AI requests"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return fmt.Errorf("failed to create ai_requests_total counter: %w", err)
	}

	mp.aiRequestDuration, err = mp.meter.Float64Histogram(
		"ai_request_duration_seconds",
		metric.WithDescription("AI request duration in seconds"),
		metric.WithUnit("s"),
		metric.WithExplicitBucketBoundaries(0.1, 0.5, 1, 2, 5, 10, 20, 30, 60),
	)
	if err != nil {
		return fmt.Errorf("failed to create ai_request_duration histogram: %w", err)
	}

	// Browser metrics
	mp.browserSessionsActive, err = mp.meter.Int64UpDownCounter(
		"browser_sessions_active",
		metric.WithDescription("Number of active browser sessions"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return fmt.Errorf("failed to create browser_sessions_active gauge: %w", err)
	}

	// Web3 metrics
	mp.web3TransactionsTotal, err = mp.meter.Int64Counter(
		"web3_transactions_total",
		metric.WithDescription("Total number of Web3 transactions"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return fmt.Errorf("failed to create web3_transactions_total counter: %w", err)
	}

	// Error rate gauge
	mp.errorRate, err = mp.meter.Float64Gauge(
		"error_rate",
		metric.WithDescription("Current error rate percentage"),
		metric.WithUnit("%"),
	)
	if err != nil {
		return fmt.Errorf("failed to create error_rate gauge: %w", err)
	}

	// System resource usage
	mp.systemResourceUsage, err = mp.meter.Float64Gauge(
		"system_resource_usage",
		metric.WithDescription("System resource usage percentage"),
		metric.WithUnit("%"),
	)
	if err != nil {
		return fmt.Errorf("failed to create system_resource_usage gauge: %w", err)
	}

	return nil
}

// HTTP Metrics Methods

// RecordHTTPRequest records an HTTP request metric
func (mp *MetricsProvider) RecordHTTPRequest(ctx context.Context, method, path, status string, duration time.Duration) {
	if mp.httpRequestsTotal == nil {
		return
	}

	attrs := []attribute.KeyValue{
		attribute.String("method", method),
		attribute.String("path", path),
		attribute.String("status", status),
	}

	mp.httpRequestsTotal.Add(ctx, 1, metric.WithAttributes(attrs...))
	mp.httpRequestDuration.Record(ctx, duration.Seconds(), metric.WithAttributes(attrs...))
}

// Workflow Metrics Methods

// RecordWorkflowExecution records a workflow execution metric
func (mp *MetricsProvider) RecordWorkflowExecution(ctx context.Context, workflowType, status string, duration time.Duration) {
	if mp.workflowExecutions == nil {
		return
	}

	attrs := []attribute.KeyValue{
		attribute.String("workflow_type", workflowType),
		attribute.String("status", status),
	}

	mp.workflowExecutions.Add(ctx, 1, metric.WithAttributes(attrs...))
	mp.workflowDuration.Record(ctx, duration.Seconds(), metric.WithAttributes(attrs...))
}

// AI Metrics Methods

// RecordAIRequest records an AI request metric
func (mp *MetricsProvider) RecordAIRequest(ctx context.Context, provider, model, operation string, duration time.Duration, success bool) {
	if mp.aiRequestsTotal == nil {
		return
	}

	status := "success"
	if !success {
		status = "error"
	}

	attrs := []attribute.KeyValue{
		attribute.String("provider", provider),
		attribute.String("model", model),
		attribute.String("operation", operation),
		attribute.String("status", status),
	}

	mp.aiRequestsTotal.Add(ctx, 1, metric.WithAttributes(attrs...))
	mp.aiRequestDuration.Record(ctx, duration.Seconds(), metric.WithAttributes(attrs...))
}

// Browser Metrics Methods

// IncrementBrowserSessions increments active browser sessions
func (mp *MetricsProvider) IncrementBrowserSessions(ctx context.Context) {
	if mp.browserSessionsActive == nil {
		return
	}
	mp.browserSessionsActive.Add(ctx, 1)
}

// DecrementBrowserSessions decrements active browser sessions
func (mp *MetricsProvider) DecrementBrowserSessions(ctx context.Context) {
	if mp.browserSessionsActive == nil {
		return
	}
	mp.browserSessionsActive.Add(ctx, -1)
}

// Web3 Metrics Methods

// RecordWeb3Transaction records a Web3 transaction metric
func (mp *MetricsProvider) RecordWeb3Transaction(ctx context.Context, chain, txType, status string) {
	if mp.web3TransactionsTotal == nil {
		return
	}

	attrs := []attribute.KeyValue{
		attribute.String("chain", chain),
		attribute.String("type", txType),
		attribute.String("status", status),
	}

	mp.web3TransactionsTotal.Add(ctx, 1, metric.WithAttributes(attrs...))
}

// System Metrics Methods

// UpdateErrorRate updates the current error rate
func (mp *MetricsProvider) UpdateErrorRate(ctx context.Context, rate float64) {
	if mp.errorRate == nil {
		return
	}
	mp.errorRate.Record(ctx, rate)
}

// UpdateSystemResourceUsage updates system resource usage
func (mp *MetricsProvider) UpdateSystemResourceUsage(ctx context.Context, resourceType string, usage float64) {
	if mp.systemResourceUsage == nil {
		return
	}

	attrs := []attribute.KeyValue{
		attribute.String("resource", resourceType),
	}

	mp.systemResourceUsage.Record(ctx, usage, metric.WithAttributes(attrs...))
}

// StartMetricsServer starts the Prometheus metrics HTTP server
func (mp *MetricsProvider) StartMetricsServer(port int) error {
	if mp.registry == nil {
		return fmt.Errorf("metrics not enabled")
	}

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(mp.registry, promhttp.HandlerOpts{
		EnableOpenMetrics: true,
	}))

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	return server.ListenAndServe()
}

// Shutdown gracefully shuts down the metrics provider
func (mp *MetricsProvider) Shutdown(ctx context.Context) error {
	if mp.meterProvider == nil {
		return nil
	}
	return mp.meterProvider.Shutdown(ctx)
}
