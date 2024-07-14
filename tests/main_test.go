package tests

import (
	"regexp"
	"testing"
)

//lint:ignore U1000 main function in test file
func main() {
	testMap := map[string]func(t *testing.T){
		"TestFetchHTML":          TestFetchHTML,
		"TestSerializePokemon":   TestSerializePokemon,
		"TestDeserializePokemon": TestDeserializePokemon,
	}

	var tests []testing.InternalTest
	for testName, testFunc := range testMap {
		tests = append(tests, testing.InternalTest{
			Name: testName,
			F:    testFunc,
		})
	}

	testing.Main(matchPattern, tests, nil, nil)

}

// helper function that returns true for patterns matched in flags
// to allow user to choose tests to run
// example input: go test -run=Testxxx.*
func matchPattern(pattern, str string) (bool, error) {
	return regexp.MatchString(pattern, str)
}
