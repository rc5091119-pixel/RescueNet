package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/rc5091119-pixel/rescuenet/internal/database"
)

func(cfg *apiConfig)handlerAcceptAlert(w http.ResponseWriter,r *http.Request){
	uid,ok := r.Context().Value(userIDKey).(uuid.UUID)
	if !ok {
		respondWithError(w,http.StatusUnauthorized,"Invalid user",nil)
		return
	}

	alertIDStr := r.PathValue("id")
	alertID,err := uuid.Parse(alertIDStr)
	if err != nil {
		respondWithError(w,http.StatusBadRequest,"Invalid alert",err)
		return
	}
	count,err := cfg.db.CountAcceptedUsers(r.Context(),alertID)
	if err !=nil {
		respondWithError(w,500,"Failed to count users",err)
		return
	}

	if count >= 15 {
		respondWithError(w,400,"Limit reached",nil)
		
	}
	err = cfg.db.CreateAlertResponse(r.Context(),database.CreateAlertResponseParams{
		AlertID: alertID,
		UserID: uid,
	})
	if err != nil {
		respondWithError(w, 500, "Failed to accept alert", err)
		return
	}

	respondWithJSON(w, 200, map[string]string{
		"message": "Accepted successfully",
	})
}