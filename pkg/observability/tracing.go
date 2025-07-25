package observability

import (
	"context"
	"fmt"
	"log"

	"github.com/ai-agentic-browser/internal/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// TracingProvider manages OpenTelemetry tracing
type TracingProvider struct {
	provider *trace.TracerProvider
	tracer   oteltrace.Tracer
}

// NewTracingProvider creates a new tracing provider
func NewTracingProvider(cfg config.ObservabilityConfig) (*TracingProvider, error) {
	// Create Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(cfg.JaegerEndpoint)))
	if err != nil {
		return nil, fmt.Errorf("failed to create Jaeger exporter: %w", err)
	}

	// Create resource with service information
	res, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(cfg.ServiceName),
			semconv.ServiceVersionKey.String("1.0.0"),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Create trace provider
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(res),
		trace.WithSampler(trace.AlwaysSample()),
	)

	// Set global trace provider
	otel.SetTracerProvider(tp)

	// Set global propagator
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	tracer := tp.Tracer(cfg.ServiceName)

	return &TracingProvider{
		provider: tp,
		tracer:   tracer,
	}, nil
}

// Tracer returns the OpenTelemetry tracer
func (tp *TracingProvider) Tracer() oteltrace.Tracer {
	return tp.tracer
}

// Shutdown gracefully shuts down the tracing provider
func (tp *TracingProvider) Shutdown(ctx context.Context) error {
	return tp.provider.Shutdown(ctx)
}

// StartSpan starts a new span with the given name
func (tp *TracingProvider) StartSpan(ctx context.Context, name string, opts ...oteltrace.SpanStartOption) (context.Context, oteltrace.Span) {
	return tp.tracer.Start(ctx, name, opts...)
}

// SpanFromContext returns the span from the context
func SpanFromContext(ctx context.Context) oteltrace.Span {
	return oteltrace.SpanFromContext(ctx)
}

// AddSpanAttributes adds attributes to the span in the context
func AddSpanAttributes(ctx context.Context, attrs ...oteltrace.SpanStartOption) {
	span := SpanFromContext(ctx)
	if span.IsRecording() {
		// Note: This is a simplified version. In practice, you'd need to extract
		// attributes from SpanStartOptions or use a different approach
		log.Printf("Adding attributes to span: %v", attrs)
	}
}

// RecordError records an error in the span
func RecordError(ctx context.Context, err error) {
	span := SpanFromContext(ctx)
	if span.IsRecording() {
		span.RecordError(err)
	}
}

// SetSpanStatus sets the status of the span
func SetSpanStatus(ctx context.Context, code codes.Code, description string) {
	span := SpanFromContext(ctx)
	if span.IsRecording() {
		span.SetStatus(code, description)
	}
}
