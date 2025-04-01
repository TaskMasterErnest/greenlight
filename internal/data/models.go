package data

import (
	"database/sql"
	"errors"
)

// custom ErrRecordNotFound error; returned from Get() method
var (
	ErrRecordNotFound = errors.New("record not found")
)

// a Models struct that wraps the MovieModel
type Models struct {
	Movies MovieModel
}

// a NewModels() method which returns a Models struct containing the initialized MovieModel
func NewModels(db *sql.DB) Models {
	return Models{
		Movies: MovieModel{DB: db},
	}
}
