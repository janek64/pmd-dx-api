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

// GetAbilityList fetches a slice of all ability entries from the database.
func GetAbilityList() ([]models.NamedResourceID, error) {
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

// GetAbility fetches an ability entry and all pokemon that have it from the database by its ID or name.
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
	// Add the first pokemon to the slice
	// Check if the pokemon is not null to find ability without pokemon
	if p.ID != 0 {
		pokemon = append(pokemon, p)
	}
	// Add all other pokemon to the slice
	for rows.Next() {
		// Use a throwaway models.Ability to ignore ability data for all other rows
		var emptyAbility models.Ability
		err = rows.Scan(&emptyAbility.AbilityID, &emptyAbility.AbilityName, &emptyAbility.Description, &p.ID, &p.Name)
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

// GetCampList fetches a slice of all camp entries from the database.
func GetCampList() ([]models.NamedResourceID, error) {
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

// GetCamp fetches a camp entry and all pokemon living in it from the database by its ID or name.
func GetCamp(input SearchInput) (camp models.Camp, pokemon []models.NamedResourceID, err error) {
	var rows pgx.Rows
	// Use different query depending on search type
	if input.SearchType == ID {
		queryString := `SELECT C.*, P.dex_number, P.pokemon_name
		FROM (SELECT * FROM camp WHERE camp_ID = $1) C
		LEFT JOIN pokemon P ON C.camp_ID = P.camp_ID;`
		rows, err = dbpool.Query(context.Background(), queryString, input.ID)
	} else if input.SearchType == NAME {
		queryString := `SELECT C.*, P.dex_number, P.pokemon_name
		FROM (SELECT * FROM camp WHERE camp_name = $1) C
		LEFT JOIN pokemon P ON C.camp_ID = P.camp_ID;`
		rows, err = dbpool.Query(context.Background(), queryString, input.Name)
	} else {
		return camp, nil, fmt.Errorf("illegal search type %v", input.SearchType)
	}
	if err != nil {
		return camp, nil, err
	}
	defer rows.Close()
	var p models.NamedResourceID
	// Read the first row outside of the loop to extract camp information and check for null pokemon
	rows.Next()
	err = rows.Scan(&camp.CampID, &camp.CampName, &camp.UnlockType, &camp.Cost, &camp.Description, &p.ID, &p.Name)
	// Add the first pokemon to the slice
	// Check if the pokemon is not null to find camp without pokemon
	if p.ID != 0 {
		pokemon = append(pokemon, p)
	}
	// Add all other pokemon to the slice
	for rows.Next() {
		// Use a throwaway models.Camp to ignore camp data for all other rows
		var emptyCamp models.Camp
		err = rows.Scan(&emptyCamp.CampID, &emptyCamp.CampName, &emptyCamp.UnlockType, &emptyCamp.Cost, &emptyCamp.Description, &p.ID, &p.Name)
		if err != nil {
			return camp, nil, err
		}
		// Checking for ID==0 is not necessary since all rows after the first will not have null values
		pokemon = append(pokemon, p)
	}
	// If the CampID is zero, no entry was found
	if camp.CampID == 0 {
		if input.SearchType == ID {
			return camp, nil, &ResourceNotFoundError{ResourceType: "camp", SearchType: input.SearchType, ID: input.ID}
		} else if input.SearchType == NAME {
			return camp, nil, &ResourceNotFoundError{ResourceType: "camp", SearchType: input.SearchType, Name: input.Name}
		}
	}
	return camp, pokemon, nil
}

// GetDungeonList fetches a slice of all dungeon entries from the database.
func GetDungeonList() ([]models.NamedResourceID, error) {
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

// GetDungeon fetches a dungeon entry and all pokemon encountered in it from the database by its ID or name.
func GetDungeon(input SearchInput) (dungeon models.Dungeon, pokemon []models.DungeonPokemonID, err error) {
	var rows pgx.Rows
	// Use different query depending on search type
	if input.SearchType == ID {
		queryString := `SELECT D.*, DP.super_enemy, P.dex_number, P.pokemon_name
		FROM (SELECT * FROM dungeon WHERE dungeon_ID = $1) D
		LEFT JOIN encountered_in DP ON D.dungeon_ID = DP.dungeon_ID
		LEFT JOIN pokemon P ON DP.dex_number = P.dex_number;`
		rows, err = dbpool.Query(context.Background(), queryString, input.ID)
	} else if input.SearchType == NAME {
		queryString := `SELECT D.*, DP.super_enemy, P.dex_number, P.pokemon_name
		FROM (SELECT * FROM dungeon WHERE dungeon_name = $1) D
		LEFT JOIN encountered_in DP ON D.dungeon_ID = DP.dungeon_ID
		LEFT JOIN pokemon P ON DP.dex_number = P.dex_number;`
		rows, err = dbpool.Query(context.Background(), queryString, input.Name)
	} else {
		return dungeon, nil, fmt.Errorf("illegal search type %v", input.SearchType)
	}
	if err != nil {
		return dungeon, nil, err
	}
	defer rows.Close()
	var p models.DungeonPokemonID
	// Read the first row outside of the loop to extract dungeon information and check for null pokemon
	rows.Next()
	err = rows.Scan(&dungeon.DungeonID, &dungeon.DungeonName, &dungeon.Levels, &dungeon.StartLevel, &dungeon.TeamSize, &dungeon.ItemsAllowed, &dungeon.PokemonJoining, &dungeon.MapVisible, &p.IsSuper, &p.Pokemon.ID, &p.Pokemon.Name)
	// Add the first pokemon to the slice
	// Check if the pokemon is not null to find dungeon without pokemon
	if p.Pokemon.ID != 0 {
		pokemon = append(pokemon, p)
	}
	// Add all other pokemon to the slice
	for rows.Next() {
		// Use a throwaway models.Dungeon to ignore dungeon data for all other rows
		var emptyDungeon models.Dungeon
		err = rows.Scan(&emptyDungeon.DungeonID, &emptyDungeon.DungeonName, &emptyDungeon.Levels, &emptyDungeon.StartLevel, &emptyDungeon.TeamSize, &emptyDungeon.ItemsAllowed, &emptyDungeon.PokemonJoining, &emptyDungeon.MapVisible, &p.IsSuper, &p.Pokemon.ID, &p.Pokemon.Name)
		if err != nil {
			return dungeon, nil, err
		}
		// Checking for ID==0 is not necessary since all rows after the first will not have null values
		pokemon = append(pokemon, p)
	}
	// If the DungeonID is zero, no entry was found
	if dungeon.DungeonID == 0 {
		if input.SearchType == ID {
			return dungeon, nil, &ResourceNotFoundError{ResourceType: "dungeon", SearchType: input.SearchType, ID: input.ID}
		} else if input.SearchType == NAME {
			return dungeon, nil, &ResourceNotFoundError{ResourceType: "dungeon", SearchType: input.SearchType, Name: input.Name}
		}
	}
	return dungeon, pokemon, nil
}

// GetMoveList fetches a slice of all attack_move entries from the database.
func GetMoveList() ([]models.NamedResourceID, error) {
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

// GetMove fetches a move entry, its type and all pokemon learning it from the database by its ID or name.
func GetMove(input SearchInput) (move models.AttackMove, moveType models.NamedResourceID, pokemon []models.MovePokemonID, err error) {
	var rows pgx.Rows
	// Use different query depending on search type
	if input.SearchType == ID {
		queryString := `SELECT M.*, T.type_name, MP.learn_type, MP.cost, MP.level,
		P.dex_number, P.pokemon_name FROM attack_move M
		INNER JOIN pokemon_type T ON M.move_ID = $1 AND M.type_ID = T.type_ID
		LEFT JOIN learns MP ON MP.move_ID = M.move_ID
		LEFT JOIN pokemon P ON MP.dex_number = P.dex_number;`
		rows, err = dbpool.Query(context.Background(), queryString, input.ID)
	} else if input.SearchType == NAME {
		queryString := `SELECT M.*, T.type_name, MP.learn_type, MP.cost, MP.level,
		P.dex_number, P.pokemon_name FROM attack_move M
		INNER JOIN pokemon_type T ON M.move_name = $1 AND M.type_ID = T.type_ID
		LEFT JOIN learns MP ON MP.move_ID = M.move_ID
		LEFT JOIN pokemon P ON MP.dex_number = P.dex_number;`
		rows, err = dbpool.Query(context.Background(), queryString, input.Name)
	} else {
		return move, moveType, nil, fmt.Errorf("illegal search type %v", input.SearchType)
	}
	if err != nil {
		return move, moveType, nil, err
	}
	defer rows.Close()
	var p models.MovePokemonID
	// Read the first row outside of the loop to extract move and type information and check for null pokemon
	rows.Next()
	err = rows.Scan(&move.MoveID, &move.MoveName, &move.Category, &move.Range, &move.Target, &move.InitialPP, &move.InitialPower, &move.Accuracy, &move.Description, &moveType.ID, &moveType.Name, &p.Method, &p.Cost, &p.Level, &p.Pokemon.ID, &p.Pokemon.Name)
	// Add the first pokemon to the slice
	// Check if the pokemon is not null to find ability without pokemon
	if p.Pokemon.ID != 0 {
		pokemon = append(pokemon, p)
	}
	// Add all other pokemon to the slice
	for rows.Next() {
		// Use a throwaway models.Dungeon and models.NamedResourceID to ignore move and type data for all other rows
		var emptyMove models.AttackMove
		var emptyMoveType models.NamedResourceID
		err = rows.Scan(&emptyMove.MoveID, &emptyMove.MoveName, &emptyMove.Category, &emptyMove.Range, &emptyMove.Target, &emptyMove.InitialPP, &emptyMove.InitialPower, &emptyMove.Accuracy, &emptyMove.Description, &emptyMoveType.ID, &emptyMoveType.Name, &p.Method, &p.Cost, &p.Level, &p.Pokemon.ID, &p.Pokemon.Name)
		if err != nil {
			return move, moveType, nil, err
		}
		// Checking for ID==0 is not necessary since all rows after the first will not have null values
		pokemon = append(pokemon, p)
	}
	// If the MoveID is zero, no entry was found
	if move.MoveID == 0 {
		if input.SearchType == ID {
			return move, moveType, nil, &ResourceNotFoundError{ResourceType: "move", SearchType: input.SearchType, ID: input.ID}
		} else if input.SearchType == NAME {
			return move, moveType, nil, &ResourceNotFoundError{ResourceType: "move", SearchType: input.SearchType, Name: input.Name}
		}
	}
	return move, moveType, pokemon, nil
}

// GetPokemonList fetches a slice of all pokemon entries from the database.
func GetPokemonList() ([]models.NamedResourceID, error) {
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

// GetPokemonTypeList fetches a slice of all pokemon_type entries from the database.
func GetPokemonTypeList() ([]models.NamedResourceID, error) {
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
