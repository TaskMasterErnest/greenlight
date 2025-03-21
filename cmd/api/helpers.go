package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// create an envelope type
type envelope map[string]any

func (app *application) readIDParams(r *http.Request) (int64, error) {
	// get the parameters from the context in the request URL
	params := httprouter.ParamsFromContext(r.Context())

	// convert the params into an int
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid ID parameter")
	}

	return id, nil
}

// a writeJSON helper to help with encoding data into JSON.
// it takes in the responseWriter, the status code to send, the data to encode, any HTTP headers and returns an error
// modify the date to be of type envelope
func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	// marshal the data
	js, err := json.MarshalIndent(data, "", "    ") // 4 spaces for the indentation
	if err != nil {
		return err
	}

	// add a newline nicety
	js = append(js, '\n')

	// if there are headers available, loop through each header in the header map
	// add the headers to the responseWriter header map
	for key, value := range headers {
		w.Header()[key] = value
	}

	// add the content-type header to enable the parsing of the data as a JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(js))

	return nil
}
