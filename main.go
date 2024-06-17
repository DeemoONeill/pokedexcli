package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/deemooneill/pokedexcli/internal/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*pokeapi.Config) (string, error)
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	command_map := commands()
	conf := pokeapi.Config{}
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

func commandHelp(conf *pokeapi.Config) (string, error) {
	var help_texts []string
	help_texts = append(help_texts, "\n\n")
	help_texts = append(help_texts, "Welcome to the Pokedex!\nUsage:\n\n")
	for _, command := range commands() {
		help_texts = append(help_texts, fmt.Sprintf("%v: %s\n", command.name, command.description))
	}
	return strings.Join(help_texts, ""), nil
}

func commandExit(conf *pokeapi.Config) (string, error) {
	return "", errors.New("exiting pokedex, goodbye")
}

func commandBMap(conf *pokeapi.Config) (string, error) {
	var uri string
	if conf.PrevMap != nil {
		uri = *conf.PrevMap
	} else {
		return "No previous map", nil
	}
	body, err := pokeapi.CallApi(uri)
	if err != nil {
		return "", err
	}
	return pokeapi.ParseToMap(body, conf)
}

func commandMap(conf *pokeapi.Config) (string, error) {
	var uri string
	if conf.NextMap != nil {
		uri = *conf.NextMap
	} else {
		uri = "https://pokeapi.co/api/v2/location-area/"
	}
	body, err := pokeapi.CallApi(uri)
	if err != nil {
		return "", err
	}
	return pokeapi.ParseToMap(body, conf)

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
		"map": {
			name:        "map",
			description: "lets you run through the map",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "lets you go backwards through the map",
			callback:    commandBMap,
		},
	}
}
