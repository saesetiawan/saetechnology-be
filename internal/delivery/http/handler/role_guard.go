package handler

import (
	"net/http"

	"saetechnology-be/internal/constant"
	"saetechnology-be/internal/delivery/http/exception"
	authMiddleware "saetechnology-be/internal/delivery/http/middleware/auth"
)

func requireAdmin(r *http.Request) {
	role, _ := r.Context().Value(authMiddleware.UserRoleContextKey).(string)
	if role != constant.RoleAdmin {
		panic(exception.NewUnauthorized("admin access is required"))
	}
}
