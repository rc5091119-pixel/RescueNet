package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/rc5091119-pixel/rescuenet/internal/auth"
)

func (cfg *apiConfig) handlerLoginUsers(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		User
		Token string `json:"token"`
	}
	params := parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not decode", err)
		return
	}
	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)

	match, err := auth.CheckPasswordHash(params.Password, user.PasswordHash)
	if err != nil || !match {
		respondWithError(w, http.StatusUnauthorized, "password not match", err)
		return
	}
	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour*24)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "could not get the secretKey", err)
		return
	}

	respondWithJSON(w, 201, response{
		User: User{
			Id:        user.ID,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		},
		Token: token,
	})

}
