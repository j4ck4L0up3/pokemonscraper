package main

import (
	"fmt"
	"os"

	"github.com/jack-gaskins/pokemonscraper/pokescraper"
)

func main() {
	htmlString, fetchErr := pokescraper.FetchHTML("https://serebii.net/pokedex-sv/")
	if fetchErr != nil {
		fmt.Printf("error fetching from https://serebii.net/pokedex-sv/: %v", fetchErr)
	}

	// TODO: parse html and create parentNode pointer
	parentNode, parseErr := pokescraper.ParseHTML(htmlString)
	if parseErr != nil {
		fmt.Printf("error parsing string fetched from https://serebii.net/pokedex-sv/: %v", parseErr)
	}

	// TODO: retieve all the page urls from {elem: "option", attrKey: "value"}
	// do this in batches of 151,100, 135, 107, 156, 72, 88, 96, 120
	// for Kanto, Johto, Hoenn, Sinnoh, Unova, Kalos, Alola, Galar/Hisui, Paldea, respectively
	// store in map
	var pageUrls []string
	optionElem := "option"
	attrKey := "value"
	pokescraper.TraverseDOMAttr(parentNode, optionElem, attrKey, &pageUrls)

	var pokeIdNames []string
	for _, attrVal := range pageUrls {
		pokescraper.TraverseDOMText(parentNode, optionElem, attrKey, attrVal, &pokeIdNames)
	}

	filename := "playgroud.txt"
	file, osErr := os.Create(filename)
	if osErr != nil {
		fmt.Printf("error creating or opening file {%v}: %v", filename, osErr)
		return
	}

	defer file.Close()

	for _, attrVal := range pageUrls {
		_, osErr := file.WriteString(attrVal + "\n")
		if osErr != nil {
			fmt.Printf("error writing to file {%v}: %v", filename, osErr)
		}
	}

	// pokeByRegion := make(map[string][]string)
	// pokeGroups := make([][]string, 1025)

	// batches := []int{151,100, 135, 107, 156, 72, 88, 96, 120}
	// for i := 0; i <= len(batches); i++ {
	// 	for j := 0; j <= batches[i]; j++ {
	// 		pokeGroups[i] = append(pokeGroups[i], )
	// 	}
	// }

	// TODO: to get id, name separated in a tuple like item

	// TODO: to get types, search {elem: <img>, attrKey: alt}, then filter pattern \b-type (or something)
	// use TraverseDOMAttr

}
