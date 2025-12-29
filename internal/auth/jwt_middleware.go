package auth

import (
	"context"
	"net/http"
	"strings"
)

type key int

const UserIDKey key = 1

func JWTMiddleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "missing or invalid auth header", http.StatusUnauthorized)
				return
			}
			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := ParseJWT(tokenStr, secret)
			if err != nil {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}
			// You can extract user_id from claims here
			userIDf, ok := claims["user_id"].(float64) // JWT numbers are float64
			if !ok {
				http.Error(w, "invalid token payload", http.StatusUnauthorized)
				return
			}
			userID := int64(userIDf)
			// Store userID in request context
			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UserIDFromContext(ctx context.Context) (int64, bool) {
	uid, ok := ctx.Value(UserIDKey).(int64)
	return uid, ok
}
