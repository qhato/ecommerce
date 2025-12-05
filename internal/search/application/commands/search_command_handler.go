package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/shopspring/decimal"

	"github.com/qhato/ecommerce/internal/search/domain"
	"github.com/qhato/ecommerce/pkg/events"
	"github.com/qhato/ecommerce/pkg/logger"
)

// SearchIndexer defines the interface for indexing operations
type SearchIndexer interface {
	IndexProduct(ctx context.Context, doc *domain.ProductSearchDocument) error
	BulkIndexProducts(ctx context.Context, docs []*domain.ProductSearchDocument) error
	DeleteProduct(ctx context.Context, productID int64) error
	ReindexAll(ctx context.Context) error
}

// ProductRepository defines the interface for product data access
type ProductRepository interface {
	FindAllForIndexing(ctx context.Context) ([]IndexProductCommand, error)
	FindByCategoryForIndexing(ctx context.Context, categoryID int64) ([]IndexProductCommand, error)
}

// SearchCommandHandler handles search-related commands
type SearchCommandHandler struct {
	indexer            SearchIndexer
	productRepo        ProductRepository
	synonymRepo        domain.SearchSynonymRepository
	redirectRepo       domain.SearchRedirectRepository
	facetConfigRepo    domain.SearchFacetConfigRepository
	indexingJobRepo    domain.IndexingJobRepository
	eventBus           events.EventBus
	logger             logger.Logger
}

// NewSearchCommandHandler creates a new search command handler
func NewSearchCommandHandler(
	indexer SearchIndexer,
	productRepo ProductRepository,
	synonymRepo domain.SearchSynonymRepository,
	redirectRepo domain.SearchRedirectRepository,
	facetConfigRepo domain.SearchFacetConfigRepository,
	indexingJobRepo domain.IndexingJobRepository,
	eventBus events.EventBus,
	logger logger.Logger,
) *SearchCommandHandler {
	return &SearchCommandHandler{
		indexer:         indexer,
		productRepo:     productRepo,
		synonymRepo:     synonymRepo,
		redirectRepo:    redirectRepo,
		facetConfigRepo: facetConfigRepo,
		indexingJobRepo: indexingJobRepo,
		eventBus:        eventBus,
		logger:          logger,
	}
}

// HandleIndexProduct indexes a single product
func (h *SearchCommandHandler) HandleIndexProduct(ctx context.Context, cmd *IndexProductCommand) error {
	doc := h.commandToDocument(cmd)

	if err := h.indexer.IndexProduct(ctx, doc); err != nil {
		h.logger.Error("Failed to index product",
			logger.Field{Key: "product_id", Value: cmd.ProductID},
			logger.Field{Key: "error", Value: err.Error()},
		)
		return fmt.Errorf("failed to index product: %w", err)
	}

	h.logger.Info("Product indexed successfully",
		logger.Field{Key: "product_id", Value: cmd.ProductID},
		logger.Field{Key: "sku", Value: cmd.SKU},
	)

	// Publish event
	h.eventBus.Publish(ctx, "product.indexed", map[string]interface{}{
		"product_id": cmd.ProductID,
		"sku":        cmd.SKU,
	})

	return nil
}

// HandleBulkIndexProducts indexes multiple products
func (h *SearchCommandHandler) HandleBulkIndexProducts(ctx context.Context, cmd *BulkIndexProductsCommand) error {
	docs := make([]*domain.ProductSearchDocument, len(cmd.Products))
	for i, product := range cmd.Products {
		docs[i] = h.commandToDocument(&product)
	}

	if err := h.indexer.BulkIndexProducts(ctx, docs); err != nil {
		h.logger.Error("Failed to bulk index products",
			logger.Field{Key: "count", Value: len(cmd.Products)},
			logger.Field{Key: "error", Value: err.Error()},
		)
		return fmt.Errorf("failed to bulk index products: %w", err)
	}

	h.logger.Info("Products bulk indexed successfully",
		logger.Field{Key: "count", Value: len(cmd.Products)},
	)

	return nil
}

// HandleDeleteProductFromIndex deletes a product from the index
func (h *SearchCommandHandler) HandleDeleteProductFromIndex(ctx context.Context, cmd *DeleteProductFromIndexCommand) error {
	if err := h.indexer.DeleteProduct(ctx, cmd.ProductID); err != nil {
		h.logger.Error("Failed to delete product from index",
			logger.Field{Key: "product_id", Value: cmd.ProductID},
			logger.Field{Key: "error", Value: err.Error()},
		)
		return fmt.Errorf("failed to delete product from index: %w", err)
	}

	h.logger.Info("Product deleted from index",
		logger.Field{Key: "product_id", Value: cmd.ProductID},
	)

	// Publish event
	h.eventBus.Publish(ctx, "product.deindexed", map[string]interface{}{
		"product_id": cmd.ProductID,
	})

	return nil
}

// HandleReindexAllProducts reindexes all products
func (h *SearchCommandHandler) HandleReindexAllProducts(ctx context.Context, cmd *ReindexAllProductsCommand) (int64, error) {
	// Create indexing job
	job := &domain.IndexingJob{
		Type:       domain.IndexingJobTypeFull,
		Status:     domain.IndexingJobStatusPending,
		EntityType: "product",
		CreatedBy:  cmd.CreatedBy,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := h.indexingJobRepo.Create(job); err != nil {
		return 0, fmt.Errorf("failed to create indexing job: %w", err)
	}

	// Start job asynchronously
	go h.executeFullReindex(context.Background(), job)

	return job.ID, nil
}

// HandleReindexCategory reindexes products in a category
func (h *SearchCommandHandler) HandleReindexCategory(ctx context.Context, cmd *ReindexCategoryCommand) (int64, error) {
	// Create indexing job
	job := &domain.IndexingJob{
		Type:       domain.IndexingJobTypeIncremental,
		Status:     domain.IndexingJobStatusPending,
		EntityType: "product",
		CreatedBy:  cmd.CreatedBy,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := h.indexingJobRepo.Create(job); err != nil {
		return 0, fmt.Errorf("failed to create indexing job: %w", err)
	}

	// Start job asynchronously
	go h.executeCategoryReindex(context.Background(), job, cmd.CategoryID)

	return job.ID, nil
}

// executeFullReindex executes a full reindex
func (h *SearchCommandHandler) executeFullReindex(ctx context.Context, job *domain.IndexingJob) {
	job.Start()
	h.indexingJobRepo.Update(ctx, job)

	h.logger.Info("Starting full reindex",
		logger.Field{Key: "job_id", Value: job.ID},
	)

	// Get all products
	products, err := h.productRepo.FindAllForIndexing(ctx)
	if err != nil {
		job.Fail(fmt.Sprintf("failed to fetch products: %v", err))
		h.indexingJobRepo.Update(ctx, job)
		h.logger.Error("Failed to fetch products for reindex",
			logger.Field{Key: "job_id", Value: job.ID},
			logger.Field{Key: "error", Value: err.Error()},
		)
		return
	}

	job.TotalItems = len(products)
	h.indexingJobRepo.Update(ctx, job)

	// Index products in batches
	batchSize := 100
	for i := 0; i < len(products); i += batchSize {
		end := i + batchSize
		if end > len(products) {
			end = len(products)
		}

		batch := products[i:end]
		bulkCmd := &BulkIndexProductsCommand{Products: batch}

		if err := h.HandleBulkIndexProducts(ctx, bulkCmd); err != nil {
			for range batch {
				job.IncrementFailed()
			}
			h.logger.Error("Failed to index batch",
				logger.Field{Key: "job_id", Value: job.ID},
				logger.Field{Key: "batch_start", Value: i},
				logger.Field{Key: "error", Value: err.Error()},
			)
		} else {
			for range batch {
				job.IncrementProcessed()
			}
		}

		h.indexingJobRepo.Update(ctx, job)
	}

	job.Complete()
	h.indexingJobRepo.Update(ctx, job)

	h.logger.Info("Full reindex completed",
		logger.Field{Key: "job_id", Value: job.ID},
		logger.Field{Key: "processed", Value: job.ProcessedItems},
		logger.Field{Key: "failed", Value: job.FailedItems},
	)
}

// executeCategoryReindex executes a category reindex
func (h *SearchCommandHandler) executeCategoryReindex(ctx context.Context, job *domain.IndexingJob, categoryID int64) {
	job.Start()
	h.indexingJobRepo.Update(ctx, job)

	h.logger.Info("Starting category reindex",
		logger.Field{Key: "job_id", Value: job.ID},
		logger.Field{Key: "category_id", Value: categoryID},
	)

	// Get products in category
	products, err := h.productRepo.FindByCategoryForIndexing(ctx, categoryID)
	if err != nil {
		job.Fail(fmt.Sprintf("failed to fetch products: %v", err))
		h.indexingJobRepo.Update(ctx, job)
		return
	}

	job.TotalItems = len(products)
	h.indexingJobRepo.Update(ctx, job)

	// Index products
	bulkCmd := &BulkIndexProductsCommand{Products: products}
	if err := h.HandleBulkIndexProducts(ctx, bulkCmd); err != nil {
		job.Fail(err.Error())
	} else {
		job.ProcessedItems = len(products)
		job.Complete()
	}

	h.indexingJobRepo.Update(ctx, job)
}

// HandleCreateSynonym creates a search synonym
func (h *SearchCommandHandler) HandleCreateSynonym(ctx context.Context, cmd *CreateSynonymCommand) (int64, error) {
	synonym := &domain.SearchSynonym{
		Term:      cmd.Term,
		Synonyms:  cmd.Synonyms,
		IsActive:  cmd.IsActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := h.synonymRepo.Create(synonym); err != nil {
		return 0, fmt.Errorf("failed to create synonym: %w", err)
	}

	return synonym.ID, nil
}

// HandleUpdateSynonym updates a search synonym
func (h *SearchCommandHandler) HandleUpdateSynonym(ctx context.Context, cmd *UpdateSynonymCommand) error {
	synonym, err := h.synonymRepo.FindByID(cmd.ID)
	if err != nil {
		return fmt.Errorf("synonym not found: %w", err)
	}

	synonym.Term = cmd.Term
	synonym.Synonyms = cmd.Synonyms
	synonym.IsActive = cmd.IsActive
	synonym.UpdatedAt = time.Now()

	return h.synonymRepo.Update(synonym)
}

// HandleDeleteSynonym deletes a search synonym
func (h *SearchCommandHandler) HandleDeleteSynonym(ctx context.Context, cmd *DeleteSynonymCommand) error {
	return h.synonymRepo.Delete(cmd.ID)
}

// HandleCreateRedirect creates a search redirect
func (h *SearchCommandHandler) HandleCreateRedirect(ctx context.Context, cmd *CreateRedirectCommand) (int64, error) {
	redirect := &domain.SearchRedirect{
		SearchTerm:     cmd.SearchTerm,
		TargetURL:      cmd.TargetURL,
		Priority:       cmd.Priority,
		IsActive:       cmd.IsActive,
		ActivationDate: cmd.ActivationDate,
		ExpirationDate: cmd.ExpirationDate,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := h.redirectRepo.Create(redirect); err != nil {
		return 0, fmt.Errorf("failed to create redirect: %w", err)
	}

	return redirect.ID, nil
}

// HandleUpdateRedirect updates a search redirect
func (h *SearchCommandHandler) HandleUpdateRedirect(ctx context.Context, cmd *UpdateRedirectCommand) error {
	redirect, err := h.redirectRepo.FindByID(cmd.ID)
	if err != nil {
		return fmt.Errorf("redirect not found: %w", err)
	}

	redirect.SearchTerm = cmd.SearchTerm
	redirect.TargetURL = cmd.TargetURL
	redirect.Priority = cmd.Priority
	redirect.IsActive = cmd.IsActive
	redirect.ActivationDate = cmd.ActivationDate
	redirect.ExpirationDate = cmd.ExpirationDate
	redirect.UpdatedAt = time.Now()

	return h.redirectRepo.Update(redirect)
}

// HandleDeleteRedirect deletes a search redirect
func (h *SearchCommandHandler) HandleDeleteRedirect(ctx context.Context, cmd *DeleteRedirectCommand) error {
	return h.redirectRepo.Delete(cmd.ID)
}

// HandleCreateFacetConfig creates a facet configuration
func (h *SearchCommandHandler) HandleCreateFacetConfig(ctx context.Context, cmd *CreateFacetConfigCommand) (int64, error) {
	config := &domain.SearchFacetConfig{
		Name:             cmd.Name,
		Label:            cmd.Label,
		FieldName:        cmd.FieldName,
		FacetType:        domain.FacetType(cmd.FacetType),
		IsActive:         cmd.IsActive,
		ShowInResults:    cmd.ShowInResults,
		ShowInNavigation: cmd.ShowInNavigation,
		Priority:         cmd.Priority,
		MinDocCount:      cmd.MinDocCount,
		MaxValues:        cmd.MaxValues,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if err := h.facetConfigRepo.Create(config); err != nil {
		return 0, fmt.Errorf("failed to create facet config: %w", err)
	}

	return config.ID, nil
}

// HandleUpdateFacetConfig updates a facet configuration
func (h *SearchCommandHandler) HandleUpdateFacetConfig(ctx context.Context, cmd *UpdateFacetConfigCommand) error {
	config, err := h.facetConfigRepo.FindByID(cmd.ID)
	if err != nil {
		return fmt.Errorf("facet config not found: %w", err)
	}

	config.Name = cmd.Name
	config.Label = cmd.Label
	config.FieldName = cmd.FieldName
	config.FacetType = domain.FacetType(cmd.FacetType)
	config.IsActive = cmd.IsActive
	config.ShowInResults = cmd.ShowInResults
	config.ShowInNavigation = cmd.ShowInNavigation
	config.Priority = cmd.Priority
	config.MinDocCount = cmd.MinDocCount
	config.MaxValues = cmd.MaxValues
	config.UpdatedAt = time.Now()

	return h.facetConfigRepo.Update(config)
}

// HandleDeleteFacetConfig deletes a facet configuration
func (h *SearchCommandHandler) HandleDeleteFacetConfig(ctx context.Context, cmd *DeleteFacetConfigCommand) error {
	return h.facetConfigRepo.Delete(cmd.ID)
}

// commandToDocument converts IndexProductCommand to ProductSearchDocument
func (h *SearchCommandHandler) commandToDocument(cmd *IndexProductCommand) *domain.ProductSearchDocument {
	var salePrice *decimal.Decimal
	if cmd.SalePrice != nil {
		sp := decimal.NewFromFloat(*cmd.SalePrice)
		salePrice = &sp
	}

	return &domain.ProductSearchDocument{
		ProductID:    cmd.ProductID,
		SKU:          cmd.SKU,
		Name:         cmd.Name,
		Description:  cmd.Description,
		LongDesc:     cmd.LongDesc,
		Tags:         cmd.Tags,
		Price:        decimal.NewFromFloat(cmd.Price),
		SalePrice:    salePrice,
		OnSale:       cmd.OnSale,
		CategoryID:   cmd.CategoryID,
		CategoryName: cmd.CategoryName,
		CategoryPath: cmd.CategoryPath,
		Brand:        cmd.Brand,
		Color:        cmd.Color,
		Size:         cmd.Size,
		Attributes:   cmd.Attributes,
		IsAvailable:  cmd.IsAvailable,
		StockLevel:   cmd.StockLevel,
		ImageURL:     cmd.ImageURL,
		ThumbnailURL: cmd.ThumbnailURL,
		IsActive:     cmd.IsActive,
		IsFeatured:   cmd.IsFeatured,
		Rating:       cmd.Rating,
		ReviewCount:  cmd.ReviewCount,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		IndexedAt:    time.Now(),
	}
}
