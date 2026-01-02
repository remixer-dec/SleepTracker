package middleware

import (
	"context"
	"net/http"
	"sleeptracker/internal/auth"
	"strings"
)

type contextKey string

const (
	UserContextKey    contextKey = "user_claims"
	IsOwnerContextKey contextKey = "is_owner"
	UserIDContextKey  contextKey = "user_id"
)

type AuthMiddleware struct {
	auth *auth.Auth
}

func NewAuthMiddleware(a *auth.Auth) *AuthMiddleware {
	return &AuthMiddleware{auth: a}
}

func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			cookie, err := r.Cookie("token")
			if err != nil {
				ctx := context.WithValue(r.Context(), IsOwnerContextKey, false)
				ctx = context.WithValue(ctx, UserIDContextKey, "")
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
			authHeader = "Bearer " + cookie.Value
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			ctx := context.WithValue(r.Context(), IsOwnerContextKey, false)
			ctx = context.WithValue(ctx, UserIDContextKey, "")
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := m.auth.ValidateToken(tokenString)
		if err != nil {
			ctx := context.WithValue(r.Context(), IsOwnerContextKey, false)
			ctx = context.WithValue(ctx, UserIDContextKey, "")
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, claims)
		ctx = context.WithValue(ctx, IsOwnerContextKey, claims.IsOwner)
		ctx = context.WithValue(ctx, UserIDContextKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *AuthMiddleware) RequireOwner(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isOwner, ok := r.Context().Value(IsOwnerContextKey).(bool)
		if !ok || !isOwner {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(UserIDContextKey).(string)
		if !ok || userID == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func GetUserID(ctx context.Context) string {
	userID, _ := ctx.Value(UserIDContextKey).(string)
	return userID
}

func IsOwner(ctx context.Context) bool {
	isOwner, _ := ctx.Value(IsOwnerContextKey).(bool)
	return isOwner
}
