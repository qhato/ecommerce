package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/qhato/ecommerce/internal/offer/application"
	"github.com/qhato/ecommerce/internal/offer/domain"
	"github.com/shopspring/decimal"
)

// OfferHandler handles HTTP requests for offer operations
type OfferHandler struct {
	offerService          application.OfferService
	offerProcessorService application.OfferProcessorService
}

// NewOfferHandler creates a new OfferHandler
func NewOfferHandler(
	offerService application.OfferService,
	offerProcessorService application.OfferProcessorService,
) *OfferHandler {
	return &OfferHandler{
		offerService:          offerService,
		offerProcessorService: offerProcessorService,
	}
}

// RegisterRoutes registers all offer-related routes
func (h *OfferHandler) RegisterRoutes(router *mux.Router) {
	// Offer CRUD
	router.HandleFunc("/api/admin/offers", h.CreateOffer).Methods("POST")
	router.HandleFunc("/api/admin/offers/{id}", h.GetOffer).Methods("GET")
	router.HandleFunc("/api/admin/offers/{id}", h.UpdateOffer).Methods("PUT")
	router.HandleFunc("/api/admin/offers/{id}", h.DeleteOffer).Methods("DELETE")
	router.HandleFunc("/api/admin/offers", h.GetActiveOffers).Methods("GET")

	// Offer Codes
	router.HandleFunc("/api/admin/offers/{id}/codes", h.CreateOfferCode).Methods("POST")
	router.HandleFunc("/api/admin/offer-codes/{id}", h.GetOfferCode).Methods("GET")
	router.HandleFunc("/api/admin/offer-codes/{id}", h.UpdateOfferCode).Methods("PUT")
	router.HandleFunc("/api/admin/offer-codes/{id}", h.DeleteOfferCode).Methods("DELETE")

	// Offer Processing
	router.HandleFunc("/api/storefront/orders/{orderId}/process-offers", h.ProcessOrderOffers).Methods("POST")
	router.HandleFunc("/api/storefront/orders/{orderId}/apply-code", h.ApplyOfferCode).Methods("POST")
	router.HandleFunc("/api/storefront/orders/{orderId}/offers/{offerId}", h.RemoveOffer).Methods("DELETE")

	// Lookup
	router.HandleFunc("/api/storefront/offers/by-code/{code}", h.GetOfferByCode).Methods("GET")
}

// CreateOffer handles POST /api/admin/offers
func (h *OfferHandler) CreateOffer(w http.ResponseWriter, r *http.Request) {
	var req CreateOfferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	cmd := &application.CreateOfferCommand{
		Name:                      req.Name,
		OfferType:                 domain.OfferType(req.OfferType),
		OfferValue:                req.OfferValue,
		AdjustmentType:            domain.OfferAdjustmentType(req.AdjustmentType),
		ApplyToChildItems:         req.ApplyToChildItems,
		ApplyToSalePrice:          req.ApplyToSalePrice,
		AutomaticallyAdded:        req.AutomaticallyAdded,
		CombinableWithOtherOffers: req.CombinableWithOtherOffers,
		OfferDescription:          req.OfferDescription,
		OfferDiscountType:         domain.OfferDiscountType(req.OfferDiscountType),
		EndDate:                   req.EndDate,
		MarketingMessage:          req.MarketingMessage,
		MaxUsesPerCustomer:        req.MaxUsesPerCustomer,
		MaxUses:                   req.MaxUses,
		OrderMinTotal:             req.OrderMinTotal,
		OfferPriority:             req.OfferPriority,
		StartDate:                 req.StartDate,
	}

	offer, err := h.offerService.CreateOffer(r.Context(), cmd)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, offer)
}

// GetOffer handles GET /api/admin/offers/{id}
func (h *OfferHandler) GetOffer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid offer ID")
		return
	}

	offer, err := h.offerService.GetOfferByID(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, "Offer not found")
		return
	}

	respondJSON(w, http.StatusOK, offer)
}

// UpdateOffer handles PUT /api/admin/offers/{id}
func (h *OfferHandler) UpdateOffer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid offer ID")
		return
	}

	var req UpdateOfferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	cmd := &application.UpdateOfferCommand{
		ID:                        id,
		Name:                      req.Name,
		OfferValue:                req.OfferValue,
		ApplyToSalePrice:          req.ApplyToSalePrice,
		Archived:                  req.Archived,
		AutomaticallyAdded:        req.AutomaticallyAdded,
		CombinableWithOtherOffers: req.CombinableWithOtherOffers,
		OfferDescription:          req.OfferDescription,
		EndDate:                   req.EndDate,
		MarketingMessage:          req.MarketingMessage,
		MaxUsesPerCustomer:        req.MaxUsesPerCustomer,
		MaxUses:                   req.MaxUses,
		OrderMinTotal:             req.OrderMinTotal,
		OfferPriority:             req.OfferPriority,
	}

	offer, err := h.offerService.UpdateOffer(r.Context(), cmd)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, offer)
}

// DeleteOffer handles DELETE /api/admin/offers/{id}
func (h *OfferHandler) DeleteOffer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid offer ID")
		return
	}

	err = h.offerService.DeleteOffer(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusNoContent, nil)
}

// GetActiveOffers handles GET /api/admin/offers
func (h *OfferHandler) GetActiveOffers(w http.ResponseWriter, r *http.Request) {
	offers, err := h.offerService.GetActiveOffers(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, offers)
}

// CreateOfferCode handles POST /api/admin/offers/{id}/codes
func (h *OfferHandler) CreateOfferCode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	offerID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid offer ID")
		return
	}

	var req CreateOfferCodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	cmd := &application.CreateOfferCodeCommand{
		Code:         req.Code,
		MaxUses:      req.MaxUses,
		EmailAddress: req.EmailAddress,
		StartDate:    req.StartDate,
		EndDate:      req.EndDate,
	}

	offerCode, err := h.offerService.CreateOfferCode(r.Context(), offerID, cmd)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, offerCode)
}

// GetOfferCode handles GET /api/admin/offer-codes/{id}
func (h *OfferHandler) GetOfferCode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid offer code ID")
		return
	}

	offerCode, err := h.offerService.GetOfferCodeByID(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, "Offer code not found")
		return
	}

	respondJSON(w, http.StatusOK, offerCode)
}

// UpdateOfferCode handles PUT /api/admin/offer-codes/{id}
func (h *OfferHandler) UpdateOfferCode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid offer code ID")
		return
	}

	var req UpdateOfferCodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	cmd := &application.UpdateOfferCodeCommand{
		MaxUses:      req.MaxUses,
		EmailAddress: req.EmailAddress,
		StartDate:    req.StartDate,
		EndDate:      req.EndDate,
		Archived:     req.Archived,
	}

	offerCode, err := h.offerService.UpdateOfferCode(r.Context(), id, cmd)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, offerCode)
}

// DeleteOfferCode handles DELETE /api/admin/offer-codes/{id}
func (h *OfferHandler) DeleteOfferCode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid offer code ID")
		return
	}

	err = h.offerService.DeleteOfferCode(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusNoContent, nil)
}

// ProcessOrderOffers handles POST /api/storefront/orders/{orderId}/process-offers
func (h *OfferHandler) ProcessOrderOffers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID, err := strconv.ParseInt(vars["orderId"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid order ID")
		return
	}

	var req ProcessOffersAPIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	orderSubtotal, _ := decimal.NewFromString(req.OrderSubtotal)
	orderTotal, _ := decimal.NewFromString(req.OrderTotal)

	items := make([]application.OrderItemData, len(req.Items))
	for i, item := range req.Items {
		price, _ := decimal.NewFromString(item.Price)
		subtotal, _ := decimal.NewFromString(item.Subtotal)
		var salePrice *decimal.Decimal
		if item.SalePrice != "" {
			sp, _ := decimal.NewFromString(item.SalePrice)
			salePrice = &sp
		}

		items[i] = application.OrderItemData{
			ItemID:     item.ItemID,
			SKUID:      item.SKUID,
			CategoryID: item.CategoryID,
			Price:      price,
			SalePrice:  salePrice,
			Quantity:   item.Quantity,
			Subtotal:   subtotal,
			ProductID:  item.ProductID,
		}
	}

	request := &application.ProcessOffersRequest{
		OrderID:       orderID,
		OrderSubtotal: orderSubtotal,
		OrderTotal:    orderTotal,
		CustomerID:    req.CustomerID,
		Items:         items,
	}

	response, err := h.offerProcessorService.ProcessOrderOffers(r.Context(), request)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, response)
}

// ApplyOfferCode handles POST /api/storefront/orders/{orderId}/apply-code
func (h *OfferHandler) ApplyOfferCode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID, err := strconv.ParseInt(vars["orderId"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid order ID")
		return
	}

	var req ApplyOfferCodeAPIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	orderSubtotal, _ := decimal.NewFromString(req.OrderSubtotal)
	orderTotal, _ := decimal.NewFromString(req.OrderTotal)

	items := make([]application.OrderItemData, len(req.Items))
	for i, item := range req.Items {
		price, _ := decimal.NewFromString(item.Price)
		subtotal, _ := decimal.NewFromString(item.Subtotal)
		var salePrice *decimal.Decimal
		if item.SalePrice != "" {
			sp, _ := decimal.NewFromString(item.SalePrice)
			salePrice = &sp
		}

		items[i] = application.OrderItemData{
			ItemID:     item.ItemID,
			SKUID:      item.SKUID,
			CategoryID: item.CategoryID,
			Price:      price,
			SalePrice:  salePrice,
			Quantity:   item.Quantity,
			Subtotal:   subtotal,
			ProductID:  item.ProductID,
		}
	}

	request := &application.ApplyOfferCodeRequest{
		OrderID:       orderID,
		OfferCode:     req.OfferCode,
		OrderSubtotal: orderSubtotal,
		OrderTotal:    orderTotal,
		CustomerID:    req.CustomerID,
		Items:         items,
	}

	response, err := h.offerProcessorService.ApplyOfferCode(r.Context(), request)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, response)
}

// RemoveOffer handles DELETE /api/storefront/orders/{orderId}/offers/{offerId}
func (h *OfferHandler) RemoveOffer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID, err := strconv.ParseInt(vars["orderId"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid order ID")
		return
	}

	offerID, err := strconv.ParseInt(vars["offerId"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid offer ID")
		return
	}

	err = h.offerProcessorService.RemoveOfferFromOrder(r.Context(), orderID, offerID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusNoContent, nil)
}

// GetOfferByCode handles GET /api/storefront/offers/by-code/{code}
func (h *OfferHandler) GetOfferByCode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	code := vars["code"]

	offer, err := h.offerService.GetOfferByCode(r.Context(), code)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if offer == nil {
		respondError(w, http.StatusNotFound, "Offer not found")
		return
	}

	respondJSON(w, http.StatusOK, offer)
}

// Request/Response DTOs

type CreateOfferRequest struct {
	Name                      string     `json:"name"`
	OfferType                 string     `json:"offer_type"`
	OfferValue                float64    `json:"offer_value"`
	AdjustmentType            string     `json:"adjustment_type"`
	ApplyToChildItems         bool       `json:"apply_to_child_items"`
	ApplyToSalePrice          bool       `json:"apply_to_sale_price"`
	AutomaticallyAdded        bool       `json:"automatically_added"`
	CombinableWithOtherOffers bool       `json:"combinable_with_other_offers"`
	OfferDescription          string     `json:"offer_description"`
	OfferDiscountType         string     `json:"offer_discount_type"`
	EndDate                   *time.Time `json:"end_date"`
	MarketingMessage          string     `json:"marketing_message"`
	MaxUsesPerCustomer        *int64     `json:"max_uses_per_customer"`
	MaxUses                   *int       `json:"max_uses"`
	OrderMinTotal             float64    `json:"order_min_total"`
	OfferPriority             int        `json:"offer_priority"`
	StartDate                 time.Time  `json:"start_date"`
}

type UpdateOfferRequest struct {
	Name                      *string    `json:"name"`
	OfferValue                *float64   `json:"offer_value"`
	ApplyToSalePrice          *bool      `json:"apply_to_sale_price"`
	Archived                  *bool      `json:"archived"`
	AutomaticallyAdded        *bool      `json:"automatically_added"`
	CombinableWithOtherOffers *bool      `json:"combinable_with_other_offers"`
	OfferDescription          *string    `json:"offer_description"`
	EndDate                   *time.Time `json:"end_date"`
	MarketingMessage          *string    `json:"marketing_message"`
	MaxUsesPerCustomer        *int64     `json:"max_uses_per_customer"`
	MaxUses                   *int       `json:"max_uses"`
	OrderMinTotal             *float64   `json:"order_min_total"`
	OfferPriority             *int       `json:"offer_priority"`
}

type CreateOfferCodeRequest struct {
	Code         string     `json:"code"`
	MaxUses      *int       `json:"max_uses"`
	EmailAddress *string    `json:"email_address"`
	StartDate    *time.Time `json:"start_date"`
	EndDate      *time.Time `json:"end_date"`
}

type UpdateOfferCodeRequest struct {
	MaxUses      *int       `json:"max_uses"`
	EmailAddress *string    `json:"email_address"`
	StartDate    *time.Time `json:"start_date"`
	EndDate      *time.Time `json:"end_date"`
	Archived     *bool      `json:"archived"`
}

type ProcessOffersAPIRequest struct {
	OrderSubtotal string              `json:"order_subtotal"`
	OrderTotal    string              `json:"order_total"`
	CustomerID    *string             `json:"customer_id"`
	Items         []OrderItemDataJSON `json:"items"`
}

type ApplyOfferCodeAPIRequest struct {
	OfferCode     string              `json:"offer_code"`
	OrderSubtotal string              `json:"order_subtotal"`
	OrderTotal    string              `json:"order_total"`
	CustomerID    *string             `json:"customer_id"`
	Items         []OrderItemDataJSON `json:"items"`
}

type OrderItemDataJSON struct {
	ItemID     string  `json:"item_id"`
	SKUID      string  `json:"sku_id"`
	CategoryID *string `json:"category_id"`
	Price      string  `json:"price"`
	SalePrice  string  `json:"sale_price,omitempty"`
	Quantity   int     `json:"quantity"`
	Subtotal   string  `json:"subtotal"`
	ProductID  *string `json:"product_id"`
}

// Helper functions

func respondJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func respondError(w http.ResponseWriter, statusCode int, message string) {
	respondJSON(w, statusCode, map[string]string{"error": message})
}
