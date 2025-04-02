package data

import (
	"database/sql"
	"time"

	"github.com/TaskMasterErnest/greenlight/internal/validator"
	"github.com/lib/pq"
)

type Movie struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"` // this field not relevant to users
	Title     string    `json:"title,omitempty"`
	Year      int32     `json:"year,omitempty"`
	Runtime   Runtime   `json:"runtime,omitempty"`
	Genres    []string  `json:"genres"`
	Version   int32     `json:"version"`
}

// a ValidateMovie function that will validate all input on the movie struct
// regardless of whether it is a fresh input or an edited input
func ValidateMovie(v *validator.Validator, movie *Movie) {
	// use the Check method from the validator to execute validation checks
	// this will add errors to the errors map if the validations do not evaluate to true
	// <validating Title input>
	v.Check(movie.Title != "", "title", "must be provided")
	v.Check(len(movie.Title) <= 500, "title", "must not be more than 500 bytes long")

	// <validating Year input>
	v.Check(movie.Year != 0, "year", "must be provided")
	v.Check(movie.Year >= 1888, "year", "must be greater than 1888")
	v.Check(movie.Year <= int32(time.Now().Year()), "year", "must not be in the future")

	// <validating Runtime input>
	v.Check(movie.Runtime != 0, "runtime", "must be provided")
	v.Check(movie.Runtime > 0, "runtime", "must be a positive integer")

	// <validating Genre input>
	v.Check(movie.Genres != nil, "genres", "must be provided")
	v.Check(len(movie.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(movie.Genres) <= 5, "genres", "must not contain more than 5 genres")
	// now we check if all the genres are unique
	v.Check(validator.Unique(movie.Genres), "genres", "must not contain duplicate values")
}

// methods for performing CRUD to Movies
// a MovieModel struct that wraps an sql.DB connection pool
type MovieModel struct {
	DB *sql.DB
}

// insert a movie record into the Movie table
func (m MovieModel) Insert(movie *Movie) error {
	// defining the SQL query for inserting the new record into the movies table
	// and returning system-generated data
	query := `
			INSERT INTO movies (title, year, runtime, genres)
			VALUES ($1, $2, $3, $4)
			RETURNING id, created_at, version`

	// an args slice to contain the values for the placeholder parameters for the movie struct
	// with this, we can make it clear as to "what values are being used where" in the query
	args := []any{movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres)}

	// using the QueryRow() method to execute the SQL query on the connection pool
	// we pass in the args slice as a variadic parameter and scan the system-generated output into the movie struct
	return m.DB.QueryRow(query, args...).Scan(&movie.ID, &movie.CreatedAt, &movie.Version)
}

// fetching a movie record from the Movie table
func (m MovieModel) Get(id int64) (*Movie, error) {
	return nil, nil
}

// update a specific movie record in the Movie table
func (m MovieModel) Update(movie *Movie) error {
	return nil
}

// delete a specific record from the Movies table
func (m MovieModel) Delete(id int64) error {
	return nil
}
