package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/janek64/pmd-dx-api/api/models"
)

// SearchType represents the valid search types for resources (ID and name)
type SearchType string

const (
	ID   = "ID"
	NAME = "name"
)

// SearchInput is the input for a query searching a resource by ID or name
type SearchInput struct {
	SearchType SearchType
	ID         int
	Name       string
}

// ResourceNotFoundError - error if a requested resource was not found
type ResourceNotFoundError struct {
	ResourceType string
	SearchType   SearchType
	ID           int
	Name         string
}

// Error - implementation of the error interface.
func (e *ResourceNotFoundError) Error() string {
	if e.SearchType == ID {
		return fmt.Sprintf("resource of type '%v' with %v '%v' not found", e.ResourceType, e.SearchType, e.ID)
	} else if e.SearchType == NAME {
		return fmt.Sprintf("resource of type '%v' with %v '%v' not found", e.ResourceType, e.SearchType, e.Name)
	} else {
		return "resource not found"
	}
}

// GetAbilities fetches a slice of all ability entries from the database.
func GetAbilities() ([]models.NamedResourceID, error) {
	var abilities []models.NamedResourceID
	rows, err := dbpool.Query(context.Background(), "SELECT ability_ID, ability_name FROM ability;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// Add all abilities found to the slice
	for rows.Next() {
		var ability models.NamedResourceID
		err = rows.Scan(&ability.ID, &ability.Name)
		if err != nil {
			return nil, err
		}
		abilities = append(abilities, ability)
	}
	return abilities, nil
}

// GetAbility fetches an ability entry and all pokemon learning it from the database by its ID or name.
func GetAbility(input SearchInput) (ability models.Ability, pokemon []models.NamedResourceID, err error) {
	var rows pgx.Rows
	// Use different query depending on search type
	if input.SearchType == ID {
		queryString := `SELECT A.*, P.dex_number, P.pokemon_name
		FROM (SELECT * FROM ability WHERE ability_ID = $1) A
		LEFT JOIN pokemon_has_ability PA ON A.ability_ID = PA.ability_ID
		LEFT JOIN pokemon P on PA.dex_number = P.dex_number;`
		rows, err = dbpool.Query(context.Background(), queryString, input.ID)
	} else if input.SearchType == NAME {
		queryString := `SELECT A.*, P.dex_number, P.pokemon_name
		FROM (SELECT * FROM ability WHERE ability_name = $1) A
		LEFT JOIN pokemon_has_ability PA ON A.ability_ID = PA.ability_ID
		LEFT JOIN pokemon P on PA.dex_number = P.dex_number;`
		rows, err = dbpool.Query(context.Background(), queryString, input.Name)
	} else {
		return ability, nil, fmt.Errorf("illegal search type %v", input.SearchType)
	}
	if err != nil {
		return ability, nil, err
	}
	defer rows.Close()
	var p models.NamedResourceID
	// Read the first row outside of the loop to extract ability information and check for null pokemon
	rows.Next()
	err = rows.Scan(&ability.AbilityID, &ability.AbilityName, &ability.Description, &p.ID, &p.Name)
	// Add the pokemon to the slice
	// Check if the pokemon is not null to find ability without pokemon
	if p.ID != 0 {
		pokemon = append(pokemon, p)
	}
	// Add all other pokemon to the slice
	for rows.Next() {
		var empty [3]interface{}
		err = rows.Scan(&empty[0], &empty[1], &empty[2], &p.ID, &p.Name)
		if err != nil {
			return ability, nil, err
		}
		// Checking for ID==0 is not necessary since all rows after the first will not have null values
		pokemon = append(pokemon, p)
	}
	// If the AbilityID is zero, no entry was found
	if ability.AbilityID == 0 {
		if input.SearchType == ID {
			return ability, nil, &ResourceNotFoundError{ResourceType: "ability", SearchType: input.SearchType, ID: input.ID}
		} else if input.SearchType == NAME {
			return ability, nil, &ResourceNotFoundError{ResourceType: "ability", SearchType: input.SearchType, Name: input.Name}
		}
	}
	return ability, pokemon, nil
}

// GetCamps fetches a slice of all camp entries from the database.
func GetCamps() ([]models.NamedResourceID, error) {
	var camps []models.NamedResourceID
	// SELECT camp_ID, camp_name FROM camp;
	rows, err := dbpool.Query(context.Background(), "SELECT camp_ID, camp_name FROM camp;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// Add all camps found to the slice
	for rows.Next() {
		var camp models.NamedResourceID
		err = rows.Scan(&camp.ID, &camp.Name)
		if err != nil {
			return nil, err
		}
		camps = append(camps, camp)
	}
	return camps, nil
}

// GetDungeons fetches a slice of all dungeon entries from the database.
func GetDungeons() ([]models.NamedResourceID, error) {
	var dungeons []models.NamedResourceID
	// SELECT dungeon_ID, dungeon_name FROM dungeon;
	rows, err := dbpool.Query(context.Background(), "SELECT dungeon_ID, dungeon_name FROM dungeon;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// Add all dungeons found to the slice
	for rows.Next() {
		var dungeon models.NamedResourceID
		err = rows.Scan(&dungeon.ID, &dungeon.Name)
		if err != nil {
			return nil, err
		}
		dungeons = append(dungeons, dungeon)
	}
	return dungeons, nil
}

// GetMoves fetches a slice of all attack_move entries from the database.
func GetMoves() ([]models.NamedResourceID, error) {
	var moves []models.NamedResourceID
	// SELECT move_ID, move_name FROM attack_move;
	rows, err := dbpool.Query(context.Background(), "SELECT move_ID, move_name FROM attack_move;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// Add all moves found to the slice
	for rows.Next() {
		var move models.NamedResourceID
		err = rows.Scan(&move.ID, &move.Name)
		if err != nil {
			return nil, err
		}
		moves = append(moves, move)
	}
	return moves, nil
}

// GetPokemon fetches a slice of all pokemon entries from the database.
func GetPokemon() ([]models.NamedResourceID, error) {
	var pokemonList []models.NamedResourceID
	// SELECT dex_number, pokemon_name FROM pokemon;
	rows, err := dbpool.Query(context.Background(), "SELECT dex_number, pokemon_name FROM pokemon;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// Add all pokemon found to the slice
	for rows.Next() {
		var pokemon models.NamedResourceID
		err = rows.Scan(&pokemon.ID, &pokemon.Name)
		if err != nil {
			return nil, err
		}
		pokemonList = append(pokemonList, pokemon)
	}
	return pokemonList, nil
}

// GetPokemonTypes fetches a slice of all pokemon_type entries from the database.
func GetPokemonTypes() ([]models.NamedResourceID, error) {
	var pokemonTypes []models.NamedResourceID
	// SELECT * FROM pokemon_type;
	rows, err := dbpool.Query(context.Background(), "SELECT * FROM pokemon_type;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// Add all types found to the slice
	for rows.Next() {
		var pokemonType models.NamedResourceID
		err = rows.Scan(&pokemonType.ID, &pokemonType.Name)
		if err != nil {
			return nil, err
		}
		pokemonTypes = append(pokemonTypes, pokemonType)
	}
	return pokemonTypes, nil
}
