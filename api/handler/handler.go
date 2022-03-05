// Package handler defines all HTTP handler functions
// for routes of the pmd-dx-api.
package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/janek64/pmd-dx-api/api/db"
	"github.com/janek64/pmd-dx-api/api/models"
	"github.com/julienschmidt/httprouter"
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
		return
	}
	// Write the response
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

// AbilityListHandler handles requests on '/v1/abilities' and returns a list of all ability resources.
func AbilityListHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Fetch the ability list from the database
	abilities, err := db.GetAbilities()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Build response JSON with URLs instead of IDs and send it to the client
	answerWithListJSON(abilities, r.Host, "abilities", w)
}

// AbilityListHandler handles requests on '/v1/abilities/:searcharg' and returns information about the resource.
func AbilitySearchHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var searchInput db.SearchInput
	// Check if the search argument provided is an ID or a name
	arg := ps.ByName("searcharg")
	// strconv.Atoi will return an error for non-numeric strings (name)
	if id, convErr := strconv.Atoi(arg); convErr == nil {
		searchInput.SearchType = db.ID
		searchInput.ID = id
	} else {
		searchInput.SearchType = db.NAME
		// Convert to lowercase and then to unicode title case
		// Done on application level because SQL-level transformation disables indexes
		searchInput.Name = strings.Title(strings.ToLower(arg))
	}
	// Get the ability from the database
	ability, pokemon, err := db.GetAbility(searchInput)
	if err != nil {
		// If the error is a db.ResourceNotFoundError, return code 404 (not found)
		if _, ok := err.(*db.ResourceNotFoundError); ok {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	// Build representation of the pokemon with URL instead of ID
	var pokemonWithURL []models.NamedResourceURL
	for _, p := range pokemon {
		pokemonWithURL = append(pokemonWithURL, p.ToNamedResourceURL(r.Host, "abilities"))
	}
	// Build the response JSON with an anonymous struct
	responseJSON := struct {
		ID          int                       `json:"id"`
		Name        string                    `json:"name"`
		Description string                    `json:"description"`
		Pokemon     []models.NamedResourceURL `json:"pokemon"`
	}{
		ID:          ability.AbilityID,
		Name:        ability.AbilityName,
		Description: ability.Description,
		Pokemon:     pokemonWithURL,
	}
	// Transform the struct to JSON
	json, err := json.Marshal(responseJSON)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Write the response
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

// CampListHandler handles requests on '/v1/camps' and returns a list of all camp resources.
func CampListHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Fetch the ability list from the database
	camps, err := db.GetCamps()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Build response JSON with URLs instead of IDs and send it to the client
	answerWithListJSON(camps, r.Host, "camps", w)
}

// DungeonListHandler handles requests on '/v1/dungeons' and returns a list of all dungeon resources.
func DungeonListHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Fetch the ability list from the database
	dungeons, err := db.GetDungeons()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Build response JSON with URLs instead of IDs and send it to the client
	answerWithListJSON(dungeons, r.Host, "dungeons", w)
}

// MoveListHandler handles requests on '/v1/moves' and returns a list of all move resources.
func MoveListHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Fetch the ability list from the database
	moves, err := db.GetMoves()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Build response JSON with URLs instead of IDs and send it to the client
	answerWithListJSON(moves, r.Host, "moves", w)
}

// PokemonListHandler handles requests on '/v1/pokemon' and returns a list of all pokemon resources.
func PokemonListHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Fetch the ability list from the database
	pokemon, err := db.GetPokemon()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Build response JSON with URLs instead of IDs and send it to the client
	answerWithListJSON(pokemon, r.Host, "pokemon", w)
}

// PokemonTypeListHandler handles requests on '/v1/types' and returns a list of all pokemon type resources.
func PokemonTypeListHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Fetch the ability list from the database
	pokemonTypes, err := db.GetPokemonTypes()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Build response JSON with URLs instead of IDs and send it to the client
	answerWithListJSON(pokemonTypes, r.Host, "types", w)
}
