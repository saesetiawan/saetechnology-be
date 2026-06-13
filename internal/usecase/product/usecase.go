package product

import (
	"context"
	"strings"

	cacheDomain "github.com/bagusaditiasetiawan/saetechnology-be/internal/domain/cache"
	"github.com/bagusaditiasetiawan/saetechnology-be/internal/domain/jsonvalue"
	productDomain "github.com/bagusaditiasetiawan/saetechnology-be/internal/domain/product"
	"github.com/bagusaditiasetiawan/saetechnology-be/internal/helpers"
	"github.com/bagusaditiasetiawan/saetechnology-be/internal/pkg/logger"

	"github.com/google/uuid"
)

type UseCase interface {
	Create(ctx context.Context, payload productDomain.CreateProductDto) (*productDomain.Product, error)
	Update(ctx context.Context, payload productDomain.UpdateProductDto) (*productDomain.Product, error)
	FindByID(ctx context.Context, id string) (*productDomain.Product, error)
	FindBySlug(ctx context.Context, slug string, publicOnly bool) (*productDomain.Product, error)
	FindAll(ctx context.Context, query productDomain.ListProductQuery) ([]productDomain.Product, int64, error)
	Delete(ctx context.Context, id string) error
}

type useCaseImpl struct {
	logger     logger.Logger
	repository productDomain.Repository
	cache      cacheDomain.Repository
}

func NewUseCase(
	logger logger.Logger,
	repository productDomain.Repository,
	cacheRepository cacheDomain.Repository,
) UseCase {
	return &useCaseImpl{
		logger:     logger,
		repository: repository,
		cache:      cacheRepository,
	}
}

func (uc *useCaseImpl) Create(ctx context.Context, payload productDomain.CreateProductDto) (*productDomain.Product, error) {
	entity, err := buildEntityFromCreate(payload)
	if err != nil {
		return nil, err
	}
	if err := uc.repository.Create(ctx, entity); err != nil {
		return nil, err
	}

	helpers.DeleteCacheByPrefix(ctx, uc.cache, "products:")
	return entity, nil
}

func (uc *useCaseImpl) Update(ctx context.Context, payload productDomain.UpdateProductDto) (*productDomain.Product, error) {
	entity, err := buildEntityFromUpdate(payload)
	if err != nil {
		return nil, err
	}
	if err := uc.repository.Update(ctx, entity); err != nil {
		return nil, err
	}

	helpers.DeleteCacheByPrefix(ctx, uc.cache, "products:")
	return entity, nil
}

func (uc *useCaseImpl) FindByID(ctx context.Context, id string) (*productDomain.Product, error) {
	return uc.repository.FindByID(ctx, id)
}

func (uc *useCaseImpl) FindBySlug(ctx context.Context, slug string, publicOnly bool) (*productDomain.Product, error) {
	cacheKey := "products:slug:" + slug
	if !publicOnly {
		cacheKey += ":admin"
	}
	if cached, ok := helpers.GetJSONCache[productDomain.Product](ctx, uc.cache, cacheKey); ok {
		return cached, nil
	}

	result, err := uc.repository.FindBySlug(ctx, slug, publicOnly)
	if err != nil {
		return nil, err
	}

	helpers.SetJSONCache(ctx, uc.cache, cacheKey, result, helpers.DefaultCacheTTL)
	return result, nil
}

func (uc *useCaseImpl) FindAll(
	ctx context.Context,
	query productDomain.ListProductQuery,
) ([]productDomain.Product, int64, error) {
	cacheKey := helpers.PaginationCacheKey(
		"products:"+query.Status+":"+query.Category,
		query.Page,
		query.Limit,
		query.Search,
		"",
		query.OrderBy,
		query.OrderType,
	)
	if query.PublicOnly {
		cacheKey += ":public"
	}

	if cached, ok := helpers.GetJSONCache[helpers.PaginatedCache[productDomain.Product]](ctx, uc.cache, cacheKey); ok {
		return cached.Data, cached.Total, nil
	}

	result, total, err := uc.repository.FindAll(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	helpers.SetJSONCache(ctx, uc.cache, cacheKey, helpers.PaginatedCache[productDomain.Product]{
		Data:  result,
		Total: total,
	}, helpers.DefaultCacheTTL)

	return result, total, nil
}

func (uc *useCaseImpl) Delete(ctx context.Context, id string) error {
	if err := uc.repository.Delete(ctx, id); err != nil {
		return err
	}

	helpers.DeleteCacheByPrefix(ctx, uc.cache, "products:")
	return nil
}

func buildEntityFromCreate(payload productDomain.CreateProductDto) (*productDomain.Product, error) {
	metadata, err := jsonvalue.New(strings.TrimSpace(payload.Metadata))
	if err != nil {
		return nil, err
	}

	return &productDomain.Product{
		Slug:        strings.TrimSpace(payload.Slug),
		Name:        strings.TrimSpace(payload.Name),
		Tagline:     payload.Tagline,
		Summary:     payload.Summary,
		Description: payload.Description,
		Category:    strings.TrimSpace(payload.Category),
		Status:      strings.TrimSpace(payload.Status),
		PriceLabel:  payload.PriceLabel,
		PriceURL:    payload.PriceURL,
		ImageURL:    payload.ImageURL,
		DemoURL:     payload.DemoURL,
		CTALabel:    payload.CTALabel,
		CTAURL:      payload.CTAURL,
		SortOrder:   payload.SortOrder,
		IsFeatured:  payload.IsFeatured,
		Metadata:    metadata,
	}, nil
}

func buildEntityFromUpdate(payload productDomain.UpdateProductDto) (*productDomain.Product, error) {
	id, err := uuid.Parse(payload.ID)
	if err != nil {
		return nil, err
	}
	createPayload := productDomain.CreateProductDto{
		Slug:        payload.Slug,
		Name:        payload.Name,
		Tagline:     payload.Tagline,
		Summary:     payload.Summary,
		Description: payload.Description,
		Category:    payload.Category,
		Status:      payload.Status,
		PriceLabel:  payload.PriceLabel,
		PriceURL:    payload.PriceURL,
		ImageURL:    payload.ImageURL,
		DemoURL:     payload.DemoURL,
		CTALabel:    payload.CTALabel,
		CTAURL:      payload.CTAURL,
		SortOrder:   payload.SortOrder,
		IsFeatured:  payload.IsFeatured,
		Metadata:    payload.Metadata,
	}
	entity, err := buildEntityFromCreate(createPayload)
	if err != nil {
		return nil, err
	}
	entity.ID = id

	return entity, nil
}
