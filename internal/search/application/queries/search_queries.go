package queries

import (
	"context"
	"fmt"
	"strings"

	"github.com/qhato/ecommerce/internal/search/domain"
)

// SearchEngine defines the interface for search operations
type SearchEngine interface {
	Search(ctx context.Context, query *domain.SearchQuery) (*domain.SearchResult, error)
	Suggest(ctx context.Context, prefix string, limit int) ([]string, error)
}

// SearchQueryService provides query operations for search
type SearchQueryService struct {
	searchEngine    SearchEngine
	synonymRepo     domain.SearchSynonymRepository
	redirectRepo    domain.SearchRedirectRepository
	facetConfigRepo domain.SearchFacetConfigRepository
	indexingJobRepo domain.IndexingJobRepository
}

// NewSearchQueryService creates a new search query service
func NewSearchQueryService(
	searchEngine SearchEngine,
	synonymRepo domain.SearchSynonymRepository,
	redirectRepo domain.SearchRedirectRepository,
	facetConfigRepo domain.SearchFacetConfigRepository,
	indexingJobRepo domain.IndexingJobRepository,
) *SearchQueryService {
	return &SearchQueryService{
		searchEngine:    searchEngine,
		synonymRepo:     synonymRepo,
		redirectRepo:    redirectRepo,
		facetConfigRepo: facetConfigRepo,
		indexingJobRepo: indexingJobRepo,
	}
}

// SearchProducts searches for products
func (s *SearchQueryService) SearchProducts(ctx context.Context, query *domain.SearchQuery) (*SearchResultDTO, error) {
	// Check for search redirects first
	if query.Query != "" {
		redirect, err := s.redirectRepo.FindBySearchTerm(strings.ToLower(query.Query))
		if err == nil && redirect != nil && redirect.IsCurrentlyActive() {
			return &SearchResultDTO{
				RedirectURL: redirect.TargetURL,
				Query:       query.Query,
			}, nil
		}
	}

	// Expand query with synonyms
	expandedQuery := s.expandQueryWithSynonyms(ctx, query.Query)
	if expandedQuery != query.Query {
		query.Query = expandedQuery
	}

	// Execute search
	result, err := s.searchEngine.Search(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}

	// Convert to DTO
	return s.mapResultToDTO(result), nil
}

// Suggest provides search suggestions
func (s *SearchQueryService) Suggest(ctx context.Context, prefix string, limit int) ([]string, error) {
	suggestions, err := s.searchEngine.Suggest(ctx, prefix, limit)
	if err != nil {
		return nil, fmt.Errorf("suggest failed: %w", err)
	}
	return suggestions, nil
}

// GetSynonymByID retrieves a synonym by ID
func (s *SearchQueryService) GetSynonymByID(ctx context.Context, id int64) (*SynonymDTO, error) {
	synonym, err := s.synonymRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return &SynonymDTO{
		ID:       synonym.ID,
		Term:     synonym.Term,
		Synonyms: synonym.Synonyms,
		IsActive: synonym.IsActive,
	}, nil
}

// ListSynonyms lists all synonyms
func (s *SearchQueryService) ListSynonyms(ctx context.Context) ([]*SynonymDTO, error) {
	synonyms, err := s.synonymRepo.FindAll()
	if err != nil {
		return nil, err
	}

	dtos := make([]*SynonymDTO, len(synonyms))
	for i, syn := range synonyms {
		dtos[i] = &SynonymDTO{
			ID:       syn.ID,
			Term:     syn.Term,
			Synonyms: syn.Synonyms,
			IsActive: syn.IsActive,
		}
	}
	return dtos, nil
}

// GetRedirectByID retrieves a redirect by ID
func (s *SearchQueryService) GetRedirectByID(ctx context.Context, id int64) (*RedirectDTO, error) {
	redirect, err := s.redirectRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return &RedirectDTO{
		ID:             redirect.ID,
		SearchTerm:     redirect.SearchTerm,
		TargetURL:      redirect.TargetURL,
		Priority:       redirect.Priority,
		IsActive:       redirect.IsActive,
		ActivationDate: redirect.ActivationDate,
		ExpirationDate: redirect.ExpirationDate,
	}, nil
}

// ListRedirects lists all redirects
func (s *SearchQueryService) ListRedirects(ctx context.Context) ([]*RedirectDTO, error) {
	redirects, err := s.redirectRepo.FindAllActive()
	if err != nil {
		return nil, err
	}

	dtos := make([]*RedirectDTO, len(redirects))
	for i, redir := range redirects {
		dtos[i] = &RedirectDTO{
			ID:             redir.ID,
			SearchTerm:     redir.SearchTerm,
			TargetURL:      redir.TargetURL,
			Priority:       redir.Priority,
			IsActive:       redir.IsActive,
			ActivationDate: redir.ActivationDate,
			ExpirationDate: redir.ExpirationDate,
		}
	}
	return dtos, nil
}

// GetFacetConfigByID retrieves a facet config by ID
func (s *SearchQueryService) GetFacetConfigByID(ctx context.Context, id int64) (*FacetConfigDTO, error) {
	config, err := s.facetConfigRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return &FacetConfigDTO{
		ID:               config.ID,
		Name:             config.Name,
		Label:            config.Label,
		FieldName:        config.FieldName,
		FacetType:        string(config.FacetType),
		IsActive:         config.IsActive,
		ShowInResults:    config.ShowInResults,
		ShowInNavigation: config.ShowInNavigation,
		Priority:         config.Priority,
		MinDocCount:      config.MinDocCount,
		MaxValues:        config.MaxValues,
	}, nil
}

// ListFacetConfigs lists all facet configs
func (s *SearchQueryService) ListFacetConfigs(ctx context.Context) ([]*FacetConfigDTO, error) {
	configs, err := s.facetConfigRepo.FindActive()
	if err != nil {
		return nil, err
	}

	dtos := make([]*FacetConfigDTO, len(configs))
	for i, cfg := range configs {
		dtos[i] = &FacetConfigDTO{
			ID:               cfg.ID,
			Name:             cfg.Name,
			Label:            cfg.Label,
			FieldName:        cfg.FieldName,
			FacetType:        string(cfg.FacetType),
			IsActive:         cfg.IsActive,
			ShowInResults:    cfg.ShowInResults,
			ShowInNavigation: cfg.ShowInNavigation,
			Priority:         cfg.Priority,
			MinDocCount:      cfg.MinDocCount,
			MaxValues:        cfg.MaxValues,
		}
	}
	return dtos, nil
}

// GetIndexingJobByID retrieves an indexing job by ID
func (s *SearchQueryService) GetIndexingJobByID(ctx context.Context, id int64) (*IndexingJobDTO, error) {
	job, err := s.indexingJobRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return s.mapJobToDTO(job), nil
}

// ListRecentIndexingJobs lists recent indexing jobs
func (s *SearchQueryService) ListRecentIndexingJobs(ctx context.Context, limit int) ([]*IndexingJobDTO, error) {
	jobs, err := s.indexingJobRepo.FindRecent(limit)
	if err != nil {
		return nil, err
	}

	dtos := make([]*IndexingJobDTO, len(jobs))
	for i, job := range jobs {
		dtos[i] = s.mapJobToDTO(job)
	}
	return dtos, nil
}

// expandQueryWithSynonyms expands a query with synonyms
func (s *SearchQueryService) expandQueryWithSynonyms(ctx context.Context, query string) string {
	if query == "" {
		return query
	}

	// Get active synonyms
	synonyms, err := s.synonymRepo.FindActive()
	if err != nil {
		return query
	}

	// Check if query matches any synonym
	queryLower := strings.ToLower(query)
	for _, syn := range synonyms {
		if strings.EqualFold(syn.Term, queryLower) {
			// Expand with all synonyms
			terms := syn.GetExpandedTerms()
			return strings.Join(terms, " OR ")
		}
		for _, synTerm := range syn.Synonyms {
			if strings.EqualFold(synTerm, queryLower) {
				terms := syn.GetExpandedTerms()
				return strings.Join(terms, " OR ")
			}
		}
	}

	return query
}

// mapResultToDTO maps SearchResult to SearchResultDTO
func (s *SearchQueryService) mapResultToDTO(result *domain.SearchResult) *SearchResultDTO {
	products := make([]*ProductSearchDTO, len(result.Documents))
	for i, doc := range result.Documents {
		products[i] = &ProductSearchDTO{
			ProductID:    doc.Fields["product_id"].(int64),
			SKU:          doc.Fields["sku"].(string),
			Name:         doc.Title,
			Description:  doc.Description,
			Price:        doc.Fields["price"].(float64),
			ImageURL:     doc.Fields["image_url"].(string),
			IsAvailable:  doc.Fields["is_available"].(bool),
			Rating:       doc.Fields["rating"].(float64),
			ReviewCount:  doc.Fields["review_count"].(int),
			Score:        doc.Score,
		}
	}

	facets := make(map[string]*FacetDTO)
	for name, facet := range result.Facets {
		values := make([]*FacetValueDTO, len(facet.Values))
		for i, val := range facet.Values {
			values[i] = &FacetValueDTO{
				Value: val.Value,
				Count: val.Count,
			}
		}
		facets[name] = &FacetDTO{
			Name:   name,
			Values: values,
		}
	}

	return &SearchResultDTO{
		Products:   products,
		Total:      result.Total,
		Facets:     facets,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
		Query:      result.Query,
	}
}

// mapJobToDTO maps IndexingJob to IndexingJobDTO
func (s *SearchQueryService) mapJobToDTO(job *domain.IndexingJob) *IndexingJobDTO {
	dto := &IndexingJobDTO{
		ID:             job.ID,
		Type:           string(job.Type),
		Status:         string(job.Status),
		EntityType:     job.EntityType,
		TotalItems:     job.TotalItems,
		ProcessedItems: job.ProcessedItems,
		FailedItems:    job.FailedItems,
		ErrorMessage:   job.ErrorMessage,
		StartedAt:      job.StartedAt,
		CompletedAt:    job.CompletedAt,
		CreatedAt:      job.CreatedAt,
		Progress:       job.GetProgress(),
	}

	if !job.IsCompleted() && job.StartedAt != nil {
		duration := job.GetDuration()
		dto.Duration = duration.String()
	} else if job.CompletedAt != nil && job.StartedAt != nil {
		duration := job.GetDuration()
		dto.Duration = duration.String()
	}

	return dto
}
