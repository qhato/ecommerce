package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/qhato/ecommerce/internal/tax/application/commands"
	"github.com/qhato/ecommerce/internal/tax/application/queries"
	"github.com/qhato/ecommerce/internal/tax/domain"
)

// TaxHandler handles HTTP requests for the tax engine
type TaxHandler struct {
	commandHandler   *commands.TaxCommandHandler
	calculatorService *queries.TaxCalculatorService
}

// NewTaxHandler creates a new tax HTTP handler
func NewTaxHandler(
	commandHandler *commands.TaxCommandHandler,
	calculatorService *queries.TaxCalculatorService,
) *TaxHandler {
	return &TaxHandler{
		commandHandler:   commandHandler,
		calculatorService: calculatorService,
	}
}

// RegisterRoutes registers all tax routes
func (h *TaxHandler) RegisterRoutes(router *mux.Router) {
	// Tax Calculation Endpoints
	router.HandleFunc("/tax/calculate", h.CalculateTax).Methods("POST")
	router.HandleFunc("/tax/estimate", h.EstimateTax).Methods("POST")
	router.HandleFunc("/tax/validate-address", h.ValidateAddress).Methods("POST")

	// Tax Jurisdiction Endpoints
	router.HandleFunc("/tax/jurisdictions", h.CreateJurisdiction).Methods("POST")
	router.HandleFunc("/tax/jurisdictions", h.GetAllJurisdictions).Methods("GET")
	router.HandleFunc("/tax/jurisdictions/{id}", h.GetJurisdiction).Methods("GET")
	router.HandleFunc("/tax/jurisdictions/{id}", h.UpdateJurisdiction).Methods("PUT")
	router.HandleFunc("/tax/jurisdictions/{id}", h.DeleteJurisdiction).Methods("DELETE")
	router.HandleFunc("/tax/jurisdictions/code/{code}", h.GetJurisdictionByCode).Methods("GET")
	router.HandleFunc("/tax/jurisdictions/country/{country}", h.GetJurisdictionsByCountry).Methods("GET")

	// Tax Rate Endpoints
	router.HandleFunc("/tax/rates", h.CreateTaxRate).Methods("POST")
	router.HandleFunc("/tax/rates", h.GetAllTaxRates).Methods("GET")
	router.HandleFunc("/tax/rates/{id}", h.GetTaxRate).Methods("GET")
	router.HandleFunc("/tax/rates/{id}", h.UpdateTaxRate).Methods("PUT")
	router.HandleFunc("/tax/rates/{id}", h.DeleteTaxRate).Methods("DELETE")
	router.HandleFunc("/tax/rates/jurisdiction/{id}", h.GetTaxRatesByJurisdiction).Methods("GET")
	router.HandleFunc("/tax/rates/bulk", h.BulkCreateTaxRates).Methods("POST")

	// Tax Exemption Endpoints
	router.HandleFunc("/tax/exemptions", h.CreateTaxExemption).Methods("POST")
	router.HandleFunc("/tax/exemptions", h.GetAllTaxExemptions).Methods("GET")
	router.HandleFunc("/tax/exemptions/{id}", h.GetTaxExemption).Methods("GET")
	router.HandleFunc("/tax/exemptions/{id}", h.UpdateTaxExemption).Methods("PUT")
	router.HandleFunc("/tax/exemptions/{id}", h.DeleteTaxExemption).Methods("DELETE")
	router.HandleFunc("/tax/exemptions/customer/{customerId}", h.GetTaxExemptionsByCustomer).Methods("GET")
}

// Tax Calculation Handlers

func (h *TaxHandler) CalculateTax(w http.ResponseWriter, r *http.Request) {
	var req queries.CalculateTaxRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Convert to domain request
	domainReq := queries.ToTaxCalculationRequest(req)

	// Calculate taxes
	result, err := h.calculatorService.Calculate(r.Context(), domainReq)
	if err != nil {
		if err == domain.ErrNoApplicableJurisdictions {
			respondWithError(w, http.StatusNotFound, err.Error())
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to calculate taxes: "+err.Error())
		return
	}

	// Convert to DTO response
	response := queries.ToTaxCalculationResponse(result)
	respondWithJSON(w, http.StatusOK, response)
}

func (h *TaxHandler) EstimateTax(w http.ResponseWriter, r *http.Request) {
	var req queries.EstimateTaxRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Convert address to domain
	address := queries.ToAddressDomain(req.Address)

	// Estimate taxes
	estimatedTax, err := h.calculatorService.EstimateTax(r.Context(), address, req.Subtotal)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to estimate taxes: "+err.Error())
		return
	}

	// Calculate effective rate
	effectiveRate := estimatedTax.Div(req.Subtotal)

	response := queries.EstimateTaxResponse{
		EstimatedTax:     estimatedTax,
		EffectiveTaxRate: effectiveRate,
	}
	respondWithJSON(w, http.StatusOK, response)
}

func (h *TaxHandler) ValidateAddress(w http.ResponseWriter, r *http.Request) {
	var req queries.AddressDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Convert to domain address
	address := queries.ToAddressDomain(req)

	// Validate address
	isValid, err := h.calculatorService.ValidateAddress(r.Context(), address)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to validate address: "+err.Error())
		return
	}

	response := map[string]bool{"valid": isValid}
	respondWithJSON(w, http.StatusOK, response)
}

// Tax Jurisdiction Handlers

func (h *TaxHandler) CreateJurisdiction(w http.ResponseWriter, r *http.Request) {
	var cmd commands.CreateTaxJurisdictionCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	jurisdiction, err := h.commandHandler.HandleCreateTaxJurisdiction(r.Context(), cmd)
	if err != nil {
		if err == domain.ErrJurisdictionAlreadyExists {
			respondWithError(w, http.StatusConflict, err.Error())
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to create jurisdiction: "+err.Error())
		return
	}

	response := queries.ToTaxJurisdictionDTO(jurisdiction)
	respondWithJSON(w, http.StatusCreated, response)
}

func (h *TaxHandler) GetAllJurisdictions(w http.ResponseWriter, r *http.Request) {
	activeOnly := r.URL.Query().Get("active") == "true"

	jurisdictions, err := h.calculatorService.GetAllJurisdictions(r.Context(), activeOnly)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get jurisdictions: "+err.Error())
		return
	}

	response := make([]queries.TaxJurisdictionDTO, len(jurisdictions))
	for i, j := range jurisdictions {
		response[i] = queries.ToTaxJurisdictionDTO(j)
	}

	respondWithJSON(w, http.StatusOK, response)
}

func (h *TaxHandler) GetJurisdiction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid jurisdiction ID")
		return
	}

	jurisdiction, err := h.calculatorService.GetJurisdictionByID(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get jurisdiction: "+err.Error())
		return
	}
	if jurisdiction == nil {
		respondWithError(w, http.StatusNotFound, "Jurisdiction not found")
		return
	}

	response := queries.ToTaxJurisdictionDTO(jurisdiction)
	respondWithJSON(w, http.StatusOK, response)
}

func (h *TaxHandler) GetJurisdictionByCode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	code := vars["code"]

	jurisdiction, err := h.calculatorService.GetJurisdictionByCode(r.Context(), code)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get jurisdiction: "+err.Error())
		return
	}
	if jurisdiction == nil {
		respondWithError(w, http.StatusNotFound, "Jurisdiction not found")
		return
	}

	response := queries.ToTaxJurisdictionDTO(jurisdiction)
	respondWithJSON(w, http.StatusOK, response)
}

func (h *TaxHandler) GetJurisdictionsByCountry(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	country := vars["country"]
	activeOnly := r.URL.Query().Get("active") == "true"

	jurisdictions, err := h.calculatorService.GetJurisdictionsByCountry(r.Context(), country, activeOnly)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get jurisdictions: "+err.Error())
		return
	}

	response := make([]queries.TaxJurisdictionDTO, len(jurisdictions))
	for i, j := range jurisdictions {
		response[i] = queries.ToTaxJurisdictionDTO(j)
	}

	respondWithJSON(w, http.StatusOK, response)
}

func (h *TaxHandler) UpdateJurisdiction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid jurisdiction ID")
		return
	}

	var cmd commands.UpdateTaxJurisdictionCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	cmd.ID = id

	jurisdiction, err := h.commandHandler.HandleUpdateTaxJurisdiction(r.Context(), cmd)
	if err != nil {
		if err == domain.ErrJurisdictionNotFound {
			respondWithError(w, http.StatusNotFound, err.Error())
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to update jurisdiction: "+err.Error())
		return
	}

	response := queries.ToTaxJurisdictionDTO(jurisdiction)
	respondWithJSON(w, http.StatusOK, response)
}

func (h *TaxHandler) DeleteJurisdiction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid jurisdiction ID")
		return
	}

	cmd := commands.DeleteTaxJurisdictionCommand{ID: id}
	if err := h.commandHandler.HandleDeleteTaxJurisdiction(r.Context(), cmd); err != nil {
		if err == domain.ErrJurisdictionNotFound {
			respondWithError(w, http.StatusNotFound, err.Error())
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to delete jurisdiction: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}

// Tax Rate Handlers

func (h *TaxHandler) CreateTaxRate(w http.ResponseWriter, r *http.Request) {
	var cmd commands.CreateTaxRateCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	rate, err := h.commandHandler.HandleCreateTaxRate(r.Context(), cmd)
	if err != nil {
		if err == domain.ErrJurisdictionNotFound {
			respondWithError(w, http.StatusNotFound, err.Error())
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to create tax rate: "+err.Error())
		return
	}

	response := queries.ToTaxRateDTO(rate)
	respondWithJSON(w, http.StatusCreated, response)
}

func (h *TaxHandler) GetAllTaxRates(w http.ResponseWriter, r *http.Request) {
	activeOnly := r.URL.Query().Get("active") == "true"

	rates, err := h.calculatorService.GetAllTaxRates(r.Context(), activeOnly)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get tax rates: "+err.Error())
		return
	}

	response := make([]queries.TaxRateDTO, len(rates))
	for i, rate := range rates {
		response[i] = queries.ToTaxRateDTO(rate)
	}

	respondWithJSON(w, http.StatusOK, response)
}

func (h *TaxHandler) GetTaxRate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid tax rate ID")
		return
	}

	rate, err := h.calculatorService.GetTaxRateByID(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get tax rate: "+err.Error())
		return
	}
	if rate == nil {
		respondWithError(w, http.StatusNotFound, "Tax rate not found")
		return
	}

	response := queries.ToTaxRateDTO(rate)
	respondWithJSON(w, http.StatusOK, response)
}

func (h *TaxHandler) GetTaxRatesByJurisdiction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid jurisdiction ID")
		return
	}

	activeOnly := r.URL.Query().Get("active") == "true"

	rates, err := h.calculatorService.GetTaxRatesByJurisdiction(r.Context(), id, activeOnly)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get tax rates: "+err.Error())
		return
	}

	response := make([]queries.TaxRateDTO, len(rates))
	for i, rate := range rates {
		response[i] = queries.ToTaxRateDTO(rate)
	}

	respondWithJSON(w, http.StatusOK, response)
}

func (h *TaxHandler) UpdateTaxRate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid tax rate ID")
		return
	}

	var cmd commands.UpdateTaxRateCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	cmd.ID = id

	rate, err := h.commandHandler.HandleUpdateTaxRate(r.Context(), cmd)
	if err != nil {
		if err == domain.ErrTaxRateNotFound {
			respondWithError(w, http.StatusNotFound, err.Error())
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to update tax rate: "+err.Error())
		return
	}

	response := queries.ToTaxRateDTO(rate)
	respondWithJSON(w, http.StatusOK, response)
}

func (h *TaxHandler) DeleteTaxRate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid tax rate ID")
		return
	}

	cmd := commands.DeleteTaxRateCommand{ID: id}
	if err := h.commandHandler.HandleDeleteTaxRate(r.Context(), cmd); err != nil {
		if err == domain.ErrTaxRateNotFound {
			respondWithError(w, http.StatusNotFound, err.Error())
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to delete tax rate: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}

func (h *TaxHandler) BulkCreateTaxRates(w http.ResponseWriter, r *http.Request) {
	var cmd commands.BulkCreateTaxRatesCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	rates, err := h.commandHandler.HandleBulkCreateTaxRates(r.Context(), cmd)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to bulk create tax rates: "+err.Error())
		return
	}

	response := make([]queries.TaxRateDTO, len(rates))
	for i, rate := range rates {
		response[i] = queries.ToTaxRateDTO(rate)
	}

	respondWithJSON(w, http.StatusCreated, response)
}

// Tax Exemption Handlers

func (h *TaxHandler) CreateTaxExemption(w http.ResponseWriter, r *http.Request) {
	var cmd commands.CreateTaxExemptionCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	exemption, err := h.commandHandler.HandleCreateTaxExemption(r.Context(), cmd)
	if err != nil {
		if err == domain.ErrExemptionAlreadyExists {
			respondWithError(w, http.StatusConflict, err.Error())
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to create tax exemption: "+err.Error())
		return
	}

	response := queries.ToTaxExemptionDTO(exemption)
	respondWithJSON(w, http.StatusCreated, response)
}

func (h *TaxHandler) GetAllTaxExemptions(w http.ResponseWriter, r *http.Request) {
	activeOnly := r.URL.Query().Get("active") == "true"

	exemptions, err := h.calculatorService.GetAllExemptions(r.Context(), activeOnly)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get tax exemptions: "+err.Error())
		return
	}

	response := make([]queries.TaxExemptionDTO, len(exemptions))
	for i, exemption := range exemptions {
		response[i] = queries.ToTaxExemptionDTO(exemption)
	}

	respondWithJSON(w, http.StatusOK, response)
}

func (h *TaxHandler) GetTaxExemption(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid tax exemption ID")
		return
	}

	exemption, err := h.calculatorService.GetExemptionByID(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get tax exemption: "+err.Error())
		return
	}
	if exemption == nil {
		respondWithError(w, http.StatusNotFound, "Tax exemption not found")
		return
	}

	response := queries.ToTaxExemptionDTO(exemption)
	respondWithJSON(w, http.StatusOK, response)
}

func (h *TaxHandler) GetTaxExemptionsByCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customerID := vars["customerId"]
	activeOnly := r.URL.Query().Get("active") == "true"

	exemptions, err := h.calculatorService.GetExemptionsByCustomer(r.Context(), customerID, activeOnly)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get tax exemptions: "+err.Error())
		return
	}

	response := make([]queries.TaxExemptionDTO, len(exemptions))
	for i, exemption := range exemptions {
		response[i] = queries.ToTaxExemptionDTO(exemption)
	}

	respondWithJSON(w, http.StatusOK, response)
}

func (h *TaxHandler) UpdateTaxExemption(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid tax exemption ID")
		return
	}

	var cmd commands.UpdateTaxExemptionCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	cmd.ID = id

	exemption, err := h.commandHandler.HandleUpdateTaxExemption(r.Context(), cmd)
	if err != nil {
		if err == domain.ErrExemptionNotFound {
			respondWithError(w, http.StatusNotFound, err.Error())
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to update tax exemption: "+err.Error())
		return
	}

	response := queries.ToTaxExemptionDTO(exemption)
	respondWithJSON(w, http.StatusOK, response)
}

func (h *TaxHandler) DeleteTaxExemption(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid tax exemption ID")
		return
	}

	cmd := commands.DeleteTaxExemptionCommand{ID: id}
	if err := h.commandHandler.HandleDeleteTaxExemption(r.Context(), cmd); err != nil {
		if err == domain.ErrExemptionNotFound {
			respondWithError(w, http.StatusNotFound, err.Error())
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to delete tax exemption: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}

// Helper functions

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if payload != nil {
		json.NewEncoder(w).Encode(payload)
	}
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}
