package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetRoomInfo(
	w http.ResponseWriter,
	r *http.Request,
) {
	roomIDStr := r.PathValue("roomID")

	roomID, err := uuid.Parse(roomIDStr)
	if err != nil {
		respondWithError(w, 400, "Invalid room ID", err)
		return
	}

	room, err := cfg.db.GetRoomInfo(
		r.Context(),
		roomID,
	)
	if err != nil {
		respondWithError(w, 500, "Couldn't get room info", err)
		return
	}

	respondWithJSON(w, 200, room)
}
