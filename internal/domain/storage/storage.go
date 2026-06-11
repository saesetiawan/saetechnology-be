package storage

import "context"

type Storage interface {
	Upload(ctx context.Context, input UploadInput) (*UploadOutput, error)
	Delete(ctx context.Context, key string) error
	GetPresignedUploadURL(ctx context.Context, input PresignedUploadInput) (*PresignedUploadOutput, error)
	GetPublicURL(key string) string
}
