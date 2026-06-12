package contact

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	cacheDomain "github.com/bagusaditiasetiawan/saetechnology-be/internal/domain/cache"
	contactDomain "github.com/bagusaditiasetiawan/saetechnology-be/internal/domain/contact"
	"github.com/bagusaditiasetiawan/saetechnology-be/internal/domain/jsonvalue"
	"github.com/bagusaditiasetiawan/saetechnology-be/internal/helpers"
	"github.com/bagusaditiasetiawan/saetechnology-be/internal/pkg/logger"

	"github.com/google/uuid"
)

const captchaTTL = 5 * time.Minute

type UseCase interface {
	CreateCaptcha(ctx context.Context) (contactDomain.ContactCaptcha, error)
	Create(ctx context.Context, payload contactDomain.CreateContactMessageDto) (*contactDomain.ContactMessage, error)
	FindAll(ctx context.Context, query contactDomain.ListContactMessageQuery) ([]contactDomain.ContactMessage, int64, error)
	UpdateStatus(ctx context.Context, payload contactDomain.UpdateContactStatusDto) error
}

type useCaseImpl struct {
	logger     logger.Logger
	repository contactDomain.Repository
	cache      cacheDomain.Repository
}

func NewUseCase(
	logger logger.Logger,
	repository contactDomain.Repository,
	cacheRepository cacheDomain.Repository,
) UseCase {
	return &useCaseImpl{
		logger:     logger,
		repository: repository,
		cache:      cacheRepository,
	}
}

func (uc *useCaseImpl) CreateCaptcha(ctx context.Context) (contactDomain.ContactCaptcha, error) {
	left := rand.Intn(8) + 2
	right := rand.Intn(8) + 2
	id, err := uuid.NewV7()
	if err != nil {
		return contactDomain.ContactCaptcha{}, err
	}

	key := captchaCacheKey(id.String())
	if err := uc.cache.Set(ctx, key, fmt.Sprintf("%d", left+right), captchaTTL); err != nil {
		return contactDomain.ContactCaptcha{}, err
	}

	return contactDomain.ContactCaptcha{
		ID:       id.String(),
		Question: fmt.Sprintf("%d + %d", left, right),
	}, nil
}

func (uc *useCaseImpl) Create(
	ctx context.Context,
	payload contactDomain.CreateContactMessageDto,
) (*contactDomain.ContactMessage, error) {
	uc.logger.Info("process create contact message")

	if strings.TrimSpace(payload.Website) != "" {
		return nil, errors.New("invalid contact submission")
	}
	if !uc.verifyCaptcha(ctx, payload.CaptchaID, payload.CaptchaAnswer) {
		return nil, errors.New("invalid captcha answer")
	}

	metadata, err := jsonvalue.New("{}")
	if err != nil {
		return nil, err
	}

	entity := &contactDomain.ContactMessage{
		Name:     strings.TrimSpace(payload.Name),
		Email:    strings.TrimSpace(payload.Email),
		Phone:    strings.TrimSpace(payload.Phone),
		Company:  strings.TrimSpace(payload.Company),
		Subject:  strings.TrimSpace(payload.Subject),
		Message:  strings.TrimSpace(payload.Message),
		Source:   "home_contact",
		Status:   "new",
		Metadata: metadata,
	}

	if err := uc.repository.Create(ctx, entity); err != nil {
		return nil, err
	}

	helpers.DeleteCacheByPrefix(ctx, uc.cache, "contact-messages:")
	return entity, nil
}

func (uc *useCaseImpl) FindAll(
	ctx context.Context,
	query contactDomain.ListContactMessageQuery,
) ([]contactDomain.ContactMessage, int64, error) {
	cacheKey := helpers.PaginationCacheKey(
		"contact-messages:"+query.Status,
		query.Page,
		query.Limit,
		query.Search,
		"",
		query.OrderBy,
		query.OrderType,
	)

	if cached, ok := helpers.GetJSONCache[helpers.PaginatedCache[contactDomain.ContactMessage]](ctx, uc.cache, cacheKey); ok {
		return cached.Data, cached.Total, nil
	}

	result, total, err := uc.repository.FindAll(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	helpers.SetJSONCache(ctx, uc.cache, cacheKey, helpers.PaginatedCache[contactDomain.ContactMessage]{
		Data:  result,
		Total: total,
	}, helpers.DefaultCacheTTL)

	return result, total, nil
}

func (uc *useCaseImpl) UpdateStatus(
	ctx context.Context,
	payload contactDomain.UpdateContactStatusDto,
) error {
	if err := uc.repository.UpdateStatus(ctx, payload.ID, strings.TrimSpace(payload.Status)); err != nil {
		return err
	}

	helpers.DeleteCacheByPrefix(ctx, uc.cache, "contact-messages:")
	return nil
}

func (uc *useCaseImpl) verifyCaptcha(ctx context.Context, id string, answer string) bool {
	key := captchaCacheKey(strings.TrimSpace(id))
	expected, err := uc.cache.Get(ctx, key)
	if err != nil {
		return false
	}
	_ = uc.cache.Delete(ctx, key)

	return strings.TrimSpace(expected) == strings.TrimSpace(answer)
}

func captchaCacheKey(id string) string {
	return "contact-captcha:" + id
}
