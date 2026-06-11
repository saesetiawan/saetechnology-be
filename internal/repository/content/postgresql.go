package content

import (
	"context"
	"errors"
	"fmt"
	"time"

	contentDomain "go-platform-core/internal/domain/content"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
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
) contentDomain.Repository {
	return &PostgresqlRepository{
		db:     db,
		tracer: tracerProvider.Tracer("WebsiteContentRepository"),
	}
}

func (r *PostgresqlRepository) Create(
	ctx context.Context,
	payload *contentDomain.WebsiteContent,
) error {
	ctx, span := r.tracer.Start(ctx, "WebsiteContentRepository.Create")
	defer span.End()

	id, err := uuid.NewV7()
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed generate content id")
		return err
	}

	payload.ID = id

	span.SetAttributes(
		attribute.String("content.key", payload.Key),
		attribute.String("content.type", payload.Type),
		attribute.String("content.placement", payload.Placement),
	)

	if err := r.db.WithContext(ctx).Create(payload).Error; err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed create website content")
		return err
	}

	span.SetStatus(codes.Ok, "success create website content")
	return nil
}

func (r *PostgresqlRepository) FindByID(
	ctx context.Context,
	id string,
) (*contentDomain.WebsiteContent, error) {
	ctx, span := r.tracer.Start(ctx, "WebsiteContentRepository.FindByID")
	defer span.End()

	var result contentDomain.WebsiteContent
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&result).Error; err != nil {
		span.RecordError(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			span.SetStatus(codes.Error, "website content not found")
			return nil, errors.New("website content not found")
		}

		span.SetStatus(codes.Error, "failed find website content by id")
		return nil, err
	}

	span.SetStatus(codes.Ok, "success find website content by id")
	return &result, nil
}

func (r *PostgresqlRepository) FindByKey(
	ctx context.Context,
	key string,
) (*contentDomain.WebsiteContent, error) {
	ctx, span := r.tracer.Start(ctx, "WebsiteContentRepository.FindByKey")
	defer span.End()

	var result contentDomain.WebsiteContent
	if err := r.db.WithContext(ctx).Where("key = ?", key).First(&result).Error; err != nil {
		span.RecordError(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			span.SetStatus(codes.Error, "website content not found")
			return nil, errors.New("website content not found")
		}

		span.SetStatus(codes.Error, "failed find website content by key")
		return nil, err
	}

	span.SetStatus(codes.Ok, "success find website content by key")
	return &result, nil
}

func (r *PostgresqlRepository) FindAll(
	ctx context.Context,
	query contentDomain.ListWebsiteContentQuery,
) ([]contentDomain.WebsiteContent, int64, error) {
	ctx, span := r.tracer.Start(ctx, "WebsiteContentRepository.FindAll")
	defer span.End()

	page, limit, offset := normalizePagination(query.Page, query.Limit)
	searchBy := normalizeSearchBy(query.SearchBy)
	orderBy := normalizeOrderBy(query.OrderBy)
	orderType := normalizeOrderType(query.OrderType)

	span.SetAttributes(
		attribute.Int("pagination.page", page),
		attribute.Int("pagination.limit", limit),
		attribute.String("content.placement", query.Placement),
		attribute.String("content.type", query.Type),
	)

	dbQuery := r.db.
		WithContext(ctx).
		Model(&contentDomain.WebsiteContent{})

	if query.Search != "" {
		dbQuery = dbQuery.Where(searchBy+" ILIKE ?", "%"+query.Search+"%")
	}
	if query.Type != "" {
		dbQuery = dbQuery.Where("type = ?", query.Type)
	}
	if query.Placement != "" {
		dbQuery = dbQuery.Where("placement = ?", query.Placement)
	}
	if query.Active != nil {
		now := time.Now()
		dbQuery = dbQuery.
			Where("is_active = ?", *query.Active)
		if *query.Active {
			dbQuery = dbQuery.
				Where("(publish_start_at IS NULL OR publish_start_at <= ?)", now).
				Where("(publish_end_at IS NULL OR publish_end_at >= ?)", now)
		}
	}

	var total int64
	if err := dbQuery.Count(&total).Error; err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed count website contents")
		return nil, 0, err
	}

	var contents []contentDomain.WebsiteContent
	if err := dbQuery.
		Order(fmt.Sprintf("%s %s", orderBy, orderType)).
		Limit(limit).
		Offset(offset).
		Find(&contents).Error; err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed find website contents")
		return nil, 0, err
	}

	span.SetAttributes(
		attribute.Int("content.result_count", len(contents)),
		attribute.Int64("content.total", total),
	)
	span.SetStatus(codes.Ok, "success find website contents")

	return contents, total, nil
}

func (r *PostgresqlRepository) Update(
	ctx context.Context,
	payload *contentDomain.WebsiteContent,
) error {
	ctx, span := r.tracer.Start(ctx, "WebsiteContentRepository.Update")
	defer span.End()

	updates := map[string]interface{}{
		"key":              payload.Key,
		"type":             payload.Type,
		"placement":        payload.Placement,
		"title":            payload.Title,
		"subtitle":         payload.Subtitle,
		"body":             payload.Body,
		"image_url":        payload.ImageURL,
		"link_url":         payload.LinkURL,
		"link_label":       payload.LinkLabel,
		"sort_order":       payload.SortOrder,
		"is_active":        payload.IsActive,
		"metadata":         payload.Metadata,
		"publish_start_at": payload.PublishStartAt,
		"publish_end_at":   payload.PublishEndAt,
		"updated_at":       time.Now(),
	}

	result := r.db.
		WithContext(ctx).
		Model(&contentDomain.WebsiteContent{}).
		Where("id = ?", payload.ID).
		Updates(updates)

	if result.Error != nil {
		span.RecordError(result.Error)
		span.SetStatus(codes.Error, "failed update website content")
		return result.Error
	}
	if result.RowsAffected != 1 {
		err := errors.New("website content not found")
		span.RecordError(err)
		span.SetStatus(codes.Error, "website content not found")
		return err
	}

	span.SetStatus(codes.Ok, "success update website content")
	return nil
}

func normalizePagination(page int, limit int) (int, int, int) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	return page, limit, (page - 1) * limit
}

func normalizeSearchBy(value string) string {
	allowed := map[string]bool{
		"key":       true,
		"title":     true,
		"type":      true,
		"placement": true,
	}
	if !allowed[value] {
		return "title"
	}

	return value
}

func normalizeOrderBy(value string) string {
	allowed := map[string]bool{
		"key":        true,
		"title":      true,
		"type":       true,
		"placement":  true,
		"sort_order": true,
		"created_at": true,
		"updated_at": true,
	}
	if !allowed[value] {
		return "sort_order"
	}

	return value
}

func normalizeOrderType(value string) string {
	if value == "desc" || value == "DESC" {
		return "DESC"
	}

	return "ASC"
}
