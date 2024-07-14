package tests

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/jack-gaskins/pokemonscraper/pokescraper"
)

func TestTraverseDOMAttr(t *testing.T) {
	// Test case: mock.html with similar structure to actual site
	t.Run("mock.html with similar structure to actual site", func(t *testing.T) {
		builder := new(strings.Builder)
		mockData, readErr := os.ReadFile("mock.html")
		if readErr != nil {
			t.Fatalf("error reading the mock html file for expected data:\n%v", readErr)
		}

		mockHtml := bytes.NewReader(mockData)
		_, ioErr := io.Copy(builder, mockHtml)
		if ioErr != nil {
			t.Fatalf("error copying byte array to strings.Builder object:\n%v", ioErr)
		}

		mockStr := builder.String()
		mockNode, parseErr := pokescraper.ParseHTML(mockStr)
		if parseErr != nil {
			t.Fatalf("expected no error, got error: %v", parseErr)
		}

		var values []string
		elem := "option"
		attrKey := "value"
		pokescraper.TraverseDOMAttr(mockNode, elem, attrKey, &values)

		if len(values) == 0 {
			t.Fatal("error retrieving values from html.Node, expected filled array, got empty array")
		}

		t.Logf("Values retrieved: %v", values)
	})
}

func TestTraverseDOMAttrBatch(t *testing.T) {
	// Test case: mock.html with batch
	t.Run("mock.html with batch", func(t *testing.T) {
		builder := new(strings.Builder)
		mockData, readErr := os.ReadFile("mock.html")
		if readErr != nil {
			t.Fatalf("error reading the mock html file for expected data:\n%v", readErr)
		}

		mockHtml := bytes.NewReader(mockData)
		_, ioErr := io.Copy(builder, mockHtml)
		if ioErr != nil {
			t.Fatalf("error copying byte array to strings.Builder object:\n%v", ioErr)
		}

		mockStr := builder.String()
		mockNode, parseErr := pokescraper.ParseHTML(mockStr)
		if parseErr != nil {
			t.Fatalf("expected no error, got error: %v", parseErr)
		}

		var values []string
		elem := "option"
		attrKey := "value"
		testBatch := 3
		pokescraper.TraverseDOMAttrBatch(mockNode, elem, attrKey, &values, testBatch)

		if len(values) == 0 {
			t.Fatal("error retrieving values from html.Node, expected filled array, got empty array")
		}

		if len(values) != testBatch {
			t.Fatalf("error filling array, expected {%v} elements, got {%v}", testBatch, len(values))
		}

		t.Logf("Values retrieved: %v", values)
	})

	// Test case: mock.html with no batch
	t.Run("mock.html with no batch", func(t *testing.T) {
		builder := new(strings.Builder)
		mockData, readErr := os.ReadFile("mock.html")
		if readErr != nil {
			t.Fatalf("error reading the mock html file for expected data:\n%v", readErr)
		}

		mockHtml := bytes.NewReader(mockData)
		_, ioErr := io.Copy(builder, mockHtml)
		if ioErr != nil {
			t.Fatalf("error copying byte array to strings.Builder object:\n%v", ioErr)
		}

		mockStr := builder.String()
		mockNode, parseErr := pokescraper.ParseHTML(mockStr)
		if parseErr != nil {
			t.Fatalf("expected no error, got error: %v", parseErr)
		}

		var values []string
		elem := "option"
		attrKey := "value"
		pokescraper.TraverseDOMAttrBatch(mockNode, elem, attrKey, &values)

		if len(values) == 0 {
			t.Fatal("error retrieving values from html.Node, expected filled array, got empty array")
		}

		if len(values) != 1 {
			t.Fatalf("error filling array, expected {%v} elements, got {%v}", 1, len(values))
		}

		t.Logf("Values retrieved: %v", values)
	})

	// Test case: mock.html with zero batch
	t.Run("mock.html with zero batch", func(t *testing.T) {
		builder := new(strings.Builder)
		mockData, readErr := os.ReadFile("mock.html")
		if readErr != nil {
			t.Fatalf("error reading the mock html file for expected data:\n%v", readErr)
		}

		mockHtml := bytes.NewReader(mockData)
		_, ioErr := io.Copy(builder, mockHtml)
		if ioErr != nil {
			t.Fatalf("error copying byte array to strings.Builder object:\n%v", ioErr)
		}

		mockStr := builder.String()
		mockNode, parseErr := pokescraper.ParseHTML(mockStr)
		if parseErr != nil {
			t.Fatalf("expected no error, got error: %v", parseErr)
		}

		var values []string
		elem := "option"
		attrKey := "value"
		testBatch := 0
		pokescraper.TraverseDOMAttrBatch(mockNode, elem, attrKey, &values, testBatch)

		if len(values) == 0 {
			t.Fatal("error retrieving values from html.Node, expected filled array, got empty array")
		}

		if len(values) != 1 {
			t.Fatalf("error filling array, expected {%v} elements, got {%v}", 1, len(values))
		}

		t.Logf("Values retrieved: %v", values)
	})

	// Test case: mock.html with negative batch
	t.Run("mock.html with negative batch", func(t *testing.T) {
		builder := new(strings.Builder)
		mockData, readErr := os.ReadFile("mock.html")
		if readErr != nil {
			t.Fatalf("error reading the mock html file for expected data:\n%v", readErr)
		}

		mockHtml := bytes.NewReader(mockData)
		_, ioErr := io.Copy(builder, mockHtml)
		if ioErr != nil {
			t.Fatalf("error copying byte array to strings.Builder object:\n%v", ioErr)
		}

		mockStr := builder.String()
		mockNode, parseErr := pokescraper.ParseHTML(mockStr)
		if parseErr != nil {
			t.Fatalf("expected no error, got error: %v", parseErr)
		}

		var values []string
		elem := "option"
		attrKey := "value"
		testBatch := -1
		pokescraper.TraverseDOMAttrBatch(mockNode, elem, attrKey, &values, testBatch)

		if len(values) == 0 {
			t.Fatal("error retrieving values from html.Node, expected filled array, got empty array")
		}

		if len(values) != 1 {
			t.Fatalf("error filling array, expected {%v} elements, got {%v}", 1, len(values))
		}

		t.Logf("Values retrieved: %v", values)
	})
}

func TestTraverseDOMText(t *testing.T) {
	// Test case: mock.html
	t.Run("mock.html", func(t *testing.T) {
		builder := new(strings.Builder)
		mockData, readErr := os.ReadFile("mock.html")
		if readErr != nil {
			t.Fatalf("error reading the mock html file for expected data:\n%v", readErr)
		}

		mockHtml := bytes.NewReader(mockData)
		_, ioErr := io.Copy(builder, mockHtml)
		if ioErr != nil {
			t.Fatalf("error copying byte array to strings.Builder object:\n%v", ioErr)
		}

		mockStr := builder.String()
		mockNode, parseErr := pokescraper.ParseHTML(mockStr)
		if parseErr != nil {
			t.Fatalf("expected no error, got error: %v", parseErr)
		}

		var values []string
		elem := "option"
		attrKey := "value"
		attrVals := []string{
			"/pokedex-sv/name1/", "/pokedex-sv/name2/", "/pokedex-sv/name3/",
			"/pokedex-sv/name4/", "/pokedex-sv/name5/", "/pokedex-sv/name6/",
			"/pokedex-sv/name7/", "/pokedex-sv/name8/", "/pokedex-sv/name9/",
		}
		for _, attrVal := range attrVals {
			pokescraper.TraverseDOMText(mockNode, elem, attrKey, attrVal, &values)
		}

		if len(values) == 0 {
			t.Fatal("error retrieving values from html.Node, expected filled array, got empty array")
		}

		if len(values) != len(attrVals) {
			t.Fatalf("error filling array, expected {%v} elements, got {%v}", len(attrVals), len(values))
		}

		t.Logf("Values retrieved: %v", values)
	})
}
