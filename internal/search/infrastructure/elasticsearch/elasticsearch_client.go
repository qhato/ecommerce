package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"

	"github.com/qhato/ecommerce/internal/search/domain"
	"github.com/qhato/ecommerce/pkg/logger"
)

const (
	// ProductIndexName is the name of the product index
	ProductIndexName = "products"
)

// ElasticsearchClient implements search operations using Elasticsearch
type ElasticsearchClient struct {
	client *elasticsearch.Client
	logger logger.Logger
}

// NewElasticsearchClient creates a new Elasticsearch client
func NewElasticsearchClient(esClient *elasticsearch.Client, logger logger.Logger) *ElasticsearchClient {
	return &ElasticsearchClient{
		client: esClient,
		logger: logger,
	}
}

// IndexProduct indexes a product document
func (c *ElasticsearchClient) IndexProduct(ctx context.Context, doc *domain.ProductSearchDocument) error {
	// Convert to JSON
	data, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("failed to marshal document: %w", err)
	}

	// Index document
	req := esapi.IndexRequest{
		Index:      ProductIndexName,
		DocumentID: fmt.Sprintf("%d", doc.ProductID),
		Body:       bytes.NewReader(data),
		Refresh:    "false",
	}

	res, err := req.Do(ctx, c.client)
	if err != nil {
		return fmt.Errorf("failed to index document: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("elasticsearch error: %s", res.String())
	}

	c.logger.Info("Product indexed in Elasticsearch",
		logger.Field{Key: "product_id", Value: doc.ProductID},
		logger.Field{Key: "sku", Value: doc.SKU},
	)

	return nil
}

// BulkIndexProducts indexes multiple products in bulk
func (c *ElasticsearchClient) BulkIndexProducts(ctx context.Context, docs []*domain.ProductSearchDocument) error {
	if len(docs) == 0 {
		return nil
	}

	var buf bytes.Buffer

	for _, doc := range docs {
		// Action line
		action := map[string]interface{}{
			"index": map[string]interface{}{
				"_index": ProductIndexName,
				"_id":    fmt.Sprintf("%d", doc.ProductID),
			},
		}
		actionData, _ := json.Marshal(action)
		buf.Write(actionData)
		buf.WriteByte('\n')

		// Document line
		docData, _ := json.Marshal(doc)
		buf.Write(docData)
		buf.WriteByte('\n')
	}

	// Perform bulk request
	req := esapi.BulkRequest{
		Body:    &buf,
		Refresh: "false",
	}

	res, err := req.Do(ctx, c.client)
	if err != nil {
		return fmt.Errorf("failed to perform bulk index: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("elasticsearch bulk error: %s", res.String())
	}

	c.logger.Info("Bulk indexed products in Elasticsearch",
		logger.Field{Key: "count", Value: len(docs)},
	)

	return nil
}

// DeleteProduct deletes a product from the index
func (c *ElasticsearchClient) DeleteProduct(ctx context.Context, productID int64) error {
	req := esapi.DeleteRequest{
		Index:      ProductIndexName,
		DocumentID: fmt.Sprintf("%d", productID),
	}

	res, err := req.Do(ctx, c.client)
	if err != nil {
		return fmt.Errorf("failed to delete document: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("elasticsearch delete error: %s", res.String())
	}

	c.logger.Info("Product deleted from Elasticsearch",
		logger.Field{Key: "product_id", Value: productID},
	)

	return nil
}

// ReindexAll is a placeholder for full reindex (handled by command handler)
func (c *ElasticsearchClient) ReindexAll(ctx context.Context) error {
	// This is typically handled by the command handler with batching
	// Here we could optionally delete and recreate the index
	return nil
}

// Search performs a search query
func (c *ElasticsearchClient) Search(ctx context.Context, query *domain.SearchQuery) (*domain.SearchResult, error) {
	// Build Elasticsearch query
	esQuery := c.buildElasticsearchQuery(query)

	// Convert to JSON
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(esQuery); err != nil {
		return nil, fmt.Errorf("failed to encode query: %w", err)
	}

	// Perform search
	res, err := c.client.Search(
		c.client.Search.WithContext(ctx),
		c.client.Search.WithIndex(ProductIndexName),
		c.client.Search.WithBody(&buf),
		c.client.Search.WithFrom((query.Page-1)*query.PageSize),
		c.client.Search.WithSize(query.PageSize),
	)
	if err != nil {
		return nil, fmt.Errorf("elasticsearch search failed: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("elasticsearch search error: %s", res.String())
	}

	// Parse response
	var esResult ElasticsearchResponse
	if err := json.NewDecoder(res.Body).Decode(&esResult); err != nil {
		return nil, fmt.Errorf("failed to parse search response: %w", err)
	}

	// Convert to domain SearchResult
	return c.convertSearchResponse(&esResult, query), nil
}

// Suggest provides autocomplete suggestions
func (c *ElasticsearchClient) Suggest(ctx context.Context, prefix string, limit int) ([]string, error) {
	// Build suggest query
	query := map[string]interface{}{
		"suggest": map[string]interface{}{
			"product-suggest": map[string]interface{}{
				"prefix": prefix,
				"completion": map[string]interface{}{
					"field": "name.completion",
					"size":  limit,
				},
			},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, fmt.Errorf("failed to encode suggest query: %w", err)
	}

	res, err := c.client.Search(
		c.client.Search.WithContext(ctx),
		c.client.Search.WithIndex(ProductIndexName),
		c.client.Search.WithBody(&buf),
	)
	if err != nil {
		return nil, fmt.Errorf("elasticsearch suggest failed: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("elasticsearch suggest error: %s", res.String())
	}

	// Parse suggestions (simplified)
	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse suggest response: %w", err)
	}

	// Extract suggestions
	suggestions := make([]string, 0, limit)
	if suggest, ok := result["suggest"].(map[string]interface{}); ok {
		if productSuggest, ok := suggest["product-suggest"].([]interface{}); ok && len(productSuggest) > 0 {
			if first, ok := productSuggest[0].(map[string]interface{}); ok {
				if options, ok := first["options"].([]interface{}); ok {
					for _, opt := range options {
						if optMap, ok := opt.(map[string]interface{}); ok {
							if text, ok := optMap["text"].(string); ok {
								suggestions = append(suggestions, text)
							}
						}
					}
				}
			}
		}
	}

	return suggestions, nil
}

// buildElasticsearchQuery builds an Elasticsearch query from domain query
func (c *ElasticsearchClient) buildElasticsearchQuery(query *domain.SearchQuery) map[string]interface{} {
	must := make([]interface{}, 0)
	filter := make([]interface{}, 0)

	// Text search
	if query.Query != "" {
		must = append(must, map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  query.Query,
				"fields": []string{"name^3", "description^2", "long_description", "tags"},
				"type":   "best_fields",
			},
		})
	}

	// Only active products
	filter = append(filter, map[string]interface{}{
		"term": map[string]interface{}{
			"is_active": true,
		},
	})

	// Category filter
	if len(query.CategoryIDs) > 0 {
		filter = append(filter, map[string]interface{}{
			"terms": map[string]interface{}{
				"category_id": query.CategoryIDs,
			},
		})
	}

	// Availability filter
	if query.Availability != nil {
		filter = append(filter, map[string]interface{}{
			"term": map[string]interface{}{
				"is_available": *query.Availability,
			},
		})
	}

	// Price range filter
	if query.PriceMin != nil || query.PriceMax != nil {
		priceFilter := make(map[string]interface{})
		if query.PriceMin != nil {
			priceFilter["gte"] = query.PriceMin
		}
		if query.PriceMax != nil {
			priceFilter["lte"] = query.PriceMax
		}
		filter = append(filter, map[string]interface{}{
			"range": map[string]interface{}{
				"price": priceFilter,
			},
		})
	}

	// Facet filters
	for facetName, facetValues := range query.Filters {
		if len(facetValues) > 0 {
			filter = append(filter, map[string]interface{}{
				"terms": map[string]interface{}{
					facetName: facetValues,
				},
			})
		}
	}

	// Build bool query
	boolQuery := map[string]interface{}{}
	if len(must) > 0 {
		boolQuery["must"] = must
	}
	if len(filter) > 0 {
		boolQuery["filter"] = filter
	}

	// If no text query, match all
	if query.Query == "" {
		must = append(must, map[string]interface{}{
			"match_all": map[string]interface{}{},
		})
		boolQuery["must"] = must
	}

	// Build aggregations for facets
	aggs := map[string]interface{}{
		"brands": map[string]interface{}{
			"terms": map[string]interface{}{
				"field": "brand.keyword",
				"size":  20,
			},
		},
		"colors": map[string]interface{}{
			"terms": map[string]interface{}{
				"field": "color.keyword",
				"size":  20,
			},
		},
		"sizes": map[string]interface{}{
			"terms": map[string]interface{}{
				"field": "size.keyword",
				"size":  20,
			},
		},
		"price_ranges": map[string]interface{}{
			"range": map[string]interface{}{
				"field": "price",
				"ranges": []map[string]interface{}{
					{"to": 50, "key": "0-50"},
					{"from": 50, "to": 100, "key": "50-100"},
					{"from": 100, "to": 200, "key": "100-200"},
					{"from": 200, "to": 500, "key": "200-500"},
					{"from": 500, "key": "500+"},
				},
			},
		},
	}

	// Build sort
	var sort []interface{}
	switch query.SortBy {
	case "price_asc":
		sort = []interface{}{map[string]interface{}{"price": "asc"}}
	case "price_desc":
		sort = []interface{}{map[string]interface{}{"price": "desc"}}
	case "name":
		sort = []interface{}{map[string]interface{}{"name.keyword": "asc"}}
	case "created":
		sort = []interface{}{map[string]interface{}{"created_at": "desc"}}
	default: // relevance
		sort = []interface{}{map[string]interface{}{"_score": "desc"}}
	}

	return map[string]interface{}{
		"query": map[string]interface{}{
			"bool": boolQuery,
		},
		"aggs": aggs,
		"sort": sort,
	}
}

// convertSearchResponse converts Elasticsearch response to domain SearchResult
func (c *ElasticsearchClient) convertSearchResponse(esResult *ElasticsearchResponse, query *domain.SearchQuery) *domain.SearchResult {
	// Convert hits to documents
	documents := make([]*domain.SearchDocument, len(esResult.Hits.Hits))
	for i, hit := range esResult.Hits.Hits {
		var doc domain.ProductSearchDocument
		json.Unmarshal(hit.Source, &doc)

		documents[i] = &domain.SearchDocument{
			ID:          hit.ID,
			Type:        "product",
			Title:       doc.Name,
			Description: doc.Description,
			Score:       hit.Score,
			Fields: map[string]interface{}{
				"product_id":   doc.ProductID,
				"sku":          doc.SKU,
				"price":        doc.Price,
				"image_url":    doc.ImageURL,
				"is_available": doc.IsAvailable,
				"rating":       doc.Rating,
				"review_count": doc.ReviewCount,
			},
		}
	}

	// Convert aggregations to facets
	facets := make(map[string]*domain.Facet)

	if brands, ok := esResult.Aggregations["brands"].(map[string]interface{}); ok {
		facets["brand"] = c.convertTermsAggregation(brands, "Brand")
	}
	if colors, ok := esResult.Aggregations["colors"].(map[string]interface{}); ok {
		facets["color"] = c.convertTermsAggregation(colors, "Color")
	}
	if sizes, ok := esResult.Aggregations["sizes"].(map[string]interface{}); ok {
		facets["size"] = c.convertTermsAggregation(sizes, "Size")
	}
	if priceRanges, ok := esResult.Aggregations["price_ranges"].(map[string]interface{}); ok {
		facets["price_range"] = c.convertRangeAggregation(priceRanges, "Price Range")
	}

	total := esResult.Hits.Total.Value
	totalPages := (total + query.PageSize - 1) / query.PageSize

	return &domain.SearchResult{
		Documents:  documents,
		Total:      total,
		Facets:     facets,
		Page:       query.Page,
		PageSize:   query.PageSize,
		TotalPages: totalPages,
		Query:      query.Query,
	}
}

// convertTermsAggregation converts terms aggregation to Facet
func (c *ElasticsearchClient) convertTermsAggregation(agg map[string]interface{}, name string) *domain.Facet {
	values := make([]*domain.FacetValue, 0)

	if buckets, ok := agg["buckets"].([]interface{}); ok {
		for _, bucket := range buckets {
			if b, ok := bucket.(map[string]interface{}); ok {
				key := b["key"].(string)
				count := int(b["doc_count"].(float64))
				values = append(values, &domain.FacetValue{
					Value: key,
					Count: count,
				})
			}
		}
	}

	return &domain.Facet{
		Name:   name,
		Values: values,
	}
}

// convertRangeAggregation converts range aggregation to Facet
func (c *ElasticsearchClient) convertRangeAggregation(agg map[string]interface{}, name string) *domain.Facet {
	values := make([]*domain.FacetValue, 0)

	if buckets, ok := agg["buckets"].([]interface{}); ok {
		for _, bucket := range buckets {
			if b, ok := bucket.(map[string]interface{}); ok {
				key := b["key"].(string)
				count := int(b["doc_count"].(float64))
				values = append(values, &domain.FacetValue{
					Value: key,
					Count: count,
				})
			}
		}
	}

	return &domain.Facet{
		Name:   name,
		Values: values,
	}
}

// ElasticsearchResponse represents Elasticsearch search response
type ElasticsearchResponse struct {
	Hits struct {
		Total struct {
			Value int `json:"value"`
		} `json:"total"`
		Hits []struct {
			ID     string          `json:"_id"`
			Score  float64         `json:"_score"`
			Source json.RawMessage `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
	Aggregations map[string]interface{} `json:"aggregations"`
}

// CreateProductIndex creates the product index with proper mappings
func (c *ElasticsearchClient) CreateProductIndex(ctx context.Context) error {
	mapping := `{
		"mappings": {
			"properties": {
				"product_id": {"type": "long"},
				"sku": {"type": "keyword"},
				"name": {
					"type": "text",
					"fields": {
						"keyword": {"type": "keyword"},
						"completion": {"type": "completion"}
					}
				},
				"description": {"type": "text"},
				"long_description": {"type": "text"},
				"tags": {"type": "keyword"},
				"price": {"type": "float"},
				"sale_price": {"type": "float"},
				"on_sale": {"type": "boolean"},
				"category_id": {"type": "long"},
				"category_name": {
					"type": "text",
					"fields": {"keyword": {"type": "keyword"}}
				},
				"category_path": {"type": "keyword"},
				"brand": {
					"type": "text",
					"fields": {"keyword": {"type": "keyword"}}
				},
				"color": {"type": "keyword"},
				"size": {"type": "keyword"},
				"attributes": {"type": "object"},
				"is_available": {"type": "boolean"},
				"stock_level": {"type": "integer"},
				"image_url": {"type": "keyword"},
				"thumbnail_url": {"type": "keyword"},
				"is_active": {"type": "boolean"},
				"is_featured": {"type": "boolean"},
				"rating": {"type": "float"},
				"review_count": {"type": "integer"},
				"view_count": {"type": "integer"},
				"created_at": {"type": "date"},
				"updated_at": {"type": "date"},
				"indexed_at": {"type": "date"}
			}
		}
	}`

	req := esapi.IndicesCreateRequest{
		Index: ProductIndexName,
		Body:  strings.NewReader(mapping),
	}

	res, err := req.Do(ctx, c.client)
	if err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() && !strings.Contains(res.String(), "resource_already_exists_exception") {
		return fmt.Errorf("elasticsearch create index error: %s", res.String())
	}

	c.logger.Info("Product index created in Elasticsearch")
	return nil
}
