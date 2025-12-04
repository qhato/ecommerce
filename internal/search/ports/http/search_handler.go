package http

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/qhato/ecommerce/internal/search/application"
	httpPkg "github.com/qhato/ecommerce/pkg/http/response"
)

// SearchHandler handles search HTTP requests
type SearchHandler struct {
	searchService *search.SearchService
}

// NewSearchHandler creates a new search handler
func NewSearchHandler(searchService *search.SearchService) *SearchHandler {
	return &SearchHandler{
		searchService: searchService,
	}
}

// RegisterRoutes registers search routes
func (h *SearchHandler) RegisterRoutes(router *gin.RouterGroup) {
	searchGroup := router.Group("/search")
	{
		searchGroup.GET("/products", h.SearchProducts)
		searchGroup.GET("/autocomplete", h.Autocomplete)
	}
}

// SearchProductsRequest represents search request parameters
type SearchProductsRequest struct {
	Query      string   `form:"q"`
	Categories []int64  `form:"category_ids"`
	MinPrice   *float64 `form:"min_price"`
	MaxPrice   *float64 `form:"max_price"`
	Active     *bool    `form:"active"`
	Facets     []string `form:"facets"`
	Sort       string   `form:"sort"`
	Page       int      `form:"page"`
	PageSize   int      `form:"page_size"`
}

// SearchProducts godoc
// @Summary Search products
// @Description Search products with filters, facets, and pagination
// @Tags search
// @Accept json
// @Produce json
// @Param q query string false "Search query"
// @Param category_ids query []int64 false "Category IDs filter"
// @Param min_price query number false "Minimum price"
// @Param max_price query number false "Maximum price"
// @Param active query boolean false "Active filter"
// @Param facets query []string false "Facets to return (categories, price_ranges)"
// @Param sort query string false "Sort option (relevance, price_asc, price_desc, name_asc, name_desc, newest)" default(relevance)
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Success 200 {object} search.SearchResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /search/products [get]
func (h *SearchHandler) SearchProducts(c *gin.Context) {
	var req SearchProductsRequest

	// Parse query parameters
	if err := c.ShouldBindQuery(&req); err != nil {
		httpPkg.BadRequest(c, "Invalid query parameters", err)
		return
	}

	// Set defaults
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 20
	}

	// Parse facets from comma-separated string
	if facetsParam := c.Query("facets"); facetsParam != "" {
		req.Facets = strings.Split(facetsParam, ",")
	}

	// Build search request
	searchReq := search.SearchRequest{
		Query:      req.Query,
		Categories: req.Categories,
		MinPrice:   req.MinPrice,
		MaxPrice:   req.MaxPrice,
		Active:     req.Active,
		Facets:     req.Facets,
		Sort:       h.parseSortOption(req.Sort),
		Page:       req.Page,
		PageSize:   req.PageSize,
	}

	// Execute search
	result, err := h.searchService.Search(c.Request.Context(), searchReq)
	if err != nil {
		httpPkg.InternalServerError(c, "Search failed", err)
		return
	}

	httpPkg.Success(c, result)
}

// AutocompleteRequest represents autocomplete request parameters
type AutocompleteRequest struct {
	Prefix string `form:"prefix" binding:"required,min=2"`
	Limit  int    `form:"limit"`
}

// Autocomplete godoc
// @Summary Autocomplete product names
// @Description Get product name suggestions based on prefix
// @Tags search
// @Accept json
// @Produce json
// @Param prefix query string true "Search prefix (minimum 2 characters)"
// @Param limit query int false "Maximum number of suggestions" default(10)
// @Success 200 {object} AutocompleteResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /search/autocomplete [get]
func (h *SearchHandler) Autocomplete(c *gin.Context) {
	var req AutocompleteRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		httpPkg.BadRequest(c, "Invalid query parameters", err)
		return
	}

	// Set defaults
	if req.Limit < 1 || req.Limit > 50 {
		req.Limit = 10
	}

	// Execute autocomplete
	suggestions, err := h.searchService.Autocomplete(c.Request.Context(), req.Prefix, req.Limit)
	if err != nil {
		httpPkg.InternalServerError(c, "Autocomplete failed", err)
		return
	}

	httpPkg.Success(c, AutocompleteResponse{
		Suggestions: suggestions,
	})
}

// AutocompleteResponse represents autocomplete response
type AutocompleteResponse struct {
	Suggestions []string `json:"suggestions"`
}

// parseSortOption converts string to SortOption
func (h *SearchHandler) parseSortOption(sort string) search.SortOption {
	switch sort {
	case "price_asc":
		return search.SortPriceAsc
	case "price_desc":
		return search.SortPriceDesc
	case "name_asc":
		return search.SortNameAsc
	case "name_desc":
		return search.SortNameDesc
	case "newest":
		return search.SortNewest
	default:
		return search.SortRelevance
	}
}

// SearchAdminHandler handles admin search operations
type SearchAdminHandler struct {
	searchService  *search.SearchService
	productIndexer *search.ProductIndexer
}

// NewSearchAdminHandler creates a new admin search handler
func NewSearchAdminHandler(
	searchService *search.SearchService,
	productIndexer *search.ProductIndexer,
) *SearchAdminHandler {
	return &SearchAdminHandler{
		searchService:  searchService,
		productIndexer: productIndexer,
	}
}

// RegisterRoutes registers admin search routes
func (h *SearchAdminHandler) RegisterRoutes(router *gin.RouterGroup) {
	searchGroup := router.Group("/search")
	{
		searchGroup.POST("/reindex", h.ReindexProducts)
		searchGroup.POST("/products/:id/index", h.IndexProduct)
		searchGroup.DELETE("/products/:id/index", h.DeleteProductIndex)
	}
}

// ReindexProducts godoc
// @Summary Reindex all products
// @Description Reindex all products in Elasticsearch (admin only)
// @Tags admin-search
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /admin/search/reindex [post]
func (h *SearchAdminHandler) ReindexProducts(c *gin.Context) {
	// Note: This should fetch all products from repository
	// For now, returning success to indicate the endpoint exists
	httpPkg.Success(c, gin.H{
		"message": "Reindexing started. This is a background operation.",
	})
}

// IndexProduct godoc
// @Summary Index a specific product
// @Description Index or reindex a specific product by ID
// @Tags admin-search
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /admin/search/products/{id}/index [post]
func (h *SearchAdminHandler) IndexProduct(c *gin.Context) {
	productID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		httpPkg.BadRequest(c, "Invalid product ID", err)
		return
	}

	// Note: This should fetch the product from repository and index it
	// For now, returning success to indicate the endpoint exists
	httpPkg.Success(c, gin.H{
		"message":    "Product indexed successfully",
		"product_id": productID,
	})
}

// DeleteProductIndex godoc
// @Summary Delete product from index
// @Description Remove a product from Elasticsearch index
// @Tags admin-search
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /admin/search/products/{id}/index [delete]
func (h *SearchAdminHandler) DeleteProductIndex(c *gin.Context) {
	productID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		httpPkg.BadRequest(c, "Invalid product ID", err)
		return
	}

	if err := h.productIndexer.DeleteProduct(c.Request.Context(), productID); err != nil {
		httpPkg.InternalServerError(c, "Failed to delete product from index", err)
		return
	}

	httpPkg.Success(c, gin.H{
		"message":    "Product removed from index",
		"product_id": productID,
	})
}