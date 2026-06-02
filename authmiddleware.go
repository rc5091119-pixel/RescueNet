package main

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/rc5091119-pixel/rescuenet/internal/auth"
	"github.com/rc5091119-pixel/rescuenet/internal/database"
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

func (cfg *apiConfig) VerifyRoomMember(ctx context.Context, roomID uuid.UUID, userID uuid.UUID) error {
	isMember, err := cfg.db.IsRoomMember(ctx, database.IsRoomMemberParams{
		RoomID: roomID,
		UserID: userID,
	})
	if err != nil {
		return err
	}

	if !isMember {
		return errors.New("user is not a room member")
	}

	return nil
}
