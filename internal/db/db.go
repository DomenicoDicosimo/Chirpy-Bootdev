package db

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/DomenicoDicosimo/Chirpy-Bootdev/internal/models"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]models.Chirp `json:"chirps"`
}

// NewDB creates a new database connection
// and creates the database file if it doesn't exist
func NewDB(path string) (*DB, error) {
	db := &DB{
		path: path,
		mux:  &sync.RWMutex{},
	}

	if err := db.ensureDB(); err != nil {
		return nil, err
	}

	return db, nil
}

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

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error {
	if _, err := os.Stat(db.path); os.IsNotExist(err) {
		file, err := os.Create(db.path)
		if err != nil {
			return err
		}
		defer file.Close()

		initialData := DBStructure{Chirps: make(map[int]models.Chirp)}
		jsonData, err := json.Marshal(initialData)
		if err != nil {
			return err
		}

		if _, err := file.Write(jsonData); err != nil {
			return err
		}
	}
	return nil
}

// loadDB reads the database file into memory
func (db *DB) loadDB() (DBStructure, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	file, err := os.Open(db.path)
	if err != nil {
		return DBStructure{}, err
	}
	defer file.Close()

	var dbStructure DBStructure
	if err := json.NewDecoder(file).Decode(&dbStructure); err != nil {
		return DBStructure{}, err
	}

	return dbStructure, nil
}

// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure *DBStructure) error {
	db.mux.Lock()
	defer db.mux.Unlock()

	file, err := os.OpenFile(db.path, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(dbStructure); err != nil {
		return err
	}

	return nil
}
