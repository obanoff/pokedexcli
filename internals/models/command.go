package models

import "errors"

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

type CommandRegistry map[string]cliCommand

func NewCommandRegistry() CommandRegistry {
	return make(map[string]cliCommand)
}

func (cd CommandRegistry) AddCommand(name, description string, callback func() error) {
	cd[name] = cliCommand{
		name:        name,
		description: description,
		callback:    callback,
	}
}

func (cd CommandRegistry) Run(name string) error {
	if _, ok := cd[name]; !ok {
		return errors.New("command not found")
	}
	return cd[name].callback()
}
