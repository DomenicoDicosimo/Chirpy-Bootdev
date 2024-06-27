package main

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"

	"github.com/DomenicoDicosimo/Chirpy-Bootdev/internal/db"
)

type apiConfig struct {
	fileserverHits int
	db             *db.DB
}

func myhandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	body := ([]byte)("OK")
	w.Write(body)
}

func (cfg *apiConfig) chirpHandler(w http.ResponseWriter, r *http.Request) {
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

func (cfg *apiConfig) getChirpsHandler(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetChirps()
	if err != nil {
		sendJSONError(w, 500, "Could not get chirp")
		return
	}

	respondWithJSON(w, http.StatusOK, chirps)
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
	db, err := db.NewDB("./db.json")
	if err != nil {
		panic(err)
	}

	cfg := &apiConfig{db: db}

	mux.HandleFunc("GET /api/healthz", myhandler)

	fileServer := http.FileServer(http.Dir("."))
	mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app", fileServer)))

	mux.HandleFunc("GET /admin/metrics", cfg.metricsHandler)
	mux.HandleFunc("GET /api/reset", cfg.resetHandler)
	mux.HandleFunc("POST /api/chirps", cfg.chirpHandler)
	mux.HandleFunc("GET /api/chirps", cfg.getChirpsHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	server.ListenAndServe()

}
