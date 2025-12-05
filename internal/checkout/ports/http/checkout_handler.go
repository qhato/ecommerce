package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/qhato/ecommerce/internal/checkout/application/commands"
	"github.com/qhato/ecommerce/internal/checkout/application/queries"
	"github.com/qhato/ecommerce/internal/checkout/domain"
)

// CheckoutHandler handles HTTP requests for checkout
type CheckoutHandler struct {
	commandHandler *commands.CheckoutCommandHandler
	queryService   *queries.CheckoutQueryService
}

// NewCheckoutHandler creates a new checkout HTTP handler
func NewCheckoutHandler(
	commandHandler *commands.CheckoutCommandHandler,
	queryService *queries.CheckoutQueryService,
) *CheckoutHandler {
	return &CheckoutHandler{
		commandHandler: commandHandler,
		queryService:   queryService,
	}
}

// RegisterRoutes registers all checkout routes
func (h *CheckoutHandler) RegisterRoutes(router *mux.Router) {
	// Checkout Session Endpoints
	router.HandleFunc("/checkout/initiate", h.InitiateCheckout).Methods("POST")
	router.HandleFunc("/checkout/{sessionId}", h.GetCheckoutSession).Methods("GET")
	router.HandleFunc("/checkout/order/{orderId}", h.GetCheckoutByOrderID).Methods("GET")

	// Checkout Steps
	router.HandleFunc("/checkout/{sessionId}/customer-info", h.AddCustomerInfo).Methods("POST")
	router.HandleFunc("/checkout/{sessionId}/shipping-address", h.AddShippingAddress).Methods("POST")
	router.HandleFunc("/checkout/{sessionId}/shipping-method", h.SelectShippingMethod).Methods("POST")
	router.HandleFunc("/checkout/{sessionId}/billing-address", h.AddBillingAddress).Methods("POST")
	router.HandleFunc("/checkout/{sessionId}/payment-method", h.AddPaymentMethod).Methods("POST")

	// Coupons
	router.HandleFunc("/checkout/{sessionId}/coupons", h.ApplyCoupon).Methods("POST")
	router.HandleFunc("/checkout/{sessionId}/coupons/{code}", h.RemoveCoupon).Methods("DELETE")

	// Submission
	router.HandleFunc("/checkout/{sessionId}/submit", h.SubmitCheckout).Methods("POST")
	router.HandleFunc("/checkout/{sessionId}/confirm", h.ConfirmCheckout).Methods("POST")
	router.HandleFunc("/checkout/{sessionId}/cancel", h.CancelCheckout).Methods("POST")

	// Session Management
	router.HandleFunc("/checkout/{sessionId}/extend", h.ExtendSession).Methods("POST")

	// Shipping Options
	router.HandleFunc("/checkout/shipping-options", h.GetAllShippingOptions).Methods("GET")
	router.HandleFunc("/checkout/shipping-options/available", h.GetAvailableShippingOptions).Methods("GET")
}

// InitiateCheckout handles initiating a new checkout session
func (h *CheckoutHandler) InitiateCheckout(w http.ResponseWriter, r *http.Request) {
	var cmd commands.InitiateCheckoutCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	session, err := h.commandHandler.HandleInitiateCheckout(r.Context(), cmd)
	if err != nil {
		if err == domain.ErrCheckoutSessionAlreadyExists {
			respondWithError(w, http.StatusConflict, err.Error())
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to initiate checkout: "+err.Error())
		return
	}

	response := queries.ToCheckoutSessionDTO(session)
	respondWithJSON(w, http.StatusCreated, response)
}

// GetCheckoutSession handles retrieving a checkout session
func (h *CheckoutHandler) GetCheckoutSession(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := vars["sessionId"]

	session, err := h.queryService.GetCheckoutSession(r.Context(), sessionID)
	if err != nil {
		if err == domain.ErrCheckoutSessionNotFound {
			respondWithError(w, http.StatusNotFound, err.Error())
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to get checkout session: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, session)
}

// GetCheckoutByOrderID handles retrieving checkout by order ID
func (h *CheckoutHandler) GetCheckoutByOrderID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID, err := strconv.ParseInt(vars["orderId"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid order ID")
		return
	}

	session, err := h.queryService.GetCheckoutByOrderID(r.Context(), orderID)
	if err != nil {
		if err == domain.ErrCheckoutSessionNotFound {
			respondWithError(w, http.StatusNotFound, err.Error())
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to get checkout session: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, session)
}

// AddCustomerInfo handles adding customer information
func (h *CheckoutHandler) AddCustomerInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := vars["sessionId"]

	var cmd commands.AddCustomerInfoCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	cmd.SessionID = sessionID

	session, err := h.commandHandler.HandleAddCustomerInfo(r.Context(), cmd)
	if err != nil {
		h.handleCheckoutError(w, err)
		return
	}

	response := queries.ToCheckoutSessionDTO(session)
	respondWithJSON(w, http.StatusOK, response)
}

// AddShippingAddress handles adding shipping address
func (h *CheckoutHandler) AddShippingAddress(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := vars["sessionId"]

	var cmd commands.AddShippingAddressCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	cmd.SessionID = sessionID

	session, err := h.commandHandler.HandleAddShippingAddress(r.Context(), cmd)
	if err != nil {
		h.handleCheckoutError(w, err)
		return
	}

	response := queries.ToCheckoutSessionDTO(session)
	respondWithJSON(w, http.StatusOK, response)
}

// SelectShippingMethod handles selecting shipping method
func (h *CheckoutHandler) SelectShippingMethod(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := vars["sessionId"]

	var cmd commands.SelectShippingMethodCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	cmd.SessionID = sessionID

	session, err := h.commandHandler.HandleSelectShippingMethod(r.Context(), cmd)
	if err != nil {
		h.handleCheckoutError(w, err)
		return
	}

	response := queries.ToCheckoutSessionDTO(session)
	respondWithJSON(w, http.StatusOK, response)
}

// AddBillingAddress handles adding billing address
func (h *CheckoutHandler) AddBillingAddress(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := vars["sessionId"]

	var cmd commands.AddBillingAddressCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	cmd.SessionID = sessionID

	session, err := h.commandHandler.HandleAddBillingAddress(r.Context(), cmd)
	if err != nil {
		h.handleCheckoutError(w, err)
		return
	}

	response := queries.ToCheckoutSessionDTO(session)
	respondWithJSON(w, http.StatusOK, response)
}

// AddPaymentMethod handles adding payment method
func (h *CheckoutHandler) AddPaymentMethod(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := vars["sessionId"]

	var cmd commands.AddPaymentMethodCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	cmd.SessionID = sessionID

	session, err := h.commandHandler.HandleAddPaymentMethod(r.Context(), cmd)
	if err != nil {
		h.handleCheckoutError(w, err)
		return
	}

	response := queries.ToCheckoutSessionDTO(session)
	respondWithJSON(w, http.StatusOK, response)
}

// ApplyCoupon handles applying a coupon code
func (h *CheckoutHandler) ApplyCoupon(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := vars["sessionId"]

	var cmd commands.ApplyCouponCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	cmd.SessionID = sessionID

	session, err := h.commandHandler.HandleApplyCoupon(r.Context(), cmd)
	if err != nil {
		h.handleCheckoutError(w, err)
		return
	}

	response := queries.ToCheckoutSessionDTO(session)
	respondWithJSON(w, http.StatusOK, response)
}

// RemoveCoupon handles removing a coupon code
func (h *CheckoutHandler) RemoveCoupon(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := vars["sessionId"]
	code := vars["code"]

	cmd := commands.RemoveCouponCommand{
		SessionID:  sessionID,
		CouponCode: code,
	}

	session, err := h.commandHandler.HandleRemoveCoupon(r.Context(), cmd)
	if err != nil {
		h.handleCheckoutError(w, err)
		return
	}

	response := queries.ToCheckoutSessionDTO(session)
	respondWithJSON(w, http.StatusOK, response)
}

// SubmitCheckout handles submitting the checkout
func (h *CheckoutHandler) SubmitCheckout(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := vars["sessionId"]

	cmd := commands.SubmitCheckoutCommand{
		SessionID: sessionID,
	}

	session, err := h.commandHandler.HandleSubmitCheckout(r.Context(), cmd)
	if err != nil {
		h.handleCheckoutError(w, err)
		return
	}

	response := queries.ToCheckoutSessionDTO(session)
	respondWithJSON(w, http.StatusOK, response)
}

// ConfirmCheckout handles confirming checkout after payment
func (h *CheckoutHandler) ConfirmCheckout(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := vars["sessionId"]

	var cmd commands.ConfirmCheckoutCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	cmd.SessionID = sessionID

	session, err := h.commandHandler.HandleConfirmCheckout(r.Context(), cmd)
	if err != nil {
		h.handleCheckoutError(w, err)
		return
	}

	response := queries.ToCheckoutSessionDTO(session)
	respondWithJSON(w, http.StatusOK, response)
}

// CancelCheckout handles cancelling a checkout
func (h *CheckoutHandler) CancelCheckout(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := vars["sessionId"]

	var req struct {
		Reason string `json:"reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	cmd := commands.CancelCheckoutCommand{
		SessionID: sessionID,
		Reason:    req.Reason,
	}

	if err := h.commandHandler.HandleCancelCheckout(r.Context(), cmd); err != nil {
		h.handleCheckoutError(w, err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}

// ExtendSession handles extending session expiration
func (h *CheckoutHandler) ExtendSession(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := vars["sessionId"]

	var req struct {
		Hours int `json:"hours"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	cmd := commands.ExtendSessionCommand{
		SessionID: sessionID,
		Hours:     req.Hours,
	}

	session, err := h.commandHandler.HandleExtendSession(r.Context(), cmd)
	if err != nil {
		h.handleCheckoutError(w, err)
		return
	}

	response := queries.ToCheckoutSessionDTO(session)
	respondWithJSON(w, http.StatusOK, response)
}

// GetAllShippingOptions handles retrieving all shipping options
func (h *CheckoutHandler) GetAllShippingOptions(w http.ResponseWriter, r *http.Request) {
	activeOnly := r.URL.Query().Get("active") == "true"

	options, err := h.queryService.GetAllShippingOptions(r.Context(), activeOnly)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get shipping options: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, options)
}

// GetAvailableShippingOptions handles retrieving available shipping options for an address
func (h *CheckoutHandler) GetAvailableShippingOptions(w http.ResponseWriter, r *http.Request) {
	country := r.URL.Query().Get("country")
	stateProvince := r.URL.Query().Get("stateProvince")
	postalCode := r.URL.Query().Get("postalCode")

	if country == "" {
		respondWithError(w, http.StatusBadRequest, "country parameter is required")
		return
	}

	options, err := h.queryService.GetAvailableShippingOptions(r.Context(), country, stateProvince, postalCode)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get shipping options: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, options)
}

// Helper methods

func (h *CheckoutHandler) handleCheckoutError(w http.ResponseWriter, err error) {
	switch err {
	case domain.ErrCheckoutSessionNotFound:
		respondWithError(w, http.StatusNotFound, err.Error())
	case domain.ErrCheckoutSessionExpired:
		respondWithError(w, http.StatusGone, err.Error())
	case domain.ErrCheckoutNotReady:
		respondWithError(w, http.StatusPreconditionFailed, err.Error())
	case domain.ErrCheckoutNotSubmitted:
		respondWithError(w, http.StatusPreconditionFailed, err.Error())
	case domain.ErrShippingMethodNotFound:
		respondWithError(w, http.StatusNotFound, err.Error())
	case domain.ErrShippingMethodUnavailable:
		respondWithError(w, http.StatusConflict, err.Error())
	default:
		respondWithError(w, http.StatusInternalServerError, "Checkout operation failed: "+err.Error())
	}
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if payload != nil {
		json.NewEncoder(w).Encode(payload)
	}
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}
