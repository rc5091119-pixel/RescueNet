package main

import (
	"context"
	"net/http"
	"strings"

	"github.com/rc5091119-pixel/rescuenet/internal/auth"
)

type contextKey string

const userIDKey contextKey = "user_id"

func (cfg *apiConfig) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			respondWithError(w, http.StatusUnauthorized, "Missing token", nil)
			return
		}
		if !strings.HasPrefix(authHeader, "Bearer ") {
			respondWithError(w, http.StatusUnauthorized, "Invalid token format", nil)
			return
		}
		tokenstr := strings.TrimPrefix(authHeader, "Bearer ")
		userID, err := auth.ValidateJWT(tokenstr, cfg.jwtSecret)
		if err != nil {
			respondWithError(w, 401, "invalid token", err)
			return
		}
		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
