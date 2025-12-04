package tracing

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Config contains tracing configuration
type Config struct {
	ServiceName    string
	ServiceVersion string
	Environment    string
	ExporterType   string // "jaeger", "otlp", or "noop"
	JaegerEndpoint string // e.g., "http://localhost:14268/api/traces"
	OTLPEndpoint   string // e.g., "localhost:4317"
	SamplingRate   float64
}

// Provider wraps the tracer provider
type Provider struct {
	tp     *sdktrace.TracerProvider
	tracer trace.Tracer
}

// Init initializes the tracing provider
func Init(cfg Config) (*Provider, error) {
	// Create resource
	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(cfg.ServiceName),
			semconv.ServiceVersion(cfg.ServiceVersion),
			attribute.String("environment", cfg.Environment),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Create exporter based on type
	var exporter sdktrace.SpanExporter
	switch cfg.ExporterType {
	case "jaeger":
		exporter, err = jaeger.New(
			jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(cfg.JaegerEndpoint)),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create Jaeger exporter: %w", err)
		}

	case "otlp":
		ctx := context.Background()
		conn, err := grpc.DialContext(ctx, cfg.OTLPEndpoint,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithBlock(),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
		}

		exporter, err = otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
		if err != nil {
			return nil, fmt.Errorf("failed to create OTLP exporter: %w", err)
		}

	case "noop":
		// No exporter - useful for testing or disabled tracing
		exporter = nil

	default:
		return nil, fmt.Errorf("unknown exporter type: %s", cfg.ExporterType)
	}

	// Create trace provider options
	opts := []sdktrace.TracerProviderOption{
		sdktrace.WithResource(res),
	}

	// Add exporter if not noop
	if exporter != nil {
		opts = append(opts, sdktrace.WithBatcher(exporter))
	}

	// Add sampler
	if cfg.SamplingRate > 0 && cfg.SamplingRate <= 1.0 {
		opts = append(opts, sdktrace.WithSampler(sdktrace.TraceIDRatioBased(cfg.SamplingRate)))
	} else {
		// Default to always sample
		opts = append(opts, sdktrace.WithSampler(sdktrace.AlwaysSample()))
	}

	// Create tracer provider
	tp := sdktrace.NewTracerProvider(opts...)

	// Set global tracer provider
	otel.SetTracerProvider(tp)

	// Set global propagator (for distributed tracing)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	return &Provider{
		tp:     tp,
		tracer: tp.Tracer(cfg.ServiceName),
	}, nil
}

// Tracer returns the tracer instance
func (p *Provider) Tracer() trace.Tracer {
	return p.tracer
}

// Shutdown gracefully shuts down the tracer provider
func (p *Provider) Shutdown(ctx context.Context) error {
	return p.tp.Shutdown(ctx)
}

// StartSpan starts a new span with the given name and options
func StartSpan(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return otel.Tracer("ecommerce").Start(ctx, spanName, opts...)
}

// SpanFromContext returns the current span from context
func SpanFromContext(ctx context.Context) trace.Span {
	return trace.SpanFromContext(ctx)
}

// AddEvent adds an event to the current span
func AddEvent(ctx context.Context, name string, attrs ...attribute.KeyValue) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent(name, trace.WithAttributes(attrs...))
}

// SetAttributes sets attributes on the current span
func SetAttributes(ctx context.Context, attrs ...attribute.KeyValue) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attrs...)
}

// RecordError records an error on the current span
func RecordError(ctx context.Context, err error, opts ...trace.EventOption) {
	span := trace.SpanFromContext(ctx)
	span.RecordError(err, opts...)
}

// Common attribute keys
var (
	AttrUserID       = attribute.Key("user.id")
	AttrCustomerID   = attribute.Key("customer.id")
	AttrOrderID      = attribute.Key("order.id")
	AttrProductID    = attribute.Key("product.id")
	AttrPaymentID    = attribute.Key("payment.id")
	AttrShipmentID   = attribute.Key("shipment.id")
	AttrDBOperation  = attribute.Key("db.operation")
	AttrDBTable      = attribute.Key("db.table")
	AttrCacheKey     = attribute.Key("cache.key")
	AttrCacheHit     = attribute.Key("cache.hit")
	AttrEventType    = attribute.Key("event.type")
)
