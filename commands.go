package main

import (
	"fmt"
	"math/rand"
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
		"catch": {
			name:        "catch [pokemon-name]",
			description: "Try to catch a pokemon!",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect [pokemon-name]",
			description: "Get info about a pokemon you've caught",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "See all the pokemon you've caught",
			callback:    commandPokedex,
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

func commandMap(cfg *config, _ ...string) error {
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

func commandMapb(cfg *config, _ ...string) error {
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

func commandCatch(_ *config, args ...string) error {
	if len(args) == 0 {
		fmt.Println("You need to specify a pokemon.")
		return nil
	}

	pokemon, err := pokeutils.GetPokemon(args[0])
	if err != nil {
		return err
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon.Name)

	if catchRate, catchRoll := 1.0/(float32(pokemon.Base_Experience)/20.0), rand.Float32(); catchRoll <= catchRate {
		fmt.Printf("%s was caught!\n", pokemon.Name)
		pokedex[pokemon.Name] = pokemon
	} else {
		fmt.Printf("%s escaped!\n", pokemon.Name)
	}

	return nil
}

func commandInspect(_ *config, args ...string) error {
	if len(args) == 0 {
		fmt.Println("You need to specify a pokemon.")
		return nil
	}

	pokemon, ok := pokedex[args[0]]
	if !ok {
		fmt.Println("you have not caught that pokemon")
		return nil
	}

	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf("  -%s: %d\n", stat.Stat.Name, stat.Base_Stat)
	}

	fmt.Println("Types:")
	for _, t := range pokemon.Types {
		fmt.Printf("  - %s\n", t.Type.Name)
	}

	return nil
}

func commandPokedex(*config, ...string) error {
	if len(pokedex) == 0 {
		fmt.Println("You haven't caught any Pokemon yet!")
		return nil
	}

	fmt.Println("Your Pokedex:")

	for k, _ := range pokedex {
		fmt.Printf(" - %s\n", k)
	}

	return nil
}
