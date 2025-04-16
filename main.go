package main

import (
	"fmt"
	"os"
	"rss-aggregator/internal/config"
)

func main() {
	commands := config.Commands{
		Map: make(map[string]func(*config.State, config.CommandInput) error),
	}
	commands.Register("login", config.HandlerLogin)

	state := &config.State{
		Config: config.Read(),
	}
	args := config.CleanArgs(os.Args)
	_, ok := commands.Map[args[0]]
	if !ok {
		fmt.Println("Command does not exist")
		os.Exit(1)
	}

	commInput := config.CommandInput{
		Name: args[0],
		Args: args,
	}
	err := commands.Run(state, commInput)
	if err != nil {
		fmt.Printf("Error when running command:\t %s\n", err)
	}
}
