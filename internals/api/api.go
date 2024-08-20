package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/obanoff/pokedexcli/internals/cache"
)

var apiCache *cache.Cache = cache.NewCache(20 * time.Second)

type Location struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Locations struct {
	Locations []Location `json:"results"`
	Next      string     `json:"next"`
	Prev      string     `json:"previous"`
}

func GetLocations(url string) (*Locations, error) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	if len(url) == 0 {
		url = "https://pokeapi.co/api/v2/location-area/?limit=20"
	}

	var locations *Locations

	data, ok := apiCache.Get(url)
	if !ok {
		resp, err := client.Get(url)
		if err != nil {
			return nil, fmt.Errorf("error making request: %w", err)
		}
		defer resp.Body.Close()

		data, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading data from response: %w", err)
		}

	}

	// overwrite cache after the visit data to prolong existence
	apiCache.Add(url, data)

	err := json.Unmarshal(data, &locations)
	if err != nil {
		return nil, fmt.Errorf("error decoding json: %w", err)
	}

	return locations, nil

}

type Pokemons struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

func GetPokemonsByLocation(location string) (*Pokemons, error) {
	if len(location) == 0 {
		return nil, fmt.Errorf("location area not provided")
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", location)

	data, ok := apiCache.Get(url)
	if !ok {
		resp, err := client.Get(url)
		if err != nil {
			return nil, fmt.Errorf("error making request: %w", err)
		}
		defer resp.Body.Close()

		data, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading data from response: %w", err)
		}
	}

	apiCache.Add(url, data)

	var pokemons *Pokemons

	err := json.Unmarshal(data, &pokemons)
	if err != nil {
		return nil, fmt.Errorf("error decoding json: %w", err)
	}

	return pokemons, nil
}

type Pokemon struct {
	Name           string `json:"name"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	BaseExperience int    `json:"base_experience"`
	Stats          []struct {
		Name string `json:"name"`
	} `json:"stats"`
	Types []struct {
		Name string `json:"name"`
	} `json:"types"`
}

func GetPokemonByName(name string) (*Pokemon, error) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", name)

	data, ok := apiCache.Get(url)
	if !ok {
		resp, err := client.Get(url)
		if err != nil {
			return nil, fmt.Errorf("error making request: %w", err)
		}
		defer resp.Body.Close()

		data, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading data from response: %w", err)
		}
	}

	apiCache.Add(url, data)

	var pokemon *Pokemon

	err := json.Unmarshal(data, &pokemon)
	if err != nil {
		return nil, fmt.Errorf("error decoding json: %w", err)
	}

	return pokemon, nil
}
