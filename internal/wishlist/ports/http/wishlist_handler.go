package http

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/qhato/ecommerce/internal/wishlist/application/commands"
	"github.com/qhato/ecommerce/internal/wishlist/application/queries"
)

type WishlistHandler struct {
	commandHandler *commands.WishlistCommandHandler
	queryService   *queries.WishlistQueryService
}

func NewWishlistHandler(
	commandHandler *commands.WishlistCommandHandler,
	queryService *queries.WishlistQueryService,
) *WishlistHandler {
	return &WishlistHandler{
		commandHandler: commandHandler,
		queryService:   queryService,
	}
}

func (h *WishlistHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/wishlists", h.CreateWishlist).Methods("POST")
	router.HandleFunc("/wishlists/{id}", h.GetWishlist).Methods("GET")
	router.HandleFunc("/wishlists/{id}", h.UpdateWishlist).Methods("PUT")
	router.HandleFunc("/wishlists/{id}", h.DeleteWishlist).Methods("DELETE")
	router.HandleFunc("/wishlists/{id}/set-default", h.SetDefaultWishlist).Methods("POST")
	router.HandleFunc("/wishlists/customer/{customerId}", h.GetCustomerWishlists).Methods("GET")
	router.HandleFunc("/wishlists/customer/{customerId}/default", h.GetDefaultWishlist).Methods("GET")
	router.HandleFunc("/wishlists/{id}/public", h.GetPublicWishlist).Methods("GET")
	router.HandleFunc("/wishlists/{id}/items", h.AddItem).Methods("POST")
	router.HandleFunc("/wishlists/{id}/items", h.GetWishlistItems).Methods("GET")
	router.HandleFunc("/wishlist-items/{id}", h.GetWishlistItem).Methods("GET")
	router.HandleFunc("/wishlist-items/{id}", h.UpdateItem).Methods("PUT")
	router.HandleFunc("/wishlist-items/{id}", h.RemoveItem).Methods("DELETE")
	router.HandleFunc("/wishlist-items/{id}/move", h.MoveItem).Methods("POST")
}

type CreateWishlistRequest struct {
	CustomerID string `json:"customer_id" validate:"required"`
	Name       string `json:"name"`
	IsDefault  bool   `json:"is_default"`
	IsPublic   bool   `json:"is_public"`
}

type UpdateWishlistRequest struct {
	Name     string `json:"name" validate:"required"`
	IsPublic bool   `json:"is_public"`
}

type SetDefaultWishlistRequest struct {
	CustomerID string `json:"customer_id" validate:"required"`
}

type AddItemRequest struct {
	ProductID string  `json:"product_id" validate:"required"`
	SKUID     *string `json:"sku_id"`
	Quantity  int     `json:"quantity"`
	Priority  int     `json:"priority"`
	Notes     string  `json:"notes"`
}

type UpdateItemRequest struct {
	Quantity int    `json:"quantity"`
	Priority int    `json:"priority"`
	Notes    string `json:"notes"`
}

type MoveItemRequest struct {
	TargetWishlistID string `json:"target_wishlist_id" validate:"required"`
}

func (h *WishlistHandler) CreateWishlist(w http.ResponseWriter, r *http.Request) {
	var req CreateWishlistRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cmd := commands.CreateWishlistCommand{
		CustomerID: req.CustomerID,
		Name:       req.Name,
		IsDefault:  req.IsDefault,
		IsPublic:   req.IsPublic,
	}

	wishlist, err := h.commandHandler.HandleCreateWishlist(r.Context(), cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(wishlist)
}

func (h *WishlistHandler) GetWishlist(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	wishlist, err := h.queryService.GetWishlist(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(wishlist)
}

func (h *WishlistHandler) UpdateWishlist(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req UpdateWishlistRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cmd := commands.UpdateWishlistCommand{
		ID:       id,
		Name:     req.Name,
		IsPublic: req.IsPublic,
	}

	wishlist, err := h.commandHandler.HandleUpdateWishlist(r.Context(), cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(wishlist)
}

func (h *WishlistHandler) DeleteWishlist(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	cmd := commands.DeleteWishlistCommand{ID: id}
	if err := h.commandHandler.HandleDeleteWishlist(r.Context(), cmd); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *WishlistHandler) SetDefaultWishlist(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req SetDefaultWishlistRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cmd := commands.SetDefaultWishlistCommand{
		ID:         id,
		CustomerID: req.CustomerID,
	}

	wishlist, err := h.commandHandler.HandleSetDefaultWishlist(r.Context(), cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(wishlist)
}

func (h *WishlistHandler) GetCustomerWishlists(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customerID := vars["customerId"]

	wishlists, err := h.queryService.GetCustomerWishlists(r.Context(), customerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(wishlists)
}

func (h *WishlistHandler) GetDefaultWishlist(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customerID := vars["customerId"]

	wishlist, err := h.queryService.GetDefaultWishlist(r.Context(), customerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(wishlist)
}

func (h *WishlistHandler) GetPublicWishlist(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	wishlist, err := h.queryService.GetPublicWishlist(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(wishlist)
}

func (h *WishlistHandler) AddItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	wishlistID := vars["id"]

	var req AddItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cmd := commands.AddItemCommand{
		WishlistID: wishlistID,
		ProductID:  req.ProductID,
		SKUID:      req.SKUID,
		Quantity:   req.Quantity,
		Priority:   req.Priority,
		Notes:      req.Notes,
	}

	item, err := h.commandHandler.HandleAddItem(r.Context(), cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

func (h *WishlistHandler) GetWishlistItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	item, err := h.queryService.GetWishlistItem(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

func (h *WishlistHandler) GetWishlistItems(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	wishlistID := vars["id"]

	items, err := h.queryService.GetWishlistItems(r.Context(), wishlistID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func (h *WishlistHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req UpdateItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cmd := commands.UpdateItemCommand{
		ID:       id,
		Quantity: req.Quantity,
		Priority: req.Priority,
		Notes:    req.Notes,
	}

	item, err := h.commandHandler.HandleUpdateItem(r.Context(), cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

func (h *WishlistHandler) RemoveItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	cmd := commands.RemoveItemCommand{ID: id}
	if err := h.commandHandler.HandleRemoveItem(r.Context(), cmd); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *WishlistHandler) MoveItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req MoveItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cmd := commands.MoveItemCommand{
		ItemID:           id,
		TargetWishlistID: req.TargetWishlistID,
	}

	item, err := h.commandHandler.HandleMoveItem(r.Context(), cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}
