package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	// import pq driver so it can register itself with the sql package
	// alias to blank identifier to stop Go from complaining that it is not being used
	_ "github.com/lib/pq"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct { // db struct field to hold config settings for db connection pool
		dsn string
		// add db connection pool config
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  time.Duration
	}
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

	// read the DB dsn command-line flag from the config struct
	// default to a DSN for local development
	// create an env-var in .zshrc export GREENLIGHT_DB_DSN='postgres://greenlight:gr33n%23Light@localhost/greenlight?sslmode=disable'
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("GREENLIGHT_DB_DSN"), "Postgres DSN")

	// read the connection pool config params from command line, but have defaults
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.DurationVar(&cfg.db.maxIdleTime, "db-max-idle-time", 15*time.Minute, "PostgreSQL max connection idle time")

	flag.Parse()

	// initialize a new logger instance
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// call the openDB helper function to create the connection pool by passing in the config struct
	// if error occurs, log error and exit application immediately
	db, err := openDB(cfg)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// defer a call to close the db connection pool before the main function exits
	defer db.Close()

	// log message that DB connection pool has been successfully established
	logger.Info("database connection pool established")

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

	err = server.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)

}

/*
**
The openDB() function that returns an sql.DB connection pool
**
*/
func openDB(cfg config) (*sql.DB, error) {
	// create an empty DB connection pool with sql.Open() using the DB DSN from the config
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	// set the database connection pool configurations
	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)
	db.SetConnMaxIdleTime(cfg.db.maxIdleTime)

	// create a context with a 5-second timeout deadline
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// use PingContext to establish a connection to the database, passing in the context we created as a param
	// if connection could not be established within the timeout deadline, return an error
	// if any error is received, we close the connection pool and return the error
	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}

	// return the connection pool
	return db, nil
}
