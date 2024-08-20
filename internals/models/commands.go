package models

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
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

		result, err := api.GetPokemonsByLocation(param)
		if err != nil {
			fmt.Println("area not found")
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

		result, err := api.GetPokemonByName(param)
		if err != nil {
			fmt.Println("pokemon not found")
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
	if _, ok := cd.commands[name]; !ok {
		return errors.New("\ncommand not found\n")
	}
	return cd.commands[name].callback(param)
}
