package main

import (
	"encoding/json"
	"net/http"

	"github.com/DomenicoDicosimo/Chirpy-Bootdev/internal/auth"
)

func (cfg *apiConfig) handlerWebhooks(w http.ResponseWriter, r *http.Request) {

	key, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find API Key")
		return
	}

	if key != cfg.polkaAPI {
		respondWithError(w, http.StatusUnauthorized, "Invalid Key")
		return
	}

	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID int `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	if params.Event != "user.upgraded" {
		respondWithError(w, http.StatusNoContent, "Incorrect envent code")
		return
	}

	err = cfg.DB.UpgradeUser(params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "User not found")
		return
	}

	respondWithJSON(w, http.StatusNoContent, "User upgraded")
}
