package pokescraper

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestFetchHTML(t *testing.T) {
	// Test case: Successful fetch with html
	t.Run("successful fetch", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			htmlFile, readErr := os.ReadFile("mock.html")
			if readErr != nil {
				t.Fatalf("unable to read mock html file for mock server: %v", readErr)
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
			t.Fatalf("error reading the mock html file for expected data: %v", readErr)
		}

		expectedHtml := bytes.NewReader(expectedData)
		_, ioErr := io.Copy(builder, expectedHtml)
		if ioErr != nil {
			t.Fatalf("error copying byte array to strings.Builder object: %v", ioErr)
		}

		expectedStr := builder.String()

		actualStr, fetchErr := FetchHTML(mockServer.URL)
		if fetchErr != nil {
			t.Fatalf("expected no error, got error: %v", fetchErr)
		}

		if !reflect.DeepEqual(expectedStr, actualStr) {
			t.Fatalf("actual string does not match expected string.\nactual string: %v\nexpected string: %v", actualStr, expectedStr)
		}

	})

}
