package middleware

import (
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// Tracing is a middleware that adds OpenTelemetry tracing to HTTP requests
func Tracing(serviceName string) func(next http.Handler) http.Handler {
	tracer := otel.Tracer(serviceName)
	propagator := otel.GetTextMapPropagator()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract trace context from incoming request
			ctx := propagator.Extract(r.Context(), propagation.HeaderCarrier(r.Header))

			// Start span
			spanName := r.Method + " " + r.URL.Path
			ctx, span := tracer.Start(
				ctx,
				spanName,
				trace.WithSpanKind(trace.SpanKindServer),
				trace.WithAttributes(
					attribute.String("http.method", r.Method),
					attribute.String("http.target", r.URL.Path),
					attribute.String("http.route", r.URL.Path),
					attribute.String("http.scheme", r.URL.Scheme),
					attribute.String("http.host", r.Host),
					attribute.String("http.user_agent", r.UserAgent()),
					attribute.String("net.peer.ip", r.RemoteAddr),
				),
			)
			defer span.End()

			// Wrap response writer to capture status code
			wrapped := &statusCodeCapture{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			// Process request with trace context
			next.ServeHTTP(wrapped, r.WithContext(ctx))

			// Set status code attribute
			span.SetAttributes(attribute.Int("http.status_code", wrapped.statusCode))

			// Mark span as error if status code is 5xx
			if wrapped.statusCode >= 500 {
				span.SetAttributes(attribute.Bool("error", true))
			}
		})
	}
}

// statusCodeCapture wraps ResponseWriter to capture status code
type statusCodeCapture struct {
	http.ResponseWriter
	statusCode int
}

func (s *statusCodeCapture) WriteHeader(code int) {
	s.statusCode = code
	s.ResponseWriter.WriteHeader(code)
}
