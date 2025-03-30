package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/TaskMasterErnest/greenlight/internal/data"
	"github.com/TaskMasterErnest/greenlight/internal/validator"
)

// showMovieHandler
func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	// get the ID params from the context
	// the ID is a integer and the params are strings, convert them to int
	id, err := app.readIDParams(r)
	if err != nil {
		app.notFoundResponse(w, r)
	}

	// initialize a Movie struct instance
	//using ID from context and some dummy data
	movie := data.Movie{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Casablanca",
		Runtime:   102,
		Genres:    []string{"drama", "romance", "war"},
		Version:   1,
	}

	// envelope the movie in the envelope type
	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.logger.Error(err.Error())
		app.serverErrorResponse(w, r, err)
	}
}

// createMovieHandler
func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	// create struct to hold movie data
	var input struct {
		Title   string       `json:"title"`
		Year    int32        `json:"year"`
		Runtime data.Runtime `json:"runtime"` // make this field a data.Runtime type
		Genres  []string     `json:"genres"`
	}

	// initialize a json.Decoder instance to decode the data from client into the input struct
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// initialize a new Validator instance for input Validation
	v := validator.New()

	// use the Check method from the validator to execute validation checks
	// this will add errors to the errors map if the validations do not evaluate to true
	// <validating Title input>
	v.Check(input.Title != "", "title", "must be provided")
	v.Check(len(input.Title) <= 500, "title", "must not be more than 500 bytes long")

	// <validating Year input>
	v.Check(input.Year != 0, "year", "must be provided")
	v.Check(input.Year >= 1888, "year", "must be greater than 1888")
	v.Check(input.Year <= int32(time.Now().Year()), "year", "must not be in the future")

	// <validating Runtime input>
	v.Check(input.Runtime != 0, "runtime", "must be provided")
	v.Check(input.Runtime > 0, "runtime", "must be a positive integer")

	// <validating Genre input>
	v.Check(input.Genres != nil, "genres", "must be provided")
	v.Check(len(input.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(input.Genres) <= 5, "genres", "must not contain more than 5 genres")
	// now we check if all the genres are unique
	v.Check(validator.Unique(input.Genres), "genres", "must not contain duplicate values")

	// check if any of the checks failed with the Valid() method.
	// if affirmative, then use the failedValidationResponse helper to send a response to the client
	// pass in the v.Errors map, map of errors accumulated
	if !v.Valid() {
		app.FailedValidationResponse(w, r, v.Errors)
		return
	}

	// respond with the contents of the input struct
	fmt.Fprintf(w, "%+v\n", input)
}
