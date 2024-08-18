package models

import (
	"errors"
	"fmt"
	"os"

	"github.com/obanoff/pokedexcli/internals/api"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

type params struct {
	locations *api.Locations
}

type CommandRegistry struct {
	commands map[string]cliCommand
	params   params
}

func NewCommandRegistry() CommandRegistry {
	cmdr := CommandRegistry{
		commands: make(map[string]cliCommand),
		params: params{
			locations: nil,
		},
	}

	// help command
	cmdr.addCommand("help", "prints a help message describing how to use the REPL", func() error {
		fmt.Print("Usege: \n\nhelp: Displays a help message\nexit: Exit the Pokedex\n\n")
		return nil
	})

	// exit command
	cmdr.addCommand("exit", "exits the program", func() error {
		os.Exit(0)
		return nil
	})

	// map command
	cmdr.addCommand("map", "displays the names of 20 location in the Pokemon world; each subsequent call displays the next 20", func() error {
		var err error

		if cmdr.params.locations == nil {
			cmdr.params.locations, err = api.GetLocations("")
			if err != nil {
				return err
			}
		} else {
			if cmdr.params.locations.Next == "" {
				return errors.New("end of list: no locations found")
			}

			cmdr.params.locations, err = api.GetLocations(cmdr.params.locations.Next)
			if err != nil {
				return err
			}
		}

		for _, l := range cmdr.params.locations.Locations {
			fmt.Println(l.Name)
		}

		return nil
	})

	// mapb command
	cmdr.addCommand("mapb", "displays the names of 20 previous locations in the Pokemon world", func() error {
		var err error

		if cmdr.params.locations == nil || cmdr.params.locations.Prev == "" {
			return errors.New("no previous locations")
		}

		cmdr.params.locations, err = api.GetLocations(cmdr.params.locations.Prev)
		if err != nil {
			return err
		}

		for _, l := range cmdr.params.locations.Locations {
			fmt.Println(l.Name)
		}

		return nil
	})

	return cmdr
}

func (cd *CommandRegistry) addCommand(name, description string, callback func() error) {
	cd.commands[name] = cliCommand{
		name:        name,
		description: description,
		callback:    callback,
	}
}

func (cd *CommandRegistry) Run(name string) error {
	if _, ok := cd.commands[name]; !ok {
		return errors.New("\ncommand not found\n")
	}
	return cd.commands[name].callback()
}
