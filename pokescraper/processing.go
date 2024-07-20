package pokescraper

import (
	// "fmt"
	"regexp"
)

// returns a string containing only the desired elements
func reduceHTMLString(expr string, htmlStr string) [][]string {
	htmlRe := regexp.MustCompile(expr)
	match := htmlRe.FindAllStringSubmatch(htmlStr, -1)
	return match
}

// yield html string containing one region's pokemon html
func BatchHTMLString(htmlRawStr string, numRegions int) <-chan string {
	batch := make(chan string)

	// reduce the raw string with regexp
	exprMain := `(?s)<main>(.*?)<\/main>`
	reducedStrMain := reduceHTMLString(exprMain, htmlRawStr)

	// reduce <main> string into matrix of <table> elements
	exprTable := `(?s)<table (.*?)>(.*?)</table>`
	reducedStrTable := reduceHTMLString(exprTable, reducedStrMain[0][0])

	// reduce to matrix of <select> elements
	exprSelect := `(?s)<SELECT (.*?)>(.*?)</SELECT>`
	reducedStrSelect := reduceHTMLString(exprSelect, reducedStrTable[1][0])

	go func() {
		defer close(batch)
		for i := 0; i < numRegions; i++ {
			batch <- reducedStrSelect[i][0]
		}
	}()

	// return processed string batch
	return batch
}
