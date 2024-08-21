package models

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/obanoff/pokedexcli/internals/api"
)

type cliCommand struct {
	name        string
	description string
	callback    func(string) error
}

type data struct {
	locations *api.Locations
	pokemons  map[string]*api.Pokemon
}

type CommandRegistry struct {
	commands map[string]cliCommand
	data     data
	rand     *rand.Rand
}

func NewCommandRegistry() CommandRegistry {
	cmdr := CommandRegistry{
		commands: make(map[string]cliCommand),
		data: data{
			locations: nil,
			pokemons:  make(map[string]*api.Pokemon),
		},
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	// help command
	cmdr.addCommand("help", "prints a help message describing how to use the REPL", func(param string) error {
		fmt.Print("Usege: \n\nhelp: Displays a help message\nexit: Exit the Pokedex\n\n")
		return nil
	})

	// exit command
	cmdr.addCommand("exit", "exits the program", func(param string) error {
		os.Exit(0)
		return nil
	})

	// map command
	cmdr.addCommand("map", "displays the names of 20 param in the Pokemon world; each subsequent call displays the next 20", func(location string) error {
		var err error

		if cmdr.data.locations == nil {
			cmdr.data.locations, err = api.GetLocations("")
			if err != nil {
				return err
			}
		} else {
			if cmdr.data.locations.Next == "" {
				return errors.New("end of list: no locations found")
			}

			cmdr.data.locations, err = api.GetLocations(cmdr.data.locations.Next)
			if err != nil {
				return err
			}
		}

		for _, l := range cmdr.data.locations.Locations {
			fmt.Println(l.Name)
		}

		return nil
	})

	// mapb command
	cmdr.addCommand("mapb", "displays the names of 20 previous locations in the Pokemon world", func(param string) error {
		var err error

		if cmdr.data.locations == nil || cmdr.data.locations.Prev == "" {
			return errors.New("no previous locations")
		}

		cmdr.data.locations, err = api.GetLocations(cmdr.data.locations.Prev)
		if err != nil {
			return err
		}

		for _, l := range cmdr.data.locations.Locations {
			fmt.Println(l.Name)
		}

		return nil
	})

	// explore command requires location area as its parameter
	cmdr.addCommand("explore", "displays pokemons in a given area", func(param string) error {
		var err error

		if len(param) == 0 {
			return errors.New("location are not provided")
		}

		result, err := api.GetPokemonsByLocation(param)
		if err != nil {
			return err
		}

		fmt.Printf("Exploring %s...\nFound Pokemon:\n", param)

		for _, pe := range result.PokemonEncounters {
			fmt.Printf(" - %s\n", pe.Pokemon.Name)
		}

		return nil
	})

	// catch command requires name as its parameter
	cmdr.addCommand("catch", "tries to catch a pokemon by its name", func(param string) error {
		var err error

		if len(param) == 0 {
			return errors.New("pokemon name not provided")
		}

		result, err := api.GetPokemonByName(param)
		if err != nil {
			return err
		}

		// catch pokemon logic
		baseExp := result.BaseExperience
		if baseExp <= 0 {
			baseExp = 1
		}
		catchChance := 1 / (float64(baseExp) / float64(50))
		proc := cmdr.rand.Float64()

		fmt.Printf("Throwing a Pokeball at %s...\n", param)

		if proc <= catchChance {
			cmdr.data.pokemons[param] = result
			fmt.Printf("%s was caught!\n", param)
		} else {
			fmt.Printf("%s escaped!\n", param)
		}

		return nil
	})

	// inspect caught pokemon by name
	cmdr.addCommand("inspect", "inspect caught pokemon by name", func(param string) error {
		if len(param) == 0 {
			return errors.New("pokemon not provided")
		}

		pokemon, ok := cmdr.data.pokemons[param]
		if !ok {
			return errors.New("you have not caught that pokemon")
		}

		fmt.Printf(`Name: %s
Height: %v
Weight: %v
`, pokemon.Name, pokemon.Height, pokemon.Weight)

		fmt.Println("Stats:")
		for _, s := range pokemon.Stats {
			fmt.Printf("  -%s: %v\n", s.Stat.Name, s.Value)
		}

		fmt.Println("Types:")
		for _, t := range pokemon.Types {
			fmt.Printf("  -%s\n", t.Type.Name)
		}

		return nil
	})

	return cmdr
}

func (cd *CommandRegistry) addCommand(name, description string, callback func(param string) error) {
	cd.commands[name] = cliCommand{
		name:        name,
		description: description,
		callback:    callback,
	}
}

func (cd *CommandRegistry) Run(name, param string) error {
	name, param = strings.TrimSpace(name), strings.TrimSpace(param)

	if _, ok := cd.commands[name]; !ok {
		return errors.New("command not found")
	}
	return cd.commands[name].callback(param)
}
