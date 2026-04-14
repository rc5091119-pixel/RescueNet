package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/rc5091119-pixel/rescuenet/internal/auth"
	"github.com/rc5091119-pixel/rescuenet/internal/database"
)

type User struct {
	Id        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

func (cfg *apiConfig) handlerCreateUsers(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		User
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not Allow", http.StatusMethodNotAllowed)
		return
	}

	params := parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not get the email or password", err)
		return
	}

	UuId := uuid.New()
	hashedPassword,err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w,http.StatusBadRequest,"could not hash password",err)
		return 
	}

	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		ID:           UuId,
		Email:        params.Email,
		PasswordHash: hashedPassword,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not get user", err)
		return
	}

	respondWithJSON(w, 201, response{
		User: User{
			Id:        user.ID,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		},
	})
}
