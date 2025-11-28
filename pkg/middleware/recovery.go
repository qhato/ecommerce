package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/qhato/ecommerce/pkg/errors"
	"github.com/qhato/ecommerce/pkg/logger"
)

// Recovery recovers from panics and returns 500 error
func Recovery() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					// Log the panic with stack trace
					logger.WithFields(logger.Fields{
						"panic":       rec,
						"stack_trace": string(debug.Stack()),
						"method":      r.Method,
						"path":        r.URL.Path,
					}).Error("Panic recovered in HTTP handler")

					// Return error response
					err := errors.Internal(fmt.Sprintf("Internal server error: %v", rec))
					errors.HandleHTTPError(w, err)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
