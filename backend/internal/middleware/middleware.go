package middleware

import (
	"context"
	"net/http"
	"strings"

	"kidtask/internal/auth"
	"kidtask/internal/respond"
)

type contextKey string

const ClaimsKey contextKey = "claims"

func Auth(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if !strings.HasPrefix(header, "Bearer ") {
				respond.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "missing token")
				return
			}
			claims, err := auth.ParseToken(strings.TrimPrefix(header, "Bearer "), secret)
			if err != nil {
				respond.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "invalid token")
				return
			}
			ctx := context.WithValue(r.Context(), ClaimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireRole(role string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := GetClaims(r)
		if claims == nil || claims.Role != role {
			respond.Error(w, http.StatusForbidden, "FORBIDDEN", "forbidden")
			return
		}
		next(w, r)
	}
}

func GetClaims(r *http.Request) *auth.Claims {
	claims, _ := r.Context().Value(ClaimsKey).(*auth.Claims)
	return claims
}