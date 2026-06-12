package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"saetechnology-be/internal/delivery/http/exception"
	"saetechnology-be/internal/delivery/http/response"
	contentDomain "saetechnology-be/internal/domain/content"
	"saetechnology-be/internal/pkg/validator"
	contentUsecase "saetechnology-be/internal/usecase/content"

	"github.com/julienschmidt/httprouter"
	"go.opentelemetry.io/otel/trace"
)

type ContentHandler interface {
	Create(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	FindByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	FindByKey(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	FindAll(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	FindPublic(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	Update(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
}

type ContentHandlerImpl struct {
	usecase         contentUsecase.UseCase
	successResponse response.SuccessResponse
	validator       validator.Validator
	tracer          trace.Tracer
}

func NewContentHandler(
	validator validator.Validator,
	usecase contentUsecase.UseCase,
	successResponse response.SuccessResponse,
	tracerProvider trace.TracerProvider,
) ContentHandler {
	return &ContentHandlerImpl{
		validator:       validator,
		usecase:         usecase,
		successResponse: successResponse,
		tracer:          tracerProvider.Tracer("ContentHandler"),
	}
}

func (h *ContentHandlerImpl) Create(
	w http.ResponseWriter,
	r *http.Request,
	_ httprouter.Params,
) {
	requireAdmin(r)

	ctx, span := h.tracer.Start(r.Context(), "ContentHandler.Create")
	defer span.End()

	var payload contentDomain.CreateWebsiteContentDto
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

func (h *ContentHandlerImpl) FindByID(
	w http.ResponseWriter,
	r *http.Request,
	ps httprouter.Params,
) {
	requireAdmin(r)

	ctx, span := h.tracer.Start(r.Context(), "ContentHandler.FindByID")
	defer span.End()

	result, err := h.usecase.FindByID(ctx, ps.ByName("id"))
	if err != nil {
		span.RecordError(err)
		panic(exception.NewBadRequestException(err.Error()))
	}

	h.successResponse.Send(ctx, w, http.StatusOK, result)
}

func (h *ContentHandlerImpl) FindByKey(
	w http.ResponseWriter,
	r *http.Request,
	ps httprouter.Params,
) {
	ctx, span := h.tracer.Start(r.Context(), "ContentHandler.FindByKey")
	defer span.End()

	result, err := h.usecase.FindByKey(ctx, ps.ByName("key"))
	if err != nil {
		span.RecordError(err)
		panic(exception.NewBadRequestException(err.Error()))
	}

	h.successResponse.Send(ctx, w, http.StatusOK, result)
}

func (h *ContentHandlerImpl) FindAll(
	w http.ResponseWriter,
	r *http.Request,
	_ httprouter.Params,
) {
	requireAdmin(r)

	ctx, span := h.tracer.Start(r.Context(), "ContentHandler.FindAll")
	defer span.End()

	query := parseContentQuery(r, false)
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

func (h *ContentHandlerImpl) FindPublic(
	w http.ResponseWriter,
	r *http.Request,
	_ httprouter.Params,
) {
	ctx, span := h.tracer.Start(r.Context(), "ContentHandler.FindPublic")
	defer span.End()

	active := true
	query := parseContentQuery(r, true)
	query.Active = &active
	query.OrderBy = "sort_order"
	query.OrderType = "asc"

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

func (h *ContentHandlerImpl) Update(
	w http.ResponseWriter,
	r *http.Request,
	ps httprouter.Params,
) {
	requireAdmin(r)

	ctx, span := h.tracer.Start(r.Context(), "ContentHandler.Update")
	defer span.End()

	var payload contentDomain.UpdateWebsiteContentDto
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

	result, err := h.usecase.Update(ctx, payload)
	if err != nil {
		span.RecordError(err)
		panic(exception.NewBadRequestException(err.Error()))
	}

	h.successResponse.Send(ctx, w, http.StatusOK, result)
}

func (h *ContentHandlerImpl) Delete(
	w http.ResponseWriter,
	r *http.Request,
	ps httprouter.Params,
) {
	requireAdmin(r)

	ctx, span := h.tracer.Start(r.Context(), "ContentHandler.Delete")
	defer span.End()

	id := ps.ByName("id")
	if id == "" {
		panic(exception.NewBadRequestException("id is required"))
	}

	if err := h.usecase.Delete(ctx, id); err != nil {
		span.RecordError(err)
		panic(exception.NewBadRequestException(err.Error()))
	}

	h.successResponse.Send(ctx, w, http.StatusOK, map[string]interface{}{
		"message": "website content deleted successfully",
	})
}

func parseContentQuery(
	r *http.Request,
	public bool,
) contentDomain.ListWebsiteContentQuery {
	query := r.URL.Query()
	page, _ := strconv.Atoi(query.Get("page"))
	limit, _ := strconv.Atoi(query.Get("limit"))

	if public && limit <= 0 {
		limit = 20
	}

	result := contentDomain.ListWebsiteContentQuery{
		Page:      page,
		Limit:     limit,
		Search:    query.Get("search"),
		SearchBy:  query.Get("search_by"),
		OrderBy:   query.Get("order_by"),
		OrderType: query.Get("order_type"),
		Type:      query.Get("type"),
		Placement: query.Get("placement"),
	}

	if activeValue := query.Get("is_active"); activeValue != "" {
		active, _ := strconv.ParseBool(activeValue)
		result.Active = &active
	}

	return result
}
