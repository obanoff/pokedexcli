package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/obanoff/pokedexcli/internals/cache"
)

var locationCache *cache.Cache = cache.NewCache(10 * time.Second)

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

	data, ok := locationCache.Get(url)
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

		locationCache.Add(url, data)
	}

	err := json.Unmarshal(data, &locations)
	if err != nil {
		return nil, fmt.Errorf("error decoding json: %w", err)
	}

	return locations, nil

}
