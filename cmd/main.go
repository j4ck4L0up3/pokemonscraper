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

	// TODO: first retieve all the page urls from {elem: "option", attrKey: "value"}
	// do this in batches of 151,100, 135, 107, 156, 72, 88, 96, 120
	// for Kanto, Johto, Hoenn, Sinnoh, Unova, Kalos, Alola, Galar/Hisui, Paldea, respectively

	// TODO: to get id, name

	// TODO: to get types, search {elem: <img>, attrKey: alt}, then filter pattern \b-type (or something)
	// use TraverseDOMAttr

}
