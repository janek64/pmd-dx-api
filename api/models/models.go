// Package main contains all data types for database
// entries and necessary custom types.
package models

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// NullInt64 - extended custom type of sql.NullInt64
type NullInt64 sql.NullInt64

// https://medium.com/aubergine-solutions/how-i-handled-null-possible-values-from-database-rows-in-golang-521fb0ee267
// https://husobee.github.io/golang/database/2015/06/12/scanner-valuer.html

// Value - Implementation of Valuer from database/sql/driver
func (n *NullInt64) Value() (driver.Value, error) {
	return int64(n.Int64), nil
}

// Scan - Implementation of Scanner from database/sql
func (n *NullInt64) Scan(src interface{}) error {
	// Scan the Input with the database/sql Scan function
	var i sql.NullInt64
	if err := i.Scan(src); err != nil {
		return err
	}

	// Use a type switch to check for nil values
	switch src.(type) {
	case nil:
		n = &NullInt64{i.Int64, false}
	case int64:
		n = &NullInt64{i.Int64, true}
	default:
		return errors.New("failed to scan NullInt64")
	}
	return nil
}

// MarshalJSON - Implementation of Marshaler from encoding/json
func (n *NullInt64) MarshalJSON() ([]byte, error) {
	// If there is a null value, return "null" as output
	if !n.Valid {
		return []byte("null"), nil
	}
	// Else, encode the int64
	return json.Marshal(n.Int64)
}

// Pokemon represents a pokemon entry from the database
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
