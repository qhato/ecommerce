package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/qhato/ecommerce/pkg/metrics"
)

// metricsResponseWriter wraps http.ResponseWriter to capture status code and size
type metricsResponseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int64
}

func (rw *metricsResponseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *metricsResponseWriter) Write(b []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(b)
	rw.size += int64(size)
	return size, err
}

// Metrics is a middleware that records Prometheus metrics for HTTP requests
func Metrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap response writer to capture status code and size
		wrapped := &metricsResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
			size:           0,
		}

		// Get request size
		requestSize := r.ContentLength
		if requestSize < 0 {
			requestSize = 0
		}

		// Process request
		next.ServeHTTP(wrapped, r)

		// Record metrics
		duration := time.Since(start)
		status := strconv.Itoa(wrapped.statusCode)

		metrics.RecordHTTPRequest(
			r.Method,
			r.URL.Path,
			status,
			duration,
			requestSize,
			wrapped.size,
		)

		// Record errors (4xx and 5xx)
		if wrapped.statusCode >= 400 {
			errorType := "client_error"
			if wrapped.statusCode >= 500 {
				errorType = "server_error"
			}
			metrics.RecordHTTPError(r.Method, r.URL.Path, errorType)
		}
	})
}
