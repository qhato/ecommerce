package http

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/qhato/ecommerce/internal/payment/application/commands"
	"github.com/qhato/ecommerce/internal/payment/application/queries"
)

type PaymentTokenHandler struct {
	commandHandler *commands.PaymentTokenCommandHandler
	queryService   *queries.PaymentTokenQueryService
}

func NewPaymentTokenHandler(
	commandHandler *commands.PaymentTokenCommandHandler,
	queryService *queries.PaymentTokenQueryService,
) *PaymentTokenHandler {
	return &PaymentTokenHandler{
		commandHandler: commandHandler,
		queryService:   queryService,
	}
}

func (h *PaymentTokenHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/payment-tokens", h.CreatePaymentToken).Methods("POST")
	router.HandleFunc("/payment-tokens/{id}", h.GetPaymentToken).Methods("GET")
	router.HandleFunc("/payment-tokens/{id}", h.DeletePaymentToken).Methods("DELETE")
	router.HandleFunc("/payment-tokens/{id}/set-default", h.SetDefaultToken).Methods("POST")
	router.HandleFunc("/payment-tokens/{id}/deactivate", h.DeactivateToken).Methods("POST")
	router.HandleFunc("/customers/{customerId}/payment-tokens", h.GetCustomerTokens).Methods("GET")
	router.HandleFunc("/customers/{customerId}/payment-tokens/active", h.GetCustomerActiveTokens).Methods("GET")
	router.HandleFunc("/customers/{customerId}/payment-tokens/default", h.GetDefaultToken).Methods("GET")
}

type CreatePaymentTokenRequest struct {
	CustomerID  string  `json:"customer_id" validate:"required"`
	Token       string  `json:"token" validate:"required"`
	GatewayName string  `json:"gateway_name" validate:"required"`
	TokenType   string  `json:"token_type" validate:"required"`
	Last4Digits *string `json:"last_4_digits"`
	CardBrand   *string `json:"card_brand"`
	ExpiryMonth *int    `json:"expiry_month"`
	ExpiryYear  *int    `json:"expiry_year"`
	IsDefault   bool    `json:"is_default"`
}

type SetDefaultTokenRequest struct {
	CustomerID string `json:"customer_id" validate:"required"`
}

func (h *PaymentTokenHandler) CreatePaymentToken(w http.ResponseWriter, r *http.Request) {
	var req CreatePaymentTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cmd := commands.CreatePaymentTokenCommand{
		CustomerID:  req.CustomerID,
		Token:       req.Token,
		GatewayName: req.GatewayName,
		TokenType:   req.TokenType,
		Last4Digits: req.Last4Digits,
		CardBrand:   req.CardBrand,
		ExpiryMonth: req.ExpiryMonth,
		ExpiryYear:  req.ExpiryYear,
		IsDefault:   req.IsDefault,
	}

	token, err := h.commandHandler.HandleCreatePaymentToken(r.Context(), cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(token)
}

func (h *PaymentTokenHandler) GetPaymentToken(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	token, err := h.queryService.GetToken(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(token)
}

func (h *PaymentTokenHandler) DeletePaymentToken(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	cmd := commands.DeleteTokenCommand{TokenID: id}
	if err := h.commandHandler.HandleDeleteToken(r.Context(), cmd); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *PaymentTokenHandler) SetDefaultToken(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req SetDefaultTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cmd := commands.SetDefaultTokenCommand{
		TokenID:    id,
		CustomerID: req.CustomerID,
	}

	token, err := h.commandHandler.HandleSetDefaultToken(r.Context(), cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(token)
}

func (h *PaymentTokenHandler) DeactivateToken(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	cmd := commands.DeactivateTokenCommand{TokenID: id}
	if err := h.commandHandler.HandleDeactivateToken(r.Context(), cmd); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *PaymentTokenHandler) GetCustomerTokens(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customerID := vars["customerId"]

	tokens, err := h.queryService.GetCustomerTokens(r.Context(), customerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokens)
}

func (h *PaymentTokenHandler) GetCustomerActiveTokens(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customerID := vars["customerId"]

	tokens, err := h.queryService.GetCustomerActiveTokens(r.Context(), customerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokens)
}

func (h *PaymentTokenHandler) GetDefaultToken(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customerID := vars["customerId"]

	token, err := h.queryService.GetDefaultToken(r.Context(), customerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(token)
}
