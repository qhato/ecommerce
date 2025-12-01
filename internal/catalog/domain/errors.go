package domain

// DomainError represents a business rule validation error
type DomainError struct {
	Message string
}

func (e *DomainError) Error() string {
	return e.Message
}

// NewDomainError creates a new DomainError
func NewDomainError(message string) error {
	return &DomainError{Message: message}
}
