package pokescraper

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func FetchHTML(url string) (string, error) {
	resp, reqErr := http.Get(url)
	if reqErr != nil {
		return "", fmt.Errorf("unable to retrieve HTML response: %v", reqErr)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected response status: %v", resp.StatusCode)
	}

	htmlBuilder := new(strings.Builder)
	_, ioErr := io.Copy(htmlBuilder, resp.Body)
	if ioErr != nil {
		return "", fmt.Errorf("unexpected I/O copy error: %v", ioErr)
	}

	return htmlBuilder.String(), nil
}
