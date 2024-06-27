package main

import (
	"net/http"

	"github.com/DomenicoDicosimo/Chirpy-Bootdev/internal/db"
)

type apiConfig struct {
	fileserverHits int
	db             *db.DB
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
	mux.HandleFunc("POST /api/chirps", cfg.handlerChirpsCreate)
	mux.HandleFunc("GET /api/chirps", cfg.handlerChirpsGet)
	mux.HandleFunc("GET /api/chirps/{id}", cfg.handlerChirpsRetrieve)
	mux.HandleFunc("POST /api/users", cfg.usersHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	server.ListenAndServe()

}
