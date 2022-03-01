// Package main is the entrypoint of the pmd-dx-api,
// setting up the server and starting it.
package main

import (
	"net/http"
	"os"
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
	// Get port from environment
	port := getEnv("PORT", "3000")

	// Create a new ServeMux that will handle requests
	mux := http.NewServeMux()

	// Start the server with the created ServeMux and specified port
	http.ListenAndServe(":"+port, mux)
}
