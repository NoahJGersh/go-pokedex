package main

import (
	"fmt"
	"os"
	pokeutils "poke-utils"
)

type config struct {
	next     string
	previous string
}

type cliCommand struct {
	name          string
	description   string
	callback      func(*config) error
	defaultConfig config
}

var commands map[string]cliCommand

func initCommands() {
	commands = map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Displays the names of locations in the world",
			callback:    commandMap,
			defaultConfig: config{
				next:     "https://pokeapi.co/api/v2/location-area/",
				previous: "",
			},
		},
	}
}

func commandExit(*config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)

	return nil
}

func commandHelp(*config) error {
	fmt.Println("Welcome to the Pokedex!\nUsage:")
	fmt.Println()

	for _, command := range commands {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}

	return nil
}

func commandMap(cfg *config) error {
	areas, next, _, err := pokeutils.GetLocationAreas(cfg.next)
	if err != nil {
		return err
	}

	cfg.previous = cfg.next
	cfg.next = next

	for _, area := range areas {
		fmt.Println(area)
	}

	return nil
}
