package main

import (
	"fmt"
	"net/http"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// inscribing a json text raw
	js := `{"status": "available", "Environment": %q, "Version": %q}`
	js = fmt.Sprintf(js, app.config.env, version)

	// setting header to recognize and parse json
	w.Header().Set("Content-Type", "application/json")

	w.Write([]byte(js))
}
