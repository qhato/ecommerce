package http

import (
	"encoding/json"
	"net/http"

	"github.com/qhato/ecommerce/pkg/logger"
)

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

// Meta represents pagination and other metadata
type Meta struct {
	Page       int `json:"page,omitempty"`
	PerPage    int `json:"per_page,omitempty"`
	Total      int `json:"total,omitempty"`
	TotalPages int `json:"total_pages,omitempty"`
}

// WriteJSON writes a JSON response
func WriteJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger.WithError(err).Error("Failed to encode JSON response")
	}
}

// WriteSuccess writes a successful JSON response
func WriteSuccess(w http.ResponseWriter, data interface{}) {
	WriteJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    data,
	})
}

// WriteCreated writes a created response (201)
func WriteCreated(w http.ResponseWriter, data interface{}) {
	WriteJSON(w, http.StatusCreated, Response{
		Success: true,
		Data:    data,
		Message: "Resource created successfully",
	})
}

// WriteNoContent writes a no content response (204)
func WriteNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

// WritePaginated writes a paginated response
func WritePaginated(w http.ResponseWriter, data interface{}, meta Meta) {
	WriteJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    data,
		Meta:    &meta,
	})
}

// WriteMessage writes a response with just a message
func WriteMessage(w http.ResponseWriter, statusCode int, message string) {
	WriteJSON(w, statusCode, Response{
		Success: statusCode >= 200 && statusCode < 300,
		Message: message,
	})
}

// RespondJSON responds with JSON data
func RespondJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger.WithError(err).Error("Failed to encode JSON response")
	}
}

// RespondError responds with an error
func RespondError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	// Determine status code based on error type
	statusCode := http.StatusInternalServerError
	if appErr, ok := err.(interface{ StatusCode() int }); ok {
		statusCode = appErr.StatusCode()
	}

	w.WriteHeader(statusCode)
	response := map[string]interface{}{
		"error": err.Error(),
	}

	if encodeErr := json.NewEncoder(w).Encode(response); encodeErr != nil {
		logger.WithError(encodeErr).Error("Failed to encode error response")
	}
}

// ValidationError represents a validation error
type ValidationError struct {
	Message string
}

func (e ValidationError) Error() string {
	return e.Message
}

func (e ValidationError) StatusCode() int {
	return http.StatusBadRequest
}

// NewValidationError creates a new validation error
func NewValidationError(message string) error {
	return ValidationError{Message: message}
}

// NotFoundError represents a not found error
type NotFoundError struct {
	Message string
}

func (e NotFoundError) Error() string {
	return e.Message
}

func (e NotFoundError) StatusCode() int {
	return http.StatusNotFound
}

// NewNotFoundError creates a new not found error
func NewNotFoundError(message string) error {
	return NotFoundError{Message: message}
}
