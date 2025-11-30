package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/qhato/ecommerce/internal/customer/application/commands"
	"github.com/qhato/ecommerce/internal/customer/application/queries"
	httpPkg "github.com/qhato/ecommerce/pkg/http"
	"github.com/qhato/ecommerce/pkg/logger"
	"github.com/qhato/ecommerce/pkg/validator"
	"github.com/qhato/ecommerce/pkg/errors" // Import pkg/errors
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
	var cmd commands.RegisterCustomerCommand // Use commands.RegisterCustomerCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		httpPkg.RespondError(w, errors.BadRequest("invalid request body").WithInternal(err))
		return
	}

	if err := h.validator.Validate(cmd); err != nil {
		httpPkg.RespondError(w, errors.ValidationError("validation failed").WithInternal(err))
		return
	}

	customerID, err := h.commandHandler.HandleRegisterCustomer(r.Context(), &cmd) // Call HandleRegisterCustomer
	if err != nil {
		// Differentiate between user-facing errors (e.g., conflict) and internal errors
		if errors.IsConflict(err) {
			httpPkg.RespondError(w, errors.Conflict(err.Error()))
		} else {
			httpPkg.RespondError(w, errors.Internal("failed to register customer").WithInternal(err))
		}
		return
	}

	httpPkg.RespondJSON(w, http.StatusCreated, map[string]interface{}{
		"id": customerID,
	})
}

// GetCustomer retrieves a customer by ID
func (h *AdminCustomerHandler) GetCustomer(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpPkg.RespondError(w, errors.BadRequest("invalid customer ID").WithInternal(err))
		return
	}

	query := &queries.GetCustomerByIDQuery{ID: id} // Use query struct
	customer, err := h.queryHandler.HandleGetCustomerByID(r.Context(), query) // Call HandleGetCustomerByID
	if err != nil {
		if errors.IsNotFound(err) {
			httpPkg.RespondError(w, errors.NotFound(err.Error()))
		} else {
			httpPkg.RespondError(w, errors.Internal("failed to get customer").WithInternal(err))
		}
		return
	}

	httpPkg.RespondJSON(w, http.StatusOK, customer) // Removed redundant application.ToCustomerDTO(customer)
}

// GetCustomerByEmail retrieves a customer by email
func (h *AdminCustomerHandler) GetCustomerByEmail(w http.ResponseWriter, r *http.Request) {
	email := chi.URLParam(r, "email")
	if email == "" {
		httpPkg.RespondError(w, errors.BadRequest("email is required"))
		return
	}

	query := &queries.GetCustomerByEmailQuery{Email: email} // Use query struct
	customer, err := h.queryHandler.HandleGetCustomerByEmail(r.Context(), query) // Call HandleGetCustomerByEmail
	if err != nil {
		if errors.IsNotFound(err) {
			httpPkg.RespondError(w, errors.NotFound(err.Error()))
		} else {
			httpPkg.RespondError(w, errors.Internal("failed to get customer by email").WithInternal(err))
		}
		return
	}

	httpPkg.RespondJSON(w, http.StatusOK, customer) // Removed redundant application.ToCustomerDTO(customer)
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
	includeArchived := r.URL.Query().Get("include_archived") == "true"
	activeOnly := r.URL.Query().Get("active_only") == "true"
	registeredOnly := r.URL.Query().Get("registered_only") == "true"
	searchQuery := r.URL.Query().Get("q")


	query := &queries.ListCustomersQuery{ // Use query struct
		Page:            page,
		PageSize:        pageSize,
		SortBy:          sortBy,
		SortOrder:       sortOrder,
		IncludeArchived: includeArchived,
		ActiveOnly:      activeOnly,
		RegisteredOnly:  registeredOnly,
		SearchQuery:     searchQuery,
	}

	result, err := h.queryHandler.HandleListCustomers(r.Context(), query) // Call HandleListCustomers
	if err != nil {
		httpPkg.RespondError(w, errors.Internal("failed to list customers").WithInternal(err))
		return
	}
	// The application.PaginatedResponse should handle the pagination details now.
	// No need to manually calculate totalPages and create PaginatedCustomerResponse.

	httpPkg.RespondJSON(w, http.StatusOK, result)
}

// UpdateCustomer updates a customer's profile
func (h *AdminCustomerHandler) UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpPkg.RespondError(w, errors.BadRequest("invalid customer ID").WithInternal(err))
		return
	}

	var cmd commands.UpdateCustomerCommand // Use commands.UpdateCustomerCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		httpPkg.RespondError(w, errors.BadRequest("invalid request body").WithInternal(err))
		return
	}

	if err := h.validator.Validate(cmd); err != nil {
		httpPkg.RespondError(w, errors.ValidationError("validation failed").WithInternal(err))
		return
	}
	cmd.ID = id // Set ID from URL param

	err = h.commandHandler.HandleUpdateCustomer(r.Context(), &cmd) // Call HandleUpdateCustomer
	if err != nil {
		if errors.IsNotFound(err) {
			httpPkg.RespondError(w, errors.NotFound(err.Error()))
		} else if errors.IsConflict(err) {
			httpPkg.RespondError(w, errors.Conflict(err.Error()))
		} else {
			httpPkg.RespondError(w, errors.Internal("failed to update customer").WithInternal(err))
		}
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
		httpPkg.RespondError(w, errors.BadRequest("invalid customer ID").WithInternal(err))
		return
	}

	var cmd commands.ChangePasswordCommand // Use commands.ChangePasswordCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		httpPkg.RespondError(w, errors.BadRequest("invalid request body").WithInternal(err))
		return
	}

	if err := h.validator.Validate(cmd); err != nil {
		httpPkg.RespondError(w, errors.ValidationError("validation failed").WithInternal(err))
		return
	}
	cmd.CustomerID = id // Set customer ID

	err = h.commandHandler.HandleChangePassword(r.Context(), &cmd) // Call HandleChangePassword
	if err != nil {
		if errors.IsUnauthorized(err) {
			httpPkg.RespondError(w, errors.Unauthorized(err.Error()))
		} else if errors.IsNotFound(err) {
			httpPkg.RespondError(w, errors.NotFound(err.Error()))
		} else {
			httpPkg.RespondError(w, errors.Internal("failed to change password").WithInternal(err))
		}
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
		httpPkg.RespondError(w, errors.BadRequest("invalid customer ID").WithInternal(err))
		return
	}

	cmd := &commands.DeactivateCustomerCommand{ID: id} // Use commands.DeactivateCustomerCommand
	err = h.commandHandler.HandleDeactivateCustomer(r.Context(), cmd) // Call HandleDeactivateCustomer
	if err != nil {
		if errors.IsNotFound(err) {
			httpPkg.RespondError(w, errors.NotFound(err.Error()))
		} else if errors.IsConflict(err) {
			httpPkg.RespondError(w, errors.Conflict(err.Error()))
		} else {
			httpPkg.RespondError(w, errors.Internal("failed to deactivate customer").WithInternal(err))
		}
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
		httpPkg.RespondError(w, errors.BadRequest("invalid customer ID").WithInternal(err))
		return
	}

	cmd := &commands.ActivateCustomerCommand{ID: id} // Use commands.ActivateCustomerCommand
	err = h.commandHandler.HandleActivateCustomer(r.Context(), cmd) // Call HandleActivateCustomer
	if err != nil {
		if errors.IsNotFound(err) {
			httpPkg.RespondError(w, errors.NotFound(err.Error()))
		} else if errors.IsConflict(err) {
			httpPkg.RespondError(w, errors.Conflict(err.Error()))
		} else {
			httpPkg.RespondError(w, errors.Internal("failed to activate customer").WithInternal(err))
		}
		return
	}

	// Invalidate cache
	h.queryHandler.InvalidateCache(r.Context(), id)

	httpPkg.RespondJSON(w, http.StatusOK, map[string]string{"message": "customer activated successfully"})
}
