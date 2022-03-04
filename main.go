// Package main is the entrypoint of the pmd-dx-api,
// setting up the server and starting it.
package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/janek64/pmd-dx-api/api/db"
	"github.com/janek64/pmd-dx-api/api/handler"
)

// getEnv returns a value from the environment or a default value if it is not defined.
func getEnv(key string, defaultValue string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	return value
}

func main() {
	// Setup the database connection pool
	err := db.InitDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	// Close the connection pool when exiting the program
	defer func() {
		err = db.CloseDB()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to close database connection pool: %v\n", err)
			os.Exit(1)
		}
	}()

	// Get port from environment
	port := getEnv("PORT", "3000")

	// Create a new ServeMux that will handle requests
	mux := http.NewServeMux()

	// Register all handlers
	mux.HandleFunc("/v1/abilities", handler.AbilityListHandler)
	mux.HandleFunc("/v1/camps", handler.CampListHandler)
	mux.HandleFunc("/v1/dungeons", handler.DungeonListHandler)
	mux.HandleFunc("/v1/moves", handler.MoveListHandler)
	mux.HandleFunc("/v1/pokemon", handler.PokemonListHandler)
	mux.HandleFunc("/v1/types", handler.PokemonTypeListHandler)

	// Start the server with the created ServeMux and specified port
	fmt.Printf("pmd-dx-api listening on port %v", port)
	http.ListenAndServe(":"+port, mux)
}
