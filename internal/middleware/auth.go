package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/rhajizada/gazette/internal/oauth"
)

// contextKey is used to avoid key collisions in context.
type contextKey string

// UserContextKey is the key under which JWT claims are stored in the request context.
var UserContextKey = contextKey("user")

// APIAuthMiddleware returns a middleware that:
// 1) extracts the JWT from the Authorization header
// 2) verifies it using the provided secret
// 3) injects the resulting ApplicationClaims into the request context
// and calls the next handler if successful, or returns 401 otherwise.
func APIAuthMiddleware(secret []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rawToken, err := oauth.ExtractTokenFromHeaders(r)
			if err != nil {
				http.Error(w, "unauthorized: token not found", http.StatusUnauthorized)
				return
			}

			claims, err := oauth.VerifyToken(rawToken, secret)
			if err != nil {
				msg := fmt.Sprintf("unauthorized: %v", err)
				http.Error(w, msg, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserClaims(r *http.Request) *oauth.ApplicationClaims {
	if claims, ok := r.Context().Value(UserContextKey).(*oauth.ApplicationClaims); ok {
		return claims
	}
	return nil
}
