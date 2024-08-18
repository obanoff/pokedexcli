package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/obanoff/pokedexcli/internals/config"
)

func main() {
	app := config.NewAppConfig()
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")

		scanner.Scan()
		err := scanner.Err()
		if err != nil {
			log.Fatal(err)
		}

		cmd := scanner.Text()

		parts := strings.Fields(cmd)

		switch len(parts) {
		case 2:
			err = app.Commands.Run(parts[0], parts[1])
		default:
			err = app.Commands.Run(cmd, "")
		}
		if err != nil {
			fmt.Println(err)
		}

	}
}
