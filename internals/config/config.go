package config

import (
	"github.com/obanoff/pokedexcli/internals/models"
)

type AppConfig struct {
	Commands models.CommandRegistry
}

func NewAppConfig() *AppConfig {

	cmdsr := models.NewCommandRegistry()

	return &AppConfig{
		Commands: cmdsr,
	}
}
