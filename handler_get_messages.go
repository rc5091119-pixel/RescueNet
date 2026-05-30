package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/rc5091119-pixel/rescuenet/internal/database"
)

func (cfg *apiConfig) handlerGetMessage(w http.ResponseWriter, r *http.Request) {
	roomIDstr := r.PathValue("roomID")
	roomID, err := uuid.Parse(roomIDstr)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid RoomId", err)
		return
	}

	uid, ok := r.Context().Value(userIDKey).(uuid.UUID)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Invalid user", nil)
		return
	}

	isMember, err := cfg.db.IsRoomMember(r.Context(), database.IsRoomMemberParams{
		RoomID: roomID,
		UserID: uid,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError,
			"Failed to verify room membership", err)
		return
	}

	if !isMember {
		respondWithError(w, http.StatusForbidden,
			"You are not a member of this room", nil)
		return
	}
	messages, err := cfg.db.GetRoomMessages(r.Context(), roomID)
	if err != nil {
		respondWithError(w,
			http.StatusInternalServerError,
			"Can't fetch messages",
			err,
		)
		return
	}

	respondWithJSON(w, http.StatusOK, messages)
}
