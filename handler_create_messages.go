package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/rc5091119-pixel/rescuenet/internal/database"
)

func (cfg *apiConfig) handlerCreateMessage(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Content string `json:"content"`
	}
	var params parameters

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
	err = json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request", err)
		return
	}
	message, err := cfg.db.CreateMessage(r.Context(), database.CreateMessageParams{
		RoomID:   roomID,
		SenderID: uid,
		Content:  params.Content,
	})

	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"Failed to create message",
			err,
		)
		return
	}

	respondWithJSON(w, http.StatusCreated, message)

}
