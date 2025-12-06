package commands

import (
	"context"
	"fmt"

	"github.com/qhato/ecommerce/internal/media/domain"
)

type MediaCommandHandler struct {
	mediaRepo domain.MediaRepository
}

func NewMediaCommandHandler(mediaRepo domain.MediaRepository) *MediaCommandHandler {
	return &MediaCommandHandler{mediaRepo: mediaRepo}
}

func (h *MediaCommandHandler) HandleCreateMedia(ctx context.Context, cmd CreateMediaCommand) (*domain.Media, error) {
	media, err := domain.NewMedia(cmd.Name, cmd.MimeType, cmd.FilePath, cmd.UploadedBy, cmd.FileSize)
	if err != nil {
		return nil, err
	}

	media.Title = cmd.Title
	media.Description = cmd.Description

	if cmd.EntityType != nil && cmd.EntityID != nil {
		media.AttachToEntity(*cmd.EntityType, *cmd.EntityID)
	}

	for _, tag := range cmd.Tags {
		media.AddTag(tag)
	}

	if err := h.mediaRepo.Create(ctx, media); err != nil {
		return nil, fmt.Errorf("failed to create media: %w", err)
	}

	return media, nil
}

func (h *MediaCommandHandler) HandleUpdateMedia(ctx context.Context, cmd UpdateMediaCommand) (*domain.Media, error) {
	media, err := h.mediaRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find media: %w", err)
	}
	if media == nil {
		return nil, domain.ErrMediaNotFound
	}

	media.UpdateInfo(cmd.Title, cmd.Description)

	// Update tags
	media.Tags = cmd.Tags

	if err := h.mediaRepo.Update(ctx, media); err != nil {
		return nil, fmt.Errorf("failed to update media: %w", err)
	}

	return media, nil
}

func (h *MediaCommandHandler) HandleSetMediaURL(ctx context.Context, cmd SetMediaURLCommand) (*domain.Media, error) {
	media, err := h.mediaRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find media: %w", err)
	}
	if media == nil {
		return nil, domain.ErrMediaNotFound
	}

	media.SetURL(cmd.URL)
	if cmd.ThumbnailURL != "" {
		media.SetThumbnail(cmd.ThumbnailURL)
	}

	if err := h.mediaRepo.Update(ctx, media); err != nil {
		return nil, fmt.Errorf("failed to update media URL: %w", err)
	}

	return media, nil
}

func (h *MediaCommandHandler) HandleActivateMedia(ctx context.Context, cmd ActivateMediaCommand) (*domain.Media, error) {
	media, err := h.mediaRepo.FindByID(ctx, cmd.ID)
	if err != nil || media == nil {
		return nil, domain.ErrMediaNotFound
	}

	if err := media.Activate(); err != nil {
		return nil, err
	}

	if err := h.mediaRepo.Update(ctx, media); err != nil {
		return nil, fmt.Errorf("failed to activate media: %w", err)
	}

	return media, nil
}

func (h *MediaCommandHandler) HandleArchiveMedia(ctx context.Context, cmd ArchiveMediaCommand) (*domain.Media, error) {
	media, err := h.mediaRepo.FindByID(ctx, cmd.ID)
	if err != nil || media == nil {
		return nil, domain.ErrMediaNotFound
	}

	media.Archive()

	if err := h.mediaRepo.Update(ctx, media); err != nil {
		return nil, fmt.Errorf("failed to archive media: %w", err)
	}

	return media, nil
}

func (h *MediaCommandHandler) HandleDeleteMedia(ctx context.Context, cmd DeleteMediaCommand) error {
	media, err := h.mediaRepo.FindByID(ctx, cmd.ID)
	if err != nil || media == nil {
		return domain.ErrMediaNotFound
	}

	media.Delete()
	return h.mediaRepo.Update(ctx, media)
}
