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

// Common errors
var (
	ErrInvalidBundleItem      = NewDomainError("bundle item must have either product or SKU")
	ErrInvalidQuantity        = NewDomainError("quantity must be greater than zero")
	ErrInvalidPrice           = NewDomainError("price must be non-negative")
	ErrBundleNotFound         = NewDomainError("product bundle not found")
	ErrRelationshipNotFound   = NewDomainError("product relationship not found")
	ErrDuplicateRelationship  = NewDomainError("relationship already exists")
	ErrSelfRelationship       = NewDomainError("product cannot be related to itself")
)
