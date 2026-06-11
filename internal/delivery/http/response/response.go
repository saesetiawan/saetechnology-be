package response

import "net/http"

type Responder interface {
	JSON(w http.ResponseWriter, status int, data interface{}) error
}
