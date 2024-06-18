package pokeapi

import (
	"fmt"
	"strings"
)

func PokeEncounters(location Encounters) (string, error) {
	poke_offset := 2
	outputs := make([]string, len(location.PokemonEncounters)+poke_offset)
	outputs[0] = fmt.Sprintf("Exploring %s...", location.Location.Name)
	outputs[1] = "Found Pokemon:"
	for i, encounter := range location.PokemonEncounters {
		outputs[i+poke_offset] = fmt.Sprintf(" - %s", encounter.Pokemon.Name)
	}
	return strings.Join(outputs, "\n"), nil
}
