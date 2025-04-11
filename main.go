package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

func main() {
	initCommands()
	scanner := bufio.NewScanner(os.Stdin)

	var cfg *config

	for {
		fmt.Print("Pokedex > ")
		hasInput := scanner.Scan()
		if !hasInput {
			continue
		}

		input := scanner.Text()

		cleaned := cleanInput(input)
		if len(cleaned) == 0 {
			continue
		}

		cmd, ok := commands[cleaned[0]]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}

		if cfg == nil && cmd.defaultConfig != (config{}) {
			newConfig := cmd.defaultConfig
			cfg = &newConfig
		}

		err := cmd.callback(cfg)
		if err != nil {
			fmt.Println("Error:", err)
		}
	}
}
