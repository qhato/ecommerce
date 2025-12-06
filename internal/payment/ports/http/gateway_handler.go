package http

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/qhato/ecommerce/internal/payment/application/queries"
)

type GatewayHandler struct {
	queryService *queries.GatewayQueryService
}

func NewGatewayHandler(queryService *queries.GatewayQueryService) *GatewayHandler {
	return &GatewayHandler{
		queryService: queryService,
	}
}

func (h *GatewayHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/payment-gateways", h.GetAllGateways).Methods("GET")
	router.HandleFunc("/payment-gateways/enabled", h.GetEnabledGateways).Methods("GET")
	router.HandleFunc("/payment-gateways/{name}", h.GetGateway).Methods("GET")
}

func (h *GatewayHandler) GetAllGateways(w http.ResponseWriter, r *http.Request) {
	gateways, err := h.queryService.GetAllGatewayConfigs(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(gateways)
}

func (h *GatewayHandler) GetEnabledGateways(w http.ResponseWriter, r *http.Request) {
	gateways, err := h.queryService.GetEnabledGatewayConfigs(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(gateways)
}

func (h *GatewayHandler) GetGateway(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	gateway, err := h.queryService.GetGatewayConfig(r.Context(), name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(gateway)
}
