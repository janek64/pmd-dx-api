// Package handler defines all HTTP handler functions
// for routes of the pmd-dx-api.
package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/iancoleman/orderedmap"
	"github.com/janek64/pmd-dx-api/api/db"
	"github.com/janek64/pmd-dx-api/api/logger"
	"github.com/janek64/pmd-dx-api/api/models"
	"github.com/julienschmidt/httprouter"
)

// ContextKey defines alls valid Context keys for requests to this API.
type ContextKey int

const (
	ResourceListParamsKey ContextKey = iota
	FieldLimitingParamsKey
)

// ResourceListParams contains the parsed parameter values for requests to resource lists.
type ResourceListParams struct {
	Sort       db.SortInput
	Pagination db.Pagination
}

// FieldLimitingParams contains the parsed parameter values for requests to resource lists.
type FieldLimitingParams struct {
	FieldLimitingEnabled bool
	Fields               []string
}

// ErrorAndLog500 is a wrapper around http.Error() that
// writes the error message to the error log instead of returning
// it to the client. Should only be used for internal server errors.
func ErrorAndLog500(w http.ResponseWriter, err error) {
	// Use http.Error() with default message
	http.Error(w, "Something went wrong on our side. Please contact the administrator.", http.StatusInternalServerError)
	// Gather caller information to pass it to the logger
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		fmt.Fprintf(os.Stderr, "ErrorAndLog500: failed to fetch caller information")
		return
	}
	caller := logger.CallerInformation{Pc: pc, File: file, Line: line}
	// Write to the error logger
	logErr := logger.LogError(err, caller)
	if logErr != nil {
		fmt.Fprintf(os.Stderr, "Writing to the error log failed: %v", err)
	}
}

// answerWithListJSON transforms the provided resources to a list with URLs, packages
// them in a JSON and sends it as a response with the provided ResponseWriter.
func answerWithListJSON(count int, resources []models.NamedResourceID, resourceTypeName string, pagination db.Pagination, w http.ResponseWriter, r *http.Request) {
	// Build representation with URL instead of ID
	var resourcesWithURL []models.NamedResourceURL
	for _, resource := range resources {
		resourcesWithURL = append(resourcesWithURL, resource.ToNamedResourceURL(r.Host, resourceTypeName))
	}
	// Build the response JSON as a map
	responseJSON := orderedmap.New()
	responseJSON.Set("count", count)
	responseJSON.Set("results", resourcesWithURL)
	// Extract the FieldLimitingParams from the context with a type assertion
	fieldLimitParams, ok := r.Context().Value(FieldLimitingParamsKey).(FieldLimitingParams)
	if !ok {
		ErrorAndLog500(w, errors.New("missing FieldLimitingParams"))
		return
	}
	// Perform field limiting if necessary
	limitResultFields(responseJSON, fieldLimitParams)
	// Transform the map to JSON
	json, err := json.Marshal(responseJSON)
	if err != nil {
		ErrorAndLog500(w, err)
		return
	}
	// Generate the headers for pagination
	// Calculate the page numbers
	lastPage := count/pagination.PerPage + 1
	if count%pagination.PerPage == 0 {
		lastPage -= 1
	}
	nextPage := pagination.Page + 1
	previousPage := pagination.Page - 1
	// Generate the URLs
	requestURL := r.Host + r.URL.String()
	// If no page URL parameter was provided, add it
	match, err := regexp.Match(`.+[?&]page=\d*(&.+)?`, []byte(requestURL))
	if err != nil {
		ErrorAndLog500(w, err)
		return
	}
	if !match {
		// Check if there is already a question mark followed by characters
		match, err = regexp.Match(`.+\?.+`, []byte(requestURL))
		if err != nil {
			ErrorAndLog500(w, err)
			return
		}
		if match {
			requestURL = fmt.Sprintf("%v&page=%v", requestURL, pagination.Page)
		} else {
			requestURL = fmt.Sprintf("%v?page=%v", requestURL, pagination.Page)
		}

	}
	re, err := regexp.Compile(`([?&])page=\d*`)
	if err != nil {
		ErrorAndLog500(w, err)
		return
	}
	nextURL := re.ReplaceAllString(requestURL, fmt.Sprintf("${1}page=%v", nextPage))
	previousURL := re.ReplaceAllString(requestURL, fmt.Sprintf("${1}page=%v", previousPage))
	lastURL := re.ReplaceAllString(requestURL, fmt.Sprintf("${1}page=%v", lastPage))
	// Set null values when links should not be provided
	if pagination.Page == 1 {
		previousURL = "null"
	}
	if pagination.Page == lastPage {
		nextURL = "null"
	} else if pagination.Page > lastPage {
		nextURL = "null"
		previousURL = "null"
	}
	// Set the Link header
	linkHeader := fmt.Sprintf("<%v>; rel=\"next\", <%v>; rel=\"previous\", <%v>; rel=\"last\"", nextURL, previousURL, lastURL)
	w.Header().Set("Link", linkHeader)
	// Write the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
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

// limitResultFields checks if field limiting is necessary and removes all fields
// from the responseJSON that should not be displayed if this is the case.
func limitResultFields(responseJSON *orderedmap.OrderedMap, params FieldLimitingParams) {
	// Check if field limiting is not enabled
	if !params.FieldLimitingEnabled {
		return
	}
	// Loop through the JSON and check which parameters need to be removed
	deleteKeys := make(map[string]bool)
	keys := responseJSON.Keys()
	for _, k := range keys {
		deleteKeys[k] = true
		for _, v := range params.Fields {
			if k == v {
				deleteKeys[k] = false
				break
			}
		}
	}
	// Delete all keys marked for deletion
	// Needs to be done separately since deleting while looping over the keys
	// caused keys to be skipped and others to be used multiple times
	for k, v := range deleteKeys {
		if v {
			responseJSON.Delete(k)
		}
	}
}

// AbilityListHandler handles requests on '/v1/abilities' and returns a list of all ability resources.
func AbilityListHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Extract the ResourceListParams from the context with a type assertion
	params, ok := r.Context().Value(ResourceListParamsKey).(ResourceListParams)
	if !ok {
		ErrorAndLog500(w, errors.New("missing ResourceListParams"))
		return
	}
	// Fetch the ability list from the database
	count, abilities, err := db.GetAbilityList(params.Sort, params.Pagination)
	if err != nil {
		ErrorAndLog500(w, err)
		return
	}
	// Build response JSON with URLs instead of IDs and send it to the client
	answerWithListJSON(count, abilities, "abilities", params.Pagination, w, r)
}

// AbilitySearchHandler handles requests on '/v1/abilities/:searcharg' and returns information about the desired ability.
func AbilitySearchHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Extract the FieldLimitingParams from the context with a type assertion
	fieldLimitParams, ok := r.Context().Value(FieldLimitingParamsKey).(FieldLimitingParams)
	if !ok {
		ErrorAndLog500(w, errors.New("missing FieldLimitingParams"))
		return
	}
	// Generate the input for the db search
	searchInput := generateSearchInput(ps.ByName("searcharg"))
	// Get the ability from the database
	ability, pokemon, err := db.GetAbility(searchInput)
	if err != nil {
		// If the error is a db.ResourceNotFoundError, return code 404 (not found)
		if _, ok := err.(*db.ResourceNotFoundError); ok {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			ErrorAndLog500(w, err)
		}
		return
	}
	// Build representation of the pokemon with URL instead of ID
	pokemonWithURL := transformToURLResources(pokemon, r.Host, "pokemon")
	// Build the response JSON with a map
	responseJSON := orderedmap.New()
	responseJSON.Set("id", ability.AbilityID)
	responseJSON.Set("name", ability.AbilityName)
	responseJSON.Set("description", ability.Description)
	responseJSON.Set("pokemon", pokemonWithURL)
	// Perform field limiting if necessary
	limitResultFields(responseJSON, fieldLimitParams)
	// Transform the map to JSON
	json, err := json.Marshal(responseJSON)
	if err != nil {
		ErrorAndLog500(w, err)
		return
	}
	// Write the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

// CampListHandler handles requests on '/v1/camps' and returns a list of all camp resources.
func CampListHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Extract the ResourceListParams from the context with a type assertion
	params, ok := r.Context().Value(ResourceListParamsKey).(ResourceListParams)
	if !ok {
		ErrorAndLog500(w, errors.New("missing ResourceListParams"))
		return
	}
	// Fetch the ability list from the database
	count, camps, err := db.GetCampList(params.Sort, params.Pagination)
	if err != nil {
		ErrorAndLog500(w, err)
		return
	}
	// Build response JSON with URLs instead of IDs and send it to the client
	answerWithListJSON(count, camps, "camps", params.Pagination, w, r)
}

// CampSearchHandler handles requests on '/v1/camps/:searcharg' and returns information about the desired camp.
func CampSearchHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Extract the FieldLimitingParams from the context with a type assertion
	fieldLimitParams, ok := r.Context().Value(FieldLimitingParamsKey).(FieldLimitingParams)
	if !ok {
		ErrorAndLog500(w, errors.New("missing FieldLimitingParams"))
		return
	}
	// Generate the input for the db search
	searchInput := generateSearchInput(ps.ByName("searcharg"))
	// Get the ability from the database
	camp, pokemon, err := db.GetCamp(searchInput)
	if err != nil {
		// If the error is a db.ResourceNotFoundError, return code 404 (not found)
		if _, ok := err.(*db.ResourceNotFoundError); ok {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			ErrorAndLog500(w, err)
		}
		return
	}
	// Build representation of the pokemon with URL instead of ID
	pokemonWithURL := transformToURLResources(pokemon, r.Host, "pokemon")
	// Build the response JSON with a map
	responseJSON := orderedmap.New()
	responseJSON.Set("id", camp.CampID)
	responseJSON.Set("name", camp.CampName)
	responseJSON.Set("description", camp.Description)
	responseJSON.Set("unlockType", camp.UnlockType)
	responseJSON.Set("cost", camp.Cost)
	responseJSON.Set("pokemon", pokemonWithURL)
	// Perform field limiting if necessary
	limitResultFields(responseJSON, fieldLimitParams)
	// Transform the map to JSON
	json, err := json.Marshal(responseJSON)
	if err != nil {
		ErrorAndLog500(w, err)
		return
	}
	// Write the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

// DungeonListHandler handles requests on '/v1/dungeons' and returns a list of all dungeon resources.
func DungeonListHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Extract the ResourceListParams from the context with a type assertion
	params, ok := r.Context().Value(ResourceListParamsKey).(ResourceListParams)
	if !ok {
		ErrorAndLog500(w, errors.New("missing ResourceListParams"))
		return
	}
	// Fetch the ability list from the database
	count, dungeons, err := db.GetDungeonList(params.Sort, params.Pagination)
	if err != nil {
		ErrorAndLog500(w, err)
		return
	}
	// Build response JSON with URLs instead of IDs and send it to the client
	answerWithListJSON(count, dungeons, "dungeons", params.Pagination, w, r)
}

// DungeonSearchHandler handles requests on '/v1/dungeons/:searcharg' and returns information about the desired dungeon.
func DungeonSearchHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Extract the FieldLimitingParams from the context with a type assertion
	fieldLimitParams, ok := r.Context().Value(FieldLimitingParamsKey).(FieldLimitingParams)
	if !ok {
		ErrorAndLog500(w, errors.New("missing FieldLimitingParams"))
		return
	}
	// Generate the input for the db search
	searchInput := generateSearchInput(ps.ByName("searcharg"))
	// Get the ability from the database
	dungeon, pokemon, err := db.GetDungeon(searchInput)
	if err != nil {
		// If the error is a db.ResourceNotFoundError, return code 404 (not found)
		if _, ok := err.(*db.ResourceNotFoundError); ok {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			ErrorAndLog500(w, err)
		}
		return
	}
	// Build representation of the pokemon with URL instead of ID
	var pokemonWithURL []models.DungeonPokemonURL
	for _, p := range pokemon {
		pokemonWithURL = append(pokemonWithURL, p.ToDungeonPokemonURL(r.Host))
	}
	// Build the response JSON with a map
	responseJSON := orderedmap.New()
	responseJSON.Set("id", dungeon.DungeonID)
	responseJSON.Set("name", dungeon.DungeonName)
	responseJSON.Set("levels", dungeon.Levels)
	responseJSON.Set("startLevel", dungeon.StartLevel)
	responseJSON.Set("teamSize", dungeon.TeamSize)
	responseJSON.Set("itemsAllowed", dungeon.ItemsAllowed)
	responseJSON.Set("pokemonJoining", dungeon.PokemonJoining)
	responseJSON.Set("mapVisible", dungeon.MapVisible)
	responseJSON.Set("pokemon", pokemonWithURL)
	// Perform field limiting if necessary
	limitResultFields(responseJSON, fieldLimitParams)
	// Transform the map to JSON
	json, err := json.Marshal(responseJSON)
	if err != nil {
		ErrorAndLog500(w, err)
		return
	}
	// Write the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

// MoveListHandler handles requests on '/v1/moves' and returns a list of all move resources.
func MoveListHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Extract the ResourceListParams from the context with a type assertion
	params, ok := r.Context().Value(ResourceListParamsKey).(ResourceListParams)
	if !ok {
		ErrorAndLog500(w, errors.New("missing ResourceListParams"))
		return
	}
	// Fetch the ability list from the database
	count, moves, err := db.GetMoveList(params.Sort, params.Pagination)
	if err != nil {
		ErrorAndLog500(w, err)
		return
	}
	// Build response JSON with URLs instead of IDs and send it to the client
	answerWithListJSON(count, moves, "moves", params.Pagination, w, r)
}

// MoveSearchHandler handles requests on '/v1/moves/:searcharg' and returns information about the desired move.
func MoveSearchHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Extract the FieldLimitingParams from the context with a type assertion
	fieldLimitParams, ok := r.Context().Value(FieldLimitingParamsKey).(FieldLimitingParams)
	if !ok {
		ErrorAndLog500(w, errors.New("missing FieldLimitingParams"))
		return
	}
	// Generate the input for the db search
	searchInput := generateSearchInput(ps.ByName("searcharg"))
	// Get the ability from the database
	move, moveType, pokemon, err := db.GetMove(searchInput)
	if err != nil {
		// If the error is a db.ResourceNotFoundError, return code 404 (not found)
		if _, ok := err.(*db.ResourceNotFoundError); ok {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			ErrorAndLog500(w, err)
		}
		return
	}
	// Build representation of the pokemon with URL instead of ID
	var pokemonWithURL []models.MovePokemonURL
	for _, p := range pokemon {
		pokemonWithURL = append(pokemonWithURL, p.ToMovePokemonURL(r.Host))
	}
	// Build the response JSON with a map
	responseJSON := orderedmap.New()
	responseJSON.Set("id", move.MoveID)
	responseJSON.Set("name", move.MoveName)
	responseJSON.Set("category", move.Category)
	responseJSON.Set("range", move.Range)
	responseJSON.Set("target", move.Target)
	responseJSON.Set("initialPP", move.InitialPP)
	responseJSON.Set("initialPower", move.InitialPower)
	responseJSON.Set("accuracy", move.Accuracy)
	responseJSON.Set("description", move.Description)
	responseJSON.Set("type", moveType.ToNamedResourceURL(r.Host, "moves"))
	responseJSON.Set("pokemon", pokemonWithURL)
	// Perform field limiting if necessary
	limitResultFields(responseJSON, fieldLimitParams)
	// Transform the map to JSON
	json, err := json.Marshal(responseJSON)
	if err != nil {
		ErrorAndLog500(w, err)
		return
	}
	// Write the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

// PokemonListHandler handles requests on '/v1/pokemon' and returns a list of all pokemon resources.
func PokemonListHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Extract the ResourceListParams from the context with a type assertion
	params, ok := r.Context().Value(ResourceListParamsKey).(ResourceListParams)
	if !ok {
		ErrorAndLog500(w, errors.New("missing ResourceListParams"))
		return
	}
	// Fetch the ability list from the database
	count, pokemon, err := db.GetPokemonList(params.Sort, params.Pagination)
	if err != nil {
		ErrorAndLog500(w, err)
		return
	}
	// Build response JSON with URLs instead of IDs and send it to the client
	answerWithListJSON(count, pokemon, "pokemon", params.Pagination, w, r)
}

// PokemonSearchHandler handles requests on '/v1/pokemon/:searcharg' and returns information about the desired pokemon.
func PokemonSearchHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Extract the FieldLimitingParams from the context with a type assertion
	fieldLimitParams, ok := r.Context().Value(FieldLimitingParamsKey).(FieldLimitingParams)
	if !ok {
		ErrorAndLog500(w, errors.New("missing FieldLimitingParams"))
		return
	}
	// Generate the input for the db search
	searchInput := generateSearchInput(ps.ByName("searcharg"))
	// Get the ability from the database
	pokemon, camp, abilities, dungeons, moves, pokemonTypes, err := db.GetPokemon(searchInput)
	if err != nil {
		// If the error is a db.ResourceNotFoundError, return code 404 (not found)
		if _, ok := err.(*db.ResourceNotFoundError); ok {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			ErrorAndLog500(w, err)
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
	responseJSON := orderedmap.New()
	responseJSON.Set("id", pokemon.DexNumber)
	responseJSON.Set("name", pokemon.PokemonName)
	responseJSON.Set("classification", pokemon.Classification)
	responseJSON.Set("evolutionStage", pokemon.EvolutionStage)
	responseJSON.Set("evolveCondition", pokemon.EvolveCondition)
	responseJSON.Set("evolveLevel", pokemon.EvolveLevel)
	responseJSON.Set("evolveCrystals", pokemon.EvolveCrystals)
	responseJSON.Set("camp", camp.ToNamedResourceURL(r.Host, "camps"))
	responseJSON.Set("abilities", abilitiesWithURL)
	responseJSON.Set("dungeons", dungeonsWithURL)
	responseJSON.Set("moves", movesWithURL)
	responseJSON.Set("types", pokemonTypesWithURL)
	// Perform field limiting if necessary
	limitResultFields(responseJSON, fieldLimitParams)
	// Transform the map to JSON
	json, err := json.Marshal(responseJSON)
	if err != nil {
		ErrorAndLog500(w, err)
		return
	}
	// Write the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

// PokemonTypeListHandler handles requests on '/v1/types' and returns a list of all pokemon type resources.
func PokemonTypeListHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Extract the ResourceListParams from the context with a type assertion
	params, ok := r.Context().Value(ResourceListParamsKey).(ResourceListParams)
	if !ok {
		ErrorAndLog500(w, errors.New("missing ResourceListParams"))
		return
	}
	// Fetch the ability list from the database
	count, pokemonTypes, err := db.GetPokemonTypeList(params.Sort, params.Pagination)
	if err != nil {
		ErrorAndLog500(w, err)
		return
	}
	// Build response JSON with URLs instead of IDs and send it to the client
	answerWithListJSON(count, pokemonTypes, "types", params.Pagination, w, r)
}

// PokemonTypeSearchHandler handles requests on '/v1/types/:searcharg' and returns information about the desired pokemonType.
func PokemonTypeSearchHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Extract the FieldLimitingParams from the context with a type assertion
	fieldLimitParams, ok := r.Context().Value(FieldLimitingParamsKey).(FieldLimitingParams)
	if !ok {
		ErrorAndLog500(w, errors.New("missing FieldLimitingParams"))
		return
	}
	// Generate the input for the db search
	searchInput := generateSearchInput(ps.ByName("searcharg"))
	// Get the ability from the database
	pokemonType, interactions, err := db.GetPokemonType(searchInput)
	if err != nil {
		// If the error is a db.ResourceNotFoundError, return code 404 (not found)
		if _, ok := err.(*db.ResourceNotFoundError); ok {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			ErrorAndLog500(w, err)
		}
		return
	}
	// Build representation of the interactions with URL instead of ID
	var interactionsWithURL []models.TypeInteractionURL
	for _, i := range interactions {
		interactionsWithURL = append(interactionsWithURL, i.ToTypeInteractionURL(r.Host))
	}
	// Build the response JSON with a map
	responseJSON := orderedmap.New()
	responseJSON.Set("id", pokemonType.TypeID)
	responseJSON.Set("name", pokemonType.TypeName)
	responseJSON.Set("interactions", interactionsWithURL)
	// Perform field limiting if necessary
	limitResultFields(responseJSON, fieldLimitParams)
	// Transform the map to JSON
	json, err := json.Marshal(responseJSON)
	if err != nil {
		ErrorAndLog500(w, err)
		return
	}
	// Write the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}
