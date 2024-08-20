package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/obanoff/pokedexcli/internals/models"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	cmdRegistry := models.NewCommandRegistry()

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
			err = cmdRegistry.Run(parts[0], parts[1])
		default:
			err = cmdRegistry.Run(cmd, "")
		}
		if err != nil {
			fmt.Println(err)
		}

	}
}
