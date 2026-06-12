package logger

import (
	"github.com/julienschmidt/httprouter"
	"saetechnology-be/internal/pkg/logger"
	"net/http"
	"time"
)

func Middleware(log logger.Logger, next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		start := time.Now()

		next(w, r, ps)

		duration := time.Since(start)

		log.Info("http_request",
			logger.Field{Key: "method", Value: r.Method},
			logger.Field{Key: "path", Value: r.URL.Path},
			logger.Field{Key: "duration_ms", Value: duration.Milliseconds()},
			logger.Field{Key: "remote_ip", Value: r.RemoteAddr},
		)
	}
}
