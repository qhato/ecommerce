package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/qhato/ecommerce/internal/media/application/commands"
	"github.com/qhato/ecommerce/internal/media/application/queries"
)

type MediaHandler struct {
	commandHandler *commands.MediaCommandHandler
	queryService   *queries.MediaQueryService
}

func NewMediaHandler(
	commandHandler *commands.MediaCommandHandler,
	queryService *queries.MediaQueryService,
) *MediaHandler {
	return &MediaHandler{
		commandHandler: commandHandler,
		queryService:   queryService,
	}
}

func (h *MediaHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/media", h.CreateMedia).Methods("POST")
	router.HandleFunc("/media/{id}", h.GetMedia).Methods("GET")
	router.HandleFunc("/media/{id}", h.UpdateMedia).Methods("PUT")
	router.HandleFunc("/media/{id}", h.DeleteMedia).Methods("DELETE")
	router.HandleFunc("/media/{id}/activate", h.ActivateMedia).Methods("POST")
	router.HandleFunc("/media/{id}/archive", h.ArchiveMedia).Methods("POST")
	router.HandleFunc("/media/entity/{entityType}/{entityId}", h.GetMediaByEntity).Methods("GET")
	router.HandleFunc("/media/type/{type}", h.GetMediaByType).Methods("GET")
	router.HandleFunc("/media", h.ListMedia).Methods("GET")
}

type CreateMediaRequest struct {
	Name        string   `json:"name" validate:"required"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	MimeType    string   `json:"mime_type" validate:"required"`
	FilePath    string   `json:"file_path" validate:"required"`
	FileSize    int64    `json:"file_size"`
	UploadedBy  string   `json:"uploaded_by"`
	EntityType  *string  `json:"entity_type"`
	EntityID    *string  `json:"entity_id"`
	Tags        []string `json:"tags"`
}

type UpdateMediaRequest struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

func (h *MediaHandler) CreateMedia(w http.ResponseWriter, r *http.Request) {
	var req CreateMediaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cmd := commands.CreateMediaCommand{
		Name:        req.Name,
		Title:       req.Title,
		Description: req.Description,
		MimeType:    req.MimeType,
		FilePath:    req.FilePath,
		FileSize:    req.FileSize,
		UploadedBy:  req.UploadedBy,
		EntityType:  req.EntityType,
		EntityID:    req.EntityID,
		Tags:        req.Tags,
	}

	media, err := h.commandHandler.HandleCreateMedia(r.Context(), cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(media)
}

func (h *MediaHandler) GetMedia(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	media, err := h.queryService.GetMedia(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(media)
}

func (h *MediaHandler) UpdateMedia(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req UpdateMediaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cmd := commands.UpdateMediaCommand{
		ID:          id,
		Title:       req.Title,
		Description: req.Description,
		Tags:        req.Tags,
	}

	media, err := h.commandHandler.HandleUpdateMedia(r.Context(), cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(media)
}

func (h *MediaHandler) DeleteMedia(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	cmd := commands.DeleteMediaCommand{ID: id}
	if err := h.commandHandler.HandleDeleteMedia(r.Context(), cmd); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *MediaHandler) ActivateMedia(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	cmd := commands.ActivateMediaCommand{ID: id}
	media, err := h.commandHandler.HandleActivateMedia(r.Context(), cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(media)
}

func (h *MediaHandler) ArchiveMedia(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	cmd := commands.ArchiveMediaCommand{ID: id}
	media, err := h.commandHandler.HandleArchiveMedia(r.Context(), cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(media)
}

func (h *MediaHandler) GetMediaByEntity(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	entityType := vars["entityType"]
	entityID := vars["entityId"]

	medias, err := h.queryService.GetMediaByEntityID(r.Context(), entityType, entityID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(medias)
}

func (h *MediaHandler) GetMediaByType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	mediaType := vars["type"]

	medias, err := h.queryService.GetMediaByType(r.Context(), mediaType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(medias)
}

func (h *MediaHandler) ListMedia(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 50
	offset := 0

	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	if offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	medias, err := h.queryService.ListMedia(r.Context(), limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(medias)
}
