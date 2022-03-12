// Package main is the entrypoint of the pmd-dx-api,
// setting up the server and starting it.
package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/janek64/pmd-dx-api/api/cache"
	"github.com/janek64/pmd-dx-api/api/db"
	"github.com/janek64/pmd-dx-api/api/handler"
	"github.com/janek64/pmd-dx-api/api/logger"
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

	// Initialize the logger
	err := logger.InitLogger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error setting up logger: %v\n", err)
		os.Exit(1)
	}
	// Close the logs files when exiting the program
	defer logger.CloseLogger()

	// Setup the database connection pool
	err = db.InitDB()
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

	// Initialize the redis connection
	err = cache.InitRedis()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to redis: %v\n", err)
		os.Exit(1)
	}
	// Close the redis connection when exiting the program
	defer func() {
		err = cache.CloseRedis()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to close redis connection: %v\n", err)
			os.Exit(1)
		}
	}()

	// Get port from environment
	port := getEnv("PORT", "3000")

	// Create a new httprouter that will handle requests
	router := httprouter.New()

	// Define the middleware chains
	defaultMiddleware := func(h httprouter.Handle) httprouter.Handle {
		return middleware.LogRequest(middleware.CacheResponse(middleware.FieldLimitingParams(h)))
	}
	resourceListMiddleware := func(h httprouter.Handle) httprouter.Handle {
		return defaultMiddleware(middleware.ResourceListParams(h))
	}

	// Register all handlers
	router.GET("/v1/abilities", resourceListMiddleware(handler.AbilityListHandler))
	router.GET("/v1/abilities/:searcharg", defaultMiddleware(handler.AbilitySearchHandler))
	router.GET("/v1/camps", resourceListMiddleware(handler.CampListHandler))
	router.GET("/v1/camps/:searcharg", defaultMiddleware(handler.CampSearchHandler))
	router.GET("/v1/dungeons", resourceListMiddleware(handler.DungeonListHandler))
	router.GET("/v1/dungeons/:searcharg", defaultMiddleware(handler.DungeonSearchHandler))
	router.GET("/v1/moves", resourceListMiddleware(handler.MoveListHandler))
	router.GET("/v1/moves/:searcharg", defaultMiddleware(handler.MoveSearchHandler))
	router.GET("/v1/pokemon", resourceListMiddleware(handler.PokemonListHandler))
	router.GET("/v1/pokemon/:searcharg", defaultMiddleware(handler.PokemonSearchHandler))
	router.GET("/v1/types", resourceListMiddleware(handler.PokemonTypeListHandler))
	router.GET("/v1/types/:searcharg", defaultMiddleware(handler.PokemonTypeSearchHandler))

	// Start the server with the created router and specified port
	fmt.Printf("pmd-dx-api listening on port %v\n", port)
	http.ListenAndServe(":"+port, router)
}
