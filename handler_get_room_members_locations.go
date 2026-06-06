package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetRoomMemberLocations(w http.ResponseWriter, r *http.Request) {
	roomIDString := r.PathValue("roomID")

	roomID, err := uuid.Parse(roomIDString)
	if err != nil {
		respondWithError(
			w,
			http.StatusBadRequest,
			"Invalid room ID",
			err,
		)
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
	locations, err := cfg.db.GetRoomMemberLocations(
		r.Context(),
		roomID,
	)
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"Couldn't get room locations",
			err,
		)
		return
	}

	respondWithJSON(
		w,
		http.StatusOK,
		locations,
	)
}


