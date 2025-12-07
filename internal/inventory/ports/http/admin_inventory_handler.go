package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/qhato/ecommerce/internal/inventory/application/commands"
	"github.com/qhato/ecommerce/internal/inventory/application/queries"
	"github.com/qhato/ecommerce/pkg/logger"
	"github.com/qhato/ecommerce/pkg/validator"
)

// AdminInventoryHandler handles admin inventory HTTP requests using CQRS
type AdminInventoryHandler struct {
	commandHandler *commands.InventoryCommandHandler
	queryService   *queries.InventoryQueryService
	validator      *validator.Validator
	log            logger.Logger
}

// NewAdminInventoryHandler creates a new admin inventory HTTP handler
func NewAdminInventoryHandler(
	commandHandler *commands.InventoryCommandHandler,
	queryService *queries.InventoryQueryService,
	validator *validator.Validator,
	log logger.Logger,
) *AdminInventoryHandler {
	return &AdminInventoryHandler{
		commandHandler: commandHandler,
		queryService:   queryService,
		validator:      validator,
		log:            log,
	}
}

// RegisterRoutes registers all admin inventory routes
func (h *AdminInventoryHandler) RegisterRoutes(r chi.Router) {
	r.Route("/admin/inventory", func(r chi.Router) {
		// Inventory Levels
		r.Post("/levels", h.CreateInventoryLevel)
		r.Get("/levels/{id}", h.GetInventoryLevel)
		r.Get("/levels/sku/{skuId}", h.GetInventoryBySKU)
		r.Get("/levels/warehouse/{warehouseId}", h.GetInventoryByWarehouse)
		r.Put("/levels/{id}", h.UpdateInventoryLevel)
		r.Delete("/levels/{id}", h.DeleteInventoryLevel)

		// Inventory Operations
		r.Post("/levels/{id}/adjust", h.AdjustInventory)
		r.Post("/levels/sku/{skuId}/set", h.SetInventory)
		r.Post("/bulk-adjust", h.BulkAdjustInventory)
		r.Post("/transfer", h.TransferInventory)

		// Inventory Queries
		r.Get("/low-stock", h.GetLowStockItems)
		r.Get("/backorderable", h.GetBackorderableItems)
		r.Post("/check-availability", h.CheckAvailability)

		// Reservations
		r.Post("/reservations", h.ReserveInventory)
		r.Get("/reservations/{id}", h.GetReservation)
		r.Get("/reservations/order/{orderId}", h.GetReservationsByOrder)
		r.Get("/reservations/expired", h.GetExpiredReservations)
		r.Post("/reservations/{id}/confirm", h.ConfirmReservation)
		r.Post("/reservations/{id}/release", h.ReleaseReservation)
		r.Post("/reservations/{id}/fulfill", h.FulfillReservation)
		r.Post("/reservations/{id}/extend", h.ExtendReservation)
		r.Post("/reservations/expire-old", h.ExpireReservations)
		r.Post("/reservations/order/{orderId}/release", h.ReleaseOrderReservations)
	})
}

// Inventory Level Management

func (h *AdminInventoryHandler) CreateInventoryLevel(w http.ResponseWriter, r *http.Request) {
	var cmd commands.CreateInventoryLevelCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.validator.Validate(cmd); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	level, err := h.commandHandler.HandleCreateInventoryLevel(r.Context(), cmd)
	if err != nil {
		h.log.WithError(err).Error("Failed to create inventory level")
		respondWithError(w, http.StatusInternalServerError, "Failed to create inventory level")
		return
	}

	respondWithJSON(w, http.StatusCreated, level)
}

func (h *AdminInventoryHandler) GetInventoryLevel(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	query := queries.GetInventoryLevelQuery{ID: id}
	level, err := h.queryService.GetInventoryLevel(r.Context(), query)
	if err != nil {
		h.log.WithError(err).Error("Failed to get inventory level")
		respondWithError(w, http.StatusNotFound, "Inventory level not found")
		return
	}

	respondWithJSON(w, http.StatusOK, level)
}

func (h *AdminInventoryHandler) GetInventoryBySKU(w http.ResponseWriter, r *http.Request) {
	skuID := chi.URLParam(r, "skuId")

	query := queries.GetInventoryBySKUQuery{SKUID: skuID}
	level, err := h.queryService.GetInventoryBySKU(r.Context(), query)
	if err != nil {
		h.log.WithError(err).Error("Failed to get inventory by SKU")
		respondWithError(w, http.StatusNotFound, "Inventory not found")
		return
	}

	respondWithJSON(w, http.StatusOK, level)
}

func (h *AdminInventoryHandler) GetInventoryByWarehouse(w http.ResponseWriter, r *http.Request) {
	warehouseID := chi.URLParam(r, "warehouseId")

	query := queries.GetInventoryByWarehouseQuery{WarehouseID: warehouseID}
	levels, err := h.queryService.GetInventoryByWarehouse(r.Context(), query)
	if err != nil {
		h.log.WithError(err).Error("Failed to get inventory by warehouse")
		respondWithError(w, http.StatusInternalServerError, "Failed to get inventory")
		return
	}

	respondWithJSON(w, http.StatusOK, levels)
}

func (h *AdminInventoryHandler) UpdateInventoryLevel(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var cmd commands.UpdateInventoryLevelCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	cmd.ID = id

	if err := h.validator.Validate(cmd); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	level, err := h.commandHandler.HandleUpdateInventoryLevel(r.Context(), cmd)
	if err != nil {
		h.log.WithError(err).Error("Failed to update inventory level")
		respondWithError(w, http.StatusInternalServerError, "Failed to update inventory level")
		return
	}

	respondWithJSON(w, http.StatusOK, level)
}

func (h *AdminInventoryHandler) DeleteInventoryLevel(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	cmd := commands.DeleteInventoryLevelCommand{ID: id}
	if err := h.commandHandler.HandleDeleteInventoryLevel(r.Context(), cmd); err != nil {
		h.log.WithError(err).Error("Failed to delete inventory level")
		respondWithError(w, http.StatusInternalServerError, "Failed to delete inventory level")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Inventory Operations

func (h *AdminInventoryHandler) AdjustInventory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req struct {
		Adjustment  int    `json:"adjustment"`
		Reason      string `json:"reason"`
		WarehouseID string `json:"warehouse_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// First get the SKU ID from the inventory level
	levelQuery := queries.GetInventoryLevelQuery{ID: id}
	level, err := h.queryService.GetInventoryLevel(r.Context(), levelQuery)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Inventory level not found")
		return
	}

	cmd := commands.AdjustInventoryCommand{
		SKUID:       level.SKUID,
		WarehouseID: req.WarehouseID,
		Adjustment:  req.Adjustment,
		Reason:      req.Reason,
	}

	updatedLevel, err := h.commandHandler.HandleAdjustInventory(r.Context(), cmd)
	if err != nil {
		h.log.WithError(err).Error("Failed to adjust inventory")
		respondWithError(w, http.StatusInternalServerError, "Failed to adjust inventory")
		return
	}

	respondWithJSON(w, http.StatusOK, updatedLevel)
}

func (h *AdminInventoryHandler) SetInventory(w http.ResponseWriter, r *http.Request) {
	skuID := chi.URLParam(r, "skuId")

	var req struct {
		NewQuantity int    `json:"new_quantity"`
		WarehouseID string `json:"warehouse_id"`
		Reason      string `json:"reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	cmd := commands.SetInventoryCommand{
		SKUID:       skuID,
		WarehouseID: req.WarehouseID,
		NewQuantity: req.NewQuantity,
		Reason:      req.Reason,
	}

	level, err := h.commandHandler.HandleSetInventory(r.Context(), cmd)
	if err != nil {
		h.log.WithError(err).Error("Failed to set inventory")
		respondWithError(w, http.StatusInternalServerError, "Failed to set inventory")
		return
	}

	respondWithJSON(w, http.StatusOK, level)
}

func (h *AdminInventoryHandler) BulkAdjustInventory(w http.ResponseWriter, r *http.Request) {
	var cmd commands.BulkAdjustInventoryCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.commandHandler.HandleBulkAdjustInventory(r.Context(), cmd); err != nil {
		h.log.WithError(err).Error("Failed to bulk adjust inventory")
		respondWithError(w, http.StatusInternalServerError, "Failed to bulk adjust inventory")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"status": "success"})
}

func (h *AdminInventoryHandler) TransferInventory(w http.ResponseWriter, r *http.Request) {
	var cmd commands.TransferInventoryCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.validator.Validate(cmd); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.commandHandler.HandleTransferInventory(r.Context(), cmd); err != nil {
		h.log.WithError(err).Error("Failed to transfer inventory")
		respondWithError(w, http.StatusInternalServerError, "Failed to transfer inventory")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"status": "transferred"})
}

// Inventory Queries

func (h *AdminInventoryHandler) GetLowStockItems(w http.ResponseWriter, r *http.Request) {
	warehouseID := r.URL.Query().Get("warehouse_id")
	limitStr := r.URL.Query().Get("limit")

	var limit int
	if limitStr != "" {
		parsed, err := strconv.Atoi(limitStr)
		if err == nil {
			limit = parsed
		}
	}

	var warehousePtr *string
	if warehouseID != "" {
		warehousePtr = &warehouseID
	}

	query := queries.GetLowStockItemsQuery{
		WarehouseID: warehousePtr,
		Limit:       limit,
	}

	items, err := h.queryService.GetLowStockItems(r.Context(), query)
	if err != nil {
		h.log.WithError(err).Error("Failed to get low stock items")
		respondWithError(w, http.StatusInternalServerError, "Failed to get low stock items")
		return
	}

	respondWithJSON(w, http.StatusOK, items)
}

func (h *AdminInventoryHandler) GetBackorderableItems(w http.ResponseWriter, r *http.Request) {
	warehouseID := r.URL.Query().Get("warehouse_id")

	var warehousePtr *string
	if warehouseID != "" {
		warehousePtr = &warehouseID
	}

	query := queries.GetBackorderableItemsQuery{
		WarehouseID: warehousePtr,
	}

	items, err := h.queryService.GetBackorderableItems(r.Context(), query)
	if err != nil {
		h.log.WithError(err).Error("Failed to get backorderable items")
		respondWithError(w, http.StatusInternalServerError, "Failed to get backorderable items")
		return
	}

	respondWithJSON(w, http.StatusOK, items)
}

func (h *AdminInventoryHandler) CheckAvailability(w http.ResponseWriter, r *http.Request) {
	var req struct {
		SKUID    string `json:"sku_id"`
		Quantity int    `json:"quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	query := queries.CheckInventoryAvailabilityQuery{
		SKUID:    req.SKUID,
		Quantity: req.Quantity,
	}

	availability, err := h.queryService.CheckInventoryAvailability(r.Context(), query)
	if err != nil {
		h.log.WithError(err).Error("Failed to check inventory availability")
		respondWithError(w, http.StatusInternalServerError, "Failed to check availability")
		return
	}

	respondWithJSON(w, http.StatusOK, availability)
}

// Reservation Management

func (h *AdminInventoryHandler) ReserveInventory(w http.ResponseWriter, r *http.Request) {
	var req struct {
		SKUID       string `json:"sku_id"`
		Quantity    int    `json:"quantity"`
		OrderID     string `json:"order_id"`
		OrderItemID string `json:"order_item_id"`
		TTLSeconds  int64  `json:"ttl_seconds"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	ttl := time.Duration(req.TTLSeconds) * time.Second
	if ttl == 0 {
		ttl = 24 * time.Hour // Default 24 hours
	}

	cmd := commands.ReserveInventoryCommand{
		SKUID:       req.SKUID,
		Quantity:    req.Quantity,
		OrderID:     req.OrderID,
		OrderItemID: req.OrderItemID,
		TTL:         ttl,
	}

	reservation, err := h.commandHandler.HandleReserveInventory(r.Context(), cmd)
	if err != nil {
		h.log.WithError(err).Error("Failed to reserve inventory")
		respondWithError(w, http.StatusInternalServerError, "Failed to reserve inventory")
		return
	}

	respondWithJSON(w, http.StatusCreated, reservation)
}

func (h *AdminInventoryHandler) GetReservation(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	query := queries.GetReservationQuery{ID: id}
	reservation, err := h.queryService.GetReservation(r.Context(), query)
	if err != nil {
		h.log.WithError(err).Error("Failed to get reservation")
		respondWithError(w, http.StatusNotFound, "Reservation not found")
		return
	}

	respondWithJSON(w, http.StatusOK, reservation)
}

func (h *AdminInventoryHandler) GetReservationsByOrder(w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "orderId")

	query := queries.GetReservationsByOrderQuery{OrderID: orderID}
	reservations, err := h.queryService.GetReservationsByOrder(r.Context(), query)
	if err != nil {
		h.log.WithError(err).Error("Failed to get reservations by order")
		respondWithError(w, http.StatusInternalServerError, "Failed to get reservations")
		return
	}

	respondWithJSON(w, http.StatusOK, reservations)
}

func (h *AdminInventoryHandler) GetExpiredReservations(w http.ResponseWriter, r *http.Request) {
	query := queries.GetExpiredReservationsQuery{}
	reservations, err := h.queryService.GetExpiredReservations(r.Context(), query)
	if err != nil {
		h.log.WithError(err).Error("Failed to get expired reservations")
		respondWithError(w, http.StatusInternalServerError, "Failed to get expired reservations")
		return
	}

	respondWithJSON(w, http.StatusOK, reservations)
}

func (h *AdminInventoryHandler) ConfirmReservation(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	cmd := commands.ConfirmReservationCommand{ReservationID: id}
	reservation, err := h.commandHandler.HandleConfirmReservation(r.Context(), cmd)
	if err != nil {
		h.log.WithError(err).Error("Failed to confirm reservation")
		respondWithError(w, http.StatusInternalServerError, "Failed to confirm reservation")
		return
	}

	respondWithJSON(w, http.StatusOK, reservation)
}

func (h *AdminInventoryHandler) ReleaseReservation(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	cmd := commands.ReleaseReservationCommand{ReservationID: id}
	if err := h.commandHandler.HandleReleaseReservation(r.Context(), cmd); err != nil {
		h.log.WithError(err).Error("Failed to release reservation")
		respondWithError(w, http.StatusInternalServerError, "Failed to release reservation")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"status": "released"})
}

func (h *AdminInventoryHandler) FulfillReservation(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	cmd := commands.FulfillReservationCommand{ReservationID: id}
	if err := h.commandHandler.HandleFulfillReservation(r.Context(), cmd); err != nil {
		h.log.WithError(err).Error("Failed to fulfill reservation")
		respondWithError(w, http.StatusInternalServerError, "Failed to fulfill reservation")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"status": "fulfilled"})
}

func (h *AdminInventoryHandler) ExtendReservation(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req struct {
		AdditionalSeconds int64 `json:"additional_seconds"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	cmd := commands.ExtendReservationCommand{
		ReservationID:  id,
		AdditionalTime: time.Duration(req.AdditionalSeconds) * time.Second,
	}

	reservation, err := h.commandHandler.HandleExtendReservation(r.Context(), cmd)
	if err != nil {
		h.log.WithError(err).Error("Failed to extend reservation")
		respondWithError(w, http.StatusInternalServerError, "Failed to extend reservation")
		return
	}

	respondWithJSON(w, http.StatusOK, reservation)
}

func (h *AdminInventoryHandler) ExpireReservations(w http.ResponseWriter, r *http.Request) {
	cmd := commands.ExpireReservationsCommand{}
	count, err := h.commandHandler.HandleExpireReservations(r.Context(), cmd)
	if err != nil {
		h.log.WithError(err).Error("Failed to expire reservations")
		respondWithError(w, http.StatusInternalServerError, "Failed to expire reservations")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]int{"expired_count": count})
}

func (h *AdminInventoryHandler) ReleaseOrderReservations(w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "orderId")

	cmd := commands.ReleaseOrderReservationsCommand{OrderID: orderID}
	if err := h.commandHandler.HandleReleaseOrderReservations(r.Context(), cmd); err != nil {
		h.log.WithError(err).Error("Failed to release order reservations")
		respondWithError(w, http.StatusInternalServerError, "Failed to release order reservations")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"status": "released"})
}

// Helper functions

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"Failed to marshal response"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
