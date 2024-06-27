package main

import (
	"net/http"
	"strconv"

	"github.com/DomenicoDicosimo/Chirpy-Bootdev/internal/models"
)

func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetChirps()
	if err != nil {
		sendJSONError(w, 500, "Could not get chirp")
		return
	}

	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {

	idStr := r.PathValue("id")
	if idStr == "" {
		sendJSONError(w, http.StatusBadRequest, "Missing Chirp ID")
		return
	}
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		sendJSONError(w, http.StatusBadRequest, "Id not an integer")
	}

	chirp, err := cfg.db.GetChirpByID(idInt)
	if err != nil {
		sendJSONError(w, 500, "Could not get chirp")
		return
	}

	if chirp == (models.Chirp{}) {
		sendJSONError(w, http.StatusNotFound, "Chirp not found")
		return
	}

	respondWithJSON(w, http.StatusOK, chirp)
}
