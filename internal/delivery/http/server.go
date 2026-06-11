package http

import (
	"fmt"
	"go-platform-core/internal/config"
	corsMiddleware "go-platform-core/internal/delivery/http/middleware/cors"
	middlewareRecover "go-platform-core/internal/delivery/http/middleware/recover"
	"go-platform-core/internal/delivery/http/middleware/tracing"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

func NewServer(
	recoverMiddleware *middlewareRecover.Middleware,
	config *config.Config,
	router *httprouter.Router,
	tracerProvider trace.TracerProvider,
) *http.Server {
	handler := tracing.NewMiddleware(router, tracerProvider)
	recoverMiddleware.Handler = handler
	return &http.Server{
		Addr:         fmt.Sprintf(":%s", config.AppPort),
		Handler:      corsMiddleware.NewMiddleware(recoverMiddleware),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}
