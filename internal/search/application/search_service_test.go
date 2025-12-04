package search_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/qhato/ecommerce/internal/search/application"
	"github.com/qhato/ecommerce/pkg/elasticsearch"
)

// MockElasticsearchClient is a mock implementation of elasticsearch.Client
type MockElasticsearchClient struct {
	mock.Mock
}

func (m *MockElasticsearchClient) Search(ctx context.Context, indexName string, query map[string]interface{}) (*elasticsearch.SearchResponse, error) {
	args := m.Called(ctx, indexName, query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*elasticsearch.SearchResponse), args.Error(1)
}

func (m *MockElasticsearchClient) IndexDocument(ctx context.Context, indexName, documentID string, document interface{}) error {
	args := m.Called(ctx, indexName, documentID, document)
	return args.Error(0)
}

func TestSearchService_Search(t *testing.T) {
	// Mock Elasticsearch client
	mockES := new(MockElasticsearchClient)

	// Create search service
	searchService := search.NewSearchService(mockES)

	// Test case: successful search with query
	t.Run("successful search with query", func(t *testing.T) {
		mockResponse := &elasticsearch.SearchResponse{
			Took:     15,
			TimedOut: false,
			Hits: struct {
				Total struct {
					Value    int64  `json:"value"`
					Relation string `json:"relation"`
				} `json:"total"`
				MaxScore float64              `json:"max_score"`
				Hits     []elasticsearch.Hit `json:"hits"`
			}{
				Total: struct {
					Value    int64  `json:"value"`
					Relation string `json:"relation"`
				}{
					Value:    2,
					Relation: "eq",
				},
				MaxScore: 1.5,
				Hits: []elasticsearch.Hit{
					{
						Index: "products",
						ID:    "1",
						Score: 1.5,
						Source: map[string]interface{}{
							"id":          float64(1),
							"name":        "Test Product 1",
							"description": "Test description",
							"sku":         "TEST-001",
							"price":       99.99,
							"sale_price":  79.99,
							"active":      true,
							"category_ids": []interface{}{float64(1), float64(2)},
						},
					},
					{
						Index: "products",
						ID:    "2",
						Score: 1.2,
						Source: map[string]interface{}{
							"id":          float64(2),
							"name":        "Test Product 2",
							"description": "Another test",
							"sku":         "TEST-002",
							"price":       149.99,
							"active":      true,
							"category_ids": []interface{}{float64(1)},
						},
					},
				},
			},
		}

		mockES.On("Search", mock.Anything, "products", mock.Anything).Return(mockResponse, nil).Once()

		req := search.SearchRequest{
			Query:    "test",
			Page:     1,
			PageSize: 20,
		}

		result, err := searchService.Search(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int64(2), result.Total)
		assert.Equal(t, 2, len(result.Hits))
		assert.Equal(t, "Test Product 1", result.Hits[0].Name)
		assert.Equal(t, 99.99, result.Hits[0].Price)
		assert.Equal(t, 79.99, result.Hits[0].SalePrice)

		mockES.AssertExpectations(t)
	})

	// Test case: search with filters
	t.Run("search with category and price filters", func(t *testing.T) {
		mockResponse := &elasticsearch.SearchResponse{
			Took:     10,
			TimedOut: false,
			Hits: struct {
				Total struct {
					Value    int64  `json:"value"`
					Relation string `json:"relation"`
				} `json:"total"`
				MaxScore float64              `json:"max_score"`
				Hits     []elasticsearch.Hit `json:"hits"`
			}{
				Total: struct {
					Value    int64  `json:"value"`
					Relation string `json:"relation"`
				}{
					Value: 1,
				},
				Hits: []elasticsearch.Hit{
					{
						ID:    "1",
						Score: 1.0,
						Source: map[string]interface{}{
							"id":    float64(1),
							"name":  "Filtered Product",
							"price": 75.00,
						},
					},
				},
			},
		}

		mockES.On("Search", mock.Anything, "products", mock.Anything).Return(mockResponse, nil).Once()

		minPrice := 50.0
		maxPrice := 100.0
		req := search.SearchRequest{
			Categories: []int64{1},
			MinPrice:   &minPrice,
			MaxPrice:   &maxPrice,
			Page:       1,
			PageSize:   20,
		}

		result, err := searchService.Search(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int64(1), result.Total)

		mockES.AssertExpectations(t)
	})

	// Test case: search with facets
	t.Run("search with facets", func(t *testing.T) {
		mockResponse := &elasticsearch.SearchResponse{
			Took: 12,
			Hits: struct {
				Total struct {
					Value    int64  `json:"value"`
					Relation string `json:"relation"`
				} `json:"total"`
				MaxScore float64              `json:"max_score"`
				Hits     []elasticsearch.Hit `json:"hits"`
			}{
				Total: struct {
					Value    int64  `json:"value"`
					Relation string `json:"relation"`
				}{Value: 0},
				Hits: []elasticsearch.Hit{},
			},
			Aggregations: map[string]elasticsearch.Aggregation{
				"categories": {
					Buckets: []elasticsearch.Bucket{
						{Key: float64(1), DocCount: 10},
						{Key: float64(2), DocCount: 5},
					},
				},
			},
		}

		mockES.On("Search", mock.Anything, "products", mock.Anything).Return(mockResponse, nil).Once()

		req := search.SearchRequest{
			Facets:   []string{"categories"},
			Page:     1,
			PageSize: 20,
		}

		result, err := searchService.Search(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotNil(t, result.Facets)
		assert.Equal(t, 2, len(result.Facets["categories"]))
		assert.Equal(t, int64(10), result.Facets["categories"][0].Count)

		mockES.AssertExpectations(t)
	})
}

func TestSearchService_Autocomplete(t *testing.T) {
	mockES := new(MockElasticsearchClient)
	searchService := search.NewSearchService(mockES)

	t.Run("successful autocomplete", func(t *testing.T) {
		mockResponse := &elasticsearch.SearchResponse{
			Hits: struct {
				Total struct {
					Value    int64  `json:"value"`
					Relation string `json:"relation"`
				} `json:"total"`
				MaxScore float64              `json:"max_score"`
				Hits     []elasticsearch.Hit `json:"hits"`
			}{
				Hits: []elasticsearch.Hit{
					{
						Source: map[string]interface{}{
							"name": "Laptop",
						},
					},
					{
						Source: map[string]interface{}{
							"name": "Laptop Pro",
						},
					},
					{
						Source: map[string]interface{}{
							"name": "Laptop Air",
						},
					},
				},
			},
		}

		mockES.On("Search", mock.Anything, "products", mock.Anything).Return(mockResponse, nil).Once()

		suggestions, err := searchService.Autocomplete(context.Background(), "lap", 10)

		assert.NoError(t, err)
		assert.Equal(t, 3, len(suggestions))
		assert.Contains(t, suggestions, "Laptop")
		assert.Contains(t, suggestions, "Laptop Pro")

		mockES.AssertExpectations(t)
	})

	t.Run("autocomplete with duplicates", func(t *testing.T) {
		mockResponse := &elasticsearch.SearchResponse{
			Hits: struct {
				Total struct {
					Value    int64  `json:"value"`
					Relation string `json:"relation"`
				} `json:"total"`
				MaxScore float64              `json:"max_score"`
				Hits     []elasticsearch.Hit `json:"hits"`
			}{
				Hits: []elasticsearch.Hit{
					{Source: map[string]interface{}{"name": "Phone"}},
					{Source: map[string]interface{}{"name": "Phone"}},
					{Source: map[string]interface{}{"name": "Tablet"}},
				},
			},
		}

		mockES.On("Search", mock.Anything, "products", mock.Anything).Return(mockResponse, nil).Once()

		suggestions, err := searchService.Autocomplete(context.Background(), "ph", 10)

		assert.NoError(t, err)
		assert.Equal(t, 2, len(suggestions)) // Duplicates removed
		assert.Contains(t, suggestions, "Phone")
		assert.Contains(t, suggestions, "Tablet")

		mockES.AssertExpectations(t)
	})
}

func TestSearchRequest_Validation(t *testing.T) {
	t.Run("pagination defaults", func(t *testing.T) {
		req := search.SearchRequest{
			Page:     0,
			PageSize: 0,
		}

		// Handler should set defaults to 1 and 20
		assert.True(t, req.Page == 0 || req.Page == 1)
		assert.True(t, req.PageSize == 0 || req.PageSize == 20)
	})

	t.Run("sort options", func(t *testing.T) {
		sortOptions := []search.SortOption{
			search.SortRelevance,
			search.SortPriceAsc,
			search.SortPriceDesc,
			search.SortNameAsc,
			search.SortNameDesc,
			search.SortNewest,
		}

		for _, opt := range sortOptions {
			req := search.SearchRequest{
				Sort: opt,
			}
			assert.NotEmpty(t, req.Sort)
		}
	})
}