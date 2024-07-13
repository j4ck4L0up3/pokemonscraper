package pokescraper

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestSerializePokemon(t *testing.T) {
	tests := []struct {
		name         string
		expectedPoke *Pokemon
		actualPoke   *Pokemon
		wantErr      bool
	}{
		{
			name:         "Pikachu test",
			expectedPoke: &Pokemon{ID: "025", Name: "pikachu", Type: []string{"electic"}},
			actualPoke:   &Pokemon{},
			wantErr:      false,
		},
		{
			name:         "Charizard test",
			expectedPoke: &Pokemon{ID: "006", Name: "charizard", Type: []string{"fire", "flying"}},
			actualPoke:   &Pokemon{},
			wantErr:      false,
		},
		{
			name:         "empty Pokemon struct test",
			expectedPoke: &Pokemon{},
			actualPoke:   &Pokemon{},
			wantErr:      true,
		},
		{
			name:         "partial Pokemon struct test",
			expectedPoke: &Pokemon{ID: "001", Name: "bulbasaur"},
			actualPoke:   &Pokemon{},
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, serializeErr := SerializePokemon(*tt.expectedPoke)
			if (serializeErr != nil) != tt.wantErr {
				t.Fatalf("Error during serialization: %v\n", serializeErr)
			}

			emptyPoke := Pokemon{}
			if serializeErr != nil && tt.wantErr && got == "" &&
				(IsEmptyPokemon(emptyPoke, *tt.expectedPoke) ||
					IsPartialPokemon(*tt.expectedPoke)) {
				t.Log("Intentional serializeErr occurred\n")
				return
			}

			jsonErr := json.Unmarshal([]byte(got), tt.actualPoke)
			if jsonErr != nil {
				t.Fatalf("Error Unmarshaling serialized data.\n Actual Poke: %v\n Expected: %v\n Error: %v\n", *tt.actualPoke, *tt.expectedPoke, jsonErr)
			}

			if !reflect.DeepEqual(*tt.expectedPoke, *tt.actualPoke) {
				t.Fatalf("Actual Pokemon does not match expected.\nActual Pokemon: %v\nExpected Pokemon: %v\n", *tt.expectedPoke, *tt.actualPoke)
			}
		})
	}
}

func TestDeserializePokemon(t *testing.T) {
	tests := []struct {
		name       string
		inputStr   *string
		inputPoke  *Pokemon
		actualPoke *Pokemon
		wantErr    bool
	}{
		{
			name:       "Pikachu test",
			inputStr:   new(string),
			inputPoke:  &Pokemon{ID: "025", Name: "pikachu", Type: []string{"electic"}},
			actualPoke: &Pokemon{},
			wantErr:    false,
		},
		{
			name:       "Charizard test",
			inputStr:   new(string),
			inputPoke:  &Pokemon{ID: "006", Name: "charizard", Type: []string{"fire", "flying"}},
			actualPoke: &Pokemon{},
			wantErr:    false,
		},
		{
			name:       "empty Pokemon struct test",
			inputStr:   new(string),
			inputPoke:  &Pokemon{},
			actualPoke: &Pokemon{},
			wantErr:    true,
		},
		{
			name:       "partial Pokemon struct test",
			inputStr:   new(string),
			inputPoke:  &Pokemon{ID: "001", Name: "bulbasaur"},
			actualPoke: &Pokemon{},
			wantErr:    true,
		},
		// TODO: add test case for extra struct field(s)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expectedPokeJSON, serialErr := json.MarshalIndent(*tt.inputPoke, "", "    ")
			if serialErr != nil {
				t.Fatalf("Error Marshaling deserialized data.\nInput Poke: %v\nError: %v\n", *tt.inputPoke, serialErr)
			}

			inputStr := string(expectedPokeJSON)
			deserialErr := DeserializePokemon(inputStr, tt.actualPoke)
			if (deserialErr != nil) != tt.wantErr {
				t.Fatalf("Error during deserialization: %v\nActual Poke (after expected change): %v\n", deserialErr, *tt.actualPoke)
			}

			emptyPoke := Pokemon{}
			if deserialErr != nil && tt.wantErr && inputStr == "" &&
				(IsEmptyPokemon(emptyPoke, *tt.inputPoke) || IsPartialPokemon(*tt.inputPoke)) {
				t.Log("Intentional deserialErr has occurred.\n")
				return
			}

			if !reflect.DeepEqual(*tt.inputPoke, *tt.actualPoke) {
				t.Fatalf("Actual Pokemon does not match expected.\nActual Pokemon: %v\nExpected Pokemon: %v\n", *tt.inputPoke, *tt.actualPoke)
			}

			actualPokeJSON, reserialErr := json.MarshalIndent(*tt.actualPoke, "", "    ")
			if reserialErr != nil {
				t.Fatalf("Error Marshaling deserialized data.\nActual Poke: %v\nError: %v\n", *tt.actualPoke, reserialErr)
			}

			outputStr := string(actualPokeJSON)
			if outputStr != inputStr {
				t.Fatalf("Final string does not match original.\nInput String: %v\nOutput String: %v\n", inputStr, outputStr)
			}
		})
	}
}
