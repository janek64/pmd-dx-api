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

// ContextKey defines alls valid Context keys for requests to this API.
type ContextKey int

const (
	ResourceListParamsKey ContextKey = iota
)

// ResourceListParams contains the parsed parameter values for requests to resource lists.
type ResourceListParams struct {
	Sort db.SortInput
}

// answerWithListJSON transforms the provided resources to a list with URLs, packages
// them in a JSON and sends it as a response with the provided ResponseWriter.
func answerWithListJSON(resources []models.NamedResourceID, requestedBaseURL string, resourceTypeName string, w http.ResponseWriter) {
	// Build representation with URL instead of ID
	var resourcesWithURL []models.NamedResourceURL
	for _, r := range resources {
		resourcesWithURL = append(resourcesWithURL, r.ToNamedResourceURL(requestedBaseURL, resourceTypeName))
	}
	// Build the response JSON as a map
	responseJSON := map[string]interface{}{
		"count":   len(resourcesWithURL),
		"results": resourcesWithURL,
	}
	// Transform the map to JSON
	json, err := json.Marshal(responseJSON)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Write the response
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

// generateSearchInput decides if a db search argument is an ID or a name and generates the corresponding db.SearchInput.
func generateSearchInput(arg string) db.SearchInput {
	var searchInput db.SearchInput
	// Check if the search argument provided is an ID or a name
	// strconv.Atoi will return an error for non-numeric strings (name)
	if id, convErr := strconv.Atoi(arg); convErr == nil {
		searchInput.SearchType = db.ID
		searchInput.ID = id
	} else {
		searchInput.SearchType = db.Name
		// Convert to lowercase and then to unicode title case
		// Done on application level because SQL-level transformation disables indexes
		searchInput.Name = strings.Title(strings.ToLower(arg))
	}
	return searchInput
}

// transformToURLResources transforms a slice of NamedResources with IDs to NamedResources with URLs and returns it.
func transformToURLResources(resources []models.NamedResourceID, instanceURL string, resourceTypeName string) []models.NamedResourceURL {
	var resourcesWithURL []models.NamedResourceURL
	for _, p := range resources {
		resourcesWithURL = append(resourcesWithURL, p.ToNamedResourceURL(instanceURL, resourceTypeName))
	}
	return resourcesWithURL
}

// AbilityListHandler handles requests on '/v1/abilities' and returns a list of all ability resources.
func AbilityListHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Extract the ResourceListParams from the context with a type assertion
	params, ok := r.Context().Value(ResourceListParamsKey).(ResourceListParams)
	if !ok {
		http.Error(w, "Missing ResourceListParams", http.StatusInternalServerError)
		return
	}
	// Fetch the ability list from the database
	abilities, err := db.GetAbilityList(params.Sort)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Build response JSON with URLs instead of IDs and send it to the client
	answerWithListJSON(abilities, r.Host, "abilities", w)
}

// AbilitySearchHandler handles requests on '/v1/abilities/:searcharg' and returns information about the desired ability.
func AbilitySearchHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Generate the input for the db search
	searchInput := generateSearchInput(ps.ByName("searcharg"))
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
	pokemonWithURL := transformToURLResources(pokemon, r.Host, "pokemon")
	// Build the response JSON with a map
	responseJSON := map[string]interface{}{
		"id":          ability.AbilityID,
		"name":        ability.AbilityName,
		"description": ability.Description,
		"pokemon":     pokemonWithURL,
	}
	// Transform the map to JSON
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
	// Extract the ResourceListParams from the context with a type assertion
	params, ok := r.Context().Value(ResourceListParamsKey).(ResourceListParams)
	if !ok {
		http.Error(w, "Missing ResourceListParams", http.StatusInternalServerError)
		return
	}
	// Fetch the ability list from the database
	camps, err := db.GetCampList(params.Sort)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Build response JSON with URLs instead of IDs and send it to the client
	answerWithListJSON(camps, r.Host, "camps", w)
}

// CampSearchHandler handles requests on '/v1/camps/:searcharg' and returns information about the desired camp.
func CampSearchHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Generate the input for the db search
	searchInput := generateSearchInput(ps.ByName("searcharg"))
	// Get the ability from the database
	camp, pokemon, err := db.GetCamp(searchInput)
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
	pokemonWithURL := transformToURLResources(pokemon, r.Host, "pokemon")
	// Build the response JSON with a map
	responseJSON := map[string]interface{}{
		"id":          camp.CampID,
		"name":        camp.CampName,
		"description": camp.Description,
		"unlockType":  camp.UnlockType,
		"cost":        camp.Cost,
		"pokemon":     pokemonWithURL,
	}
	// Transform the map to JSON
	json, err := json.Marshal(responseJSON)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Write the response
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

// DungeonListHandler handles requests on '/v1/dungeons' and returns a list of all dungeon resources.
func DungeonListHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Extract the ResourceListParams from the context with a type assertion
	params, ok := r.Context().Value(ResourceListParamsKey).(ResourceListParams)
	if !ok {
		http.Error(w, "Missing ResourceListParams", http.StatusInternalServerError)
		return
	}
	// Fetch the ability list from the database
	dungeons, err := db.GetDungeonList(params.Sort)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Build response JSON with URLs instead of IDs and send it to the client
	answerWithListJSON(dungeons, r.Host, "dungeons", w)
}

// DungeonSearchHandler handles requests on '/v1/dungeons/:searcharg' and returns information about the desired dungeon.
func DungeonSearchHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Generate the input for the db search
	searchInput := generateSearchInput(ps.ByName("searcharg"))
	// Get the ability from the database
	dungeon, pokemon, err := db.GetDungeon(searchInput)
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
	var pokemonWithURL []models.DungeonPokemonURL
	for _, p := range pokemon {
		pokemonWithURL = append(pokemonWithURL, p.ToDungeonPokemonURL(r.Host))
	}
	// Build the response JSON with a map
	responseJSON := map[string]interface{}{
		"id":             dungeon.DungeonID,
		"name":           dungeon.DungeonName,
		"levels":         dungeon.Levels,
		"startLevel":     dungeon.StartLevel,
		"teamSize":       dungeon.TeamSize,
		"itemsAllowed":   dungeon.ItemsAllowed,
		"pokemonJoining": dungeon.PokemonJoining,
		"mapVisible":     dungeon.MapVisible,
		"pokemon":        pokemonWithURL,
	}
	// Transform the map to JSON
	json, err := json.Marshal(responseJSON)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Write the response
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

// MoveListHandler handles requests on '/v1/moves' and returns a list of all move resources.
func MoveListHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Extract the ResourceListParams from the context with a type assertion
	params, ok := r.Context().Value(ResourceListParamsKey).(ResourceListParams)
	if !ok {
		http.Error(w, "Missing ResourceListParams", http.StatusInternalServerError)
		return
	}
	// Fetch the ability list from the database
	moves, err := db.GetMoveList(params.Sort)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Build response JSON with URLs instead of IDs and send it to the client
	answerWithListJSON(moves, r.Host, "moves", w)
}

// MoveSearchHandler handles requests on '/v1/moves/:searcharg' and returns information about the desired move.
func MoveSearchHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Generate the input for the db search
	searchInput := generateSearchInput(ps.ByName("searcharg"))
	// Get the ability from the database
	move, moveType, pokemon, err := db.GetMove(searchInput)
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
	var pokemonWithURL []models.MovePokemonURL
	for _, p := range pokemon {
		pokemonWithURL = append(pokemonWithURL, p.ToMovePokemonURL(r.Host))
	}
	// Build the response JSON with a map
	responseJSON := map[string]interface{}{
		"id":           move.MoveID,
		"name":         move.MoveName,
		"category":     move.Category,
		"range":        move.Range,
		"target":       move.Target,
		"initialPP":    move.InitialPP,
		"initialPower": move.InitialPower,
		"accuracy":     move.Accuracy,
		"description":  move.Description,
		"type":         moveType.ToNamedResourceURL(r.Host, "moves"),
		"pokemon":      pokemonWithURL,
	}
	// Transform the map to JSON
	json, err := json.Marshal(responseJSON)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Write the response
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

// PokemonListHandler handles requests on '/v1/pokemon' and returns a list of all pokemon resources.
func PokemonListHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Extract the ResourceListParams from the context with a type assertion
	params, ok := r.Context().Value(ResourceListParamsKey).(ResourceListParams)
	if !ok {
		http.Error(w, "Missing ResourceListParams", http.StatusInternalServerError)
		return
	}
	// Fetch the ability list from the database
	pokemon, err := db.GetPokemonList(params.Sort)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Build response JSON with URLs instead of IDs and send it to the client
	answerWithListJSON(pokemon, r.Host, "pokemon", w)
}

// PokemonSearchHandler handles requests on '/v1/pokemon/:searcharg' and returns information about the desired pokemon.
func PokemonSearchHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Generate the input for the db search
	searchInput := generateSearchInput(ps.ByName("searcharg"))
	// Get the ability from the database
	pokemon, camp, abilities, dungeons, moves, pokemonTypes, err := db.GetPokemon(searchInput)
	if err != nil {
		// If the error is a db.ResourceNotFoundError, return code 404 (not found)
		if _, ok := err.(*db.ResourceNotFoundError); ok {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	// Build representation of the abilities with URL instead of ID
	abilitiesWithURL := transformToURLResources(abilities, r.Host, "abilities")
	// Build representation of the dungeons with URL instead of ID
	var dungeonsWithURL []models.PokemonDungeonURL
	for _, d := range dungeons {
		dungeonsWithURL = append(dungeonsWithURL, d.ToPokemonDungeonURL(r.Host))
	}
	// Build representation of the moves with URL instead of ID
	var movesWithURL []models.PokemonMoveURL
	for _, m := range moves {
		movesWithURL = append(movesWithURL, m.ToPokemonMoveURL(r.Host))
	}
	// Build representation of the types with URL instead of ID
	pokemonTypesWithURL := transformToURLResources(pokemonTypes, r.Host, "types")
	// Build the response JSON with a map
	responseJSON := map[string]interface{}{
		"id":              pokemon.DexNumber,
		"name":            pokemon.PokemonName,
		"classification":  pokemon.Classification,
		"evolutionStage":  pokemon.EvolutionStage,
		"evolveCondition": pokemon.EvolveCondition,
		"evolveLevel":     pokemon.EvolveLevel,
		"evolveCrystals":  pokemon.EvolveCrystals,
		"camp":            camp.ToNamedResourceURL(r.Host, "camps"),
		"abilities":       abilitiesWithURL,
		"dungeons":        dungeonsWithURL,
		"moves":           movesWithURL,
		"types":           pokemonTypesWithURL,
	}
	// Transform the map to JSON
	json, err := json.Marshal(responseJSON)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Write the response
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

// PokemonTypeListHandler handles requests on '/v1/types' and returns a list of all pokemon type resources.
func PokemonTypeListHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Extract the ResourceListParams from the context with a type assertion
	params, ok := r.Context().Value(ResourceListParamsKey).(ResourceListParams)
	if !ok {
		http.Error(w, "Missing ResourceListParams", http.StatusInternalServerError)
		return
	}
	// Fetch the ability list from the database
	pokemonTypes, err := db.GetPokemonTypeList(params.Sort)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Build response JSON with URLs instead of IDs and send it to the client
	answerWithListJSON(pokemonTypes, r.Host, "types", w)
}

// PokemonTypeSearchHandler handles requests on '/v1/types/:searcharg' and returns information about the desired pokemonType.
func PokemonTypeSearchHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Generate the input for the db search
	searchInput := generateSearchInput(ps.ByName("searcharg"))
	// Get the ability from the database
	pokemonType, interactions, err := db.GetPokemonType(searchInput)
	if err != nil {
		// If the error is a db.ResourceNotFoundError, return code 404 (not found)
		if _, ok := err.(*db.ResourceNotFoundError); ok {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	// Build representation of the interactions with URL instead of ID
	var interactionsWithURL []models.TypeInteractionURL
	for _, i := range interactions {
		interactionsWithURL = append(interactionsWithURL, i.ToTypeInteractionURL(r.Host))
	}
	// Build the response JSON with a map
	responseJSON := map[string]interface{}{
		"id":           pokemonType.TypeID,
		"name":         pokemonType.TypeName,
		"interactions": interactionsWithURL,
	}
	// Transform the map to JSON
	json, err := json.Marshal(responseJSON)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Write the response
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}
