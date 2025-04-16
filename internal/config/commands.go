package config

import (
	"errors"
	"fmt"
	"os"
)

type CommandInput struct {
	Name string
	Args []string
}
type State struct {
	Config *Config
}

type Commands struct {
	Map map[string]func(*State, CommandInput) error
}

func CleanArgs(args []string) []string {
	if len(args) <= 1 { // first is the program name
		fmt.Println("Not enough arguments")
		os.Exit(1)
	}
	return args[1:]
}

func HandlerLogin(s *State, cmd CommandInput) error {
	if len(cmd.Args) == 1 {
		fmt.Println("Username is required")
		os.Exit(1)
	}
	user := cmd.Args[1]
	s.Config.SetUser(user)
	fmt.Printf("User has been set to: %s\n", user)
	return nil
}

func (c *Commands) Register(name string, f func(*State, CommandInput) error) {
	c.Map[name] = f
}

func (c *Commands) Run(s *State, cmdInput CommandInput) error {
	cmd, ok := c.Map[cmdInput.Name]
	if !ok {
		return errors.New("command not found")
	}
	err := cmd(s, cmdInput)
	if err != nil {
		return err
	}
	return nil
}
