package pokeapi

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
)

type Config struct {
	NextMap *string
	PrevMap *string
}

type PokePage struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func CallApi(uri string) ([]byte, error) {
	res, err := http.Get(uri)
	if err != nil {
		var zero []byte
		return zero, err
	}

	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}

	if err != nil {
		log.Fatalf("%s", err)
	}
	return body, nil
}

func ParseToMap(body []byte, conf *Config) (string, error) {
	page := PokePage{}
	json.Unmarshal(body, &page)

	outputs := make([]string, len(page.Results))
	for i, place := range page.Results {
		outputs[i] = place.Name
	}
	conf.NextMap = page.Next
	conf.PrevMap = page.Previous

	return strings.Join(outputs, "\n"), nil
}
