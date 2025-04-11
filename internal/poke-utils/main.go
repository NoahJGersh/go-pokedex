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

type locationAreasBody struct {
	Count    int                `json:"count"`
	Next     string             `json:"next"`
	Previous string             `json:"previous"`
	Results  []namedAPIResource `json:"results"`
}

func GetLocationAreas(url string) (locationAreas []string, next, previous string, err error) {
	body, ok := pokeapiCache.Get(url)

	if !ok {
		res, err := http.Get(url)
		if err != nil {
			return []string{}, "", "", fmt.Errorf("no response from API: %w", err)
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
