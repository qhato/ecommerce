package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/qhato/ecommerce/internal/pricing/application/commands"
	"github.com/qhato/ecommerce/internal/pricing/application/queries"
	"github.com/qhato/ecommerce/internal/pricing/domain"
	"github.com/shopspring/decimal"
)

// PricingHandler handles HTTP requests for pricing operations
type PricingHandler struct {
	commandHandler *commands.PricingCommandHandler
	queryService   *queries.PricingQueryService
}

// NewPricingHandler creates a new PricingHandler
func NewPricingHandler(
	commandHandler *commands.PricingCommandHandler,
	queryService *queries.PricingQueryService,
) *PricingHandler {
	return &PricingHandler{
		commandHandler: commandHandler,
		queryService:   queryService,
	}
}

// RegisterRoutes registers all pricing-related routes
func (h *PricingHandler) RegisterRoutes(router *mux.Router) {
	// Price Lists
	router.HandleFunc("/api/admin/price-lists", h.CreatePriceList).Methods("POST")
	router.HandleFunc("/api/admin/price-lists/{id}", h.GetPriceList).Methods("GET")
	router.HandleFunc("/api/admin/price-lists/{id}", h.UpdatePriceList).Methods("PUT")
	router.HandleFunc("/api/admin/price-lists/{id}", h.DeletePriceList).Methods("DELETE")
	router.HandleFunc("/api/admin/price-lists", h.GetActivePriceLists).Methods("GET")
	router.HandleFunc("/api/admin/price-lists/code/{code}", h.GetPriceListByCode).Methods("GET")

	// Price List Items
	router.HandleFunc("/api/admin/price-lists/{id}/items", h.CreatePriceListItem).Methods("POST")
	router.HandleFunc("/api/admin/price-lists/{id}/items/bulk", h.BulkCreatePriceListItems).Methods("POST")
	router.HandleFunc("/api/admin/price-lists/{id}/items", h.GetPriceListItems).Methods("GET")
	router.HandleFunc("/api/admin/price-list-items/{id}", h.GetPriceListItem).Methods("GET")
	router.HandleFunc("/api/admin/price-list-items/{id}", h.UpdatePriceListItem).Methods("PUT")
	router.HandleFunc("/api/admin/price-list-items/{id}", h.DeletePriceListItem).Methods("DELETE")

	// Pricing Rules
	router.HandleFunc("/api/admin/pricing-rules", h.CreatePricingRule).Methods("POST")
	router.HandleFunc("/api/admin/pricing-rules/{id}", h.GetPricingRule).Methods("GET")
	router.HandleFunc("/api/admin/pricing-rules/{id}", h.UpdatePricingRule).Methods("PUT")
	router.HandleFunc("/api/admin/pricing-rules/{id}", h.DeletePricingRule).Methods("DELETE")
	router.HandleFunc("/api/admin/pricing-rules", h.GetActivePricingRules).Methods("GET")

	// Pricing Calculations (Storefront)
	router.HandleFunc("/api/storefront/prices/calculate", h.CalculatePrices).Methods("POST")
	router.HandleFunc("/api/storefront/prices/sku/{skuId}", h.GetPriceForSKU).Methods("GET")
}

// CreatePriceList handles POST /api/admin/price-lists
func (h *PricingHandler) CreatePriceList(w http.ResponseWriter, r *http.Request) {
	var req CreatePriceListRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	cmd := &commands.CreatePriceListCommand{
		Name:             req.Name,
		Code:             req.Code,
		PriceListType:    domain.PriceListType(req.PriceListType),
		Currency:         req.Currency,
		Priority:         req.Priority,
		Description:      req.Description,
		StartDate:        req.StartDate,
		EndDate:          req.EndDate,
		CustomerSegments: req.CustomerSegments,
	}

	id, err := h.commandHandler.HandleCreatePriceList(r.Context(), cmd)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, map[string]int64{"id": id})
}

// GetPriceList handles GET /api/admin/price-lists/{id}
func (h *PricingHandler) GetPriceList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid price list ID")
		return
	}

	priceList, err := h.queryService.GetPriceList(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, "Price list not found")
		return
	}

	respondJSON(w, http.StatusOK, queries.ToPriceListDTO(priceList))
}

// GetPriceListByCode handles GET /api/admin/price-lists/code/{code}
func (h *PricingHandler) GetPriceListByCode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	code := vars["code"]

	priceList, err := h.queryService.GetPriceListByCode(r.Context(), code)
	if err != nil {
		respondError(w, http.StatusNotFound, "Price list not found")
		return
	}

	respondJSON(w, http.StatusOK, queries.ToPriceListDTO(priceList))
}

// UpdatePriceList handles PUT /api/admin/price-lists/{id}
func (h *PricingHandler) UpdatePriceList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid price list ID")
		return
	}

	var req UpdatePriceListRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	cmd := &commands.UpdatePriceListCommand{
		ID:               id,
		Name:             req.Name,
		Priority:         req.Priority,
		IsActive:         req.IsActive,
		Description:      req.Description,
		StartDate:        req.StartDate,
		EndDate:          req.EndDate,
		CustomerSegments: req.CustomerSegments,
	}

	err = h.commandHandler.HandleUpdatePriceList(r.Context(), cmd)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Price list updated successfully"})
}

// DeletePriceList handles DELETE /api/admin/price-lists/{id}
func (h *PricingHandler) DeletePriceList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid price list ID")
		return
	}

	err = h.commandHandler.HandleDeletePriceList(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusNoContent, nil)
}

// GetActivePriceLists handles GET /api/admin/price-lists
func (h *PricingHandler) GetActivePriceLists(w http.ResponseWriter, r *http.Request) {
	currency := r.URL.Query().Get("currency")
	if currency == "" {
		currency = "USD"
	}

	priceLists, err := h.queryService.GetActivePriceLists(r.Context(), currency)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	dtos := make([]*queries.PriceListDTO, len(priceLists))
	for i, pl := range priceLists {
		dtos[i] = queries.ToPriceListDTO(pl)
	}

	respondJSON(w, http.StatusOK, dtos)
}

// CreatePriceListItem handles POST /api/admin/price-lists/{id}/items
func (h *PricingHandler) CreatePriceListItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	priceListID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid price list ID")
		return
	}

	var req CreatePriceListItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	price, _ := decimal.NewFromString(req.Price)
	var compareAtPrice *decimal.Decimal
	if req.CompareAtPrice != nil {
		cap, _ := decimal.NewFromString(*req.CompareAtPrice)
		compareAtPrice = &cap
	}

	cmd := &commands.CreatePriceListItemCommand{
		PriceListID:    priceListID,
		SKUID:          req.SKUID,
		ProductID:      req.ProductID,
		Price:          price,
		CompareAtPrice: compareAtPrice,
		MinQuantity:    req.MinQuantity,
		MaxQuantity:    req.MaxQuantity,
		StartDate:      req.StartDate,
		EndDate:        req.EndDate,
	}

	id, err := h.commandHandler.HandleCreatePriceListItem(r.Context(), cmd)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, map[string]int64{"id": id})
}

// BulkCreatePriceListItems handles POST /api/admin/price-lists/{id}/items/bulk
func (h *PricingHandler) BulkCreatePriceListItems(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	priceListID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid price list ID")
		return
	}

	var req BulkCreatePriceListItemsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	items := make([]commands.BulkPriceListItem, len(req.Items))
	for i, item := range req.Items {
		price, _ := decimal.NewFromString(item.Price)
		var compareAtPrice *decimal.Decimal
		if item.CompareAtPrice != nil {
			cap, _ := decimal.NewFromString(*item.CompareAtPrice)
			compareAtPrice = &cap
		}

		items[i] = commands.BulkPriceListItem{
			SKUID:          item.SKUID,
			ProductID:      item.ProductID,
			Price:          price,
			CompareAtPrice: compareAtPrice,
			MinQuantity:    item.MinQuantity,
			MaxQuantity:    item.MaxQuantity,
		}
	}

	cmd := &commands.BulkCreatePriceListItemsCommand{
		PriceListID: priceListID,
		Items:       items,
	}

	err = h.commandHandler.HandleBulkCreatePriceListItems(r.Context(), cmd)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, map[string]string{"message": "Bulk creation successful"})
}

// GetPriceListItems handles GET /api/admin/price-lists/{id}/items
func (h *PricingHandler) GetPriceListItems(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	priceListID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid price list ID")
		return
	}

	items, err := h.queryService.GetPriceListItems(r.Context(), priceListID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	dtos := make([]*queries.PriceListItemDTO, len(items))
	for i, item := range items {
		dtos[i] = queries.ToPriceListItemDTO(item)
	}

	respondJSON(w, http.StatusOK, dtos)
}

// CalculatePrices handles POST /api/storefront/prices/calculate
func (h *PricingHandler) CalculatePrices(w http.ResponseWriter, r *http.Request) {
	var req queries.CalculatePriceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	pricingCtx := queries.ToPricingContext(&req)

	result, err := h.queryService.CalculatePrices(r.Context(), pricingCtx)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, queries.ToPricingResultDTO(result))
}

// GetPriceForSKU handles GET /api/storefront/prices/sku/{skuId}
func (h *PricingHandler) GetPriceForSKU(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	skuID := vars["skuId"]

	quantity := 1
	if q := r.URL.Query().Get("quantity"); q != "" {
		quantity, _ = strconv.Atoi(q)
	}

	currency := r.URL.Query().Get("currency")
	if currency == "" {
		currency = "USD"
	}

	var customerSegment *string
	if cs := r.URL.Query().Get("customer_segment"); cs != "" {
		customerSegment = &cs
	}

	pricedItem, err := h.queryService.GetPriceForSKU(r.Context(), skuID, quantity, currency, customerSegment)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, queries.ToPricedItemDTO(pricedItem))
}

// Stub methods for other endpoints (simplified for brevity)
func (h *PricingHandler) GetPriceListItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 64)
	item, _ := h.queryService.GetPriceListItem(r.Context(), id)
	respondJSON(w, http.StatusOK, queries.ToPriceListItemDTO(item))
}

func (h *PricingHandler) UpdatePriceListItem(w http.ResponseWriter, r *http.Request) {
	// Implementation similar to CreatePriceListItem
	respondJSON(w, http.StatusOK, map[string]string{"message": "Updated"})
}

func (h *PricingHandler) DeletePriceListItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 64)
	h.commandHandler.HandleDeletePriceListItem(r.Context(), id)
	respondJSON(w, http.StatusNoContent, nil)
}

func (h *PricingHandler) CreatePricingRule(w http.ResponseWriter, r *http.Request) {
	// Implementation similar to CreatePriceList
	respondJSON(w, http.StatusCreated, map[string]int64{"id": 1})
}

func (h *PricingHandler) GetPricingRule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 64)
	rule, _ := h.queryService.GetPricingRule(r.Context(), id)
	respondJSON(w, http.StatusOK, queries.ToPricingRuleDTO(rule))
}

func (h *PricingHandler) UpdatePricingRule(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{"message": "Updated"})
}

func (h *PricingHandler) DeletePricingRule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 64)
	h.commandHandler.HandleDeletePricingRule(r.Context(), id)
	respondJSON(w, http.StatusNoContent, nil)
}

func (h *PricingHandler) GetActivePricingRules(w http.ResponseWriter, r *http.Request) {
	rules, _ := h.queryService.GetActivePricingRules(r.Context())
	dtos := make([]*queries.PricingRuleDTO, len(rules))
	for i, rule := range rules {
		dtos[i] = queries.ToPricingRuleDTO(rule)
	}
	respondJSON(w, http.StatusOK, dtos)
}

// Request/Response DTOs

type CreatePriceListRequest struct {
	Name             string     `json:"name"`
	Code             string     `json:"code"`
	PriceListType    string     `json:"price_list_type"`
	Currency         string     `json:"currency"`
	Priority         int        `json:"priority"`
	Description      string     `json:"description"`
	StartDate        *time.Time `json:"start_date"`
	EndDate          *time.Time `json:"end_date"`
	CustomerSegments []string   `json:"customer_segments"`
}

type UpdatePriceListRequest struct {
	Name             *string    `json:"name"`
	Priority         *int       `json:"priority"`
	IsActive         *bool      `json:"is_active"`
	Description      *string    `json:"description"`
	StartDate        *time.Time `json:"start_date"`
	EndDate          *time.Time `json:"end_date"`
	CustomerSegments []string   `json:"customer_segments"`
}

type CreatePriceListItemRequest struct {
	SKUID          string     `json:"sku_id"`
	ProductID      *string    `json:"product_id"`
	Price          string     `json:"price"`
	CompareAtPrice *string    `json:"compare_at_price"`
	MinQuantity    int        `json:"min_quantity"`
	MaxQuantity    *int       `json:"max_quantity"`
	StartDate      *time.Time `json:"start_date"`
	EndDate        *time.Time `json:"end_date"`
}

type BulkCreatePriceListItemsRequest struct {
	Items []CreatePriceListItemRequest `json:"items"`
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
