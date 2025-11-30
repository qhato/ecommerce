package http

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/qhato/ecommerce/internal/catalog/application/commands"
	"github.com/qhato/ecommerce/internal/catalog/application/queries"
	pkghttp "github.com/qhato/ecommerce/pkg/http"
	"github.com/qhato/ecommerce/pkg/logger"
)

// AdminCategoryHandler handles admin category HTTP requests
type AdminCategoryHandler struct {
	commandHandler *commands.CategoryCommandHandler
	queryHandler   *queries.CategoryQueryHandler
	logger         *logger.Logger
}

// NewAdminCategoryHandler creates a new admin category handler
func NewAdminCategoryHandler(
	commandHandler *commands.CategoryCommandHandler,
	queryHandler *queries.CategoryQueryHandler,
	logger *logger.Logger,
) *AdminCategoryHandler {
	return &AdminCategoryHandler{
		commandHandler: commandHandler,
		queryHandler:   queryHandler,
		logger:         logger,
	}
}

// RegisterRoutes registers admin category routes
func (h *AdminCategoryHandler) RegisterRoutes(r chi.Router) {
	r.Route("/admin/categories", func(r chi.Router) {
		r.Post("/", h.CreateCategory)
		r.Get("/", h.ListCategories)
		r.Get("/root", h.ListRootCategories)
		r.Get("/{id}", h.GetCategory)
		r.Put("/{id}", h.UpdateCategory)
		r.Delete("/{id}", h.DeleteCategory)
		r.Get("/{id}/children", h.ListChildCategories)
		r.Get("/{id}/path", h.GetCategoryPath)
	})
}

// CreateCategory creates a new category
func (h *AdminCategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var cmd commands.CreateCategoryCommand
	if err := pkghttp.DecodeJSON(r, &cmd); err != nil {
		pkghttp.RespondError(w, err)
		return
	}

	categoryID, err := h.commandHandler.HandleCreateCategory(r.Context(), &cmd)
	if err != nil {
		h.logger.WithError(err).Error("failed to create category")
		pkghttp.RespondError(w, err)
		return
	}

	pkghttp.RespondJSON(w, http.StatusCreated, map[string]interface{}{
		"id": categoryID,
	})
}

// ListCategories lists all categories with pagination
func (h *AdminCategoryHandler) ListCategories(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	includeArchived := r.URL.Query().Get("include_archived") == "true"
	activeOnly := r.URL.Query().Get("active_only") == "true"
	sortBy := r.URL.Query().Get("sort_by")
	sortOrder := r.URL.Query().Get("sort_order")

	query := &queries.ListCategoriesQuery{
		Page:            page,
		PageSize:        pageSize,
		IncludeArchived: includeArchived,
		ActiveOnly:      activeOnly,
		SortBy:          sortBy,
		SortOrder:       sortOrder,
	}

	result, err := h.queryHandler.HandleListCategories(r.Context(), query)
	if err != nil {
		h.logger.WithError(err).Error("failed to list categories")
		pkghttp.RespondError(w, err)
		return
	}

	pkghttp.RespondJSON(w, http.StatusOK, result)
}

// ListRootCategories lists root categories
func (h *AdminCategoryHandler) ListRootCategories(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	includeArchived := r.URL.Query().Get("include_archived") == "true"
	activeOnly := r.URL.Query().Get("active_only") == "true"
	sortBy := r.URL.Query().Get("sort_by")
	sortOrder := r.URL.Query().Get("sort_order")

	query := &queries.ListRootCategoriesQuery{
		Page:            page,
		PageSize:        pageSize,
		IncludeArchived: includeArchived,
		ActiveOnly:      activeOnly,
		SortBy:          sortBy,
		SortOrder:       sortOrder,
	}

	result, err := h.queryHandler.HandleListRootCategories(r.Context(), query)
	if err != nil {
		h.logger.WithError(err).Error("failed to list root categories")
		pkghttp.RespondError(w, err)
		return
	}

	pkghttp.RespondJSON(w, http.StatusOK, result)
}

// GetCategory retrieves a category by ID
func (h *AdminCategoryHandler) GetCategory(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		pkghttp.RespondError(w, pkghttp.NewValidationError("invalid category ID"))
		return
	}

	query := &queries.GetCategoryByIDQuery{ID: id}
	category, err := h.queryHandler.HandleGetCategoryByID(r.Context(), query)
	if err != nil {
		h.logger.WithError(err).WithField("category_id", id).Error("failed to get category")
		pkghttp.RespondError(w, err)
		return
	}

	pkghttp.RespondJSON(w, http.StatusOK, category)
}

// UpdateCategory updates an existing category
func (h *AdminCategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		pkghttp.RespondError(w, pkghttp.NewValidationError("invalid category ID"))
		return
	}

	var cmd commands.UpdateCategoryCommand
	if err := pkghttp.DecodeJSON(r, &cmd); err != nil {
		pkghttp.RespondError(w, err)
		return
	}
	cmd.ID = id

	if err := h.commandHandler.HandleUpdateCategory(r.Context(), &cmd); err != nil {
		h.logger.WithError(err).WithField("category_id", id).Error("failed to update category")
		pkghttp.RespondError(w, err)
		return
	}

	pkghttp.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "category updated successfully",
	})
}

// DeleteCategory deletes a category
func (h *AdminCategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		pkghttp.RespondError(w, pkghttp.NewValidationError("invalid category ID"))
		return
	}

	cmd := &commands.DeleteCategoryCommand{ID: id}
	if err := h.commandHandler.HandleDeleteCategory(r.Context(), cmd); err != nil {
		h.logger.WithError(err).WithField("category_id", id).Error("failed to delete category")
		pkghttp.RespondError(w, err)
		return
	}

	pkghttp.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "category deleted successfully",
	})
}

// ListChildCategories lists child categories of a parent
func (h *AdminCategoryHandler) ListChildCategories(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		pkghttp.RespondError(w, pkghttp.NewValidationError("invalid category ID"))
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	includeArchived := r.URL.Query().Get("include_archived") == "true"
	activeOnly := r.URL.Query().Get("active_only") == "true"
	sortBy := r.URL.Query().Get("sort_by")
	sortOrder := r.URL.Query().Get("sort_order")

	query := &queries.ListCategoriesByParentQuery{
		ParentID:        id,
		Page:            page,
		PageSize:        pageSize,
		IncludeArchived: includeArchived,
		ActiveOnly:      activeOnly,
		SortBy:          sortBy,
		SortOrder:       sortOrder,
	}

	result, err := h.queryHandler.HandleListCategoriesByParent(r.Context(), query)
	if err != nil {
		h.logger.WithError(err).WithField("parent_id", id).Error("failed to list child categories")
		pkghttp.RespondError(w, err)
		return
	}

	pkghttp.RespondJSON(w, http.StatusOK, result)
}

// GetCategoryPath retrieves the full path from root to category
func (h *AdminCategoryHandler) GetCategoryPath(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		pkghttp.RespondError(w, pkghttp.NewValidationError("invalid category ID"))
		return
	}

	query := &queries.GetCategoryPathQuery{CategoryID: id}
	path, err := h.queryHandler.HandleGetCategoryPath(r.Context(), query)
	if err != nil {
		h.logger.WithError(err).WithField("category_id", id).Error("failed to get category path")
		pkghttp.RespondError(w, err)
		return
	}

	pkghttp.RespondJSON(w, http.StatusOK, path)
}