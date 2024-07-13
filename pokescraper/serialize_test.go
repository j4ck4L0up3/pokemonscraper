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
				t.Fatalf("error during serialization: %v", serializeErr)
			}

			emptyPoke := Pokemon{}
			if serializeErr != nil && tt.wantErr && got == "" &&
				(IsEmptyPokemon(emptyPoke, *tt.expectedPoke) ||
					IsPartialPokemon(*tt.expectedPoke)) {
				t.Log("intentional serializeErr occurred")
				return
			}

			jsonErr := json.Unmarshal([]byte(got), tt.actualPoke)
			if jsonErr != nil {
				t.Fatalf("error Unmarshaling serialized data.\n actual Poke: %v\n expected: %v\n error: %v", *tt.actualPoke, *tt.expectedPoke, jsonErr)
			}

			if !reflect.DeepEqual(*tt.expectedPoke, *tt.actualPoke) {
				t.Fatalf("actual Pokemon does not match expected.\nactual Pokemon: %v\nexpected Pokemon: %v", *tt.expectedPoke, *tt.actualPoke)
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
			name:       "pikachu test",
			inputStr:   new(string),
			inputPoke:  &Pokemon{ID: "025", Name: "pikachu", Type: []string{"electic"}},
			actualPoke: &Pokemon{},
			wantErr:    false,
		},
		{
			name:       "charizard test",
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expectedPokeJSON, serialErr := json.MarshalIndent(*tt.inputPoke, "", "    ")
			if serialErr != nil {
				t.Fatalf("error Marshaling deserialized data.\ninput Poke: %v\nerror: %v", *tt.inputPoke, serialErr)
			}

			inputStr := string(expectedPokeJSON)
			deserialErr := DeserializePokemon(inputStr, tt.actualPoke)
			if (deserialErr != nil) != tt.wantErr {
				t.Fatalf("error during deserialization: %v\nactual Poke (after expected change): %v", deserialErr, *tt.actualPoke)
			}

			emptyPoke := Pokemon{}
			if deserialErr != nil && tt.wantErr && inputStr == "" &&
				(IsEmptyPokemon(emptyPoke, *tt.inputPoke) || IsPartialPokemon(*tt.inputPoke)) {
				t.Log("intentional deserialErr has occurred.")
				return
			}

			if !reflect.DeepEqual(*tt.inputPoke, *tt.actualPoke) {
				t.Fatalf("actual Pokemon does not match expected.\nactual Pokemon: %v\nexpected Pokemon: %v", *tt.inputPoke, *tt.actualPoke)
			}

			actualPokeJSON, reserialErr := json.MarshalIndent(*tt.actualPoke, "", "    ")
			if reserialErr != nil {
				t.Fatalf("error Marshaling deserialized data.\nactual Poke: %v\nerror: %v", *tt.actualPoke, reserialErr)
			}

			outputStr := string(actualPokeJSON)
			if outputStr != inputStr {
				t.Fatalf("final string does not match original.\ninput String: %v\noutput String: %v", inputStr, outputStr)
			}
		})
	}
}
