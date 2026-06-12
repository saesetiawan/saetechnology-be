package auth

import (
	"context"
	"fmt"
	"saetechnology-be/internal/config"
	"saetechnology-be/internal/constant"
	"saetechnology-be/internal/delivery/http/exception"
	"saetechnology-be/internal/domain/user"
	hash2 "saetechnology-be/internal/pkg/hash"
	"saetechnology-be/internal/pkg/jwt"
	"saetechnology-be/internal/pkg/logger"
	"saetechnology-be/internal/usecase/publish_register"

	gojwt "github.com/golang-jwt/jwt/v5"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type UserLoginPayload struct {
	Email     string `json:"email"`
	FullName  string `json:"full_name"`
	Phone     string `json:"phone"`
	AvatarUrl string `json:"avatar_url"`
	Role      string `json:"role"`
}

type LoginResponse struct {
	Type         string           `json:"type"`
	User         UserLoginPayload `json:"user"`
	AccessToken  string           `json:"accessToken"`
	RefreshToken string           `json:"refreshToken"`
}

type ProfileResponse struct {
	ID        string `json:"id"`
	FullName  string `json:"full_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	AvatarURL string `json:"avatar_url"`
	Role      string `json:"role"`
	Status    string `json:"status"`
}

type Usecase interface {
	Login(ctx context.Context, email string, password string, allowedRoles ...string) LoginResponse
	RefreshToken(ctx context.Context, refreshToken string) LoginResponse
	Logout(ctx context.Context, userID string)
	ActivateAccount(ctx context.Context, activationToken string)
	Register(ctx context.Context, payload user.CreateUserDto) string
	GetProfile(ctx context.Context, userID string) ProfileResponse
	UpdateProfile(ctx context.Context, userID string, payload user.UpdateProfileDto) ProfileResponse
	ChangePassword(ctx context.Context, userID string, payload user.ChangePasswordDto)
}

type usecaseImpl struct {
	userRepository  user.Repository
	secret          string
	hasher          hash2.Hasher
	jwtService      jwt.JWT
	logger          logger.Logger
	publishRegister publish_register.UseCase
	Trace           trace.Tracer
}

func NewUseCase(
	publisher publish_register.UseCase,
	logger logger.Logger,
	jwtService jwt.JWT,
	userRepository user.Repository,
	config *config.Config,
	hasher hash2.Hasher,
	traceProvider trace.TracerProvider,
) Usecase {
	return &usecaseImpl{
		jwtService:      jwtService,
		userRepository:  userRepository,
		secret:          config.Secret,
		hasher:          hasher,
		logger:          logger,
		publishRegister: publisher,
		Trace:           traceProvider.Tracer("AuthUsecase"),
	}
}

func (u *usecaseImpl) Login(ctx context.Context, email string, password string, allowedRoles ...string) LoginResponse {
	ctx, span := u.Trace.Start(ctx, "AuthUsecase.Login")
	defer span.End()

	roleScope := fmt.Sprintf("%v", allowedRoles)
	span.SetAttributes(
		attribute.String("auth.email", email),
		attribute.String("auth.allowed_roles", roleScope),
	)

	u.logger.Info(fmt.Sprintf("process login check email existing %s", email))

	userFounded, err := u.userRepository.FindByEmail(ctx, email)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "email is not found")

		u.logger.Error("email is not found", logger.Field{
			Key:   "email",
			Value: email,
		})

		panic(exception.NewUnauthorized("invalid credentials"))
	}

	span.SetAttributes(
		attribute.String("user.id", userFounded.ID.String()),
		attribute.String("user.email", userFounded.Email),
		attribute.String("user.role", userFounded.Role),
	)

	if len(allowedRoles) > 0 && !isAllowedRole(userFounded.Role, allowedRoles) {
		err := fmt.Errorf("invalid role")
		span.RecordError(err)
		span.SetStatus(codes.Error, "invalid role")

		panic(exception.NewUnauthorized("invalid credentials"))
	}

	if userFounded.Status != constant.StatusActive {
		err := fmt.Errorf("account is not active")
		span.RecordError(err)
		span.SetStatus(codes.Error, "account is not active")

		panic(exception.NewUnauthorized("account is not active, please activate your email"))
	}

	err = u.hasher.Compare(userFounded.Password, password)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "invalid password")

		panic(exception.NewUnauthorized("invalid credentials"))
	}

	accessToken, err := u.jwtService.GenerateAccessToken(userFounded.ID, userFounded.Role)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed generate jwt token")

		panic(exception.NewInternalServiceException("internal service error"))
	}

	refreshToken, err := u.jwtService.GenerateRefreshToken(userFounded.ID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed generate refresh token")

		panic(exception.NewInternalServiceException("internal service error"))
	}

	loginResponse := LoginResponse{
		Type:         constant.BearerToken,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: UserLoginPayload{
			Email:     email,
			FullName:  userFounded.FullName,
			AvatarUrl: userFounded.AvatarURL,
			Phone:     userFounded.Phone,
			Role:      userFounded.Role,
		},
	}

	span.SetStatus(codes.Ok, "login success")

	return loginResponse
}

func (u *usecaseImpl) RefreshToken(ctx context.Context, refreshToken string) LoginResponse {
	ctx, span := u.Trace.Start(ctx, "AuthUsecase.RefreshToken")
	defer span.End()

	token, err := u.jwtService.Verify(refreshToken)
	if err != nil || token == nil || !token.Valid {
		if err != nil {
			span.RecordError(err)
		}
		span.SetStatus(codes.Error, "invalid refresh token")

		panic(exception.NewUnauthorized("invalid refresh token"))
	}

	claims, ok := token.Claims.(gojwt.MapClaims)
	if !ok {
		span.SetStatus(codes.Error, "invalid refresh token claims")

		panic(exception.NewUnauthorized("invalid refresh token"))
	}

	tokenType, ok := claims["token_type"].(string)
	if !ok || tokenType != constant.RefreshToken {
		span.SetStatus(codes.Error, "invalid refresh token type")

		panic(exception.NewUnauthorized("invalid refresh token"))
	}

	userID, ok := claims["id"].(string)
	if !ok || userID == "" {
		span.SetStatus(codes.Error, "invalid refresh token user id")

		panic(exception.NewUnauthorized("invalid refresh token"))
	}

	userFounded, err := u.userRepository.FindByID(ctx, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "user is not found")

		panic(exception.NewUnauthorized("invalid refresh token"))
	}

	accessToken, err := u.jwtService.GenerateAccessToken(userFounded.ID, userFounded.Role)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed generate jwt token")

		panic(exception.NewInternalServiceException("internal service error"))
	}

	newRefreshToken, err := u.jwtService.GenerateRefreshToken(userFounded.ID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed generate refresh token")

		panic(exception.NewInternalServiceException("internal service error"))
	}

	span.SetStatus(codes.Ok, "refresh token success")

	return LoginResponse{
		Type:         constant.BearerToken,
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		User: UserLoginPayload{
			Email:     userFounded.Email,
			FullName:  userFounded.FullName,
			AvatarUrl: userFounded.AvatarURL,
			Phone:     userFounded.Phone,
			Role:      userFounded.Role,
		},
	}
}

func (u *usecaseImpl) Logout(ctx context.Context, userID string) {
	ctx, span := u.Trace.Start(ctx, "AuthUsecase.Logout")
	defer span.End()

	span.SetAttributes(attribute.String("user.id", userID))
	span.SetStatus(codes.Ok, "logout success")
}

func (u *usecaseImpl) GetProfile(ctx context.Context, userID string) ProfileResponse {
	ctx, span := u.Trace.Start(ctx, "AuthUsecase.GetProfile")
	defer span.End()

	userFounded, err := u.userRepository.FindByID(ctx, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "user is not found")
		panic(exception.NewUnauthorized("user not found"))
	}

	span.SetStatus(codes.Ok, "profile found")
	return toProfileResponse(userFounded)
}

func (u *usecaseImpl) UpdateProfile(
	ctx context.Context,
	userID string,
	payload user.UpdateProfileDto,
) ProfileResponse {
	ctx, span := u.Trace.Start(ctx, "AuthUsecase.UpdateProfile")
	defer span.End()

	currentUser, err := u.userRepository.FindByID(ctx, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "user is not found")
		panic(exception.NewUnauthorized("user not found"))
	}

	if currentUser.Role != constant.RoleCustomer {
		err := fmt.Errorf("only customer can update customer profile")
		span.RecordError(err)
		span.SetStatus(codes.Error, "invalid role")
		panic(exception.NewUnauthorized(err.Error()))
	}

	updatedUser, err := u.userRepository.UpdateProfile(ctx, userID, payload)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed update profile")
		panic(exception.NewBadRequestException(err.Error()))
	}

	span.SetStatus(codes.Ok, "profile updated")
	return toProfileResponse(updatedUser)
}

func (u *usecaseImpl) ChangePassword(
	ctx context.Context,
	userID string,
	payload user.ChangePasswordDto,
) {
	ctx, span := u.Trace.Start(ctx, "AuthUsecase.ChangePassword")
	defer span.End()

	if payload.NewPassword != payload.NewPasswordConfirmation {
		err := fmt.Errorf("new password confirmation does not match")
		span.RecordError(err)
		span.SetStatus(codes.Error, "password confirmation mismatch")
		panic(exception.NewBadRequestException(err.Error()))
	}

	currentUser, err := u.userRepository.FindByID(ctx, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "user is not found")
		panic(exception.NewUnauthorized("user not found"))
	}

	if currentUser.Role != constant.RoleCustomer {
		err := fmt.Errorf("only customer can change customer password")
		span.RecordError(err)
		span.SetStatus(codes.Error, "invalid role")
		panic(exception.NewUnauthorized(err.Error()))
	}

	if err := u.hasher.Compare(currentUser.Password, payload.CurrentPassword); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "invalid current password")
		panic(exception.NewBadRequestException("password lama tidak sesuai"))
	}

	hashedPassword, err := u.hasher.Hash(payload.NewPassword)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed hash password")
		panic(exception.NewInternalServiceException("internal server error"))
	}

	if err := u.userRepository.UpdatePassword(ctx, userID, hashedPassword); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed update password")
		panic(exception.NewInternalServiceException("failed update password"))
	}

	span.SetStatus(codes.Ok, "password updated")
}

func (u *usecaseImpl) ActivateAccount(ctx context.Context, activationToken string) {
	ctx, span := u.Trace.Start(ctx, "AuthUsecase.ActivateAccount")
	defer span.End()

	token, err := u.jwtService.Verify(activationToken)
	if err != nil || token == nil || !token.Valid {
		if err != nil {
			span.RecordError(err)
		}
		span.SetStatus(codes.Error, "invalid activation token")

		panic(exception.NewUnauthorized("invalid activation token"))
	}

	claims, ok := token.Claims.(gojwt.MapClaims)
	if !ok {
		span.SetStatus(codes.Error, "invalid activation token claims")

		panic(exception.NewUnauthorized("invalid activation token"))
	}

	tokenType, ok := claims["token_type"].(string)
	if !ok || tokenType != constant.ActivationToken {
		span.SetStatus(codes.Error, "invalid activation token type")

		panic(exception.NewUnauthorized("invalid activation token"))
	}

	userID, ok := claims["id"].(string)
	if !ok || userID == "" {
		span.SetStatus(codes.Error, "invalid activation token user id")

		panic(exception.NewUnauthorized("invalid activation token"))
	}

	userFounded, err := u.userRepository.FindByID(ctx, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "user is not found")

		panic(exception.NewUnauthorized("invalid activation token"))
	}

	if userFounded.Status == constant.StatusActive {
		span.SetStatus(codes.Ok, "account already active")
		return
	}

	if err := u.userRepository.Activate(ctx, userID); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed activate account")

		panic(exception.NewInternalServiceException("failed activate account"))
	}

	span.SetStatus(codes.Ok, "account activated")
}

func isAllowedRole(role string, allowedRoles []string) bool {
	for _, allowedRole := range allowedRoles {
		if role == allowedRole {
			return true
		}
	}

	return false
}

func toProfileResponse(userData *user.User) ProfileResponse {
	return ProfileResponse{
		ID:        userData.ID.String(),
		FullName:  userData.FullName,
		Email:     userData.Email,
		Phone:     userData.Phone,
		AvatarURL: userData.AvatarURL,
		Role:      userData.Role,
		Status:    userData.Status,
	}
}

func (u *usecaseImpl) Register(ctx context.Context, payload user.CreateUserDto) string {
	ctx, span := u.Trace.Start(ctx, "AuthUsecase.Register")
	defer span.End()

	span.SetAttributes(
		attribute.String("user.email", payload.Email),
		attribute.String("user.full_name", payload.FullName),
		attribute.String("user.phone", payload.Phone),
		attribute.String("user.role", payload.Role),
	)

	existingUser, err := u.userRepository.FindByEmail(ctx, payload.Email)
	if err == nil && existingUser != nil {
		err := fmt.Errorf("email already registered")
		span.RecordError(err)
		span.SetStatus(codes.Error, "email already registered")

		u.logger.Error(fmt.Sprintf("email already registered %s", existingUser.Email))

		panic(exception.NewBadRequestException("user with this email already exists"))
	}

	hashedPassword, err := u.hasher.Hash(payload.Password)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed hash password")

		u.logger.Error(err.Error())

		panic(exception.NewInternalServiceException("internal server error"))
	}

	status := constant.StatusPending
	if payload.Role == constant.RoleAdmin {
		status = constant.StatusActive
	}

	dataUser := user.User{
		Status:   status,
		Email:    payload.Email,
		Phone:    payload.Phone,
		FullName: payload.FullName,
		Password: hashedPassword,
		Role:     payload.Role,
	}

	err = u.userRepository.Create(ctx, &dataUser)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed create user")

		u.logger.Error(err.Error())

		panic(exception.NewInternalServiceException("internal server error"))
	}

	span.SetAttributes(
		attribute.String("created_user.id", dataUser.ID.String()),
		attribute.String("created_user.email", dataUser.Email),
	)

	if dataUser.Role != constant.RoleAdmin {
		activationToken, err := u.jwtService.GenerateActivationToken(dataUser.ID)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, "failed generate activation token")

			panic(exception.NewInternalServiceException("internal server error"))
		}

		u.publishRegister.SendToQueue(ctx, dataUser, activationToken)
	}

	span.SetStatus(codes.Ok, "register success")

	return dataUser.Email
}
