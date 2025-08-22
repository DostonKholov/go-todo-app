package middleware

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"go.mood/pkg"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type ctxKey string

const (
	ctxKeyUserID   ctxKey = "user_id"
	ctxKeyUserRole ctxKey = "user_role"
)

// AuthMiddleware — проверяет Authorization: Bearer <token>
// кладёт user_id и role в context.
func AuthMiddleware(next http.Handler) http.Handler {
	secret := os.Getenv("JWT_SECRET")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if secret == "" {
			pkg.WriteJSONResponse(w, http.StatusInternalServerError, errors.New("server misconfigured: JWT_SECRET not set"))
			return
		}

		auth := r.Header.Get("Authorization")
		if auth == "" {
			pkg.WriteJSONResponse(w, http.StatusUnauthorized, errors.New("требуется авторизация"))
			return
		}

		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			pkg.WriteJSONResponse(w, http.StatusUnauthorized, errors.New("неправильный заголовок Authorization"))
			return
		}
		tokenStr := parts[1]

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			// можно проверить метод подписи
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			pkg.WriteJSONResponse(w, http.StatusUnauthorized, errors.New("неверный или просроченный токен"))
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			pkg.WriteJSONResponse(w, http.StatusUnauthorized, errors.New("неверные claims"))
			return
		}

		// извлечь user_id (возможно float64) и role
		var userID int64
		switch v := claims["user_id"].(type) {
		case float64:
			userID = int64(v)
		case int64:
			userID = v
		case string:
			id, _ := strconv.ParseInt(v, 10, 64)
			userID = id
		default:
			pkg.WriteJSONResponse(w, http.StatusUnauthorized, errors.New("неверный user_id в токене"))
			return
		}

		role, _ := claims["role"].(string)

		// запишем в context
		ctx := context.WithValue(r.Context(), ctxKeyUserID, userID)
		ctx = context.WithValue(ctx, ctxKeyUserRole, role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireRole возвращает middleware, который разрешает доступ только роли requiredRole (например "admin")
func RequireRole(requiredRole string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, err := GetUserRole(r.Context())
			if err != nil {
				pkg.WriteJSONResponse(w, http.StatusUnauthorized, errors.New("требуется авторизация"))
				return
			}
			if role != requiredRole {
				pkg.WriteJSONResponse(w, http.StatusForbidden, errors.New("доступ запрещён"))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// GetUserID — helper: доставать user id из context
func GetUserID(ctx context.Context) (int64, error) {
	v := ctx.Value(ctxKeyUserID)
	if v == nil {
		return 0, errors.New("user id not found in context")
	}
	id, ok := v.(int64)
	if !ok {
		return 0, errors.New("invalid user id type in context")
	}
	return id, nil
}

// GetUserRole — helper: доставать роль из context
func GetUserRole(ctx context.Context) (string, error) {
	v := ctx.Value(ctxKeyUserRole)
	if v == nil {
		return "", errors.New("user role not found in context")
	}
	role, ok := v.(string)
	if !ok {
		return "", errors.New("invalid role type in context")
	}
	return role, nil
}
