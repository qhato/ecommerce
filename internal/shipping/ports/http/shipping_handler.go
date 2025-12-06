package http

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/qhato/ecommerce/internal/shipping/application"
	"github.com/shopspring/decimal"
)

type ShippingHandler struct {
	service *application.ShippingService
}

func NewShippingHandler(service *application.ShippingService) *ShippingHandler {
	return &ShippingHandler{service: service}
}

func (h *ShippingHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/shipping/rates", h.CalculateRates).Methods("POST")
}

type CalculateRatesRequest struct {
	Weight     decimal.Decimal `json:"weight" validate:"required"`
	OrderTotal decimal.Decimal `json:"order_total" validate:"required"`
	Quantity   int             `json:"quantity" validate:"required,min=1"`
	Country    string          `json:"country" validate:"required"`
	ZipCode    string          `json:"zip_code"`
}

func (h *ShippingHandler) CalculateRates(w http.ResponseWriter, r *http.Request) {
	var req CalculateRatesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	svcReq := application.ShippingRateRequest{
		Weight:     req.Weight,
		OrderTotal: req.OrderTotal,
		Quantity:   req.Quantity,
		Country:    req.Country,
		ZipCode:    req.ZipCode,
	}

	rates, err := h.service.CalculateRates(r.Context(), svcReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rates)
}
