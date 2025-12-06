package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/qhato/ecommerce/internal/promotionmsg/application/commands"
	"github.com/qhato/ecommerce/internal/promotionmsg/application/queries"
)

type PromotionMessageHandler struct {
	commandHandler *commands.PromotionMessageCommandHandler
	queryService   *queries.PromotionMessageQueryService
}

func NewPromotionMessageHandler(
	commandHandler *commands.PromotionMessageCommandHandler,
	queryService *queries.PromotionMessageQueryService,
) *PromotionMessageHandler {
	return &PromotionMessageHandler{
		commandHandler: commandHandler,
		queryService:   queryService,
	}
}

func (h *PromotionMessageHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/promotion-messages", h.CreatePromotionMessage).Methods("POST")
	router.HandleFunc("/promotion-messages/{id}", h.GetMessage).Methods("GET")
	router.HandleFunc("/promotion-messages/{id}", h.UpdatePromotionMessage).Methods("PUT")
	router.HandleFunc("/promotion-messages/{id}", h.DeleteMessage).Methods("DELETE")
	router.HandleFunc("/promotion-messages/type/{type}", h.GetMessagesByType).Methods("GET")
	router.HandleFunc("/promotion-messages/status/{status}", h.GetMessagesByStatus).Methods("GET")
	router.HandleFunc("/promotion-messages/active", h.GetActiveMessages).Methods("GET")
	router.HandleFunc("/promotion-messages/placement/{placement}", h.GetMessagesByPlacement).Methods("GET")
	router.HandleFunc("/promotion-messages/event/{event}", h.GetMessagesByEvent).Methods("POST")
	router.HandleFunc("/promotion-messages/match", h.GetMatchingMessages).Methods("POST")
	router.HandleFunc("/promotion-messages/{id}/activate", h.ActivateMessage).Methods("POST")
	router.HandleFunc("/promotion-messages/{id}/deactivate", h.DeactivateMessage).Methods("POST")
	router.HandleFunc("/promotion-messages/{id}/view", h.IncrementView).Methods("POST")
	router.HandleFunc("/promotion-messages/{id}/click", h.IncrementClick).Methods("POST")
}

func (h *PromotionMessageHandler) CreatePromotionMessage(w http.ResponseWriter, r *http.Request) {
	var cmd commands.CreatePromotionMessageCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	message, err := h.commandHandler.HandleCreatePromotionMessage(r.Context(), cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(message)
}

func (h *PromotionMessageHandler) GetMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	message, err := h.queryService.GetMessage(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(message)
}

func (h *PromotionMessageHandler) UpdatePromotionMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var cmd commands.UpdatePromotionMessageCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	cmd.ID = id

	message, err := h.commandHandler.HandleUpdatePromotionMessage(r.Context(), cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(message)
}

func (h *PromotionMessageHandler) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	cmd := commands.DeleteMessageCommand{ID: id}
	if err := h.commandHandler.HandleDeleteMessage(r.Context(), cmd); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *PromotionMessageHandler) GetMessagesByType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	messageType := vars["type"]

	messages, err := h.queryService.GetMessagesByType(r.Context(), messageType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

func (h *PromotionMessageHandler) GetMessagesByStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	status := vars["status"]

	messages, err := h.queryService.GetMessagesByStatus(r.Context(), status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

func (h *PromotionMessageHandler) GetActiveMessages(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit := 10
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	messages, err := h.queryService.GetActiveMessages(r.Context(), limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

func (h *PromotionMessageHandler) GetMessagesByPlacement(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	placement := vars["placement"]

	messages, err := h.queryService.GetMessagesByPlacement(r.Context(), placement)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

func (h *PromotionMessageHandler) GetMessagesByEvent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	event := vars["event"]

	var req struct {
		Context map[string]interface{} `json:"context"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	messages, err := h.queryService.GetMessagesByEvent(r.Context(), event, req.Context)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

func (h *PromotionMessageHandler) GetMatchingMessages(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Placement string                 `json:"placement"`
		Context   map[string]interface{} `json:"context"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	messages, err := h.queryService.GetMatchingMessages(r.Context(), req.Placement, req.Context)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

func (h *PromotionMessageHandler) ActivateMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	cmd := commands.ActivateMessageCommand{ID: id}
	message, err := h.commandHandler.HandleActivateMessage(r.Context(), cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(message)
}

func (h *PromotionMessageHandler) DeactivateMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	cmd := commands.DeactivateMessageCommand{ID: id}
	message, err := h.commandHandler.HandleDeactivateMessage(r.Context(), cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(message)
}

func (h *PromotionMessageHandler) IncrementView(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	cmd := commands.IncrementViewCommand{ID: id}
	if err := h.commandHandler.HandleIncrementView(r.Context(), cmd); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *PromotionMessageHandler) IncrementClick(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	cmd := commands.IncrementClickCommand{ID: id}
	if err := h.commandHandler.HandleIncrementClick(r.Context(), cmd); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
