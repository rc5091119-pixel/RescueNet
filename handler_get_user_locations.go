package main

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/rc5091119-pixel/rescuenet/internal/database"
)

func (cfg *apiConfig) handlerCreateAlerts(w http.ResponseWriter, r *http.Request) {

	uid, ok := r.Context().Value(userIDKey).(uuid.UUID)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Invalid user ID", nil)
		return
	}
	userLoc, err := cfg.db.GetUserLocationByUserID(r.Context(), uid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, 400, "User location not found", nil)
			return
		}
		respondWithError(w, 500, "DB error", err)
		return
	}

	alert, err := cfg.db.CreateAlert(r.Context(), database.CreateAlertParams{
		UserID:    uid,
		Latitude:  userLoc.Latitude,
		Longitude: userLoc.Longitude,
	})
	if err != nil {
		respondWithError(w, 500, "Failed to create alert", err)
		return
	}

	users, err := cfg.db.GetUserLocations(r.Context(), database.GetUserLocationsParams{
		Lat: userLoc.Latitude,
		Lng: userLoc.Longitude,
	})

	var nearbyUsers []database.GetUserLocationsRow

	for _, u := range users {
		if u.UserID == uid {
			continue
		}
		distance := haversine(userLoc.Latitude, userLoc.Longitude, u.Latitude, u.Longitude)
		if distance <= 1.0 {
			nearbyUsers = append(nearbyUsers, u)
		}
	}
	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"alert_id":     alert.ID,
		"nearby_users": nearbyUsers,
	})
}
