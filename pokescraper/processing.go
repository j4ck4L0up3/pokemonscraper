package pokescraper

import (
	"fmt"
	"regexp"
)

// returns a string containing only the desired elements
func reduceHTMLString(expr string, htmlStr string) [][]string {
	htmlRe := regexp.MustCompile(expr)
	match := htmlRe.FindAllStringSubmatch(htmlStr, -1)
	return match
}

func ProcessHTML(url string) [][]string {
	// get raw html string
	htmlRawString, fetchErr := FetchHTML(url)
	if fetchErr != nil {
		fmt.Printf("error fetching from {%v}: %v", url, fetchErr)
	}
	// reduce the raw string with regexp
	expr := `(?s)<main>(.*?)<\/main>`
	reducedStr := reduceHTMLString(expr, htmlRawString)
	// parse the html from the new string

	// get the attribute values from the string

	// get the text nodes from the elements that have text

	// return processed text
	return reducedStr
}
