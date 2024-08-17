package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/obanoff/pokedexcli/internals/models"
)

func main() {
	cmdRegistry := models.NewCommandRegistry()

	cmdRegistry.AddCommand("help", "prints a help message describing how to use the REPL", func() error {
		fmt.Println("Usege: ")
		return nil
	})

	cmdRegistry.AddCommand("exit", "exits the program", func() error {
		os.Exit(0)
		return nil
	})

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("pokedex > ")

		scanner.Scan()
		err := scanner.Err()
		if err != nil {
			log.Fatal(err)
		}

		cmd := scanner.Text()

		err = cmdRegistry.Run(cmd)
		if err != nil {
			fmt.Println(err)
		}

	}
}
