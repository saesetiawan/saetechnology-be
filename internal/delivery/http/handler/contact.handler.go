package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"saetechnology-be/internal/delivery/http/exception"
	"saetechnology-be/internal/delivery/http/response"
	contactDomain "saetechnology-be/internal/domain/contact"
	"saetechnology-be/internal/pkg/validator"
	contactUsecase "saetechnology-be/internal/usecase/contact"

	"github.com/julienschmidt/httprouter"
	"go.opentelemetry.io/otel/trace"
)

type ContactHandler interface {
	CreateCaptcha(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	Create(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	FindAll(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	UpdateStatus(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
}

type ContactHandlerImpl struct {
	usecase         contactUsecase.UseCase
	successResponse response.SuccessResponse
	validator       validator.Validator
	tracer          trace.Tracer
}

func NewContactHandler(
	validator validator.Validator,
	usecase contactUsecase.UseCase,
	successResponse response.SuccessResponse,
	tracerProvider trace.TracerProvider,
) ContactHandler {
	return &ContactHandlerImpl{
		validator:       validator,
		usecase:         usecase,
		successResponse: successResponse,
		tracer:          tracerProvider.Tracer("ContactHandler"),
	}
}

func (h *ContactHandlerImpl) CreateCaptcha(
	w http.ResponseWriter,
	r *http.Request,
	_ httprouter.Params,
) {
	ctx, span := h.tracer.Start(r.Context(), "ContactHandler.CreateCaptcha")
	defer span.End()

	result, err := h.usecase.CreateCaptcha(ctx)
	if err != nil {
		span.RecordError(err)
		panic(exception.NewInternalServiceException(err.Error()))
	}

	h.successResponse.Send(ctx, w, http.StatusOK, result)
}

func (h *ContactHandlerImpl) Create(
	w http.ResponseWriter,
	r *http.Request,
	_ httprouter.Params,
) {
	ctx, span := h.tracer.Start(r.Context(), "ContactHandler.Create")
	defer span.End()

	var payload contactDomain.CreateContactMessageDto
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		span.RecordError(err)
		panic(exception.NewBadRequestException("invalid request body"))
	}

	if errorMessage := h.validator.Validate(payload); errorMessage != "" {
		span.RecordError(errors.New(errorMessage))
		panic(exception.NewBadRequestException(errorMessage))
	}

	result, err := h.usecase.Create(ctx, payload)
	if err != nil {
		span.RecordError(err)
		panic(exception.NewBadRequestException(err.Error()))
	}

	h.successResponse.Send(ctx, w, http.StatusCreated, result)
}

func (h *ContactHandlerImpl) FindAll(
	w http.ResponseWriter,
	r *http.Request,
	_ httprouter.Params,
) {
	requireAdmin(r)

	ctx, span := h.tracer.Start(r.Context(), "ContactHandler.FindAll")
	defer span.End()

	query := parseContactQuery(r)
	result, total, err := h.usecase.FindAll(ctx, query)
	if err != nil {
		span.RecordError(err)
		panic(exception.NewInternalServiceException(err.Error()))
	}

	h.successResponse.Send(ctx, w, http.StatusOK, map[string]interface{}{
		"data":  result,
		"total": total,
		"page":  query.Page,
		"limit": query.Limit,
	})
}

func (h *ContactHandlerImpl) UpdateStatus(
	w http.ResponseWriter,
	r *http.Request,
	ps httprouter.Params,
) {
	requireAdmin(r)

	ctx, span := h.tracer.Start(r.Context(), "ContactHandler.UpdateStatus")
	defer span.End()

	var payload contactDomain.UpdateContactStatusDto
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		span.RecordError(err)
		panic(exception.NewBadRequestException("invalid request body"))
	}

	payload.ID = ps.ByName("id")
	if payload.ID == "" {
		panic(exception.NewBadRequestException("id is required"))
	}

	if errorMessage := h.validator.Validate(payload); errorMessage != "" {
		span.RecordError(errors.New(errorMessage))
		panic(exception.NewBadRequestException(errorMessage))
	}

	if err := h.usecase.UpdateStatus(ctx, payload); err != nil {
		span.RecordError(err)
		panic(exception.NewBadRequestException(err.Error()))
	}

	h.successResponse.Send(ctx, w, http.StatusOK, map[string]interface{}{
		"message": "contact status updated successfully",
	})
}

func parseContactQuery(r *http.Request) contactDomain.ListContactMessageQuery {
	query := r.URL.Query()
	page, _ := strconv.Atoi(query.Get("page"))
	limit, _ := strconv.Atoi(query.Get("limit"))

	return contactDomain.ListContactMessageQuery{
		Page:      page,
		Limit:     limit,
		Search:    query.Get("search"),
		Status:    query.Get("status"),
		OrderBy:   query.Get("order_by"),
		OrderType: query.Get("order_type"),
	}
}
