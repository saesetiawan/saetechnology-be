package upload

import (
	"context"
	"strings"

	"go-platform-core/internal/domain/storage"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type FileUploadUseCase interface {
	Execute(ctx context.Context, input storage.UploadInput) (*storage.UploadOutput, error)
}

type FileUploadUseCaseImpl struct {
	Storage storage.Storage
	Trace   trace.Tracer
}

func NewFileUploadUseCase(
	storage storage.Storage,
	traceProvider trace.TracerProvider,
) FileUploadUseCase {
	return &FileUploadUseCaseImpl{
		Storage: storage,
		Trace:   traceProvider.Tracer("FileUploadUseCase"),
	}
}

func (uc *FileUploadUseCaseImpl) Execute(
	ctx context.Context,
	input storage.UploadInput,
) (*storage.UploadOutput, error) {
	ctx, span := uc.Trace.Start(ctx, "FileUploadUseCase.Execute")
	defer span.End()

	span.SetAttributes(
		attribute.String("file.name", strings.ToValidUTF8(input.Filename, "_")),
		attribute.String("file.content_type", strings.ToValidUTF8(input.ContentType, "_")),
	)

	result, err := uc.Storage.Upload(ctx, input)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed upload file")
		return nil, err
	}

	if result != nil {
		span.SetAttributes(
			attribute.String("file.url", strings.ToValidUTF8(result.URL, "_")),
		)
	}

	span.SetStatus(codes.Ok, "success upload file")

	return result, nil
}
