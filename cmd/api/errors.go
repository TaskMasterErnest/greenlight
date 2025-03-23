package main

import (
	"fmt"
	"net/http"
)

// the logError helper to log an error message with the method used and the URL requested
func (app *application) logError(r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)

	// log the error with the components
	app.logger.Error(err.Error(), "method", method, "URI", uri)
}

// errorResponse is a generic helper function that writes the error message to the user in JSON format
// and with a given status code. The error message type is not strict to as to allow flexibility
func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	// envelope the error message
	env := envelope{"error": message}

	// convert to JSON and write to user
	// if this fails, log a 500 error to the user
	err := app.writeJSON(w, status, env, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(500)
	}
}

// a detailed serverErrorResponse() method to log server errors at runtime
func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	// log the error gotten
	app.logError(r, err)

	// craft a message
	message := "The server encountered a problem and could not process your request"
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

// a detailed notFoundResponse
func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {

	message := "The requested resource could not be found"

	app.errorResponse(w, r, http.StatusNotFound, message)
}

// a detailed methodNotAllowedResponse error response
func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request, err error) {
	message := fmt.Sprintf("The %s method is not supported for this resource", r.Method)

	app.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}
