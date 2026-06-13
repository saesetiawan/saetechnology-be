package website_setting

import (
	"context"
	"errors"
	"time"

	websiteSettingDomain "github.com/bagusaditiasetiawan/saetechnology-be/internal/domain/website_setting"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

type PostgresqlRepository struct {
	db     *gorm.DB
	tracer trace.Tracer
}

func NewPostgresqlRepository(
	db *gorm.DB,
	tracerProvider trace.TracerProvider,
) websiteSettingDomain.Repository {
	return &PostgresqlRepository{
		db:     db,
		tracer: tracerProvider.Tracer("WebsiteSettingRepository"),
	}
}

func (r *PostgresqlRepository) Find(
	ctx context.Context,
) (*websiteSettingDomain.WebsiteSetting, error) {
	ctx, span := r.tracer.Start(ctx, "WebsiteSettingRepository.Find")
	defer span.End()

	var result websiteSettingDomain.WebsiteSetting
	if err := r.db.WithContext(ctx).Order("created_at ASC").First(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			span.SetStatus(codes.Error, "website setting not found")
			return nil, errors.New("website setting not found")
		}

		span.RecordError(err)
		span.SetStatus(codes.Error, "failed find website setting")
		return nil, err
	}

	span.SetStatus(codes.Ok, "success find website setting")
	return &result, nil
}

func (r *PostgresqlRepository) Upsert(
	ctx context.Context,
	payload *websiteSettingDomain.WebsiteSetting,
) error {
	ctx, span := r.tracer.Start(ctx, "WebsiteSettingRepository.Upsert")
	defer span.End()

	existing, err := r.Find(ctx)
	if err != nil && err.Error() != "website setting not found" {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed find existing website setting")
		return err
	}

	if existing == nil {
		id, err := uuid.NewV7()
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, "failed generate website setting id")
			return err
		}

		payload.ID = id
		if err := r.db.WithContext(ctx).Create(payload).Error; err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, "failed create website setting")
			return err
		}

		span.SetStatus(codes.Ok, "success create website setting")
		return nil
	}

	updates := map[string]interface{}{
		"site_name":            payload.SiteName,
		"tagline":              payload.Tagline,
		"logo_url":             payload.LogoURL,
		"favicon_url":          payload.FaviconURL,
		"primary_image_url":    payload.PrimaryImageURL,
		"secondary_image_url":  payload.SecondaryImageURL,
		"background_image_url": payload.BackgroundImageURL,
		"email":                payload.Email,
		"phone":                payload.Phone,
		"address":              payload.Address,
		"facebook_url":         payload.FacebookURL,
		"instagram_url":        payload.InstagramURL,
		"tiktok_url":           payload.TiktokURL,
		"primary_color":        payload.PrimaryColor,
		"secondary_color":      payload.SecondaryColor,
		"accent_color":         payload.AccentColor,
		"background_color":     payload.BackgroundColor,
		"surface_color":        payload.SurfaceColor,
		"text_color":           payload.TextColor,
		"muted_text_color":     payload.MutedTextColor,
		"border_color":         payload.BorderColor,
		"primary_contrast_color": payload.PrimaryContrastColor,
		"accent_contrast_color":  payload.AccentContrastColor,
		"surface_contrast_color": payload.SurfaceContrastColor,
		"success_color":          payload.SuccessColor,
		"warning_color":          payload.WarningColor,
		"danger_color":           payload.DangerColor,
		"info_color":             payload.InfoColor,
		"label_color":            payload.LabelColor,
		"label_background_color": payload.LabelBackgroundColor,
		"font_family":            payload.FontFamily,
		"heading_font_family":    payload.HeadingFontFamily,
		"border_radius":          payload.BorderRadius,
		"button_radius":          payload.ButtonRadius,
		"shadow_style":           payload.ShadowStyle,
		"metadata":             payload.Metadata,
		"updated_at":           time.Now(),
	}

	result := r.db.
		WithContext(ctx).
		Model(&websiteSettingDomain.WebsiteSetting{}).
		Where("id = ?", existing.ID).
		Updates(updates)

	if result.Error != nil {
		span.RecordError(result.Error)
		span.SetStatus(codes.Error, "failed update website setting")
		return result.Error
	}
	if result.RowsAffected != 1 {
		err := errors.New("website setting not found")
		span.RecordError(err)
		span.SetStatus(codes.Error, "website setting not found")
		return err
	}

	payload.ID = existing.ID
	span.SetStatus(codes.Ok, "success update website setting")
	return nil
}
