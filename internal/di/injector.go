//go:build wireinject
// +build wireinject

package di

import (
	"net/http"

	"github.com/google/wire"
	"github.com/julienschmidt/httprouter"
	"gorm.io/gorm"

	"github.com/bagusaditiasetiawan/saetechnology-be/internal/config"
	httpDelivery "github.com/bagusaditiasetiawan/saetechnology-be/internal/delivery/http"
	"github.com/bagusaditiasetiawan/saetechnology-be/internal/delivery/http/handler"
	recoverMiddleware "github.com/bagusaditiasetiawan/saetechnology-be/internal/delivery/http/middleware/recover"
	"github.com/bagusaditiasetiawan/saetechnology-be/internal/delivery/http/response"
	brokerDomain "github.com/bagusaditiasetiawan/saetechnology-be/internal/domain/broker"
	"github.com/bagusaditiasetiawan/saetechnology-be/internal/infrastructure/broker"
	"github.com/bagusaditiasetiawan/saetechnology-be/internal/infrastructure/cache"
	"github.com/bagusaditiasetiawan/saetechnology-be/internal/infrastructure/database"
	"github.com/bagusaditiasetiawan/saetechnology-be/internal/infrastructure/email"
	"github.com/bagusaditiasetiawan/saetechnology-be/internal/infrastructure/storage"
	"github.com/bagusaditiasetiawan/saetechnology-be/internal/pkg/hash"
	"github.com/bagusaditiasetiawan/saetechnology-be/internal/pkg/jwt"
	"github.com/bagusaditiasetiawan/saetechnology-be/internal/pkg/logger"
	"github.com/bagusaditiasetiawan/saetechnology-be/internal/pkg/tracing"
	"github.com/bagusaditiasetiawan/saetechnology-be/internal/pkg/validator"
	"github.com/bagusaditiasetiawan/saetechnology-be/internal/repository/contact"
	"github.com/bagusaditiasetiawan/saetechnology-be/internal/repository/content"
	"github.com/bagusaditiasetiawan/saetechnology-be/internal/repository/product"
	"github.com/bagusaditiasetiawan/saetechnology-be/internal/repository/user"
	"github.com/bagusaditiasetiawan/saetechnology-be/internal/repository/website_setting"
	"github.com/bagusaditiasetiawan/saetechnology-be/internal/usecase/auth"
	contactUsecase "github.com/bagusaditiasetiawan/saetechnology-be/internal/usecase/contact"
	contentUsecase "github.com/bagusaditiasetiawan/saetechnology-be/internal/usecase/content"
	"github.com/bagusaditiasetiawan/saetechnology-be/internal/usecase/email_register"
	"github.com/bagusaditiasetiawan/saetechnology-be/internal/usecase/publish_register"
	"github.com/bagusaditiasetiawan/saetechnology-be/internal/usecase/upload"
	productUsecase "github.com/bagusaditiasetiawan/saetechnology-be/internal/usecase/product"
	websiteSettingUsecase "github.com/bagusaditiasetiawan/saetechnology-be/internal/usecase/website_setting"
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

var productHandlerSet = wire.NewSet(
	product.NewPostgresqlRepository,
	productUsecase.NewUseCase,
	handler.NewProductHandler,
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
		productHandlerSet,
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
