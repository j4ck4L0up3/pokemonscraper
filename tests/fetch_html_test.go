package tests

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/jack-gaskins/pokemonscraper/pokescraper"
)

func TestFetchHTML(t *testing.T) {
	// Test case: Successful fetch with html
	t.Run("successful fetch with mock HTML", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			htmlFile, readErr := os.ReadFile("mock.html")
			if readErr != nil {
				t.Fatalf("unable to read mock html file for mock server:\n%v", readErr)
			}

			w.Header().Set("Content-Type", "text/html")

			w.WriteHeader(http.StatusOK)
			_, writeErr := w.Write(htmlFile)
			if writeErr != nil {
				t.Fatalf("failed to write mock html file to mock server.\nwrite error: %v", writeErr)
			}
		}))

		defer mockServer.Close()

		builder := new(strings.Builder)
		expectedData, readErr := os.ReadFile("mock.html")
		if readErr != nil {
			t.Fatalf("error reading the mock html file for expected data:\n%v", readErr)
		}

		expectedHtml := bytes.NewReader(expectedData)
		_, ioErr := io.Copy(builder, expectedHtml)
		if ioErr != nil {
			t.Fatalf("error copying byte array to strings.Builder object:\n%v", ioErr)
		}

		expectedStr := builder.String()

		actualStr, fetchErr := pokescraper.FetchHTML(mockServer.URL)
		if fetchErr != nil {
			t.Fatalf("expected no error, got error: %v", fetchErr)
		}

		if !reflect.DeepEqual(expectedStr, actualStr) {
			t.Fatalf("actual string does not match expected string.\nactual string: %v\nexpected string: %v", actualStr, expectedStr)
		}

		t.Logf("HTML string retrieved:\n%v", actualStr)
	})

	// Test case: Successful fetch with bytes
	t.Run("succcesful fetch with bytes", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, writeErr := w.Write([]byte("Hello, World!"))
			if writeErr != nil {
				t.Fatalf("failed to write bytes to mock server:\n%v", writeErr)
			}
		}))

		defer mockServer.Close()

		expectedStr := "Hello, World!"
		actualStr, fetchErr := pokescraper.FetchHTML(mockServer.URL)
		if fetchErr != nil {
			t.Fatalf("expected no error, got error:\n%v", fetchErr)
		}

		if actualStr != expectedStr {
			t.Fatalf("actual string does not match expected string.\nactual string: %v\nexpected string: %v", actualStr, expectedStr)
		}

		t.Logf("Actual string retrieved:\n%v", actualStr)
	})

	// Test case: HTTP error
	t.Run("failed fetch with HTTP error", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))

		defer mockServer.Close()

		actualStr, fetchErr := pokescraper.FetchHTML(mockServer.URL)
		if fetchErr == nil {
			t.Fatal("expected error, received nil")
		}

		if actualStr != "" {
			t.Fatalf("expected empty string, received non-empty string:\n%v", actualStr)
		}

		t.Logf("FetchHTML error received:\n%v", fetchErr)
	})

	// Test case: Network error
	t.Run("failed fetch with network error", func(t *testing.T) {
		url := "http://invalid.url"
		actualStr, fetchErr := pokescraper.FetchHTML(url)
		if fetchErr == nil {
			t.Fatal("expected error, received nil")
		}

		if actualStr != "" {
			t.Fatalf("expected empty string, received non-empty string:\n%v", actualStr)
		}

		t.Logf("FetchHTML error received:\n%v", fetchErr)
	})

}

func TestParseHTML(t *testing.T) {
	// Test case: mock.html converted into string by FetchHTML
	t.Run("mock.html as string", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			htmlFile, readErr := os.ReadFile("mock.html")
			if readErr != nil {
				t.Fatalf("unable to read mock html file for mock server:\n%v", readErr)
			}

			w.Header().Set("Content-Type", "text/html")

			w.WriteHeader(http.StatusOK)
			_, writeErr := w.Write(htmlFile)
			if writeErr != nil {
				t.Fatalf("failed to write mock html file to mock server.\nwrite error: %v", writeErr)
			}
		}))

		defer mockServer.Close()

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

		if mockNode.Type == 0 {
			t.Fatal("no html.Node detected")
		}

		if mockNode == nil {
			t.Fatal("expected *html.Node, got nil")
		}

	})

	// Test case: Valid html content string
	t.Run("valid html content string", func(t *testing.T) {
		htmlContent := "<html><body><h1>Hello, World!</h1></body></html>"
		validNode, err := pokescraper.ParseHTML(htmlContent)
		if err != nil {
			t.Fatalf("expected no error, got error:\n%v", err)
		}

		if validNode == nil {
			t.Fatal("expected html.Node value, got nil")
		}

		t.Logf("HTML Node pointer NodeType: %v", validNode.Type)
	})

	// Test case: Empty string
	t.Run("empty string", func(t *testing.T) {
		emptyStr := ""
		emptyNode, err := pokescraper.ParseHTML(emptyStr)
		if err == nil {
			t.Fatalf("expected error, got nil\nReturned NodeType: %v", emptyNode.Type)
		}

		if emptyNode != nil {
			t.Fatalf("expected nil, got:\n%v", emptyNode)
		}
	})
}

func TestTraverseDOM(t *testing.T) {
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
		pokescraper.TraverseDOM(mockNode, elem, attrKey, &values)

		if len(values) == 0 {
			t.Fatal("error retrieving values from html.Node, expected filled array, got empty array")
		}

		t.Logf("Values retrieved: %v", values)
	})
}
