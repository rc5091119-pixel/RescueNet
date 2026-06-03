package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetUserRooms(
	w http.ResponseWriter,
	r *http.Request,
) {
	uid, ok := r.Context().Value(userIDKey).(uuid.UUID)
	if !ok {
		respondWithError(w, 401, "Invalid user", nil)
		return
	}

	rooms, err := cfg.db.GetUserRooms(
		r.Context(),
		uid,
	)
	if err != nil {
		respondWithError(
			w,
			500,
			"Couldn't get rooms",
			err,
		)
		return
	}

	respondWithJSON(w, 200, rooms)
}
