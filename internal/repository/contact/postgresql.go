package contact

import (
	"context"
	"errors"
	"fmt"
	"time"

	contactDomain "saetechnology-be/internal/domain/contact"

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
) contactDomain.Repository {
	return &PostgresqlRepository{
		db:     db,
		tracer: tracerProvider.Tracer("ContactRepository"),
	}
}

func (r *PostgresqlRepository) Create(
	ctx context.Context,
	payload *contactDomain.ContactMessage,
) error {
	ctx, span := r.tracer.Start(ctx, "ContactRepository.Create")
	defer span.End()

	id, err := uuid.NewV7()
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed generate contact id")
		return err
	}

	payload.ID = id

	span.SetAttributes(
		attribute.String("contact.email", payload.Email),
		attribute.String("contact.status", payload.Status),
	)

	if err := r.db.WithContext(ctx).Create(payload).Error; err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed create contact message")
		return err
	}

	span.SetStatus(codes.Ok, "success create contact message")
	return nil
}

func (r *PostgresqlRepository) FindAll(
	ctx context.Context,
	query contactDomain.ListContactMessageQuery,
) ([]contactDomain.ContactMessage, int64, error) {
	ctx, span := r.tracer.Start(ctx, "ContactRepository.FindAll")
	defer span.End()

	_, limit, offset := normalizePagination(query.Page, query.Limit)
	orderBy := normalizeOrderBy(query.OrderBy)
	orderType := normalizeOrderType(query.OrderType)

	dbQuery := r.db.
		WithContext(ctx).
		Model(&contactDomain.ContactMessage{})

	if query.Search != "" {
		search := "%" + query.Search + "%"
		dbQuery = dbQuery.Where(
			"name ILIKE ? OR email ILIKE ? OR company ILIKE ? OR subject ILIKE ?",
			search,
			search,
			search,
			search,
		)
	}
	if query.Status != "" {
		dbQuery = dbQuery.Where("status = ?", query.Status)
	}

	var total int64
	if err := dbQuery.Count(&total).Error; err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed count contact messages")
		return nil, 0, err
	}

	var messages []contactDomain.ContactMessage
	if err := dbQuery.
		Order(fmt.Sprintf("%s %s", orderBy, orderType)).
		Limit(limit).
		Offset(offset).
		Find(&messages).Error; err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed find contact messages")
		return nil, 0, err
	}

	span.SetAttributes(attribute.Int64("contact.total", total))
	span.SetStatus(codes.Ok, "success find contact messages")
	return messages, total, nil
}

func (r *PostgresqlRepository) UpdateStatus(
	ctx context.Context,
	id string,
	status string,
) error {
	ctx, span := r.tracer.Start(ctx, "ContactRepository.UpdateStatus")
	defer span.End()

	result := r.db.
		WithContext(ctx).
		Model(&contactDomain.ContactMessage{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": time.Now(),
		})

	if result.Error != nil {
		span.RecordError(result.Error)
		span.SetStatus(codes.Error, "failed update contact status")
		return result.Error
	}
	if result.RowsAffected != 1 {
		err := errors.New("contact message not found")
		span.RecordError(err)
		span.SetStatus(codes.Error, "contact message not found")
		return err
	}

	span.SetStatus(codes.Ok, "success update contact status")
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
		"name":       true,
		"email":      true,
		"subject":    true,
		"status":     true,
		"created_at": true,
		"updated_at": true,
	}
	if !allowed[value] {
		return "created_at"
	}

	return value
}

func normalizeOrderType(value string) string {
	if value == "asc" || value == "ASC" {
		return "ASC"
	}

	return "DESC"
}
