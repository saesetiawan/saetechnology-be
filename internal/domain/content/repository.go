package content

import "context"

type Repository interface {
	Create(ctx context.Context, payload *WebsiteContent) error
	FindByID(ctx context.Context, id string) (*WebsiteContent, error)
	FindByKey(ctx context.Context, key string) (*WebsiteContent, error)
	FindAll(ctx context.Context, query ListWebsiteContentQuery) ([]WebsiteContent, int64, error)
	Update(ctx context.Context, payload *WebsiteContent) error
}
