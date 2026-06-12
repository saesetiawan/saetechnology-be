package contact

import "context"

type Repository interface {
	Create(ctx context.Context, payload *ContactMessage) error
	FindAll(ctx context.Context, query ListContactMessageQuery) ([]ContactMessage, int64, error)
	UpdateStatus(ctx context.Context, id string, status string) error
}
