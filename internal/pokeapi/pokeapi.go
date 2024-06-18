package pokeapi

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
)

type Config struct {
	PokeDex map[string]Pokemon
	NextMap *string
	PrevMap *string
}

func CallApi[T any](uri string, cache *Cache, into *T) error {

	if res, ok := cache.Get(uri); ok {
		json.Unmarshal(res, &into)
		return nil
	}
	res, err := http.Get(uri)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		// log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
		return errors.New("not found")
	}

	if err != nil {
		// log.Fatalf("%s", err)
		return errors.New("not found")
	}
	cache.Add(uri, body)
	json.Unmarshal(body, &into)
	return nil
}

func ParseToLocations(page Locations, conf *Config) (string, error) {

	outputs := make([]string, len(page.Results))
	for i, place := range page.Results {
		outputs[i] = place.Name
	}
	conf.NextMap = page.Next
	conf.PrevMap = page.Previous

	return strings.Join(outputs, "\n"), nil
}
