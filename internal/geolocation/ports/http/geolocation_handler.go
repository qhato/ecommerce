package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/qhato/ecommerce/internal/geolocation/application/commands"
	"github.com/qhato/ecommerce/internal/geolocation/application/queries"
	"github.com/qhato/ecommerce/internal/geolocation/domain"
)

type GeolocationHandler struct {
	commandHandler *commands.GeolocationCommandHandler
	queryService   *queries.GeolocationQueryService
}

func NewGeolocationHandler(
	commandHandler *commands.GeolocationCommandHandler,
	queryService *queries.GeolocationQueryService,
) *GeolocationHandler {
	return &GeolocationHandler{
		commandHandler: commandHandler,
		queryService:   queryService,
	}
}

func (h *GeolocationHandler) RegisterRoutes(router *mux.Router) {
	// Location endpoints
	router.HandleFunc("/geolocation/locations", h.CreateLocation).Methods("POST")
	router.HandleFunc("/geolocation/locations/{id}", h.GetLocation).Methods("GET")
	router.HandleFunc("/geolocation/locations/{id}", h.UpdateLocation).Methods("PUT")
	router.HandleFunc("/geolocation/locations/{id}", h.DeleteLocation).Methods("DELETE")
	router.HandleFunc("/geolocation/locations/search", h.SearchLocations).Methods("GET")
	router.HandleFunc("/geolocation/locations/nearby", h.FindNearbyLocations).Methods("GET")
	router.HandleFunc("/geolocation/locations/radius", h.FindLocationsInRadius).Methods("GET")
	router.HandleFunc("/geolocation/locations/city/{city}", h.FindLocationsByCity).Methods("GET")
	router.HandleFunc("/geolocation/locations/postal/{postalCode}", h.FindLocationsByPostalCode).Methods("GET")
	router.HandleFunc("/geolocation/distance", h.CalculateDistance).Methods("GET")

	// GeoZone endpoints
	router.HandleFunc("/geolocation/zones", h.CreateGeoZone).Methods("POST")
	router.HandleFunc("/geolocation/zones/{id}", h.GetGeoZone).Methods("GET")
	router.HandleFunc("/geolocation/zones/{id}", h.UpdateGeoZone).Methods("PUT")
	router.HandleFunc("/geolocation/zones/{id}", h.DeleteGeoZone).Methods("DELETE")
	router.HandleFunc("/geolocation/zones/type/{type}", h.GetGeoZonesByType).Methods("GET")
	router.HandleFunc("/geolocation/zones/active", h.GetActiveGeoZones).Methods("GET")
	router.HandleFunc("/geolocation/zones/location/{locationId}", h.GetGeoZonesForLocation).Methods("GET")
}

// Location handlers

func (h *GeolocationHandler) CreateLocation(w http.ResponseWriter, r *http.Request) {
	var cmd commands.CreateLocationCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	location, err := h.commandHandler.HandleCreateLocation(r.Context(), cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(queries.ToLocationDTO(location))
}

func (h *GeolocationHandler) GetLocation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid location ID", http.StatusBadRequest)
		return
	}

	location, err := h.queryService.GetLocation(r.Context(), id)
	if err != nil {
		if err == domain.ErrLocationNotFound {
			http.Error(w, "Location not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(location)
}

func (h *GeolocationHandler) UpdateLocation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid location ID", http.StatusBadRequest)
		return
	}

	var cmd commands.UpdateLocationCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	cmd.ID = id

	location, err := h.commandHandler.HandleUpdateLocation(r.Context(), cmd)
	if err != nil {
		if err == domain.ErrLocationNotFound {
			http.Error(w, "Location not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(queries.ToLocationDTO(location))
}

func (h *GeolocationHandler) DeleteLocation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid location ID", http.StatusBadRequest)
		return
	}

	cmd := commands.DeleteLocationCommand{ID: id}
	if err := h.commandHandler.HandleDeleteLocation(r.Context(), cmd); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *GeolocationHandler) SearchLocations(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	limitStr := r.URL.Query().Get("limit")
	limit := 50
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil {
			limit = parsedLimit
		}
	}

	locations, err := h.queryService.SearchLocations(r.Context(), query, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(locations)
}

func (h *GeolocationHandler) FindNearbyLocations(w http.ResponseWriter, r *http.Request) {
	latStr := r.URL.Query().Get("lat")
	lngStr := r.URL.Query().Get("lng")
	limitStr := r.URL.Query().Get("limit")

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		http.Error(w, "Invalid latitude", http.StatusBadRequest)
		return
	}

	lng, err := strconv.ParseFloat(lngStr, 64)
	if err != nil {
		http.Error(w, "Invalid longitude", http.StatusBadRequest)
		return
	}

	limit := 10
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil {
			limit = parsedLimit
		}
	}

	locations, err := h.queryService.FindNearbyLocations(r.Context(), lat, lng, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(locations)
}

func (h *GeolocationHandler) FindLocationsInRadius(w http.ResponseWriter, r *http.Request) {
	latStr := r.URL.Query().Get("lat")
	lngStr := r.URL.Query().Get("lng")
	radiusStr := r.URL.Query().Get("radius")

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		http.Error(w, "Invalid latitude", http.StatusBadRequest)
		return
	}

	lng, err := strconv.ParseFloat(lngStr, 64)
	if err != nil {
		http.Error(w, "Invalid longitude", http.StatusBadRequest)
		return
	}

	radius, err := strconv.ParseFloat(radiusStr, 64)
	if err != nil {
		http.Error(w, "Invalid radius", http.StatusBadRequest)
		return
	}

	locations, err := h.queryService.FindLocationsInRadius(r.Context(), lat, lng, radius)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(locations)
}

func (h *GeolocationHandler) FindLocationsByCity(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	city := vars["city"]
	state := r.URL.Query().Get("state")
	country := r.URL.Query().Get("country")

	locations, err := h.queryService.FindLocationsByCity(r.Context(), city, state, country)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(locations)
}

func (h *GeolocationHandler) FindLocationsByPostalCode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postalCode := vars["postalCode"]
	country := r.URL.Query().Get("country")

	locations, err := h.queryService.FindLocationsByPostalCode(r.Context(), postalCode, country)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(locations)
}

func (h *GeolocationHandler) CalculateDistance(w http.ResponseWriter, r *http.Request) {
	fromIDStr := r.URL.Query().Get("from")
	toIDStr := r.URL.Query().Get("to")
	unit := r.URL.Query().Get("unit")

	fromID, err := strconv.ParseInt(fromIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid from location ID", http.StatusBadRequest)
		return
	}

	toID, err := strconv.ParseInt(toIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid to location ID", http.StatusBadRequest)
		return
	}

	if unit == "" {
		unit = "KM"
	}

	result, err := h.queryService.CalculateDistance(r.Context(), fromID, toID, unit)
	if err != nil {
		if err == domain.ErrLocationNotFound {
			http.Error(w, "Location not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// GeoZone handlers

func (h *GeolocationHandler) CreateGeoZone(w http.ResponseWriter, r *http.Request) {
	var cmd commands.CreateGeoZoneCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	zone, err := h.commandHandler.HandleCreateGeoZone(r.Context(), cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(queries.ToGeoZoneDTO(zone))
}

func (h *GeolocationHandler) GetGeoZone(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid zone ID", http.StatusBadRequest)
		return
	}

	zone, err := h.queryService.GetGeoZone(r.Context(), id)
	if err != nil {
		if err == domain.ErrGeoZoneNotFound {
			http.Error(w, "Geo zone not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(zone)
}

func (h *GeolocationHandler) UpdateGeoZone(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid zone ID", http.StatusBadRequest)
		return
	}

	var cmd commands.UpdateGeoZoneCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	cmd.ID = id

	zone, err := h.commandHandler.HandleUpdateGeoZone(r.Context(), cmd)
	if err != nil {
		if err == domain.ErrGeoZoneNotFound {
			http.Error(w, "Geo zone not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(queries.ToGeoZoneDTO(zone))
}

func (h *GeolocationHandler) DeleteGeoZone(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid zone ID", http.StatusBadRequest)
		return
	}

	cmd := commands.DeleteGeoZoneCommand{ID: id}
	if err := h.commandHandler.HandleDeleteGeoZone(r.Context(), cmd); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *GeolocationHandler) GetGeoZonesByType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	zoneType := vars["type"]

	zones, err := h.queryService.GetGeoZonesByType(r.Context(), zoneType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(zones)
}

func (h *GeolocationHandler) GetActiveGeoZones(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit := 50
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil {
			limit = parsedLimit
		}
	}

	zones, err := h.queryService.GetActiveGeoZones(r.Context(), limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(zones)
}

func (h *GeolocationHandler) GetGeoZonesForLocation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	locationID, err := strconv.ParseInt(vars["locationId"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid location ID", http.StatusBadRequest)
		return
	}

	zones, err := h.queryService.GetGeoZonesForLocation(r.Context(), locationID)
	if err != nil {
		if err == domain.ErrLocationNotFound {
			http.Error(w, "Location not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(zones)
}
