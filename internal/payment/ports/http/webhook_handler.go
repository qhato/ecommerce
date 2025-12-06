package http

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/qhato/ecommerce/internal/payment/application/commands"
	"github.com/qhato/ecommerce/internal/payment/application/queries"
)

type WebhookHandler struct {
	commandHandler *commands.WebhookCommandHandler
	queryService   *queries.WebhookQueryService
}

func NewWebhookHandler(
	commandHandler *commands.WebhookCommandHandler,
	queryService *queries.WebhookQueryService,
) *WebhookHandler {
	return &WebhookHandler{
		commandHandler: commandHandler,
		queryService:   queryService,
	}
}

func (h *WebhookHandler) RegisterRoutes(router *mux.Router) {
	// Webhook endpoints for each gateway
	router.HandleFunc("/webhooks/stripe", h.HandleStripeWebhook).Methods("POST")
	router.HandleFunc("/webhooks/paypal", h.HandlePayPalWebhook).Methods("POST")
	router.HandleFunc("/webhooks/authorizenet", h.HandleAuthorizeNetWebhook).Methods("POST")

	// Admin endpoints
	router.HandleFunc("/admin/webhooks/{id}", h.GetWebhook).Methods("GET")
	router.HandleFunc("/admin/webhooks/pending", h.GetPendingWebhooks).Methods("GET")
	router.HandleFunc("/admin/webhooks/status/{status}", h.GetWebhooksByStatus).Methods("GET")
	router.HandleFunc("/admin/webhooks/{id}/retry", h.RetryWebhook).Methods("POST")
}

func (h *WebhookHandler) HandleStripeWebhook(w http.ResponseWriter, r *http.Request) {
	h.handleGenericWebhook(w, r, "Stripe")
}

func (h *WebhookHandler) HandlePayPalWebhook(w http.ResponseWriter, r *http.Request) {
	h.handleGenericWebhook(w, r, "PayPal")
}

func (h *WebhookHandler) HandleAuthorizeNetWebhook(w http.ResponseWriter, r *http.Request) {
	h.handleGenericWebhook(w, r, "Authorize.Net")
}

func (h *WebhookHandler) handleGenericWebhook(w http.ResponseWriter, r *http.Request, gatewayName string) {
	// Read raw body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// Get signature header (varies by gateway)
	signature := r.Header.Get("Stripe-Signature")
	if signature == "" {
		signature = r.Header.Get("PayPal-Transmission-Sig")
	}
	if signature == "" {
		signature = r.Header.Get("X-ANET-Signature")
	}

	// Get client IP
	ip := r.RemoteAddr

	// Parse payload to extract event ID and type
	// This is a simplified version - in reality, you'd parse the specific gateway format
	var webhookData map[string]interface{}
	if err := json.Unmarshal(body, &webhookData); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// Extract event ID and type (varies by gateway)
	eventID, _ := webhookData["id"].(string)
	if eventID == "" {
		eventID, _ = webhookData["event_id"].(string)
	}

	eventType, _ := webhookData["type"].(string)
	if eventType == "" {
		eventType, _ = webhookData["event_type"].(string)
	}

	if eventID == "" || eventType == "" {
		http.Error(w, "Missing event ID or type", http.StatusBadRequest)
		return
	}

	cmd := commands.ProcessWebhookCommand{
		GatewayName: gatewayName,
		EventID:     eventID,
		EventType:   eventType,
		Payload:     string(body),
		Signature:   &signature,
		IPAddress:   &ip,
	}

	event, err := h.commandHandler.HandleProcessWebhook(r.Context(), cmd)
	if err != nil {
		// Log error but return 200 to prevent retries for invalid webhooks
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":     "success",
		"webhook_id": event.ID,
	})
}

func (h *WebhookHandler) GetWebhook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	webhook, err := h.queryService.GetWebhook(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(webhook)
}

func (h *WebhookHandler) GetPendingWebhooks(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit := 100
	if limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	webhooks, err := h.queryService.GetPendingWebhooks(r.Context(), limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(webhooks)
}

func (h *WebhookHandler) GetWebhooksByStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	status := vars["status"]

	limitStr := r.URL.Query().Get("limit")
	limit := 100
	if limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	webhooks, err := h.queryService.GetWebhooksByStatus(r.Context(), status, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(webhooks)
}

func (h *WebhookHandler) RetryWebhook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	cmd := commands.RetryWebhookCommand{WebhookID: id}
	event, err := h.commandHandler.HandleRetryWebhook(r.Context(), cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(event)
}
