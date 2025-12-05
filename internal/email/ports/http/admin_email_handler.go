package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/qhato/ecommerce/internal/email/application"
	"github.com/qhato/ecommerce/internal/email/application/commands"
	pkghttp "github.com/qhato/ecommerce/pkg/http"
	"github.com/qhato/ecommerce/pkg/logger"
)

// AdminEmailHandler handles admin email HTTP requests
type AdminEmailHandler struct {
	emailService *application.EmailService
	logger       logger.Logger
}

// NewAdminEmailHandler creates a new admin email handler
func NewAdminEmailHandler(emailService *application.EmailService, logger logger.Logger) *AdminEmailHandler {
	return &AdminEmailHandler{
		emailService: emailService,
		logger:       logger,
	}
}

// RegisterRoutes registers the admin email routes
func (h *AdminEmailHandler) RegisterRoutes(r chi.Router) {
	r.Get("/emails", h.ListEmails)
	r.Get("/emails/{id}", h.GetEmail)
	r.Get("/emails/status/{status}", h.ListEmailsByStatus)
	r.Get("/emails/type/{type}", h.ListEmailsByType)
	r.Get("/emails/order/{orderId}", h.ListEmailsByOrder)
	r.Get("/emails/customer/{customerId}", h.ListEmailsByCustomer)
	r.Get("/emails/stats", h.GetEmailStats)
	r.Post("/emails/send", h.SendEmail)
	r.Post("/emails/schedule", h.ScheduleEmail)
	r.Post("/emails/{id}/cancel", h.CancelEmail)
	r.Post("/emails/{id}/retry", h.RetryEmail)
}

// ListEmails lists all emails
func (h *AdminEmailHandler) ListEmails(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	if limit == 0 {
		limit = 20
	}

	var emails interface{}
	var err error

	if status != "" {
		emails, err = h.emailService.ListEmailsByStatus(r.Context(), status, offset, limit)
	} else {
		// Default to listing pending emails
		emails, err = h.emailService.ListEmailsByStatus(r.Context(), "PENDING", offset, limit)
	}

	if err != nil {
		h.logger.Error("Failed to list emails",
			logger.Field{Key: "error", Value: err.Error()},
		)
		pkghttp.RespondWithError(w, http.StatusInternalServerError, "Failed to list emails")
		return
	}

	pkghttp.RespondWithJSON(w, http.StatusOK, emails)
}

// GetEmail retrieves a single email
func (h *AdminEmailHandler) GetEmail(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		pkghttp.RespondWithError(w, http.StatusBadRequest, "Invalid email ID")
		return
	}

	email, err := h.emailService.GetEmailByID(r.Context(), id)
	if err != nil {
		h.logger.Error("Failed to get email",
			logger.Field{Key: "id", Value: id},
			logger.Field{Key: "error", Value: err.Error()},
		)
		pkghttp.RespondWithError(w, http.StatusNotFound, "Email not found")
		return
	}

	pkghttp.RespondWithJSON(w, http.StatusOK, email)
}

// ListEmailsByStatus lists emails by status
func (h *AdminEmailHandler) ListEmailsByStatus(w http.ResponseWriter, r *http.Request) {
	status := chi.URLParam(r, "status")
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	if limit == 0 {
		limit = 20
	}

	emails, err := h.emailService.ListEmailsByStatus(r.Context(), status, offset, limit)
	if err != nil {
		h.logger.Error("Failed to list emails by status",
			logger.Field{Key: "status", Value: status},
			logger.Field{Key: "error", Value: err.Error()},
		)
		pkghttp.RespondWithError(w, http.StatusInternalServerError, "Failed to list emails")
		return
	}

	pkghttp.RespondWithJSON(w, http.StatusOK, emails)
}

// ListEmailsByType lists emails by type
func (h *AdminEmailHandler) ListEmailsByType(w http.ResponseWriter, r *http.Request) {
	emailType := chi.URLParam(r, "type")
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	if limit == 0 {
		limit = 20
	}

	emails, err := h.emailService.ListEmailsByType(r.Context(), emailType, offset, limit)
	if err != nil {
		h.logger.Error("Failed to list emails by type",
			logger.Field{Key: "type", Value: emailType},
			logger.Field{Key: "error", Value: err.Error()},
		)
		pkghttp.RespondWithError(w, http.StatusInternalServerError, "Failed to list emails")
		return
	}

	pkghttp.RespondWithJSON(w, http.StatusOK, emails)
}

// ListEmailsByOrder lists emails for an order
func (h *AdminEmailHandler) ListEmailsByOrder(w http.ResponseWriter, r *http.Request) {
	orderIDStr := chi.URLParam(r, "orderId")
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		pkghttp.RespondWithError(w, http.StatusBadRequest, "Invalid order ID")
		return
	}

	emails, err := h.emailService.ListEmailsByOrderID(r.Context(), orderID)
	if err != nil {
		h.logger.Error("Failed to list emails by order",
			logger.Field{Key: "order_id", Value: orderID},
			logger.Field{Key: "error", Value: err.Error()},
		)
		pkghttp.RespondWithError(w, http.StatusInternalServerError, "Failed to list emails")
		return
	}

	pkghttp.RespondWithJSON(w, http.StatusOK, emails)
}

// ListEmailsByCustomer lists emails for a customer
func (h *AdminEmailHandler) ListEmailsByCustomer(w http.ResponseWriter, r *http.Request) {
	customerIDStr := chi.URLParam(r, "customerId")
	customerID, err := strconv.ParseInt(customerIDStr, 10, 64)
	if err != nil {
		pkghttp.RespondWithError(w, http.StatusBadRequest, "Invalid customer ID")
		return
	}

	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	if limit == 0 {
		limit = 20
	}

	emails, err := h.emailService.ListEmailsByCustomerID(r.Context(), customerID, offset, limit)
	if err != nil {
		h.logger.Error("Failed to list emails by customer",
			logger.Field{Key: "customer_id", Value: customerID},
			logger.Field{Key: "error", Value: err.Error()},
		)
		pkghttp.RespondWithError(w, http.StatusInternalServerError, "Failed to list emails")
		return
	}

	pkghttp.RespondWithJSON(w, http.StatusOK, emails)
}

// GetEmailStats retrieves email statistics
func (h *AdminEmailHandler) GetEmailStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.emailService.GetEmailStats(r.Context())
	if err != nil {
		h.logger.Error("Failed to get email stats",
			logger.Field{Key: "error", Value: err.Error()},
		)
		pkghttp.RespondWithError(w, http.StatusInternalServerError, "Failed to get email stats")
		return
	}

	pkghttp.RespondWithJSON(w, http.StatusOK, stats)
}

// SendEmailRequest represents a send email request
type SendEmailRequest struct {
	Type         string                 `json:"type"`
	From         string                 `json:"from,omitempty"`
	To           []string               `json:"to"`
	CC           []string               `json:"cc,omitempty"`
	BCC          []string               `json:"bcc,omitempty"`
	ReplyTo      string                 `json:"reply_to,omitempty"`
	Subject      string                 `json:"subject"`
	Body         string                 `json:"body,omitempty"`
	HTMLBody     string                 `json:"html_body,omitempty"`
	TemplateName string                 `json:"template_name,omitempty"`
	TemplateData map[string]interface{} `json:"template_data,omitempty"`
	Priority     int                    `json:"priority,omitempty"`
	OrderID      *int64                 `json:"order_id,omitempty"`
	CustomerID   *int64                 `json:"customer_id,omitempty"`
}

// SendEmail sends an email
func (h *AdminEmailHandler) SendEmail(w http.ResponseWriter, r *http.Request) {
	var req SendEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkghttp.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	cmd := &commands.SendEmailCommand{
		Type:         req.Type,
		From:         req.From,
		To:           req.To,
		CC:           req.CC,
		BCC:          req.BCC,
		ReplyTo:      req.ReplyTo,
		Subject:      req.Subject,
		Body:         req.Body,
		HTMLBody:     req.HTMLBody,
		TemplateName: req.TemplateName,
		TemplateData: req.TemplateData,
		Priority:     req.Priority,
		OrderID:      req.OrderID,
		CustomerID:   req.CustomerID,
	}

	emailID, err := h.emailService.SendEmail(r.Context(), cmd)
	if err != nil {
		h.logger.Error("Failed to send email",
			logger.Field{Key: "error", Value: err.Error()},
		)
		pkghttp.RespondWithError(w, http.StatusInternalServerError, "Failed to send email")
		return
	}

	pkghttp.RespondWithJSON(w, http.StatusCreated, map[string]interface{}{
		"email_id": emailID,
		"message":  "Email queued successfully",
	})
}

// ScheduleEmailRequest represents a schedule email request
type ScheduleEmailRequest struct {
	SendEmailRequest
	ScheduledAt time.Time `json:"scheduled_at"`
}

// ScheduleEmail schedules an email
func (h *AdminEmailHandler) ScheduleEmail(w http.ResponseWriter, r *http.Request) {
	var req ScheduleEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkghttp.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	cmd := &commands.ScheduleEmailCommand{
		SendEmailCommand: commands.SendEmailCommand{
			Type:         req.Type,
			From:         req.From,
			To:           req.To,
			CC:           req.CC,
			BCC:          req.BCC,
			ReplyTo:      req.ReplyTo,
			Subject:      req.Subject,
			Body:         req.Body,
			HTMLBody:     req.HTMLBody,
			TemplateName: req.TemplateName,
			TemplateData: req.TemplateData,
			Priority:     req.Priority,
			OrderID:      req.OrderID,
			CustomerID:   req.CustomerID,
		},
		ScheduledAt: req.ScheduledAt,
	}

	emailID, err := h.emailService.ScheduleEmail(r.Context(), cmd)
	if err != nil {
		h.logger.Error("Failed to schedule email",
			logger.Field{Key: "error", Value: err.Error()},
		)
		pkghttp.RespondWithError(w, http.StatusInternalServerError, "Failed to schedule email")
		return
	}

	pkghttp.RespondWithJSON(w, http.StatusCreated, map[string]interface{}{
		"email_id": emailID,
		"message":  "Email scheduled successfully",
	})
}

// CancelEmail cancels a pending email
func (h *AdminEmailHandler) CancelEmail(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		pkghttp.RespondWithError(w, http.StatusBadRequest, "Invalid email ID")
		return
	}

	cmd := &commands.CancelEmailCommand{
		EmailID: id,
	}

	if err := h.emailService.CancelEmail(r.Context(), cmd); err != nil {
		h.logger.Error("Failed to cancel email",
			logger.Field{Key: "id", Value: id},
			logger.Field{Key: "error", Value: err.Error()},
		)
		pkghttp.RespondWithError(w, http.StatusInternalServerError, "Failed to cancel email")
		return
	}

	pkghttp.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Email cancelled successfully",
	})
}

// RetryEmail retries a failed email
func (h *AdminEmailHandler) RetryEmail(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		pkghttp.RespondWithError(w, http.StatusBadRequest, "Invalid email ID")
		return
	}

	cmd := &commands.RetryFailedEmailCommand{
		EmailID: id,
	}

	if err := h.emailService.RetryFailedEmail(r.Context(), cmd); err != nil {
		h.logger.Error("Failed to retry email",
			logger.Field{Key: "id", Value: id},
			logger.Field{Key: "error", Value: err.Error()},
		)
		pkghttp.RespondWithError(w, http.StatusInternalServerError, "Failed to retry email")
		return
	}

	pkghttp.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Email queued for retry",
	})
}
