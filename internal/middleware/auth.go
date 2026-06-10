package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/internships-backend/test-backend-bibihoya-1/pkg/jwt"
)

type contextKey string

const (
	UserIDKey contextKey = "user_id"
	RoleKey   contextKey = "role"
)

func AuthMiddleware(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				sendAuthError(w, "UNAUTHORIZED", "missing authorization", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				sendAuthError(w, "UNAUTHORIZED", "invalid authorization header", http.StatusUnauthorized)
				return
			}

			claims, err := jwt.ParseToken(parts[1], jwtSecret)
			if err != nil {
				sendAuthError(w, "UNAUTHORIZED", "invalid token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserId)
			ctx = context.WithValue(ctx, RoleKey, claims.Role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func AdminCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, ok := r.Context().Value(RoleKey).(string)
		if !ok || role != "admin" {
			sendAuthError(w, "FORBIDDEN", "admin role required", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func UserCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, ok := r.Context().Value(RoleKey).(string)
		if !ok || role != "user" {
			sendAuthError(w, "FORBIDDEN", "user role required", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func sendAuthError(w http.ResponseWriter, code, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": map[string]string{
			"code":    code,
			"message": message,
		},
	})
}

func GetUserID(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDKey).(string)
	return userID, ok
}

func GetRole(ctx context.Context) (string, bool) {
	role, ok := ctx.Value(RoleKey).(string)
	return role, ok
}
