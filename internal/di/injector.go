//go:build wireinject
// +build wireinject

package di

import (
	"net/http"

	"github.com/google/wire"
	"github.com/julienschmidt/httprouter"
	"gorm.io/gorm"

	"saetechnology-be/internal/config"
	httpDelivery "saetechnology-be/internal/delivery/http"
	"saetechnology-be/internal/delivery/http/handler"
	recoverMiddleware "saetechnology-be/internal/delivery/http/middleware/recover"
	"saetechnology-be/internal/delivery/http/response"
	brokerDomain "saetechnology-be/internal/domain/broker"
	"saetechnology-be/internal/infrastructure/broker"
	"saetechnology-be/internal/infrastructure/cache"
	"saetechnology-be/internal/infrastructure/database"
	"saetechnology-be/internal/infrastructure/email"
	"saetechnology-be/internal/infrastructure/storage"
	"saetechnology-be/internal/pkg/hash"
	"saetechnology-be/internal/pkg/jwt"
	"saetechnology-be/internal/pkg/logger"
	"saetechnology-be/internal/pkg/tracing"
	"saetechnology-be/internal/pkg/validator"
	"saetechnology-be/internal/repository/contact"
	"saetechnology-be/internal/repository/content"
	"saetechnology-be/internal/repository/user"
	"saetechnology-be/internal/repository/website_setting"
	"saetechnology-be/internal/usecase/auth"
	contactUsecase "saetechnology-be/internal/usecase/contact"
	contentUsecase "saetechnology-be/internal/usecase/content"
	"saetechnology-be/internal/usecase/email_register"
	"saetechnology-be/internal/usecase/publish_register"
	"saetechnology-be/internal/usecase/upload"
	websiteSettingUsecase "saetechnology-be/internal/usecase/website_setting"
)

var loggerSet = wire.NewSet(logger.NewLogrus)

var cacheSet = wire.NewSet(
	config.LoadRedisConfig,
	cache.NewRedisClient,
	cache.NewRedisRepository,
)

var registerEmailSet = wire.NewSet(
	broker.NewRabbitConnection,
	broker.NewRabbitPublisher,
	publish_register.NewUseCase,
)

var authHandlerSet = wire.NewSet(
	user.NewPostgresqlRepository,
	hash.NewArgon2Hasher,
	registerEmailSet,
	auth.NewUseCase,
	handler.NewAuthHandler,
)

var contentHandlerSet = wire.NewSet(
	content.NewPostgresqlRepository,
	contentUsecase.NewUseCase,
	handler.NewContentHandler,
)

var contactHandlerSet = wire.NewSet(
	contact.NewPostgresqlRepository,
	contactUsecase.NewUseCase,
	handler.NewContactHandler,
)

var websiteSettingHandlerSet = wire.NewSet(
	website_setting.NewPostgresqlRepository,
	websiteSettingUsecase.NewUseCase,
	handler.NewWebsiteSettingHandler,
)

var storageHandlerSet = wire.NewSet(
	config.LoadS3Config,
	storage.NewS3Storage,
	upload.NewFileUploadUseCase,
	handler.NewStorageHandlerImpl,
)

var httpSet = wire.NewSet(
	response.NewJSONResponder,
	response.NewSuccessResponse,
	handler.NewHealthHandler,
	httpDelivery.NewRouter,
	recoverMiddleware.NewMiddleware,
	wire.Bind(new(http.Handler), new(*httprouter.Router)),
	httpDelivery.NewServer,
)

func InitServer() *http.Server {
	wire.Build(
		tracing.NewTracerProvider,
		config.Load,
		validator.NewValidator,
		jwt.NewJWT,
		database.NewPostgresql,
		loggerSet,
		cacheSet,
		authHandlerSet,
		contentHandlerSet,
		contactHandlerSet,
		websiteSettingHandlerSet,
		storageHandlerSet,
		httpSet,
	)

	return nil
}

func InitPostgresql() (*gorm.DB, error) {
	wire.Build(config.Load, database.NewPostgresql)
	return nil, nil
}

func InitConsumer() brokerDomain.Consumer {
	wire.Build(
		tracing.NewTracerProvider,
		config.Load,
		broker.NewRabbitConnection,
		broker.NewRabbitConsumer,
	)

	return nil
}

func InitRegisterEmailConsumer() *broker.UserConsumer {
	wire.Build(
		loggerSet,
		config.Load,
		email.NewEmailSender,
		email_register.NewUseCase,
		broker.NewUserConsumer,
	)

	return nil
}
