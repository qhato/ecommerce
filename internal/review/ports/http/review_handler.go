package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/qhato/ecommerce/internal/review/application/commands"
	"github.com/qhato/ecommerce/internal/review/application/queries"
)

type ReviewHandler struct {
	commandHandler *commands.ReviewCommandHandler
	queryService   *queries.ReviewQueryService
}

func NewReviewHandler(
	commandHandler *commands.ReviewCommandHandler,
	queryService   *queries.ReviewQueryService,
) *ReviewHandler {
	return &ReviewHandler{
		commandHandler: commandHandler,
		queryService:   queryService,
	}
}

func (h *ReviewHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/reviews", h.CreateReview).Methods("POST")
	router.HandleFunc("/reviews/{id}", h.GetReview).Methods("GET")
	router.HandleFunc("/reviews/{id}", h.UpdateReview).Methods("PUT")
	router.HandleFunc("/reviews/{id}", h.DeleteReview).Methods("DELETE")
	router.HandleFunc("/reviews/{id}/approve", h.ApproveReview).Methods("POST")
	router.HandleFunc("/reviews/{id}/reject", h.RejectReview).Methods("POST")
	router.HandleFunc("/reviews/{id}/flag", h.FlagReview).Methods("POST")
	router.HandleFunc("/reviews/{id}/response", h.AddResponse).Methods("POST")
	router.HandleFunc("/reviews/{id}/helpful", h.MarkHelpful).Methods("POST")
	router.HandleFunc("/reviews/{id}/not-helpful", h.MarkNotHelpful).Methods("POST")
	router.HandleFunc("/reviews/product/{productId}", h.GetProductReviews).Methods("GET")
	router.HandleFunc("/reviews/product/{productId}/summary", h.GetRatingSummary).Methods("GET")
	router.HandleFunc("/reviews/customer/{customerId}", h.GetCustomerReviews).Methods("GET")
	router.HandleFunc("/reviews/pending", h.GetPendingReviews).Methods("GET")
}

type CreateReviewRequest struct {
	ProductID     string  `json:"product_id" validate:"required"`
	CustomerID    string  `json:"customer_id" validate:"required"`
	CustomerName  string  `json:"customer_name" validate:"required"`
	ReviewerEmail string  `json:"reviewer_email" validate:"required,email"`
	Rating        int     `json:"rating" validate:"required,min=1,max=5"`
	Title         string  `json:"title"`
	Comment       string  `json:"comment" validate:"required"`
	OrderID       *string `json:"order_id"`
}

type UpdateReviewRequest struct {
	Title   string `json:"title"`
	Comment string `json:"comment" validate:"required"`
	Rating  int    `json:"rating" validate:"required,min=1,max=5"`
}

type AddResponseRequest struct {
	ResponseText string `json:"response_text" validate:"required"`
}

func (h *ReviewHandler) CreateReview(w http.ResponseWriter, r *http.Request) {
	var req CreateReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cmd := commands.CreateReviewCommand{
		ProductID:     req.ProductID,
		CustomerID:    req.CustomerID,
		CustomerName:  req.CustomerName,
		ReviewerEmail: req.ReviewerEmail,
		Rating:        req.Rating,
		Title:         req.Title,
		Comment:       req.Comment,
		OrderID:       req.OrderID,
	}

	review, err := h.commandHandler.HandleCreateReview(r.Context(), cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(review)
}

func (h *ReviewHandler) GetReview(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	review, err := h.queryService.GetReview(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(review)
}

func (h *ReviewHandler) UpdateReview(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req UpdateReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cmd := commands.UpdateReviewCommand{
		ID:      id,
		Title:   req.Title,
		Comment: req.Comment,
		Rating:  req.Rating,
	}

	review, err := h.commandHandler.HandleUpdateReview(r.Context(), cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(review)
}

func (h *ReviewHandler) DeleteReview(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	cmd := commands.DeleteReviewCommand{ID: id}
	if err := h.commandHandler.HandleDeleteReview(r.Context(), cmd); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ReviewHandler) ApproveReview(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	cmd := commands.ApproveReviewCommand{ID: id}
	review, err := h.commandHandler.HandleApproveReview(r.Context(), cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(review)
}

func (h *ReviewHandler) RejectReview(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	cmd := commands.RejectReviewCommand{ID: id}
	review, err := h.commandHandler.HandleRejectReview(r.Context(), cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(review)
}

func (h *ReviewHandler) FlagReview(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	cmd := commands.FlagReviewCommand{ID: id}
	review, err := h.commandHandler.HandleFlagReview(r.Context(), cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(review)
}

func (h *ReviewHandler) AddResponse(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req AddResponseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cmd := commands.AddResponseCommand{
		ID:           id,
		ResponseText: req.ResponseText,
	}

	review, err := h.commandHandler.HandleAddResponse(r.Context(), cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(review)
}

func (h *ReviewHandler) MarkHelpful(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	cmd := commands.MarkHelpfulCommand{ID: id}
	review, err := h.commandHandler.HandleMarkHelpful(r.Context(), cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(review)
}

func (h *ReviewHandler) MarkNotHelpful(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	cmd := commands.MarkNotHelpfulCommand{ID: id}
	review, err := h.commandHandler.HandleMarkNotHelpful(r.Context(), cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(review)
}

func (h *ReviewHandler) GetProductReviews(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["productId"]

	status := r.URL.Query().Get("status")
	var statusPtr *string
	if status != "" {
		statusPtr = &status
	}

	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 20
	offset := 0

	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	if offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	reviews, err := h.queryService.GetProductReviews(r.Context(), productID, statusPtr, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reviews)
}

func (h *ReviewHandler) GetRatingSummary(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["productId"]

	summary, err := h.queryService.GetProductRatingSummary(r.Context(), productID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
}

func (h *ReviewHandler) GetCustomerReviews(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customerID := vars["customerId"]

	reviews, err := h.queryService.GetCustomerReviews(r.Context(), customerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reviews)
}

func (h *ReviewHandler) GetPendingReviews(w http.ResponseWriter, r *http.Request) {
	reviews, err := h.queryService.GetPendingReviews(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reviews)
}
