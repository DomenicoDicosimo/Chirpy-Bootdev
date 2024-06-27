package main

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"
)

type apiConfig struct {
	fileserverHits int
}

func myhandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	body := ([]byte)("OK")
	w.Write(body)
}

func validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type returnVals struct {
		Cleaned_Body string `json:"cleaned_body"`
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

	respondWithJSON(w, http.StatusOK, returnVals{
		Cleaned_Body: scrubString(params.Body),
	})

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

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(dat)
}

func sendJSONError(w http.ResponseWriter, statusCode int, errorMessage string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	// Create the error message object
	response := map[string]string{"error": errorMessage}

	// Convert it to JSON
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Internal server error"}`))
		return
	}

	// Write it to the response
	w.Write(jsonResponse)
}

func main() {
	mux := http.NewServeMux()

	cfg := &apiConfig{}

	mux.HandleFunc("GET /api/healthz", myhandler)

	fileServer := http.FileServer(http.Dir("."))
	mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app", fileServer)))

	mux.HandleFunc("GET /admin/metrics", cfg.metricsHandler)
	mux.HandleFunc("GET /api/reset", cfg.resetHandler)
	mux.HandleFunc("POST /api/validate_chirp", validateChirpHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	server.ListenAndServe()

}
