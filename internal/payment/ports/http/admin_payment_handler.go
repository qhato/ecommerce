package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/qhato/ecommerce/internal/payment/application"
	"github.com/qhato/ecommerce/internal/payment/application/commands"
	"github.com/qhato/ecommerce/internal/payment/application/queries"
	"github.com/qhato/ecommerce/internal/payment/domain"
	"github.com/qhato/ecommerce/pkg/errors"
	httpPkg "github.com/qhato/ecommerce/pkg/http"
	"github.com/qhato/ecommerce/pkg/logger"
	"github.com/qhato/ecommerce/pkg/validator"
)

// AdminPaymentHandler handles admin payment HTTP requests
type AdminPaymentHandler struct {
	commandHandler *commands.PaymentCommandHandler
	queryHandler   *queries.PaymentQueryHandler
	validator      *validator.Validator
	log            *logger.Logger
}

// NewAdminPaymentHandler creates a new AdminPaymentHandler
func NewAdminPaymentHandler(
	commandHandler *commands.PaymentCommandHandler,
	queryHandler *queries.PaymentQueryHandler,
	validator *validator.Validator,
	log *logger.Logger,
) *AdminPaymentHandler {
	return &AdminPaymentHandler{
		commandHandler: commandHandler,
		queryHandler:   queryHandler,
		validator:      validator,
		log:            log,
	}
}

// RegisterRoutes registers payment routes
func (h *AdminPaymentHandler) RegisterRoutes(r chi.Router) {
	r.Route("/payments", func(r chi.Router) {
		r.Post("/", h.CreatePayment)
		r.Get("/", h.ListPayments)
		r.Get("/{id}", h.GetPayment)
		r.Post("/{id}/authorize", h.AuthorizePayment)
		r.Post("/{id}/capture", h.CapturePayment)
		r.Post("/{id}/complete", h.CompletePayment)
		r.Post("/{id}/fail", h.FailPayment)
		r.Post("/{id}/refund", h.RefundPayment)
		r.Post("/{id}/cancel", h.CancelPayment)
		r.Get("/order/{orderId}", h.GetPaymentsByOrder)
		r.Get("/transaction/{transactionId}", h.GetPaymentByTransaction)
	})
}

// CreatePayment creates a new payment
func (h *AdminPaymentHandler) CreatePayment(w http.ResponseWriter, r *http.Request) {
	var req application.CreatePaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpPkg.RespondError(w, errors.BadRequest("invalid request body").WithInternal(err))
		return
	}

	if err := h.validator.Validate(req); err != nil {
		httpPkg.RespondError(w, errors.BadRequest("validation failed").WithInternal(err))
		return
	}

	payment, err := h.commandHandler.CreatePayment(
		r.Context(),
		req.OrderID,
		req.CustomerID,
		domain.PaymentMethod(req.PaymentMethod),
		req.Amount,
		req.CurrencyCode,
	)
	if err != nil {
		httpPkg.RespondError(w, errors.InternalWrap(err, "failed to create payment"))
		return
	}

	httpPkg.RespondJSON(w, http.StatusCreated, application.ToPaymentDTO(payment))
}

// GetPayment retrieves a payment by ID
func (h *AdminPaymentHandler) GetPayment(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpPkg.RespondError(w, errors.BadRequest("invalid payment ID").WithInternal(err))
		return
	}

	payment, err := h.queryHandler.GetByID(r.Context(), id)
	if err != nil {
		httpPkg.RespondError(w, errors.NotFound("payment not found").WithInternal(err))
		return
	}

	httpPkg.RespondJSON(w, http.StatusOK, application.ToPaymentDTO(payment))
}

// GetPaymentByTransaction retrieves a payment by transaction ID
func (h *AdminPaymentHandler) GetPaymentByTransaction(w http.ResponseWriter, r *http.Request) {
	transactionID := chi.URLParam(r, "transactionId")

	payment, err := h.queryHandler.GetByTransactionID(r.Context(), transactionID)
	if err != nil {
		httpPkg.RespondError(w, errors.NotFound("payment not found").WithInternal(err))
		return
	}

	httpPkg.RespondJSON(w, http.StatusOK, application.ToPaymentDTO(payment))
}

// GetPaymentsByOrder retrieves payments for an order
func (h *AdminPaymentHandler) GetPaymentsByOrder(w http.ResponseWriter, r *http.Request) {
	orderIDStr := chi.URLParam(r, "orderId")
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		httpPkg.RespondError(w, errors.BadRequest("invalid order ID").WithInternal(err))
		return
	}

	payments, err := h.queryHandler.ListByOrder(r.Context(), orderID)
	if err != nil {
		httpPkg.RespondError(w, errors.InternalWrap(err, "failed to list payments"))
		return
	}

	httpPkg.RespondJSON(w, http.StatusOK, application.ToPaymentDTOs(payments))
}

// ListPayments lists all payments
func (h *AdminPaymentHandler) ListPayments(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	paymentMethod := r.URL.Query().Get("payment_method")
	customerIDStr := r.URL.Query().Get("customer_id")
	orderIDStr := r.URL.Query().Get("order_id")
	sortBy := r.URL.Query().Get("sort_by")
	sortOrder := r.URL.Query().Get("sort_order")

	var customerID, orderID int64
	if customerIDStr != "" {
		customerID, _ = strconv.ParseInt(customerIDStr, 10, 64)
	}
	if orderIDStr != "" {
		orderID, _ = strconv.ParseInt(orderIDStr, 10, 64)
	}

	filter := &domain.PaymentFilter{
		Page:          page,
		PageSize:      pageSize,
		PaymentMethod: domain.PaymentMethod(paymentMethod),
		CustomerID:    customerID,
		OrderID:       orderID,
		SortBy:        sortBy,
		SortOrder:     sortOrder,
	}

	payments, total, err := h.queryHandler.List(r.Context(), filter)
	if err != nil {
		httpPkg.RespondError(w, errors.InternalWrap(err, "failed to list payments"))
		return
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	response := application.PaginatedPaymentResponse{
		Data:       application.ToPaymentDTOs(payments),
		Page:       page,
		PageSize:   pageSize,
		TotalItems: total,
		TotalPages: totalPages,
	}

	httpPkg.RespondJSON(w, http.StatusOK, response)
}

// AuthorizePayment authorizes a payment
func (h *AdminPaymentHandler) AuthorizePayment(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpPkg.RespondError(w, errors.BadRequest("invalid payment ID").WithInternal(err))
		return
	}

	var req application.AuthorizePaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpPkg.RespondError(w, errors.BadRequest("invalid request body").WithInternal(err))
		return
	}

	if err := h.validator.Validate(req); err != nil {
		httpPkg.RespondError(w, errors.BadRequest("validation failed").WithInternal(err))
		return
	}

	err = h.commandHandler.AuthorizePayment(r.Context(), id, req.AuthorizationCode, req.TransactionID)
	if err != nil {
		httpPkg.RespondError(w, errors.InternalWrap(err, "failed to authorize payment"))
		return
	}

	payment, _ := h.queryHandler.GetByID(r.Context(), id)
	if payment != nil {
		h.queryHandler.InvalidateCache(r.Context(), payment.ID, payment.TransactionID)
	}

	httpPkg.RespondJSON(w, http.StatusOK, map[string]string{"message": "payment authorized successfully"})
}

// CapturePayment captures an authorized payment
func (h *AdminPaymentHandler) CapturePayment(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpPkg.RespondError(w, errors.BadRequest("invalid payment ID").WithInternal(err))
		return
	}

	var req application.CapturePaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpPkg.RespondError(w, errors.BadRequest("invalid request body").WithInternal(err))
		return
	}

	err = h.commandHandler.CapturePayment(r.Context(), id, req.TransactionID)
	if err != nil {
		httpPkg.RespondError(w, errors.InternalWrap(err, "failed to capture payment"))
		return
	}

	payment, _ := h.queryHandler.GetByID(r.Context(), id)
	if payment != nil {
		h.queryHandler.InvalidateCache(r.Context(), payment.ID, payment.TransactionID)
	}

	httpPkg.RespondJSON(w, http.StatusOK, map[string]string{"message": "payment captured successfully"})
}

// CompletePayment completes a payment
func (h *AdminPaymentHandler) CompletePayment(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpPkg.RespondError(w, errors.BadRequest("invalid payment ID").WithInternal(err))
		return
	}

	var req application.ProcessPaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpPkg.RespondError(w, errors.BadRequest("invalid request body").WithInternal(err))
		return
	}

	if err := h.validator.Validate(req); err != nil {
		httpPkg.RespondError(w, errors.BadRequest("validation failed").WithInternal(err))
		return
	}

	err = h.commandHandler.CompletePayment(r.Context(), id, req.TransactionID)
	if err != nil {
		httpPkg.RespondError(w, errors.InternalWrap(err, "failed to complete payment"))
		return
	}

	payment, _ := h.queryHandler.GetByID(r.Context(), id)
	if payment != nil {
		h.queryHandler.InvalidateCache(r.Context(), payment.ID, payment.TransactionID)
	}

	httpPkg.RespondJSON(w, http.StatusOK, map[string]string{"message": "payment completed successfully"})
}

// FailPayment marks a payment as failed
func (h *AdminPaymentHandler) FailPayment(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpPkg.RespondError(w, errors.BadRequest("invalid payment ID").WithInternal(err))
		return
	}

	var req struct {
		Reason string `json:"reason" validate:"required"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpPkg.RespondError(w, errors.BadRequest("invalid request body").WithInternal(err))
		return
	}

	if err := h.validator.Validate(req); err != nil {
		httpPkg.RespondError(w, errors.BadRequest("validation failed").WithInternal(err))
		return
	}

	err = h.commandHandler.FailPayment(r.Context(), id, req.Reason)
	if err != nil {
		httpPkg.RespondError(w, errors.InternalWrap(err, "failed to mark payment as failed"))
		return
	}

	payment, _ := h.queryHandler.GetByID(r.Context(), id)
	if payment != nil {
		h.queryHandler.InvalidateCache(r.Context(), payment.ID, payment.TransactionID)
	}

	httpPkg.RespondJSON(w, http.StatusOK, map[string]string{"message": "payment marked as failed successfully"})
}

// RefundPayment refunds a payment
func (h *AdminPaymentHandler) RefundPayment(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpPkg.RespondError(w, errors.BadRequest("invalid payment ID").WithInternal(err))
		return
	}

	var req application.RefundPaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpPkg.RespondError(w, errors.BadRequest("invalid request body").WithInternal(err))
		return
	}

	if err := h.validator.Validate(req); err != nil {
		httpPkg.RespondError(w, errors.BadRequest("validation failed").WithInternal(err))
		return
	}

	err = h.commandHandler.RefundPayment(r.Context(), id, req.Amount)
	if err != nil {
		httpPkg.RespondError(w, errors.InternalWrap(err, "failed to refund payment"))
		return
	}

	payment, _ := h.queryHandler.GetByID(r.Context(), id)
	if payment != nil {
		h.queryHandler.InvalidateCache(r.Context(), payment.ID, payment.TransactionID)
	}

	httpPkg.RespondJSON(w, http.StatusOK, map[string]string{"message": "payment refunded successfully"})
}

// CancelPayment cancels a payment
func (h *AdminPaymentHandler) CancelPayment(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpPkg.RespondError(w, errors.BadRequest("invalid payment ID").WithInternal(err))
		return
	}

	err = h.commandHandler.CancelPayment(r.Context(), id)
	if err != nil {
		httpPkg.RespondError(w, errors.InternalWrap(err, "failed to cancel payment"))
		return
	}

	payment, _ := h.queryHandler.GetByID(r.Context(), id)
	if payment != nil {
		h.queryHandler.InvalidateCache(r.Context(), payment.ID, payment.TransactionID)
	}

	httpPkg.RespondJSON(w, http.StatusOK, map[string]string{"message": "payment cancelled successfully"})
}
