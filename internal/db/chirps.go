package db

import (
	"github.com/DomenicoDicosimo/Chirpy-Bootdev/internal/models"
)

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) (models.Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return models.Chirp{}, err
	}

	newID := len(dbStructure.Chirps) + 1
	chirp := models.Chirp{ID: newID, Body: body}
	dbStructure.Chirps[newID] = chirp

	if err := db.writeDB(&dbStructure); err != nil {
		return models.Chirp{}, err
	}

	return chirp, nil

}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]models.Chirp, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	chirps := make([]models.Chirp, 0, len(dbStructure.Chirps))
	for _, chirp := range dbStructure.Chirps {
		chirps = append(chirps, chirp)
	}

	return chirps, nil
}

// GetChirps returns a chirp in the database by ID
func (db *DB) GetChirpByID(id int) (models.Chirp, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	dbStructure, err := db.loadDB()
	if err != nil {
		return models.Chirp{}, err
	}

	chirp, exists := dbStructure.Chirps[id]
	if !exists {
		return models.Chirp{}, nil
	}

	return chirp, nil
}
