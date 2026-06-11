package content

import (
	"context"
	"strings"
	"time"

	cacheDomain "go-platform-core/internal/domain/cache"
	contentDomain "go-platform-core/internal/domain/content"
	"go-platform-core/internal/domain/jsonvalue"
	"go-platform-core/internal/helpers"
	"go-platform-core/internal/pkg/logger"

	"github.com/google/uuid"
)

type UseCase interface {
	Create(ctx context.Context, payload contentDomain.CreateWebsiteContentDto) (*contentDomain.WebsiteContent, error)
	Update(ctx context.Context, payload contentDomain.UpdateWebsiteContentDto) (*contentDomain.WebsiteContent, error)
	FindByID(ctx context.Context, id string) (*contentDomain.WebsiteContent, error)
	FindByKey(ctx context.Context, key string) (*contentDomain.WebsiteContent, error)
	FindAll(ctx context.Context, query contentDomain.ListWebsiteContentQuery) ([]contentDomain.WebsiteContent, int64, error)
}

type useCaseImpl struct {
	logger     logger.Logger
	repository contentDomain.Repository
	cache      cacheDomain.Repository
}

func NewUseCase(
	logger logger.Logger,
	repository contentDomain.Repository,
	cacheRepository cacheDomain.Repository,
) UseCase {
	return &useCaseImpl{
		logger:     logger,
		repository: repository,
		cache:      cacheRepository,
	}
}

func (uc *useCaseImpl) Create(
	ctx context.Context,
	payload contentDomain.CreateWebsiteContentDto,
) (*contentDomain.WebsiteContent, error) {
	uc.logger.Info("process create website content")

	entity, err := buildEntityFromCreate(payload)
	if err != nil {
		return nil, err
	}

	if err := uc.repository.Create(ctx, entity); err != nil {
		return nil, err
	}

	helpers.DeleteCacheByPrefix(ctx, uc.cache, "website-contents:")
	return entity, nil
}

func (uc *useCaseImpl) Update(
	ctx context.Context,
	payload contentDomain.UpdateWebsiteContentDto,
) (*contentDomain.WebsiteContent, error) {
	uc.logger.Info("process update website content")

	entity, err := buildEntityFromUpdate(payload)
	if err != nil {
		return nil, err
	}

	if err := uc.repository.Update(ctx, entity); err != nil {
		return nil, err
	}

	helpers.DeleteCacheByPrefix(ctx, uc.cache, "website-contents:")
	return entity, nil
}

func (uc *useCaseImpl) FindByID(
	ctx context.Context,
	id string,
) (*contentDomain.WebsiteContent, error) {
	return uc.repository.FindByID(ctx, id)
}

func (uc *useCaseImpl) FindByKey(
	ctx context.Context,
	key string,
) (*contentDomain.WebsiteContent, error) {
	cacheKey := "website-contents:key:" + key
	if cached, ok := helpers.GetJSONCache[contentDomain.WebsiteContent](ctx, uc.cache, cacheKey); ok {
		return cached, nil
	}

	result, err := uc.repository.FindByKey(ctx, key)
	if err != nil {
		return nil, err
	}

	helpers.SetJSONCache(ctx, uc.cache, cacheKey, result, helpers.DefaultCacheTTL)
	return result, nil
}

func (uc *useCaseImpl) FindAll(
	ctx context.Context,
	query contentDomain.ListWebsiteContentQuery,
) ([]contentDomain.WebsiteContent, int64, error) {
	cacheKey := helpers.PaginationCacheKey(
		"website-contents:"+query.Type+":"+query.Placement,
		query.Page,
		query.Limit,
		query.Search,
		query.SearchBy,
		query.OrderBy,
		query.OrderType,
	)
	if query.Active != nil {
		cacheKey += ":active:" + boolCacheValue(*query.Active)
	}

	if cached, ok := helpers.GetJSONCache[helpers.PaginatedCache[contentDomain.WebsiteContent]](ctx, uc.cache, cacheKey); ok {
		return cached.Data, cached.Total, nil
	}

	result, total, err := uc.repository.FindAll(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	helpers.SetJSONCache(ctx, uc.cache, cacheKey, helpers.PaginatedCache[contentDomain.WebsiteContent]{
		Data:  result,
		Total: total,
	}, helpers.DefaultCacheTTL)

	return result, total, nil
}

func buildEntityFromCreate(
	payload contentDomain.CreateWebsiteContentDto,
) (*contentDomain.WebsiteContent, error) {
	metadata, err := normalizeMetadata(payload.Metadata)
	if err != nil {
		return nil, err
	}

	publishStartAt, err := parseOptionalTime(payload.PublishStartAt)
	if err != nil {
		return nil, err
	}

	publishEndAt, err := parseOptionalTime(payload.PublishEndAt)
	if err != nil {
		return nil, err
	}

	return &contentDomain.WebsiteContent{
		Key:            strings.TrimSpace(payload.Key),
		Type:           strings.TrimSpace(payload.Type),
		Placement:      strings.TrimSpace(payload.Placement),
		Title:          strings.TrimSpace(payload.Title),
		Subtitle:       payload.Subtitle,
		Body:           payload.Body,
		ImageURL:       payload.ImageURL,
		LinkURL:        payload.LinkURL,
		LinkLabel:      payload.LinkLabel,
		SortOrder:      payload.SortOrder,
		IsActive:       payload.IsActive,
		Metadata:       metadata,
		PublishStartAt: publishStartAt,
		PublishEndAt:   publishEndAt,
	}, nil
}

func buildEntityFromUpdate(
	payload contentDomain.UpdateWebsiteContentDto,
) (*contentDomain.WebsiteContent, error) {
	id, err := uuid.Parse(payload.ID)
	if err != nil {
		return nil, err
	}

	metadata, err := normalizeMetadata(payload.Metadata)
	if err != nil {
		return nil, err
	}

	publishStartAt, err := parseOptionalTime(payload.PublishStartAt)
	if err != nil {
		return nil, err
	}

	publishEndAt, err := parseOptionalTime(payload.PublishEndAt)
	if err != nil {
		return nil, err
	}

	return &contentDomain.WebsiteContent{
		ID:             id,
		Key:            strings.TrimSpace(payload.Key),
		Type:           strings.TrimSpace(payload.Type),
		Placement:      strings.TrimSpace(payload.Placement),
		Title:          strings.TrimSpace(payload.Title),
		Subtitle:       payload.Subtitle,
		Body:           payload.Body,
		ImageURL:       payload.ImageURL,
		LinkURL:        payload.LinkURL,
		LinkLabel:      payload.LinkLabel,
		SortOrder:      payload.SortOrder,
		IsActive:       payload.IsActive,
		Metadata:       metadata,
		PublishStartAt: publishStartAt,
		PublishEndAt:   publishEndAt,
	}, nil
}

func normalizeMetadata(value string) (jsonvalue.JSON, error) {
	return jsonvalue.New(strings.TrimSpace(value))
}

func parseOptionalTime(value string) (*time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, nil
	}

	parsed, err := time.Parse(time.RFC3339, value)
	if err == nil {
		return &parsed, nil
	}

	parsed, err = time.Parse("2006-01-02", value)
	if err != nil {
		return nil, err
	}

	return &parsed, nil
}

func boolCacheValue(value bool) string {
	if value {
		return "true"
	}

	return "false"
}
