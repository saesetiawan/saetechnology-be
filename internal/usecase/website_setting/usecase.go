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

	primaryContrastColor := strings.TrimSpace(payload.PrimaryContrastColor)
	if primaryContrastColor == "" {
		primaryContrastColor = "#ffffff"
	}

	accentContrastColor := strings.TrimSpace(payload.AccentContrastColor)
	if accentContrastColor == "" {
		accentContrastColor = "#ffffff"
	}

	surfaceContrastColor := strings.TrimSpace(payload.SurfaceContrastColor)
	if surfaceContrastColor == "" {
		surfaceContrastColor = textColor
	}

	successColor := strings.TrimSpace(payload.SuccessColor)
	if successColor == "" {
		successColor = "#10b981"
	}

	warningColor := strings.TrimSpace(payload.WarningColor)
	if warningColor == "" {
		warningColor = "#f59e0b"
	}

	dangerColor := strings.TrimSpace(payload.DangerColor)
	if dangerColor == "" {
		dangerColor = "#ef4444"
	}

	infoColor := strings.TrimSpace(payload.InfoColor)
	if infoColor == "" {
		infoColor = "#3b82f6"
	}

	labelColor := strings.TrimSpace(payload.LabelColor)
	if labelColor == "" {
		labelColor = textColor
	}

	labelBackgroundColor := strings.TrimSpace(payload.LabelBackgroundColor)
	if labelBackgroundColor == "" {
		labelBackgroundColor = secondaryColor
	}

	fontFamily := strings.TrimSpace(payload.FontFamily)
	if fontFamily == "" {
		fontFamily = "Plus Jakarta Sans, ui-sans-serif, system-ui, sans-serif"
	}

	headingFontFamily := strings.TrimSpace(payload.HeadingFontFamily)
	if headingFontFamily == "" {
		headingFontFamily = fontFamily
	}

	borderRadius := strings.TrimSpace(payload.BorderRadius)
	if borderRadius == "" {
		borderRadius = "2rem"
	}

	buttonRadius := strings.TrimSpace(payload.ButtonRadius)
	if buttonRadius == "" {
		buttonRadius = "999px"
	}

	shadowStyle := strings.TrimSpace(payload.ShadowStyle)
	if shadowStyle == "" {
		shadowStyle = "0 18px 60px rgba(15, 23, 42, 0.10)"
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
		PrimaryContrastColor: primaryContrastColor,
		AccentContrastColor:  accentContrastColor,
		SurfaceContrastColor: surfaceContrastColor,
		SuccessColor:         successColor,
		WarningColor:         warningColor,
		DangerColor:          dangerColor,
		InfoColor:            infoColor,
		LabelColor:           labelColor,
		LabelBackgroundColor: labelBackgroundColor,
		FontFamily:           fontFamily,
		HeadingFontFamily:    headingFontFamily,
		BorderRadius:         borderRadius,
		ButtonRadius:         buttonRadius,
		ShadowStyle:          shadowStyle,
		Metadata:           metadata,
	}, nil
}

func normalizeMetadata(value string) (jsonvalue.JSON, error) {
	return jsonvalue.New(strings.TrimSpace(value))
}
