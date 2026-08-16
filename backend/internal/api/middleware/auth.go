package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/Diaku49/AI-visibility-tracker/internal/pkg"
	"github.com/google/uuid"
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

func GetUserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	val := ctx.Value(userIDKey)
	if val == nil {
		return uuid.Nil, errors.New("user ID not found in context")
	}

	userID, ok := val.(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.New("user ID in context has wrong type")
	}

	return userID, nil
}
