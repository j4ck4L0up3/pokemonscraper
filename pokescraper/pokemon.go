package pokescraper

import (
	"errors"
)

type Pokemon struct {
	ID   string   `json:"id"`
	Name string   `json:"name"`
	Type []string `json:"type"`
}

// returns false for invalid Pokemon structs, true for valid ones
func ValidatePokemon(poke Pokemon) error {
	emptyPoke := Pokemon{}
	if IsEmptyPokemon(emptyPoke, poke) {
		err := errors.New("Empty Pokemon struct given")
		return err
	}

	if IsPartialPokemon(poke) {
		err := errors.New("Partial Pokemon struct given")
		return err
	}

	return nil
}

// returns true if Pokemon struct is empty, false if not
func IsEmptyPokemon(empty, poke Pokemon) bool {
	return empty.ID == poke.ID &&
		empty.Name == poke.Name &&
		len(empty.Type) == len(poke.Type)
}

// returns true if Pokemon struct is partial, false if not
func IsPartialPokemon(poke Pokemon) bool {
	return poke.ID == "" || poke.Name == "" || len(poke.Type) == 0
}
