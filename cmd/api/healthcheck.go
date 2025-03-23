package main

import (
	"net/http"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// enveloping the healthcheck response
	env := envelope{
		"Status": "Available",
		"System_Info": map[string]string{
			"Environment": app.config.env,
			"Version":     version,
		},
	}

	// call the writeJSON helper method to convert data to JSON
	err := app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		// log the error
		app.logger.Error(err.Error())
		app.serverErrorResponse(w, r, err)
	}
}
