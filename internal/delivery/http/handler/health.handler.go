package handler

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type HealthHandler interface {
	Liveness(w http.ResponseWriter, r *http.Request, _ httprouter.Params)
}
