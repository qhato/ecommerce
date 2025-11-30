package application

import (
	"time"

	"github.com/qhato/ecommerce/internal/customer/domain"
)

// CustomerDTO represents a customer data transfer object
type CustomerDTO struct {
	ID           int64     `json:"id"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	EmailAddress string    `json:"email_address"`
	UserName     string    `json:"user_name"`
	ReceiveEmail bool      `json:"receive_email"`
	Deactivated  bool      `json:"deactivated"`
	Archived     bool      `json:"archived"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ToCustomerDTO converts a domain Customer to CustomerDTO
func ToCustomerDTO(c *domain.Customer) *CustomerDTO {
	return &CustomerDTO{
		ID:           c.ID,
		FirstName:    c.FirstName,
		LastName:     c.LastName,
		EmailAddress: c.EmailAddress,
		UserName:     c.UserName,
		ReceiveEmail: c.ReceiveEmail,
		Deactivated:  c.Deactivated,
		Archived:     c.Archived,
		CreatedAt:    c.CreatedAt,
		UpdatedAt:    c.UpdatedAt,
	}
}

// PaginatedResponse represents a paginated response (reusing structure if not imported)
// Ideally this should be shared, but defining here for independence or using the one from catalog if imported.
// customer_queries.go was trying to use application.PaginatedResponse.
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalItems int64       `json:"total_items"`
	TotalPages int64       `json:"total_pages"`
}

// NewPaginatedResponse creates a new paginated response
func NewPaginatedResponse(data interface{}, page, pageSize int, totalItems int64) *PaginatedResponse {
	totalPages := totalItems / int64(pageSize)
	if totalItems%int64(pageSize) > 0 {
		totalPages++
	}

	return &PaginatedResponse{
		Data:       data,
		Page:       page,
		PageSize:   pageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}
}
