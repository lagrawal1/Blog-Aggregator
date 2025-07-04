package main

import (
	"fmt"
	"gator/internal/config"
	"os"
)

type State struct {
	cfg *config.Config
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

	err := s.cfg.SetUser(cmd.arguments[0])

	if err != nil {
		return err
	}

	fmt.Println("User has been set to", cmd.arguments[0])
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
