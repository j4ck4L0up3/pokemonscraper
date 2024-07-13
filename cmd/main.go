package main

import (
	"fmt"
	"reflect"

	"github.com/jack-gaskins/pokemonscraper/pokescraper"
)

func main() {
	pikachu := pokescraper.Pokemon{
		ID:   "025",
		Name: "pikachu",
		Type: []string{"electric"},
	}

	pikachuJSON, serialErr := pokescraper.SerializePokemon(pikachu)
	if serialErr != nil {
		fmt.Printf("Could not serialize Pokemon: %v\n", serialErr)
	}

	fmt.Println(pikachuJSON)
	fmt.Println(reflect.TypeOf(pikachuJSON))

	charizardJSON := `{"id": "006", "name": "charizard", "type": ["fire","flying"]}`
	charizard := pokescraper.Pokemon{}

	deserialErr := pokescraper.DeserializePokemon(charizardJSON, &charizard)
	if deserialErr != nil {
		fmt.Printf("Could not deserialize Pokemon: %v\n", deserialErr)
	}

	fmt.Println(charizard)
	fmt.Println(reflect.TypeOf(charizard))
}
