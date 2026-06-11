package handler

import (
	"encoding/json"
	"errors"
	"go-platform-core/internal/constant"
	"go-platform-core/internal/delivery/http/exception"
	authMiddleware "go-platform-core/internal/delivery/http/middleware/auth"
	"go-platform-core/internal/delivery/http/response"
	"go-platform-core/internal/domain/user"
	"go-platform-core/internal/pkg/validator"
	"go-platform-core/internal/usecase/auth"
	"go.opentelemetry.io/otel/trace"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type AuthHandler interface {
	Login(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	Register(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	RefreshToken(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	Logout(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	ActivateAccount(w http.ResponseWriter, r *http.Request, ps httprouter.Params)

	CustomerLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	CustomerRegister(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	CustomerProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	UpdateCustomerProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	ChangeCustomerPassword(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
}

type authHandler struct {
	usecase   auth.Usecase
	resp      response.Responder
	validator validator.Validator
	Trace     trace.Tracer
}

func NewAuthHandler(validator validator.Validator, u auth.Usecase, resp response.Responder,
	traceProvider trace.TracerProvider) AuthHandler {
	return &authHandler{
		usecase:   u,
		resp:      resp,
		validator: validator,
		Trace:     traceProvider.Tracer("AuthHandler"),
	}
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type refreshTokenRequest struct {
	RefreshToken      string `json:"refreshToken"`
	RefreshTokenSnake string `json:"refresh_token"`
}

type activateAccountRequest struct {
	Token string `json:"token"`
}

func (r refreshTokenRequest) Token() string {
	if r.RefreshToken != "" {
		return r.RefreshToken
	}

	return r.RefreshTokenSnake
}

func (h *authHandler) Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()
	ctx, span := h.Trace.Start(ctx, "AuthHandler.Login")
	defer span.End()
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		span.RecordError(err)
		panic(exception.NewBadRequestException("malformed request"))
	}
	errorMessage := h.validator.Validate(req)
	if errorMessage != "" {
		span.RecordError(errors.New(errorMessage))
		panic(exception.NewBadRequestException(errorMessage))
	}
	responsePayload := h.usecase.Login(ctx, req.Email, req.Password, constant.RoleAdmin)

	h.resp.JSON(w, 200, responsePayload)
}
func (h *authHandler) Register(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()
	ctx, span := h.Trace.Start(ctx, "AuthHandler.Register")
	defer span.End()
	var req user.CreateUserDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		span.RecordError(err)
		panic(exception.NewBadRequestException("malformed request"))
	}
	req.Role = constant.RoleAdmin
	errorMessage := h.validator.Validate(req)
	if errorMessage != "" {
		span.RecordError(errors.New(errorMessage))
		panic(exception.NewBadRequestException(errorMessage))
	}
	h.usecase.Register(ctx, req)
	h.resp.JSON(w, 200, map[string]interface{}{
		"message": "successfully registered user",
	})
}

func (h *authHandler) RefreshToken(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()
	ctx, span := h.Trace.Start(ctx, "AuthHandler.RefreshToken")
	defer span.End()
	var req refreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		span.RecordError(err)
		panic(exception.NewBadRequestException("malformed request"))
	}
	refreshToken := req.Token()
	if refreshToken == "" {
		err := errors.New("RefreshToken is required")
		span.RecordError(err)
		panic(exception.NewBadRequestException(err.Error()))
	}
	responsePayload := h.usecase.RefreshToken(ctx, refreshToken)

	h.resp.JSON(w, 200, responsePayload)
}

func (h *authHandler) Logout(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()
	ctx, span := h.Trace.Start(ctx, "AuthHandler.Logout")
	defer span.End()

	userID, ok := ctx.Value(authMiddleware.UserIDContextKey).(string)
	if !ok || userID == "" {
		err := errors.New("user id is required")
		span.RecordError(err)
		panic(exception.NewUnauthorized(err.Error()))
	}

	h.usecase.Logout(ctx, userID)

	h.resp.JSON(w, 200, map[string]interface{}{
		"message": "successfully logged out",
	})
}

func (h *authHandler) ActivateAccount(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()
	ctx, span := h.Trace.Start(ctx, "AuthHandler.ActivateAccount")
	defer span.End()

	token := r.URL.Query().Get("token")
	if token == "" {
		var req activateAccountRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			span.RecordError(err)
			panic(exception.NewBadRequestException("malformed request"))
		}

		token = req.Token
	}

	if token == "" {
		err := errors.New("token is required")
		span.RecordError(err)
		panic(exception.NewBadRequestException(err.Error()))
	}

	h.usecase.ActivateAccount(ctx, token)

	h.resp.JSON(w, 200, map[string]interface{}{
		"message": "account activated successfully",
	})
}

func (h *authHandler) CustomerLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()
	ctx, span := h.Trace.Start(ctx, "AuthHandler.CustomerLogin")
	defer span.End()
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		span.RecordError(err)
		panic(exception.NewBadRequestException("malformed request"))
	}
	req.Role = constant.RoleCustomer
	errorMessage := h.validator.Validate(req)
	if errorMessage != "" {
		span.RecordError(errors.New(errorMessage))
		panic(exception.NewBadRequestException(errorMessage))
	}
	loginResponse := h.usecase.Login(ctx, req.Email, req.Password, req.Role)

	h.resp.JSON(w, 200, loginResponse)
}
func (h *authHandler) CustomerRegister(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()
	ctx, span := h.Trace.Start(ctx, "AuthHandler.CustomerRegister")
	defer span.End()
	var req user.CreateUserDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		span.RecordError(err)
		panic(exception.NewBadRequestException("malformed request"))
	}
	req.Role = constant.RoleCustomer
	errorMessage := h.validator.Validate(req)
	if errorMessage != "" {
		span.RecordError(errors.New(errorMessage))
		panic(exception.NewBadRequestException(errorMessage))
	}
	h.usecase.Register(ctx, req)
	h.resp.JSON(w, 200, map[string]interface{}{
		"message": "registration successful, please activate your account via email",
	})
}

func (h *authHandler) CustomerProfile(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()
	ctx, span := h.Trace.Start(ctx, "AuthHandler.CustomerProfile")
	defer span.End()

	userID, ok := ctx.Value(authMiddleware.UserIDContextKey).(string)
	if !ok || userID == "" {
		err := errors.New("user id is required")
		span.RecordError(err)
		panic(exception.NewUnauthorized(err.Error()))
	}

	result := h.usecase.GetProfile(ctx, userID)

	h.resp.JSON(w, http.StatusOK, map[string]interface{}{
		"result": result,
	})
}

func (h *authHandler) UpdateCustomerProfile(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()
	ctx, span := h.Trace.Start(ctx, "AuthHandler.UpdateCustomerProfile")
	defer span.End()

	userID, ok := ctx.Value(authMiddleware.UserIDContextKey).(string)
	if !ok || userID == "" {
		err := errors.New("user id is required")
		span.RecordError(err)
		panic(exception.NewUnauthorized(err.Error()))
	}

	var req user.UpdateProfileDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		span.RecordError(err)
		panic(exception.NewBadRequestException("malformed request"))
	}
	if errorMessage := h.validator.Validate(req); errorMessage != "" {
		span.RecordError(errors.New(errorMessage))
		panic(exception.NewBadRequestException(errorMessage))
	}

	result := h.usecase.UpdateProfile(ctx, userID, req)
	h.resp.JSON(w, http.StatusOK, map[string]interface{}{
		"result": result,
	})
}

func (h *authHandler) ChangeCustomerPassword(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()
	ctx, span := h.Trace.Start(ctx, "AuthHandler.ChangeCustomerPassword")
	defer span.End()

	userID, ok := ctx.Value(authMiddleware.UserIDContextKey).(string)
	if !ok || userID == "" {
		err := errors.New("user id is required")
		span.RecordError(err)
		panic(exception.NewUnauthorized(err.Error()))
	}

	var req user.ChangePasswordDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		span.RecordError(err)
		panic(exception.NewBadRequestException("malformed request"))
	}
	if errorMessage := h.validator.Validate(req); errorMessage != "" {
		span.RecordError(errors.New(errorMessage))
		panic(exception.NewBadRequestException(errorMessage))
	}

	h.usecase.ChangePassword(ctx, userID, req)
	h.resp.JSON(w, http.StatusOK, map[string]interface{}{
		"message": "password updated successfully",
	})
}
