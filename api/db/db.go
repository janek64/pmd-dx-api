// Package db contains all functions for interactions with
// the database in the pmd-dx-api, consisting of functions
// for connecting to the db and executing queries.
package db

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

// DBConnectionError- type for database connection error
type DBConnectionError struct {
	MissingVar string
}

// Error - implementation of the error interface.
func (e *DBConnectionError) Error() string {
	return fmt.Sprintf("connecting to database failed because of missing environment variable '%v'", e.MissingVar)
}

// dbpool is the global connection pool for the database.
var dbpool *pgxpool.Pool

// InitDB connects to the database and sets the connection pool global variable.
func InitDB() error {
	// Get connection data from environment
	dbuser, ok := os.LookupEnv("DB_USER")
	if !ok {
		return &DBConnectionError{"DB_USER"}
	}
	dbpassword, ok := os.LookupEnv("DB_PASSWORD")
	if !ok {
		return &DBConnectionError{"DB_PASSWORD"}
	}
	dburl, ok := os.LookupEnv("DB_URL")
	if !ok {
		return &DBConnectionError{"DB_URL"}
	}
	dbname, ok := os.LookupEnv("DB_NAME")
	if !ok {
		return &DBConnectionError{"DB_NAME"}
	}

	// Establish the database connection
	databaseURL := fmt.Sprintf("postgres://%v:%v@%v/%v", dbuser, dbpassword, dburl, dbname)
	var err error
	dbpool, err = pgxpool.Connect(context.Background(), databaseURL)
	if err != nil {
		return err
	}
	// Test the connection pool and return the result
	return dbpool.Ping(context.Background())
}

// CloseDB closes the connection pool to the database stored in the global variable.
func CloseDB() error {
	if dbpool == nil {
		return errors.New("no connection pool to close")
	}
	dbpool.Close()
	return nil
}
