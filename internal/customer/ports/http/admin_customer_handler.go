package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/qhato/ecommerce/internal/customer/application"
	"github.com/qhato/ecommerce/internal/customer/application/commands"
	"github.com/qhato/ecommerce/internal/customer/application/queries"
	httpPkg "github.com/qhato/ecommerce/pkg/http"
	"github.com/qhato/ecommerce/pkg/logger"
	"github.com/qhato/ecommerce/pkg/validator"
)

// AdminCustomerHandler handles admin customer HTTP requests
type AdminCustomerHandler struct {
	commandHandler *commands.CustomerCommandHandler
	queryHandler   *queries.CustomerQueryHandler
	validator      *validator.Validator
	log            *logger.Logger
}

// NewAdminCustomerHandler creates a new AdminCustomerHandler
func NewAdminCustomerHandler(
	commandHandler *commands.CustomerCommandHandler,
	queryHandler *queries.CustomerQueryHandler,
	validator *validator.Validator,
	log *logger.Logger,
) *AdminCustomerHandler {
	return &AdminCustomerHandler{
		commandHandler: commandHandler,
		queryHandler:   queryHandler,
		validator:      validator,
		log:            log,
	}
}

// RegisterRoutes registers customer routes
func (h *AdminCustomerHandler) RegisterRoutes(r chi.Router) {
	r.Route("/customers", func(r chi.Router) {
		r.Post("/", h.RegisterCustomer)
		r.Get("/", h.ListCustomers)
		r.Get("/{id}", h.GetCustomer)
		r.Put("/{id}", h.UpdateCustomer)
		r.Put("/{id}/password", h.ChangePassword)
		r.Post("/{id}/deactivate", h.DeactivateCustomer)
		r.Post("/{id}/activate", h.ActivateCustomer)
		r.Get("/email/{email}", h.GetCustomerByEmail)
	})
}

// RegisterCustomer registers a new customer
func (h *AdminCustomerHandler) RegisterCustomer(w http.ResponseWriter, r *http.Request) {
	var req application.RegisterCustomerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpPkg.RespondError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	if err := h.validator.Validate(req); err != nil {
		httpPkg.RespondError(w, http.StatusBadRequest, "validation failed", err)
		return
	}

	customer, err := h.commandHandler.RegisterCustomer(
		r.Context(),
		req.EmailAddress,
		req.UserName,
		req.Password,
		req.FirstName,
		req.LastName,
	)
	if err != nil {
		httpPkg.RespondError(w, http.StatusInternalServerError, "failed to register customer", err)
		return
	}

	httpPkg.RespondJSON(w, http.StatusCreated, application.ToCustomerDTO(customer))
}

// GetCustomer retrieves a customer by ID
func (h *AdminCustomerHandler) GetCustomer(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpPkg.RespondError(w, http.StatusBadRequest, "invalid customer ID", err)
		return
	}

	customer, err := h.queryHandler.GetByID(r.Context(), id)
	if err != nil {
		httpPkg.RespondError(w, http.StatusNotFound, "customer not found", err)
		return
	}

	httpPkg.RespondJSON(w, http.StatusOK, application.ToCustomerDTO(customer))
}

// GetCustomerByEmail retrieves a customer by email
func (h *AdminCustomerHandler) GetCustomerByEmail(w http.ResponseWriter, r *http.Request) {
	email := chi.URLParam(r, "email")
	if email == "" {
		httpPkg.RespondError(w, http.StatusBadRequest, "email is required", nil)
		return
	}

	customer, err := h.queryHandler.GetByEmail(r.Context(), email)
	if err != nil {
		httpPkg.RespondError(w, http.StatusNotFound, "customer not found", err)
		return
	}

	httpPkg.RespondJSON(w, http.StatusOK, application.ToCustomerDTO(customer))
}

// ListCustomers lists all customers
func (h *AdminCustomerHandler) ListCustomers(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	sortBy := r.URL.Query().Get("sort_by")
	sortOrder := r.URL.Query().Get("sort_order")

	filter := &application.CustomerFilter{
		Page:      page,
		PageSize:  pageSize,
		SortBy:    sortBy,
		SortOrder: sortOrder,
	}

	customers, total, err := h.queryHandler.List(r.Context(), filter)
	if err != nil {
		httpPkg.RespondError(w, http.StatusInternalServerError, "failed to list customers", err)
		return
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	response := application.PaginatedCustomerResponse{
		Data:       application.ToCustomerDTOs(customers),
		Page:       page,
		PageSize:   pageSize,
		TotalItems: total,
		TotalPages: totalPages,
	}

	httpPkg.RespondJSON(w, http.StatusOK, response)
}

// UpdateCustomer updates a customer's profile
func (h *AdminCustomerHandler) UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpPkg.RespondError(w, http.StatusBadRequest, "invalid customer ID", err)
		return
	}

	var req application.UpdateCustomerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpPkg.RespondError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	if err := h.validator.Validate(req); err != nil {
		httpPkg.RespondError(w, http.StatusBadRequest, "validation failed", err)
		return
	}

	err = h.commandHandler.UpdateCustomer(r.Context(), id, req.FirstName, req.LastName, req.EmailAddress)
	if err != nil {
		httpPkg.RespondError(w, http.StatusInternalServerError, "failed to update customer", err)
		return
	}

	// Invalidate cache
	h.queryHandler.InvalidateCache(r.Context(), id)

	httpPkg.RespondJSON(w, http.StatusOK, map[string]string{"message": "customer updated successfully"})
}

// ChangePassword changes a customer's password
func (h *AdminCustomerHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpPkg.RespondError(w, http.StatusBadRequest, "invalid customer ID", err)
		return
	}

	var req application.ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpPkg.RespondError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	if err := h.validator.Validate(req); err != nil {
		httpPkg.RespondError(w, http.StatusBadRequest, "validation failed", err)
		return
	}

	err = h.commandHandler.ChangePassword(r.Context(), id, req.NewPassword)
	if err != nil {
		httpPkg.RespondError(w, http.StatusInternalServerError, "failed to change password", err)
		return
	}

	// Invalidate cache
	h.queryHandler.InvalidateCache(r.Context(), id)

	httpPkg.RespondJSON(w, http.StatusOK, map[string]string{"message": "password changed successfully"})
}

// DeactivateCustomer deactivates a customer account
func (h *AdminCustomerHandler) DeactivateCustomer(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpPkg.RespondError(w, http.StatusBadRequest, "invalid customer ID", err)
		return
	}

	err = h.commandHandler.DeactivateCustomer(r.Context(), id)
	if err != nil {
		httpPkg.RespondError(w, http.StatusInternalServerError, "failed to deactivate customer", err)
		return
	}

	// Invalidate cache
	h.queryHandler.InvalidateCache(r.Context(), id)

	httpPkg.RespondJSON(w, http.StatusOK, map[string]string{"message": "customer deactivated successfully"})
}

// ActivateCustomer activates a customer account
func (h *AdminCustomerHandler) ActivateCustomer(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpPkg.RespondError(w, http.StatusBadRequest, "invalid customer ID", err)
		return
	}

	err = h.commandHandler.ActivateCustomer(r.Context(), id)
	if err != nil {
		httpPkg.RespondError(w, http.StatusInternalServerError, "failed to activate customer", err)
		return
	}

	// Invalidate cache
	h.queryHandler.InvalidateCache(r.Context(), id)

	httpPkg.RespondJSON(w, http.StatusOK, map[string]string{"message": "customer activated successfully"})
}
