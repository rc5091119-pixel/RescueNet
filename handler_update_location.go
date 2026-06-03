package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/rc5091119-pixel/rescuenet/internal/database"
)

func (cfg *apiConfig) handlerUpdateLocation(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Lat float64 `json:"lat"`
		Lng float64 `json:"lng"`
	}

	var params parameters
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Body", err)
		return
	}

	uid, ok := r.Context().Value(userIDKey).(uuid.UUID)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Invalid user ID", nil)
		return
	}

	err = cfg.db.UpdateUserLocations(r.Context(), database.UpdateUserLocationsParams{
		UserID:    uid,
		Latitude:  params.Lat,
		Longitude: params.Lng,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update location", err)
		return
	}

	respondWithJSON(w, 200, map[string]string{
		"message": "Location updated successfully",
	})
}
