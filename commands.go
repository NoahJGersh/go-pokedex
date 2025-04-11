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
	callback      func(*config, ...string) error
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
			description: "Displays the names of the next 20 locations in the world",
			callback:    commandMap,
			defaultConfig: config{
				next:     "https://pokeapi.co/api/v2/location-area/",
				previous: "",
			},
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the names of the previous 20 locations in the world",
			callback:    commandMapb,
			defaultConfig: config{
				next:     "",
				previous: "https://pokeapi.co/api/v2/location-area/",
			},
		},
		"explore": {
			name:        "explore [location-name]",
			description: "Lists the pokemon you can find at a given location",
			callback:    commandExplore,
		},
	}
}

func commandExit(*config, ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)

	return nil
}

func commandHelp(*config, ...string) error {
	fmt.Println("Welcome to the Pokedex!\nUsage:")
	fmt.Println()

	for _, command := range commands {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}

	return nil
}

func commandMap(cfg *config, args ...string) error {
	areas, next, previous, err := pokeutils.GetLocationAreas(cfg.next)
	if err != nil {
		return err
	}

	cfg.previous = previous
	cfg.next = next

	for _, area := range areas {
		fmt.Println(area)
	}

	return nil
}

func commandMapb(cfg *config, args ...string) error {
	if cfg.previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}

	areas, next, previous, err := pokeutils.GetLocationAreas(cfg.previous)
	if err != nil {
		return err
	}

	cfg.next = next
	cfg.previous = previous

	for _, area := range areas {
		fmt.Println(area)
	}

	return nil
}

func commandExplore(_ *config, args ...string) error {
	if len(args) == 0 {
		fmt.Println("You need to specify a location.")
		return nil
	}

	area, err := pokeutils.GetLocationArea(args[0])
	if err != nil {
		return err
	}

	if len(area.Pokemon_Encounters) == 0 {
		fmt.Println("No pokemon were found for that location.")
		return nil
	}

	fmt.Println("Found Pokemon:")
	for _, encounter := range area.Pokemon_Encounters {
		fmt.Println(" - ", encounter.Pokemon.Name)
	}

	return nil
}
