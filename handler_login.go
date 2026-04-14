package main

import (
	"go/types"
	"net/http"
)

func(cfg *apiConfig)handlerLoginUsers(w http.ResponseWriter,r *http.Request){
	type parameters struct{
		Email string `json:"email"`
		Password string `json:"password"`
	}
	type response struct{
		User
	}

	
}