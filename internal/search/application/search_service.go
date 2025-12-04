package search

import (
	"context"
	"fmt"

	"github.com/qhato/ecommerce/pkg/elasticsearch"
)

// SearchService handles product search operations
type SearchService struct {
	esClient *elasticsearch.Client
}

// NewSearchService creates a new search service
func NewSearchService(esClient *elasticsearch.Client) *SearchService {
	return &SearchService{
		esClient: esClient,
	}
}

// SearchRequest represents a search request
type SearchRequest struct {
	Query      string
	Categories []int64
	MinPrice   *float64
	MaxPrice   *float64
	Active     *bool
	Facets     []string
	Sort       SortOption
	Page       int
	PageSize   int
}

// SortOption represents sorting options
type SortOption string

const (
	SortRelevance    SortOption = "relevance"
	SortPriceAsc     SortOption = "price_asc"
	SortPriceDesc    SortOption = "price_desc"
	SortNameAsc      SortOption = "name_asc"
	SortNameDesc     SortOption = "name_desc"
	SortNewest       SortOption = "newest"
)

// SearchResponse represents search results
type SearchResponse struct {
	Hits       []ProductHit       `json:"hits"`
	Total      int64              `json:"total"`
	Page       int                `json:"page"`
	PageSize   int                `json:"page_size"`
	TotalPages int                `json:"total_pages"`
	Facets     map[string][]Facet `json:"facets,omitempty"`
	Took       int64              `json:"took_ms"`
}

// ProductHit represents a search result
type ProductHit struct {
	ID          int64                  `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	SKU         string                 `json:"sku"`
	Price       float64                `json:"price"`
	SalePrice   float64                `json:"sale_price,omitempty"`
	Active      bool                   `json:"active"`
	Categories  []int64                `json:"category_ids"`
	Attributes  map[string]interface{} `json:"attributes,omitempty"`
	Score       float64                `json:"score"`
}

// Facet represents a facet value
type Facet struct {
	Key   string `json:"key"`
	Count int64  `json:"count"`
}

// Search performs a product search
func (s *SearchService) Search(ctx context.Context, req SearchRequest) (*SearchResponse, error) {
	// Build Elasticsearch query
	esQuery := s.buildSearchQuery(req)

	// Execute search
	res, err := s.esClient.Search(ctx, ProductIndexName, esQuery)
	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}

	// Parse results
	return s.parseSearchResponse(res, req)
}

// buildSearchQuery constructs the Elasticsearch query
func (s *SearchService) buildSearchQuery(req SearchRequest) map[string]interface{} {
	query := map[string]interface{}{
		"from": (req.Page - 1) * req.PageSize,
		"size": req.PageSize,
	}

	// Build bool query
	boolQuery := map[string]interface{}{
		"must":   []interface{}{},
		"filter": []interface{}{},
	}

	// Text search
	if req.Query != "" {
		boolQuery["must"] = append(boolQuery["must"].([]interface{}), map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  req.Query,
				"fields": []string{"name^3", "description", "sku^2", "category_names"},
				"type":   "best_fields",
				"fuzziness": "AUTO",
			},
		})
	} else {
		// Match all if no query
		boolQuery["must"] = append(boolQuery["must"].([]interface{}), map[string]interface{}{
			"match_all": map[string]interface{}{},
		})
	}

	// Category filter
	if len(req.Categories) > 0 {
		boolQuery["filter"] = append(boolQuery["filter"].([]interface{}), map[string]interface{}{
			"terms": map[string]interface{}{
				"category_ids": req.Categories,
			},
		})
	}

	// Price range filter
	if req.MinPrice != nil || req.MaxPrice != nil {
		priceFilter := map[string]interface{}{}
		if req.MinPrice != nil {
			priceFilter["gte"] = *req.MinPrice
		}
		if req.MaxPrice != nil {
			priceFilter["lte"] = *req.MaxPrice
		}
		boolQuery["filter"] = append(boolQuery["filter"].([]interface{}), map[string]interface{}{
			"range": map[string]interface{}{
				"price": priceFilter,
			},
		})
	}

	// Active filter
	if req.Active != nil {
		boolQuery["filter"] = append(boolQuery["filter"].([]interface{}), map[string]interface{}{
			"term": map[string]interface{}{
				"active": *req.Active,
			},
		})
		boolQuery["filter"] = append(boolQuery["filter"].([]interface{}), map[string]interface{}{
			"term": map[string]interface{}{
				"archived": false,
			},
		})
	}

	query["query"] = map[string]interface{}{
		"bool": boolQuery,
	}

	// Add sorting
	query["sort"] = s.buildSort(req.Sort)

	// Add aggregations (facets)
	if len(req.Facets) > 0 {
		query["aggs"] = s.buildAggregations(req.Facets)
	}

	return query
}

// buildSort constructs the sort clause
func (s *SearchService) buildSort(sort SortOption) []interface{} {
	switch sort {
	case SortPriceAsc:
		return []interface{}{
			map[string]interface{}{"price": "asc"},
		}
	case SortPriceDesc:
		return []interface{}{
			map[string]interface{}{"price": "desc"},
		}
	case SortNameAsc:
		return []interface{}{
			map[string]interface{}{"name.keyword": "asc"},
		}
	case SortNameDesc:
		return []interface{}{
			map[string]interface{}{"name.keyword": "desc"},
		}
	case SortNewest:
		return []interface{}{
			map[string]interface{}{"created_at": "desc"},
		}
	default: // SortRelevance
		return []interface{}{
			map[string]interface{}{"_score": "desc"},
		}
	}
}

// buildAggregations constructs aggregations for facets
func (s *SearchService) buildAggregations(facets []string) map[string]interface{} {
	aggs := make(map[string]interface{})

	for _, facet := range facets {
		switch facet {
		case "categories":
			aggs["categories"] = map[string]interface{}{
				"terms": map[string]interface{}{
					"field": "category_ids",
					"size":  50,
				},
			}
		case "price_ranges":
			aggs["price_ranges"] = map[string]interface{}{
				"range": map[string]interface{}{
					"field": "price",
					"ranges": []map[string]interface{}{
						{"key": "0-50", "to": 50},
						{"key": "50-100", "from": 50, "to": 100},
						{"key": "100-200", "from": 100, "to": 200},
						{"key": "200+", "from": 200},
					},
				},
			}
		}
	}

	return aggs
}

// parseSearchResponse converts Elasticsearch response to SearchResponse
func (s *SearchService) parseSearchResponse(esRes *elasticsearch.SearchResponse, req SearchRequest) (*SearchResponse, error) {
	hits := make([]ProductHit, len(esRes.Hits.Hits))
	for i, hit := range esRes.Hits.Hits {
		hits[i] = ProductHit{
			ID:          int64(hit.Source["id"].(float64)),
			Name:        hit.Source["name"].(string),
			Description: getStringOrEmpty(hit.Source, "description"),
			SKU:         getStringOrEmpty(hit.Source, "sku"),
			Price:       getFloat64OrZero(hit.Source, "price"),
			SalePrice:   getFloat64OrZero(hit.Source, "sale_price"),
			Active:      getBoolOrFalse(hit.Source, "active"),
			Score:       hit.Score,
		}

		// Parse categories
		if catIds, ok := hit.Source["category_ids"].([]interface{}); ok {
			hits[i].Categories = make([]int64, len(catIds))
			for j, catID := range catIds {
				hits[i].Categories[j] = int64(catID.(float64))
			}
		}

		// Parse attributes
		if attrs, ok := hit.Source["attributes"].(map[string]interface{}); ok {
			hits[i].Attributes = attrs
		}
	}

	totalPages := int((esRes.Hits.Total.Value + int64(req.PageSize) - 1) / int64(req.PageSize))

	response := &SearchResponse{
		Hits:       hits,
		Total:      esRes.Hits.Total.Value,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
		Took:       esRes.Took,
	}

	// Parse facets
	if len(esRes.Aggregations) > 0 {
		response.Facets = s.parseFacets(esRes.Aggregations)
	}

	return response, nil
}

// parseFacets converts Elasticsearch aggregations to facets
func (s *SearchService) parseFacets(aggs map[string]elasticsearch.Aggregation) map[string][]Facet {
	facets := make(map[string][]Facet)

	for name, agg := range aggs {
		if len(agg.Buckets) > 0 {
			facetValues := make([]Facet, len(agg.Buckets))
			for i, bucket := range agg.Buckets {
				facetValues[i] = Facet{
					Key:   fmt.Sprintf("%v", bucket.Key),
					Count: bucket.DocCount,
				}
			}
			facets[name] = facetValues
		}
	}

	return facets
}

// Autocomplete performs autocomplete search
func (s *SearchService) Autocomplete(ctx context.Context, prefix string, limit int) ([]string, error) {
	query := map[string]interface{}{
		"size": limit,
		"_source": []string{"name"},
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []interface{}{
					map[string]interface{}{
						"match": map[string]interface{}{
							"name.autocomplete": prefix,
						},
					},
				},
				"filter": []interface{}{
					map[string]interface{}{
						"term": map[string]interface{}{
							"active": true,
						},
					},
					map[string]interface{}{
						"term": map[string]interface{}{
							"archived": false,
						},
					},
				},
			},
		},
	}

	res, err := s.esClient.Search(ctx, ProductIndexName, query)
	if err != nil {
		return nil, fmt.Errorf("autocomplete search failed: %w", err)
	}

	suggestions := make([]string, 0, len(res.Hits.Hits))
	seen := make(map[string]bool)

	for _, hit := range res.Hits.Hits {
		if name, ok := hit.Source["name"].(string); ok {
			if !seen[name] {
				suggestions = append(suggestions, name)
				seen[name] = true
			}
		}
	}

	return suggestions, nil
}

// Helper functions
func getStringOrEmpty(m map[string]interface{}, key string) string {
	if val, ok := m[key].(string); ok {
		return val
	}
	return ""
}

func getFloat64OrZero(m map[string]interface{}, key string) float64 {
	if val, ok := m[key].(float64); ok {
		return val
	}
	return 0.0
}

func getBoolOrFalse(m map[string]interface{}, key string) bool {
	if val, ok := m[key].(bool); ok {
		return val
	}
	return false
}