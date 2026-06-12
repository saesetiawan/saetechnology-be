package upload

import (
	"context"

	"saetechnology-be/internal/domain/storage"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type GetFileUrlUseCase interface {
	Execute(ctx context.Context, key string) string
}

type GetFileUrlUseCaseImpl struct {
	Storage storage.Storage
	Trace   trace.Tracer
}

func NewGetFileUrlUseCase(
	storage storage.Storage,
	traceProvider trace.TracerProvider,
) GetFileUrlUseCase {
	return &GetFileUrlUseCaseImpl{
		Storage: storage,
		Trace:   traceProvider.Tracer("GetFileUrlUseCase"),
	}
}

func (uc *GetFileUrlUseCaseImpl) Execute(
	ctx context.Context,
	key string,
) string {
	ctx, span := uc.Trace.Start(ctx, "GetFileUrlUseCase.Execute")
	defer span.End()

	span.SetAttributes(
		attribute.String("file.key", key),
	)

	url := uc.Storage.GetPublicURL(key)

	span.SetAttributes(
		attribute.String("file.url", url),
	)

	span.SetStatus(codes.Ok, "success get file url")

	return url
}
