package errors

import (
	"encoding/json"
	"net/http"

	"github.com/qhato/ecommerce/pkg/logger"
)

// ErrorResponse represents the JSON error response sent to clients
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail contains the error details
type ErrorDetail struct {
	Code    ErrorCode              `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// HandleHTTPError handles errors in HTTP handlers
func HandleHTTPError(w http.ResponseWriter, err error) {
	var appErr *AppError
	var statusCode int
	var errorResponse ErrorResponse

	// Check if it's an AppError
	if As(err, &appErr) {
		statusCode = appErr.StatusCode
		errorResponse = ErrorResponse{
			Error: ErrorDetail{
				Code:    appErr.Code,
				Message: appErr.Message,
				Details: appErr.Details,
			},
		}

		// Log internal error if present
		if appErr.Internal != nil {
			logger.WithError(appErr.Internal).WithFields(logger.Fields{
				"code":        appErr.Code,
				"status_code": appErr.StatusCode,
			}).Error("Internal error occurred")
		}
	} else {
		// Unknown error - treat as internal server error
		statusCode = http.StatusInternalServerError
		errorResponse = ErrorResponse{
			Error: ErrorDetail{
				Code:    ErrCodeInternal,
				Message: "An internal error occurred",
			},
		}

		// Log the unknown error
		logger.WithError(err).Error("Unexpected error occurred")
	}

	// Set response headers and write JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
		logger.WithError(err).Error("Failed to encode error response")
	}
}

// RecoveryHandler creates a middleware for panic recovery
func RecoveryHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				logger.WithField("panic", rec).Error("Panic recovered")

				err := Internal("An unexpected error occurred")
				HandleHTTPError(w, err)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
