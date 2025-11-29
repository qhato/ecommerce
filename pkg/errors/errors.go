package errors

import (
	"errors"
	"fmt"
	"net/http"
)

// ErrorCode represents a unique error code
type ErrorCode string

// Common error codes
const (
	// Client errors (4xx)
	ErrCodeBadRequest          ErrorCode = "BAD_REQUEST"
	ErrCodeUnauthorized        ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden           ErrorCode = "FORBIDDEN"
	ErrCodeNotFound            ErrorCode = "NOT_FOUND"
	ErrCodeConflict            ErrorCode = "CONFLICT"
	ErrCodeValidation          ErrorCode = "VALIDATION_ERROR"
	ErrCodeUnprocessableEntity ErrorCode = "UNPROCESSABLE_ENTITY"
	ErrCodeTooManyRequests     ErrorCode = "TOO_MANY_REQUESTS"

	// Server errors (5xx)
	ErrCodeInternal       ErrorCode = "INTERNAL_ERROR"
	ErrCodeNotImplemented ErrorCode = "NOT_IMPLEMENTED"
	ErrCodeServiceUnavail ErrorCode = "SERVICE_UNAVAILABLE"
	ErrCodeGatewayTimeout ErrorCode = "GATEWAY_TIMEOUT"

	// Business logic errors
	ErrCodeInsufficientStock ErrorCode = "INSUFFICIENT_STOCK"
	ErrCodeInvalidCoupon     ErrorCode = "INVALID_COUPON"
	ErrCodePaymentFailed     ErrorCode = "PAYMENT_FAILED"
	ErrCodeProductInactive   ErrorCode = "PRODUCT_INACTIVE"
	ErrCodeOrderNotEditable  ErrorCode = "ORDER_NOT_EDITABLE"
)

// AppError represents an application error with additional context
type AppError struct {
	Code       ErrorCode              `json:"code"`
	Message    string                 `json:"message"`
	StatusCode int                    `json:"-"`
	Details    map[string]interface{} `json:"details,omitempty"`
	Internal   error                  `json:"-"` // Internal error (not exposed to client)
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Internal != nil {
		return fmt.Sprintf("%s: %s (internal: %v)", e.Code, e.Message, e.Internal)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap implements the errors.Unwrap interface
func (e *AppError) Unwrap() error {
	return e.Internal
}

// WithDetail adds a detail field to the error
func (e *AppError) WithDetail(key string, value interface{}) *AppError {
	if e.Details == nil {
		e.Details = make(map[string]interface{})
	}
	e.Details[key] = value
	return e
}

// WithInternal adds an internal error
func (e *AppError) WithInternal(err error) *AppError {
	e.Internal = err
	return e
}

// New creates a new AppError
func New(code ErrorCode, message string, statusCode int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
	}
}

// Wrap wraps an existing error into an AppError
func Wrap(err error, code ErrorCode, message string, statusCode int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
		Internal:   err,
	}
}

// Is checks if the error is of a specific type
func Is(err error, target error) bool {
	return errors.Is(err, target)
}

// As finds the first error in err's chain that matches target
func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

// GetStatusCode returns the HTTP status code for an error
func GetStatusCode(err error) int {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.StatusCode
	}
	return http.StatusInternalServerError
}

// Predefined error constructors for common cases

// BadRequest creates a bad request error (400)
func BadRequest(message string) *AppError {
	return New(ErrCodeBadRequest, message, http.StatusBadRequest)
}

// Unauthorized creates an unauthorized error (401)
func Unauthorized(message string) *AppError {
	return New(ErrCodeUnauthorized, message, http.StatusUnauthorized)
}

// Forbidden creates a forbidden error (403)
func Forbidden(message string) *AppError {
	return New(ErrCodeForbidden, message, http.StatusForbidden)
}

// NotFound creates a not found error (404)
func NotFound(resource string) *AppError {
	return New(ErrCodeNotFound, fmt.Sprintf("%s not found", resource), http.StatusNotFound)
}

// Conflict creates a conflict error (409)
func Conflict(message string) *AppError {
	return New(ErrCodeConflict, message, http.StatusConflict)
}

// ValidationError creates a validation error (422)
func ValidationError(message string) *AppError {
	return New(ErrCodeValidation, message, http.StatusUnprocessableEntity)
}

// Internal creates an internal server error (500)
func Internal(message string) *AppError {
	return New(ErrCodeInternal, message, http.StatusInternalServerError)
}

// InternalWrap wraps an error as internal server error
func InternalWrap(err error, message string) *AppError {
	return Wrap(err, ErrCodeInternal, message, http.StatusInternalServerError)
}

// NotImplemented creates a not implemented error (501)
func NotImplemented(message string) *AppError {
	return New(ErrCodeNotImplemented, message, http.StatusNotImplemented)
}

// ServiceUnavailable creates a service unavailable error (503)
func ServiceUnavailable(message string) *AppError {
	return New(ErrCodeServiceUnavail, message, http.StatusServiceUnavailable)
}

// Business logic error constructors

// InsufficientStock creates an insufficient stock error
func InsufficientStock(productID string, requested, available int) *AppError {
	return New(
		ErrCodeInsufficientStock,
		fmt.Sprintf("Insufficient stock for product %s", productID),
		http.StatusConflict,
	).WithDetail("product_id", productID).
		WithDetail("requested", requested).
		WithDetail("available", available)
}

// InvalidCoupon creates an invalid coupon error
func InvalidCoupon(couponCode string) *AppError {
	return New(
		ErrCodeInvalidCoupon,
		"Invalid or expired coupon code",
		http.StatusUnprocessableEntity,
	).WithDetail("coupon_code", couponCode)
}

// PaymentFailed creates a payment failed error
func PaymentFailed(reason string) *AppError {
	return New(
		ErrCodePaymentFailed,
		"Payment processing failed",
		http.StatusPaymentRequired,
	).WithDetail("reason", reason)
}

// ProductInactive creates a product inactive error
func ProductInactive(productID string) *AppError {
	return New(
		ErrCodeProductInactive,
		"Product is not available for purchase",
		http.StatusConflict,
	).WithDetail("product_id", productID)
}

// OrderNotEditable creates an order not editable error
func OrderNotEditable(orderID string, status string) *AppError {
	return New(
		ErrCodeOrderNotEditable,
		"Order cannot be modified in current status",
		http.StatusConflict,
	).WithDetail("order_id", orderID).
		WithDetail("status", status)
}
