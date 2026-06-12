package handler

import (
	"errors"
	"fmt"
	"github.com/bagusaditiasetiawan/saetechnology-be/internal/constant"
	"github.com/bagusaditiasetiawan/saetechnology-be/internal/delivery/http/exception"
	"github.com/bagusaditiasetiawan/saetechnology-be/internal/domain/storage"
	"io"
	"net/http"
	"strings"

	"go.opentelemetry.io/otel/trace"

	"github.com/bagusaditiasetiawan/saetechnology-be/internal/delivery/http/response"
	"github.com/bagusaditiasetiawan/saetechnology-be/internal/usecase/upload"

	"github.com/julienschmidt/httprouter"
)

type StorageHandler interface {
	UploadFile(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
}

type StorageHandlerImpl struct {
	FileUploadUseCase upload.FileUploadUseCase
	SuccessResponse   response.SuccessResponse
	Tracer            trace.Tracer
}

func NewStorageHandlerImpl(
	file upload.FileUploadUseCase,
	successResponse response.SuccessResponse,
	tracerProvider trace.TracerProvider,
) StorageHandler {
	return &StorageHandlerImpl{
		FileUploadUseCase: file,
		SuccessResponse:   successResponse,
		Tracer:            tracerProvider.Tracer("StorageHandler"),
	}
}

func (s *StorageHandlerImpl) UploadFile(
	w http.ResponseWriter,
	r *http.Request,
	ps httprouter.Params,
) {
	ctx := r.Context()

	ctx, span := s.Tracer.Start(ctx, "StorageHandler.UploadFile")
	defer span.End()

	r.Body = http.MaxBytesReader(w, r.Body, 32<<20)

	contentType := r.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "multipart/form-data") {
		err := fmt.Errorf("invalid content-type: %s", contentType)
		span.RecordError(err)
		panic(exception.NewBadRequestException("request must be multipart/form-data"))
	}

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		span.RecordError(err)

		if errors.Is(err, io.ErrUnexpectedEOF) || strings.Contains(err.Error(), "unexpected EOF") {
			panic(exception.NewBadRequestException("multipart upload body is incomplete"))
		}

		if strings.Contains(err.Error(), "request body too large") {
			panic(exception.NewBadRequestException("file is too large"))
		}

		panic(exception.NewBadRequestException("malformed multipart form"))
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		span.RecordError(err)
		panic(exception.NewBadRequestException("file is required"))
	}
	defer file.Close()

	if header.Size <= 0 {
		err := fmt.Errorf("empty file")
		span.RecordError(err)
		panic(exception.NewBadRequestException("file is empty"))
	}

	body, err := io.ReadAll(file)
	if err != nil {
		span.RecordError(err)
		panic(exception.NewBadRequestException("failed to read file"))
	}

	contentType = header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = http.DetectContentType(body)
	}

	res, err := s.FileUploadUseCase.Execute(ctx, storage.UploadInput{
		Directory:   constant.TemporaryDirectory,
		Body:        body,
		ContentType: contentType,
		Filename:    header.Filename,
	})
	if err != nil {
		span.RecordError(err)
		panic(exception.NewInternalServiceException("file upload error"))
	}

	s.SuccessResponse.Send(ctx, w, http.StatusOK, res)
}
