package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go-platform-core/internal/domain/user"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

type postgresqlRepository struct {
	db     *gorm.DB
	tracer trace.Tracer
}

func NewPostgresqlRepository(
	db *gorm.DB,
	tracerProvider trace.TracerProvider,
) user.Repository {
	return &postgresqlRepository{
		db:     db,
		tracer: tracerProvider.Tracer("UserRepository"),
	}
}

func (r *postgresqlRepository) FindByEmail(
	ctx context.Context,
	email string,
) (*user.User, error) {
	ctx, span := r.tracer.Start(ctx, "UserRepository.FindByEmail")
	defer span.End()

	span.SetAttributes(
		attribute.String("user.email", email),
	)

	var result user.User

	err := r.db.
		WithContext(ctx).
		Where("email = ?", email).
		First(&result).
		Error

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed find user by email")
		return nil, err
	}

	span.SetAttributes(
		attribute.String("user.id", result.ID.String()),
		attribute.String("user.status", result.Status),
	)

	span.SetStatus(codes.Ok, "success find user by email")

	return &result, nil
}

func (r *postgresqlRepository) Activate(
	ctx context.Context,
	id string,
) error {
	ctx, span := r.tracer.Start(ctx, "UserRepository.Activate")
	defer span.End()

	span.SetAttributes(
		attribute.String("user.id", id),
	)

	now := time.Now()

	err := r.db.
		WithContext(ctx).
		Model(&user.User{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":            "active",
			"email_verified_at": &now,
			"updated_at":        now,
		}).
		Error

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed activate user")
		return err
	}

	span.SetStatus(codes.Ok, "success activate user")

	return nil
}

func (r *postgresqlRepository) FindByID(
	ctx context.Context,
	id string,
) (*user.User, error) {
	ctx, span := r.tracer.Start(ctx, "UserRepository.FindByID")
	defer span.End()

	span.SetAttributes(
		attribute.String("user.id", id),
	)

	var result user.User

	err := r.db.
		WithContext(ctx).
		Where("id = ?", id).
		First(&result).
		Error

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed find user by id")
		return nil, err
	}

	span.SetAttributes(
		attribute.String("user.email", result.Email),
		attribute.String("user.status", result.Status),
	)

	span.SetStatus(codes.Ok, "success find user by id")

	return &result, nil
}

func (r *postgresqlRepository) FindAll(
	ctx context.Context,
	query user.ListUsersQuery,
) ([]user.User, int64, error) {
	ctx, span := r.tracer.Start(ctx, "UserRepository.FindAll")
	defer span.End()

	_, limit, offset := normalizeUserPagination(query.Page, query.Limit)
	searchBy := normalizeUserSearchBy(query.SearchBy)
	orderBy := normalizeUserOrderBy(query.OrderBy)
	orderType := normalizeUserOrderType(query.OrderType)

	dbQuery := r.db.
		WithContext(ctx).
		Model(&user.User{}).
		Where("deleted_at IS NULL")

	if query.Search != "" {
		dbQuery = dbQuery.Where(searchBy+" ILIKE ?", "%"+query.Search+"%")
	}
	if query.Role != "" {
		dbQuery = dbQuery.Where("role = ?", query.Role)
	}
	if query.Status != "" {
		dbQuery = dbQuery.Where("status = ?", query.Status)
	}

	var total int64
	if err := dbQuery.Count(&total).Error; err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed count users")
		return nil, 0, err
	}

	var users []user.User
	if err := dbQuery.
		Order(fmt.Sprintf("%s %s", orderBy, orderType)).
		Limit(limit).
		Offset(offset).
		Find(&users).Error; err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed find users")
		return nil, 0, err
	}

	span.SetStatus(codes.Ok, "success find users")
	return users, total, nil
}

func (r *postgresqlRepository) Create(
	ctx context.Context,
	payload *user.User,
) error {
	ctx, span := r.tracer.Start(ctx, "UserRepository.Create")
	defer span.End()

	span.SetAttributes(
		attribute.String("user.email", payload.Email),
		attribute.String("user.full_name", payload.FullName),
		attribute.String("user.status", payload.Status),
	)

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		id, err := uuid.NewV7()
		if err != nil {
			return err
		}

		payload.ID = id

		if err := tx.Create(payload).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed create user")
		return err
	}

	span.SetAttributes(
		attribute.String("created_user.id", payload.ID.String()),
	)

	span.SetStatus(codes.Ok, "success create user")

	return nil
}

func (r *postgresqlRepository) UpdateStatus(
	ctx context.Context,
	id string,
	status string,
) (*user.User, error) {
	ctx, span := r.tracer.Start(ctx, "UserRepository.UpdateStatus")
	defer span.End()

	parsedID, err := uuid.Parse(id)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "invalid user id")
		return nil, err
	}

	var result user.User
	err = r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		updates := map[string]interface{}{
			"status":     status,
			"updated_at": time.Now(),
		}
		if status == "active" {
			now := time.Now()
			updates["email_verified_at"] = &now
		}

		update := tx.
			Model(&user.User{}).
			Where("id = ? AND deleted_at IS NULL", parsedID).
			Updates(updates)
		if update.Error != nil {
			return update.Error
		}
		if update.RowsAffected != 1 {
			return errors.New("user not found")
		}

		return tx.Where("id = ?", parsedID).First(&result).Error
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed update user status")
		return nil, err
	}

	span.SetStatus(codes.Ok, "success update user status")
	return &result, nil
}

func (r *postgresqlRepository) UpdateProfile(
	ctx context.Context,
	id string,
	payload user.UpdateProfileDto,
) (*user.User, error) {
	ctx, span := r.tracer.Start(ctx, "UserRepository.UpdateProfile")
	defer span.End()

	parsedID, err := uuid.Parse(id)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "invalid user id")
		return nil, err
	}

	var result user.User
	err = r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		update := tx.
			Model(&user.User{}).
			Where("id = ?", parsedID).
			Updates(map[string]interface{}{
				"full_name":  payload.FullName,
				"email":      payload.Email,
				"phone":      payload.Phone,
				"updated_at": time.Now(),
			})
		if update.Error != nil {
			return update.Error
		}
		if update.RowsAffected != 1 {
			return errors.New("user not found")
		}

		return tx.Where("id = ?", parsedID).First(&result).Error
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed update profile")
		return nil, err
	}

	span.SetStatus(codes.Ok, "success update profile")
	return &result, nil
}

func (r *postgresqlRepository) UpdatePassword(
	ctx context.Context,
	id string,
	passwordHash string,
) error {
	ctx, span := r.tracer.Start(ctx, "UserRepository.UpdatePassword")
	defer span.End()

	parsedID, err := uuid.Parse(id)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "invalid user id")
		return err
	}

	update := r.db.
		WithContext(ctx).
		Model(&user.User{}).
		Where("id = ?", parsedID).
		Updates(map[string]interface{}{
			"password":   passwordHash,
			"updated_at": time.Now(),
		})
	if update.Error != nil {
		span.RecordError(update.Error)
		span.SetStatus(codes.Error, "failed update password")
		return update.Error
	}
	if update.RowsAffected != 1 {
		err := errors.New("user not found")
		span.RecordError(err)
		span.SetStatus(codes.Error, "user not found")
		return err
	}

	span.SetStatus(codes.Ok, "success update password")
	return nil
}

func normalizeUserPagination(page int, limit int) (int, int, int) {
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

func normalizeUserSearchBy(value string) string {
	allowed := map[string]bool{
		"full_name": true,
		"email":     true,
		"phone":     true,
		"role":      true,
		"status":    true,
	}
	if !allowed[value] {
		return "full_name"
	}

	return value
}

func normalizeUserOrderBy(value string) string {
	allowed := map[string]bool{
		"full_name":  true,
		"email":      true,
		"role":       true,
		"status":     true,
		"created_at": true,
		"updated_at": true,
	}
	if !allowed[value] {
		return "created_at"
	}

	return value
}

func normalizeUserOrderType(value string) string {
	if value == "asc" || value == "ASC" {
		return "ASC"
	}

	return "DESC"
}
