package cors

import "net/http"

type Middleware struct {
	Handler http.Handler
}

func NewMiddleware(handler http.Handler) *Middleware {
	return &Middleware{
		Handler: handler,
	}
}

func (m *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	m.Handler.ServeHTTP(w, r)
}
