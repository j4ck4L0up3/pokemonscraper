package main

import (
	"fmt"

	"github.com/jack-gaskins/pokemonscraper/internal/pokescraper"
)

func main() {
	pikachu := pokescraper.Pokemon{
		ID:   025,
		Name: "pikachu",
		Type: "electric",
	}

	pikachuJSON, err := pokescraper.SerializePokemon(pikachu)
	if err != nil {
		fmt.Printf("Could not serialize Pokemon: %v\n", err)
	}

	fmt.Println(pikachuJSON)
}
