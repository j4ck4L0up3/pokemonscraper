package pokescraper

import (
	"bytes"
	"fmt"
	"text/template"
)

type Pokemon struct {
	ID   int
	Name string
	Type string
}

func SerializePokemon(poke Pokemon) (string, error) {
	var pokeBuffer bytes.Buffer
	tmpl := template.New("json")

	pokeTemplate := `{
		"id": {{ .ID }},
		"name": "{{ .Name }}",
		"type": "{{ .Type }}"
	}`

	_, parseErr := tmpl.Parse(pokeTemplate)
	if parseErr != nil {
		return "", fmt.Errorf("Template could not be parsed: %v\n", parseErr)
	}

	executeErr := tmpl.Execute(&pokeBuffer, poke)
	if executeErr != nil {
		return "", fmt.Errorf("Data could not be serialized: %v\n", executeErr)
	}

	return pokeBuffer.String(), nil
}
