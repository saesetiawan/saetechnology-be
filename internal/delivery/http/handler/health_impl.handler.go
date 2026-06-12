// internal/delivery/http/handler/health_impl.go
package handler

import (
	"github.com/julienschmidt/httprouter"
	"saetechnology-be/internal/delivery/http/response"
	"net/http"
)

type healthHandler struct {
	Responder response.Responder
}

func NewHealthHandler(responder response.Responder) HealthHandler {
	return &healthHandler{
		Responder: responder,
	}
}

type HealthResponse struct {
	Message string `json:"message"`
}

func (h *healthHandler) Liveness(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	resp := HealthResponse{
		Message: "app is live",
	}
	h.Responder.JSON(w, http.StatusOK, resp)
}
