package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/bagusaditiasetiawan/saetechnology-be/internal/delivery/http/exception"
	"github.com/bagusaditiasetiawan/saetechnology-be/internal/delivery/http/response"
	productDomain "github.com/bagusaditiasetiawan/saetechnology-be/internal/domain/product"
	"github.com/bagusaditiasetiawan/saetechnology-be/internal/pkg/validator"
	productUsecase "github.com/bagusaditiasetiawan/saetechnology-be/internal/usecase/product"

	"github.com/julienschmidt/httprouter"
	"go.opentelemetry.io/otel/trace"
)

type ProductHandler interface {
	Create(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	FindAll(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	FindPublic(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	FindByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	FindPublicBySlug(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	Update(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
}

type ProductHandlerImpl struct {
	usecase         productUsecase.UseCase
	successResponse response.SuccessResponse
	validator       validator.Validator
	tracer          trace.Tracer
}

func NewProductHandler(
	validator validator.Validator,
	usecase productUsecase.UseCase,
	successResponse response.SuccessResponse,
	tracerProvider trace.TracerProvider,
) ProductHandler {
	return &ProductHandlerImpl{
		validator:       validator,
		usecase:         usecase,
		successResponse: successResponse,
		tracer:          tracerProvider.Tracer("ProductHandler"),
	}
}

func (h *ProductHandlerImpl) Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	requireAdmin(r)

	ctx, span := h.tracer.Start(r.Context(), "ProductHandler.Create")
	defer span.End()

	var payload productDomain.CreateProductDto
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

func (h *ProductHandlerImpl) FindAll(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	requireAdmin(r)
	h.findList(w, r, false)
}

func (h *ProductHandlerImpl) FindPublic(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	h.findList(w, r, true)
}

func (h *ProductHandlerImpl) FindByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	requireAdmin(r)

	ctx, span := h.tracer.Start(r.Context(), "ProductHandler.FindByID")
	defer span.End()

	result, err := h.usecase.FindByID(ctx, ps.ByName("id"))
	if err != nil {
		span.RecordError(err)
		panic(exception.NewBadRequestException(err.Error()))
	}

	h.successResponse.Send(ctx, w, http.StatusOK, result)
}

func (h *ProductHandlerImpl) FindPublicBySlug(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx, span := h.tracer.Start(r.Context(), "ProductHandler.FindPublicBySlug")
	defer span.End()

	result, err := h.usecase.FindBySlug(ctx, ps.ByName("slug"), true)
	if err != nil {
		span.RecordError(err)
		panic(exception.NewBadRequestException(err.Error()))
	}

	h.successResponse.Send(ctx, w, http.StatusOK, result)
}

func (h *ProductHandlerImpl) Update(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	requireAdmin(r)

	ctx, span := h.tracer.Start(r.Context(), "ProductHandler.Update")
	defer span.End()

	var payload productDomain.UpdateProductDto
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

func (h *ProductHandlerImpl) Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	requireAdmin(r)

	ctx, span := h.tracer.Start(r.Context(), "ProductHandler.Delete")
	defer span.End()

	if err := h.usecase.Delete(ctx, ps.ByName("id")); err != nil {
		span.RecordError(err)
		panic(exception.NewBadRequestException(err.Error()))
	}

	h.successResponse.Send(ctx, w, http.StatusOK, map[string]interface{}{
		"message": "product deleted successfully",
	})
}

func (h *ProductHandlerImpl) findList(w http.ResponseWriter, r *http.Request, public bool) {
	ctx, span := h.tracer.Start(r.Context(), "ProductHandler.FindList")
	defer span.End()

	query := parseProductQuery(r)
	query.PublicOnly = public
	if public {
		query.Status = "published"
	}

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

func parseProductQuery(r *http.Request) productDomain.ListProductQuery {
	query := r.URL.Query()
	page, _ := strconv.Atoi(query.Get("page"))
	limit, _ := strconv.Atoi(query.Get("limit"))

	var isFeatured *bool
	if value := query.Get("is_featured"); value != "" {
		parsed, _ := strconv.ParseBool(value)
		isFeatured = &parsed
	}

	return productDomain.ListProductQuery{
		Page:       page,
		Limit:      limit,
		Search:     query.Get("search"),
		Status:     query.Get("status"),
		Category:   query.Get("category"),
		IsFeatured: isFeatured,
		OrderBy:    query.Get("order_by"),
		OrderType:  query.Get("order_type"),
	}
}
