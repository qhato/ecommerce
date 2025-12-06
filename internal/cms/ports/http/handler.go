package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/qhato/ecommerce/internal/cms/application/commands"
	"github.com/qhato/ecommerce/internal/cms/application/queries"
)

type CMSHandler struct {
	commandHandler *commands.CMSCommandHandler
	queryService   *queries.CMSQueryService
}

func NewCMSHandler(commandHandler *commands.CMSCommandHandler, queryService *queries.CMSQueryService) *CMSHandler {
	return &CMSHandler{
		commandHandler: commandHandler,
		queryService:   queryService,
	}
}

func (h *CMSHandler) RegisterRoutes(router *mux.Router) {
	// Content routes
	router.HandleFunc("/content", h.CreateContent).Methods("POST")
	router.HandleFunc("/content", h.GetAllContent).Methods("GET")
	router.HandleFunc("/content/{id}", h.GetContent).Methods("GET")
	router.HandleFunc("/content/{id}", h.UpdateContent).Methods("PUT")
	router.HandleFunc("/content/{id}", h.DeleteContent).Methods("DELETE")
	router.HandleFunc("/content/{id}/publish", h.PublishContent).Methods("POST")
	router.HandleFunc("/content/{id}/unpublish", h.UnpublishContent).Methods("POST")
	router.HandleFunc("/content/{id}/archive", h.ArchiveContent).Methods("POST")
	router.HandleFunc("/content/slug/{slug}", h.GetContentBySlug).Methods("GET")
	router.HandleFunc("/content/type/{type}", h.GetContentByType).Methods("GET")
	router.HandleFunc("/content/{id}/children", h.GetContentChildren).Methods("GET")
	router.HandleFunc("/content/search", h.SearchContent).Methods("GET")
	router.HandleFunc("/content/{id}/versions", h.GetContentVersions).Methods("GET")
	router.HandleFunc("/content/{id}/versions", h.CreateContentVersion).Methods("POST")

	// Media routes
	router.HandleFunc("/media", h.CreateMedia).Methods("POST")
	router.HandleFunc("/media", h.GetAllMedia).Methods("GET")
	router.HandleFunc("/media/{id}", h.GetMedia).Methods("GET")
	router.HandleFunc("/media/{id}", h.UpdateMedia).Methods("PUT")
	router.HandleFunc("/media/{id}", h.DeleteMedia).Methods("DELETE")
	router.HandleFunc("/media/uploader/{uploaderId}", h.GetMediaByUploader).Methods("GET")

	// Version routes
	router.HandleFunc("/versions/{id}", h.GetContentVersion).Methods("GET")
}

// Content Handlers

func (h *CMSHandler) CreateContent(w http.ResponseWriter, r *http.Request) {
	var cmd commands.CreateContentCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	content, err := h.commandHandler.HandleCreateContent(r.Context(), cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(content)
}

func (h *CMSHandler) GetContent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 64)

	content, err := h.queryService.GetContent(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(content)
}

func (h *CMSHandler) GetContentBySlug(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]
	locale := r.URL.Query().Get("locale")
	if locale == "" {
		locale = "en"
	}

	content, err := h.queryService.GetContentBySlug(r.Context(), slug, locale)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(content)
}

func (h *CMSHandler) GetContentByType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	contentType := vars["type"]
	locale := r.URL.Query().Get("locale")
	if locale == "" {
		locale = "en"
	}
	publishedOnly := r.URL.Query().Get("published_only") == "true"

	contents, err := h.queryService.GetContentByType(r.Context(), contentType, locale, publishedOnly)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(contents)
}

func (h *CMSHandler) GetAllContent(w http.ResponseWriter, r *http.Request) {
	locale := r.URL.Query().Get("locale")
	if locale == "" {
		locale = "en"
	}
	publishedOnly := r.URL.Query().Get("published_only") == "true"

	contents, err := h.queryService.GetAllContent(r.Context(), locale, publishedOnly)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(contents)
}

func (h *CMSHandler) GetContentChildren(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	parentID, _ := strconv.ParseInt(vars["id"], 10, 64)

	children, err := h.queryService.GetContentChildren(r.Context(), parentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(children)
}

func (h *CMSHandler) SearchContent(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	locale := r.URL.Query().Get("locale")
	if locale == "" {
		locale = "en"
	}
	publishedOnly := r.URL.Query().Get("published_only") == "true"

	contents, err := h.queryService.SearchContent(r.Context(), query, locale, publishedOnly)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(contents)
}

func (h *CMSHandler) UpdateContent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 64)

	var cmd commands.UpdateContentCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	cmd.ID = id

	content, err := h.commandHandler.HandleUpdateContent(r.Context(), cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(content)
}

func (h *CMSHandler) PublishContent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 64)

	content, err := h.commandHandler.HandlePublishContent(r.Context(), commands.PublishContentCommand{ID: id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(content)
}

func (h *CMSHandler) UnpublishContent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 64)

	content, err := h.commandHandler.HandleUnpublishContent(r.Context(), commands.UnpublishContentCommand{ID: id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(content)
}

func (h *CMSHandler) ArchiveContent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 64)

	content, err := h.commandHandler.HandleArchiveContent(r.Context(), commands.ArchiveContentCommand{ID: id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(content)
}

func (h *CMSHandler) DeleteContent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 64)

	if err := h.commandHandler.HandleDeleteContent(r.Context(), commands.DeleteContentCommand{ID: id}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Media Handlers

func (h *CMSHandler) CreateMedia(w http.ResponseWriter, r *http.Request) {
	var cmd commands.CreateMediaCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
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

func (h *CMSHandler) GetMedia(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 64)

	media, err := h.queryService.GetMedia(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(media)
}

func (h *CMSHandler) GetAllMedia(w http.ResponseWriter, r *http.Request) {
	mimeType := r.URL.Query().Get("mime_type")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit == 0 {
		limit = 100
	}

	medias, err := h.queryService.GetAllMedia(r.Context(), mimeType, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(medias)
}

func (h *CMSHandler) GetMediaByUploader(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uploaderID, _ := strconv.ParseInt(vars["uploaderId"], 10, 64)
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit == 0 {
		limit = 100
	}

	medias, err := h.queryService.GetMediaByUploader(r.Context(), uploaderID, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(medias)
}

func (h *CMSHandler) UpdateMedia(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 64)

	var cmd commands.UpdateMediaCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	cmd.ID = id

	media, err := h.commandHandler.HandleUpdateMedia(r.Context(), cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(media)
}

func (h *CMSHandler) DeleteMedia(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 64)

	if err := h.commandHandler.HandleDeleteMedia(r.Context(), commands.DeleteMediaCommand{ID: id}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Content Version Handlers

func (h *CMSHandler) GetContentVersions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	contentID, _ := strconv.ParseInt(vars["id"], 10, 64)

	versions, err := h.queryService.GetContentVersions(r.Context(), contentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(versions)
}

func (h *CMSHandler) GetContentVersion(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	versionID, _ := strconv.ParseInt(vars["id"], 10, 64)

	version, err := h.queryService.GetContentVersion(r.Context(), versionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(version)
}

func (h *CMSHandler) CreateContentVersion(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	contentID, _ := strconv.ParseInt(vars["id"], 10, 64)

	var cmd commands.CreateContentVersionCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	cmd.ContentID = contentID

	version, err := h.commandHandler.HandleCreateContentVersion(r.Context(), cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(version)
}
