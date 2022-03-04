package db

import (
	"context"

	"github.com/janek64/pmd-dx-api/api/models"
)

// GetAbilities fetches a slice of all ability entries from the database.
func GetAbilities() ([]models.NamedResourceID, error) {
	var abilities []models.NamedResourceID
	//SELECT ability_ID, ability_name FROM ability;
	rows, err := dbpool.Query(context.Background(), "SELECT ability_ID, ability_name FROM ability;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
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

// GetCamps fetches a slice of all camp entries from the database.
func GetCamps() ([]models.NamedResourceID, error) {
	var camps []models.NamedResourceID
	// SELECT camp_ID, camp_name FROM camp;
	rows, err := dbpool.Query(context.Background(), "SELECT camp_ID, camp_name FROM camp")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
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
	rows, err := dbpool.Query(context.Background(), "SELECT dungeon_ID, dungeon_name FROM dungeon")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
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
	rows, err := dbpool.Query(context.Background(), "SELECT move_ID, move_name FROM attack_move")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
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
	rows, err := dbpool.Query(context.Background(), "SELECT dex_number, pokemon_name FROM pokemon")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
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
	rows, err := dbpool.Query(context.Background(), "SELECT * FROM pokemon_type")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
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
