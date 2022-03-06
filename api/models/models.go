// Package models contains all data types for database
// entries and necessary custom types.
package models

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

// NullInt64 - extended custom type of sql.NullInt64.
type NullInt64 sql.NullInt64

// https://medium.com/aubergine-solutions/how-i-handled-null-possible-values-from-database-rows-in-golang-521fb0ee267
// https://husobee.github.io/golang/database/2015/06/12/scanner-valuer.html

// Value - Implementation of Valuer from database/sql/driver.
func (n *NullInt64) Value() (driver.Value, error) {
	return int64(n.Int64), nil
}

// Scan - Implementation of Scanner from database/sql.
func (n *NullInt64) Scan(src interface{}) error {
	// Scan the Input with the database/sql Scan function
	var i sql.NullInt64
	if err := i.Scan(src); err != nil {
		return err
	}

	// Use a type switch to check for nil values
	switch src.(type) {
	case nil:
		*n = NullInt64{i.Int64, false}
	case int64:
		*n = NullInt64{i.Int64, true}
	default:
		return errors.New("failed to scan NullInt64")
	}
	return nil
}

// MarshalJSON - Implementation of Marshaler from encoding/json.
func (n NullInt64) MarshalJSON() ([]byte, error) {
	// If there is a null value, return "null" as output
	if !n.Valid {
		return []byte("null"), nil
	}
	// Else, encode the int64
	return json.Marshal(n.Int64)
}

// AttackMove represents an attack_move entry from the database.
type AttackMove struct {
	MoveID       int
	MoveName     string
	Category     string
	MoveRange    string
	Target       string
	InitialPP    int
	InitialPower int
	Accuracy     int
	Description  string
}

// Ability represents an ability entry from the database.
type Ability struct {
	AbilityID   int
	AbilityName string
	Description string
}

// Camp represents a camp entry from the database.
type Camp struct {
	CampID      int
	CampName    string
	UnlockType  string
	Cost        NullInt64
	Description string
}

// Dungeon represents a dungeon entry from the database.
type Dungeon struct {
	DungeonID      int
	DungeonName    string
	Levels         int
	StartLevel     NullInt64
	TeamSize       int
	ItemsAllowed   bool
	PokemonJoining bool
	MapVisible     bool
}

// Pokemon represents a pokemon entry from the database.
type Pokemon struct {
	DexNumber       int
	PokemonName     string
	EvolutionStage  int
	EvolveCondition string
	EvolveLevel     NullInt64
	EvolveCrystals  NullInt64
	Classification  string
	CampID          int
}

// PokemonType represents a pokemon_type entry from the database.
type PokemonType struct {
	TypeID   int
	TypeName string
}

// NamedResourceID is a short representation of an API resource with its name and ID (for URL construction).
type NamedResourceID struct {
	Name string
	ID   int
}

// ToNamedResourceURL returns the named resource with its URL instead of the ID.
func (n *NamedResourceID) ToNamedResourceURL(instanceURL string, resourceTypeName string) NamedResourceURL {
	url := fmt.Sprintf("%v/v1/%v/%v", instanceURL, resourceTypeName, n.ID)
	return NamedResourceURL{Name: n.Name, URL: url}
}

// NamedResourceURL is a short representation of an API resource with its name and URL.
type NamedResourceURL struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// DungeonPokemonID is a short representation of a pokemon appearing in a dungeon with its ID.
type DungeonPokemonID struct {
	Pokemon NamedResourceID
	IsSuper bool
}

// ToDungeonPokemonURL returns the DungeonPokemon with its URL instead of the ID.
func (d *DungeonPokemonID) ToDungeonPokemonURL(instanceURL string) DungeonPokemonURL {
	return DungeonPokemonURL{Pokemon: d.Pokemon.ToNamedResourceURL(instanceURL, "pokemon"), IsSuper: d.IsSuper}
}

// DungeonPokemonURL is a short representation of a pokemon appearing in a dungeon with its URL.
type DungeonPokemonURL struct {
	Pokemon NamedResourceURL
	IsSuper bool
}

// MovePokemonID is a short representation of a pokemon learning a move with its ID.
type MovePokemonID struct {
	Pokemon NamedResourceID
	Method  string
	Level   NullInt64
	Cost    NullInt64
}

// ToMovePokemonURL returns the MovePokemon with its URL instead of the ID.
func (m *MovePokemonID) ToMovePokemonURL(instanceURL string) MovePokemonURL {
	return MovePokemonURL{Pokemon: m.Pokemon.ToNamedResourceURL(instanceURL, "pokemon"), Method: m.Method, Level: m.Level, Cost: m.Cost}
}

// MovePokemonURL is a short representation of a pokemon learning a move with its URL.
type MovePokemonURL struct {
	Pokemon NamedResourceURL
	Method  string
	Level   NullInt64
	Cost    NullInt64
}

// PokemonDungeonID is a short representation of a dungeon a pokemon appears in with its ID.
type PokemonDungeonID struct {
	Dungeon NamedResourceID
	IsSuper bool
}

// ToPokemonDungeonURL returns the PokemonDungeon with its URL instead of the ID.
func (p *PokemonDungeonID) ToPokemonDungeonURL(instanceURL string) PokemonDungeonURL {
	return PokemonDungeonURL{Dungeon: p.Dungeon.ToNamedResourceURL(instanceURL, "dungeons"), IsSuper: p.IsSuper}
}

// PokemonDungeonURL is a short representation of a dungeon a pokemon appears in with its URL.
type PokemonDungeonURL struct {
	Dungeon NamedResourceURL
	IsSuper bool
}

// PokemonMoveID is a short representation of a move learned by a pokemon with its ID.
type PokemonMoveID struct {
	Move   NamedResourceID
	Method string
	Level  NullInt64
	Cost   NullInt64
}

// ToPokemonMoveURL returns the PokemonMove with its URL instead of the ID.
func (p *PokemonMoveID) ToPokemonMoveURL(instanceURL string) PokemonMoveURL {
	return PokemonMoveURL{Move: p.Move.ToNamedResourceURL(instanceURL, "moves"), Method: p.Method, Level: p.Level, Cost: p.Cost}
}

// PokemonMoveURL is a short representation of a move learned by a pokemon with its URL.
type PokemonMoveURL struct {
	Move   NamedResourceURL
	Method string
	Level  NullInt64
	Cost   NullInt64
}

// TypeInteractionID represents an interaction of a type attacking another type with its ID.
type TypeInteractionID struct {
	Defender    NamedResourceID
	Interaction string
}

// ToTypeInteractionURL returns the TypeInteraction with its URL instead of the ID.
func (t *TypeInteractionID) ToTypeInteractionURL(instanceURL string) TypeInteractionURL {
	return TypeInteractionURL{Defender: t.Defender.ToNamedResourceURL(instanceURL, "types"), Interaction: t.Interaction}
}

// TypeInteractionID represents an interaction of a type attacking another type with its URL.
type TypeInteractionURL struct {
	Defender    NamedResourceURL
	Interaction string
}
