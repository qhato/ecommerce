package domain

import (
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

// SearchDocument represents a document in the search index
type SearchDocument struct {
	ID          string
	Type        string // "product", "category", etc.
	Title       string
	Description string
	Content     string
	Fields      map[string]interface{} // Custom fields
	Facets      map[string][]string    // Facet values
	Score       float64
	IndexedAt   time.Time
}

// SearchQuery represents a search query
type SearchQuery struct {
	Query        string
	Filters      map[string][]string // Facet filters
	PriceMin     *decimal.Decimal
	PriceMax     *decimal.Decimal
	CategoryIDs  []string
	Availability *bool
	SortBy       string // "relevance", "price_asc", "price_desc", "name", "created"
	Page         int
	PageSize     int
}

// SearchResult represents search results
type SearchResult struct {
	Documents  []*SearchDocument
	Total      int
	Facets     map[string]*Facet
	Page       int
	PageSize   int
	TotalPages int
	Query      string
}

// Facet represents a facet for filtering
type Facet struct {
	Name   string
	Values []*FacetValue
}

// FacetValue represents a facet value with count
type FacetValue struct {
	Value string
	Count int
}

// SearchIndexer defines the interface for indexing documents
type SearchIndexer interface {
	Index(doc *SearchDocument) error
	Update(id string, doc *SearchDocument) error
	Delete(id string) error
	Bulk(docs []*SearchDocument) error
}

// SearchEngine defines the interface for searching
type SearchEngine interface {
	Search(query *SearchQuery) (*SearchResult, error)
	Suggest(prefix string, limit int) ([]string, error)
}

// InMemorySearchIndex is a simple in-memory search index
type InMemorySearchIndex struct {
	documents map[string]*SearchDocument
	facets    map[string]map[string]int // facet name -> value -> count
}

// NewInMemorySearchIndex creates a new in-memory search index
func NewInMemorySearchIndex() *InMemorySearchIndex {
	return &InMemorySearchIndex{
		documents: make(map[string]*SearchDocument),
		facets:    make(map[string]map[string]int),
	}
}

func (idx *InMemorySearchIndex) Index(doc *SearchDocument) error {
	doc.IndexedAt = time.Now()
	idx.documents[doc.ID] = doc

	// Update facets
	for facetName, facetValues := range doc.Facets {
		if idx.facets[facetName] == nil {
			idx.facets[facetName] = make(map[string]int)
		}
		for _, value := range facetValues {
			idx.facets[facetName][value]++
		}
	}

	return nil
}

func (idx *InMemorySearchIndex) Update(id string, doc *SearchDocument) error {
	// Remove old facets
	if oldDoc, exists := idx.documents[id]; exists {
		for facetName, facetValues := range oldDoc.Facets {
			for _, value := range facetValues {
				idx.facets[facetName][value]--
				if idx.facets[facetName][value] <= 0 {
					delete(idx.facets[facetName], value)
				}
			}
		}
	}

	return idx.Index(doc)
}

func (idx *InMemorySearchIndex) Delete(id string) error {
	if doc, exists := idx.documents[id]; exists {
		// Remove facets
		for facetName, facetValues := range doc.Facets {
			for _, value := range facetValues {
				idx.facets[facetName][value]--
				if idx.facets[facetName][value] <= 0 {
					delete(idx.facets[facetName], value)
				}
			}
		}
		delete(idx.documents, id)
	}
	return nil
}

func (idx *InMemorySearchIndex) Bulk(docs []*SearchDocument) error {
	for _, doc := range docs {
		if err := idx.Index(doc); err != nil {
			return err
		}
	}
	return nil
}

func (idx *InMemorySearchIndex) Search(query *SearchQuery) (*SearchResult, error) {
	matches := make([]*SearchDocument, 0)

	// Simple keyword search
	queryLower := strings.ToLower(query.Query)

	for _, doc := range idx.documents {
		// Check if document matches query
		if query.Query != "" {
			titleMatch := strings.Contains(strings.ToLower(doc.Title), queryLower)
			descMatch := strings.Contains(strings.ToLower(doc.Description), queryLower)
			contentMatch := strings.Contains(strings.ToLower(doc.Content), queryLower)

			if !titleMatch && !descMatch && !contentMatch {
				continue
			}

			// Calculate simple relevance score
			score := 0.0
			if titleMatch {
				score += 3.0
			}
			if descMatch {
				score += 2.0
			}
			if contentMatch {
				score += 1.0
			}
			doc.Score = score
		}

		// Apply facet filters
		if len(query.Filters) > 0 {
			matches := true
			for filterName, filterValues := range query.Filters {
				docValues, exists := doc.Facets[filterName]
				if !exists {
					matches = false
					break
				}

				hasMatch := false
				for _, filterValue := range filterValues {
					for _, docValue := range docValues {
						if docValue == filterValue {
							hasMatch = true
							break
						}
					}
					if hasMatch {
						break
					}
				}

				if !hasMatch {
					matches = false
					break
				}
			}

			if !matches {
				continue
			}
		}

		matches = append(matches, doc)
	}

	// Sort results
	// TODO: Implement sorting

	// Build facets from results
	resultFacets := idx.buildFacets(matches)

	// Pagination
	total := len(matches)
	start := (query.Page - 1) * query.PageSize
	end := start + query.PageSize

	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	paginatedDocs := matches[start:end]
	totalPages := (total + query.PageSize - 1) / query.PageSize

	return &SearchResult{
		Documents:  paginatedDocs,
		Total:      total,
		Facets:     resultFacets,
		Page:       query.Page,
		PageSize:   query.PageSize,
		TotalPages: totalPages,
		Query:      query.Query,
	}, nil
}

func (idx *InMemorySearchIndex) buildFacets(docs []*SearchDocument) map[string]*Facet {
	facetCounts := make(map[string]map[string]int)

	for _, doc := range docs {
		for facetName, facetValues := range doc.Facets {
			if facetCounts[facetName] == nil {
				facetCounts[facetName] = make(map[string]int)
			}
			for _, value := range facetValues {
				facetCounts[facetName][value]++
			}
		}
	}

	result := make(map[string]*Facet)
	for facetName, counts := range facetCounts {
		values := make([]*FacetValue, 0)
		for value, count := range counts {
			values = append(values, &FacetValue{
				Value: value,
				Count: count,
			})
		}

		result[facetName] = &Facet{
			Name:   facetName,
			Values: values,
		}
	}

	return result
}

func (idx *InMemorySearchIndex) Suggest(prefix string, limit int) ([]string, error) {
	suggestions := make([]string, 0)
	prefixLower := strings.ToLower(prefix)

	for _, doc := range idx.documents {
		if strings.HasPrefix(strings.ToLower(doc.Title), prefixLower) {
			suggestions = append(suggestions, doc.Title)
			if len(suggestions) >= limit {
				break
			}
		}
	}

	return suggestions, nil
}
