package middleware

import (
	"context"
	"net/http"
)

type MockUserMiddlewareConfig struct {
}

type key string

const (
	UserIdType key = "userId"
)

func InitAuthMiddleware() MockUserMiddlewareConfig {
	return MockUserMiddlewareConfig{}
}

func MockAuthMiddleware(userId int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), UserIdType, userId)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
