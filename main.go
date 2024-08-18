package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/obanoff/pokedexcli/internals/config"
)

func main() {
	app := config.NewAppConfig()
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("pokedex > ")

		scanner.Scan()
		err := scanner.Err()
		if err != nil {
			log.Fatal(err)
		}

		cmd := scanner.Text()

		err = app.Commands.Run(cmd)
		if err != nil {
			fmt.Println(err)
		}

	}
}
