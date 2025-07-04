package main

import (
	"fmt"
	config "gator/internal/config"
	"os"
)

func main() {

	conf := config.Read()
	curr_state := State{&conf}
	commands := Commands{
		CommandsMap: map[string]func(*State, Command) error{
			"login": handlerLogin,
		},
	}

	arg_list := os.Args

	if len(arg_list) < 2 {
		fmt.Println("too few arguments")
		os.Exit(1)
	}

	comm := Command{name: arg_list[1], arguments: arg_list[2:]}

	err := commands.run(&curr_state, comm)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
