package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/qhato/ecommerce/internal/catalog/application/commands"
	"github.com/qhato/ecommerce/internal/catalog/application/queries"
)

type ProductRelationshipHandler struct {
	commandHandler *commands.ProductRelationshipCommandHandler
	queryService   *queries.ProductRelationshipQueryService
}

func NewProductRelationshipHandler(
	commandHandler *commands.ProductRelationshipCommandHandler,
	queryService   *queries.ProductRelationshipQueryService,
) *ProductRelationshipHandler {
	return &ProductRelationshipHandler{
		commandHandler: commandHandler,
		queryService:   queryService,
	}
}

func (h *ProductRelationshipHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/product-relationships", h.CreateRelationship).Methods("POST")
	router.HandleFunc("/product-relationships/{id}", h.DeleteRelationship).Methods("DELETE")
	router.HandleFunc("/product-relationships/{id}/activate", h.ActivateRelationship).Methods("POST")
	router.HandleFunc("/product-relationships/{id}/deactivate", h.DeactivateRelationship).Methods("POST")
	router.HandleFunc("/products/{productId}/cross-sell", h.GetCrossSellProducts).Methods("GET")
	router.HandleFunc("/products/{productId}/up-sell", h.GetUpSellProducts).Methods("GET")
	router.HandleFunc("/products/{productId}/related", h.GetRelatedProducts).Methods("GET")
}

type CreateRelationshipRequest struct {
	ProductID        int64  `json:"product_id" validate:"required"`
	RelatedProductID int64  `json:"related_product_id" validate:"required"`
	RelationshipType string `json:"relationship_type" validate:"required"`
	Sequence         int    `json:"sequence"`
}

func (h *ProductRelationshipHandler) CreateRelationship(w http.ResponseWriter, r *http.Request) {
	var req CreateRelationshipRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cmd := commands.CreateProductRelationshipCommand{
		ProductID:        req.ProductID,
		RelatedProductID: req.RelatedProductID,
		RelationshipType: req.RelationshipType,
		Sequence:         req.Sequence,
	}

	relationship, err := h.commandHandler.HandleCreateProductRelationship(r.Context(), cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(relationship)
}

func (h *ProductRelationshipHandler) DeleteRelationship(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	cmd := commands.DeleteRelationshipCommand{ID: id}
	if err := h.commandHandler.HandleDeleteRelationship(r.Context(), cmd); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ProductRelationshipHandler) ActivateRelationship(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	cmd := commands.ActivateRelationshipCommand{ID: id}
	if err := h.commandHandler.HandleActivateRelationship(r.Context(), cmd); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ProductRelationshipHandler) DeactivateRelationship(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	cmd := commands.DeactivateRelationshipCommand{ID: id}
	if err := h.commandHandler.HandleDeactivateRelationship(r.Context(), cmd); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ProductRelationshipHandler) GetCrossSellProducts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID, err := strconv.ParseInt(vars["productId"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	products, err := h.queryService.GetCrossSellProducts(r.Context(), productID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func (h *ProductRelationshipHandler) GetUpSellProducts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID, err := strconv.ParseInt(vars["productId"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	products, err := h.queryService.GetUpSellProducts(r.Context(), productID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func (h *ProductRelationshipHandler) GetRelatedProducts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID, err := strconv.ParseInt(vars["productId"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	products, err := h.queryService.GetRelatedProducts(r.Context(), productID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}
