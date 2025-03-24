package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	// initialize a new httpRouter instance
	router := httprouter.New()

	// adding custom error handling for certain routes
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	// register the routes
	router.HandlerFunc(http.MethodGet, "/v1/healthz", app.healthCheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/movies", app.createMovieHandler)
	router.HandlerFunc(http.MethodGet, "/v1/movies/:id", app.showMovieHandler)

	// wrap the call to router with the recoverPanic middleware
	return app.recoverPanic(router)
}
