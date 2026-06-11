package handler

import (
	"net/http"

	"go-platform-core/internal/constant"
	"go-platform-core/internal/delivery/http/exception"
	authMiddleware "go-platform-core/internal/delivery/http/middleware/auth"
)

func requireAdmin(r *http.Request) {
	role, _ := r.Context().Value(authMiddleware.UserRoleContextKey).(string)
	if role != constant.RoleAdmin {
		panic(exception.NewUnauthorized("admin access is required"))
	}
}
