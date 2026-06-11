package auth

import (
	"context"
	"net/http"
	"strings"

	"go-platform-core/internal/delivery/http/exception"
	jwtPkg "go-platform-core/internal/pkg/jwt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/julienschmidt/httprouter"
)

type contextKey string

const UserIDContextKey contextKey = "user_id"
const UserRoleContextKey contextKey = "user_role"

func AuthGuard(
	jwtService jwtPkg.JWT,
	next httprouter.Handle,
) httprouter.Handle {
	return func(
		w http.ResponseWriter,
		r *http.Request,
		ps httprouter.Params,
	) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			panic(exception.NewUnauthorized("missing authorization header"))
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" || parts[1] == "" {
			panic(exception.NewUnauthorized("invalid authorization header"))
		}

		token, err := jwtService.Verify(parts[1])
		if err != nil || token == nil || !token.Valid {
			panic(exception.NewUnauthorized("invalid token"))
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			panic(exception.NewUnauthorized("invalid token claims"))
		}

		tokenType, ok := claims["token_type"].(string)
		if !ok || tokenType != "access" {
			panic(exception.NewUnauthorized("invalid token type"))
		}

		userID, ok := claims["id"].(string)
		if !ok || userID == "" {
			panic(exception.NewUnauthorized("invalid user id"))
		}

		userRole, _ := claims["role"].(string)

		if userRole == "" {
			panic(exception.NewUnauthorized("invalid token role"))
		}

		ctx := context.WithValue(
			r.Context(),
			UserIDContextKey,
			userID,
		)
		ctx = context.WithValue(ctx, UserRoleContextKey, userRole)

		next(w, r.WithContext(ctx), ps)
	}
}
