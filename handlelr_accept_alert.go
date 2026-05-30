package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/rc5091119-pixel/rescuenet/internal/database"
)

func (cfg *apiConfig) handlerAcceptAlert(w http.ResponseWriter, r *http.Request) {
	uid, ok := r.Context().Value(userIDKey).(uuid.UUID)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Invalid user", nil)
		return
	}

	alertIDStr := r.PathValue("id")
	alertID, err := uuid.Parse(alertIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid alert", err)
		return
	}

	count, err := cfg.db.CountAcceptedUsers(r.Context(), alertID)
	if err != nil {
		respondWithError(w, 500, "Failed to count users", err)
		return
	}

	if count >= 3 {
		respondWithError(w, 400, "Limit reached", nil)
		return
	}

	err = cfg.db.CreateAlertResponse(r.Context(), database.CreateAlertResponseParams{
		AlertID: alertID,
		UserID:  uid,
	})
	if err != nil {
		respondWithError(w, 500, "Failed to accept alert", err)
		return
	}
	count, err = cfg.db.CountAcceptedUsers(r.Context(), alertID)
	if err != nil {
		respondWithError(w, 500, "Failed to count users", err)
		return
	}
	if count == 3 {
		roomID := uuid.New()
		err = cfg.db.CreateRoom(r.Context(), database.CreateRoomParams{
			ID:      roomID,
			AlertID: alertID,
		})
		if err != nil {
			respondWithError(w, 400, "Can't created a room", err)
			return
		}
		users, err := cfg.db.GetAcceptedUsers(r.Context(), alertID)
		if err != nil {
			respondWithError(w, 500, "Failed to get users", err)
			return
		}
		for _, userID := range users {
			err := cfg.db.AddRoomMember(r.Context(), database.AddRoomMemberParams{
				RoomID: roomID,
				UserID: userID,
			})
			if err != nil {
				respondWithError(w, 400, "can't add the room members", err)
				return
			}
		}
	}
	respondWithJSON(w, 200, map[string]string{
		"message": "Accepted successfully",
	})
}
