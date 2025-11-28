package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/qhato/ecommerce/pkg/logger"
)

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	bytes      int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.bytes += n
	return n, err
}

// RequestLogger logs HTTP requests
func RequestLogger() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Generate correlation ID
			correlationID := r.Header.Get("X-Correlation-ID")
			if correlationID == "" {
				correlationID = uuid.New().String()
			}

			// Add correlation ID to context
			ctx := context.WithValue(r.Context(), "correlation_id", correlationID)
			r = r.WithContext(ctx)

			// Add correlation ID to response header
			w.Header().Set("X-Correlation-ID", correlationID)

			// Wrap response writer
			wrapped := newResponseWriter(w)

			// Log request
			logger.WithFields(logger.Fields{
				"correlation_id": correlationID,
				"method":         r.Method,
				"path":           r.URL.Path,
				"query":          r.URL.RawQuery,
				"remote_addr":    r.RemoteAddr,
				"user_agent":     r.UserAgent(),
			}).Info("HTTP request started")

			// Process request
			next.ServeHTTP(wrapped, r)

			// Log response
			duration := time.Since(start)
			logger.WithFields(logger.Fields{
				"correlation_id": correlationID,
				"method":         r.Method,
				"path":           r.URL.Path,
				"status":         wrapped.statusCode,
				"duration_ms":    duration.Milliseconds(),
				"bytes":          wrapped.bytes,
			}).Info("HTTP request completed")
		})
	}
}

// GetCorrelationID extracts correlation ID from context
func GetCorrelationID(ctx context.Context) string {
	correlationID, _ := ctx.Value("correlation_id").(string)
	return correlationID
}
