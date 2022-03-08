// Package main is the entrypoint of the pmd-dx-api,
// setting up the server and starting it.
package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/janek64/pmd-dx-api/api/db"
	"github.com/janek64/pmd-dx-api/api/handler"
	"github.com/janek64/pmd-dx-api/api/middleware"
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
	router.GET("/v1/abilities", middleware.FieldLimitingParams(middleware.ResourceListParams(handler.AbilityListHandler)))
	router.GET("/v1/abilities/:searcharg", middleware.FieldLimitingParams(handler.AbilitySearchHandler))
	router.GET("/v1/camps", middleware.FieldLimitingParams(middleware.ResourceListParams(handler.CampListHandler)))
	router.GET("/v1/camps/:searcharg", middleware.FieldLimitingParams(handler.CampSearchHandler))
	router.GET("/v1/dungeons", middleware.FieldLimitingParams(middleware.ResourceListParams(handler.DungeonListHandler)))
	router.GET("/v1/dungeons/:searcharg", middleware.FieldLimitingParams(handler.DungeonSearchHandler))
	router.GET("/v1/moves", middleware.FieldLimitingParams(middleware.ResourceListParams(handler.MoveListHandler)))
	router.GET("/v1/moves/:searcharg", middleware.FieldLimitingParams(handler.MoveSearchHandler))
	router.GET("/v1/pokemon", middleware.FieldLimitingParams(middleware.ResourceListParams(handler.PokemonListHandler)))
	router.GET("/v1/pokemon/:searcharg", middleware.FieldLimitingParams(handler.PokemonSearchHandler))
	router.GET("/v1/types", middleware.FieldLimitingParams(middleware.ResourceListParams(handler.PokemonTypeListHandler)))
	router.GET("/v1/types/:searcharg", middleware.FieldLimitingParams(handler.PokemonTypeSearchHandler))

	// Start the server with the created router and specified port
	fmt.Printf("pmd-dx-api listening on port %v\n", port)
	http.ListenAndServe(":"+port, router)
}
