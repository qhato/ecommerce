package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/qhato/ecommerce/pkg/errors"
)

// Validator wraps go-playground/validator
type Validator struct {
	validate *validator.Validate
}

// New creates a new validator instance
func New() *Validator {
	validate := validator.New()

	// Register custom validation tags here
	// Example: validate.RegisterValidation("custom_tag", customValidationFunc)

	return &Validator{
		validate: validate,
	}
}

// Validate validates a struct and returns AppError if validation fails
func (v *Validator) Validate(data interface{}) error {
	if err := v.validate.Struct(data); err != nil {
		return v.formatValidationErrors(err)
	}
	return nil
}

// ValidateVar validates a single variable
func (v *Validator) ValidateVar(field interface{}, tag string) error {
	if err := v.validate.Var(field, tag); err != nil {
		return v.formatValidationErrors(err)
	}
	return nil
}

// formatValidationErrors converts validator errors to AppError
func (v *Validator) formatValidationErrors(err error) error {
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return errors.ValidationError(err.Error())
	}

	var errorMessages []string
	appErr := errors.ValidationError("Validation failed")

	for _, fieldErr := range validationErrors {
		message := v.getErrorMessage(fieldErr)
		errorMessages = append(errorMessages, message)

		// Add field-specific details
		appErr = appErr.WithDetail(
			fieldErr.Field(),
			map[string]interface{}{
				"tag":     fieldErr.Tag(),
				"value":   fieldErr.Value(),
				"message": message,
			},
		)
	}

	// Set overall message
	appErr.Message = strings.Join(errorMessages, "; ")

	return appErr
}

// getErrorMessage generates a human-readable error message for a field error
func (v *Validator) getErrorMessage(fieldErr validator.FieldError) string {
	field := fieldErr.Field()
	tag := fieldErr.Tag()
	param := fieldErr.Param()

	switch tag {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", field, param)
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", field, param)
	case "len":
		return fmt.Sprintf("%s must be exactly %s characters", field, param)
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", field, param)
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", field, param)
	case "lt":
		return fmt.Sprintf("%s must be less than %s", field, param)
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", field, param)
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, param)
	case "url":
		return fmt.Sprintf("%s must be a valid URL", field)
	case "uuid":
		return fmt.Sprintf("%s must be a valid UUID", field)
	case "numeric":
		return fmt.Sprintf("%s must be numeric", field)
	case "alpha":
		return fmt.Sprintf("%s must contain only letters", field)
	case "alphanum":
		return fmt.Sprintf("%s must contain only letters and numbers", field)
	default:
		return fmt.Sprintf("%s failed validation: %s", field, tag)
	}
}

// Global validator instance
var defaultValidator *Validator

// init initializes the default validator
func init() {
	defaultValidator = New()
}

// Validate validates using the default validator
func Validate(data interface{}) error {
	return defaultValidator.Validate(data)
}

// ValidateVar validates a variable using the default validator
func ValidateVar(field interface{}, tag string) error {
	return defaultValidator.ValidateVar(field, tag)
}
