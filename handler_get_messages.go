package main

import (
	"net/http"

	"github.com/google/uuid"
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

	err = cfg.VerifyRoomMember(r.Context(), roomID, uid)
	if err != nil {
		respondWithError(
			w,
			http.StatusForbidden,
			"You are not a member of this room",
			nil,
		)
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
