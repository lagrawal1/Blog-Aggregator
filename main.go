package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	config "gator/internal/config"
	"gator/internal/database"
	"os"
)

func main() {

	conf := config.Read()
	db, err := sql.Open("postgres", conf.DBUrl)
	if err != nil {
		fmt.Println("unable to open postgres")
		os.Exit(1)
	}
	dbQueries := database.New(db)

	curr_state := State{&conf, dbQueries}
	commands := Commands{
		CommandsMap: map[string]func(*State, Command) error{
			"login":    handlerLogin,
			"register": handlerRegister,
		},
	}

	commands.register("reset", handlerReset)
	commands.register("users", handlerUsers)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	arg_list := os.Args

	if len(arg_list) < 2 {
		fmt.Println("too few arguments")
		os.Exit(1)
	}

	comm := Command{name: arg_list[1], arguments: arg_list[2:]}

	err = commands.run(&curr_state, comm)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
