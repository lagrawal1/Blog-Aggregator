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

	id := uuid.New()

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
	err := s.db.DropTableFeeds(context.Background())

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = s.db.DropUsers(context.Background())

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = s.db.CreateUsers(context.Background())

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = s.db.CreateTableFeeds(context.Background())

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

func handlerAddFeed(s *State, cmd Command) error {

	if len(cmd.arguments) < 2 {
		fmt.Println("not enough arguments")
		os.Exit(1)
	}
	feedname_arg := cmd.arguments[0]
	url_arg := cmd.arguments[1]

	currentUser := s.cfg.CurrentUserName
	var username sql.NullString
	username.Scan(currentUser)

	user, err := s.db.GetUser(context.Background(), username)

	var user_fk uuid.NullUUID
	user_fk.Scan(user.ID.String())

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	id := uuid.New()

	var createdAt sql.NullTime
	createdAt.Scan(time.Now())

	var updatedAt sql.NullTime
	updatedAt.Scan(time.Now())

	var url sql.NullString
	url.Scan(url_arg)

	var feedname sql.NullString
	feedname.Scan(feedname_arg)

	_, err = s.db.CreateFeed(context.Background(), database.CreateFeedParams{

		ID:        id,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		Name:      feedname,
		Url:       url,
		UserID:    user_fk,
	},
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return nil
}

func handlerFeeds(s *State, cmd Command) error {
	feeds, err := s.db.SelectAllFeeds(context.Background())
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	for _, feed := range feeds {
		user, err := s.db.GetUserById(context.Background(), feed.UserID.UUID)

		if err != nil {
			fmt.Print(err)
			os.Exit(1)
		}
		fmt.Println(feed.Name.String)
		fmt.Println("- URL:", feed.Url.String)
		fmt.Println("- Username", user.Name.String, "\n")
	}

	return nil
}

func handlerFollow(s *State, cmd Command) error {
	url_arg := cmd.arguments[0]
	var name sql.NullString
	name.Scan(s.cfg.CurrentUserName)

	user, err := s.db.GetUser(context.Background(), name)

	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	follow_id := uuid.New()

	var created_at sql.NullTime
	var updated_at sql.NullTime

	created_at.Scan(time.Now())
	updated_at.Scan(time.Now())

	var user_id uuid.NullUUID
	user_id.Scan(user.ID.String())

	var feed_id uuid.NullUUID
	var url sql.NullString
	url.Scan(url_arg)
	feed, err := s.db.SelectFeedByURL(context.Background(), url)

	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	feed_id.Scan(feed.ID.String())

	s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        follow_id,
		CreatedAt: created_at,
		UpdatedAt: updated_at,
		UserID:    user_id,
		FeedID:    feed_id,
	})

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
