package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/janek64/pmd-dx-api/api/models"
	"golang.org/x/sync/errgroup"
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
	defer rows.Close()
	if err != nil {
		return move, moveType, nil, err
	}
	var p models.MovePokemonID
	// Read the first row outside of the loop to extract move and type information and check for null pokemon
	rows.Next()
	err = rows.Scan(&move.MoveID, &move.MoveName, &move.Category, &move.Range, &move.Target, &move.InitialPP, &move.InitialPower, &move.Accuracy, &move.Description, &moveType.ID, &moveType.Name, &p.Method, &p.Cost, &p.Level, &p.Pokemon.ID, &p.Pokemon.Name)
	// Add the first pokemon to the slice
	// Check if the pokemon is not null to find move without pokemon
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

// GetMove fetches a move entry, its type and all pokemon learning it from the database by its ID or name.
func GetPokemon(input SearchInput) (pokemon models.Pokemon, camp models.NamedResourceID, abilities []models.NamedResourceID, dungeons []models.PokemonDungeonID, moves []models.PokemonMoveID, types []models.NamedResourceID, err error) {
	// Create a pgx.Rows variable for each query to be executed
	var rows [4]pgx.Rows
	// Create an errgroup.Group to wait until the goroutines have finished
	// Channels are not necessary since we work with closures
	errs, _ := errgroup.WithContext(context.Background())
	// Query 1 - pokemon, camp, dungeon
	errs.Go(func() error {
		// Use different query depending on search type
		if input.SearchType == ID {
			queryString := `SELECT P.*, C.camp_name, D.dungeon_ID, D.dungeon_name, PD.super_enemy
			FROM pokemon P INNER JOIN camp C ON P.dex_number = $1 AND P.camp_ID = C.camp_ID
			LEFT JOIN encountered_in PD ON P.dex_number = PD.dex_number
			LEFT JOIN dungeon D ON PD.dungeon_ID = D.dungeon_ID;`
			rows[0], err = dbpool.Query(context.Background(), queryString, input.ID)
			return err
		} else if input.SearchType == NAME {
			queryString := `SELECT P.*, C.camp_name, D.dungeon_ID, D.dungeon_name, PD.super_enemy
			FROM pokemon P INNER JOIN camp C ON P.pokemon_name = $1 AND P.camp_ID = C.camp_ID
			LEFT JOIN encountered_in PD ON P.dex_number = PD.dex_number
			LEFT JOIN dungeon D ON PD.dungeon_ID = D.dungeon_ID;`
			rows[0], err = dbpool.Query(context.Background(), queryString, input.Name)
			return err
		} else {
			return fmt.Errorf("illegal search type %v", input.SearchType)
		}
	})
	// Query 2 - pokemonTypes
	errs.Go(func() error {
		// Use different query depending on search type
		if input.SearchType == ID {
			queryString := `SELECT T.* FROM pokemon_type T
			INNER JOIN pokemon_has_type PT ON PT.dex_number = $1 AND PT.type_ID = T.type_ID;`
			rows[1], err = dbpool.Query(context.Background(), queryString, input.ID)
			return err
		} else if input.SearchType == NAME {
			queryString := `SELECT T.* FROM pokemon P
			INNER JOIN pokemon_has_type PT ON P.pokemon_name = $1 AND P.dex_number = PT.dex_number
			INNER JOIN pokemon_type T ON PT.type_ID = T.type_ID;`
			rows[1], err = dbpool.Query(context.Background(), queryString, input.Name)
			return err
		} else {
			return fmt.Errorf("illegal search type %v", input.SearchType)
		}
	})
	// Query 3 - abilities
	errs.Go(func() error {
		// Use different query depending on search type
		if input.SearchType == ID {
			queryString := `SELECT A.ability_ID, A.ability_name FROM ability A
			INNER JOIN pokemon_has_ability PA ON PA.dex_number = $1 AND PA.ability_ID = A.ability_ID;`
			rows[2], err = dbpool.Query(context.Background(), queryString, input.ID)
			return err
		} else if input.SearchType == NAME {
			queryString := `SELECT A.ability_ID, A.ability_name FROM pokemon P
			INNER JOIN pokemon_has_ability PA ON P.pokemon_name = $1 AND P.dex_number = PA.dex_number
			INNER JOIN ability A ON PA.ability_ID = A.ability_ID;`
			rows[2], err = dbpool.Query(context.Background(), queryString, input.Name)
			return err
		} else {
			return fmt.Errorf("illegal search type %v", input.SearchType)
		}
	})
	// Query 4 - moves
	errs.Go(func() error {
		// Use different query depending on search type
		if input.SearchType == ID {
			queryString := `SELECT M.move_ID, M.move_name, PM.learn_type, PM.cost, PM.level FROM attack_move M
			INNER JOIN learns PM ON PM.dex_number = $1 AND PM.move_ID = M.move_ID;`
			rows[3], err = dbpool.Query(context.Background(), queryString, input.ID)
			return err
		} else if input.SearchType == NAME {
			queryString := `SELECT M.move_ID, M.move_name, PM.learn_type, PM.cost, PM.level
			FROM pokemon P INNER JOIN learns PM ON P.pokemon_name = $1 AND P.dex_number = PM.dex_number
			INNER JOIN attack_move M ON PM.move_ID = M.move_ID;`
			rows[3], err = dbpool.Query(context.Background(), queryString, input.Name)
			return err
		} else {
			return fmt.Errorf("illegal search type %v", input.SearchType)
		}
	})
	// Wait for all Goroutines and check for any errors
	if err := errs.Wait(); err != nil {
		return pokemon, camp, nil, nil, nil, nil, err
	}
	// Close all rows after the function finished
	defer func() {
		for i := range rows {
			rows[i].Close()
		}
	}()
	// Read query 1
	var d models.PokemonDungeonID
	// Read the first row of query 1 outside of the loop to extract pokemon and camp information and check for null dungeon
	rows[0].Next()
	err = rows[0].Scan(&pokemon.DexNumber, &pokemon.PokemonName, &pokemon.EvolutionStage, &pokemon.EvolveCondition, &pokemon.EvolveLevel, &pokemon.EvolveCrystals, &pokemon.Classification, &camp.ID, &camp.Name, &d.Dungeon.ID, &d.Dungeon.Name, &d.IsSuper)
	// Add the first dungeon to the slice
	// Check if the dungeon is not null to find pokemon without dungeon
	if d.Dungeon.ID != 0 {
		dungeons = append(dungeons, d)
	}
	// Add all other pokemon to the slice
	for rows[0].Next() {
		// Use a throwaway models.Pokemon and models.NamedResourceID to ignore pokemon and camp data for all other rows
		var emptyPokemon models.Pokemon
		var emptyCamp models.NamedResourceID
		err = rows[0].Scan(&emptyPokemon.DexNumber, &emptyPokemon.PokemonName, &emptyPokemon.EvolutionStage, &emptyPokemon.EvolveCondition, &emptyPokemon.EvolveLevel, &emptyPokemon.EvolveCrystals, &emptyPokemon.Classification, &emptyCamp.ID, &emptyCamp.Name, &d.Dungeon.ID, &d.Dungeon.Name, &d.IsSuper)
		if err != nil {
			return pokemon, camp, nil, nil, nil, nil, err
		}
		// Checking for ID==0 is not necessary since all rows after the first will not have null values
		dungeons = append(dungeons, d)
	}
	// If the DexNumber is zero, no entry was found
	if pokemon.DexNumber == 0 {
		if input.SearchType == ID {
			return pokemon, camp, nil, nil, nil, nil, &ResourceNotFoundError{ResourceType: "pokemon", SearchType: input.SearchType, ID: input.ID}
		} else if input.SearchType == NAME {
			return pokemon, camp, nil, nil, nil, nil, &ResourceNotFoundError{ResourceType: "pokemon", SearchType: input.SearchType, Name: input.Name}
		}
	}
	// Read query 2
	for rows[1].Next() {
		var t models.NamedResourceID
		err = rows[1].Scan(&t.ID, &t.Name)
		if err != nil {
			return pokemon, camp, nil, nil, nil, nil, err
		}
		types = append(types, t)
	}
	// Read query 3
	for rows[2].Next() {
		var a models.NamedResourceID
		err = rows[2].Scan(&a.ID, &a.Name)
		if err != nil {
			return pokemon, camp, nil, nil, nil, nil, err
		}
		abilities = append(abilities, a)
	}
	// Read query 4
	for rows[3].Next() {
		var m models.PokemonMoveID
		err = rows[3].Scan(&m.Move.ID, &m.Move.Name, &m.Method, &m.Cost, &m.Level)
		if err != nil {
			return pokemon, camp, nil, nil, nil, nil, err
		}
		moves = append(moves, m)
	}
	return pokemon, camp, abilities, dungeons, moves, types, nil
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

// GetPokemonType fetches a pokemonType entry and its type interactions from the database by its ID or name.
func GetPokemonType(input SearchInput) (pokemonType models.PokemonType, interactions []models.TypeInteractionID, err error) {
	var rows pgx.Rows
	// Use different query depending on search type
	if input.SearchType == ID {
		queryString := `SELECT AT.*, TT.interaction, DT.*
		FROM (SELECT * FROM pokemon_type WHERE type_ID = $1) AT
		LEFT JOIN effectiveness TT ON AT.type_ID = TT.attacker
		LEFT JOIN pokemon_type DT ON TT.defender = DT.type_ID;`
		rows, err = dbpool.Query(context.Background(), queryString, input.ID)
	} else if input.SearchType == NAME {
		queryString := `SELECT AT.*, TT.interaction, DT.*
		FROM (SELECT * FROM pokemon_type WHERE type_name = $1) AT
		LEFT JOIN effectiveness TT ON AT.type_ID = TT.attacker
		LEFT JOIN pokemon_type DT ON TT.defender = DT.type_ID;`
		rows, err = dbpool.Query(context.Background(), queryString, input.Name)
	} else {
		return pokemonType, nil, fmt.Errorf("illegal search type %v", input.SearchType)
	}
	if err != nil {
		return pokemonType, nil, err
	}
	defer rows.Close()
	var i models.TypeInteractionID
	// Read the first row outside of the loop to extract pokemonType information and check for null interaction
	rows.Next()
	err = rows.Scan(&pokemonType.TypeID, &pokemonType.TypeName, &i.Interaction, &i.Defender.ID, &i.Defender.Name)
	// Add the first interaction to the slice
	// Check if the interaction is not null to find pokemonType without interaction
	if i.Defender.ID != 0 {
		interactions = append(interactions, i)
	}
	// Add all other pokemon to the slice
	for rows.Next() {
		// Use a throwaway models.PokemonType to ignore pokemonType data for all other rows
		var emptyPokemonType models.PokemonType
		err = rows.Scan(&emptyPokemonType.TypeID, &emptyPokemonType.TypeName, &i.Interaction, &i.Defender.ID, &i.Defender.Name)
		if err != nil {
			return pokemonType, nil, err
		}
		// Checking for ID==0 is not necessary since all rows after the first will not have null values
		interactions = append(interactions, i)
	}
	// If the TypeID is zero, no entry was found
	if pokemonType.TypeID == 0 {
		if input.SearchType == ID {
			return pokemonType, nil, &ResourceNotFoundError{ResourceType: "type", SearchType: input.SearchType, ID: input.ID}
		} else if input.SearchType == NAME {
			return pokemonType, nil, &ResourceNotFoundError{ResourceType: "type", SearchType: input.SearchType, Name: input.Name}
		}
	}
	return pokemonType, interactions, nil
}
