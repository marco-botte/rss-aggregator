package main

import (
	"database/sql"
	"fmt"
	"os"
	"rss-aggregator/internal/config"
	"rss-aggregator/internal/database"

	_ "github.com/lib/pq"
)

func main() {
	commands := config.Commands{
		Map: make(map[string]func(*config.State, config.CommandInput) error),
	}
	commands.Register("login", config.HandlerLogin)
	commands.Register("reset", config.HandlerReset)
	commands.Register("register", config.HandlerRegister)
	commands.Register("users", config.HandlerListUsers)
	commands.Register("agg", config.HandlerAgg)
	conf := config.Read()
	db, err := sql.Open("postgres", conf.DBurl)
	if err != nil {
		fmt.Printf("Error when opening db:\t %s\n", err)
	}
	dbQueries := database.New(db)
	state := &config.State{
		Db:     dbQueries,
		Config: conf,
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
	err = commands.Run(state, commInput)
	if err != nil {
		fmt.Printf("Error when running command:\t %s\n", err)
	}
}
