package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/qhato/ecommerce/internal/store/application/commands"
	"github.com/qhato/ecommerce/internal/store/application/queries"
	"github.com/qhato/ecommerce/internal/store/domain"
)

type StoreHandler struct {
	commandHandler *commands.StoreCommandHandler
	queryService   *queries.StoreQueryService
}

func NewStoreHandler(
	commandHandler *commands.StoreCommandHandler,
	queryService   *queries.StoreQueryService,
) *StoreHandler {
	return &StoreHandler{
		commandHandler: commandHandler,
		queryService:   queryService,
	}
}

func (h *StoreHandler) RegisterRoutes(router *mux.Router) {
	// Store management
	router.HandleFunc("/stores", h.CreateStore).Methods("POST")
	router.HandleFunc("/stores/{id}", h.GetStore).Methods("GET")
	router.HandleFunc("/stores/{id}", h.UpdateStore).Methods("PUT")
	router.HandleFunc("/stores/{id}", h.DeleteStore).Methods("DELETE")
	router.HandleFunc("/stores/code/{code}", h.GetStoreByCode).Methods("GET")
	router.HandleFunc("/stores/status/{status}", h.GetStoresByStatus).Methods("GET")
	router.HandleFunc("/stores/type/{type}", h.GetStoresByType).Methods("GET")
	router.HandleFunc("/stores", h.GetAllStores).Methods("GET")
	router.HandleFunc("/stores/nearby", h.GetNearbyStores).Methods("GET")
	router.HandleFunc("/stores/{id}/activate", h.ActivateStore).Methods("POST")
	router.HandleFunc("/stores/{id}/deactivate", h.DeactivateStore).Methods("POST")
	router.HandleFunc("/stores/{id}/close", h.CloseStore).Methods("POST")

	// Inventory management
	router.HandleFunc("/stores/{storeId}/inventory", h.GetStoreInventory).Methods("GET")
	router.HandleFunc("/stores/{storeId}/inventory/product/{productId}", h.GetInventoryForProduct).Methods("GET")
	router.HandleFunc("/stores/{storeId}/inventory/low-stock", h.GetLowStockItems).Methods("GET")
	router.HandleFunc("/stores/inventory/product/{productId}", h.GetInventoryByProduct).Methods("GET")
	router.HandleFunc("/stores/inventory/sku/{sku}", h.GetInventoryBySKU).Methods("GET")
	router.HandleFunc("/stores/{storeId}/inventory/update", h.UpdateInventory).Methods("POST")
	router.HandleFunc("/stores/{storeId}/inventory/reserve", h.ReserveInventory).Methods("POST")
	router.HandleFunc("/stores/{storeId}/inventory/release", h.ReleaseInventory).Methods("POST")
}

func (h *StoreHandler) CreateStore(w http.ResponseWriter, r *http.Request) {
	var cmd commands.CreateStoreCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	store, err := h.commandHandler.HandleCreateStore(r.Context(), cmd)
	if err != nil {
		if err == domain.ErrStoreCodeTaken {
			http.Error(w, "Store code already taken", http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(queries.ToStoreDTO(store))
}

func (h *StoreHandler) GetStore(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid store ID", http.StatusBadRequest)
		return
	}

	store, err := h.queryService.GetStore(r.Context(), id)
	if err != nil {
		if err == domain.ErrStoreNotFound {
			http.Error(w, "Store not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(store)
}

func (h *StoreHandler) UpdateStore(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid store ID", http.StatusBadRequest)
		return
	}

	var cmd commands.UpdateStoreCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	cmd.ID = id

	store, err := h.commandHandler.HandleUpdateStore(r.Context(), cmd)
	if err != nil {
		if err == domain.ErrStoreNotFound {
			http.Error(w, "Store not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(queries.ToStoreDTO(store))
}

func (h *StoreHandler) DeleteStore(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid store ID", http.StatusBadRequest)
		return
	}

	cmd := commands.DeleteStoreCommand{ID: id}
	if err := h.commandHandler.HandleDeleteStore(r.Context(), cmd); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *StoreHandler) GetStoreByCode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	code := vars["code"]

	store, err := h.queryService.GetStoreByCode(r.Context(), code)
	if err != nil {
		if err == domain.ErrStoreNotFound {
			http.Error(w, "Store not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(store)
}

func (h *StoreHandler) GetStoresByStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	status := vars["status"]

	stores, err := h.queryService.GetStoresByStatus(r.Context(), status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stores)
}

func (h *StoreHandler) GetStoresByType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	storeType := vars["type"]

	stores, err := h.queryService.GetStoresByType(r.Context(), storeType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stores)
}

func (h *StoreHandler) GetAllStores(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit := 100
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil {
			limit = parsedLimit
		}
	}

	stores, err := h.queryService.GetAllStores(r.Context(), limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stores)
}

func (h *StoreHandler) GetNearbyStores(w http.ResponseWriter, r *http.Request) {
	latStr := r.URL.Query().Get("lat")
	lngStr := r.URL.Query().Get("lng")
	radiusStr := r.URL.Query().Get("radius")
	limitStr := r.URL.Query().Get("limit")

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		http.Error(w, "Invalid latitude", http.StatusBadRequest)
		return
	}

	lng, err := strconv.ParseFloat(lngStr, 64)
	if err != nil {
		http.Error(w, "Invalid longitude", http.StatusBadRequest)
		return
	}

	radius := 50.0 // Default 50km
	if radiusStr != "" {
		if parsedRadius, err := strconv.ParseFloat(radiusStr, 64); err == nil {
			radius = parsedRadius
		}
	}

	limit := 10
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil {
			limit = parsedLimit
		}
	}

	stores, err := h.queryService.GetNearbyStores(r.Context(), lat, lng, radius, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stores)
}

func (h *StoreHandler) ActivateStore(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid store ID", http.StatusBadRequest)
		return
	}

	cmd := commands.ActivateStoreCommand{ID: id}
	store, err := h.commandHandler.HandleActivateStore(r.Context(), cmd)
	if err != nil {
		if err == domain.ErrStoreNotFound {
			http.Error(w, "Store not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(queries.ToStoreDTO(store))
}

func (h *StoreHandler) DeactivateStore(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid store ID", http.StatusBadRequest)
		return
	}

	cmd := commands.DeactivateStoreCommand{ID: id}
	store, err := h.commandHandler.HandleDeactivateStore(r.Context(), cmd)
	if err != nil {
		if err == domain.ErrStoreNotFound {
			http.Error(w, "Store not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(queries.ToStoreDTO(store))
}

func (h *StoreHandler) CloseStore(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid store ID", http.StatusBadRequest)
		return
	}

	cmd := commands.CloseStoreCommand{ID: id}
	store, err := h.commandHandler.HandleCloseStore(r.Context(), cmd)
	if err != nil {
		if err == domain.ErrStoreNotFound {
			http.Error(w, "Store not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(queries.ToStoreDTO(store))
}

// Inventory handlers

func (h *StoreHandler) GetStoreInventory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	storeID, err := strconv.ParseInt(vars["storeId"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid store ID", http.StatusBadRequest)
		return
	}

	inventory, err := h.queryService.GetInventoryByStore(r.Context(), storeID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(inventory)
}

func (h *StoreHandler) GetInventoryForProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	storeID, err := strconv.ParseInt(vars["storeId"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid store ID", http.StatusBadRequest)
		return
	}

	productID, err := strconv.ParseInt(vars["productId"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	inventory, err := h.queryService.GetStoreInventory(r.Context(), storeID, productID)
	if err != nil {
		if err == domain.ErrInventoryNotFound {
			http.Error(w, "Inventory not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(inventory)
}

func (h *StoreHandler) GetLowStockItems(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	storeID, err := strconv.ParseInt(vars["storeId"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid store ID", http.StatusBadRequest)
		return
	}

	inventory, err := h.queryService.GetLowStockItems(r.Context(), storeID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(inventory)
}

func (h *StoreHandler) GetInventoryByProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID, err := strconv.ParseInt(vars["productId"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	inventory, err := h.queryService.GetInventoryByProduct(r.Context(), productID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(inventory)
}

func (h *StoreHandler) GetInventoryBySKU(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sku := vars["sku"]

	inventory, err := h.queryService.GetInventoryBySKU(r.Context(), sku)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(inventory)
}

func (h *StoreHandler) UpdateInventory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	storeID, err := strconv.ParseInt(vars["storeId"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid store ID", http.StatusBadRequest)
		return
	}

	var cmd commands.UpdateInventoryCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	cmd.StoreID = storeID

	inventory, err := h.commandHandler.HandleUpdateInventory(r.Context(), cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(queries.ToStoreInventoryDTO(inventory))
}

func (h *StoreHandler) ReserveInventory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	storeID, err := strconv.ParseInt(vars["storeId"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid store ID", http.StatusBadRequest)
		return
	}

	var cmd commands.ReserveInventoryCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	cmd.StoreID = storeID

	inventory, err := h.commandHandler.HandleReserveInventory(r.Context(), cmd)
	if err != nil {
		if err == domain.ErrInsufficientInventory {
			http.Error(w, "Insufficient inventory", http.StatusConflict)
			return
		}
		if err == domain.ErrInventoryNotFound {
			http.Error(w, "Inventory not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(queries.ToStoreInventoryDTO(inventory))
}

func (h *StoreHandler) ReleaseInventory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	storeID, err := strconv.ParseInt(vars["storeId"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid store ID", http.StatusBadRequest)
		return
	}

	var cmd commands.ReleaseInventoryCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	cmd.StoreID = storeID

	inventory, err := h.commandHandler.HandleReleaseInventory(r.Context(), cmd)
	if err != nil {
		if err == domain.ErrInventoryNotFound {
			http.Error(w, "Inventory not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(queries.ToStoreInventoryDTO(inventory))
}
