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

// StorefrontCustomerHandler handles storefront customer HTTP requests
type StorefrontCustomerHandler struct {
	commandHandler *commands.CustomerCommandHandler
	queryHandler   *queries.CustomerQueryHandler
	validator      *validator.Validator
	log            *logger.Logger
}

// NewStorefrontCustomerHandler creates a new StorefrontCustomerHandler
func NewStorefrontCustomerHandler(
	commandHandler *commands.CustomerCommandHandler,
	queryHandler *queries.CustomerQueryHandler,
	validator *validator.Validator,
	log *logger.Logger,
) *StorefrontCustomerHandler {
	return &StorefrontCustomerHandler{
		commandHandler: commandHandler,
		queryHandler:   queryHandler,
		validator:      validator,
		log:            log,
	}
}

// RegisterRoutes registers customer routes
func (h *StorefrontCustomerHandler) RegisterRoutes(r chi.Router) {
	r.Route("/customers", func(r chi.Router) {
		r.Post("/register", h.RegisterCustomer)
		r.Get("/{id}/profile", h.GetProfile)
		r.Put("/{id}/profile", h.UpdateProfile)
		r.Put("/{id}/password", h.ChangePassword)
	})
}

// RegisterCustomer registers a new customer
func (h *StorefrontCustomerHandler) RegisterCustomer(w http.ResponseWriter, r *http.Request) {
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

// GetProfile retrieves a customer's profile
func (h *StorefrontCustomerHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpPkg.RespondError(w, errors.BadRequest("invalid customer ID").WithInternal(err))
		return
	}

	// TODO: In production, verify that the authenticated user matches this ID
	// or has appropriate permissions

	query := &queries.GetCustomerByIDQuery{ID: id}
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

// UpdateProfile updates a customer's profile
func (h *StorefrontCustomerHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpPkg.RespondError(w, errors.BadRequest("invalid customer ID").WithInternal(err))
		return
	}

	// TODO: In production, verify that the authenticated user matches this ID

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

	httpPkg.RespondJSON(w, http.StatusOK, map[string]string{"message": "profile updated successfully"})
}

// ChangePassword changes a customer's password
func (h *StorefrontCustomerHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpPkg.RespondError(w, errors.BadRequest("invalid customer ID").WithInternal(err))
		return
	}

	// TODO: In production, verify that the authenticated user matches this ID
	// and validate old password before changing

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
