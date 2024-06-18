package main

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/deemooneill/pokedexcli/internal/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*pokeapi.Config, *pokeapi.Cache, ...string) (string, error)
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	command_map := commands()
	conf := pokeapi.Config{PokeDex: map[string]pokeapi.Pokemon{}}
	cache := pokeapi.NewCache(5 * time.Minute)
	for {

		print("pokedex > ")
		for scanner.Scan() {
			text := strings.Fields(scanner.Text())
			command := text[0]
			rest := text[1:]
			result, ok := command_map[strings.ToLower(command)]
			if !ok {
				println("Command not recognised, type \"help\" to see the list of available commands")
				break
			}
			callback_text, err := result.callback(&conf, &cache, rest...)
			if err != nil {
				fmt.Println(err)
				return
			}
			println(callback_text)
			break
		}

	}
}

func commandHelp(conf *pokeapi.Config, cache *pokeapi.Cache, _ ...string) (string, error) {
	var help_texts []string
	help_texts = append(help_texts, "\n\n")
	help_texts = append(help_texts, "Welcome to the Pokedex!\nUsage:\n\n")
	for _, command := range commands() {
		help_texts = append(help_texts, fmt.Sprintf("%v: %s\n", command.name, command.description))
	}
	return strings.Join(help_texts, ""), nil
}

func commandExit(conf *pokeapi.Config, cache *pokeapi.Cache, _ ...string) (string, error) {
	return "", errors.New("exiting pokedex, goodbye")
}

func commandBMap(conf *pokeapi.Config, cache *pokeapi.Cache, _ ...string) (string, error) {
	var uri string
	if conf.PrevMap != nil {
		uri = *conf.PrevMap
	} else {
		return "No previous map", nil
	}

	location := pokeapi.Locations{}
	err := pokeapi.CallApi(uri, cache, &location)
	if err != nil {
		return "", err
	}
	return pokeapi.ParseToLocations(location, conf)
}

func commandMap(conf *pokeapi.Config, cache *pokeapi.Cache, _ ...string) (string, error) {
	var uri string
	if conf.NextMap != nil {
		uri = *conf.NextMap
	} else {
		uri = "https://pokeapi.co/api/v2/location-area/"
	}
	location := pokeapi.Locations{}
	err := pokeapi.CallApi(uri, cache, &location)
	if err != nil {
		return "", err
	}
	return pokeapi.ParseToLocations(location, conf)

}

func commandExplore(conf *pokeapi.Config, cache *pokeapi.Cache, params ...string) (string, error) {

	uri := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", params[0])
	encounter := pokeapi.Encounters{}
	err := pokeapi.CallApi(uri, cache, &encounter)
	if err != nil {
		return "", err
	}
	return pokeapi.PokeEncounters(encounter)

}

func commandCatch(conf *pokeapi.Config, cache *pokeapi.Cache, params ...string) (string, error) {
	uri := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", params[0])
	p := pokeapi.PokemonFull{}
	err := pokeapi.CallApi(uri, cache, &p)
	if err != nil {
		return "", err
	}
	pokemon := pokeapi.Pokemon{
		Name:           p.Name,
		Stats:          p.Stats,
		Height:         p.Height,
		Weight:         p.Weight,
		BaseExperience: p.BaseExperience,
		Types:          p.Types,
	}
	println("Throwing a Pokeball at", p.Name, "...")
	if rand.Intn(pokemon.BaseExperience) > pokemon.BaseExperience/3 {
		conf.PokeDex[p.Name] = pokemon
		return "Caught pokemon " + p.Name, nil
	}
	return p.Name + " escaped!", nil
}

func commandInspect(conf *pokeapi.Config, cache *pokeapi.Cache, params ...string) (string, error) {
	name := params[0]
	data, ok := conf.PokeDex[name]
	if !ok {
		return "you have not caught that pokemon", nil
	}
	return data.String(), nil
}
func commandPokedex(conf *pokeapi.Config, cache *pokeapi.Cache, params ...string) (string, error) {
	pokemon := make([]string, 0, len(conf.PokeDex))
	pokemon = append(pokemon, "Your Pokedex:")
	for k, _ := range conf.PokeDex {
		pokemon = append(pokemon, " - "+k)
	}
	return strings.Join(pokemon, "\n"), nil
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
		"explore": {
			name:        "explore",
			description: "lets you explore an area's pokemon.\nusage: explore area_name",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "lets you try to catch a pokemon.\nusage: catch pikachu",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "lets you inspect a caught pokemon",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "lets you see all of the pokemon you've caught",
			callback:    commandPokedex,
		},
	}
}
