package main

import (
	"errors"
	"fmt"
	"net/http"

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

	// call the Get() method to fetch specific movie data, return errors
	movie, err := app.models.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
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

	// copy movie values from the input struct to a new Movie struct
	movie := &data.Movie{
		Title:   input.Title,
		Year:    input.Year,
		Runtime: input.Runtime,
		Genres:  input.Genres,
	}

	// initialize a new Validator instance for input Validation
	v := validator.New()

	// validate movie with the ValidateMovie function and return a response
	if data.ValidateMovie(v, movie); !v.Valid() {
		app.FailedValidationResponse(w, r, v.Errors)
		return
	}

	// call the Insert method from the movies model, and pass in the pointer to the validated movie struct
	err = app.models.Movies.Insert(movie)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// in the HTTP response header, we display the location at which the client can get their newly-created resource
	// make an empty http.Header map and then use the Set() method to add a new Location header
	// here, we interpolate the system-generated ID for out new movie in the URL
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))

	// write a JSON response with a 201 Created status code, movie data in response body,
	// and the Location header
	err = app.writeJSON(w, http.StatusCreated, envelope{"movie": movie}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	// respond with the contents of the input struct
	fmt.Fprintf(w, "%+v\n", input)
}
