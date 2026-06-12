package website_setting

import (
	"context"
	"strings"

	cacheDomain "github.com/bagusaditiasetiawan/saetechnology-be/internal/domain/cache"
	"github.com/bagusaditiasetiawan/saetechnology-be/internal/domain/jsonvalue"
	websiteSettingDomain "github.com/bagusaditiasetiawan/saetechnology-be/internal/domain/website_setting"
	"github.com/bagusaditiasetiawan/saetechnology-be/internal/helpers"
	"github.com/bagusaditiasetiawan/saetechnology-be/internal/pkg/logger"
)

type UseCase interface {
	Find(ctx context.Context) (*websiteSettingDomain.WebsiteSetting, error)
	Update(ctx context.Context, payload websiteSettingDomain.UpdateWebsiteSettingDto) (*websiteSettingDomain.WebsiteSetting, error)
}

type useCaseImpl struct {
	logger     logger.Logger
	repository websiteSettingDomain.Repository
	cache      cacheDomain.Repository
}

func NewUseCase(
	logger logger.Logger,
	repository websiteSettingDomain.Repository,
	cacheRepository cacheDomain.Repository,
) UseCase {
	return &useCaseImpl{
		logger:     logger,
		repository: repository,
		cache:      cacheRepository,
	}
}

func (uc *useCaseImpl) Find(
	ctx context.Context,
) (*websiteSettingDomain.WebsiteSetting, error) {
	cacheKey := "website-settings:public"
	if cached, ok := helpers.GetJSONCache[websiteSettingDomain.WebsiteSetting](ctx, uc.cache, cacheKey); ok {
		return cached, nil
	}

	result, err := uc.repository.Find(ctx)
	if err != nil {
		return nil, err
	}

	helpers.SetJSONCache(ctx, uc.cache, cacheKey, result, helpers.DefaultCacheTTL)
	return result, nil
}

func (uc *useCaseImpl) Update(
	ctx context.Context,
	payload websiteSettingDomain.UpdateWebsiteSettingDto,
) (*websiteSettingDomain.WebsiteSetting, error) {
	uc.logger.Info("process update website setting")

	entity, err := buildEntity(payload)
	if err != nil {
		return nil, err
	}

	if err := uc.repository.Upsert(ctx, entity); err != nil {
		return nil, err
	}

	helpers.DeleteCacheByPrefix(ctx, uc.cache, "website-settings:")
	return uc.repository.Find(ctx)
}

func buildEntity(
	payload websiteSettingDomain.UpdateWebsiteSettingDto,
) (*websiteSettingDomain.WebsiteSetting, error) {
	metadata, err := normalizeMetadata(payload.Metadata)
	if err != nil {
		return nil, err
	}

	primaryColor := strings.TrimSpace(payload.PrimaryColor)
	if primaryColor == "" {
		primaryColor = "#ec4899"
	}

	accentColor := strings.TrimSpace(payload.AccentColor)
	if accentColor == "" {
		accentColor = "#06b6d4"
	}

	secondaryColor := strings.TrimSpace(payload.SecondaryColor)
	if secondaryColor == "" {
		secondaryColor = accentColor
	}

	backgroundColor := strings.TrimSpace(payload.BackgroundColor)
	if backgroundColor == "" {
		backgroundColor = "#fff7fb"
	}

	surfaceColor := strings.TrimSpace(payload.SurfaceColor)
	if surfaceColor == "" {
		surfaceColor = "#ffffff"
	}

	textColor := strings.TrimSpace(payload.TextColor)
	if textColor == "" {
		textColor = "#0f172a"
	}

	mutedTextColor := strings.TrimSpace(payload.MutedTextColor)
	if mutedTextColor == "" {
		mutedTextColor = "#64748b"
	}

	borderColor := strings.TrimSpace(payload.BorderColor)
	if borderColor == "" {
		borderColor = "#e2e8f0"
	}

	return &websiteSettingDomain.WebsiteSetting{
		SiteName:           strings.TrimSpace(payload.SiteName),
		Tagline:            strings.TrimSpace(payload.Tagline),
		LogoURL:            strings.TrimSpace(payload.LogoURL),
		FaviconURL:         strings.TrimSpace(payload.FaviconURL),
		PrimaryImageURL:    strings.TrimSpace(payload.PrimaryImageURL),
		SecondaryImageURL:  strings.TrimSpace(payload.SecondaryImageURL),
		BackgroundImageURL: strings.TrimSpace(payload.BackgroundImageURL),
		Email:              strings.TrimSpace(payload.Email),
		Phone:              strings.TrimSpace(payload.Phone),
		Address:            strings.TrimSpace(payload.Address),
		FacebookURL:        strings.TrimSpace(payload.FacebookURL),
		InstagramURL:       strings.TrimSpace(payload.InstagramURL),
		TiktokURL:          strings.TrimSpace(payload.TiktokURL),
		PrimaryColor:       primaryColor,
		SecondaryColor:     secondaryColor,
		AccentColor:        accentColor,
		BackgroundColor:    backgroundColor,
		SurfaceColor:       surfaceColor,
		TextColor:          textColor,
		MutedTextColor:     mutedTextColor,
		BorderColor:        borderColor,
		Metadata:           metadata,
	}, nil
}

func normalizeMetadata(value string) (jsonvalue.JSON, error) {
	return jsonvalue.New(strings.TrimSpace(value))
}
