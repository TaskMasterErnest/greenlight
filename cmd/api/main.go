package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
}

type application struct {
	config config
	logger *slog.Logger
}

func main() {
	// initialize a config struct instance
	var cfg config

	// populate cfg with values from the command-line arguments
	flag.IntVar(&cfg.port, "port", 4567, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment(development|staging|production)")
	flag.Parse()

	// initialize a new logger instance
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// initialize an instance of the application struct
	app := &application{
		config: cfg,
		logger: logger,
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	// start the server
	logger.Info("Starting server...", "port", server.Addr, "environment", cfg.env)

	err := server.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)

}
