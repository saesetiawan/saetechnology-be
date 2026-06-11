package recover

import (
	"go-platform-core/internal/delivery/http/exception"
	"go-platform-core/internal/delivery/http/response"
	"go-platform-core/internal/pkg/logger"
	"net/http"
)

type Middleware struct {
	Handler   http.Handler
	Responder response.Responder
	Logger    logger.Logger
}

func NewMiddleware(logger logger.Logger, handler http.Handler, responder response.Responder) *Middleware {
	return &Middleware{
		Handler:   handler,
		Responder: responder,
		Logger:    logger,
	}
}

func (resp *Middleware) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	defer func() {
		if rec := recover(); rec != nil {
			resp.Logger.Error(rec.(error).Error())
			switch err := rec.(type) {
			case exception.BadRequestException:
				resp.Responder.JSON(writer, http.StatusBadRequest, map[string]string{"message": err.Error()})
				return
			case exception.Unauthorized:
				resp.Responder.JSON(writer, http.StatusUnauthorized, map[string]string{"message": err.Error()})
				return
			default:
				resp.Responder.JSON(writer, http.StatusInternalServerError, map[string]string{"message": "internal server error"})
				return
			}
		}
	}()
	resp.Handler.ServeHTTP(writer, request)
}
