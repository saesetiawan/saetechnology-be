package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"saetechnology-be/internal/delivery/http/exception"
	"saetechnology-be/internal/delivery/http/response"
	websiteSettingDomain "saetechnology-be/internal/domain/website_setting"
	"saetechnology-be/internal/pkg/validator"
	websiteSettingUsecase "saetechnology-be/internal/usecase/website_setting"

	"github.com/julienschmidt/httprouter"
	"go.opentelemetry.io/otel/trace"
)

type WebsiteSettingHandler interface {
	Find(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	FindPublic(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	Update(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
}

type WebsiteSettingHandlerImpl struct {
	usecase         websiteSettingUsecase.UseCase
	successResponse response.SuccessResponse
	validator       validator.Validator
	tracer          trace.Tracer
}

func NewWebsiteSettingHandler(
	validator validator.Validator,
	usecase websiteSettingUsecase.UseCase,
	successResponse response.SuccessResponse,
	tracerProvider trace.TracerProvider,
) WebsiteSettingHandler {
	return &WebsiteSettingHandlerImpl{
		validator:       validator,
		usecase:         usecase,
		successResponse: successResponse,
		tracer:          tracerProvider.Tracer("WebsiteSettingHandler"),
	}
}

func (h *WebsiteSettingHandlerImpl) Find(
	w http.ResponseWriter,
	r *http.Request,
	_ httprouter.Params,
) {
	requireAdmin(r)

	ctx, span := h.tracer.Start(r.Context(), "WebsiteSettingHandler.Find")
	defer span.End()

	result, err := h.usecase.Find(ctx)
	if err != nil {
		span.RecordError(err)
		panic(exception.NewBadRequestException(err.Error()))
	}

	h.successResponse.Send(ctx, w, http.StatusOK, result)
}

func (h *WebsiteSettingHandlerImpl) FindPublic(
	w http.ResponseWriter,
	r *http.Request,
	_ httprouter.Params,
) {
	ctx, span := h.tracer.Start(r.Context(), "WebsiteSettingHandler.FindPublic")
	defer span.End()

	result, err := h.usecase.Find(ctx)
	if err != nil {
		span.RecordError(err)
		panic(exception.NewBadRequestException(err.Error()))
	}

	h.successResponse.Send(ctx, w, http.StatusOK, result)
}

func (h *WebsiteSettingHandlerImpl) Update(
	w http.ResponseWriter,
	r *http.Request,
	_ httprouter.Params,
) {
	requireAdmin(r)

	ctx, span := h.tracer.Start(r.Context(), "WebsiteSettingHandler.Update")
	defer span.End()

	var payload websiteSettingDomain.UpdateWebsiteSettingDto
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		span.RecordError(err)
		panic(exception.NewBadRequestException("invalid request body"))
	}

	if errorMessage := h.validator.Validate(payload); errorMessage != "" {
		span.RecordError(errors.New(errorMessage))
		panic(exception.NewBadRequestException(errorMessage))
	}

	result, err := h.usecase.Update(ctx, payload)
	if err != nil {
		span.RecordError(err)
		panic(exception.NewBadRequestException(err.Error()))
	}

	h.successResponse.Send(ctx, w, http.StatusOK, result)
}
