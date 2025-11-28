package queries

import (
	"context"

	"github.com/qhato/ecommerce/internal/customer/application"
	"github.com/qhato/ecommerce/internal/customer/domain"
	"github.com/qhato/ecommerce/pkg/cache"
	"github.com/qhato/ecommerce/pkg/errors"
	"github.com/qhato/ecommerce/pkg/logger"
)

// GetCustomerByIDQuery represents a query to get a customer by ID
type GetCustomerByIDQuery struct {
	ID int64 `json:"id" validate:"required"`
}

// GetCustomerByEmailQuery represents a query to get a customer by email
type GetCustomerByEmailQuery struct {
	Email string `json:"email" validate:"required,email"`
}

// ListCustomersQuery represents a query to list customers
type ListCustomersQuery struct {
	Page            int    `json:"page" validate:"min=1"`
	PageSize        int    `json:"page_size" validate:"min=1,max=100"`
	IncludeArchived bool   `json:"include_archived"`
	ActiveOnly      bool   `json:"active_only"`
	RegisteredOnly  bool   `json:"registered_only"`
	SortBy          string `json:"sort_by"`
	SortOrder       string `json:"sort_order"`
	SearchQuery     string `json:"search_query"`
}

// CustomerQueryHandler handles customer queries
type CustomerQueryHandler struct {
	repo   domain.CustomerRepository
	cache  cache.Cache
	logger *logger.Logger
}

// NewCustomerQueryHandler creates a new customer query handler
func NewCustomerQueryHandler(
	repo domain.CustomerRepository,
	cache cache.Cache,
	logger *logger.Logger,
) *CustomerQueryHandler {
	return &CustomerQueryHandler{
		repo:   repo,
		cache:  cache,
		logger: logger,
	}
}

// HandleGetCustomerByID handles the get customer by ID query
func (h *CustomerQueryHandler) HandleGetCustomerByID(ctx context.Context, query *GetCustomerByIDQuery) (*application.CustomerDTO, error) {
	// Try to get from cache first
	cacheKey := customerCacheKey(query.ID)
	var customer *domain.Customer

	if err := h.cache.Get(ctx, cacheKey, &customer); err == nil && customer != nil {
		h.logger.Debug("customer found in cache", "customer_id", query.ID)
		return application.ToCustomerDTO(customer), nil
	}

	// Get from repository
	customer, err := h.repo.FindByID(ctx, query.ID)
	if err != nil {
		return nil, errors.Wrap(err, "customer not found")
	}

	// Cache the result
	if err := h.cache.Set(ctx, cacheKey, customer, cache.DefaultTTL); err != nil {
		h.logger.Warn("failed to cache customer", "error", err, "customer_id", query.ID)
	}

	return application.ToCustomerDTO(customer), nil
}

// HandleGetCustomerByEmail handles the get customer by email query
func (h *CustomerQueryHandler) HandleGetCustomerByEmail(ctx context.Context, query *GetCustomerByEmailQuery) (*application.CustomerDTO, error) {
	customer, err := h.repo.FindByEmail(ctx, query.Email)
	if err != nil {
		return nil, errors.Wrap(err, "customer not found")
	}

	// Cache the result
	cacheKey := customerCacheKey(customer.ID)
	if err := h.cache.Set(ctx, cacheKey, customer, cache.DefaultTTL); err != nil {
		h.logger.Warn("failed to cache customer", "error", err, "customer_id", customer.ID)
	}

	return application.ToCustomerDTO(customer), nil
}

// HandleListCustomers handles the list customers query
func (h *CustomerQueryHandler) HandleListCustomers(ctx context.Context, query *ListCustomersQuery) (*application.PaginatedResponse, error) {
	// Set defaults
	if query.Page < 1 {
		query.Page = 1
	}
	if query.PageSize < 1 {
		query.PageSize = 20
	}
	if query.SortBy == "" {
		query.SortBy = "created_at"
	}
	if query.SortOrder == "" {
		query.SortOrder = "desc"
	}

	// Create filter
	filter := &domain.CustomerFilter{
		Page:            query.Page,
		PageSize:        query.PageSize,
		IncludeArchived: query.IncludeArchived,
		ActiveOnly:      query.ActiveOnly,
		RegisteredOnly:  query.RegisteredOnly,
		SortBy:          query.SortBy,
		SortOrder:       query.SortOrder,
		SearchQuery:     query.SearchQuery,
	}

	// Get from repository
	customers, total, err := h.repo.FindAll(ctx, filter)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list customers")
	}

	// Convert to DTOs
	customerDTOs := make([]*application.CustomerDTO, len(customers))
	for i, customer := range customers {
		customerDTOs[i] = application.ToCustomerDTO(customer)
	}

	return application.NewPaginatedResponse(customerDTOs, query.Page, query.PageSize, total), nil
}

// customerCacheKey generates a cache key for a customer
func customerCacheKey(id int64) string {
	return cache.Key("customer", "customer", id)
}

// PaginatedResponse represents a paginated response (reusing from catalog)
type PaginatedResponse = application.PaginatedResponse

// NewPaginatedResponse creates a new paginated response
var NewPaginatedResponse = application.NewPaginatedResponse
