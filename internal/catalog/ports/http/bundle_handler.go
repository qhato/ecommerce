package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/qhato/ecommerce/internal/catalog/application/commands"
	"github.com/qhato/ecommerce/internal/catalog/application/queries"
	"github.com/shopspring/decimal"
)

type ProductBundleHandler struct {
	commandHandler *commands.ProductBundleCommandHandler
	queryService   *queries.ProductBundleQueryService
}

func NewProductBundleHandler(
	commandHandler *commands.ProductBundleCommandHandler,
	queryService   *queries.ProductBundleQueryService,
) *ProductBundleHandler {
	return &ProductBundleHandler{
		commandHandler: commandHandler,
		queryService:   queryService,
	}
}

func (h *ProductBundleHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/bundles", h.CreateBundle).Methods("POST")
	router.HandleFunc("/bundles", h.GetAllBundles).Methods("GET")
	router.HandleFunc("/bundles/{id}", h.GetBundle).Methods("GET")
	router.HandleFunc("/bundles/{id}", h.UpdateBundle).Methods("PUT")
	router.HandleFunc("/bundles/{id}", h.DeleteBundle).Methods("DELETE")
	router.HandleFunc("/bundles/{id}/activate", h.ActivateBundle).Methods("POST")
	router.HandleFunc("/bundles/{id}/deactivate", h.DeactivateBundle).Methods("POST")
}

type CreateBundleRequest struct {
	Name        string                      `json:"name" validate:"required"`
	Description string                      `json:"description"`
	BundlePrice decimal.Decimal             `json:"bundle_price" validate:"required"`
	Items       []commands.BundleItemInput  `json:"items" validate:"required,min=1"`
}

type UpdateBundleRequest struct {
	Name        string          `json:"name" validate:"required"`
	Description string          `json:"description"`
	BundlePrice decimal.Decimal `json:"bundle_price" validate:"required"`
}

func (h *ProductBundleHandler) CreateBundle(w http.ResponseWriter, r *http.Request) {
	var req CreateBundleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cmd := commands.CreateProductBundleCommand{
		Name:        req.Name,
		Description: req.Description,
		BundlePrice: req.BundlePrice,
		Items:       req.Items,
	}

	bundle, err := h.commandHandler.HandleCreateProductBundle(r.Context(), cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(bundle)
}

func (h *ProductBundleHandler) GetBundle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	bundle, err := h.queryService.GetBundle(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bundle)
}

func (h *ProductBundleHandler) GetAllBundles(w http.ResponseWriter, r *http.Request) {
	activeOnlyStr := r.URL.Query().Get("active_only")
	activeOnly := activeOnlyStr == "true"

	bundles, err := h.queryService.GetAllBundles(r.Context(), activeOnly)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bundles)
}

func (h *ProductBundleHandler) UpdateBundle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var req UpdateBundleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cmd := commands.UpdateProductBundleCommand{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		BundlePrice: req.BundlePrice,
	}

	bundle, err := h.commandHandler.HandleUpdateProductBundle(r.Context(), cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bundle)
}

func (h *ProductBundleHandler) DeleteBundle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	cmd := commands.DeleteBundleCommand{ID: id}
	if err := h.commandHandler.HandleDeleteBundle(r.Context(), cmd); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ProductBundleHandler) ActivateBundle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	cmd := commands.ActivateBundleCommand{ID: id}
	if err := h.commandHandler.HandleActivateBundle(r.Context(), cmd); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ProductBundleHandler) DeactivateBundle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	cmd := commands.DeactivateBundleCommand{ID: id}
	if err := h.commandHandler.HandleDeactivateBundle(r.Context(), cmd); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
