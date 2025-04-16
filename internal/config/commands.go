package config

import (
	"context"
	"errors"
	"fmt"
	"os"
	"rss-aggregator/internal/database"
	"time"

	"github.com/google/uuid"
)

type CommandInput struct {
	Name string
	Args []string
}
type State struct {
	Db     *database.Queries
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
	username := cmd.Args[1]
	user, err := s.Db.GetUser(context.Background(), username)
	if err != nil {
		fmt.Println("User must exist to be logged in")
		os.Exit(1)
	}
	s.Config.SetUser(user.Name)
	return nil
}
func HandlerRegister(s *State, cmd CommandInput) error {
	if len(cmd.Args) == 1 {
		fmt.Println("Name is required")
		os.Exit(1)
	}
	user, err := s.Db.CreateUser(context.Background(), userParams(cmd.Args[1]))
	if err != nil {
		fmt.Printf("Error. User with name may already exist. %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("User has been registered:\n User:\t%v\n", user)
	s.Config.SetUser(user.Name)
	return nil
}
func HandlerReset(s *State, cmd CommandInput) error {
	err := s.Db.DeleteUsers(context.Background())
	if err != nil {
		fmt.Printf("Error while resetting database. %s\n", err)
		os.Exit(1)
	}
	fmt.Println("Database successfully resetted")
	return nil
}
func userParams(name string) database.CreateUserParams {
	now := time.Now()
	return database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Name:      name,
	}
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
