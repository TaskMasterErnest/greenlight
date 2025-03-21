package main

import (
	"net/http"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// stating the data in a map object type
	data := map[string]string{
		"Status":      "Available",
		"Environment": app.config.env,
		"Version":     version,
	}

	// call the writeJSON helper method to convert data to JSON
	err := app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		// log the error
		app.logger.Error(err.Error())
		http.Error(w, "The server encountered a problem and could not parse JSON request", http.StatusInternalServerError)
	}
}
