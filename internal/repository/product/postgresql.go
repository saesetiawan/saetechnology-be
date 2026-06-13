package product

import (
	"context"
	"errors"
	"fmt"
	"time"

	productDomain "github.com/bagusaditiasetiawan/saetechnology-be/internal/domain/product"

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
) productDomain.Repository {
	return &PostgresqlRepository{
		db:     db,
		tracer: tracerProvider.Tracer("ProductRepository"),
	}
}

func (r *PostgresqlRepository) Create(ctx context.Context, payload *productDomain.Product) error {
	ctx, span := r.tracer.Start(ctx, "ProductRepository.Create")
	defer span.End()

	id, err := uuid.NewV7()
	if err != nil {
		span.RecordError(err)
		return err
	}
	payload.ID = id

	if err := r.db.WithContext(ctx).Create(payload).Error; err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed create product")
		return err
	}

	span.SetStatus(codes.Ok, "success create product")
	return nil
}

func (r *PostgresqlRepository) FindByID(ctx context.Context, id string) (*productDomain.Product, error) {
	ctx, span := r.tracer.Start(ctx, "ProductRepository.FindByID")
	defer span.End()

	var result productDomain.Product
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&result).Error; err != nil {
		span.RecordError(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	return &result, nil
}

func (r *PostgresqlRepository) FindBySlug(
	ctx context.Context,
	slug string,
	publicOnly bool,
) (*productDomain.Product, error) {
	ctx, span := r.tracer.Start(ctx, "ProductRepository.FindBySlug")
	defer span.End()

	query := r.db.WithContext(ctx).Where("slug = ?", slug)
	if publicOnly {
		query = query.Where("status = ?", "published")
	}

	var result productDomain.Product
	if err := query.First(&result).Error; err != nil {
		span.RecordError(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	return &result, nil
}

func (r *PostgresqlRepository) FindAll(
	ctx context.Context,
	query productDomain.ListProductQuery,
) ([]productDomain.Product, int64, error) {
	ctx, span := r.tracer.Start(ctx, "ProductRepository.FindAll")
	defer span.End()

	_, limit, offset := normalizePagination(query.Page, query.Limit)
	orderBy := normalizeOrderBy(query.OrderBy)
	orderType := normalizeOrderType(query.OrderType)

	dbQuery := r.db.WithContext(ctx).Model(&productDomain.Product{})
	if query.Search != "" {
		search := "%" + query.Search + "%"
		dbQuery = dbQuery.Where("name ILIKE ? OR slug ILIKE ? OR tagline ILIKE ?", search, search, search)
	}
	if query.Status != "" {
		dbQuery = dbQuery.Where("status = ?", query.Status)
	}
	if query.Category != "" {
		dbQuery = dbQuery.Where("category = ?", query.Category)
	}
	if query.IsFeatured != nil {
		dbQuery = dbQuery.Where("is_featured = ?", *query.IsFeatured)
	}
	if query.PublicOnly {
		dbQuery = dbQuery.Where("status = ?", "published")
	}

	var total int64
	if err := dbQuery.Count(&total).Error; err != nil {
		span.RecordError(err)
		return nil, 0, err
	}

	var products []productDomain.Product
	if err := dbQuery.
		Order(fmt.Sprintf("%s %s", orderBy, orderType)).
		Limit(limit).
		Offset(offset).
		Find(&products).Error; err != nil {
		span.RecordError(err)
		return nil, 0, err
	}

	return products, total, nil
}

func (r *PostgresqlRepository) Update(ctx context.Context, payload *productDomain.Product) error {
	ctx, span := r.tracer.Start(ctx, "ProductRepository.Update")
	defer span.End()

	result := r.db.
		WithContext(ctx).
		Model(&productDomain.Product{}).
		Where("id = ?", payload.ID).
		Updates(map[string]interface{}{
			"slug":        payload.Slug,
			"name":        payload.Name,
			"tagline":     payload.Tagline,
			"summary":     payload.Summary,
			"description": payload.Description,
			"category":    payload.Category,
			"status":      payload.Status,
			"price_label": payload.PriceLabel,
			"price_url":   payload.PriceURL,
			"image_url":   payload.ImageURL,
			"demo_url":    payload.DemoURL,
			"cta_label":   payload.CTALabel,
			"cta_url":     payload.CTAURL,
			"sort_order":  payload.SortOrder,
			"is_featured": payload.IsFeatured,
			"metadata":    payload.Metadata,
			"updated_at":  time.Now(),
		})

	if result.Error != nil {
		span.RecordError(result.Error)
		return result.Error
	}
	if result.RowsAffected != 1 {
		return errors.New("product not found")
	}

	return nil
}

func (r *PostgresqlRepository) Delete(ctx context.Context, id string) error {
	ctx, span := r.tracer.Start(ctx, "ProductRepository.Delete")
	defer span.End()

	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&productDomain.Product{})
	if result.Error != nil {
		span.RecordError(result.Error)
		return result.Error
	}
	if result.RowsAffected != 1 {
		return errors.New("product not found")
	}

	return nil
}

func normalizePagination(page int, limit int) (int, int, int) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	return page, limit, (page - 1) * limit
}

func normalizeOrderBy(value string) string {
	allowed := map[string]bool{
		"slug":       true,
		"name":       true,
		"category":   true,
		"status":     true,
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
