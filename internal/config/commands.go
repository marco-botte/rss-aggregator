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

func MiddlewareLoggedIn(handler func(s *State, cmd CommandInput, user database.User) error) func(*State, CommandInput) error {
	return func(s *State, cmd CommandInput) error {
		user, err := s.Db.GetUser(context.Background(), s.Config.Username)
		if err != nil {
			fmt.Printf("Error. User may not exist. %s\n", err)
			os.Exit(1)
		}
		return handler(s, cmd, user)
	}
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
		fmt.Println("User name is required")
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
	err_1 := s.Db.DeleteUsers(context.Background())
	if err_1 != nil {
		fmt.Printf("Error while resetting users database. %s\n", err_1)
	}
	err_2 := s.Db.DeleteOrphanedFeeds(context.Background())
	if err_2 != nil {
		fmt.Printf("Error while deleting orphaned feeds. %s\n", err_2)
	}
	if err_1 != nil || err_2 != nil {
		os.Exit(1)
	}
	fmt.Println("Database successfully resetted")
	return nil
}

func HandlerListUsers(s *State, cmd CommandInput) error {
	users, err := s.Db.GetUsers(context.Background())
	if err != nil {
		fmt.Printf("Could not retrieve users. %s\n", err)
		os.Exit(1)
	}
	for _, user := range users {
		user_msg := fmt.Sprintf("* %v", user.Name)
		if user.Name == s.Config.Username {
			user_msg = user_msg + " (current)"
		}
		fmt.Println(user_msg)
	}
	return nil
}

func HandlerListFeeds(s *State, cmd CommandInput) error {
	feeds, err := s.Db.GetFeeds(context.Background())
	if err != nil {
		fmt.Printf("Could not retrieve feeds. %s\n", err)
		os.Exit(1)
	}

	for _, feed := range feeds {
		fmt.Printf("* %s,%s,%s\n", feed.Name, feed.Url, feed.Username)
	}
	return nil
}

func HandlerAgg(s *State, cmd CommandInput) error {
	feed, err := FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		fmt.Printf("Error while fetching feed. %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("%v\n", feed)
	return nil
}

func HandlerAddFeed(s *State, cmd CommandInput, user database.User) error {
	if len(cmd.Args) <= 2 {
		fmt.Println("Feed name and url are required")
		os.Exit(1)
	}
	feed, err := s.Db.CreateFeed(context.Background(), feedParams(cmd.Args[1], cmd.Args[2], user.ID))
	if err != nil {
		fmt.Printf("Error %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("Feed has been added:\n feed:\t%v\n", feed)
	feed_follow, err := s.Db.CreateFeedFollow(context.Background(), feedFollowParams(user.ID, feed.ID))
	if err != nil {
		fmt.Printf("Error %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("Feed has been created:\n feed_follow:\t%v\n", feed_follow)
	return nil
}

func HandlerFollow(s *State, cmd CommandInput, user database.User) error {
	if len(cmd.Args) == 1 {
		fmt.Println("Feed name is required")
		os.Exit(1)
	}
	feed, err := s.Db.GetFeed(context.Background(), cmd.Args[1])
	if err != nil {
		fmt.Printf("Error. Feed may not exist. %s\n", err)
		os.Exit(1)
	}
	feed_follow, err := s.Db.CreateFeedFollow(context.Background(), feedFollowParams(user.ID, feed.ID))
	if err != nil {
		fmt.Printf("Error %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("Feed %s has been followed by %s, follow_id: %s\n", feed.Name, user.Name, feed_follow.ID)
	return nil
}

func HandlerUnfollow(s *State, cmd CommandInput, user database.User) error {
	if len(cmd.Args) == 1 {
		fmt.Println("Feed name is required")
		os.Exit(1)
	}
	user, err := s.Db.GetUser(context.Background(), s.Config.Username)
	if err != nil {
		fmt.Printf("Error. User may not exist. %s\n", err)
		os.Exit(1)
	}
	feed, err := s.Db.GetFeed(context.Background(), cmd.Args[1])
	if err != nil {
		fmt.Printf("Error. Feed may not exist. %s\n", err)
		os.Exit(1)
	}
	params := database.RemoveFeedFollowParams{ID: user.ID, Url: feed.Url}
	err = s.Db.RemoveFeedFollow(context.Background(), params)
	if err != nil {
		fmt.Printf("Error removing feed %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("Feed %s has been unfollowed by %s\n", feed.Name, user.Name)
	return nil
}
func HandlerFollowing(s *State, cmd CommandInput, user database.User) error {
	feed_follows, err := s.Db.GetFeedFollowsForUser(context.Background(), user.Name)
	if err != nil {
		fmt.Printf("Error %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("%s follows:\n", user.Name)
	for _, feedFollow := range feed_follows {
		fmt.Printf("* %s\n", feedFollow.FeedName)
	}
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

func feedParams(name string, url string, userID uuid.UUID) database.CreateFeedParams {
	now := time.Now()
	return database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Name:      name,
		Url:       url,
		UserID:    userID,
	}
}

func feedFollowParams(userID uuid.UUID, feedID uuid.UUID) database.CreateFeedFollowParams {
	now := time.Now()
	return database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    userID,
		FeedID:    feedID,
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
