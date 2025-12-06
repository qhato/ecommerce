package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/qhato/ecommerce/internal/rule/application/commands"
	"github.com/qhato/ecommerce/internal/rule/application/queries"
	"github.com/qhato/ecommerce/internal/rule/domain"
)

type RuleHandler struct {
	commandHandler *commands.RuleCommandHandler
	queryService   *queries.RuleQueryService
}

func NewRuleHandler(
	commandHandler *commands.RuleCommandHandler,
	queryService *queries.RuleQueryService,
) *RuleHandler {
	return &RuleHandler{
		commandHandler: commandHandler,
		queryService:   queryService,
	}
}

func (h *RuleHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/rules", h.CreateRule).Methods("POST")
	router.HandleFunc("/rules", h.GetAllRules).Methods("GET")
	router.HandleFunc("/rules/{id}", h.GetRule).Methods("GET")
	router.HandleFunc("/rules/{id}", h.UpdateRule).Methods("PUT")
	router.HandleFunc("/rules/{id}", h.DeleteRule).Methods("DELETE")
	router.HandleFunc("/rules/{id}/activate", h.ActivateRule).Methods("POST")
	router.HandleFunc("/rules/{id}/deactivate", h.DeactivateRule).Methods("POST")
	router.HandleFunc("/rules/type/{type}", h.GetRulesByType).Methods("GET")
	router.HandleFunc("/rules/evaluate", h.EvaluateRules).Methods("POST")
}

func (h *RuleHandler) CreateRule(w http.ResponseWriter, r *http.Request) {
	var cmd commands.CreateRuleCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	rule, err := h.commandHandler.HandleCreateRule(r.Context(), cmd)
	if err != nil {
		if err == domain.ErrRuleNameTaken {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(rule)
}

func (h *RuleHandler) GetAllRules(w http.ResponseWriter, r *http.Request) {
	activeOnly := r.URL.Query().Get("active_only") == "true"

	rules, err := h.queryService.GetAllRules(r.Context(), activeOnly)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rules)
}

func (h *RuleHandler) GetRule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid rule ID", http.StatusBadRequest)
		return
	}

	rule, err := h.queryService.GetRule(r.Context(), id)
	if err != nil {
		if err == domain.ErrRuleNotFound {
			http.Error(w, "Rule not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rule)
}

func (h *RuleHandler) UpdateRule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid rule ID", http.StatusBadRequest)
		return
	}

	var cmd commands.UpdateRuleCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	cmd.ID = id

	rule, err := h.commandHandler.HandleUpdateRule(r.Context(), cmd)
	if err != nil {
		if err == domain.ErrRuleNotFound {
			http.Error(w, "Rule not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rule)
}

func (h *RuleHandler) DeleteRule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid rule ID", http.StatusBadRequest)
		return
	}

	cmd := commands.DeleteRuleCommand{ID: id}
	if err := h.commandHandler.HandleDeleteRule(r.Context(), cmd); err != nil {
		if err == domain.ErrRuleNotFound {
			http.Error(w, "Rule not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *RuleHandler) ActivateRule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid rule ID", http.StatusBadRequest)
		return
	}

	cmd := commands.ActivateRuleCommand{ID: id}
	rule, err := h.commandHandler.HandleActivateRule(r.Context(), cmd)
	if err != nil {
		if err == domain.ErrRuleNotFound {
			http.Error(w, "Rule not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rule)
}

func (h *RuleHandler) DeactivateRule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid rule ID", http.StatusBadRequest)
		return
	}

	cmd := commands.DeactivateRuleCommand{ID: id}
	rule, err := h.commandHandler.HandleDeactivateRule(r.Context(), cmd)
	if err != nil {
		if err == domain.ErrRuleNotFound {
			http.Error(w, "Rule not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rule)
}

func (h *RuleHandler) GetRulesByType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ruleType := vars["type"]
	activeOnly := r.URL.Query().Get("active_only") == "true"

	rules, err := h.queryService.GetRulesByType(r.Context(), ruleType, activeOnly)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rules)
}

func (h *RuleHandler) EvaluateRules(w http.ResponseWriter, r *http.Request) {
	var cmd commands.EvaluateRuleCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	actions, err := h.commandHandler.HandleEvaluateRule(r.Context(), cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(actions)
}
