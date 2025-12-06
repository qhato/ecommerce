package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/qhato/ecommerce/internal/notification/application/commands"
	"github.com/qhato/ecommerce/internal/notification/application/queries"
	"github.com/qhato/ecommerce/internal/notification/domain"
)

type NotificationHandler struct {
	commandHandler *commands.NotificationCommandHandler
	queryService   *queries.NotificationQueryService
}

func NewNotificationHandler(
	commandHandler *commands.NotificationCommandHandler,
	queryService *queries.NotificationQueryService,
) *NotificationHandler {
	return &NotificationHandler{
		commandHandler: commandHandler,
		queryService:   queryService,
	}
}

func (h *NotificationHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/notifications", h.CreateNotification).Methods("POST")
	router.HandleFunc("/notifications/{id}", h.GetNotification).Methods("GET")
	router.HandleFunc("/notifications/recipient/{recipientId}", h.GetNotificationsByRecipient).Methods("GET")
	router.HandleFunc("/notifications/status/{status}", h.GetNotificationsByStatus).Methods("GET")
	router.HandleFunc("/notifications/pending", h.GetPendingNotifications).Methods("GET")
	router.HandleFunc("/notifications/scheduled", h.GetScheduledNotifications).Methods("GET")
	router.HandleFunc("/notifications/failed", h.GetFailedNotifications).Methods("GET")
	router.HandleFunc("/notifications/{id}/send", h.SendNotification).Methods("POST")
	router.HandleFunc("/notifications/{id}/sent", h.MarkAsSent).Methods("POST")
	router.HandleFunc("/notifications/{id}/failed", h.MarkAsFailed).Methods("POST")
	router.HandleFunc("/notifications/{id}/retry", h.RetryNotification).Methods("POST")
	router.HandleFunc("/notifications/{id}/cancel", h.CancelNotification).Methods("POST")
	router.HandleFunc("/notifications/{id}", h.DeleteNotification).Methods("DELETE")
}

func (h *NotificationHandler) CreateNotification(w http.ResponseWriter, r *http.Request) {
	var cmd commands.CreateNotificationCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	notification, err := h.commandHandler.HandleCreateNotification(r.Context(), cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(queries.ToNotificationDTO(notification))
}

func (h *NotificationHandler) GetNotification(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid notification ID", http.StatusBadRequest)
		return
	}

	notification, err := h.queryService.GetNotification(r.Context(), id)
	if err != nil {
		if err == domain.ErrNotificationNotFound {
			http.Error(w, "Notification not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notification)
}

func (h *NotificationHandler) GetNotificationsByRecipient(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	recipientID := vars["recipientId"]

	limitStr := r.URL.Query().Get("limit")
	limit := 50
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil {
			limit = parsedLimit
		}
	}

	notifications, err := h.queryService.GetNotificationsByRecipient(r.Context(), recipientID, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notifications)
}

func (h *NotificationHandler) GetNotificationsByStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	status := vars["status"]

	limitStr := r.URL.Query().Get("limit")
	limit := 50
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil {
			limit = parsedLimit
		}
	}

	notifications, err := h.queryService.GetNotificationsByStatus(r.Context(), status, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notifications)
}

func (h *NotificationHandler) GetPendingNotifications(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit := 50
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil {
			limit = parsedLimit
		}
	}

	notifications, err := h.queryService.GetPendingNotifications(r.Context(), limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notifications)
}

func (h *NotificationHandler) GetScheduledNotifications(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit := 50
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil {
			limit = parsedLimit
		}
	}

	notifications, err := h.queryService.GetScheduledNotifications(r.Context(), limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notifications)
}

func (h *NotificationHandler) GetFailedNotifications(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit := 50
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil {
			limit = parsedLimit
		}
	}

	notifications, err := h.queryService.GetFailedNotifications(r.Context(), limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notifications)
}

func (h *NotificationHandler) SendNotification(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid notification ID", http.StatusBadRequest)
		return
	}

	cmd := commands.SendNotificationCommand{ID: id}
	notification, err := h.commandHandler.HandleSendNotification(r.Context(), cmd)
	if err != nil {
		if err == domain.ErrNotificationNotFound {
			http.Error(w, "Notification not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(queries.ToNotificationDTO(notification))
}

func (h *NotificationHandler) MarkAsSent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid notification ID", http.StatusBadRequest)
		return
	}

	cmd := commands.MarkAsSentCommand{ID: id}
	notification, err := h.commandHandler.HandleMarkAsSent(r.Context(), cmd)
	if err != nil {
		if err == domain.ErrNotificationNotFound {
			http.Error(w, "Notification not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(queries.ToNotificationDTO(notification))
}

func (h *NotificationHandler) MarkAsFailed(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid notification ID", http.StatusBadRequest)
		return
	}

	var cmd commands.MarkAsFailedCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	cmd.ID = id

	notification, err := h.commandHandler.HandleMarkAsFailed(r.Context(), cmd)
	if err != nil {
		if err == domain.ErrNotificationNotFound {
			http.Error(w, "Notification not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(queries.ToNotificationDTO(notification))
}

func (h *NotificationHandler) RetryNotification(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid notification ID", http.StatusBadRequest)
		return
	}

	cmd := commands.RetryNotificationCommand{ID: id}
	notification, err := h.commandHandler.HandleRetryNotification(r.Context(), cmd)
	if err != nil {
		if err == domain.ErrNotificationNotFound {
			http.Error(w, "Notification not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(queries.ToNotificationDTO(notification))
}

func (h *NotificationHandler) CancelNotification(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid notification ID", http.StatusBadRequest)
		return
	}

	cmd := commands.CancelNotificationCommand{ID: id}
	notification, err := h.commandHandler.HandleCancelNotification(r.Context(), cmd)
	if err != nil {
		if err == domain.ErrNotificationNotFound {
			http.Error(w, "Notification not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(queries.ToNotificationDTO(notification))
}

func (h *NotificationHandler) DeleteNotification(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid notification ID", http.StatusBadRequest)
		return
	}

	cmd := commands.DeleteNotificationCommand{ID: id}
	if err := h.commandHandler.HandleDeleteNotification(r.Context(), cmd); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
