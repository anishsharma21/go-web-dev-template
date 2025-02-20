package middleware

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/anishsharma21/go-web-dev-template/internal/auth"
	"github.com/golang-jwt/jwt/v5"
)

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			slog.WarnContext(r.Context(), "Authorization header was empty")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		_, err := auth.VerifyToken(tokenString)
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				http.Error(w, "Token expired", http.StatusUnauthorized)
				return
			}
			slog.WarnContext(r.Context(), "Failed to verify access token", "error", err, "access_token", tokenString)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
