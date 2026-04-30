package main

import (
	"fmt"
	"net/http"
)
func (cfg *apiConfig) handlerTestProtected(w http.ResponseWriter,r * http.Request){
	userID := r.Context().Value(userIDKey)

	w.Write([]byte("Protected route accesed by user : "))
	w.Write([]byte(fmt.Sprint(userID)))
}