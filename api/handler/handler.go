// Package handler defines all HTTP handler functions
// for routes of the pmd-dx-api.
package handler

import (
	"encoding/json"
	"net/http"

	"github.com/janek64/pmd-dx-api/api/db"
	"github.com/janek64/pmd-dx-api/api/models"
)

// listResponse defines the JSON structure for lists of API resources
type listResponse struct {
	Count   int                       `json:"count"`
	Results []models.NamedResourceURL `json:"results"`
}

// answerWithListJSON transforms the provided resources to a list with URLs, packages
// them in a JSON and sends it as a response with the provided ResponseWriter.
func answerWithListJSON(resources []models.NamedResourceID, requestedBaseURL string, resourceTypeName string, w http.ResponseWriter) {
	// Build representation with URL instead of ID
	var resourcesWithURL []models.NamedResourceURL
	for _, r := range resources {
		resourcesWithURL = append(resourcesWithURL, r.ToNamedResourceURL(requestedBaseURL, resourceTypeName))
	}
	// Build the response JSON as a struct
	responseJSON := listResponse{len(resourcesWithURL), resourcesWithURL}
	// Transform the struct to JSON
	json, err := json.Marshal(responseJSON)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// Write the response
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

// AbilityListHandler handles requests on '/v1/abilities' and returns a list of all ability resources.
func AbilityListHandler(w http.ResponseWriter, r *http.Request) {
	// Fetch the ability list from the database
	abilities, err := db.GetAbilities()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// Build response JSON with URLs instead of IDs and send it to the client
	answerWithListJSON(abilities, r.Host, "abilities", w)
}

// CampListHandler handles requests on '/v1/camps' and returns a list of all camp resources.
func CampListHandler(w http.ResponseWriter, r *http.Request) {
	// Fetch the ability list from the database
	camps, err := db.GetCamps()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// Build response JSON with URLs instead of IDs and send it to the client
	answerWithListJSON(camps, r.Host, "camps", w)
}

// DungeonListHandler handles requests on '/v1/dungeons' and returns a list of all dungeon resources.
func DungeonListHandler(w http.ResponseWriter, r *http.Request) {
	// Fetch the ability list from the database
	dungeons, err := db.GetDungeons()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// Build response JSON with URLs instead of IDs and send it to the client
	answerWithListJSON(dungeons, r.Host, "dungeons", w)
}

// MoveListHandler handles requests on '/v1/moves' and returns a list of all move resources.
func MoveListHandler(w http.ResponseWriter, r *http.Request) {
	// Fetch the ability list from the database
	moves, err := db.GetMoves()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// Build response JSON with URLs instead of IDs and send it to the client
	answerWithListJSON(moves, r.Host, "moves", w)
}

// PokemonListHandler handles requests on '/v1/pokemon' and returns a list of all pokemon resources.
func PokemonListHandler(w http.ResponseWriter, r *http.Request) {
	// Fetch the ability list from the database
	pokemon, err := db.GetPokemon()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// Build response JSON with URLs instead of IDs and send it to the client
	answerWithListJSON(pokemon, r.Host, "pokemon", w)
}

// PokemonTypeListHandler handles requests on '/v1/types' and returns a list of all pokemon type resources.
func PokemonTypeListHandler(w http.ResponseWriter, r *http.Request) {
	// Fetch the ability list from the database
	pokemonTypes, err := db.GetPokemonTypes()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// Build response JSON with URLs instead of IDs and send it to the client
	answerWithListJSON(pokemonTypes, r.Host, "types", w)
}
