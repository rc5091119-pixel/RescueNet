package main

import (
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) handlerUserLocations(w http.ResponseWriter,r *http.Request){
	type parameters struct{
		Lat float64 `json:"lat"`
		Lon float64 `json:"lon"`
	}

	var params parameters

	err := json.NewDecoder(r.Body).Decode(&params)

	if err != nil {
		respondWithError(w,http.StatusBadRequest,"Invalid Body",err)
		return
	}

	
}