package main

import (
	"fmt"
	"net/http"
)

func (app *application) recoverPanic(next http.Handler) http.Handler {
	// specifying the handler to return
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			// use the built-in recover to check if there have been a panic or not
			if err := recover(); err != nil {
				// if there has been a panic, set a "Connection : close" on the response header
				// this triggers the HTTP server to automatically close the current connection after the response has been sent
				w.Header().Set("Connection", "close")
				// the value returned by the recover is of type any
				// we then use the fmt.Errorf() to normalize it into an error
				// we call the serverErrorResponse helper and log this
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
