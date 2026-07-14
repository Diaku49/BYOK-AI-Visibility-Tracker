package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/Diaku49/AI-visibility-tracker/internal/pkg"
)

type contextKey string

const (
	userIDKey contextKey = "userID"
	tierKey   contextKey = "tier"
)

func Authenticate(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			token, found := strings.CutPrefix(authHeader, "Bearer ")
			if !found {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			valid, err := pkg.ValidateToken(token, []byte(secret))
			if err != nil || valid == nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, valid.UserID)
			ctx = context.WithValue(ctx, tierKey, valid.Tier)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
