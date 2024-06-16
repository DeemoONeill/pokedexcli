package main

import (
	"fmt"
	"bufio"
	"os"
	"strings"
	"errors"
	"net/http"
	"io"
	"log"
	"encoding/json"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config) (string, error)
}

type config struct {
	nextMap *string
	prevMap *string
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	command_map := commands()
	conf := config{}
	for {

		print("pokedex > ")
		for scanner.Scan() {
			text := scanner.Text()
			result, ok := command_map[strings.ToLower(text)]
			if !ok {
				println("Command not recognised, type \"help\" to see the list of available commands")
				break
			}
			callback_text, err := result.callback(&conf)	
			if err != nil {
				fmt.Println(err)
				return
			}
			println(callback_text)
			break
		}

	}
}

func commandHelp(conf *config) (string, error) {
	var help_texts []string
	help_texts = append(help_texts, "\n\n")
	help_texts = append(help_texts, "Welcome to the Pokedex!\nUsage:\n\n")
	for _, command := range commands() {
		help_texts = append(help_texts, fmt.Sprintf("%v: %s\n", command.name, command.description))
	}
	return strings.Join(help_texts, ""), nil
}

func commandExit(conf *config) (string, error) {
	return "", errors.New("exiting pokedex, goodbye!")
}

func commandMap(conf *config) (string, error) {
	var uri string
	if conf.nextMap != nil {
		uri = *conf.nextMap
	}else{
		uri = "https://pokeapi.co/api/v2/location-area/"
	}
	res, err := http.Get(uri)
	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}

	if err != nil {
		log.Fatalf("%s", err)
	}
	page := PokePage{}
	json.Unmarshal(body, &page)
	
	outputs := make([]string, len(page.Results)) 
	for i, place := range page.Results {
		outputs[i] = fmt.Sprintf("%s", place.Name)
	}
	conf.nextMap = page.Next
	conf.prevMap = page.Previous

	return strings.Join(outputs, "\n"), nil


}
func commandBMap(conf *config) (string, error) {
	var uri string
	if conf.prevMap != nil {
		uri = *conf.prevMap
	}else{
		return "No previous map", nil
	}
	res, err := http.Get(uri)
	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}

	if err != nil {
		log.Fatalf("%s", err)
	}
	page := PokePage{}
	json.Unmarshal(body, &page)
	
	outputs := make([]string, len(page.Results)) 
	for i, place := range page.Results {
		outputs[i] = fmt.Sprintf("%s", place.Name)
	}
	conf.nextMap = page.Next
	conf.prevMap = page.Previous

	return strings.Join(outputs, "\n"), nil
}


func commands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "displays this help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exits the pokedex",
			callback:    commandExit,
		},
		"map" : {
			name: "map",
			description: "lets you run through the map",
			callback: commandMap,
		},
		"mapb": {
			name: "mapb",
			description: "lets you go backwards through the map",
			callback: commandBMap,
		},
	}
}

type PokePage struct {
	Count    int    `json:"count"`
	Next     *string `json:"next"`
	Previous *string    `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

