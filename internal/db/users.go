package db

import "github.com/DomenicoDicosimo/Chirpy-Bootdev/internal/models"

// CreateUser creates a new user and saves it to disk
func (db *DB) CreateUser(email string) (models.User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return models.User{}, err
	}
	newID := len(dbStructure.Users) + 1
	user := models.User{ID: newID, Email: email}
	dbStructure.Users[newID] = user

	if err := db.writeDB(&dbStructure); err != nil {
		return models.User{}, err
	}

	return user, nil
}
