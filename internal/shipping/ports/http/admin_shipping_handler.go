package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/qhato/ecommerce/internal/shipping/application/commands"
	"github.com/qhato/ecommerce/internal/shipping/application/queries"
	"github.com/qhato/ecommerce/pkg/logger"
	"github.com/qhato/ecommerce/pkg/validator"
)

// AdminShippingHandler handles admin shipping HTTP requests using CQRS
type AdminShippingHandler struct {
	commandHandler *commands.ShippingCommandHandler
	queryService   *queries.ShippingQueryService
	validator      *validator.Validator
	log            *logger.Logger
}

// NewAdminShippingHandler creates a new admin shipping HTTP handler
func NewAdminShippingHandler(
	commandHandler *commands.ShippingCommandHandler,
	queryService *queries.ShippingQueryService,
	validator *validator.Validator,
	log *logger.Logger,
) *AdminShippingHandler {
	return &AdminShippingHandler{
		commandHandler: commandHandler,
		queryService:   queryService,
		validator:      validator,
		log:            log,
	}
}

// RegisterRoutes registers all admin shipping routes
func (h *AdminShippingHandler) RegisterRoutes(r chi.Router) {
	r.Route("/admin/shipping", func(r chi.Router) {
		// Shipping Methods
		r.Post("/methods", h.CreateShippingMethod)
		r.Get("/methods", h.GetAllShippingMethods)
		r.Get("/methods/{id}", h.GetShippingMethod)
		r.Put("/methods/{id}", h.UpdateShippingMethod)
		r.Delete("/methods/{id}", h.DeleteShippingMethod)

		// Shipping Bands
		r.Post("/methods/{methodId}/bands", h.CreateShippingBand)
		r.Get("/methods/{methodId}/bands", h.GetShippingBands)
		r.Put("/bands/{id}", h.UpdateShippingBand)
		r.Delete("/bands/{id}", h.DeleteShippingBand)

		// Shipping Rules
		r.Post("/rules", h.CreateShippingRule)
		r.Get("/rules", h.GetAllShippingRules)
		r.Get("/rules/{id}", h.GetShippingRule)
		r.Put("/rules/{id}", h.UpdateShippingRule)
		r.Delete("/rules/{id}", h.DeleteShippingRule)

		// Carrier Configs
		r.Post("/carriers", h.CreateCarrierConfig)
		r.Get("/carriers", h.GetAllCarrierConfigs)
		r.Get("/carriers/{id}", h.GetCarrierConfig)
		r.Put("/carriers/{id}", h.UpdateCarrierConfig)
		r.Delete("/carriers/{id}", h.DeleteCarrierConfig)
	})
}

// Shipping Methods

func (h *AdminShippingHandler) CreateShippingMethod(w http.ResponseWriter, r *http.Request) {
	var cmd commands.CreateShippingMethodCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.validator.Validate(cmd); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	method, err := h.commandHandler.HandleCreateShippingMethod(r.Context(), cmd)
	if err != nil {
		h.log.WithError(err).Error("Failed to create shipping method")
		respondWithError(w, http.StatusInternalServerError, "Failed to create shipping method")
		return
	}

	respondWithJSON(w, http.StatusCreated, method)
}

func (h *AdminShippingHandler) GetAllShippingMethods(w http.ResponseWriter, r *http.Request) {
	query := queries.GetAllEnabledShippingMethodsQuery{}

	methods, err := h.queryService.GetAllEnabledShippingMethods(r.Context(), query)
	if err != nil {
		h.log.WithError(err).Error("Failed to get shipping methods")
		respondWithError(w, http.StatusInternalServerError, "Failed to get shipping methods")
		return
	}

	respondWithJSON(w, http.StatusOK, methods)
}

func (h *AdminShippingHandler) GetShippingMethod(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid method ID")
		return
	}

	query := queries.GetShippingMethodQuery{ID: id}
	method, err := h.queryService.GetShippingMethod(r.Context(), query)
	if err != nil {
		h.log.WithError(err).Error("Failed to get shipping method")
		respondWithError(w, http.StatusNotFound, "Shipping method not found")
		return
	}

	respondWithJSON(w, http.StatusOK, method)
}

func (h *AdminShippingHandler) UpdateShippingMethod(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid method ID")
		return
	}

	var cmd commands.UpdateShippingMethodCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	cmd.ID = id

	if err := h.validator.Validate(cmd); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	method, err := h.commandHandler.HandleUpdateShippingMethod(r.Context(), cmd)
	if err != nil {
		h.log.WithError(err).Error("Failed to update shipping method")
		respondWithError(w, http.StatusInternalServerError, "Failed to update shipping method")
		return
	}

	respondWithJSON(w, http.StatusOK, method)
}

func (h *AdminShippingHandler) DeleteShippingMethod(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid method ID")
		return
	}

	cmd := commands.DeleteShippingMethodCommand{ID: id}
	if err := h.commandHandler.HandleDeleteShippingMethod(r.Context(), cmd); err != nil {
		h.log.WithError(err).Error("Failed to delete shipping method")
		respondWithError(w, http.StatusInternalServerError, "Failed to delete shipping method")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Shipping Bands

func (h *AdminShippingHandler) CreateShippingBand(w http.ResponseWriter, r *http.Request) {
	methodIDStr := chi.URLParam(r, "methodId")
	methodID, err := strconv.ParseInt(methodIDStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid method ID")
		return
	}

	var cmd commands.CreateShippingBandCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	cmd.MethodID = methodID

	if err := h.validator.Validate(cmd); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	band, err := h.commandHandler.HandleCreateShippingBand(r.Context(), cmd)
	if err != nil {
		h.log.WithError(err).Error("Failed to create shipping band")
		respondWithError(w, http.StatusInternalServerError, "Failed to create shipping band")
		return
	}

	respondWithJSON(w, http.StatusCreated, band)
}

func (h *AdminShippingHandler) GetShippingBands(w http.ResponseWriter, r *http.Request) {
	methodIDStr := chi.URLParam(r, "methodId")
	methodID, err := strconv.ParseInt(methodIDStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid method ID")
		return
	}

	query := queries.GetShippingBandsByMethodQuery{MethodID: methodID}
	bands, err := h.queryService.GetShippingBandsByMethod(r.Context(), query)
	if err != nil {
		h.log.WithError(err).Error("Failed to get shipping bands")
		respondWithError(w, http.StatusInternalServerError, "Failed to get shipping bands")
		return
	}

	respondWithJSON(w, http.StatusOK, bands)
}

func (h *AdminShippingHandler) UpdateShippingBand(w http.ResponseWriter, r *http.Request) {
	// Not implemented - shipping bands can only be created and deleted, not updated
	// To update, delete and recreate
	respondWithError(w, http.StatusNotImplemented, "Shipping bands cannot be updated. Delete and recreate instead.")
}

func (h *AdminShippingHandler) DeleteShippingBand(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid band ID")
		return
	}

	cmd := commands.DeleteShippingBandCommand{ID: id}
	if err := h.commandHandler.HandleDeleteShippingBand(r.Context(), cmd); err != nil {
		h.log.WithError(err).Error("Failed to delete shipping band")
		respondWithError(w, http.StatusInternalServerError, "Failed to delete shipping band")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Shipping Rules

func (h *AdminShippingHandler) CreateShippingRule(w http.ResponseWriter, r *http.Request) {
	var cmd commands.CreateShippingRuleCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.validator.Validate(cmd); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	rule, err := h.commandHandler.HandleCreateShippingRule(r.Context(), cmd)
	if err != nil {
		h.log.WithError(err).Error("Failed to create shipping rule")
		respondWithError(w, http.StatusInternalServerError, "Failed to create shipping rule")
		return
	}

	respondWithJSON(w, http.StatusCreated, rule)
}

func (h *AdminShippingHandler) GetAllShippingRules(w http.ResponseWriter, r *http.Request) {
	query := queries.GetAllEnabledShippingRulesQuery{}

	rules, err := h.queryService.GetAllEnabledShippingRules(r.Context(), query)
	if err != nil {
		h.log.WithError(err).Error("Failed to get shipping rules")
		respondWithError(w, http.StatusInternalServerError, "Failed to get shipping rules")
		return
	}

	respondWithJSON(w, http.StatusOK, rules)
}

func (h *AdminShippingHandler) GetShippingRule(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid rule ID")
		return
	}

	query := queries.GetShippingRuleQuery{ID: id}
	rule, err := h.queryService.GetShippingRule(r.Context(), query)
	if err != nil {
		h.log.WithError(err).Error("Failed to get shipping rule")
		respondWithError(w, http.StatusNotFound, "Shipping rule not found")
		return
	}

	respondWithJSON(w, http.StatusOK, rule)
}

func (h *AdminShippingHandler) UpdateShippingRule(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid rule ID")
		return
	}

	var cmd commands.UpdateShippingRuleCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	cmd.ID = id

	if err := h.validator.Validate(cmd); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	rule, err := h.commandHandler.HandleUpdateShippingRule(r.Context(), cmd)
	if err != nil {
		h.log.WithError(err).Error("Failed to update shipping rule")
		respondWithError(w, http.StatusInternalServerError, "Failed to update shipping rule")
		return
	}

	respondWithJSON(w, http.StatusOK, rule)
}

func (h *AdminShippingHandler) DeleteShippingRule(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid rule ID")
		return
	}

	cmd := commands.DeleteShippingRuleCommand{ID: id}
	if err := h.commandHandler.HandleDeleteShippingRule(r.Context(), cmd); err != nil {
		h.log.WithError(err).Error("Failed to delete shipping rule")
		respondWithError(w, http.StatusInternalServerError, "Failed to delete shipping rule")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Carrier Configs

func (h *AdminShippingHandler) CreateCarrierConfig(w http.ResponseWriter, r *http.Request) {
	var cmd commands.CreateCarrierConfigCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.validator.Validate(cmd); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	config, err := h.commandHandler.HandleCreateCarrierConfig(r.Context(), cmd)
	if err != nil {
		h.log.WithError(err).Error("Failed to create carrier config")
		respondWithError(w, http.StatusInternalServerError, "Failed to create carrier config")
		return
	}

	respondWithJSON(w, http.StatusCreated, config)
}

func (h *AdminShippingHandler) GetAllCarrierConfigs(w http.ResponseWriter, r *http.Request) {
	query := queries.GetAllCarrierConfigsQuery{
		EnabledOnly: r.URL.Query().Get("enabled_only") == "true",
	}

	configs, err := h.queryService.GetAllCarrierConfigs(r.Context(), query)
	if err != nil {
		h.log.WithError(err).Error("Failed to get carrier configs")
		respondWithError(w, http.StatusInternalServerError, "Failed to get carrier configs")
		return
	}

	respondWithJSON(w, http.StatusOK, configs)
}

func (h *AdminShippingHandler) GetCarrierConfig(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid config ID")
		return
	}

	query := queries.GetCarrierConfigQuery{ID: id}
	config, err := h.queryService.GetCarrierConfig(r.Context(), query)
	if err != nil {
		h.log.WithError(err).Error("Failed to get carrier config")
		respondWithError(w, http.StatusNotFound, "Carrier config not found")
		return
	}

	respondWithJSON(w, http.StatusOK, config)
}

func (h *AdminShippingHandler) UpdateCarrierConfig(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid config ID")
		return
	}

	var cmd commands.UpdateCarrierConfigCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	cmd.ID = id

	if err := h.validator.Validate(cmd); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	config, err := h.commandHandler.HandleUpdateCarrierConfig(r.Context(), cmd)
	if err != nil {
		h.log.WithError(err).Error("Failed to update carrier config")
		respondWithError(w, http.StatusInternalServerError, "Failed to update carrier config")
		return
	}

	respondWithJSON(w, http.StatusOK, config)
}

func (h *AdminShippingHandler) DeleteCarrierConfig(w http.ResponseWriter, r *http.Request) {
	// Not implemented - carrier configs are typically not deleted, just disabled
	// Set IsEnabled to false instead
	respondWithError(w, http.StatusNotImplemented, "Carrier configs cannot be deleted. Disable instead by setting is_enabled to false.")
}

// Helper functions

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"Failed to marshal response"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
