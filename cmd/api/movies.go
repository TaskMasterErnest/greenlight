package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/TaskMasterErnest/greenlight/internal/data"
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
		Title   string   `json:"title"`
		Year    int32    `json:"year"`
		Runtime int32    `json:"runtime"`
		Genres  []string `json:"genres"`
	}

	// initialize a json.Decoder instance to decode the data from client into the input struct
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// respond with the contents of the input struct
	fmt.Fprintf(w, "%+v\n", input)
}
