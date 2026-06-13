package http

import (
	"github.com/bagusaditiasetiawan/saetechnology-be/internal/delivery/http/handler"
	middleware "github.com/bagusaditiasetiawan/saetechnology-be/internal/delivery/http/middleware/auth"
	loggerMiddleware "github.com/bagusaditiasetiawan/saetechnology-be/internal/delivery/http/middleware/logger"
	"github.com/bagusaditiasetiawan/saetechnology-be/internal/pkg/jwt"
	"github.com/bagusaditiasetiawan/saetechnology-be/internal/pkg/logger"

	"github.com/julienschmidt/httprouter"
)

func NewRouter(
	websiteSettingHandler handler.WebsiteSettingHandler,
	contentHandler handler.ContentHandler,
	contactHandler handler.ContactHandler,
	productHandler handler.ProductHandler,
	jwtService jwt.JWT,
	logger logger.Logger,
	healthHandler handler.HealthHandler,
	authHandler handler.AuthHandler,
	storageHandler handler.StorageHandler,
) *httprouter.Router {
	router := httprouter.New()

	withLogger := func(handle httprouter.Handle) httprouter.Handle {
		return loggerMiddleware.Middleware(logger, handle)
	}

	withAuth := func(handle httprouter.Handle) httprouter.Handle {
		return middleware.AuthGuard(jwtService, withLogger(handle))
	}

	router.GET("/health", withLogger(healthHandler.Liveness))

	router.POST("/api/login", withLogger(authHandler.Login))
	router.POST("/api/register", withLogger(authHandler.Register))
	router.POST("/api/refresh-token", withLogger(authHandler.RefreshToken))
	router.POST("/api/logout", withAuth(authHandler.Logout))
	router.POST("/api/activate-account", withLogger(authHandler.ActivateAccount))
	router.GET("/api/activate-account", withLogger(authHandler.ActivateAccount))
	router.POST("/api/customer/login", withLogger(authHandler.CustomerLogin))
	router.POST("/api/customer/register", withLogger(authHandler.CustomerRegister))
	router.GET("/api/profile", withAuth(authHandler.CustomerProfile))
	router.GET("/api/customer/profile", withAuth(authHandler.CustomerProfile))
	router.PUT("/api/customer/profile", withAuth(authHandler.UpdateCustomerProfile))
	router.PATCH("/api/customer/password", withAuth(authHandler.ChangeCustomerPassword))

	router.POST("/api/storage/upload-file", withAuth(storageHandler.UploadFile))

	router.GET("/api/public/website-contents", withLogger(contentHandler.FindPublic))
	router.POST("/api/website-contents", withAuth(contentHandler.Create))
	router.GET("/api/website-contents", withAuth(contentHandler.FindAll))
	router.GET("/api/website-contents/:id", withAuth(contentHandler.FindByID))
	router.PUT("/api/website-contents/:id", withAuth(contentHandler.Update))
	router.DELETE("/api/website-contents/:id", withAuth(contentHandler.Delete))

	router.GET("/api/public/contact-captcha", withLogger(contactHandler.CreateCaptcha))
	router.POST("/api/public/contact-messages", withLogger(contactHandler.Create))
	router.GET("/api/contact-messages", withAuth(contactHandler.FindAll))
	router.PATCH("/api/contact-messages/:id/status", withAuth(contactHandler.UpdateStatus))

	router.GET("/api/public/products", withLogger(productHandler.FindPublic))
	router.GET("/api/public/products/:slug", withLogger(productHandler.FindPublicBySlug))
	router.POST("/api/products", withAuth(productHandler.Create))
	router.GET("/api/products", withAuth(productHandler.FindAll))
	router.GET("/api/products/:id", withAuth(productHandler.FindByID))
	router.PUT("/api/products/:id", withAuth(productHandler.Update))
	router.DELETE("/api/products/:id", withAuth(productHandler.Delete))

	router.GET("/api/public/website-settings", withLogger(websiteSettingHandler.FindPublic))
	router.GET("/api/website-settings", withAuth(websiteSettingHandler.Find))
	router.PUT("/api/website-settings", withAuth(websiteSettingHandler.Update))

	return router
}
