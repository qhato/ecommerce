package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/qhato/ecommerce/internal/menu/application/commands"
	"github.com/qhato/ecommerce/internal/menu/application/queries"
)

type MenuHandler struct {
	commandHandler *commands.MenuCommandHandler
	queryService   *queries.MenuQueryService
}

func NewMenuHandler(
	commandHandler *commands.MenuCommandHandler,
	queryService *queries.MenuQueryService,
) *MenuHandler {
	return &MenuHandler{
		commandHandler: commandHandler,
		queryService:   queryService,
	}
}

func (h *MenuHandler) RegisterRoutes(router *mux.Router) {
	// Menu routes
	router.HandleFunc("/menus", h.CreateMenu).Methods("POST")
	router.HandleFunc("/menus/{id}", h.GetMenu).Methods("GET")
	router.HandleFunc("/menus/{id}", h.UpdateMenu).Methods("PUT")
	router.HandleFunc("/menus/{id}", h.DeleteMenu).Methods("DELETE")
	router.HandleFunc("/menus/slug/{slug}", h.GetMenuBySlug).Methods("GET")
	router.HandleFunc("/menus/location/{location}", h.GetMenuByLocation).Methods("GET")
	router.HandleFunc("/menus/type/{type}", h.GetMenusByType).Methods("GET")
	router.HandleFunc("/menus", h.ListMenus).Methods("GET")

	// Menu item routes
	router.HandleFunc("/menu-items", h.CreateMenuItem).Methods("POST")
	router.HandleFunc("/menu-items/{id}", h.GetMenuItem).Methods("GET")
	router.HandleFunc("/menu-items/{id}", h.UpdateMenuItem).Methods("PUT")
	router.HandleFunc("/menu-items/{id}", h.DeleteMenuItem).Methods("DELETE")
	router.HandleFunc("/menu-items/{id}/move", h.MoveMenuItem).Methods("POST")
	router.HandleFunc("/menus/{menuId}/items", h.GetMenuItems).Methods("GET")
	router.HandleFunc("/menus/{menuId}/tree", h.GetMenuTree).Methods("GET")
}

type CreateMenuRequest struct {
	Name        string `json:"name" validate:"required"`
	Slug        string `json:"slug" validate:"required"`
	Description string `json:"description"`
	Location    string `json:"location"`
}

type UpdateMenuRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	Location    string `json:"location"`
}

type CreateMenuItemRequest struct {
	MenuID      int64   `json:"menu_id" validate:"required"`
	ParentID    *int64  `json:"parent_id"`
	Title       string  `json:"title" validate:"required"`
	URL         string  `json:"url"`
	Target      string  `json:"target"`
	Icon        string  `json:"icon"`
	CSSClass    string  `json:"css_class"`
	SortOrder   int     `json:"sort_order"`
	Permissions *string `json:"permissions"`
}

type UpdateMenuItemRequest struct {
	Title       string  `json:"title" validate:"required"`
	URL         string  `json:"url"`
	Target      string  `json:"target"`
	Icon        string  `json:"icon"`
	CSSClass    string  `json:"css_class"`
	SortOrder   int     `json:"sort_order"`
	Permissions *string `json:"permissions"`
}

type MoveMenuItemRequest struct {
	ParentID  *int64 `json:"parent_id"`
	SortOrder int    `json:"sort_order"`
}

func (h *MenuHandler) CreateMenu(w http.ResponseWriter, r *http.Request) {
	var req CreateMenuRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cmd := commands.CreateMenuCommand{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		Location:    req.Location,
	}

	menu, err := h.commandHandler.HandleCreateMenu(r.Context(), cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(menu)
}

func (h *MenuHandler) GetMenu(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	menu, err := h.queryService.GetMenu(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(menu)
}

func (h *MenuHandler) UpdateMenu(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var req UpdateMenuRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cmd := commands.UpdateMenuCommand{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		Location:    req.Location,
	}

	menu, err := h.commandHandler.HandleUpdateMenu(r.Context(), cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(menu)
}

func (h *MenuHandler) DeleteMenu(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	cmd := commands.DeleteMenuCommand{ID: id}
	if err := h.commandHandler.HandleDeleteMenu(r.Context(), cmd); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *MenuHandler) GetMenuBySlug(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]

	menu, err := h.queryService.GetMenuBySlug(r.Context(), slug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(menu)
}

func (h *MenuHandler) GetMenuByLocation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	location := vars["location"]

	menu, err := h.queryService.GetMenuByLocation(r.Context(), location)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(menu)
}

func (h *MenuHandler) GetMenusByType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	menuType := vars["type"]

	menus, err := h.queryService.GetMenusByType(r.Context(), menuType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(menus)
}

func (h *MenuHandler) ListMenus(w http.ResponseWriter, r *http.Request) {
	activeOnly := r.URL.Query().Get("active") == "true"

	menus, err := h.queryService.GetAllMenus(r.Context(), activeOnly)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(menus)
}

func (h *MenuHandler) CreateMenuItem(w http.ResponseWriter, r *http.Request) {
	var req CreateMenuItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cmd := commands.CreateMenuItemCommand{
		MenuID:      req.MenuID,
		ParentID:    req.ParentID,
		Title:       req.Title,
		URL:         req.URL,
		Target:      req.Target,
		Icon:        req.Icon,
		CSSClass:    req.CSSClass,
		SortOrder:   req.SortOrder,
		Permissions: req.Permissions,
	}

	item, err := h.commandHandler.HandleCreateMenuItem(r.Context(), cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

func (h *MenuHandler) GetMenuItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	item, err := h.queryService.GetMenuItem(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

func (h *MenuHandler) UpdateMenuItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var req UpdateMenuItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cmd := commands.UpdateMenuItemCommand{
		ID:          id,
		Title:       req.Title,
		URL:         req.URL,
		Target:      req.Target,
		Icon:        req.Icon,
		CSSClass:    req.CSSClass,
		SortOrder:   req.SortOrder,
		Permissions: req.Permissions,
	}

	item, err := h.commandHandler.HandleUpdateMenuItem(r.Context(), cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

func (h *MenuHandler) DeleteMenuItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	cmd := commands.DeleteMenuItemCommand{ID: id}
	if err := h.commandHandler.HandleDeleteMenuItem(r.Context(), cmd); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *MenuHandler) MoveMenuItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var req MoveMenuItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cmd := commands.MoveMenuItemCommand{
		ID:        id,
		ParentID:  req.ParentID,
		SortOrder: req.SortOrder,
	}

	item, err := h.commandHandler.HandleMoveMenuItem(r.Context(), cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

func (h *MenuHandler) GetMenuItems(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	menuID, err := strconv.ParseInt(vars["menuId"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid menu ID", http.StatusBadRequest)
		return
	}

	items, err := h.queryService.GetMenuItems(r.Context(), menuID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func (h *MenuHandler) GetMenuTree(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	menuID, err := strconv.ParseInt(vars["menuId"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid menu ID", http.StatusBadRequest)
		return
	}

	tree, err := h.queryService.GetMenuTree(r.Context(), menuID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tree)
}
