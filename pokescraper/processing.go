package pokescraper

import (
	"fmt"
	"regexp"
)

// returns a string containing only the desired elements
func reduceString(expr string, str string) [][]string {
	htmlRe := regexp.MustCompile(expr)
	match := htmlRe.FindAllStringSubmatch(str, -1)
	return match
}

// yield html string containing one region's pokemon html
func batchHTMLString(htmlRawStr string, numRegions int) <-chan string {
	batch := make(chan string)

	// reduce the raw string with regexp
	exprMain := `(?s)<main>(.*?)<\/main>`
	reducedStrMain := reduceString(exprMain, htmlRawStr)

	// reduce <main> string into matrix of <table> elements
	exprTable := `(?s)<table (.*?)>(.*?)</table>`
	reducedStrTable := reduceString(exprTable, reducedStrMain[0][0])

	// reduce to matrix of <select> elements
	exprSelect := `(?s)<SELECT (.*?)>(.*?)</SELECT>`
	reducedStrSelect := reduceString(exprSelect, reducedStrTable[1][0])

	go func() {
		defer close(batch)
		for i := 0; i < numRegions; i++ {
			batch <- reducedStrSelect[i][0]
		}
	}()

	// return processed string batch
	return batch
}

// return string of text from html nodes
func parseHTMLText(htmlStr string) []string {
	// parse strings into html node pointers
	node, parseErr := ParseHTML(htmlStr)
	if parseErr != nil {
		fmt.Printf("error parsing string: %v", parseErr)
	}

	// get attribute values for parsing text nodes
	var pageUrls []string
	optionElem := "option"
	attrKey := "value"
	GetDOMAttrVals(node, optionElem, attrKey, &pageUrls)

	// parse text nodes in list
	var pokeIdNames []string
	for _, attrVal := range pageUrls {
		GetDOMText(node, optionElem, attrKey, attrVal, &pokeIdNames)
	}

	return pokeIdNames
}

// return a matrix of scraper pokemon IDs and names per region
func processPokemonMatrix(url string, numRegions int) [][]string {
	// convert page into string
	htmlRawStr, fetchErr := FetchHTML(url)
	if fetchErr != nil {
		fmt.Printf("FetchHTML Error:\n%v", fetchErr)
	}

	// get html region batches
	htmlStrBatches := []string{}
	for htmlStr := range batchHTMLString(htmlRawStr, numRegions) {
		htmlStrBatches = append(htmlStrBatches, htmlStr)
	}

	// create matrix of pokemon per region
	regionPokeIdNames := [][]string{}
	for _, htmlStr := range htmlStrBatches {
		pokeIdNames := parseHTMLText(htmlStr)
		regionPokeIdNames = append(regionPokeIdNames, pokeIdNames)
	}

	// get region name only as first element
	for i := 0; i < len(regionPokeIdNames); i++ {
		regionExpr := `^([^:]+)`
		match := reduceString(regionExpr, regionPokeIdNames[i][0])
		regionPokeIdNames[i][0] = match[0][0]
	}

	return regionPokeIdNames
}

// remove and return the first element in a list of strings
func popFirst(strList []string) (string, []string) {
	if len(strList) > 0 {
		firstElem := strList[0]
		strList := strList[1:]
		return firstElem, strList
	}
	return "", strList
}

// return map of regions and their pokemon IDs & names
func ProcessPokemonMap(url string, numRegions int) map[string][]string {
	pokemonMap := make(map[string][]string)
	pokemonMatrix := processPokemonMatrix(url, numRegions)
	for _, strList := range pokemonMatrix {
		region, pokeList := popFirst(strList)
		pokemonMap[region] = pokeList
	}

	return pokemonMap
}
