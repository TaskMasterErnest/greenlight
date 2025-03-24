package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

// a readJSON helper function to help with reading JSON input
// we use this to also triage errors regarding JSON input adn provide a suitable error message
func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dest any) error {
	// decode request into the target destination
	err := json.NewDecoder(r.Body).Decode(dest)
	if err != nil {
		// triage errors, using common;y gotten errors
		var syntaxError *json.SyntaxError
		var unMarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		// check whether the error has the type json.SyntaxError
		// return a plain-English error message which includes the location of the problem
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

			// in some circumstances, it may return an io.ErrUnexpectedEOF for syntax errors in the JSON
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

			// catch any UnmarshalTypeError when the JSON type is wrong for the target destination
		case errors.As(err, unMarshalTypeError):
			if unMarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unMarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unMarshalTypeError.Offset)

			// check if the request body is empty, this gives a io.EOF error
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

			// throw an InvalidUnmarshalError if we pass something that is not a non-nil pointer to Decode()
			// catch this and panic instead of returning an error to the Handler
		case errors.As(err, &invalidUnmarshalError):
			panic(err)

			// for anything else, return the error message as is
		default:
			return err
		}
	}

	return nil
}
