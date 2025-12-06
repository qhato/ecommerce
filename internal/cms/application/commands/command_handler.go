package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/qhato/ecommerce/internal/cms/domain"
)

type CMSCommandHandler struct {
	contentRepo domain.ContentRepository
	mediaRepo   domain.MediaRepository
	versionRepo domain.ContentVersionRepository
}

func NewCMSCommandHandler(
	contentRepo domain.ContentRepository,
	mediaRepo domain.MediaRepository,
	versionRepo domain.ContentVersionRepository,
) *CMSCommandHandler {
	return &CMSCommandHandler{
		contentRepo: contentRepo,
		mediaRepo:   mediaRepo,
		versionRepo: versionRepo,
	}
}

// Content Commands

func (h *CMSCommandHandler) HandleCreateContent(ctx context.Context, cmd CreateContentCommand) (*domain.Content, error) {
	content, err := domain.NewContent(cmd.Title, cmd.Slug, domain.ContentType(cmd.Type), cmd.Body, cmd.AuthorID, cmd.Locale)
	if err != nil {
		return nil, err
	}

	content.Excerpt = cmd.Excerpt
	content.FeaturedImage = cmd.FeaturedImage
	content.MetaTitle = cmd.MetaTitle
	content.MetaDescription = cmd.MetaDescription
	content.MetaKeywords = cmd.MetaKeywords
	content.Template = cmd.Template
	content.ParentID = cmd.ParentID
	content.SortOrder = cmd.SortOrder
	content.CustomFields = cmd.CustomFields

	if err := h.contentRepo.Create(ctx, content); err != nil {
		return nil, fmt.Errorf("failed to create content: %w", err)
	}

	// Create initial version
	version := domain.NewContentVersion(content.ID, content.Body, content.AuthorID, "Initial version")
	if err := h.versionRepo.Create(ctx, version); err != nil {
		return nil, fmt.Errorf("failed to create initial version: %w", err)
	}

	return content, nil
}

func (h *CMSCommandHandler) HandleUpdateContent(ctx context.Context, cmd UpdateContentCommand) (*domain.Content, error) {
	content, err := h.contentRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}
	if content == nil {
		return nil, domain.ErrContentNotFound
	}

	if err := content.Update(cmd.Title, cmd.Slug, cmd.Body); err != nil {
		return nil, err
	}

	content.Excerpt = cmd.Excerpt
	content.FeaturedImage = cmd.FeaturedImage
	content.MetaTitle = cmd.MetaTitle
	content.MetaDescription = cmd.MetaDescription
	content.MetaKeywords = cmd.MetaKeywords
	content.Template = cmd.Template
	content.CustomFields = cmd.CustomFields
	content.Version++

	if err := h.contentRepo.Update(ctx, content); err != nil {
		return nil, fmt.Errorf("failed to update content: %w", err)
	}

	return content, nil
}

func (h *CMSCommandHandler) HandlePublishContent(ctx context.Context, cmd PublishContentCommand) (*domain.Content, error) {
	content, err := h.contentRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}
	if content == nil {
		return nil, domain.ErrContentNotFound
	}

	publishedAt := time.Now()
	if cmd.PublishedAt != nil {
		publishedAt = *cmd.PublishedAt
	}

	if err := content.Publish(publishedAt); err != nil {
		return nil, err
	}

	if err := h.contentRepo.Update(ctx, content); err != nil {
		return nil, fmt.Errorf("failed to publish content: %w", err)
	}

	return content, nil
}

func (h *CMSCommandHandler) HandleUnpublishContent(ctx context.Context, cmd UnpublishContentCommand) (*domain.Content, error) {
	content, err := h.contentRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}
	if content == nil {
		return nil, domain.ErrContentNotFound
	}

	content.Unpublish()

	if err := h.contentRepo.Update(ctx, content); err != nil {
		return nil, fmt.Errorf("failed to unpublish content: %w", err)
	}

	return content, nil
}

func (h *CMSCommandHandler) HandleArchiveContent(ctx context.Context, cmd ArchiveContentCommand) (*domain.Content, error) {
	content, err := h.contentRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}
	if content == nil {
		return nil, domain.ErrContentNotFound
	}

	content.Archive()

	if err := h.contentRepo.Update(ctx, content); err != nil {
		return nil, fmt.Errorf("failed to archive content: %w", err)
	}

	return content, nil
}

func (h *CMSCommandHandler) HandleDeleteContent(ctx context.Context, cmd DeleteContentCommand) error {
	return h.contentRepo.Delete(ctx, cmd.ID)
}

// Media Commands

func (h *CMSCommandHandler) HandleCreateMedia(ctx context.Context, cmd CreateMediaCommand) (*domain.Media, error) {
	media, err := domain.NewMedia(cmd.FileName, cmd.FilePath, cmd.MimeType, cmd.FileSize, cmd.UploadedBy)
	if err != nil {
		return nil, err
	}

	media.Title = cmd.Title
	media.AltText = cmd.AltText
	media.Caption = cmd.Caption

	if err := h.mediaRepo.Create(ctx, media); err != nil {
		return nil, fmt.Errorf("failed to create media: %w", err)
	}

	return media, nil
}

func (h *CMSCommandHandler) HandleUpdateMedia(ctx context.Context, cmd UpdateMediaCommand) (*domain.Media, error) {
	media, err := h.mediaRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}
	if media == nil {
		return nil, domain.ErrMediaNotFound
	}

	if err := media.UpdateMetadata(cmd.Title, cmd.AltText, cmd.Caption); err != nil {
		return nil, err
	}

	if err := h.mediaRepo.Update(ctx, media); err != nil {
		return nil, fmt.Errorf("failed to update media: %w", err)
	}

	return media, nil
}

func (h *CMSCommandHandler) HandleDeleteMedia(ctx context.Context, cmd DeleteMediaCommand) error {
	return h.mediaRepo.Delete(ctx, cmd.ID)
}

// Content Version Commands

func (h *CMSCommandHandler) HandleCreateContentVersion(ctx context.Context, cmd CreateContentVersionCommand) (*domain.ContentVersion, error) {
	version := domain.NewContentVersion(cmd.ContentID, cmd.Body, cmd.CreatedBy, cmd.Comment)

	if err := h.versionRepo.Create(ctx, version); err != nil {
		return nil, fmt.Errorf("failed to create content version: %w", err)
	}

	return version, nil
}
