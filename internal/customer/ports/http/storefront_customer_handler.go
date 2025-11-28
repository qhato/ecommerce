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

	// Don't return password in response
	dto := application.ToCustomerDTO(customer)

	httpPkg.RespondJSON(w, http.StatusCreated, dto)
}

// GetProfile retrieves a customer's profile
func (h *StorefrontCustomerHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpPkg.RespondError(w, http.StatusBadRequest, "invalid customer ID", err)
		return
	}

	// TODO: In production, verify that the authenticated user matches this ID
	// or has appropriate permissions

	customer, err := h.queryHandler.GetByID(r.Context(), id)
	if err != nil {
		httpPkg.RespondError(w, http.StatusNotFound, "customer not found", err)
		return
	}

	httpPkg.RespondJSON(w, http.StatusOK, application.ToCustomerDTO(customer))
}

// UpdateProfile updates a customer's profile
func (h *StorefrontCustomerHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpPkg.RespondError(w, http.StatusBadRequest, "invalid customer ID", err)
		return
	}

	// TODO: In production, verify that the authenticated user matches this ID

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
		httpPkg.RespondError(w, http.StatusInternalServerError, "failed to update profile", err)
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
		httpPkg.RespondError(w, http.StatusBadRequest, "invalid customer ID", err)
		return
	}

	// TODO: In production, verify that the authenticated user matches this ID
	// and validate old password before changing

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
