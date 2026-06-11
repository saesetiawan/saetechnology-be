//go:build wireinject
// +build wireinject

package di

import (
	"net/http"

	"github.com/google/wire"
	"github.com/julienschmidt/httprouter"
	"gorm.io/gorm"

	"go-platform-core/internal/config"
	httpDelivery "go-platform-core/internal/delivery/http"
	"go-platform-core/internal/delivery/http/handler"
	recoverMiddleware "go-platform-core/internal/delivery/http/middleware/recover"
	"go-platform-core/internal/delivery/http/response"
	brokerDomain "go-platform-core/internal/domain/broker"
	"go-platform-core/internal/infrastructure/broker"
	"go-platform-core/internal/infrastructure/cache"
	"go-platform-core/internal/infrastructure/database"
	"go-platform-core/internal/infrastructure/email"
	"go-platform-core/internal/infrastructure/storage"
	"go-platform-core/internal/pkg/hash"
	"go-platform-core/internal/pkg/jwt"
	"go-platform-core/internal/pkg/logger"
	"go-platform-core/internal/pkg/tracing"
	"go-platform-core/internal/pkg/validator"
	"go-platform-core/internal/repository/content"
	"go-platform-core/internal/repository/user"
	"go-platform-core/internal/repository/website_setting"
	"go-platform-core/internal/usecase/auth"
	contentUsecase "go-platform-core/internal/usecase/content"
	"go-platform-core/internal/usecase/email_register"
	"go-platform-core/internal/usecase/publish_register"
	"go-platform-core/internal/usecase/upload"
	websiteSettingUsecase "go-platform-core/internal/usecase/website_setting"
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
