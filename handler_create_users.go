package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rc5091119-pixel/rescuenet/internal/auth"
	"github.com/rc5091119-pixel/rescuenet/internal/database"
)

type User struct {
	Id        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

func (cfg *apiConfig) handlerCreateUsers(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name     string `json:"name"`
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
	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not hash password", err)
		return
	}
	email := strings.TrimSpace(params.Email)
	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		ID: UuId,
		Name: sql.NullString{
			String: params.Name,
			Valid:  true,
		},
		Email:        email,
		PasswordHash: hashedPassword,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not get user", err)
		return
	}

	respondWithJSON(w, 201, response{
		User: User{
			Id:        user.ID,
			Name:      user.Name.String,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		},
	})
}
