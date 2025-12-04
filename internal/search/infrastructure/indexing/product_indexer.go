package search

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/shopspring/decimal"
	"github.com/qhato/ecommerce/internal/catalog/domain"
	"github.com/qhato/ecommerce/pkg/elasticsearch"
)

const (
	ProductIndexName = "products"
)

// ProductDocument represents a product in Elasticsearch
type ProductDocument struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	SKU         string    `json:"sku"`
	Price       float64   `json:"price"`
	SalePrice   float64   `json:"sale_price,omitempty"`
	Active      bool      `json:"active"`
	Archived    bool      `json:"archived"`
	Categories  []int64   `json:"category_ids"`
	CategoryNames []string `json:"category_names"`
	Attributes  map[string]interface{} `json:"attributes"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ProductIndexer handles product indexing operations
type ProductIndexer struct {
	esClient *elasticsearch.Client
}

// NewProductIndexer creates a new product indexer
func NewProductIndexer(esClient *elasticsearch.Client) *ProductIndexer {
	return &ProductIndexer{
		esClient: esClient,
	}
}

// CreateProductIndex creates the product index with mapping
func (pi *ProductIndexer) CreateProductIndex(ctx context.Context) error {
	// Check if index already exists
	exists, err := pi.esClient.IndexExists(ctx, ProductIndexName)
	if err != nil {
		return fmt.Errorf("error checking index existence: %w", err)
	}

	if exists {
		return nil // Index already exists
	}

	mapping := map[string]interface{}{
		"settings": map[string]interface{}{
			"number_of_shards":   1,
			"number_of_replicas": 1,
			"analysis": map[string]interface{}{
				"analyzer": map[string]interface{}{
					"product_analyzer": map[string]interface{}{
						"type":      "custom",
						"tokenizer": "standard",
						"filter": []string{
							"lowercase",
							"asciifolding",
							"stop",
						},
					},
					"autocomplete_analyzer": map[string]interface{}{
						"type":      "custom",
						"tokenizer": "standard",
						"filter": []string{
							"lowercase",
							"asciifolding",
							"autocomplete_filter",
						},
					},
				},
				"filter": map[string]interface{}{
					"autocomplete_filter": map[string]interface{}{
						"type":     "edge_ngram",
						"min_gram": 2,
						"max_gram": 20,
					},
				},
			},
		},
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"id": map[string]interface{}{
					"type": "long",
				},
				"name": map[string]interface{}{
					"type":     "text",
					"analyzer": "product_analyzer",
					"fields": map[string]interface{}{
						"keyword": map[string]interface{}{
							"type": "keyword",
						},
						"autocomplete": map[string]interface{}{
							"type":     "text",
							"analyzer": "autocomplete_analyzer",
						},
					},
				},
				"description": map[string]interface{}{
					"type":     "text",
					"analyzer": "product_analyzer",
				},
				"sku": map[string]interface{}{
					"type": "keyword",
				},
				"price": map[string]interface{}{
					"type": "double",
				},
				"sale_price": map[string]interface{}{
					"type": "double",
				},
				"active": map[string]interface{}{
					"type": "boolean",
				},
				"archived": map[string]interface{}{
					"type": "boolean",
				},
				"category_ids": map[string]interface{}{
					"type": "long",
				},
				"category_names": map[string]interface{}{
					"type":     "text",
					"analyzer": "product_analyzer",
					"fields": map[string]interface{}{
						"keyword": map[string]interface{}{
							"type": "keyword",
						},
					},
				},
				"attributes": map[string]interface{}{
					"type": "object",
					"dynamic": true,
				},
				"created_at": map[string]interface{}{
					"type": "date",
				},
				"updated_at": map[string]interface{}{
					"type": "date",
				},
			},
		},
	}

	return pi.esClient.CreateIndex(ctx, ProductIndexName, mapping)
}

// IndexProduct indexes a single product
func (pi *ProductIndexer) IndexProduct(ctx context.Context, product *domain.Product) error {
	doc := pi.productToDocument(product)
	return pi.esClient.IndexDocument(ctx, ProductIndexName, strconv.FormatInt(product.ID, 10), doc)
}

// UpdateProduct updates a product in the index
func (pi *ProductIndexer) UpdateProduct(ctx context.Context, product *domain.Product) error {
	doc := pi.productToDocument(product)
	return pi.esClient.UpdateDocument(ctx, ProductIndexName, strconv.FormatInt(product.ID, 10), doc)
}

// DeleteProduct deletes a product from the index
func (pi *ProductIndexer) DeleteProduct(ctx context.Context, productID int64) error {
	return pi.esClient.DeleteDocument(ctx, ProductIndexName, strconv.FormatInt(productID, 10))
}

// BulkIndexProducts indexes multiple products in bulk
func (pi *ProductIndexer) BulkIndexProducts(ctx context.Context, products []*domain.Product) error {
	if len(products) == 0 {
		return nil
	}

	bulkDocs := make([]elasticsearch.BulkDocument, len(products))
	for i, product := range products {
		bulkDocs[i] = elasticsearch.BulkDocument{
			ID:       strconv.FormatInt(product.ID, 10),
			Document: pi.productToDocument(product),
		}
	}

	return pi.esClient.BulkIndexDocuments(ctx, ProductIndexName, bulkDocs)
}

// ReindexAllProducts reindexes all products (for full refresh)
func (pi *ProductIndexer) ReindexAllProducts(ctx context.Context, products []*domain.Product) error {
	// Delete existing index
	_ = pi.esClient.DeleteIndex(ctx, ProductIndexName)

	// Create new index
	if err := pi.CreateProductIndex(ctx); err != nil {
		return fmt.Errorf("error creating index: %w", err)
	}

	// Bulk index all products
	return pi.BulkIndexProducts(ctx, products)
}

// productToDocument converts a product entity to an Elasticsearch document
func (pi *ProductIndexer) productToDocument(product *domain.Product) ProductDocument {
	doc := ProductDocument{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		SKU:         product.SKU,
		Active:      product.Active,
		Archived:    product.Archived,
		Attributes:  make(map[string]interface{}),
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}

	// Convert price to float64
	if product.DefaultPrice != nil {
		price, _ := product.DefaultPrice.Float64()
		doc.Price = price
	}

	if product.SalePrice != nil {
		salePrice, _ := product.SalePrice.Float64()
		doc.SalePrice = salePrice
	}

	// Convert attributes
	for key, value := range product.Attributes {
		doc.Attributes[key] = value
	}

	return doc
}

// CategoryIndexer handles category indexing (similar pattern)
type CategoryIndexer struct {
	esClient *elasticsearch.Client
}

func NewCategoryIndexer(esClient *elasticsearch.Client) *CategoryIndexer {
	return &CategoryIndexer{
		esClient: esClient,
	}
}

// CategoryDocument represents a category in Elasticsearch
type CategoryDocument struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Active      bool      `json:"active"`
	ParentID    *int64    `json:"parent_id,omitempty"`
	Level       int       `json:"level"`
	Path        string    `json:"path"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// IndexSyncService handles automatic synchronization
type IndexSyncService struct {
	productIndexer  *ProductIndexer
	categoryIndexer *CategoryIndexer
}

func NewIndexSyncService(esClient *elasticsearch.Client) *IndexSyncService {
	return &IndexSyncService{
		productIndexer:  NewProductIndexer(esClient),
		categoryIndexer: NewCategoryIndexer(esClient),
	}
}

// OnProductCreated handles product creation event
func (s *IndexSyncService) OnProductCreated(ctx context.Context, product *domain.Product) error {
	return s.productIndexer.IndexProduct(ctx, product)
}

// OnProductUpdated handles product update event
func (s *IndexSyncService) OnProductUpdated(ctx context.Context, product *domain.Product) error {
	return s.productIndexer.UpdateProduct(ctx, product)
}

// OnProductDeleted handles product deletion event
func (s *IndexSyncService) OnProductDeleted(ctx context.Context, productID int64) error {
	return s.productIndexer.DeleteProduct(ctx, productID)
}

// OnProductArchived handles product archival (update active status)
func (s *IndexSyncService) OnProductArchived(ctx context.Context, product *domain.Product) error {
	return s.productIndexer.UpdateProduct(ctx, product)
}