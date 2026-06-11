package response

import (
	"encoding/json"
	"net/http"
)

type jsonResponder struct{}

func NewJSONResponder() Responder {
	return &jsonResponder{}
}

func (r *jsonResponder) JSON(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(data)
	return err
}
