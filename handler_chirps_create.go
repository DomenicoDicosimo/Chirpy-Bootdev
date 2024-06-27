package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
)

func (cfg *apiConfig) handlerChirpsCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		sendJSONError(w, 400, "Something went wrong")
		return
	}
	if len(params.Body) > 140 {
		sendJSONError(w, 400, "Chirp is too long")
		return
	}
	scrubbedString := scrubString(params.Body)

	chirp, err := cfg.db.CreateChirp(scrubbedString)
	if err != nil {
		sendJSONError(w, 500, "Could not create chirp")
		return
	}

	respondWithJSON(w, http.StatusCreated, chirp)
}

func scrubString(message string) string {
	badwords := [3]string{"kerfuffle", "sharbert", "fornax"}
	words := strings.Split(message, " ")

	lowercaseMessage := strings.ToLower(message)
	lowerWords := strings.Split(lowercaseMessage, " ")

	for i, word := range lowerWords {
		if slices.Contains(badwords[:], word) {
			words[i] = "****"
		}
	}
	cleanedWords := strings.Join(words, " ")
	return cleanedWords
}
