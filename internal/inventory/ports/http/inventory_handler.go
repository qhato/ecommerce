package http

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/qhato/ecommerce/internal/inventory/application"
)

type InventoryHandler struct {
	inventoryService application.InventoryService
}

func NewInventoryHandler(inventoryService application.InventoryService) *InventoryHandler {
	return &InventoryHandler{
		inventoryService: inventoryService,
	}
}

func (h *InventoryHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/inventory", h.CreateInventoryLevel).Methods("POST")
	router.HandleFunc("/inventory/{id}", h.GetInventoryLevel).Methods("GET")
	router.HandleFunc("/inventory/sku/{skuId}", h.GetInventoryBySKU).Methods("GET")
	router.HandleFunc("/inventory/{id}", h.DeleteInventoryLevel).Methods("DELETE")
	router.HandleFunc("/inventory/{id}/increment", h.IncrementInventory).Methods("POST")
	router.HandleFunc("/inventory/{id}/decrement", h.DecrementInventory).Methods("POST")
	router.HandleFunc("/inventory/{id}/reserve", h.ReserveInventory).Methods("POST")
	router.HandleFunc("/inventory/{id}/release", h.ReleaseInventory).Methods("POST")
	router.HandleFunc("/inventory/{id}/quantities", h.UpdateQuantities).Methods("PUT")
}

type CreateInventoryLevelRequest struct {
	SKUID          string `json:"sku_id" validate:"required"`
	QuantityOnHand int    `json:"quantity_on_hand" validate:"min=0"`
}

type QuantityOperationRequest struct {
	Quantity int `json:"quantity" validate:"required,min=1"`
}

type UpdateQuantitiesRequest struct {
	QuantityOnHand   int `json:"quantity_on_hand" validate:"min=0"`
	QuantityReserved int `json:"quantity_reserved" validate:"min=0"`
}

func (h *InventoryHandler) CreateInventoryLevel(w http.ResponseWriter, r *http.Request) {
	var req CreateInventoryLevelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cmd := &application.CreateInventoryLevelCommand{
		SKUID:          req.SKUID,
		QuantityOnHand: req.QuantityOnHand,
	}

	level, err := h.inventoryService.CreateInventoryLevel(r.Context(), cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(level)
}

func (h *InventoryHandler) GetInventoryLevel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	level, err := h.inventoryService.GetInventoryLevelByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(level)
}

func (h *InventoryHandler) GetInventoryBySKU(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	skuID := vars["skuId"]

	level, err := h.inventoryService.GetInventoryLevelBySKUID(r.Context(), skuID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(level)
}

func (h *InventoryHandler) DeleteInventoryLevel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.inventoryService.DeleteInventoryLevel(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *InventoryHandler) IncrementInventory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req QuantityOperationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	level, err := h.inventoryService.IncrementInventory(r.Context(), id, req.Quantity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(level)
}

func (h *InventoryHandler) DecrementInventory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req QuantityOperationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	level, err := h.inventoryService.DecrementInventory(r.Context(), id, req.Quantity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(level)
}

func (h *InventoryHandler) ReserveInventory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req QuantityOperationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	level, err := h.inventoryService.ReserveInventory(r.Context(), id, req.Quantity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(level)
}

func (h *InventoryHandler) ReleaseInventory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req QuantityOperationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	level, err := h.inventoryService.ReleaseInventory(r.Context(), id, req.Quantity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(level)
}

func (h *InventoryHandler) UpdateQuantities(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req UpdateQuantitiesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	level, err := h.inventoryService.UpdateInventoryQuantities(r.Context(), id, req.QuantityOnHand, req.QuantityReserved)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(level)
}
