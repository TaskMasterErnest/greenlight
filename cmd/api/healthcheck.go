package main

import (
	"encoding/json"
	"net/http"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// stating the data in a map object type
	data := map[string]string{
		"Status":      "Available",
		"Environment": app.config.env,
		"Version":     version,
	}

	// marshal the data struct into JSON
	js, err := json.Marshal(data)
	if err != nil {
		// this uses the error logger to log the err from marshalling the json
		app.logger.Error(err.Error())
		http.Error(w, "Error processing JSON marshalling", http.StatusInternalServerError)
		return
	}

	// a nicety to make this viewable in terminals
	js = append(js, '\n')

	// setting header to recognize and parse json
	w.Header().Set("Content-Type", "application/json")

	w.Write([]byte(js))
}
