package product

import "context"

type Repository interface {
	Create(ctx context.Context, payload *Product) error
	FindByID(ctx context.Context, id string) (*Product, error)
	FindBySlug(ctx context.Context, slug string, publicOnly bool) (*Product, error)
	FindAll(ctx context.Context, query ListProductQuery) ([]Product, int64, error)
	Update(ctx context.Context, payload *Product) error
	Delete(ctx context.Context, id string) error
}
