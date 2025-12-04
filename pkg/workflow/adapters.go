package workflow

import (
	"context"
	"time"

	"github.com/qhato/ecommerce/pkg/logging"
	"github.com/qhato/ecommerce/pkg/metrics"
	"github.com/qhato/ecommerce/pkg/tracing"
)

// LoggerAdapter adapts our logging.Logger to workflow.Logger
type LoggerAdapter struct {
	logger logging.Logger
}

// NewLoggerAdapter creates a new logger adapter
func NewLoggerAdapter(logger logging.Logger) *LoggerAdapter {
	return &LoggerAdapter{logger: logger}
}

func (l *LoggerAdapter) Debug(msg string, fields ...interface{}) {
	l.logger.Debug(msg, convertFields(fields...)...)
}

func (l *LoggerAdapter) Info(msg string, fields ...interface{}) {
	l.logger.Info(msg, convertFields(fields...)...)
}

func (l *LoggerAdapter) Warn(msg string, fields ...interface{}) {
	l.logger.Warn(msg, convertFields(fields...)...)
}

func (l *LoggerAdapter) Error(msg string, fields ...interface{}) {
	l.logger.Error(msg, convertFields(fields...)...)
}

func (l *LoggerAdapter) WithContext(ctx context.Context) Logger {
	return &LoggerAdapter{logger: l.logger.WithContext(ctx)}
}

// convertFields converts key-value pairs to logging fields
func convertFields(fields ...interface{}) []logging.Field {
	result := make([]logging.Field, 0, len(fields)/2)
	for i := 0; i < len(fields)-1; i += 2 {
		key, ok := fields[i].(string)
		if !ok {
			continue
		}
		value := fields[i+1]
		result = append(result, logging.Any(key, value))
	}
	return result
}

// MetricsAdapter adapts Prometheus metrics to workflow.MetricsRecorder
type MetricsAdapter struct {
	namespace string
}

// NewMetricsAdapter creates a new metrics adapter
func NewMetricsAdapter(namespace string) *MetricsAdapter {
	return &MetricsAdapter{namespace: namespace}
}

func (m *MetricsAdapter) RecordWorkflowExecution(workflowName string, duration time.Duration, status Status) {
	// Record workflow duration
	if metrics.HTTP != nil {
		// Use existing metrics infrastructure
		// In production, you would create workflow-specific metrics
		metrics.RecordHTTPRequest("WORKFLOW", workflowName, string(status), duration, 0, 0)
	}
}

func (m *MetricsAdapter) RecordActivityExecution(workflowName, activityName string, duration time.Duration, status Status) {
	// Record activity duration
	if metrics.Database != nil {
		// Use database metrics as placeholder
		metrics.RecordDatabaseQuery(workflowName, activityName, duration)
	}
}

func (m *MetricsAdapter) IncrementWorkflowCounter(workflowName string, status Status) {
	// Increment workflow counter based on status
	switch status {
	case StatusCompleted:
		// Success counter
	case StatusFailed:
		// Failure counter
	case StatusCompensated:
		// Compensation counter
	}
}

func (m *MetricsAdapter) IncrementActivityCounter(workflowName, activityName string, status Status) {
	// Increment activity counter
}

// TracerAdapter adapts OpenTelemetry tracer to workflow.Tracer
type TracerAdapter struct{}

// NewTracerAdapter creates a new tracer adapter
func NewTracerAdapter() *TracerAdapter {
	return &TracerAdapter{}
}

func (t *TracerAdapter) StartSpan(ctx context.Context, name string) (context.Context, Span) {
	ctx, span := tracing.StartSpan(ctx, name)
	return ctx, &SpanAdapter{span: span}
}

// SpanAdapter adapts OpenTelemetry span to workflow.Span
type SpanAdapter struct {
	span tracing.Span
}

func (s *SpanAdapter) End() {
	s.span.End()
}

func (s *SpanAdapter) SetAttribute(key string, value interface{}) {
	// Convert to OpenTelemetry attribute
	switch v := value.(type) {
	case string:
		s.span.SetAttributes(tracing.Attribute.String(key, v))
	case int:
		s.span.SetAttributes(tracing.Attribute.Int(key, v))
	case int64:
		s.span.SetAttributes(tracing.Attribute.Int64(key, v))
	case bool:
		s.span.SetAttributes(tracing.Attribute.Bool(key, v))
	default:
		// Use Any for other types
		s.span.SetAttributes(tracing.Attribute.Any(key, v))
	}
}

func (s *SpanAdapter) RecordError(err error) {
	s.span.RecordError(err)
}

// Note: This is a placeholder implementation
// In the actual tracing package, you would have these methods
type AttributeHelper struct{}

var Attribute AttributeHelper

func (AttributeHelper) String(key, value string) interface{} {
	return struct{ key, value string }{key, value}
}

func (AttributeHelper) Int(key string, value int) interface{} {
	return struct {
		key   string
		value int
	}{key, value}
}

func (AttributeHelper) Int64(key string, value int64) interface{} {
	return struct {
		key   string
		value int64
	}{key, value}
}

func (AttributeHelper) Bool(key string, value bool) interface{} {
	return struct {
		key   string
		value bool
	}{key, value}
}

func (AttributeHelper) Any(key string, value interface{}) interface{} {
	return struct {
		key   string
		value interface{}
	}{key, value}
}
