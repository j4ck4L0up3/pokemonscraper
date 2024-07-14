package main

func main() {
	// TODO: fetch html of https://serebii.net/pokedex-sv/

	// TODO: parse html and create parentNode pointer

	// TODO: retieve all the page urls from {elem: "option", attrKey: "value"}
	// do this in batches of 151,100, 135, 107, 156, 72, 88, 96, 120
	// for Kanto, Johto, Hoenn, Sinnoh, Unova, Kalos, Alola, Galar/Hisui, Paldea, respectively
	// store in map

	// TODO: to get id, name

	// TODO: to get types, search {elem: <img>, attrKey: alt}, then filter pattern \b-type (or something)
	// use TraverseDOMAttr

}
