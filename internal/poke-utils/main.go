package pokeutils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

var pokeapiCache = NewCache(5 * time.Second)

type namedAPIResource struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type pokemonEncounter struct {
	Pokemon namedAPIResource `json:"pokemon"`
}

type pokemonStat struct {
	Stat      namedAPIResource `json:"stat"`
	Base_Stat int              `json:"base_stat"`
}

type pokemonType struct {
	Slot int              `json:"slot"`
	Type namedAPIResource `json:"type"`
}

type Pokemon struct {
	Name            string        `json:"name"`
	Base_Experience int           `json:"base_experience"`
	Height          int           `json:"height"`
	Weight          int           `json:"weight"`
	Stats           []pokemonStat `json:"stats"`
	Types           []pokemonType `json:"types"`
}

type locationAreasBody struct {
	Next     string             `json:"next"`
	Previous string             `json:"previous"`
	Results  []namedAPIResource `json:"results"`
}

type locationArea struct {
	Pokemon_Encounters []pokemonEncounter `json:"pokemon_encounters"`
}

func GetLocationAreas(url string) (locationAreas []string, next, previous string, err error) {
	body, ok := pokeapiCache.Get(url)

	if !ok {
		res, err := http.Get(url)
		if err != nil {
			return []string{}, "", "", fmt.Errorf("bad response from API: %w", err)
		}
		defer res.Body.Close()

		body, err = io.ReadAll(res.Body)
		if err != nil {
			return []string{}, "", "", fmt.Errorf("could not read response body: %w", err)
		}

		if res.StatusCode > 299 {
			return []string{}, "", "", fmt.Errorf("response failed with status code: %d and\nbody: %s", res.StatusCode, body)
		}

		pokeapiCache.Add(url, body)
	}

	areas := locationAreasBody{}
	err = json.Unmarshal(body, &areas)
	if err != nil {
		return []string{}, "", "", fmt.Errorf("could not unmarshal body: %w", err)
	}

	var areaNames []string
	for _, area := range areas.Results {
		areaNames = append(areaNames, area.Name)
	}

	return areaNames, areas.Next, areas.Previous, nil
}

func GetLocationArea(name string) (locationArea, error) {
	url := "https://pokeapi.co/api/v2/location-area/" + name

	body, ok := pokeapiCache.Get(url)
	if !ok {
		res, err := http.Get(url)
		if err != nil {
			return locationArea{}, fmt.Errorf("bad response from API: %w", err)
		}
		defer res.Body.Close()

		body, err = io.ReadAll(res.Body)
		if err != nil {
			return locationArea{}, fmt.Errorf("could not read response body: %w", err)
		}

		if res.StatusCode > 299 {
			return locationArea{}, fmt.Errorf("response failed with status code: %d and\nbody: %s", res.StatusCode, body)
		}

		pokeapiCache.Add(url, body)
	}

	area := locationArea{}
	err := json.Unmarshal(body, &area)
	if err != nil {
		return locationArea{}, fmt.Errorf("could not unmarshal body: %w", err)
	}

	return area, nil
}

func GetPokemon(name string) (Pokemon, error) {
	url := "https://pokeapi.co/api/v2/pokemon/" + name

	body, ok := pokeapiCache.Get(url)
	if !ok {
		res, err := http.Get(url)
		if err != nil {
			return Pokemon{}, fmt.Errorf("bad response from server: %w", err)
		}
		defer res.Body.Close()

		body, err = io.ReadAll(res.Body)
		if err != nil {
			return Pokemon{}, fmt.Errorf("could not read body: %w", err)
		}

		if res.StatusCode > 299 {
			return Pokemon{}, fmt.Errorf("response failed with status code: %d and\nbody: %s", res.StatusCode, body)
		}

		pokeapiCache.Add(url, body)
	}

	pokemon := Pokemon{}
	err := json.Unmarshal(body, &pokemon)
	if err != nil {
		return Pokemon{}, fmt.Errorf("could not unmarshal body: %w", err)
	}

	return pokemon, nil
}
