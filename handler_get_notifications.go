package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetNotifications(w http.ResponseWriter, r *http.Request) {
	uid, ok := r.Context().Value(userIDKey).(uuid.UUID)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Invalid user ID", nil)
		return
	}

	notifications, err := cfg.db.GetPendingAlertsForUser(r.Context(), uid)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get notifications", err)
		return
	}

	respondWithJSON(w, http.StatusOK, notifications)
}
