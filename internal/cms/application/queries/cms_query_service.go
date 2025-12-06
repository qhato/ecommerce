package queries

import (
	"context"
	"fmt"

	"github.com/qhato/ecommerce/internal/cms/domain"
)

// CMSQueryService handles CMS-related queries
type CMSQueryService struct {
	contentRepo domain.ContentRepository
	mediaRepo   domain.MediaRepository
	versionRepo domain.ContentVersionRepository
}

// NewCMSQueryService creates a new CMS query service
func NewCMSQueryService(
	contentRepo domain.ContentRepository,
	mediaRepo domain.MediaRepository,
	versionRepo domain.ContentVersionRepository,
) *CMSQueryService {
	return &CMSQueryService{
		contentRepo: contentRepo,
		mediaRepo:   mediaRepo,
		versionRepo: versionRepo,
	}
}

// GetContent retrieves content by ID
func (s *CMSQueryService) GetContent(ctx context.Context, id int64) (*ContentDTO, error) {
	content, err := s.contentRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find content: %w", err)
	}
	if content == nil {
		return nil, domain.ErrContentNotFound
	}

	return ToContentDTO(content), nil
}

// GetContentBySlug retrieves content by slug
func (s *CMSQueryService) GetContentBySlug(ctx context.Context, slug, locale string) (*ContentDTO, error) {
	content, err := s.contentRepo.FindBySlug(ctx, slug, locale)
	if err != nil {
		return nil, fmt.Errorf("failed to find content: %w", err)
	}
	if content == nil {
		return nil, domain.ErrContentNotFound
	}

	// Increment view count
	s.contentRepo.IncrementViewCount(ctx, content.ID)

	return ToContentDTO(content), nil
}

// GetContentByType retrieves content by type
func (s *CMSQueryService) GetContentByType(ctx context.Context, contentType string, locale string, publishedOnly bool) ([]*ContentDTO, error) {
	contents, err := s.contentRepo.FindByType(ctx, domain.ContentType(contentType), locale, publishedOnly)
	if err != nil {
		return nil, fmt.Errorf("failed to find content: %w", err)
	}

	dtos := make([]*ContentDTO, len(contents))
	for i, content := range contents {
		dtos[i] = ToContentDTO(content)
	}

	return dtos, nil
}

// GetAllContent retrieves all content
func (s *CMSQueryService) GetAllContent(ctx context.Context, locale string, publishedOnly bool) ([]*ContentDTO, error) {
	contents, err := s.contentRepo.FindAll(ctx, locale, publishedOnly)
	if err != nil {
		return nil, fmt.Errorf("failed to find content: %w", err)
	}

	dtos := make([]*ContentDTO, len(contents))
	for i, content := range contents {
		dtos[i] = ToContentDTO(content)
	}

	return dtos, nil
}

// GetContentHierarchy retrieves content hierarchy
func (s *CMSQueryService) GetContentHierarchy(ctx context.Context, contentType string, locale string) ([]*ContentDTO, error) {
	contents, err := s.contentRepo.FindHierarchy(ctx, domain.ContentType(contentType), locale)
	if err != nil {
		return nil, fmt.Errorf("failed to find content hierarchy: %w", err)
	}

	dtos := make([]*ContentDTO, len(contents))
	for i, content := range contents {
		dtos[i] = ToContentDTO(content)
	}

	return dtos, nil
}

// GetContentChildren retrieves child content
func (s *CMSQueryService) GetContentChildren(ctx context.Context, parentID int64) ([]*ContentDTO, error) {
	contents, err := s.contentRepo.FindChildren(ctx, parentID)
	if err != nil {
		return nil, fmt.Errorf("failed to find children: %w", err)
	}

	dtos := make([]*ContentDTO, len(contents))
	for i, content := range contents {
		dtos[i] = ToContentDTO(content)
	}

	return dtos, nil
}

// SearchContent searches content by query
func (s *CMSQueryService) SearchContent(ctx context.Context, query, locale string, publishedOnly bool) ([]*ContentDTO, error) {
	contents, err := s.contentRepo.Search(ctx, query, locale, publishedOnly)
	if err != nil {
		return nil, fmt.Errorf("failed to search content: %w", err)
	}

	dtos := make([]*ContentDTO, len(contents))
	for i, content := range contents {
		dtos[i] = ToContentDTO(content)
	}

	return dtos, nil
}

// GetContentVersions retrieves all versions of content
func (s *CMSQueryService) GetContentVersions(ctx context.Context, contentID int64) ([]*ContentVersionDTO, error) {
	versions, err := s.versionRepo.FindByContentID(ctx, contentID)
	if err != nil {
		return nil, fmt.Errorf("failed to find versions: %w", err)
	}

	dtos := make([]*ContentVersionDTO, len(versions))
	for i, version := range versions {
		dtos[i] = ToContentVersionDTO(version)
	}

	return dtos, nil
}

// GetContentVersion retrieves a specific version
func (s *CMSQueryService) GetContentVersion(ctx context.Context, versionID int64) (*ContentVersionDTO, error) {
	version, err := s.versionRepo.FindByID(ctx, versionID)
	if err != nil {
		return nil, fmt.Errorf("failed to find version: %w", err)
	}
	if version == nil {
		return nil, domain.ErrVersionNotFound
	}

	return ToContentVersionDTO(version), nil
}

// GetMedia retrieves media by ID
func (s *CMSQueryService) GetMedia(ctx context.Context, id int64) (*MediaDTO, error) {
	media, err := s.mediaRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find media: %w", err)
	}
	if media == nil {
		return nil, domain.ErrMediaNotFound
	}

	return ToMediaDTO(media), nil
}

// GetAllMedia retrieves all media
func (s *CMSQueryService) GetAllMedia(ctx context.Context, mimeType string, limit int) ([]*MediaDTO, error) {
	medias, err := s.mediaRepo.FindAll(ctx, mimeType, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to find media: %w", err)
	}

	dtos := make([]*MediaDTO, len(medias))
	for i, media := range medias {
		dtos[i] = ToMediaDTO(media)
	}

	return dtos, nil
}

// GetMediaByUploader retrieves media by uploader
func (s *CMSQueryService) GetMediaByUploader(ctx context.Context, uploaderID int64, limit int) ([]*MediaDTO, error) {
	medias, err := s.mediaRepo.FindByUploader(ctx, uploaderID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to find media: %w", err)
	}

	dtos := make([]*MediaDTO, len(medias))
	for i, media := range medias {
		dtos[i] = ToMediaDTO(media)
	}

	return dtos, nil
}
