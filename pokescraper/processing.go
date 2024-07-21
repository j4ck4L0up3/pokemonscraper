package pokescraper

import (
	"fmt"
	"golang.org/x/net/html"
	"regexp"
	"strings"
	"time"
)

// returns a string containing only the desired elements
func reduceString(expr string, str string) [][]string {
	re := regexp.MustCompile(expr)
	match := re.FindAllStringSubmatch(str, -1)
	return match
}

// remove a string that matches a specific pattern
func removeString(expr string, str string) string {
	re := regexp.MustCompile(expr)
	cleanedStr := re.ReplaceAllString(str, "")
	return cleanedStr
}

// yield html string containing one region's pokemon html
func batchPokeHtmlString(htmlRawStr string, numRegions int) <-chan string {
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

// return string of text from html nodes on pokedex page and page urls
func parsePokeHtmlText(htmlStr string) ([]string, []string) {
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

	return pokeIdNames, pageUrls
}

// return string of text from html nodes on type page
func parseTypeHtmlText(htmlRawStr string, pokeHrefs []string) []string {
	// reduce the raw string with regexp
	exprMain := `(?s)<main>(.*?)<\/main>`
	reducedStrMain := reduceString(exprMain, htmlRawStr)

	// remove all img tags b/c they ruin this
	exprImg := `(?s)<table\s+class="pkmn"[^>]*>.*?</table>`
	cleanedStr := removeString(exprImg, reducedStrMain[0][0])

	// parse strings into html node pointers
	node, parseErr := ParseHTML(cleanedStr)
	if parseErr != nil {
		fmt.Printf("error parsing string: %v", parseErr)
	}

	var pokeItems []string
	tableDataElem := "td"
	attrKey := "class"
	attrVal := "fooinfo"
	GetDOMText(node, tableDataElem, attrKey, attrVal, &pokeItems)

	pokeIds := []string{}
	for _, items := range pokeItems {
		exprRemove := `\n`
		items = removeString(exprRemove, items)
		exprID := `#\d{4}`
		match := reduceString(exprID, items)
		if len(match) > 0 {
			item := match[0][0]
			pokeIds = append(pokeIds, item)
		}
	}

	return pokeIds
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
	for htmlStr := range batchPokeHtmlString(htmlRawStr, numRegions) {
		htmlStrBatches = append(htmlStrBatches, htmlStr)
	}

	// create matrix of pokemon per region
	regionPokeIdNames := [][]string{}
	for _, htmlStr := range htmlStrBatches {
		pokeIdNames, _ := parsePokeHtmlText(htmlStr)
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

func ParseTypePageUrls(baseUrl string) []string {
	typeUris := []string{}
	typeUrls := []string{}

	// get the <a> href attributes
	htmlRawStr, fetchErr := FetchHTML(baseUrl)
	if fetchErr != nil {
		fmt.Printf("FetchHTML Error:\n%v\n", fetchErr)
	}

	node, parseErr := ParseHTML(htmlRawStr)
	if parseErr != nil {
		fmt.Printf("ParseHTML Error:\n%v\n", parseErr)
	}

	var typeNodes []*html.Node
	imgElem := "img"
	classKey := "class"
	classVal := "typeimg"
	GetDOMParentNode(node, imgElem, classKey, classVal, &typeNodes)

	aElem := "a"
	attrKey := "href"
	for _, parentNode := range typeNodes {
		GetDOMAttrVals(parentNode, aElem, attrKey, &typeUris)
	}

	// get type urls
	for _, uri := range typeUris {
		expr := `[^/]+$`
		match := reduceString(expr, uri)
		url := baseUrl + match[0][0]
		typeUrls = append(typeUrls, url)
	}

	return typeUrls
}

// return list of pokemon names that exist in a certain type
func ProcessPokemonTypeMap(baseUrl string) map[string][]string {
	typeMap := map[string][]string{}

	htmlRawStr, fetchErr := FetchHTML(baseUrl)
	if fetchErr != nil {
		fmt.Printf("FetchHTML Error:\n%v\n", fetchErr)
	}

	// get pokeHrefs for attrVals
	_, pokeHrefs := parsePokeHtmlText(htmlRawStr)

	// get type urls
	typeUrls := ParseTypePageUrls(baseUrl)

	for _, typeUrl := range typeUrls {
		// convert type page into string
		htmlRawStr, fetchErr := FetchHTML(typeUrl)
		if fetchErr != nil {
			fmt.Printf("FetchHTML Error:\n%v", fetchErr)
		}

		// grab type name from typeUrl
		expr := `(?:.*\/)?([^\/]+)\.shtml?$`
		match := reduceString(expr, typeUrl)
		typeName := match[0][1]

		// grab list of pokemon on page
		pokeList := parseTypeHtmlText(htmlRawStr, pokeHrefs)

		typeMap[typeName] = pokeList

		// sleep for server
		fmt.Print("Sleeping for server...\n")
		time.Sleep(5 * time.Second)
	}

	return typeMap
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

// split strings in a list of strings, returns list of lists of strings
func splitStrInList(strList []string) [][]string {
	stringMatrix := [][]string{}
	for _, str := range strList {
		splitList := strings.Split(str, " ")
		stringMatrix = append(stringMatrix, splitList)
	}
	return stringMatrix
}

// return map of regions and their pokemon IDs & names
func ProcessPokemonMap(url string, numRegions int) map[string][][]string {
	pokemonMap := make(map[string][][]string)
	pokemonMatrix := processPokemonMatrix(url, numRegions)
	for i, strList := range pokemonMatrix {
		region, pokeList := popFirst(strList)
		pokeMatrix := splitStrInList(pokeList)
		// TODO: figure out how to add a lending digit to IDs of length 3
		if len(pokeMatrix[i][0]) == 3 {
			pokeMatrix[i][0] = "0" + pokeMatrix[i][0]
		}
		pokemonMap[region] = pokeMatrix
	}

	return pokemonMap
}

// returns true if item is in slice, false if not
func sliceContains(slice []string, searchItem string) bool {
	for _, item := range slice {
		if item == searchItem {
			return true
		}
	}
	return false
}

func SetPokemon(url string, numRegions int) []Pokemon {
	typeMap := ProcessPokemonTypeMap(url)
	pokeMap := ProcessPokemonMap(url, numRegions)
	pokeList := []Pokemon{}

	for region, matrix := range pokeMap {
		for _, idName := range matrix {
			id := idName[0]
			name := idName[1]

			var types []string
			for typeName, idList := range typeMap {
				for _, typeId := range idList {
					if strings.Contains(typeId, id) && !sliceContains(types, typeName) {
						types = append(types, typeName)
					}
				}
			}

			pokemon := Pokemon{
				ID:     id,
				Name:   name,
				Type:   types,
				Region: region,
			}

			pokeList = append(pokeList, pokemon)
		}
	}

	return pokeList
}
