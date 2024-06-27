package main

import (
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) usersHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		sendJSONError(w, 400, "Something went wrong")
		return
	}
	user, err := cfg.db.CreateUser(params.Email)
	if err != nil {
		sendJSONError(w, 500, "Could not create user")
		return
	}

	respondWithJSON(w, http.StatusCreated, user)
}
