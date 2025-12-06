package queries

import (
	"context"
	"fmt"

	"github.com/qhato/ecommerce/internal/media/domain"
)

type MediaQueryService struct {
	mediaRepo domain.MediaRepository
}

func NewMediaQueryService(mediaRepo domain.MediaRepository) *MediaQueryService {
	return &MediaQueryService{mediaRepo: mediaRepo}
}

func (s *MediaQueryService) GetMedia(ctx context.Context, id string) (*MediaDTO, error) {
	media, err := s.mediaRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find media: %w", err)
	}
	if media == nil {
		return nil, domain.ErrMediaNotFound
	}

	return ToMediaDTO(media), nil
}

func (s *MediaQueryService) GetMediaByEntityID(ctx context.Context, entityType, entityID string) ([]*MediaDTO, error) {
	medias, err := s.mediaRepo.FindByEntityID(ctx, entityType, entityID)
	if err != nil {
		return nil, fmt.Errorf("failed to find media by entity: %w", err)
	}

	dtos := make([]*MediaDTO, len(medias))
	for i, media := range medias {
		dtos[i] = ToMediaDTO(media)
	}

	return dtos, nil
}

func (s *MediaQueryService) GetMediaByType(ctx context.Context, mediaType string) ([]*MediaDTO, error) {
	medias, err := s.mediaRepo.FindByType(ctx, domain.MediaType(mediaType))
	if err != nil {
		return nil, fmt.Errorf("failed to find media by type: %w", err)
	}

	dtos := make([]*MediaDTO, len(medias))
	for i, media := range medias {
		dtos[i] = ToMediaDTO(media)
	}

	return dtos, nil
}

func (s *MediaQueryService) GetMediaByTags(ctx context.Context, tags []string) ([]*MediaDTO, error) {
	medias, err := s.mediaRepo.FindByTags(ctx, tags)
	if err != nil {
		return nil, fmt.Errorf("failed to find media by tags: %w", err)
	}

	dtos := make([]*MediaDTO, len(medias))
	for i, media := range medias {
		dtos[i] = ToMediaDTO(media)
	}

	return dtos, nil
}

func (s *MediaQueryService) ListMedia(ctx context.Context, limit, offset int) ([]*MediaDTO, error) {
	medias, err := s.mediaRepo.FindAll(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list media: %w", err)
	}

	dtos := make([]*MediaDTO, len(medias))
	for i, media := range medias {
		dtos[i] = ToMediaDTO(media)
	}

	return dtos, nil
}
