// Package main is the entrypoint of the pmd-dx-api,
// setting up the server and starting it.
package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/janek64/pmd-dx-api/api/db"
	"github.com/janek64/pmd-dx-api/api/handler"
	"github.com/julienschmidt/httprouter"
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

	// Create a new httprouter that will handle requests
	router := httprouter.New()

	// Register all handlers
	router.GET("/v1/abilities", handler.AbilityListHandler)
	router.GET("/v1/abilities/:searcharg", handler.AbilitySearchHandler)
	router.GET("/v1/camps", handler.CampListHandler)
	router.GET("/v1/camps/:searcharg", handler.CampSearchHandler)
	router.GET("/v1/dungeons", handler.DungeonListHandler)
	router.GET("/v1/dungeons/:searcharg", handler.DungeonSearchHandler)
	router.GET("/v1/moves", handler.MoveListHandler)
	router.GET("/v1/moves/:searcharg", handler.MoveSearchHandler)
	router.GET("/v1/pokemon", handler.PokemonListHandler)
	router.GET("/v1/types", handler.PokemonTypeListHandler)

	// Start the server with the created router and specified port
	fmt.Printf("pmd-dx-api listening on port %v\n", port)
	http.ListenAndServe(":"+port, router)
}
