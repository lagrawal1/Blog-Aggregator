package main

import (
	"context"
	"database/sql"
	"fmt"
	"gator/internal/config"
	"gator/internal/database"
	"os"
	"time"

	"github.com/google/uuid"
)

type State struct {
	cfg *config.Config
	db  *database.Queries
}

type Command struct {
	name      string
	arguments []string
}

type Commands struct {
	CommandsMap map[string]func(*State, Command) error
}

func handlerLogin(s *State, cmd Command) error {
	if len(cmd.arguments) == 0 {
		fmt.Println("you need to enter a username")
		os.Exit(1)
	}

	var name sql.NullString
	name.Scan(cmd.arguments[0])

	if _, err := s.db.GetUser(context.Background(), name); err != nil {
		fmt.Println("Can't login as nonexistent user.")
		print(err)
		os.Exit(1)
	}

	err := s.cfg.SetUser(cmd.arguments[0])

	if err != nil {
		return err
	}

	fmt.Println("User has been set to", cmd.arguments[0])
	return nil

}

func handlerRegister(s *State, cmd Command) error {
	if len(cmd.arguments) < 1 {
		fmt.Println("A name argument must be passed")
		os.Exit(1)
	}

	var id uuid.UUID

	id = uuid.New()

	var CreatedAt sql.NullTime
	CreatedAt.Scan(time.Now())

	var UpdatedAt sql.NullTime
	UpdatedAt.Scan(time.Now())

	var name sql.NullString
	name.Scan(cmd.arguments[0])

	_, err := s.db.GetUser(context.Background(), name)

	if err == nil {
		fmt.Println("User already exists")
		os.Exit(1)
	}

	name.Scan(cmd.arguments[0])
	user, err := s.db.CreateUser(context.Background(),
		database.CreateUserParams{
			ID:        id,
			CreatedAt: CreatedAt,
			UpdatedAt: UpdatedAt,
			Name:      name,
		})

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	s.cfg.SetUser(user.Name.String)
	fmt.Println("user was created", user)

	return nil

}

func handlerReset(s *State, cmd Command) error {
	err := s.db.DropUsers(context.Background())

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = s.db.CreateUsers(context.Background())

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return nil
}

func handlerUsers(s *State, cmd Command) error {
	users_list, err := s.db.GetUsers(context.Background())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, val := range users_list {
		fmt.Print(val.String)
		if val.String == s.cfg.CurrentUserName {
			fmt.Print(" (current)")
		}
		fmt.Print("\n")
	}
	return nil

}

func handlerAgg(s *State, cmd Command) error {
	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Print(feed)
	return nil
}

func (c *Commands) run(s *State, cmd Command) error {
	function, ok := c.CommandsMap[cmd.name]
	if !ok {
		return fmt.Errorf("dictionary access didn't work")

	}
	return function(s, cmd)
}

func (c *Commands) register(name string, f func(*State, Command) error) {
	c.CommandsMap[name] = f
}
