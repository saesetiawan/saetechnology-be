package website_setting

import "context"

type Repository interface {
	Find(ctx context.Context) (*WebsiteSetting, error)
	Upsert(ctx context.Context, payload *WebsiteSetting) error
}
