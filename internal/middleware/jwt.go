package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/blihor/todo-app/internal/auth"
)

type AuthService interface {
	ValidateToken(tokenString string) (*auth.JwtClaims, error)
}

type jwt struct {
	authService AuthService
	logger      *slog.Logger
}

func NewMiddleware(authService AuthService, logger *slog.Logger) *jwt {
	return &jwt{
		authService: authService,
		logger:      logger,
	}
}

func (m *jwt) Protect(next func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		claims, err := m.authService.ValidateToken(tokenString)
		if err != nil {
			m.logger.Warn("JWT validation failed", "error", err)
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "UserID", claims.UserID)
		next(w, r.WithContext(ctx))
	})
}
